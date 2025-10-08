package cpu

// Register represents a CPU register
type Register uint8

const (
	RegPC      Register = 0 // Program counter
	RegR1      Register = 1 // General purpose register 1 (stack pointer)
	RegR2      Register = 2 // General purpose register 2
	RegR3      Register = 3 // General purpose register 3
	RegBanking Register = 4 // Banking control register
	RegR5      Register = 5 // General purpose register 5
	RegR6      Register = 6 // General purpose register 6
	RegR7      Register = 7 // General purpose register 7
	RegFlags   Register = 8 // Flags register
	RegInt     Register = 9 // Interrupt register
)

// Operation represents an ALU operation
type Operation uint8

const (
	OpCopy    Operation = 0x0 // Copy
	OpAnd     Operation = 0x1 // Bitwise AND
	OpOr      Operation = 0x2 // Bitwise OR
	OpXor     Operation = 0x3 // Bitwise XOR
	OpInv     Operation = 0x4 // Bitwise invert
	OpShl     Operation = 0x5 // Shift left
	OpShr     Operation = 0x6 // Shift right (respects signed mode)
	OpSwap    Operation = 0x7 // Swap MSB and LSB
	OpMsb     Operation = 0x8 // MSB only (A & 0xFF00)
	OpLsb     Operation = 0x9 // LSB only (A & 0x00FF)
	OpAdd     Operation = 0xA // Add
	OpSub     Operation = 0xB // Subtract
	OpMult    Operation = 0xC // Multiply (LSW)
	OpMultMsw Operation = 0xD // Multiply (MSW)
	OpAddc    Operation = 0xE // Add with carry
)

// Effect represents conditional execution
type Effect uint8

const (
	EffectZero      Effect = 0 // Store if zero
	EffectNotZero   Effect = 1 // Store if not zero
	EffectNegative  Effect = 2 // Store if negative
	EffectPositive  Effect = 3 // Store if positive
	EffectAlways    Effect = 4 // Store always
	EffectOverflow  Effect = 5 // Store if overflow
	EffectInterrupt Effect = 6 // Store if interrupt
	EffectFlagsOnly Effect = 7 // Set flags only, don't store
)

// SourceType represents the type of source operand
type SourceType uint8

const (
	SourcePC  SourceType = 0 // PC register value
	SourceMem SourceType = 1 // Memory at register address
	SourceVal SourceType = 2 // Immediate value
	SourceReg SourceType = 3 // Register value
)

// DestType represents the type of destination operand
type DestType uint8

const (
	DestReg DestType = 0 // Register
	DestMem DestType = 1 // Memory at register address
)

// Flags represents processor flags
type Flags uint16

const (
	FlagZero     Flags = 1 << 0  // Zero flag
	FlagNegative Flags = 1 << 1  // Negative flag
	FlagCarry    Flags = 1 << 2  // Carry flag
	FlagOverflow Flags = 1 << 3  // Overflow flag
	FlagInterrupt Flags = 1 << 8  // Interrupt flag
	FlagSigned   Flags = 1 << 11 // Signed mode flag
	FlagBlocked  Flags = 1 << 13 // Blocked memory writable
	FlagResume   Flags = 1 << 14 // Resume on interrupt
	FlagHalt     Flags = 1 << 15 // Halt flag
)

// Source represents an instruction source operand
type Source struct {
	Type SourceType
	Reg  Register
}

// Destination represents an instruction destination operand
type Destination struct {
	Type DestType
	Reg  Register
}

// Instruction represents a decoded uCISC instruction
type Instruction struct {
	Raw         uint32      // Raw instruction bytes
	Source      Source      // Source operand
	Destination Destination // Destination operand
	Increment   bool        // Push/pop flag
	Effect      Effect      // Conditional effect
	Operation   Operation   // ALU operation
	Immediate   int16       // Immediate value (signed)
	Offset      uint8       // Offset for memory destinations (0-15)
}

// Result contains the result of an ALU operation
type Result struct {
	Value    uint16
	Zero     bool
	Negative bool
	Carry    bool
	Overflow bool
}

// String returns the register name
func (r Register) String() string {
	switch r {
	case RegPC:
		return "pc"
	case RegR1:
		return "r1"
	case RegR2:
		return "r2"
	case RegR3:
		return "r3"
	case RegBanking:
		return "banking"
	case RegR5:
		return "r5"
	case RegR6:
		return "r6"
	case RegR7:
		return "r7"
	case RegFlags:
		return "flags"
	case RegInt:
		return "interrupt"
	default:
		return "unknown"
	}
}

// String returns the operation name
func (o Operation) String() string {
	switch o {
	case OpCopy:
		return "copy"
	case OpAnd:
		return "and"
	case OpOr:
		return "or"
	case OpXor:
		return "xor"
	case OpInv:
		return "inv"
	case OpShl:
		return "shl"
	case OpShr:
		return "shr"
	case OpSwap:
		return "swap"
	case OpMsb:
		return "msb"
	case OpLsb:
		return "lsb"
	case OpAdd:
		return "add"
	case OpSub:
		return "sub"
	case OpMult:
		return "mult"
	case OpMultMsw:
		return "multsw"
	case OpAddc:
		return "addc"
	default:
		return "unknown"
	}
}

// IsSet checks if a flag is set
func (f Flags) IsSet(flag Flags) bool {
	return f&flag != 0
}

// Set sets a flag
func (f Flags) Set(flag Flags) Flags {
	return f | flag
}

// Clear clears a flag
func (f Flags) Clear(flag Flags) Flags {
	return f &^ flag
}

// Zero returns the zero flag
func (f Flags) Zero() bool {
	return f.IsSet(FlagZero)
}

// Negative returns the negative flag
func (f Flags) Negative() bool {
	return f.IsSet(FlagNegative)
}

// Carry returns the carry flag
func (f Flags) Carry() bool {
	return f.IsSet(FlagCarry)
}

// Overflow returns the overflow flag
func (f Flags) Overflow() bool {
	return f.IsSet(FlagOverflow)
}

// Interrupt returns the interrupt flag
func (f Flags) Interrupt() bool {
	return f.IsSet(FlagInterrupt)
}

// Signed returns the signed mode flag
func (f Flags) Signed() bool {
	return f.IsSet(FlagSigned)
}

// Halted returns the halt flag
func (f Flags) Halted() bool {
	return f.IsSet(FlagHalt)
}
