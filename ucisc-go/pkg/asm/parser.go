package asm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grokthis/ucisc-go/pkg/cpu"
)

// Parser parses uCISC assembly source into an AST
type Parser struct {
	lexer   *Lexer
	curTok  Token
	peekTok Token
	errors  []string
}

// NewParser creates a new parser for the given input
func NewParser(input string) *Parser {
	p := &Parser{
		lexer:  NewLexer(input),
		errors: []string{},
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

// parseStatement parses a single statement
func (p *Parser) parseStatement() Statement {
	// Check for label (identifier followed by colon)
	if p.curTok.Type == TokIdent && p.peekTok.Type == TokColon {
		return p.parseLabel()
	}

	// Otherwise, it's an instruction
	return p.parseInstruction()
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

	// Parse register
	reg, ok := p.parseRegister()
	if !ok {
		p.addError("expected register in destination")
		return nil, false
	}
	op.Register = reg

	// If it's a register destination (&r1), we're done
	if op.IsAddress {
		op.Type = OpRegister
		return op, false
	}

	// Otherwise it's memory destination, check for offset
	if p.curTok.Type == TokSlash {
		p.nextToken() // consume /
		offset, ok := p.parseNumber()
		if !ok {
			p.addError("expected number after /")
			return nil, false
		}
		op.Offset = int16(offset)
	}

	op.Type = OpMemory

	// Check for push modifier
	if p.curTok.Type == TokPush {
		push = true
		p.nextToken()
	}

	return op, push
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

	// Parse register
	reg, ok := p.parseRegister()
	if !ok {
		p.addError("expected register or val in source")
		return nil, false
	}
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
			// Label reference
			op.LabelRef = p.curTok.Value
			p.nextToken()
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

// parseNumber parses a number token
func (p *Parser) parseNumber() (int, bool) {
	if p.curTok.Type != TokNumber {
		p.addError(fmt.Sprintf("expected number, got %v", p.curTok.Type))
		return 0, false
	}

	value := p.curTok.Value
	p.nextToken()

	// Handle different number formats
	var num int64
	var err error

	if strings.HasPrefix(value, "%") {
		// %FFFF format
		num, err = strconv.ParseInt(value[1:], 16, 32)
	} else if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		// 0xFFFF format
		num, err = strconv.ParseInt(value[2:], 16, 32)
	} else {
		// Decimal
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
