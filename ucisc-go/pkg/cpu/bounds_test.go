package cpu

import (
	"strings"
	"testing"
)

// TestStrictBounds_PCOutOfBounds verifies strict bounds checking for PC
func TestStrictBounds_PCOutOfBounds(t *testing.T) {
	cpu := New(10) // Small memory: 10 words
	cpu.SetStrictBounds(true)

	// Set PC beyond memory
	cpu.SetPC(15)

	// Should error in strict mode
	err := cpu.Step()
	if err == nil {
		t.Fatal("Expected error for PC out of bounds in strict mode")
	}

	if !strings.Contains(err.Error(), "out of bounds") {
		t.Errorf("Expected 'out of bounds' error, got: %v", err)
	}
}

// TestNonStrictBounds_PCWraps verifies wrap-around in non-strict mode
func TestNonStrictBounds_PCWraps(t *testing.T) {
	cpu := New(10) // Small memory: 10 words
	cpu.SetStrictBounds(false)

	// Load a HALT instruction at wrapped position
	// PC=15, memSize=10, so 15%10=5, (15+1)%10=6
	// HALT: pc <- copy pc/0 encodes as 0x0040 0x0000
	cpu.memory.WriteLocal(5, 0x0040) // HALT instruction word 1
	cpu.memory.WriteLocal(6, 0x0000) // HALT instruction word 2

	// Let's verify what we wrote
	w1 := cpu.memory.ReadLocal(5)
	w2 := cpu.memory.ReadLocal(6)
	t.Logf("Wrote HALT at 5-6: %04X %04X", w1, w2)

	cpu.SetPC(15)

	// Should wrap and fetch from address 5-6 in non-strict mode
	err := cpu.Step()
	if err != nil {
		t.Fatalf("Unexpected error in non-strict mode: %v", err)
	}

	// Debug: check what was fetched
	t.Logf("PC was %d, memory size %d", 15, cpu.memory.Size())
	t.Logf("Should fetch from addresses %d and %d", 15%10, (15+1)%10)
	t.Logf("CPU halted: %v", cpu.IsHalted())

	// Should have halted (fetched HALT from wrapped address)
	if !cpu.IsHalted() {
		t.Error("Expected CPU to halt after fetching wrapped HALT instruction")
	}
}

// TestStrictBounds_EdgeCase verifies boundary at memory edge
func TestStrictBounds_EdgeCase(t *testing.T) {
	cpu := New(10)
	cpu.SetStrictBounds(true)

	// PC at last valid position for 2-word fetch
	cpu.SetPC(8) // Can fetch words 8 and 9

	// Load HALT
	cpu.memory.WriteLocal(8, 0x0040)
	cpu.memory.WriteLocal(9, 0x0000)

	// Should succeed
	err := cpu.Step()
	if err != nil {
		t.Errorf("Unexpected error at valid boundary: %v", err)
	}

	// Reset and try PC=9 (would need to fetch word 10, which doesn't exist)
	cpu2 := New(10)
	cpu2.SetStrictBounds(true)
	cpu2.SetPC(9)

	err = cpu2.Step()
	if err == nil {
		t.Error("Expected error when PC+1 exceeds memory")
	}
}

// TestDefaultStrictMode verifies strict mode is default
func TestDefaultStrictMode(t *testing.T) {
	cpu := New(10)

	// Should be in strict mode by default
	cpu.SetPC(15)

	err := cpu.Step()
	if err == nil {
		t.Error("Expected strict mode to be default")
	}
}
