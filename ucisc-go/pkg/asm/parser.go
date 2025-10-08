package asm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grokthis/ucisc-go/pkg/cpu"
)

// ParseContext provides the minimal interface needed by element parsers
// This follows the Interface Segregation Principle - parsers only depend on what they need
type ParseContext interface {
	// Token stream access
	CurrentToken() Token
	PeekToken() Token
	NextToken()

	// Parser helpers
	ParseNumber() (int, bool)
	ParseNumberInContext(ctx NumberContext) (int, bool)
	ParseRegister() (cpu.Register, bool)
	ParseQualifiedName() (string, bool)
	ParseAdditionalOffset(op *Operand) bool
	ParseEffect() cpu.Effect
	ParseOperation() cpu.Operation

	// Symbol table access
	Symbols() *SymbolTable

	// Error reporting
	AddError(msg string)
}

// Parser parses uCISC assembly source into an AST
type Parser struct {
	lexer   *Lexer
	curTok  Token
	peekTok Token
	errors  []string
	symtab  *SymbolTable // Symbol table for def/var resolution
}

// NewParser creates a new parser for the given input
func NewParser(input string) *Parser {
	p := &Parser{
		lexer:  NewLexer(input),
		errors: []string{},
		symtab: NewSymbolTable(),
	}
	// Initialize element parsers once
	if len(elementParsers) == 0 {
		InitElementParsers()
	}
	// Read two tokens to initialize curTok and peekTok
	p.nextToken()
	p.nextToken()
	return p
}

// Parse parses the input and returns the AST
func (p *Parser) Parse() (*Program, error) {
	program := &Program{
		Statements: []Statement{},
	}

	for p.curTok.Type != TokEOF {
		// Skip empty lines and comments
		if p.curTok.Type == TokNewline || p.curTok.Type == TokComment {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		} else {
			// Parse error occurred, skip to next line to avoid infinite loop
			for p.curTok.Type != TokNewline && p.curTok.Type != TokEOF {
				p.nextToken()
			}
		}

		// Consume trailing comments and newlines
		for p.curTok.Type == TokComment || p.curTok.Type == TokNewline {
			p.nextToken()
		}
	}

	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parse errors:\n%s", strings.Join(p.errors, "\n"))
	}

	return program, nil
}

// Errors returns the list of parse errors
func (p *Parser) Errors() []string {
	return p.errors
}

// nextToken advances to the next token
func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.lexer.NextToken()
}

// parseStatement parses a single statement using element parsers
func (p *Parser) parseStatement() Statement {
	// Try each element parser in order
	for _, ep := range elementParsers {
		// Save parser state before calling Matches
		savedCur := p.curTok
		savedPeek := p.peekTok
		savedMark := p.lexer.Mark()

		// Check if this parser matches
		matches := ep.Matches(p)

		// Restore parser state after Matches
		p.curTok = savedCur
		p.peekTok = savedPeek
		p.lexer.Restore(savedMark)

		if matches {
			return ep.Parse(p)
		}
	}

	// No parser matched - this shouldn't happen if InstructionParser is last
	p.addError(fmt.Sprintf("unexpected token %s", p.curTok))
	// Consume the token to avoid infinite loop
	p.nextToken()
	return nil
}

// parseCallStmt parses a standalone function call statement
// Example: stack.double(val/2) -> result
// Returns nil if this is not a function call (no error added)
func (p *Parser) parseCallStmt() Statement {
	// First, peek ahead to see if this matches the pattern: IDENT DOT IDENT LPAREN
	// We have curTok=IDENT and peekTok=DOT already checked by caller
	// We need to peek further to see IDENT LPAREN
	// Since we can't peek that far, we'll speculatively parse and backtrack if needed

	// Save parser state
	savedCur := p.curTok
	savedPeek := p.peekTok

	callStmt := &CallStmt{
		Line: p.curTok.Line,
		Col:  p.curTok.Column,
	}

	// Parse the function call: stack.funcName(args)
	name := p.curTok.Value
	p.nextToken()

	// Expect dot
	if p.curTok.Type != TokDot {
		// Not a function call, restore and return nil
		p.curTok = savedCur
		p.peekTok = savedPeek
		return nil
	}
	p.nextToken() // consume '.'

	// Get function name
	if p.curTok.Type != TokIdent {
		// Not a function call, restore and return nil
		p.curTok = savedCur
		p.peekTok = savedPeek
		return nil
	}
	name = name + "." + p.curTok.Value
	p.nextToken()

	if p.curTok.Type != TokLParen {
		// Not a function call, restore and return nil
		p.curTok = savedCur
		p.peekTok = savedPeek
		return nil
	}

	fc := p.parseFunctionCall(name)
	if fc == nil {
		return nil
	}
	callStmt.Call = fc

	// The -> result part creates a variable pointing to the return value
	// The return value is at offset 0 on the stack
	if p.curTok.Type == TokArrow {
		p.nextToken() // consume ->
		if p.curTok.Type == TokIdent {
			resultName := p.curTok.Value
			// Define stack.result pointing to offset 0 (where return value is)
			parts := strings.Split(name, ".")
			if len(parts) == 2 {
				defName := parts[0]
				if err := p.symtab.DefineVar(defName, resultName, 0, callStmt.Line, callStmt.Col); err != nil {
					p.addError(err.Error())
					return nil
				}
			}
			p.nextToken() // consume result name
		}
	}

	return callStmt
}

// parseLabel parses a label declaration
func (p *Parser) parseLabel() Statement {
	label := &LabelStmt{
		Name: p.curTok.Value,
		Line: p.curTok.Line,
		Col:  p.curTok.Column,
	}

	p.nextToken() // consume label name
	p.nextToken() // consume colon

	return label
}

// parseInstruction parses an instruction
func (p *Parser) parseInstruction() Statement {
	inst := &InstructionStmt{
		Line: p.curTok.Line,
		Col:  p.curTok.Column,
	}

	// Parse destination
	dest, push := p.parseDestination()
	if dest == nil {
		return nil
	}
	inst.Destination = dest
	inst.Push = push

	// Parse effect
	effect := p.parseEffect()
	inst.Effect = effect

	// Parse operation
	op := p.parseOperation()
	inst.Operation = op

	// Parse source
	source, pop := p.parseSource()
	if source == nil {
		return nil
	}
	inst.Source = source
	inst.Pop = pop

	return inst
}

// parseDestination parses the destination operand
func (p *Parser) parseDestination() (*Operand, bool) {
	push := false
	op := &Operand{
		Line: p.curTok.Line,
		Col:  p.curTok.Column,
	}

	// Check for & prefix (register destination)
	if p.curTok.Type == TokAmpersand {
		op.IsAddress = true
		p.nextToken()
	}

	// Try to parse as register first
	reg, ok := p.parseRegister()
	if ok {
		op.Register = reg

		// Check if there's a slash (offset) - indicates memory access
		hasOffset := p.curTok.Type == TokSlash

		if hasOffset {
			p.nextToken() // consume /
			offset, ok := p.parseNumber()
			if !ok {
				p.addError("expected number after /")
				return nil, false
			}
			op.Offset = int16(offset)
		}

		// Determine type:
		// - If & prefix OR no offset: register destination
		// - If offset: memory destination
		if op.IsAddress || !hasOffset {
			op.Type = OpRegister
			return op, false
		} else {
			op.Type = OpMemory
		}

		// Check for push modifier
		if p.curTok.Type == TokPush {
			push = true
			p.nextToken()
		}

		return op, push
	}

	// Not a register, try to resolve as def/var name
	name, ok := p.parseQualifiedName()
	if ok {
		// Try to resolve the name
		resolvedReg, resolvedOffset, ok := p.symtab.ResolveOperand(name)
		if !ok {
			p.addError(fmt.Sprintf("undefined identifier %q", name))
			return nil, false
		}

		op.Register = resolvedReg
		op.Offset = resolvedOffset

		// Check for additional offset (/N)
		if !p.parseAdditionalOffset(op) {
			return nil, false
		}

		// Determine type based on & prefix
		if op.IsAddress {
			op.Type = OpRegister
		} else {
			op.Type = OpMemory
		}

		// Check for push modifier
		if p.curTok.Type == TokPush {
			push = true
			p.nextToken()
		}

		return op, push
	}

	p.addError("expected register or identifier in destination")
	return nil, false
}

// parseSource parses the source operand
func (p *Parser) parseSource() (*Operand, bool) {
	pop := false
	op := &Operand{
		Line: p.curTok.Line,
		Col:  p.curTok.Column,
	}

	// Check for & prefix (register value source)
	isRegisterValue := false
	if p.curTok.Type == TokAmpersand {
		isRegisterValue = true
		p.nextToken() // consume &
	}

	// Check for val (immediate)
	if p.curTok.Type == TokVal {
		p.nextToken() // consume val
		if p.curTok.Type != TokSlash {
			p.addError("expected / after val")
			return nil, false
		}
		p.nextToken() // consume /

		// Check for number or label
		if p.curTok.Type == TokNumber {
			num, ok := p.parseNumber()
			if !ok {
				return nil, false
			}
			op.Offset = int16(num)
			op.Type = OpImmediate
		} else if p.curTok.Type == TokIdent {
			// Label reference for immediate
			op.LabelRef = p.curTok.Value
			op.Type = OpImmediate
			p.nextToken()
		} else {
			p.addError("expected number or label after val/")
			return nil, false
		}
		return op, false
	}

	// Try to parse as register
	reg, ok := p.parseRegister()
	if ok {
		op.Register = reg

		// Check for offset
		if p.curTok.Type == TokSlash {
			p.nextToken() // consume /

			// Could be number or label
			if p.curTok.Type == TokNumber {
				offset, ok := p.parseNumber()
				if !ok {
					return nil, false
				}
				op.Offset = int16(offset)
			} else if p.curTok.Type == TokIdent {
				// Label reference (could be qualified like stack.funcName)
				op.LabelRef = p.curTok.Value
				p.nextToken()

				// Check for qualified label (e.g., stack.funcName)
				if p.curTok.Type == TokDot {
					p.nextToken() // consume dot
					if p.curTok.Type != TokIdent {
						p.addError("expected identifier after '.' in label reference")
						return nil, false
					}
					op.LabelRef = op.LabelRef + "." + p.curTok.Value
					p.nextToken()
				}
			} else {
				p.addError("expected number or label after /")
				return nil, false
			}
		}

		// Determine type
		if isRegisterValue {
			// &r1 as source means direct register value
			op.Type = OpRegister
		} else if reg == cpu.RegPC {
			op.Type = OpPCRelative
		} else {
			// All non-PC registers default to memory access (r1 = r1/0)
			// This aligns with uCISC philosophy: registers are pointers
			op.Type = OpMemory
		}

		// Check for pop modifier
		if p.curTok.Type == TokPop {
			pop = true
			p.nextToken()
		}

		return op, pop
	}

	// Not a register, try to resolve as def/var name or function call
	name, ok := p.parseQualifiedName()
	if ok {
		// Check if this is a function call: name(args)
		if p.curTok.Type == TokLParen {
			funcCall := p.parseFunctionCall(name)
			if funcCall == nil {
				return nil, false
			}
			op.FunctionCall = funcCall
			op.Type = OpImmediate // Placeholder type
			return op, false
		}

		// Try to resolve the name
		resolvedReg, resolvedOffset, ok := p.symtab.ResolveOperand(name)
		if !ok {
			p.addError(fmt.Sprintf("undefined identifier %q", name))
			return nil, false
		}

		op.Register = resolvedReg
		op.Offset = resolvedOffset

		// Check for additional offset (/N)
		if !p.parseAdditionalOffset(op) {
			return nil, false
		}

		// Determine type
		if isRegisterValue {
			op.Type = OpRegister
		} else if resolvedReg == cpu.RegPC {
			op.Type = OpPCRelative
		} else {
			op.Type = OpMemory
		}

		// Check for pop modifier
		if p.curTok.Type == TokPop {
			pop = true
			p.nextToken()
		}

		return op, pop
	}

	p.addError("expected register, val, or identifier in source")
	return nil, false
}

// parseFunctionCall parses a function call: defName.funcName(arg1, arg2)
func (p *Parser) parseFunctionCall(qualifiedName string) *FunctionCall {
	// qualifiedName should already be in "def.func" format from parseSource/parseCallStmt
	// which combine the parts with dots
	parts := strings.Split(qualifiedName, ".")
	if len(parts) != 2 {
		p.addError("function call must be qualified like 'stack.funcName'")
		return nil
	}

	fc := &FunctionCall{
		DefName:  parts[0],
		FuncName: parts[1],
		Args:     []*Operand{},
		Line:     p.curTok.Line,
		Col:      p.curTok.Column,
	}

	p.nextToken() // consume '('

	// Parse arguments
	for p.curTok.Type != TokRParen && p.curTok.Type != TokEOF {
		// Parse argument as a source operand
		arg, _ := p.parseSource()
		if arg == nil {
			return nil
		}
		fc.Args = append(fc.Args, arg)

		if p.curTok.Type == TokComma {
			p.nextToken()
		}
	}

	if p.curTok.Type != TokRParen {
		p.addError("expected ')' after function arguments")
		return nil
	}
	p.nextToken() // consume ')'

	return fc
}

// parseEffect parses an effect token
func (p *Parser) parseEffect() cpu.Effect {
	var effect cpu.Effect

	switch p.curTok.Type {
	case TokEffectAlways:
		effect = cpu.EffectAlways
	case TokEffectFlags:
		effect = cpu.EffectFlagsOnly
	case TokEffectZero:
		effect = cpu.EffectZero
	case TokEffectNotZero:
		effect = cpu.EffectNotZero
	case TokEffectPositive:
		effect = cpu.EffectPositive
	case TokEffectNegative:
		effect = cpu.EffectNegative
	case TokEffectOverflow:
		effect = cpu.EffectOverflow
	case TokEffectInterrupt:
		effect = cpu.EffectInterrupt
	default:
		p.addError(fmt.Sprintf("expected effect, got %v", p.curTok.Type))
		return cpu.EffectAlways
	}

	p.nextToken()
	return effect
}

// parseOperation parses an operation token
func (p *Parser) parseOperation() cpu.Operation {
	var op cpu.Operation

	switch p.curTok.Type {
	case TokCopy:
		op = cpu.OpCopy
	case TokAnd:
		op = cpu.OpAnd
	case TokOr:
		op = cpu.OpOr
	case TokXor:
		op = cpu.OpXor
	case TokInv:
		op = cpu.OpInv
	case TokShl:
		op = cpu.OpShl
	case TokShr:
		op = cpu.OpShr
	case TokSwap:
		op = cpu.OpSwap
	case TokMsb:
		op = cpu.OpMsb
	case TokLsb:
		op = cpu.OpLsb
	case TokAdd:
		op = cpu.OpAdd
	case TokAddc:
		op = cpu.OpAddc
	case TokSub:
		op = cpu.OpSub
	case TokMult:
		op = cpu.OpMult
	case TokMultsw:
		op = cpu.OpMultMsw
	default:
		p.addError(fmt.Sprintf("expected operation, got %v", p.curTok.Type))
		return cpu.OpCopy
	}

	p.nextToken()
	return op
}

// parseRegister parses a register token
func (p *Parser) parseRegister() (cpu.Register, bool) {
	var reg cpu.Register

	switch p.curTok.Type {
	case TokR1:
		reg = cpu.RegR1
	case TokR2:
		reg = cpu.RegR2
	case TokR3:
		reg = cpu.RegR3
	case TokR5:
		reg = cpu.RegR5
	case TokR6:
		reg = cpu.RegR6
	case TokR7:
		reg = cpu.RegR7
	case TokPC:
		reg = cpu.RegPC
	case TokControl, TokBanking:
		reg = cpu.RegBanking
	case TokFlags:
		reg = cpu.RegFlags
	case TokInterrupt:
		reg = cpu.RegInt
	default:
		return 0, false
	}

	p.nextToken()
	return reg, true
}

// NumberContext indicates the parsing context for numbers
type NumberContext int

const (
	NumberContextNormal NumberContext = iota // Instructions, offsets, etc.
	NumberContextData                        // Inside data blocks (% ...)
)

// parseNumber parses a number token (normal context)
func (p *Parser) parseNumber() (int, bool) {
	return p.parseNumberInContext(NumberContextNormal)
}

// parseNumberInContext parses a number with context-specific rules
func (p *Parser) parseNumberInContext(ctx NumberContext) (int, bool) {
	// Handle identifiers in data blocks (e.g., ABCD, FFFF)
	if ctx == NumberContextData && p.curTok.Type == TokIdent {
		value := p.curTok.Value
		p.nextToken()

		// Try to parse as hex
		num, err := strconv.ParseInt(value, 16, 32)
		if err != nil {
			p.addError(fmt.Sprintf("invalid hex value %q in data block: %v", value, err))
			return 0, false
		}
		return int(num), true
	}

	if p.curTok.Type != TokNumber {
		p.addError(fmt.Sprintf("expected number, got %v", p.curTok.Type))
		return 0, false
	}

	value := p.curTok.Value
	p.nextToken()

	var num int64
	var err error

	// Check if it already has a prefix
	if strings.HasPrefix(value, "%") {
		// %FFFF format
		num, err = strconv.ParseInt(value[1:], 16, 32)
	} else if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		// 0xFFFF format
		num, err = strconv.ParseInt(value[2:], 16, 32)
	} else if ctx == NumberContextData {
		// In data blocks, assume hex by default
		num, err = strconv.ParseInt(value, 16, 32)
		if err != nil {
			// If hex parsing fails, try decimal as fallback
			num, err = strconv.ParseInt(value, 10, 32)
		}
	} else {
		// Normal context: decimal
		num, err = strconv.ParseInt(value, 10, 32)
	}

	if err != nil {
		p.addError(fmt.Sprintf("invalid number %q: %v", value, err))
		return 0, false
	}

	return int(num), true
}

// addError adds a parse error
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("line %d:%d: %s",
		p.curTok.Line, p.curTok.Column, msg))
}

// parseQualifiedName parses a qualified identifier (e.g., stack.value)
// Returns the full name and whether parsing succeeded
func (p *Parser) parseQualifiedName() (string, bool) {
	if p.curTok.Type != TokIdent {
		return "", false
	}

	name := p.curTok.Value
	p.nextToken()

	// Check for qualified name (e.g., stack.output)
	if p.curTok.Type == TokDot {
		p.nextToken() // consume dot
		if p.curTok.Type != TokIdent {
			p.addError("expected identifier after '.'")
			return "", false
		}
		name = name + "." + p.curTok.Value
		p.nextToken()
	}

	return name, true
}

// parseAdditionalOffset parses an optional /N offset and adds it to the operand
func (p *Parser) parseAdditionalOffset(op *Operand) bool {
	if p.curTok.Type == TokSlash {
		p.nextToken()
		additionalOffset, ok := p.parseNumber()
		if !ok {
			return false
		}
		op.Offset += int16(additionalOffset)
	}
	return true
}

// parseDef parses a def statement
// def name/register [<- copy val/init]
func (p *Parser) parseDef() Statement {
	defStmt := &DefStmt{
		Line: p.curTok.Line,
		Col:  p.curTok.Column,
	}

	p.nextToken() // consume 'def'

	// Parse name/register
	if p.curTok.Type != TokIdent {
		p.addError("expected identifier after 'def'")
		return nil
	}
	defStmt.Name = p.curTok.Value
	p.nextToken()

	// Expect /
	if p.curTok.Type != TokSlash {
		p.addError("expected '/' after def name")
		return nil
	}
	p.nextToken()

	// Parse register
	reg, ok := p.parseRegister()
	if !ok {
		p.addError("expected register after '/'")
		return nil
	}
	defStmt.Register = reg

	// Register the definition in symbol table
	if err := p.symtab.Define(defStmt.Name, defStmt.Register, defStmt.Line, defStmt.Col); err != nil {
		p.addError(err.Error())
		return nil
	}

	// Check for optional initializer
	if p.curTok.Type == TokEffectAlways || p.curTok.Type == TokEquals {
		if p.curTok.Type == TokEquals {
			// Convert = to <- for consistency
			p.nextToken() // consume =
		}
		// Parse the instruction
		defStmt.Initializer = p.parseInstructionFromDef(defStmt.Name, defStmt.Register)
	}

	return defStmt
}

// parseInstructionFromDef parses an instruction that starts with a def/var
// The destination is already parsed, so we just need effect + operation + source
func (p *Parser) parseInstructionFromDef(name string, reg cpu.Register) *InstructionStmt {
	inst := &InstructionStmt{
		Line: p.curTok.Line,
		Col:  p.curTok.Column,
		Destination: &Operand{
			Type:      OpRegister,
			Register:  reg,
			IsAddress: true, // def/var always targets the register directly
		},
	}

	// Parse effect
	inst.Effect = p.parseEffect()

	// Parse operation
	inst.Operation = p.parseOperation()

	// Parse source
	source, pop := p.parseSource()
	if source == nil {
		return nil
	}
	inst.Source = source
	inst.Pop = pop

	return inst
}

// parseVar parses a var statement
// var def.name/offset [push] [<- operation source]
func (p *Parser) parseVar() Statement {
	varStmt := &VarStmt{
		Line: p.curTok.Line,
		Col:  p.curTok.Column,
	}

	p.nextToken() // consume 'var'

	// Parse qualified name: def.varname
	if p.curTok.Type != TokIdent {
		p.addError("expected identifier after 'var'")
		return nil
	}

	varStmt.DefName = p.curTok.Value
	p.nextToken()

	// Expect dot
	if p.curTok.Type != TokDot {
		p.addError("expected '.' after def name in var declaration")
		return nil
	}
	p.nextToken() // consume '.'

	// Parse variable name
	if p.curTok.Type != TokIdent {
		p.addError("expected variable name after '.'")
		return nil
	}
	varStmt.VarName = p.curTok.Value
	p.nextToken()

	// Expect /
	if p.curTok.Type != TokSlash {
		p.addError("expected '/' after var name")
		return nil
	}
	p.nextToken()

	// Parse offset
	offset, ok := p.parseNumber()
	if !ok {
		return nil
	}
	varStmt.Offset = int16(offset)

	// Register the variable in symbol table
	if err := p.symtab.DefineVar(varStmt.DefName, varStmt.VarName, varStmt.Offset, varStmt.Line, varStmt.Col); err != nil {
		p.addError(err.Error())
		return nil
	}

	// Check for push modifier
	if p.curTok.Type == TokPush {
		varStmt.Push = true
		p.nextToken()
	}

	// Check for optional initializer
	if p.curTok.Type == TokEffectAlways || p.curTok.Type == TokEquals {
		if p.curTok.Type == TokEquals {
			p.nextToken() // consume =
		}

		// Look up the register for this def
		def, ok := p.symtab.LookupDef(varStmt.DefName)
		if !ok {
			p.addError(fmt.Sprintf("undefined definition %q", varStmt.DefName))
			return nil
		}

		// Parse instruction with def.varname as destination
		inst := &InstructionStmt{
			Line: p.curTok.Line,
			Col:  p.curTok.Column,
			Destination: &Operand{
				Type:     OpMemory,
				Register: def.Register,
				Offset:   varStmt.Offset,
			},
			Push: varStmt.Push,
		}

		// Parse effect
		inst.Effect = p.parseEffect()

		// Parse operation
		inst.Operation = p.parseOperation()

		// Parse source
		source, pop := p.parseSource()
		if source == nil {
			return nil
		}
		inst.Source = source
		inst.Pop = pop

		varStmt.Initializer = inst
	}

	return varStmt
}

// parseData parses a data block
// % 0042 00FF 1234 ABCD
func (p *Parser) parseData() Statement {
	dataStmt := &DataStmt{
		Line:   p.curTok.Line,
		Col:    p.curTok.Column,
		Values: []uint16{},
	}

	p.nextToken() // consume '%'

	// Read hex numbers until end of line
	// In data blocks, numbers are assumed to be hex by default
	for p.curTok.Type != TokNewline && p.curTok.Type != TokEOF && p.curTok.Type != TokComment {
		if p.curTok.Type == TokNumber || p.curTok.Type == TokIdent {
			// Use data context: assumes hex, accepts identifiers like ABCD
			num, ok := p.parseNumberInContext(NumberContextData)
			if !ok {
				return nil
			}
			dataStmt.Values = append(dataStmt.Values, uint16(num))
		} else {
			// Skip whitespace and continue
			p.nextToken()
		}
	}

	return dataStmt
}


// parseString parses a string literal
func (p *Parser) parseString() Statement {
	stringStmt := &StringStmt{
		Line:  p.curTok.Line,
		Col:   p.curTok.Column,
		Value: p.curTok.Value,
	}

	p.nextToken() // consume string

	return stringStmt
}

// parseFun parses a function declaration
// fun stack.name(arg1, arg2) -> ret1, ret2 { body }
func (p *Parser) parseFun() Statement {
	funStmt := &FunStmt{
		Line:    p.curTok.Line,
		Col:     p.curTok.Column,
		Args:    []string{},
		Returns: []string{},
		Body:    []Statement{},
	}

	p.nextToken() // consume 'fun'

	// Parse qualified name: stack.functionName
	if p.curTok.Type != TokIdent {
		p.addError("expected identifier after 'fun'")
		return nil
	}

	funStmt.DefName = p.curTok.Value
	p.nextToken()

	// Expect dot
	if p.curTok.Type != TokDot {
		p.addError("expected '.' after stack name in function declaration")
		return nil
	}
	p.nextToken() // consume '.'

	// Parse function name
	if p.curTok.Type != TokIdent {
		p.addError("expected function name after '.'")
		return nil
	}
	funStmt.Name = p.curTok.Value
	p.nextToken()

	// Parse arguments (arg1, arg2, ...)
	if p.curTok.Type == TokLParen {
		p.nextToken() // consume '('

		for p.curTok.Type != TokRParen && p.curTok.Type != TokEOF {
			if p.curTok.Type != TokIdent {
				p.addError("expected argument name")
				return nil
			}
			funStmt.Args = append(funStmt.Args, p.curTok.Value)
			p.nextToken()

			if p.curTok.Type == TokComma {
				p.nextToken()
			}
		}

		if p.curTok.Type != TokRParen {
			p.addError("expected ')' after arguments")
			return nil
		}
		p.nextToken() // consume ')'
	}

	// Parse return values (-> ret1, ret2)
	if p.curTok.Type == TokArrow {
		p.nextToken() // consume '->'

		for p.curTok.Type == TokIdent {
			funStmt.Returns = append(funStmt.Returns, p.curTok.Value)
			p.nextToken()

			if p.curTok.Type == TokComma {
				p.nextToken()
			} else {
				break
			}
		}
	}

	// Parse body { ... }
	if p.curTok.Type != TokLBrace {
		p.addError("expected '{' to start function body")
		return nil
	}

	// Before parsing the body, we need to set up the function's stack frame
	// Push a new scope for the function
	p.nextToken() // consume '{'
	p.symtab.PushScope()

	// According to uCISC calling convention:
	// Stack layout (from stack pointer upward):
	//   offset 0: arguments (pushed last, so at top of stack)
	//   offset 1: return address
	//   offset 2+: return values (reserved by caller, pushed first)

	// Define variables for arguments, return address, and return values
	offset := int16(0)

	// Arguments are at the top (offset 0+)
	for _, argName := range funStmt.Args {
		if err := p.symtab.DefineVar(funStmt.DefName, argName, offset, funStmt.Line, funStmt.Col); err != nil {
			p.addError(err.Error())
			return nil
		}
		offset++
	}

	// Return address is next
	if err := p.symtab.DefineVar(funStmt.DefName, "return", offset, funStmt.Line, funStmt.Col); err != nil {
		p.addError(err.Error())
		return nil
	}
	offset++

	// Return values are at the bottom (highest offsets)
	for _, retName := range funStmt.Returns {
		if err := p.symtab.DefineVar(funStmt.DefName, retName, offset, funStmt.Line, funStmt.Col); err != nil {
			p.addError(err.Error())
			return nil
		}
		offset++
	}

	// Parse statements until '}'
	for p.curTok.Type != TokRBrace && p.curTok.Type != TokEOF {
		// Skip empty lines and comments
		if p.curTok.Type == TokNewline || p.curTok.Type == TokComment {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			funStmt.Body = append(funStmt.Body, stmt)
		} else {
			// Parse error occurred, skip to next line to avoid infinite loop
			for p.curTok.Type != TokNewline && p.curTok.Type != TokEOF {
				p.nextToken()
			}
		}

		// Consume trailing comments and newlines
		for p.curTok.Type == TokComment || p.curTok.Type == TokNewline {
			p.nextToken()
		}
	}

	if p.curTok.Type != TokRBrace {
		p.addError("expected '}' to close function body")
		return nil
	}
	p.nextToken() // consume '}'

	// Pop the function scope
	p.symtab.PopScope()

	return funStmt
}

// parseScope parses a scope block { ... }
func (p *Parser) parseScope() Statement {
	scopeStmt := &ScopeStmt{
		Line:       p.curTok.Line,
		Col:        p.curTok.Column,
		Statements: []Statement{},
	}

	p.nextToken() // consume '{'

	// Push a new scope
	p.symtab.PushScope()

	// Parse statements until '}'
	for p.curTok.Type != TokRBrace && p.curTok.Type != TokEOF {
		// Skip empty lines and comments
		if p.curTok.Type == TokNewline || p.curTok.Type == TokComment {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			scopeStmt.Statements = append(scopeStmt.Statements, stmt)
		} else {
			// Parse error occurred, skip to next line to avoid infinite loop
			for p.curTok.Type != TokNewline && p.curTok.Type != TokEOF {
				p.nextToken()
			}
		}

		// Consume trailing comments and newlines
		for p.curTok.Type == TokComment || p.curTok.Type == TokNewline {
			p.nextToken()
		}
	}

	if p.curTok.Type != TokRBrace {
		p.addError("expected '}' to close scope")
		return nil
	}
	p.nextToken() // consume '}'

	// Pop the scope
	p.symtab.PopScope()

	return scopeStmt
}

// ParseContext interface implementation
// These methods allow Parser to be used as a ParseContext

func (p *Parser) CurrentToken() Token { return p.curTok }
func (p *Parser) PeekToken() Token    { return p.peekTok }
func (p *Parser) NextToken()          { p.nextToken() }

func (p *Parser) ParseNumber() (int, bool) { return p.parseNumber() }
func (p *Parser) ParseNumberInContext(ctx NumberContext) (int, bool) {
	return p.parseNumberInContext(ctx)
}
func (p *Parser) ParseRegister() (cpu.Register, bool) { return p.parseRegister() }
func (p *Parser) ParseQualifiedName() (string, bool)  { return p.parseQualifiedName() }
func (p *Parser) ParseAdditionalOffset(op *Operand) bool {
	return p.parseAdditionalOffset(op)
}
func (p *Parser) ParseEffect() cpu.Effect       { return p.parseEffect() }
func (p *Parser) ParseOperation() cpu.Operation { return p.parseOperation() }

func (p *Parser) Symbols() *SymbolTable { return p.symtab }
func (p *Parser) AddError(msg string)   { p.addError(msg) }
