package testbench

import (
	"testing"
)

// TestFlagOperations tests CPU flag behavior
func TestFlagOperations(t *testing.T) {
	suite := NewTestSuite("Flag Operations", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Zero flag on zero result",
		Description: "Zero flag is set when result is 0",
		Source: `
			&r1 <- copy val/5
			&r1 <~ sub val/5       # Flags only, should set zero flag
			debug.assertHasFlags(val/0x0001)  # FlagZero
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Zero flag cleared on non-zero",
		Description: "Zero flag is cleared when result is not 0",
		Source: `
			&r1 <- copy val/5
			&r1 <~ sub val/3       # Flags only, result = 2
			debug.assertNoFlags(val/0x0001)  # FlagZero should not be set
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Negative flag in signed mode",
		Description: "Negative flag is set for negative result in signed mode",
		Source: `
			&flags <- copy val/0x0800  # Set signed mode
			&r1 <- copy val/1
			&r1 <~ sub val/2          # Flags only: 1-2 = -1 (signed)
			debug.assertHasFlags(val/0x0002)  # FlagNegative
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Carry flag on addition overflow",
		Description: "Carry flag is set when addition overflows",
		Source: `
			&r1 <- copy val/65535
			&r1 <~ add val/1          # Flags only
			debug.assertHasFlags(val/0x0004)  # FlagCarry
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Overflow flag in signed mode",
		Description: "Overflow flag is set when signed addition overflows",
		Source: `
			&flags <- copy val/0x0800  # Set signed mode
			&r1 <- copy val/32767
			&r1 <~ add val/1          # Flags only
			debug.assertHasFlags(val/0x0008)  # FlagOverflow
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Multiple flags set",
		Description: "Multiple flags can be set simultaneously",
		Source: `
			&r1 <- copy val/65535
			&r1 <~ add val/1          # Should set Zero and Carry
			debug.assertHasFlags(val/0x0005)  # FlagZero | FlagCarry
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

// TestFlagsOnlyOperations tests operations that only affect flags
func TestFlagsOnlyOperations(t *testing.T) {
	suite := NewTestSuite("Flags-Only Operations", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Flags-only doesn't modify destination",
		Description: "Flags-only operation leaves destination unchanged",
		Source: `
			&r1 <- copy val/10
			&r1 <~ add val/5          # Flags only
			debug.assertEqual(&r1, val/10)  # r1 unchanged
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Comparison using flags-only subtract",
		Description: "Compare two values using flags-only subtract",
		Source: `
			&r1 <- copy val/10
			&r1 <~ sub val/10         # Compare: 10 == 10
			debug.assertHasFlags(val/0x0001)  # Should be zero
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Test bit using flags-only AND",
		Description: "Test if bit is set using flags-only AND",
		Source: `
			&r1 <- copy val/0x00F0
			&r1 <~ and val/0x0010     # Test if bit 4 is set
			debug.assertNoFlags(val/0x0001)  # Result not zero, bit is set
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

// TestSignedUnsignedMode tests signed vs unsigned mode differences
func TestSignedUnsignedMode(t *testing.T) {
	suite := NewTestSuite("Signed vs Unsigned Mode", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Unsigned mode no negative flag",
		Description: "In unsigned mode, negative flag is not set",
		Source: `
			# Unsigned mode (default)
			&r1 <- copy val/1
			&r1 <~ sub val/2          # Result wraps to 65535
			debug.assertNoFlags(val/0x0002)  # FlagNegative should not be set
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Signed mode sets negative flag",
		Description: "In signed mode, negative flag is set for negative results",
		Source: `
			&flags <- copy val/0x0800  # Set signed mode
			&r1 <- copy val/1
			&r1 <~ sub val/2          # Result is -1
			debug.assertHasFlags(val/0x0002)  # FlagNegative should be set
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Unsigned mode no overflow",
		Description: "In unsigned mode, overflow flag is not set",
		Source: `
			# Unsigned mode (default)
			&r1 <- copy val/32767
			&r1 <~ add val/1          # No overflow in unsigned
			debug.assertNoFlags(val/0x0008)  # FlagOverflow should not be set
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Signed overflow positive to negative",
		Description: "Overflow when positive + positive = negative",
		Source: `
			&flags <- copy val/0x0800  # Set signed mode
			&r1 <- copy val/32767     # Max positive signed 16-bit
			&r1 <~ add val/1          # Wraps to -32768
			debug.assertHasFlags(val/0x0008)  # FlagOverflow
			debug.assertHasFlags(val/0x0002)  # FlagNegative
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Signed overflow negative to positive",
		Description: "Overflow when negative + negative = positive",
		Source: `
			&flags <- copy val/0x0800  # Set signed mode
			&r1 <- copy val/-32768    # Min negative signed 16-bit
			&r1 <~ sub val/1          # Wraps to 32767
			debug.assertHasFlags(val/0x0008)  # FlagOverflow
			debug.assertNoFlags(val/0x0002)   # FlagNegative should not be set
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

// TestFlagPersistence tests that flags persist across operations
func TestFlagPersistence(t *testing.T) {
	suite := NewTestSuite("Flag Persistence", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Flags persist until cleared",
		Description: "Flags remain set until explicitly cleared by another operation",
		Source: `
			&r1 <- copy val/0
			&r1 <~ copy val/0         # Set zero flag
			debug.assertHasFlags(val/0x0001)

			# Do a non-flags-only operation that doesn't affect flags
			&r2 <- copy val/42

			# Zero flag should still be set
			debug.assertHasFlags(val/0x0001)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Flags updated by subsequent operation",
		Description: "New operation clears old flags and sets new ones",
		Source: `
			&r1 <- copy val/0
			&r1 <~ copy val/0         # Set zero flag
			debug.assertHasFlags(val/0x0001)

			&r1 <~ copy val/5         # Clear zero flag
			debug.assertNoFlags(val/0x0001)
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}
