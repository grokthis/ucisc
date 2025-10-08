package testbench

import (
	"testing"
	"time"
)

// TestALUOperations tests all ALU operations thoroughly
func TestALUOperations(t *testing.T) {
	suite := NewTestSuite("ALU Operations", NewEmulatorExecutor(65536))

	// Addition tests
	suite.AddTest(&Test{
		Name:        "Add positive numbers",
		Description: "5 + 3 = 8",
		Source: `
			&r1 <- copy val/5
			&r1 <- add val/3
			debug.assertEqual(&r1, val/8)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Add with zero",
		Description: "42 + 0 = 42",
		Source: `
			&r1 <- copy val/42
			&r1 <- add val/0
			debug.assertEqual(&r1, val/42)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Add with carry",
		Description: "65535 + 1 = 0 (wraps around)",
		Source: `
			&r1 <- copy val/65535
			&r1 <- add val/1
			debug.assertEqual(&r1, val/0)
			debug.assertHasFlags(val/0x0005)  # Zero + Carry flags
			pc <- copy &pc
		`,
	})

	// Subtraction tests
	suite.AddTest(&Test{
		Name:        "Subtract positive numbers",
		Description: "10 - 3 = 7",
		Source: `
			&r1 <- copy val/10
			&r1 <- sub val/3
			debug.assertEqual(&r1, val/7)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Subtract to zero",
		Description: "5 - 5 = 0",
		Source: `
			&r1 <- copy val/5
			&r1 <- sub val/5
			debug.assertEqual(&r1, val/0)
			debug.assertHasFlags(val/0x0001)  # Zero flag
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Subtract with underflow",
		Description: "0 - 1 = 65535 (wraps around)",
		Source: `
			&r1 <- copy val/0
			&r1 <- sub val/1
			debug.assertEqual(&r1, val/65535)
			pc <- copy &pc
		`,
	})

	// Multiplication tests
	suite.AddTest(&Test{
		Name:        "Multiply small numbers",
		Description: "6 * 7 = 42",
		Source: `
			&r1 <- copy val/6
			&r1 <- mult val/7
			debug.assertEqual(&r1, val/42)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Multiply by zero",
		Description: "100 * 0 = 0",
		Source: `
			&r1 <- copy val/100
			&r1 <- mult val/0
			debug.assertEqual(&r1, val/0)
			debug.assertHasFlags(val/0x0001)  # Zero flag
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Multiply by one",
		Description: "42 * 1 = 42",
		Source: `
			&r1 <- copy val/42
			&r1 <- mult val/1
			debug.assertEqual(&r1, val/42)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Multiply with overflow",
		Description: "256 * 256 = 0 (wraps around)",
		Source: `
			&r1 <- copy val/256
			&r1 <- mult val/256
			debug.assertEqual(&r1, val/0)
			pc <- copy &pc
		`,
	})

	// Bitwise operations
	suite.AddTest(&Test{
		Name:        "Bitwise AND",
		Description: "0xFF & 0x0F = 0x0F",
		Source: `
			&r1 <- copy val/0xFF
			&r1 <- and val/0x0F
			debug.assertEqual(&r1, val/0x0F)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Bitwise AND to zero",
		Description: "0xFF & 0x00 = 0x00",
		Source: `
			&r1 <- copy val/0xFF
			&r1 <- and val/0x00
			debug.assertEqual(&r1, val/0)
			debug.assertHasFlags(val/0x0001)  # Zero flag
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Bitwise OR",
		Description: "0xF0 | 0x0F = 0xFF",
		Source: `
			&r1 <- copy val/0xF0
			&r1 <- or val/0x0F
			debug.assertEqual(&r1, val/0xFF)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Bitwise XOR",
		Description: "0xFF ^ 0xAA = 0x55",
		Source: `
			&r1 <- copy val/0xFF
			&r1 <- xor val/0xAA
			debug.assertEqual(&r1, val/0x55)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Bitwise XOR same value",
		Description: "0x42 ^ 0x42 = 0",
		Source: `
			&r1 <- copy val/0x42
			&r1 <- xor val/0x42
			debug.assertEqual(&r1, val/0)
			debug.assertHasFlags(val/0x0001)  # Zero flag
			pc <- copy &pc
		`,
	})


	// Shift operations
	suite.AddTest(&Test{
		Name:        "Shift left",
		Description: "1 << 4 = 16",
		Source: `
			&r1 <- copy val/1
			&r1 <- shl val/4
			debug.assertEqual(&r1, val/16)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Shift right",
		Description: "16 >> 2 = 4",
		Source: `
			&r1 <- copy val/16
			&r1 <- shr val/2
			debug.assertEqual(&r1, val/4)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Shift left overflow",
		Description: "0x8000 << 1 = 0",
		Source: `
			&r1 <- copy val/0x8000
			&r1 <- shl val/1
			debug.assertEqual(&r1, val/0)
			pc <- copy &pc
		`,
	})

	// Comparison operations (flags only)
	suite.AddTest(&Test{
		Name:        "Compare equal values",
		Description: "5 - 5 sets zero flag",
		Source: `
			&r1 <- copy val/5
			&r1 <~ sub val/5  # Flags only
			debug.assertHasFlags(val/0x0001)  # Zero flag
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Compare less than",
		Description: "3 - 5 in unsigned mode",
		Source: `
			&r1 <- copy val/3
			&r1 <~ sub val/5  # Flags only
			debug.assertNoFlags(val/0x0001)  # Not zero
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

// TestALUWithRegisters tests ALU operations using register operands
func TestALUWithRegisters(t *testing.T) {
	suite := NewTestSuite("ALU with Register Operands", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Add two registers",
		Description: "r1 + r2 = r3",
		Source: `
			&r1 <- copy val/10
			&r2 <- copy val/20
			&r3 <- copy &r1
			&r3 <- add &r2
			debug.assertEqual(&r3, val/30)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Multiply registers",
		Description: "r1 * r2 = r1",
		Source: `
			&r1 <- copy val/6
			&r2 <- copy val/7
			&r1 <- mult &r2
			debug.assertEqual(&r1, val/42)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Bitwise AND registers",
		Description: "r1 & r2 = r3",
		Source: `
			&r1 <- copy val/0xFF00
			&r2 <- copy val/0x0FFF
			&r3 <- copy &r1
			&r3 <- and &r2
			debug.assertEqual(&r3, val/0x0F00)
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

// TestALUSignedMode tests signed arithmetic
func TestALUSignedMode(t *testing.T) {
	suite := NewTestSuite("Signed Arithmetic", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Signed addition positive",
		Description: "Add two positive signed numbers",
		Source: `
			&flags <- copy val/0x0800  # Set signed mode
			&r1 <- copy val/10
			&r1 <- add val/20
			debug.assertEqual(&r1, val/30)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Signed addition negative",
		Description: "Add positive and negative",
		Source: `
			&flags <- copy val/0x0800  # Set signed mode
			&r1 <- copy val/10
			&r1 <- add val/-5
			debug.assertEqual(&r1, val/5)
			pc <- copy &pc
		`,
		Timeout: 5 * time.Second,
	})

	suite.AddTest(&Test{
		Name:        "Signed overflow positive",
		Description: "32767 + 1 overflows",
		Source: `
			&flags <- copy val/0x0800  # Set signed mode
			&r1 <- copy val/32767
			&r1 <~ add val/1  # Flags only
			debug.assertHasFlags(val/0x0008)  # Overflow flag
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Signed overflow negative",
		Description: "-32768 - 1 overflows",
		Source: `
			&flags <- copy val/0x0800  # Set signed mode
			&r1 <- copy val/-32768
			&r1 <~ sub val/1  # Flags only
			debug.assertHasFlags(val/0x0008)  # Overflow flag
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Signed negative flag",
		Description: "Negative result sets negative flag",
		Source: `
			&flags <- copy val/0x0800  # Set signed mode
			&r1 <- copy val/5
			&r1 <~ sub val/10  # Flags only
			debug.assertHasFlags(val/0x0002)  # Negative flag
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}
