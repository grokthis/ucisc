package asm

import (
	"fmt"
	"strconv"
	"strings"
)

// ElementParser is an interface for parsing different statement types
// Uses ParseContext to follow Interface Segregation Principle
type ElementParser interface {
	// Matches returns true if this parser can handle the current token sequence
	// It should not consume tokens or add errors
	Matches(ctx ParseContext) bool

	// Parse parses the statement and returns it, or nil on error
	// It should consume tokens and may add errors
	Parse(ctx ParseContext) Statement
}

// =============================================================================
// Label Parser
// =============================================================================

// LabelParser parses label statements
type LabelParser struct{}

func (lp *LabelParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokIdent && ctx.PeekToken().Type == TokColon
}

func (lp *LabelParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	label := &LabelStmt{
		Name: curTok.Value,
		Line: curTok.Line,
		Col:  curTok.Column,
	}

	ctx.NextToken() // consume label name
	ctx.NextToken() // consume colon

	return label
}

// =============================================================================
// String Parser
// =============================================================================

// StringParser parses string literal statements
type StringParser struct{}

func (sp *StringParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokString
}

func (sp *StringParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	stringStmt := &StringStmt{
		Value: curTok.Value,
		Line:  curTok.Line,
		Col:   curTok.Column,
	}

	ctx.NextToken() // consume string

	return stringStmt
}

// =============================================================================
// Data Parser
// =============================================================================

// DataParser parses data block statements
type DataParser struct{}

func (dp *DataParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokPercent
}

func (dp *DataParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	dataStmt := &DataStmt{
		Line:   curTok.Line,
		Col:    curTok.Column,
		Values: []uint16{},
	}

	ctx.NextToken() // consume %

	// Parse hex numbers until newline or EOF
	for ctx.CurrentToken().Type != TokNewline && ctx.CurrentToken().Type != TokEOF {
		curTok = ctx.CurrentToken()
		if curTok.Type == TokNumber {
			// Parse as hex
			value := curTok.Value
			// Remove any prefix
			if len(value) > 0 && value[0] == '%' {
				value = value[1:]
			}
			num, err := strconv.ParseInt(value, 16, 32)
			if err != nil {
				ctx.AddError(fmt.Sprintf("invalid hex number %q: %v", curTok.Value, err))
				return nil
			}
			dataStmt.Values = append(dataStmt.Values, uint16(num))
			ctx.NextToken()
		} else if curTok.Type == TokIdent {
			// In data blocks, identifiers like ABCD are treated as hex
			value := curTok.Value
			num, err := strconv.ParseInt(value, 16, 32)
			if err != nil {
				ctx.AddError(fmt.Sprintf("invalid hex identifier %q in data block", value))
				return nil
			}
			dataStmt.Values = append(dataStmt.Values, uint16(num))
			ctx.NextToken()
		} else {
			ctx.AddError(fmt.Sprintf("expected hex number in data block, got %v", curTok.Type))
			return nil
		}
	}

	return dataStmt
}

// =============================================================================
// Def Parser
// =============================================================================

// DefParser parses def statements
type DefParser struct{}

func (dp *DefParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokDef
}

func (dp *DefParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	defStmt := &DefStmt{
		Line: curTok.Line,
		Col:  curTok.Column,
	}

	ctx.NextToken() // consume 'def'

	// Parse name/register
	curTok = ctx.CurrentToken()
	if curTok.Type != TokIdent {
		ctx.AddError("expected identifier after 'def'")
		return nil
	}
	defStmt.Name = curTok.Value
	ctx.NextToken()

	// Expect /
	if ctx.CurrentToken().Type != TokSlash {
		ctx.AddError("expected '/' after def name")
		return nil
	}
	ctx.NextToken()

	// Parse register
	reg, ok := ctx.ParseRegister()
	if !ok {
		ctx.AddError("expected register after '/'")
		return nil
	}
	defStmt.Register = reg

	// Register the definition in symbol table
	if err := ctx.Symbols().Define(defStmt.Name, defStmt.Register, defStmt.Line, defStmt.Col); err != nil {
		ctx.AddError(err.Error())
		return nil
	}

	// Check for optional initializer
	if ctx.CurrentToken().Type == TokEffectAlways || ctx.CurrentToken().Type == TokEquals {
		if ctx.CurrentToken().Type == TokEquals {
			// Convert = to <- for consistency
			ctx.NextToken() // consume =
		}
		// Parse the instruction
		defStmt.Initializer = ctx.ParseInstructionFromDef(defStmt.Name, defStmt.Register)
	}

	return defStmt
}

// =============================================================================
// Var Parser
// =============================================================================

// VarParser parses var statements
type VarParser struct{}

func (vp *VarParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokVar
}

func (vp *VarParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	varStmt := &VarStmt{
		Line: curTok.Line,
		Col:  curTok.Column,
	}

	ctx.NextToken() // consume 'var'

	// Parse qualified name: def.varname
	curTok = ctx.CurrentToken()
	if curTok.Type != TokIdent {
		ctx.AddError("expected identifier after 'var'")
		return nil
	}

	varStmt.DefName = curTok.Value
	ctx.NextToken()

	// Expect dot
	if ctx.CurrentToken().Type != TokDot {
		ctx.AddError("expected '.' after def name in var declaration")
		return nil
	}
	ctx.NextToken() // consume '.'

	// Parse variable name
	curTok = ctx.CurrentToken()
	if curTok.Type != TokIdent {
		ctx.AddError("expected variable name after '.'")
		return nil
	}
	varStmt.VarName = curTok.Value
	ctx.NextToken()

	// Expect /
	if ctx.CurrentToken().Type != TokSlash {
		ctx.AddError("expected '/' after var name")
		return nil
	}
	ctx.NextToken()

	// Parse offset
	offset, ok := ctx.ParseNumber()
	if !ok {
		return nil
	}
	varStmt.Offset = int16(offset)

	// Register the variable in symbol table
	if err := ctx.Symbols().DefineVar(varStmt.DefName, varStmt.VarName, varStmt.Offset, varStmt.Line, varStmt.Col); err != nil {
		ctx.AddError(err.Error())
		return nil
	}

	// Check for push modifier
	if ctx.CurrentToken().Type == TokPush {
		varStmt.Push = true
		ctx.NextToken()
	}

	// Check for optional initializer
	if ctx.CurrentToken().Type == TokEffectAlways || ctx.CurrentToken().Type == TokEquals {
		if ctx.CurrentToken().Type == TokEquals {
			ctx.NextToken() // consume =
		}

		// Look up the register for this def
		def, ok := ctx.Symbols().LookupDef(varStmt.DefName)
		if !ok {
			ctx.AddError("undefined definition \"" + varStmt.DefName + "\"")
			return nil
		}

		// Parse instruction with def.varname as destination
		curTok = ctx.CurrentToken()
		inst := &InstructionStmt{
			Line: curTok.Line,
			Col:  curTok.Column,
			Destination: &Operand{
				Type:     OpMemory,
				Register: def.Register,
				Offset:   varStmt.Offset,
			},
			Push: varStmt.Push,
		}

		// Parse effect
		inst.Effect = ctx.ParseEffect()

		// Parse operation
		inst.Operation = ctx.ParseOperation()

		// Parse source
		source, pop := ctx.ParseSource()
		if source == nil {
			return nil
		}
		inst.Source = source
		inst.Pop = pop

		varStmt.Initializer = inst
	}

	return varStmt
}

// =============================================================================
// Scope Parser
// =============================================================================

// ScopeParser parses scope block statements
type ScopeParser struct{}

func (sp *ScopeParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokLBrace
}

func (sp *ScopeParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	scopeStmt := &ScopeStmt{
		Line:       curTok.Line,
		Col:        curTok.Column,
		Statements: []Statement{},
	}

	ctx.NextToken() // consume '{'

	// Push a new scope
	ctx.Symbols().PushScope()

	// Parse statements until '}'
	for ctx.CurrentToken().Type != TokRBrace && ctx.CurrentToken().Type != TokEOF {
		// Skip empty lines and comments
		if ctx.CurrentToken().Type == TokNewline || ctx.CurrentToken().Type == TokComment {
			ctx.NextToken()
			continue
		}

		stmt := ctx.ParseStatement()
		if stmt != nil {
			scopeStmt.Statements = append(scopeStmt.Statements, stmt)
		} else {
			// Parse error occurred, skip to next line to avoid infinite loop
			for ctx.CurrentToken().Type != TokNewline && ctx.CurrentToken().Type != TokEOF {
				ctx.NextToken()
			}
		}

		// Consume trailing comments and newlines
		for ctx.CurrentToken().Type == TokComment || ctx.CurrentToken().Type == TokNewline {
			ctx.NextToken()
		}
	}

	if ctx.CurrentToken().Type != TokRBrace {
		ctx.AddError("expected '}' to close scope")
		return nil
	}
	ctx.NextToken() // consume '}'

	// Pop the scope
	ctx.Symbols().PopScope()

	return scopeStmt
}

// =============================================================================
// Call Parser
// =============================================================================

// CallParser parses standalone function call statements
type CallParser struct{}

func (cp *CallParser) Matches(ctx ParseContext) bool {
	// Pattern: IDENT DOT IDENT LPAREN
	// Parser will restore state after this, so we can consume tokens freely

	// Check: IDENT
	if ctx.CurrentToken().Type != TokIdent {
		return false
	}
	ctx.NextToken()

	// Check: DOT
	if ctx.CurrentToken().Type != TokDot {
		return false
	}
	ctx.NextToken()

	// Check: IDENT
	if ctx.CurrentToken().Type != TokIdent {
		return false
	}
	ctx.NextToken()

	// Check: LPAREN
	return ctx.CurrentToken().Type == TokLParen
}

func (cp *CallParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	callStmt := &CallStmt{
		Line: curTok.Line,
		Col:  curTok.Column,
	}

	// Parse: stack.funcName(args) -> result
	name := curTok.Value
	ctx.NextToken() // consume identifier

	ctx.NextToken() // consume dot

	name = name + "." + ctx.CurrentToken().Value
	ctx.NextToken() // consume function name

	// Now at (
	fc := ctx.ParseFunctionCall(name)
	if fc == nil {
		return nil
	}
	callStmt.Call = fc

	// Handle -> result part
	if ctx.CurrentToken().Type == TokArrow {
		ctx.NextToken() // consume ->
		if ctx.CurrentToken().Type == TokIdent {
			resultName := ctx.CurrentToken().Value
			// Define stack.result pointing to offset 0 (where return value is)
			parts := strings.Split(name, ".")
			if len(parts) == 2 {
				defName := parts[0]
				if err := ctx.Symbols().DefineVar(defName, resultName, 0, callStmt.Line, callStmt.Col); err != nil {
					ctx.AddError(err.Error())
					return nil
				}
			}
			ctx.NextToken() // consume result name
		}
	}

	return callStmt
}

// =============================================================================
// Fun Parser
// =============================================================================

// FunParser parses fun statements
type FunParser struct{}

func (fp *FunParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokFun
}

func (fp *FunParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	funStmt := &FunStmt{
		Line:    curTok.Line,
		Col:     curTok.Column,
		Args:    []string{},
		Returns: []string{},
		Body:    []Statement{},
	}

	ctx.NextToken() // consume 'fun'

	// Parse qualified name: stack.functionName
	curTok = ctx.CurrentToken()
	if curTok.Type != TokIdent {
		ctx.AddError("expected identifier after 'fun'")
		return nil
	}

	funStmt.DefName = curTok.Value
	ctx.NextToken()

	// Expect dot
	if ctx.CurrentToken().Type != TokDot {
		ctx.AddError("expected '.' after stack name in function declaration")
		return nil
	}
	ctx.NextToken() // consume '.'

	// Parse function name
	curTok = ctx.CurrentToken()
	if curTok.Type != TokIdent {
		ctx.AddError("expected function name after '.'")
		return nil
	}
	funStmt.Name = curTok.Value
	ctx.NextToken()

	// Parse arguments (arg1, arg2, ...)
	if ctx.CurrentToken().Type == TokLParen {
		ctx.NextToken() // consume '('

		for ctx.CurrentToken().Type != TokRParen && ctx.CurrentToken().Type != TokEOF {
			if ctx.CurrentToken().Type != TokIdent {
				ctx.AddError("expected argument name")
				return nil
			}
			funStmt.Args = append(funStmt.Args, ctx.CurrentToken().Value)
			ctx.NextToken()

			if ctx.CurrentToken().Type == TokComma {
				ctx.NextToken()
			}
		}

		if ctx.CurrentToken().Type != TokRParen {
			ctx.AddError("expected ')' after arguments")
			return nil
		}
		ctx.NextToken() // consume ')'
	}

	// Parse return values (-> ret1, ret2)
	if ctx.CurrentToken().Type == TokArrow {
		ctx.NextToken() // consume '->'

		for ctx.CurrentToken().Type == TokIdent {
			funStmt.Returns = append(funStmt.Returns, ctx.CurrentToken().Value)
			ctx.NextToken()

			if ctx.CurrentToken().Type == TokComma {
				ctx.NextToken()
			} else {
				break
			}
		}
	}

	// Parse body { ... }
	if ctx.CurrentToken().Type != TokLBrace {
		ctx.AddError("expected '{' to start function body")
		return nil
	}

	// Before parsing the body, we need to set up the function's stack frame
	// Push a new scope for the function
	ctx.NextToken() // consume '{'
	ctx.Symbols().PushScope()

	// According to uCISC calling convention:
	// Stack layout (from stack pointer upward):
	//   offset 0: arguments (pushed last, so at top of stack)
	//   offset 1: return address
	//   offset 2+: return values (reserved by caller, pushed first)

	// Define variables for arguments, return address, and return values
	offset := int16(0)

	// Arguments are at the top (offset 0+)
	for _, argName := range funStmt.Args {
		if err := ctx.Symbols().DefineVar(funStmt.DefName, argName, offset, funStmt.Line, funStmt.Col); err != nil {
			ctx.AddError(err.Error())
			return nil
		}
		offset++
	}

	// Return address is next
	if err := ctx.Symbols().DefineVar(funStmt.DefName, "return", offset, funStmt.Line, funStmt.Col); err != nil {
		ctx.AddError(err.Error())
		return nil
	}
	offset++

	// Return values are at the bottom (highest offsets)
	for _, retName := range funStmt.Returns {
		if err := ctx.Symbols().DefineVar(funStmt.DefName, retName, offset, funStmt.Line, funStmt.Col); err != nil {
			ctx.AddError(err.Error())
			return nil
		}
		offset++
	}

	// Parse statements until '}'
	for ctx.CurrentToken().Type != TokRBrace && ctx.CurrentToken().Type != TokEOF {
		// Skip empty lines and comments
		if ctx.CurrentToken().Type == TokNewline || ctx.CurrentToken().Type == TokComment {
			ctx.NextToken()
			continue
		}

		stmt := ctx.ParseStatement()
		if stmt != nil {
			funStmt.Body = append(funStmt.Body, stmt)
		} else {
			// Parse error occurred, skip to next line to avoid infinite loop
			for ctx.CurrentToken().Type != TokNewline && ctx.CurrentToken().Type != TokEOF {
				ctx.NextToken()
			}
		}

		// Consume trailing comments and newlines
		for ctx.CurrentToken().Type == TokComment || ctx.CurrentToken().Type == TokNewline {
			ctx.NextToken()
		}
	}

	if ctx.CurrentToken().Type != TokRBrace {
		ctx.AddError("expected '}' to close function body")
		return nil
	}
	ctx.NextToken() // consume '}'

	// Pop the function scope
	ctx.Symbols().PopScope()

	return funStmt
}

// =============================================================================
// Instruction Parser (catch-all)
// =============================================================================

// InstructionParser parses instruction statements (catch-all)
type InstructionParser struct{}

func (ip *InstructionParser) Matches(ctx ParseContext) bool {
	// Always matches - this should be last in the list
	return true
}

func (ip *InstructionParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	inst := &InstructionStmt{
		Line: curTok.Line,
		Col:  curTok.Column,
	}

	// Parse destination
	dest, push := ctx.ParseDestination()
	if dest == nil {
		return nil
	}
	inst.Destination = dest
	inst.Push = push

	// Parse effect
	effect := ctx.ParseEffect()
	inst.Effect = effect

	// Parse operation
	op := ctx.ParseOperation()
	inst.Operation = op

	// Parse source
	source, pop := ctx.ParseSource()
	if source == nil {
		return nil
	}
	inst.Source = source
	inst.Pop = pop

	return inst
}
