package cpu

import "fmt"

// =============================================================================
// Instruction Encoding
// =============================================================================

// EncodeInstruction encodes an instruction from its components
func EncodeInstruction(src Source, dest Destination, inc bool, eff Effect, op Operation, imm int16, off uint8) uint32 {
	word1 := uint16(0)
	word2 := uint16(0)

	// Encode source (top 4 bits)
	srcBits := encodeSourceBits(src)
	word1 |= uint16(srcBits) << 12

	// Encode destination (next 4 bits)
	destBits := encodeDestinationBits(dest)
	word1 |= uint16(destBits) << 8

	// Encode increment (bit 7)
	if inc {
		word1 |= 0x0080
	}

	// Encode effect (bits 4-6)
	word1 |= uint16(eff) << 4

	// Encode operation (bits 0-3)
	word1 |= uint16(op)

	// Encode immediate/offset
	if dest.Type == DestMem {
		// Memory destination: 4-bit offset, 12-bit immediate
		word2 = uint16(off)<<12 | (uint16(imm) & 0x0FFF)
	} else {
		// Register destination: 16-bit immediate
		word2 = uint16(imm)
	}

	return uint32(word1)<<16 | uint32(word2)
}

func encodeSourceBits(src Source) uint8 {
	switch src.Type {
	case SourceReg:
		switch src.Reg {
		case RegPC:
			return 0
		case RegR1:
			return 5
		case RegR2:
			return 6
		case RegR3:
			return 7
		case RegFlags:
			return 8
		case RegR5:
			return 13
		case RegR6:
			return 14
		case RegR7:
			return 15
		case RegInt:
			return 12
		}
	case SourceMem:
		switch src.Reg {
		case RegR1:
			return 1
		case RegR2:
			return 2
		case RegR3:
			return 3
		case RegR5:
			return 9
		case RegR6:
			return 10
		case RegR7:
			return 11
		}
	case SourceVal:
		return 4
	}
	return 0
}

func encodeDestinationBits(dest Destination) uint8 {
	switch dest.Type {
	case DestReg:
		switch dest.Reg {
		case RegPC:
			return 0
		case RegR1:
			return 5
		case RegR2:
			return 6
		case RegR3:
			return 7
		case RegBanking:
			return 4
		case RegFlags:
			return 8
		case RegR5:
			return 13
		case RegR6:
			return 14
		case RegR7:
			return 15
		case RegInt:
			return 12
		}
	case DestMem:
		switch dest.Reg {
		case RegR1:
			return 1
		case RegR2:
			return 2
		case RegR3:
			return 3
		case RegR5:
			return 9
		case RegR6:
			return 10
		case RegR7:
			return 11
		}
	}
	return 0
}

// =============================================================================
// Instruction Decoding
// =============================================================================

// Decode decodes a 32-bit instruction into its components
func Decode(raw uint32) (*Instruction, error) {
	// Extract bit fields from instruction
	// Format: SSSSDDDD MEEEAAAA IIIIIIII IIIIIIII (non-mem dest)
	//         SSSSDDDD MEEEAAAA OOOOIIII IIIIIIII (mem dest)
	// Word 1: SSSSDDDD MEEEAAAA (high 16 bits)
	// Word 2: IIIIIIII IIIIIIII or OOOOIIII IIIIIIII (low 16 bits)

	word1 := uint16(raw >> 16)
	word2 := uint16(raw & 0xFFFF)

	srcBits := uint8((word1 >> 12) & 0xF)        // Top 4 bits of word1
	destBits := uint8((word1 >> 8) & 0xF)        // Next 4 bits of word1
	increment := (word1 & 0x0080) != 0           // Bit 7 of word1 (M bit)
	effectBits := uint8((word1 >> 4) & 0x7)      // Bits 4-6 of word1 (EEE)
	opBits := uint8(word1 & 0xF)                 // Bottom 4 bits of word1 (AAAA)

	// Decode source
	src, err := decodeSource(srcBits)
	if err != nil {
		return nil, fmt.Errorf("invalid source: %w", err)
	}

	// Decode destination
	dest, err := decodeDestination(destBits)
	if err != nil {
		return nil, fmt.Errorf("invalid destination: %w", err)
	}

	// Decode effect
	effect := Effect(effectBits)
	if effect > EffectFlagsOnly {
		return nil, fmt.Errorf("invalid effect: %d", effectBits)
	}

	// Decode operation
	operation := Operation(opBits)
	if operation > OpAddc {
		return nil, fmt.Errorf("invalid operation: %d", opBits)
	}

	// Decode immediate and offset
	var immediate int16
	var offset uint8

	if dest.Type == DestMem {
		// Memory destination: 4-bit offset, 12-bit immediate
		offset = uint8((word2 >> 12) & 0xF)
		imm12 := word2 & 0x0FFF
		// Sign extend 12-bit to 16-bit
		if imm12&0x0800 != 0 {
			immediate = int16(imm12 | 0xF000)
		} else {
			immediate = int16(imm12)
		}
	} else {
		// Register destination: no offset, 16-bit immediate
		offset = 0
		immediate = int16(word2)
	}

	inst := &Instruction{
		Raw:         raw,
		Source:      src,
		Destination: dest,
		Increment:   increment,
		Effect:      effect,
		Operation:   operation,
		Immediate:   immediate,
		Offset:      offset,
	}

	return inst, nil
}

// decodeSource decodes the source field bits
func decodeSource(bits uint8) (Source, error) {
	switch bits {
	case 0: // 0000 - PC register
		return Source{Type: SourceReg, Reg: RegPC}, nil
	case 1: // 0001 - r1 memory
		return Source{Type: SourceMem, Reg: RegR1}, nil
	case 2: // 0010 - r2 memory
		return Source{Type: SourceMem, Reg: RegR2}, nil
	case 3: // 0011 - r3 memory
		return Source{Type: SourceMem, Reg: RegR3}, nil
	case 4: // 0100 - immediate value
		return Source{Type: SourceVal, Reg: 0}, nil
	case 5: // 0101 - r1 register
		return Source{Type: SourceReg, Reg: RegR1}, nil
	case 6: // 0110 - r2 register
		return Source{Type: SourceReg, Reg: RegR2}, nil
	case 7: // 0111 - r3 register
		return Source{Type: SourceReg, Reg: RegR3}, nil
	case 8: // 1000 - flags register
		return Source{Type: SourceReg, Reg: RegFlags}, nil
	case 9: // 1001 - r5 memory
		return Source{Type: SourceMem, Reg: RegR5}, nil
	case 10: // 1010 - r6 memory
		return Source{Type: SourceMem, Reg: RegR6}, nil
	case 11: // 1011 - r7 memory
		return Source{Type: SourceMem, Reg: RegR7}, nil
	case 12: // 1100 - interrupt register
		return Source{Type: SourceReg, Reg: RegInt}, nil
	case 13: // 1101 - r5 register
		return Source{Type: SourceReg, Reg: RegR5}, nil
	case 14: // 1110 - r6 register
		return Source{Type: SourceReg, Reg: RegR6}, nil
	case 15: // 1111 - r7 register
		return Source{Type: SourceReg, Reg: RegR7}, nil
	default:
		return Source{}, fmt.Errorf("invalid source bits: %d", bits)
	}
}

// decodeDestination decodes the destination field bits
func decodeDestination(bits uint8) (Destination, error) {
	switch bits {
	case 0: // 0000 - PC register
		return Destination{Type: DestReg, Reg: RegPC}, nil
	case 1: // 0001 - r1 memory
		return Destination{Type: DestMem, Reg: RegR1}, nil
	case 2: // 0010 - r2 memory
		return Destination{Type: DestMem, Reg: RegR2}, nil
	case 3: // 0011 - r3 memory
		return Destination{Type: DestMem, Reg: RegR3}, nil
	case 4: // 0100 - banking register
		return Destination{Type: DestReg, Reg: RegBanking}, nil
	case 5: // 0101 - r1 register
		return Destination{Type: DestReg, Reg: RegR1}, nil
	case 6: // 0110 - r2 register
		return Destination{Type: DestReg, Reg: RegR2}, nil
	case 7: // 0111 - r3 register
		return Destination{Type: DestReg, Reg: RegR3}, nil
	case 8: // 1000 - flags register
		return Destination{Type: DestReg, Reg: RegFlags}, nil
	case 9: // 1001 - r5 memory
		return Destination{Type: DestMem, Reg: RegR5}, nil
	case 10: // 1010 - r6 memory
		return Destination{Type: DestMem, Reg: RegR6}, nil
	case 11: // 1011 - r7 memory
		return Destination{Type: DestMem, Reg: RegR7}, nil
	case 12: // 1100 - interrupt register
		return Destination{Type: DestReg, Reg: RegInt}, nil
	case 13: // 1101 - r5 register
		return Destination{Type: DestReg, Reg: RegR5}, nil
	case 14: // 1110 - r6 register
		return Destination{Type: DestReg, Reg: RegR6}, nil
	case 15: // 1111 - r7 register
		return Destination{Type: DestReg, Reg: RegR7}, nil
	default:
		return Destination{}, fmt.Errorf("invalid destination bits: %d", bits)
	}
}

// =============================================================================
// Instruction Helpers
// =============================================================================

// IsHalt checks if this instruction is a HALT instruction
// HALT is: copy pc to pc with immediate 0 and effect always (0x0000 0x0400 0x0000 0x0000)
func (i *Instruction) IsHalt() bool {
	// HALT is: pc <- copy pc (register source)
	// OR: pc <- copy pc/0 (memory source with immediate 0)
	if i.Operation != OpCopy ||
		i.Destination.Type != DestReg ||
		i.Destination.Reg != RegPC ||
		i.Effect != EffectAlways {
		return false
	}

	// Check for register source (pc <- copy pc)
	if i.Source.Type == SourceReg && i.Source.Reg == RegPC && i.Immediate == 0 {
		return true
	}

	// Check for memory source with PC and immediate 0 (pc <- copy pc/0)
	if i.Source.Type == SourceMem && i.Source.Reg == RegPC && i.Immediate == 0 {
		return true
	}

	return false
}
