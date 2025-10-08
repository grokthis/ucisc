package asm

import "fmt"

// Value range constants for safety checks
const (
	MinInt16  = -32768
	MaxInt16  = 32767
	MaxUint16 = 65535
)

// CheckInt16Range verifies a value fits in int16 range
func CheckInt16Range(value int64, context string) error {
	if value < MinInt16 || value > MaxInt16 {
		return fmt.Errorf("%s: value %d out of int16 range (%d to %d)", context, value, MinInt16, MaxInt16)
	}
	return nil
}

// CheckUint16Range verifies a value fits in uint16 range
func CheckUint16Range(value int64, context string) error {
	if value < 0 || value > MaxUint16 {
		return fmt.Errorf("%s: value %d out of uint16 range (0 to %d)", context, value, MaxUint16)
	}
	return nil
}

// CheckOffsetRange verifies an offset is valid for instruction encoding
func CheckOffsetRange(offset int16) error {
	// Offsets are stored in 12 bits signed (-2048 to 2047)
	if offset < -2048 || offset > 2047 {
		return fmt.Errorf("offset %d out of range (-2048 to 2047)", offset)
	}
	return nil
}

// CheckImmediateRange verifies an immediate value fits in 12 bits
func CheckImmediateRange(value uint16) error {
	// Immediate values are 12 bits unsigned (0 to 4095)
	if value > 4095 {
		return fmt.Errorf("immediate value %d out of range (0 to 4095)", value)
	}
	return nil
}
