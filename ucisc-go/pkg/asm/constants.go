package asm

// Instruction encoding constants
const (
	// InstructionSizeWords is the size of each uCISC instruction in 16-bit words
	// All instructions are encoded as 32 bits (2 words)
	InstructionSizeWords = 2
)

// Data encoding constants
const (
	// StringNullTerminator is the number of words added for string null terminator
	StringNullTerminator = 1
)
