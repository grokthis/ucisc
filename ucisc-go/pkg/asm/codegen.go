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
	funcCallCount int                  // Counter for generating unique labels
	encoder       *InstructionEncoder // Instruction encoder
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

	// Create encoder with label table
	g.encoder = NewInstructionEncoder(g.labels)

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

// alignTo2Words aligns address to 2-word boundary by adding 1 if odd
func alignTo2Words(addr uint16) uint16 {
	if addr%2 == 1 {
		return addr + 1
	}
	return addr
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
			// Align to 2-word boundary before instruction
			addr = alignTo2Words(addr)
			// Each instruction is 2 words (32 bits)
			addr += InstructionSizeWords

		case *DefStmt:
			// If def has an initializer, it generates code
			if s.Initializer != nil {
				// Align to 2-word boundary before instruction
				addr = alignTo2Words(addr)
				addr += InstructionSizeWords
			}

		case *VarStmt:
			// If var has an initializer, it generates code
			if s.Initializer != nil {
				// Align to 2-word boundary before instruction
				addr = alignTo2Words(addr)
				addr += InstructionSizeWords
			}

		case *DataStmt:
			// Data blocks consume space
			addr += uint16(len(s.Values))

		case *StringStmt:
			// Strings are converted to words (1 char per word) + null terminator
			addr += uint16(len(s.Value) + StringNullTerminator)

		case *ScopeStmt:
			// Recursively process scope
			newAddr, err := g.pass1Statements(s.Statements, addr)
			if err != nil {
				return 0, err
			}
			addr = newAddr

		case *FunStmt:
			// Align to 2-word boundary before function (which contains instructions)
			addr = alignTo2Words(addr)

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
			// Add padding word if needed for alignment
			if addr%2 == 1 {
				code = append(code, 0x0000)
				addr++
			}
			words, err := g.encodeInstructionAtAddr(s, addr)
			if err != nil {
				return nil, fmt.Errorf("line %d:%d: %w", s.Line, s.Col, err)
			}
			code = append(code, words...)
			addr += uint16(len(words))

		case *DefStmt:
			// If def has an initializer, emit it
			if s.Initializer != nil {
				// Add padding word if needed for alignment
				if addr%2 == 1 {
					code = append(code, 0x0000)
					addr++
				}
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
				// Add padding word if needed for alignment
				if addr%2 == 1 {
					code = append(code, 0x0000)
					addr++
				}
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
			// Add padding word if needed for alignment
			if addr%2 == 1 {
				code = append(code, 0x0000)
				addr++
			}
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

// buildFunctionCallInstructions generates the instruction sequence for a function call
// This implements the uCISC calling convention:
//   1. Push space for return value
//   2. Push return address
//   3. Push arguments
//   4. Jump to function
func (g *CodeGenerator) buildFunctionCallInstructions(fc *FunctionCall, addr uint16) ([]uint16, error) {
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

// expandStandaloneCall expands a standalone function call (no assignment, without PC-relative support)
func (g *CodeGenerator) expandStandaloneCall(fc *FunctionCall) ([]uint16, error) {
	return g.expandStandaloneCallAtAddr(fc, 0)
}

// expandStandaloneCallAtAddr expands a standalone function call (no assignment)
func (g *CodeGenerator) expandStandaloneCallAtAddr(fc *FunctionCall, addr uint16) ([]uint16, error) {
	return g.buildFunctionCallInstructions(fc, addr)
}

// stringToWords converts a string to uCISC data words
// Strings are packed as 1 character per word and null-terminated
func (g *CodeGenerator) stringToWords(s string) []uint16 {
	words := make([]uint16, 0, len(s)+StringNullTerminator)

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

	// Use the instruction encoder
	return g.encoder.Encode(inst, addr)
}

// expandFunctionCall expands a function call (without PC-relative support)
func (g *CodeGenerator) expandFunctionCall(inst *InstructionStmt) ([]uint16, error) {
	return g.expandFunctionCallAtAddr(inst, 0)
}

// expandFunctionCallAtAddr expands a function call into the calling convention
// &dest <- stack.func(arg1) becomes the same instruction sequence as a standalone call.
// The original destination is ignored - caller must pop result manually.
func (g *CodeGenerator) expandFunctionCallAtAddr(inst *InstructionStmt, addr uint16) ([]uint16, error) {
	return g.buildFunctionCallInstructions(inst.Source.FunctionCall, addr)
}
