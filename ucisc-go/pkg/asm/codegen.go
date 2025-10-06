package asm

import (
	"fmt"

	"github.com/grokthis/ucisc-go/pkg/cpu"
)

// CodeGenerator generates machine code from an AST
type CodeGenerator struct {
	program *Program
	labels  map[string]uint16 // label name -> address
	errors  []string
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
	addr := uint16(0)

	for _, stmt := range g.program.Statements {
		switch s := stmt.(type) {
		case *LabelStmt:
			// Record label address (points to next instruction)
			if _, exists := g.labels[s.Name]; exists {
				return fmt.Errorf("line %d:%d: duplicate label %q", s.Line, s.Col, s.Name)
			}
			g.labels[s.Name] = addr

		case *InstructionStmt:
			// Each instruction is 2 words (32 bits)
			addr += 2
		}
	}

	return nil
}

// pass2 generates the actual machine code
func (g *CodeGenerator) pass2() ([]uint16, error) {
	code := []uint16{}

	for _, stmt := range g.program.Statements {
		switch s := stmt.(type) {
		case *LabelStmt:
			// Labels don't generate code
			continue

		case *InstructionStmt:
			words, err := g.encodeInstruction(s)
			if err != nil {
				return nil, fmt.Errorf("line %d:%d: %w", s.Line, s.Col, err)
			}
			code = append(code, words...)
		}
	}

	return code, nil
}

// encodeInstruction encodes a single instruction into machine code
func (g *CodeGenerator) encodeInstruction(inst *InstructionStmt) ([]uint16, error) {
	// Resolve source and destination
	src, immediate, err := g.encodeSource(inst.Source, inst.Pop)
	if err != nil {
		return nil, err
	}

	dest, offset, err := g.encodeDestination(inst.Destination, inst.Push)
	if err != nil {
		return nil, err
	}

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

// encodeSource encodes a source operand and returns the immediate value
func (g *CodeGenerator) encodeSource(op *Operand, pop bool) (cpu.Source, int16, error) {
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
			addr, ok := g.labels[op.LabelRef]
			if !ok {
				return src, 0, fmt.Errorf("undefined label %q", op.LabelRef)
			}
			immediate = int16(addr)
		}

		// For PC-relative, use SourcePC
		src.Type = cpu.SourcePC
		src.Reg = cpu.RegPC
		return src, immediate, nil
	}

	return src, 0, fmt.Errorf("unknown source type")
}

// encodeDestination encodes a destination operand
func (g *CodeGenerator) encodeDestination(op *Operand, push bool) (cpu.Destination, uint8, error) {
	var dest cpu.Destination
	var offset uint8

	if op.IsAddress {
		// Register destination (&r1)
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
