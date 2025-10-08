package asm

import (
	"fmt"

	"github.com/grokthis/ucisc-go/pkg/cpu"
)

// InstructionEncoder encodes AST instructions into machine code
type InstructionEncoder struct {
	labels map[string]uint16 // Label table for resolving references
}

// NewInstructionEncoder creates a new instruction encoder
func NewInstructionEncoder(labels map[string]uint16) *InstructionEncoder {
	return &InstructionEncoder{
		labels: labels,
	}
}

// Encode encodes an instruction to machine code
func (e *InstructionEncoder) Encode(inst *InstructionStmt, addr uint16) ([]uint16, error) {
	// Normal instruction encoding
	src, immediate, err := e.encodeSource(inst.Source, inst.Pop, addr)
	if err != nil {
		return nil, err
	}

	dest, offset, err := e.encodeDestination(inst.Destination, inst.Push)
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

	// Encode using the CPU encoder function
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
// currentAddr is the address of the current instruction (for PC-relative calculations)
func (e *InstructionEncoder) encodeSource(op *Operand, pop bool, currentAddr uint16) (cpu.Source, int16, error) {
	var src cpu.Source
	immediate := op.Offset

	switch op.Type {
	case OpImmediate:
		// val/immediate or val/label
		src.Type = cpu.SourceVal
		src.Reg = cpu.RegR1 // Doesn't matter for immediate

		// Resolve label if present
		if op.LabelRef != "" {
			addr, ok := e.labels[op.LabelRef]
			if !ok {
				return src, 0, fmt.Errorf("undefined label %q", op.LabelRef)
			}
			immediate = int16(addr)
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
			labelAddr, ok := e.labels[op.LabelRef]
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
func (e *InstructionEncoder) encodeDestination(op *Operand, push bool) (cpu.Destination, uint8, error) {
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
		addr, ok := e.labels[op.LabelRef]
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
