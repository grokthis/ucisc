package cpu

// Instruction format constants
const (
	// InstructionSizeWords is the size of each uCISC instruction in 16-bit words
	// All instructions are encoded as 32 bits (2 words)
	InstructionSizeWords = 2
)

// Stack operation constants
const (
	// PopBonusIncrement is the free increment on pop operations
	// "You always get 1 for free" - pop increments by immediate+1
	PopBonusIncrement = 1

	// PushDecrement is the amount to decrement the stack pointer before push
	PushDecrement = 1
)
