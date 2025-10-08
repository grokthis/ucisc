package testbench

import (
	"testing"
	"time"
)

// TestConditionalExecution tests conditional instruction execution
func TestConditionalExecution(t *testing.T) {
	suite := NewTestSuite("Conditional Execution", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Execute on zero flag",
		Description: "Instruction executes when zero flag is set",
		Source: `
			&r2 <- copy val/0
			&r1 <- copy val/0
			&r1 <~ copy val/0      # Set zero flag
			&r2 <0 copy val/100    # Execute only if zero
			debug.assertEqual(&r2, val/100)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Skip on not zero",
		Description: "Instruction skipped when zero flag is not set",
		Source: `
			&r2 <- copy val/0
			&r1 <- copy val/5
			&r1 <~ copy val/5      # Clear zero flag (5 != 0)
			&r2 <0 copy val/100    # Execute only if zero (should skip)
			debug.assertEqual(&r2, val/0)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Execute on not zero flag",
		Description: "Instruction executes when zero flag is not set",
		Source: `
			&r2 <- copy val/0
			&r1 <- copy val/5
			&r1 <~ copy val/5      # Clear zero flag
			&r2 <1 copy val/100    # Execute if NOT zero
			debug.assertEqual(&r2, val/100)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Execute on negative flag",
		Description: "Instruction executes when negative flag is set",
		Source: `
			&flags <- copy val/0x0800  # Signed mode
			&r2 <- copy val/0
			&r1 <- copy val/1
			&r1 <~ sub val/2          # 1-2 = -1, sets negative
			&r2 <n copy val/100       # Execute if negative
			debug.assertEqual(&r2, val/100)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Execute on positive flag",
		Description: "Instruction executes when positive flag is set",
		Source: `
			&flags <- copy val/0x0800  # Signed mode
			&r2 <- copy val/0
			&r1 <- copy val/10
			&r1 <~ sub val/5          # 10-5 = 5, positive
			&r2 <p copy val/100       # Execute if positive
			debug.assertEqual(&r2, val/100)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Execute on overflow flag",
		Description: "Instruction executes when overflow flag is set",
		Source: `
			&flags <- copy val/0x0800  # Signed mode
			&r2 <- copy val/0
			&r1 <- copy val/32767
			&r1 <~ add val/1          # Sets overflow flag
			&r2 <o copy val/100       # Execute if overflow
			debug.assertEqual(&r2, val/100)
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

// TestJumpsAndBranches tests jump and branch instructions
func TestJumpsAndBranches(t *testing.T) {
	suite := NewTestSuite("Jumps and Branches", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Unconditional jump forward",
		Description: "Jump to label forward in code",
		Source: `
			&r1 <- copy val/1
			&pc <- copy pc/skip
			&r1 <- copy val/99  # Should be skipped
			skip:
			debug.assertEqual(&r1, val/1)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Unconditional jump backward",
		Description: "Jump backward to label",
		Source: `
			&r1 <- copy val/0
			&pc <- copy pc/start
			done:
			debug.assertEqual(&r1, val/1)
			pc <- copy &pc
			start:
			&r1 <- copy val/1
			&pc <- copy pc/done
		`,
	})

	suite.AddTest(&Test{
		Name:        "Conditional branch taken",
		Description: "Branch when condition is true",
		Source: `
			&r1 <- copy val/5
			&r1 <~ sub val/5      # Set zero flag
			&r2 <- copy val/0
			&pc <0 copy pc/taken  # Branch if zero
			&r2 <- copy val/99    # Should be skipped
			taken:
			debug.assertEqual(&r2, val/0)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Conditional branch not taken",
		Description: "Don't branch when condition is false",
		Source: `
			&r1 <- copy val/5
			&r1 <~ sub val/3      # Clear zero flag (result != 0)
			&r2 <- copy val/0
			&pc <0 copy pc/skip   # Branch if zero (should not branch)
			&r2 <- copy val/100
			skip:
			debug.assertEqual(&r2, val/100)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "PC-relative addressing",
		Description: "Jump using PC-relative offset",
		Source: `
			&r1 <- copy val/0
			&pc <- copy pc/2      # Skip next instruction
			&r1 <- copy val/99    # Skipped
			&r1 <- copy val/42    # Executed
			debug.assertEqual(&r1, val/42)
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

// TestLoops tests loop constructs
func TestLoops(t *testing.T) {
	suite := NewTestSuite("Loop Constructs", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Simple counting loop",
		Description: "Count from 0 to 5",
		Source: `
			&r1 <- copy val/0
			loop:
			&r1 <- add val/1
			&r1 <~ sub val/5      # Compare with 5
			&pc <1 copy pc/loop   # Continue if not zero
			debug.assertEqual(&r1, val/5)
			pc <- copy &pc
		`,
		Timeout: 1 * time.Second,
	})

	suite.AddTest(&Test{
		Name:        "Loop with counter",
		Description: "Sum numbers 1 to 10",
		Source: `
			&r1 <- copy val/0     # Sum
			&r2 <- copy val/1     # Counter
			loop:
			&r1 <- add &r2
			&r2 <- add val/1
			&r2 <~ sub val/11     # Compare counter with 11
			&pc <1 copy pc/loop   # Continue if not zero
			debug.assertEqual(&r1, val/55)  # 1+2+...+10 = 55
			pc <- copy &pc
		`,
		Timeout: 1 * time.Second,
	})

	suite.AddTest(&Test{
		Name:        "Nested loop",
		Description: "Multiply using nested loops (3 * 4 = 12)",
		Source: `
			&r1 <- copy val/0     # Result
			&r2 <- copy val/3     # Outer counter
			outer:
			&r3 <- copy val/4     # Inner counter
			inner:
			&r1 <- add val/1
			&r3 <- sub val/1
			&r3 <~ sub val/0
			&pc <1 copy pc/inner  # Continue inner if not zero

			&r2 <- sub val/1
			&r2 <~ sub val/0
			&pc <1 copy pc/outer  # Continue outer if not zero

			debug.assertEqual(&r1, val/12)
			pc <- copy &pc
		`,
		Timeout: 1 * time.Second,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

// TestSubroutines tests subroutine call/return patterns
func TestSubroutines(t *testing.T) {
	suite := NewTestSuite("Subroutine Patterns", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Simple subroutine",
		Description: "Call and return from subroutine",
		Source: `
			&r7 <- copy val/0x1000    # Stack pointer
			&r1 <- copy val/5

			# Call subroutine
			r7/0 push <- copy pc/return_addr
			&pc <- copy pc/double
			return_addr:

			debug.assertEqual(&r1, val/10)
			pc <- copy &pc

			# Subroutine: double r1
			double:
			&r1 <- add &r1
			&pc <- copy r7/0 pop
		`,
		Timeout: 1 * time.Second,
	})

	suite.AddTest(&Test{
		Name:        "Nested subroutine calls",
		Description: "Call subroutine from subroutine",
		Source: `
			&r7 <- copy val/0x1000
			&r1 <- copy val/3

			# Call add_five
			r7/0 push <- copy pc/ret1
			&pc <- copy pc/add_five
			ret1:

			debug.assertEqual(&r1, val/13)  # 3 + 5 + 5
			pc <- copy &pc

			# add_five: adds 5 twice by calling add_val
			add_five:
			r7/0 push <- copy pc/ret2
			&pc <- copy pc/add_val
			ret2:
			r7/0 push <- copy pc/ret3
			&pc <- copy pc/add_val
			ret3:
			&pc <- copy r7/0 pop

			# add_val: adds 5 to r1
			add_val:
			&r1 <- add val/5
			&pc <- copy r7/0 pop
		`,
		Timeout: 1 * time.Second,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}
