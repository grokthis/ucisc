package asm

import (
	"fmt"

	"github.com/grokthis/ucisc-go/pkg/cpu"
)

// CodeGenerator generates machine code from an AST
type CodeGenerator struct {
	program       *Program
	labels        map[string]uint16 // label name -> address
	errors        []string
	funcCallCount int               // Counter for generating unique labels
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator(program *Program) *CodeGenerator {
	return &CodeGenerator{
		program: program,
		labels:  make(map[string]uint16),
		errors:  []string{},
	}
}

// Generate generates machine code from the AST
func (g *CodeGenerator) Generate() ([]uint16, error) {
	// Pass 1: Calculate addresses and build label table
	if err := g.pass1(); err != nil {
		return nil, err
	}

	// Pass 2: Generate code
	code, err := g.pass2()
	if err != nil {
		return nil, err
	}

	return code, nil
}

// pass1 calculates instruction addresses and builds the label table
func (g *CodeGenerator) pass1() error {
	_, err := g.pass1Statements(g.program.Statements, 0)
	return err
}

// pass1Statements recursively processes statements to calculate addresses
// Returns the next available address after processing all statements
func (g *CodeGenerator) pass1Statements(stmts []Statement, addr uint16) (uint16, error) {
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *LabelStmt:
			// Record label address (points to next instruction)
			if _, exists := g.labels[s.Name]; exists {
				return 0, fmt.Errorf("line %d:%d: duplicate label %q", s.Line, s.Col, s.Name)
			}
			g.labels[s.Name] = addr

		case *InstructionStmt:
			// Each instruction is 2 words (32 bits)
			addr += 2

		case *DefStmt:
			// If def has an initializer, it generates code
			if s.Initializer != nil {
				addr += 2
			}

		case *VarStmt:
			// If var has an initializer, it generates code
			if s.Initializer != nil {
				addr += 2
			}

		case *DataStmt:
			// Data blocks consume space
			addr += uint16(len(s.Values))

		case *StringStmt:
			// Strings are converted to words (1 char per word) + null terminator
			addr += uint16(len(s.Value) + 1)

		case *ScopeStmt:
			// Recursively process scope
			newAddr, err := g.pass1Statements(s.Statements, addr)
			if err != nil {
				return 0, err
			}
			addr = newAddr

		case *FunStmt:
			// Register function as a label: defName.funcName
			funcLabel := fmt.Sprintf("%s.%s", s.DefName, s.Name)
			if _, exists := g.labels[funcLabel]; exists {
				return 0, fmt.Errorf("line %d:%d: duplicate label %q", s.Line, s.Col, funcLabel)
			}
			g.labels[funcLabel] = addr

			// Functions generate their body code
			newAddr, err := g.pass1Statements(s.Body, addr)
			if err != nil {
				return 0, err
			}
			addr = newAddr
		}
	}

	return addr, nil
}

// pass2 generates the actual machine code
func (g *CodeGenerator) pass2() ([]uint16, error) {
	return g.pass2Statements(g.program.Statements, 0)
}

// pass2Statements recursively generates code for statements
// startAddr is the address where this block of statements begins
func (g *CodeGenerator) pass2Statements(stmts []Statement, startAddr uint16) ([]uint16, error) {
	code := []uint16{}
	addr := startAddr

	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *LabelStmt:
			// Labels don't generate code
			continue

		case *InstructionStmt:
			words, err := g.encodeInstructionAtAddr(s, addr)
			if err != nil {
				return nil, fmt.Errorf("line %d:%d: %w", s.Line, s.Col, err)
			}
			code = append(code, words...)
			addr += uint16(len(words))

		case *DefStmt:
			// If def has an initializer, emit it
			if s.Initializer != nil {
				words, err := g.encodeInstructionAtAddr(s.Initializer, addr)
				if err != nil {
					return nil, fmt.Errorf("line %d:%d: %w", s.Line, s.Col, err)
				}
				code = append(code, words...)
				addr += uint16(len(words))
			}

		case *VarStmt:
			// If var has an initializer, emit it
			if s.Initializer != nil {
				words, err := g.encodeInstructionAtAddr(s.Initializer, addr)
				if err != nil {
					return nil, fmt.Errorf("line %d:%d: %w", s.Line, s.Col, err)
				}
				code = append(code, words...)
				addr += uint16(len(words))
			}

		case *DataStmt:
			// Emit data words directly
			code = append(code, s.Values...)
			addr += uint16(len(s.Values))

		case *StringStmt:
			// Convert string to words
			words := g.stringToWords(s.Value)
			code = append(code, words...)
			addr += uint16(len(words))

		case *ScopeStmt:
			// Recursively process scope
			scopeCode, err := g.pass2Statements(s.Statements, addr)
			if err != nil {
				return nil, err
			}
			code = append(code, scopeCode...)
			addr += uint16(len(scopeCode))

		case *FunStmt:
			// Functions are expanded into their body
			funCode, err := g.pass2Statements(s.Body, addr)
			if err != nil {
				return nil, err
			}
			code = append(code, funCode...)
			addr += uint16(len(funCode))

		case *CallStmt:
			// Expand standalone function call
			callCode, err := g.expandStandaloneCallAtAddr(s.Call, addr)
			if err != nil {
				return nil, fmt.Errorf("line %d:%d: %w", s.Line, s.Col, err)
			}
			code = append(code, callCode...)
			addr += uint16(len(callCode))
		}
	}

	return code, nil
}

// expandStandaloneCall expands a standalone function call (no assignment, without PC-relative support)
func (g *CodeGenerator) expandStandaloneCall(fc *FunctionCall) ([]uint16, error) {
	return g.expandStandaloneCallAtAddr(fc, 0)
}

// expandStandaloneCallAtAddr expands a standalone function call (no assignment)
// stack.func(arg1) -> result becomes:
//   1. stack/0 push <- copy val/0 (reserve space for return value)
//   2. stack/0 push <- copy pc/N (push return address)
//   3. stack/0 push <- copy arg1 (push arguments)
//   4. pc <- copy pc/stack.func (jump to function)
func (g *CodeGenerator) expandStandaloneCallAtAddr(fc *FunctionCall, addr uint16) ([]uint16, error) {
	code := []uint16{}
	currentAddr := addr

	// 1. Push space for return value
	pushInst := &InstructionStmt{
		Destination: &Operand{
			Type:     OpMemory,
			Register: cpu.RegR1,
			Offset:   0,
		},
		Push:      true,
		Effect:    cpu.EffectAlways,
		Operation: cpu.OpCopy,
		Source: &Operand{
			Type:   OpImmediate,
			Offset: 0,
		},
	}
	words, err := g.encodeInstructionAtAddr(pushInst, currentAddr)
	if err != nil {
		return nil, err
	}
	code = append(code, words...)
	currentAddr += uint16(len(words))

	// 2. Calculate and push return address
	returnOffset := int16(2 + (len(fc.Args) * 2) + 2)

	returnAddrInst := &InstructionStmt{
		Destination: &Operand{
			Type:     OpMemory,
			Register: cpu.RegR1,
			Offset:   0,
		},
		Push:      true,
		Effect:    cpu.EffectAlways,
		Operation: cpu.OpCopy,
		Source: &Operand{
			Type:     OpPCRelative,
			Register: cpu.RegPC,
			Offset:   returnOffset,
		},
	}
	words, err = g.encodeInstructionAtAddr(returnAddrInst, currentAddr)
	if err != nil {
		return nil, err
	}
	code = append(code, words...)
	currentAddr += uint16(len(words))

	// 3. Push arguments
	for _, arg := range fc.Args {
		pushArgInst := &InstructionStmt{
			Destination: &Operand{
				Type:     OpMemory,
				Register: cpu.RegR1,
				Offset:   0,
			},
			Push:      true,
			Effect:    cpu.EffectAlways,
			Operation: cpu.OpCopy,
			Source:    arg,
		}
		words, err := g.encodeInstructionAtAddr(pushArgInst, currentAddr)
		if err != nil {
			return nil, err
		}
		code = append(code, words...)
		currentAddr += uint16(len(words))
	}

	// 4. Jump to function
	funcLabel := fmt.Sprintf("%s.%s", fc.DefName, fc.FuncName)
	jumpInst := &InstructionStmt{
		Destination: &Operand{
			Type:      OpRegister,
			Register:  cpu.RegPC,
			IsAddress: true,
		},
		Effect:    cpu.EffectAlways,
		Operation: cpu.OpCopy,
		Source: &Operand{
			Type:     OpPCRelative,
			Register: cpu.RegPC,
			LabelRef: funcLabel,
		},
	}
	words, err = g.encodeInstructionAtAddr(jumpInst, currentAddr)
	if err != nil {
		return nil, err
	}
	code = append(code, words...)

	return code, nil
}

// stringToWords converts a string to uCISC data words
// Strings are packed as 1 character per word and null-terminated
func (g *CodeGenerator) stringToWords(s string) []uint16 {
	words := make([]uint16, 0, len(s)+1)

	// Process escape sequences
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '\\' && i+1 < len(s) {
			// Handle escape sequences
			i++
			switch s[i] {
			case 'n':
				words = append(words, uint16('\n'))
			case 't':
				words = append(words, uint16('\t'))
			case 'r':
				words = append(words, uint16('\r'))
			case '\\':
				words = append(words, uint16('\\'))
			case '"':
				words = append(words, uint16('"'))
			default:
				// Unknown escape, keep the character
				words = append(words, uint16(s[i]))
			}
		} else {
			words = append(words, uint16(ch))
		}
	}

	// Add null terminator
	words = append(words, 0)

	return words
}

// encodeInstruction encodes a single instruction into machine code (without PC-relative support)
// For instructions with PC-relative addressing, use encodeInstructionAtAddr instead
func (g *CodeGenerator) encodeInstruction(inst *InstructionStmt) ([]uint16, error) {
	return g.encodeInstructionAtAddr(inst, 0)
}

// encodeInstructionAtAddr encodes a single instruction into machine code
// addr is the address where this instruction will be placed (for PC-relative calculations)
func (g *CodeGenerator) encodeInstructionAtAddr(inst *InstructionStmt, addr uint16) ([]uint16, error) {
	// Check if source is a function call - expand it
	if inst.Source.FunctionCall != nil {
		return g.expandFunctionCallAtAddr(inst, addr)
	}

	// Normal instruction encoding
	src, immediate, err := g.encodeSource(inst.Source, inst.Pop, addr)
	if err != nil {
		return nil, err
	}

	dest, offset, err := g.encodeDestination(inst.Destination, inst.Push)
	if err != nil {
		return nil, err
	}

	// Validate immediate range based on destination type
	if dest.Type == cpu.DestMem {
		// Memory destination: 12-bit immediate (-2048 to 2047)
		if immediate < -2048 || immediate > 2047 {
			return nil, fmt.Errorf("immediate value %d out of range for memory destination (-2048 to 2047)", immediate)
		}
	}
	// Register destination has full 16-bit range (-32768 to 32767), which int16 already enforces

	// Encode using the encoder function
	encoded := cpu.EncodeInstruction(
		src,
		dest,
		inst.Push || inst.Pop,
		inst.Effect,
		inst.Operation,
		immediate,
		offset,
	)

	// Split into two 16-bit words
	word1 := uint16(encoded >> 16)
	word2 := uint16(encoded & 0xFFFF)

	return []uint16{word1, word2}, nil
}

// expandFunctionCall expands a function call (without PC-relative support)
func (g *CodeGenerator) expandFunctionCall(inst *InstructionStmt) ([]uint16, error) {
	return g.expandFunctionCallAtAddr(inst, 0)
}

// expandFunctionCallAtAddr expands a function call into the calling convention
// &dest <- stack.func(arg1) becomes:
//   1. stack/0 push <- copy val/0 (reserve space for return value)
//   2. stack/0 push <- copy pc/N (push return address)
//   3. stack/0 push <- copy arg1 (push arguments)
//   4. pc <- copy pc/stack.func (jump to function)
// That's it. The original destination is ignored - caller must pop result manually.
func (g *CodeGenerator) expandFunctionCallAtAddr(inst *InstructionStmt, addr uint16) ([]uint16, error) {
	fc := inst.Source.FunctionCall
	code := []uint16{}
	currentAddr := addr

	// 1. Push space for return value
	pushInst := &InstructionStmt{
		Destination: &Operand{
			Type:     OpMemory,
			Register: cpu.RegR1,
			Offset:   0,
		},
		Push:      true,
		Effect:    cpu.EffectAlways,
		Operation: cpu.OpCopy,
		Source: &Operand{
			Type:   OpImmediate,
			Offset: 0,
		},
	}
	words, err := g.encodeInstructionAtAddr(pushInst, currentAddr)
	if err != nil {
		return nil, err
	}
	code = append(code, words...)
	currentAddr += uint16(len(words))

	// 2. Calculate and push return address
	returnOffset := int16(2 + (len(fc.Args) * 2) + 2)

	returnAddrInst := &InstructionStmt{
		Destination: &Operand{
			Type:     OpMemory,
			Register: cpu.RegR1,
			Offset:   0,
		},
		Push:      true,
		Effect:    cpu.EffectAlways,
		Operation: cpu.OpCopy,
		Source: &Operand{
			Type:     OpPCRelative,
			Register: cpu.RegPC,
			Offset:   returnOffset,
		},
	}
	words, err = g.encodeInstructionAtAddr(returnAddrInst, currentAddr)
	if err != nil {
		return nil, err
	}
	code = append(code, words...)
	currentAddr += uint16(len(words))

	// 3. Push arguments
	for _, arg := range fc.Args {
		pushArgInst := &InstructionStmt{
			Destination: &Operand{
				Type:     OpMemory,
				Register: cpu.RegR1,
				Offset:   0,
			},
			Push:      true,
			Effect:    cpu.EffectAlways,
			Operation: cpu.OpCopy,
			Source:    arg,
		}
		words, err := g.encodeInstructionAtAddr(pushArgInst, currentAddr)
		if err != nil {
			return nil, err
		}
		code = append(code, words...)
		currentAddr += uint16(len(words))
	}

	// 4. Jump to function
	funcLabel := fmt.Sprintf("%s.%s", fc.DefName, fc.FuncName)
	jumpInst := &InstructionStmt{
		Destination: &Operand{
			Type:      OpRegister,
			Register:  cpu.RegPC,
			IsAddress: true,
		},
		Effect:    cpu.EffectAlways,
		Operation: cpu.OpCopy,
		Source: &Operand{
			Type:     OpPCRelative,
			Register: cpu.RegPC,
			LabelRef: funcLabel,
		},
	}
	words, err = g.encodeInstructionAtAddr(jumpInst, currentAddr)
	if err != nil {
		return nil, err
	}
	code = append(code, words...)

	return code, nil
}

// encodeSource encodes a source operand and returns the immediate value
// currentAddr is the address of the current instruction (for PC-relative calculations)
func (g *CodeGenerator) encodeSource(op *Operand, pop bool, currentAddr uint16) (cpu.Source, int16, error) {
	var src cpu.Source
	immediate := op.Offset

	switch op.Type {
	case OpImmediate:
		// val/immediate or val/label
		src.Type = cpu.SourceVal
		src.Reg = cpu.RegR1 // Doesn't matter for immediate

		// Resolve label if present
		if op.LabelRef != "" {
			addr, ok := g.labels[op.LabelRef]
			if !ok {
				return src, 0, fmt.Errorf("undefined label %q", op.LabelRef)
			}
			immediate = int16(addr)
			// Debug
			_ = addr
		}

		return src, immediate, nil

	case OpRegister:
		// Bare register (e.g., r1) - direct register value
		src.Type = cpu.SourceReg
		src.Reg = op.Register
		return src, 0, nil

	case OpMemory:
		// r1/offset
		src.Type = cpu.SourceMem
		src.Reg = op.Register
		return src, immediate, nil

	case OpPCRelative:
		// pc/offset or pc/label
		if op.LabelRef != "" {
			// Resolve label
			labelAddr, ok := g.labels[op.LabelRef]
			if !ok {
				return src, 0, fmt.Errorf("undefined label %q", op.LabelRef)
			}
			// PC-relative addressing: immediate is offset from current PC
			// The CPU reads this as: value = PC + immediate
			// So we need: immediate = target_address - current_pc
			immediate = int16(labelAddr) - int16(currentAddr)
		}

		// For PC-relative, use SourceReg with PC register
		// The CPU will read this as: value = PC_register + immediate
		src.Type = cpu.SourceReg
		src.Reg = cpu.RegPC
		return src, immediate, nil
	}

	return src, 0, fmt.Errorf("unknown source type")
}

// encodeDestination encodes a destination operand
func (g *CodeGenerator) encodeDestination(op *Operand, push bool) (cpu.Destination, uint8, error) {
	var dest cpu.Destination
	var offset uint8

	// Check if it's a register destination
	// Either: op.IsAddress (has & prefix) OR op.Type == OpRegister (bare register name)
	if op.IsAddress || op.Type == OpRegister {
		// Register destination (&r1 or pc)
		dest.Type = cpu.DestReg
		dest.Reg = op.Register
		return dest, 0, nil
	}

	// Memory destination (r1/offset)
	dest.Type = cpu.DestMem
	dest.Reg = op.Register

	// Handle label references in destination offset
	if op.LabelRef != "" {
		addr, ok := g.labels[op.LabelRef]
		if !ok {
			return dest, 0, fmt.Errorf("undefined label %q", op.LabelRef)
		}
		op.Offset = int16(addr)
	}

	// Offset for memory destination is 4 bits (0-15)
	if op.Offset < 0 || op.Offset > 15 {
		return dest, 0, fmt.Errorf("memory destination offset %d out of range (0-15)", op.Offset)
	}
	offset = uint8(op.Offset)

	return dest, offset, nil
}
