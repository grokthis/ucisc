package testbench

import (
	"testing"
	"time"
)

func TestEmulatorExecutor(t *testing.T) {
	suite := NewTestSuite("Basic ALU Operations", NewEmulatorExecutor(65536))

	// Test: Addition
	suite.AddTest(&Test{
		Name:        "Add two numbers",
		Description: "Verify that 5 + 3 = 8",
		Source: `
			&r1 <- copy val/5
			&r1 <- add val/3
			debug.assertEqual(&r1, val/8)
			pc <- copy &pc  # halt
		`,
	})

	// Test: Subtraction
	suite.AddTest(&Test{
		Name:        "Subtract two numbers",
		Description: "Verify that 10 - 3 = 7",
		Source: `
			&r1 <- copy val/10
			&r1 <- sub val/3
			debug.assertEqual(&r1, val/7)
			pc <- copy &pc
		`,
	})

	// Test: Multiplication
	suite.AddTest(&Test{
		Name:        "Multiply two numbers",
		Description: "Verify that 6 * 7 = 42",
		Source: `
			&r1 <- copy val/6
			&r1 <- mult val/7
			debug.assertEqual(&r1, val/42)
			pc <- copy &pc
		`,
	})

	// Test: Bitwise operations
	suite.AddTest(&Test{
		Name:        "Bitwise AND",
		Description: "Verify that 0xFF & 0x0F = 0x0F",
		Source: `
			&r1 <- copy val/0xFF
			&r1 <- and val/0x0F
			debug.assertEqual(&r1, val/0x0F)
			pc <- copy &pc
		`,
	})

	// Test: Memory operations
	suite.AddTest(&Test{
		Name:        "Memory write and read",
		Description: "Write 123 to memory and read it back",
		Source: `
			&r1 <- copy val/0x100  # Use address 0x100
			&r1 <- copy val/123   # Store 123 to memory at r1
			&r2 <- copy &r1        # Read from memory at r1
			debug.assertEqual(&r2, val/123)
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

func TestFlagAssertions(t *testing.T) {
	suite := NewTestSuite("Flag Assertions", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Zero flag set",
		Description: "Verify zero flag is set after result is 0",
		Source: `
			&r1 <- copy val/5
			&r1 <~ sub val/5       # Flags only, should set zero flag
			debug.assertHasFlags(val/0x0001)  # FlagZero
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Negative flag set",
		Description: "Verify negative flag is set for negative result",
		Source: `
			&flags <- copy val/0x0800  # Set signed mode
			&r1 <- copy val/1
			&r1 <~ sub val/2          # Flags only: 1-2 = -1 (signed)
			debug.assertHasFlags(val/0x0002)  # FlagNegative
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

func TestBasicStackOperations(t *testing.T) {
	suite := NewTestSuite("Basic Stack Operations", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Push and Pop",
		Description: "Push value to stack and pop it back",
		Source: `
			&r7 <- copy val/0x1000   # Initialize stack pointer
			r7/0 push <- copy val/42 # Push 42 to stack
			&r1 <- copy r7/0 pop     # Pop from stack to r1
			debug.assertEqual(&r1, val/42)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Multiple push/pop",
		Description: "Push multiple values and pop in reverse order",
		Source: `
			&r7 <- copy val/0x1000

			# Push 10, 20, 30
			r7/0 push <- copy val/10
			r7/0 push <- copy val/20
			r7/0 push <- copy val/30

			# Pop and verify (30, 20, 10)
			&r1 <- copy r7/0 pop
			debug.assertEqual(&r1, val/30)

			&r1 <- copy r7/0 pop
			debug.assertEqual(&r1, val/20)

			&r1 <- copy r7/0 pop
			debug.assertEqual(&r1, val/10)

			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

func TestBasicConditionalExecution(t *testing.T) {
	suite := NewTestSuite("Basic Conditional Execution", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Conditional on zero",
		Description: "Execute instruction only if zero flag is set",
		Source: `
			&r1 <- copy val/0
			&r1 <~ copy val/0      # Set zero flag (flags only)
			&r2 <0 copy val/100    # Execute only if zero
			debug.assertEqual(&r2, val/100)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Conditional on not zero",
		Description: "Skip instruction if zero flag is set",
		Source: `
			&r2 <- copy val/0
			&r1 <- copy val/0
			r1 <~ copy val/0      # Set zero flag
			&r2 <1 copy val/100    # Execute only if NOT zero (should skip)
			debug.assertEqual(&r2, val/0)
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

func TestEdgeCases(t *testing.T) {
	suite := NewTestSuite("Edge Cases", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Overflow detection",
		Description: "32767 + 1 should overflow in signed mode",
		Source: `
			&flags <- copy val/0x0800  # Set signed mode
			&r1 <- copy val/32767
			&r1 <~ add val/1           # Flags only
			debug.assertHasFlags(val/0x0008)  # FlagOverflow
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Maximum value",
		Description: "Handle maximum 16-bit value",
		Source: `
			&r1 <- copy val/65535
			debug.assertEqual(&r1, val/65535)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Negative numbers (signed)",
		Description: "Verify -1 is represented correctly",
		Source: `
			&r1 <- copy val/-1
			debug.assertEqual(&r1, val/65535)  # -1 as uint16
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

func TestTimeout(t *testing.T) {
	suite := NewTestSuite("Timeout Test", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Infinite loop detection",
		Description: "Detect infinite loop via timeout",
		Source: `
			&r1 <- copy val/1
			loop:
			&r1 <- add val/1     # Increment r1
			&pc <- copy pc/loop # Jump back to loop
		`,
		Timeout: 100 * time.Millisecond,
	})

	result := suite.Run()

	// This test SHOULD fail with a timeout
	if !result.HasFailures() {
		t.Error("Expected timeout test to fail, but it passed")
	}

	// Verify it's a timeout error, not an assertion failure
	if result.Results[0].Error == nil {
		t.Error("Expected timeout error, got assertion failure")
	}
}

func TestAssemblyError(t *testing.T) {
	suite := NewTestSuite("Assembly Error Test", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Invalid syntax",
		Description: "Verify assembly errors are caught",
		Source:      `invalid syntax here`,
	})

	result := suite.Run()

	// This test SHOULD fail with assembly error
	if !result.HasFailures() {
		t.Error("Expected assembly error test to fail")
	}

	if result.Results[0].Error == nil {
		t.Error("Expected assembly error")
	}
}
