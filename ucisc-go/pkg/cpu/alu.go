package cpu

// ALU performs arithmetic and logic operations
type ALU struct {
	signed bool // Current signed/unsigned mode
}

// NewALU creates a new ALU
func NewALU() *ALU {
	return &ALU{signed: false}
}

// SetSigned sets the signed mode
func (a *ALU) SetSigned(signed bool) {
	a.signed = signed
}

// Execute performs an ALU operation and returns the result with flags
func (a *ALU) Execute(op Operation, dest, src uint16, carryIn bool) Result {
	switch op {
	case OpCopy:
		return a.copy(src)
	case OpAnd:
		return a.and(dest, src)
	case OpOr:
		return a.or(dest, src)
	case OpXor:
		return a.xor(dest, src)
	case OpInv:
		return a.inv(src)
	case OpShl:
		return a.shl(dest, src)
	case OpShr:
		return a.shr(dest, src)
	case OpSwap:
		return a.swap(src)
	case OpMsb:
		return a.msb(src)
	case OpLsb:
		return a.lsb(src)
	case OpAdd:
		return a.add(dest, src, false)
	case OpAddc:
		return a.add(dest, src, carryIn)
	case OpSub:
		return a.sub(dest, src)
	case OpMult:
		return a.mult(dest, src)
	case OpMultMsw:
		return a.multMsw(dest, src)
	default:
		return Result{Value: 0}
	}
}

// copy operation (no flags set except zero/negative)
func (a *ALU) copy(src uint16) Result {
	return Result{
		Value:    src,
		Zero:     src == 0,
		Negative: a.isNegative(src),
		Carry:    false,
		Overflow: false,
	}
}

// and operation
func (a *ALU) and(dest, src uint16) Result {
	value := dest & src
	return Result{
		Value:    value,
		Zero:     value == 0,
		Negative: a.isNegative(value),
		Carry:    false,
		Overflow: false,
	}
}

// or operation
func (a *ALU) or(dest, src uint16) Result {
	value := dest | src
	return Result{
		Value:    value,
		Zero:     value == 0,
		Negative: a.isNegative(value),
		Carry:    false,
		Overflow: false,
	}
}

// xor operation
func (a *ALU) xor(dest, src uint16) Result {
	value := dest ^ src
	return Result{
		Value:    value,
		Zero:     value == 0,
		Negative: a.isNegative(value),
		Carry:    false,
		Overflow: false,
	}
}

// inv operation (bitwise invert)
func (a *ALU) inv(src uint16) Result {
	value := ^src
	return Result{
		Value:    value,
		Zero:     value == 0,
		Negative: a.isNegative(value),
		Carry:    false,
		Overflow: false,
	}
}

// shl operation (shift left, zero extend)
func (a *ALU) shl(dest, src uint16) Result {
	// Only use lower 4 bits of shift amount (max shift 15)
	shift := src & 0x0F
	value := dest << shift
	carry := shift > 0 && (dest&(1<<(16-shift))) != 0

	return Result{
		Value:    value,
		Zero:     value == 0,
		Negative: a.isNegative(value),
		Carry:    carry,
		Overflow: false,
	}
}

// shr operation (shift right, respect signed mode)
func (a *ALU) shr(dest, src uint16) Result {
	// Only use lower 4 bits of shift amount (max shift 15)
	shift := src & 0x0F
	var value uint16
	var carry bool

	if a.signed && a.isNegative(dest) {
		// Arithmetic right shift (sign extend)
		signedDest := int16(dest)
		value = uint16(signedDest >> shift)
		carry = shift > 0 && (dest&(1<<(shift-1))) != 0
	} else {
		// Logical right shift (zero extend)
		value = dest >> shift
		carry = shift > 0 && (dest&(1<<(shift-1))) != 0
	}

	return Result{
		Value:    value,
		Zero:     value == 0,
		Negative: a.isNegative(value),
		Carry:    carry,
		Overflow: false,
	}
}

// swap operation (swap MSB and LSB bytes)
func (a *ALU) swap(src uint16) Result {
	value := (src << 8) | (src >> 8)
	return Result{
		Value:    value,
		Zero:     value == 0,
		Negative: a.isNegative(value),
		Carry:    false,
		Overflow: false,
	}
}

// msb operation (MSB only: A & 0xFF00)
func (a *ALU) msb(src uint16) Result {
	value := src & 0xFF00
	return Result{
		Value:    value,
		Zero:     value == 0,
		Negative: a.isNegative(value),
		Carry:    false,
		Overflow: false,
	}
}

// lsb operation (LSB only: A & 0x00FF)
func (a *ALU) lsb(src uint16) Result {
	value := src & 0x00FF
	return Result{
		Value:    value,
		Zero:     value == 0,
		Negative: false, // LSB operation always results in positive (MSB is 0)
		Carry:    false,
		Overflow: false,
	}
}

// add operation (with optional carry in)
func (a *ALU) add(dest, src uint16, carryIn bool) Result {
	carry := uint16(0)
	if carryIn {
		carry = 1
	}

	if a.signed {
		// Signed addition
		signedDest := int16(dest)
		signedSrc := int16(src)
		signedCarry := int16(carry)

		result32 := int32(signedDest) + int32(signedSrc) + int32(signedCarry)
		value := uint16(result32)

		// Overflow occurs when signs are same but result sign differs
		overflow := (signedDest >= 0 && signedSrc >= 0 && int16(value) < 0) ||
			(signedDest < 0 && signedSrc < 0 && int16(value) >= 0)

		// Carry in signed mode: unsigned overflow
		unsignedResult := uint32(dest) + uint32(src) + uint32(carry)
		carryOut := unsignedResult > 0xFFFF

		return Result{
			Value:    value,
			Zero:     value == 0,
			Negative: a.isNegative(value),
			Carry:    carryOut,
			Overflow: overflow,
		}
	} else {
		// Unsigned addition
		result32 := uint32(dest) + uint32(src) + uint32(carry)
		value := uint16(result32)
		carryOut := result32 > 0xFFFF

		return Result{
			Value:    value,
			Zero:     value == 0,
			Negative: a.isNegative(value),
			Carry:    carryOut,
			Overflow: carryOut,
		}
	}
}

// sub operation (dest - src)
func (a *ALU) sub(dest, src uint16) Result {
	if a.signed {
		// Signed subtraction
		signedDest := int16(dest)
		signedSrc := int16(src)

		result32 := int32(signedDest) - int32(signedSrc)
		value := uint16(result32)

		// Overflow occurs when signs differ and result sign differs from dest
		overflow := (signedDest >= 0 && signedSrc < 0 && int16(value) < 0) ||
			(signedDest < 0 && signedSrc >= 0 && int16(value) >= 0)

		// Carry/borrow in subtraction: dest < src
		borrow := dest < src

		return Result{
			Value:    value,
			Zero:     value == 0,
			Negative: a.isNegative(value),
			Carry:    borrow,
			Overflow: overflow,
		}
	} else {
		// Unsigned subtraction
		value := dest - src
		borrow := dest < src

		return Result{
			Value:    value,
			Zero:     value == 0,
			Negative: a.isNegative(value),
			Carry:    borrow,
			Overflow: borrow,
		}
	}
}

// mult operation (multiply, return LSW)
func (a *ALU) mult(dest, src uint16) Result {
	if a.signed {
		// Signed multiplication
		signedDest := int16(dest)
		signedSrc := int16(src)

		result32 := int32(signedDest) * int32(signedSrc)
		value := uint16(result32)

		// Overflow if result doesn't fit in 16 bits
		overflow := result32 > 32767 || result32 < -32768

		return Result{
			Value:    value,
			Zero:     value == 0,
			Negative: a.isNegative(value),
			Carry:    false,
			Overflow: overflow,
		}
	} else {
		// Unsigned multiplication
		result32 := uint32(dest) * uint32(src)
		value := uint16(result32)
		overflow := result32 > 0xFFFF

		return Result{
			Value:    value,
			Zero:     value == 0,
			Negative: a.isNegative(value),
			Carry:    false,
			Overflow: overflow,
		}
	}
}

// multMsw operation (multiply, return MSW)
func (a *ALU) multMsw(dest, src uint16) Result {
	if a.signed {
		// Signed multiplication
		signedDest := int16(dest)
		signedSrc := int16(src)

		result32 := int32(signedDest) * int32(signedSrc)
		value := uint16(result32 >> 16)

		// Overflow if MSW is non-zero (doesn't fit in LSW)
		overflow := result32 > 32767 || result32 < -32768

		return Result{
			Value:    value,
			Zero:     value == 0,
			Negative: a.isNegative(value),
			Carry:    false,
			Overflow: overflow,
		}
	} else {
		// Unsigned multiplication
		result32 := uint32(dest) * uint32(src)
		value := uint16(result32 >> 16)
		overflow := value != 0

		return Result{
			Value:    value,
			Zero:     value == 0,
			Negative: a.isNegative(value),
			Carry:    false,
			Overflow: overflow,
		}
	}
}

// isNegative checks if a value is negative (MSB set)
// Only returns true in signed mode
func (a *ALU) isNegative(value uint16) bool {
	if !a.signed {
		return false
	}
	return value&0x8000 != 0
}
