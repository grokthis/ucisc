package testbench

import (
	"testing"
)

// TestMemoryOperations tests memory read/write operations
func TestMemoryOperations(t *testing.T) {
	suite := NewTestSuite("Memory Operations", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Write and read memory",
		Description: "Store value to memory and read it back",
		Source: `
			&r1 <- copy val/0x100
			r1/0 <- copy val/42
			&r2 <- copy r1/0
			debug.assertEqual(&r2, val/42)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Memory with offset",
		Description: "Write to memory with offset",
		Source: `
			&r1 <- copy val/0x100
			r1/10 <- copy val/99
			&r2 <- copy r1/10
			debug.assertEqual(&r2, val/99)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Copy between memory locations",
		Description: "Copy value from one memory location to another",
		Source: `
			&r1 <- copy val/0x100
			&r2 <- copy val/0x200
			r1 <- copy val/123
			r2 <- copy r1
			&r3 <- copy r2
			debug.assertEqual(&r3, val/123)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Multiple writes to same location",
		Description: "Overwrite memory location",
		Source: `
			&r1 <- copy val/0x100
			r1 <- copy val/10
			r1 <- copy val/20
			r1 <- copy val/30
			&r2 <- copy r1
			debug.assertEqual(&r2, val/30)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Memory at address zero",
		Description: "Read/write at address 0",
		Source: `
			&r1 <- copy val/0
			r1 <- copy val/77
			&r2 <- copy r1
			debug.assertEqual(&r2, val/77)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Sequential memory writes",
		Description: "Write to consecutive memory locations",
		Source: `
			&r1 <- copy val/0x100
			r1/0 <- copy val/10
			r1/1 <- copy val/20
			r1/2 <- copy val/30

			&r2 <- copy r1/0
			debug.assertEqual(&r2, val/10)
			&r2 <- copy r1/1
			debug.assertEqual(&r2, val/20)
			&r2 <- copy r1/2
			debug.assertEqual(&r2, val/30)
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

// TestIndirectAddressing tests indirect memory addressing
func TestIndirectAddressing(t *testing.T) {
	suite := NewTestSuite("Indirect Addressing", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "Indirect read",
		Description: "Read memory through pointer",
		Source: `
			&r1 <- copy val/0x100
			&r2 <- copy val/0x200

			# Store pointer value (0x200) at memory location r1
			r1/0 <- copy val/0x200

			# Write value 42 to memory at address 0x200
			r2/0 <- copy val/42

			# Read pointer from memory at r1, store in r3
			&r3 <- copy r1/0

			# Now r3 = 0x200, read from memory at r3
			&r5 <- copy r3/0
			debug.assertEqual(&r5, val/42)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Indirect write",
		Description: "Write memory through pointer",
		Source: `
			&r1 <- copy val/0x100
			&r2 <- copy val/0x200

			# Store pointer value (0x200) at memory location r1
			r1/0 <- copy val/0x200

			# Read pointer from memory at r1 into r3
			&r3 <- copy r1/0

			# Write value through pointer (r3 now contains 0x200)
			r3/0 <- copy val/99

			# Read directly from memory at 0x200 using r2
			&r5 <- copy r2/0
			debug.assertEqual(&r5, val/99)
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}

// TestMemoryBoundaries tests edge cases in memory addressing
func TestMemoryBoundaries(t *testing.T) {
	suite := NewTestSuite("Memory Boundaries", NewEmulatorExecutor(65536))

	suite.AddTest(&Test{
		Name:        "High memory address",
		Description: "Write to high memory address",
		Source: `
			&r1 <- copy val/0xFFF0
			r1 <- copy val/123
			&r2 <- copy r1
			debug.assertEqual(&r2, val/123)
			pc <- copy &pc
		`,
	})

	suite.AddTest(&Test{
		Name:        "Maximum offset",
		Description: "Use maximum offset value (15)",
		Source: `
			&r1 <- copy val/0x1000
			r1/15 <- copy val/42
			&r2 <- copy r1/15
			debug.assertEqual(&r2, val/42)
			pc <- copy &pc
		`,
	})

	result := suite.Run()
	if result.HasFailures() {
		result.Print()
		t.Errorf("Test suite failed: %d/%d tests passed", result.Passed, len(result.Results))
	}
}
