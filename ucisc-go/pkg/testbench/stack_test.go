package testbench

import (
	"testing"
)

// TestStackOperations tests stack push and pop operations
func TestStackOperations(t *testing.T) {
	suite := NewTestSuite("Stack Operations", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Push and pop single value",
		Description: "Push 42 to stack and pop it back",
		Source: `
			&r7 <- copy val/0x1000
			r7/0 push <- copy val/42
			&r1 <- copy r7/0 pop
			debug.assertEqual(&r1, val/42)
			debug.assertEqual(&r7, val/0x1000)  # Stack pointer restored
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Push and pop multiple values",
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

			debug.assertEqual(&r7, val/0x1000)  # Stack pointer restored
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Push register value",
		Description: "Push register contents to stack",
		Source: `
			&r7 <- copy val/0x1000
			&r1 <- copy val/99
			r7/0 push <- copy &r1
			&r2 <- copy r7/0 pop
			debug.assertEqual(&r2, val/99)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Stack pointer decrements",
		Description: "Verify stack pointer decrements on push",
		Source: `
			&r7 <- copy val/0x1000
			debug.assertEqual(&r7, val/0x1000)
			r7/0 push <- copy val/42
			debug.assertEqual(&r7, val/0x0FFF)
			r7/0 push <- copy val/43
			debug.assertEqual(&r7, val/0x0FFE)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Stack pointer increments",
		Description: "Verify stack pointer increments on pop",
		Source: `
			&r7 <- copy val/0x1000
			r7/0 push <- copy val/1
			r7/0 push <- copy val/2
			debug.assertEqual(&r7, val/0x0FFE)
			&r1 <- copy r7/0 pop
			debug.assertEqual(&r7, val/0x0FFF)
			&r1 <- copy r7/0 pop
			debug.assertEqual(&r7, val/0x1000)
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

// TestStackPatterns tests common stack usage patterns
func TestStackPatterns(t *testing.T) {
	suite := NewTestSuite("Stack Usage Patterns", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Save and restore registers",
		Description: "Use stack to preserve register values",
		Source: `
			&r7 <- copy val/0x1000
			&r1 <- copy val/10
			&r2 <- copy val/20
			&r3 <- copy val/30

			# Save registers
			r7/0 push <- copy &r1
			r7/0 push <- copy &r2
			r7/0 push <- copy &r3

			# Modify registers
			&r1 <- copy val/0
			&r2 <- copy val/0
			&r3 <- copy val/0

			# Restore registers (reverse order)
			&r3 <- copy r7/0 pop
			&r2 <- copy r7/0 pop
			&r1 <- copy r7/0 pop

			debug.assertEqual(&r1, val/10)
			debug.assertEqual(&r2, val/20)
			debug.assertEqual(&r3, val/30)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Stack as temporary storage",
		Description: "Use stack to swap register values",
		Source: `
			&r7 <- copy val/0x1000
			&r1 <- copy val/100
			&r2 <- copy val/200

			# Swap r1 and r2 using stack
			r7/0 push <- copy &r1
			r7/0 push <- copy &r2
			&r1 <- copy r7/0 pop   # r1 = old r2
			&r2 <- copy r7/0 pop   # r2 = old r1

			debug.assertEqual(&r1, val/200)
			debug.assertEqual(&r2, val/100)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Deep stack",
		Description: "Push many values to stack",
		Source: `
			&r7 <- copy val/0x1000

			# Push 10 values
			r7/0 push <- copy val/1
			r7/0 push <- copy val/2
			r7/0 push <- copy val/3
			r7/0 push <- copy val/4
			r7/0 push <- copy val/5
			r7/0 push <- copy val/6
			r7/0 push <- copy val/7
			r7/0 push <- copy val/8
			r7/0 push <- copy val/9
			r7/0 push <- copy val/10

			# Pop last value
			&r1 <- copy r7/0 pop
			debug.assertEqual(&r1, val/10)

			# Stack should be 9 deep
			debug.assertEqual(&r7, val/0x0FF7)
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

// TestStackEdgeCases tests edge cases and error conditions
func TestStackEdgeCases(t *testing.T) {
	suite := NewTestSuite("Stack Edge Cases", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Push zero",
		Description: "Push zero value to stack",
		Source: `
			&r7 <- copy val/0x1000
			r7/0 push <- copy val/0
			&r1 <- copy r7/0 pop
			debug.assertEqual(&r1, val/0)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Push maximum value",
		Description: "Push 65535 to stack",
		Source: `
			&r7 <- copy val/0x1000
			r7/0 push <- copy val/65535
			&r1 <- copy r7/0 pop
			debug.assertEqual(&r1, val/65535)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Alternate push and pop",
		Description: "Interleave push and pop operations",
		Source: `
			&r7 <- copy val/0x1000
			r7/0 push <- copy val/1
			&r1 <- copy r7/0 pop
			debug.assertEqual(&r1, val/1)

			r7/0 push <- copy val/2
			&r1 <- copy r7/0 pop
			debug.assertEqual(&r1, val/2)

			r7/0 push <- copy val/3
			&r1 <- copy r7/0 pop
			debug.assertEqual(&r1, val/3)

			debug.assertEqual(&r7, val/0x1000)
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}
