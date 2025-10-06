package cpu

// Helper functions for encoding instructions (useful for testing)

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
