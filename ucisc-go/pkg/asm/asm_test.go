package asm

import (
	"strings"
	"testing"

	"github.com/grokthis/ucisc-go/pkg/cpu"
)

func TestAssemble_SimpleProgram(t *testing.T) {
	// Simple test program - just add some numbers to r1 register
	source := `
# Simple addition test

    &r1 <- copy val/5
    &r1 <- add val/3
    &r1 <- add val/2

    pc <- copy pc/0       # halt
`

	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Should generate 4 instructions * 2 words = 8 words
	if len(code) != 8 {
		t.Errorf("Expected 8 words, got %d", len(code))
	}

	// Test by running in emulator
	c := cpu.New(1024)
	err = c.LoadProgram(code, 0)
	if err != nil {
		t.Fatalf("LoadProgram failed: %v", err)
	}

	err = c.Run()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// r1 should contain 10 (5+3+2)
	if c.GetR1() != 10 {
		t.Errorf("Expected r1=10, got %d", c.GetR1())
	}
}

func TestAssemble_SimpleLoop(t *testing.T) {
	// Simple loop that adds 1 five times
	source := `
    &r1 <- copy val/0     # result
    &r2 <- copy val/5     # counter

loop:
    &r1 <- add val/1      # r1++
    &r2 <- sub val/1      # r2--
    &r2 <~ sub val/0      # test r2
    pc <1 copy val/loop   # if r2 != 0, jump to loop

    pc <- copy pc/0       # halt
`

	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Run in emulator
	c := cpu.New(1024)
	c.LoadProgram(code, 0)
	c.Run()

	if c.GetR1() != 5 {
		t.Errorf("Expected r1=5, got %d", c.GetR1())
	}
}

func TestAssemble_PushPop(t *testing.T) {
	source := `
# Test push and pop
    &r1 <- copy val/1000        # stack pointer

    r1/0 push <- copy val/42    # push 42
    r1/0 push <- copy val/99    # push 99

    &r2 <- copy r1/0 pop        # pop into r2 (should be 99)
    &r3 <- copy r1/0 pop        # pop into r3 (should be 42)

    pc <- copy pc/0             # halt
`

	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Run in emulator
	c := cpu.New(1024)
	c.LoadProgram(code, 0)
	c.Run()

	if c.GetR2() != 99 {
		t.Errorf("Expected r2=99, got %d", c.GetR2())
	}
	if c.GetR3() != 42 {
		t.Errorf("Expected r3=42, got %d", c.GetR3())
	}
	// Stack pointer should be back to 1000
	if c.GetR1() != 1000 {
		t.Errorf("Expected r1=1000, got %d", c.GetR1())
	}
}

func TestAssemble_AllOperations(t *testing.T) {
	source := `
    &r1 <- copy val/5
    &r1 <- and val/3
    &r1 <- or val/8
    &r1 <- xor val/15
    &r1 <- inv val/0
    &r1 <- shl val/1
    &r1 <- shr val/1
    &r1 <- add val/10
    &r1 <- sub val/5
    pc <- copy pc/0
`

	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Just verify it assembles and runs without error
	c := cpu.New(1024)
	c.LoadProgram(code, 0)
	err = c.Run()
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}
}

func TestAssemble_AllEffects(t *testing.T) {
	source := `
    # Test zero flag
    &r1 <- copy val/0       # Set r1 to 0
    &r1 <~ sub val/0        # Set flags (0 - 0 = 0, sets zero flag)
    &r2 <0 copy val/1       # Should execute (zero flag is set)
    &r3 <1 copy val/2       # Should not execute (zero flag is set, want not-zero)

    # Test negative flag
    &r1 <- copy val/5       # Set r1 to 5
    &r1 <~ sub val/10       # Set flags only (5 - 10 = -5, sets negative flag)
    &r5 <n copy val/3       # Should execute (negative flag is set)
    &r6 <+ copy val/4       # Should not execute (negative flag is set, want positive)

    pc <- copy pc/0         # halt
`

	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	c := cpu.New(1024)
	c.LoadProgram(code, 0)
	c.Run()

	if c.GetR2() != 1 {
		t.Errorf("Expected r2=1 (zero condition), got %d", c.GetR2())
	}
	if c.GetR3() != 0 {
		t.Errorf("Expected r3=0 (not-zero should not execute), got %d", c.GetR3())
	}
	if c.GetR5() != 3 {
		t.Errorf("Expected r5=3 (negative condition), got %d", c.GetR5())
	}
	if c.GetR6() != 0 {
		t.Errorf("Expected r6=0 (positive should not execute), got %d", c.GetR6())
	}
}

func TestAssembleToHex(t *testing.T) {
	source := "&r1 <- copy val/42"

	hex, err := AssembleToHex(source)
	if err != nil {
		t.Fatalf("AssembleToHex failed: %v", err)
	}

	// Should produce "4540 002A"
	expected := "4540 002A"
	if hex != expected {
		t.Errorf("Expected %q, got %q", expected, hex)
	}
}

func TestAssemble_SyntaxError(t *testing.T) {
	source := "&r1 <- invalid_op val/42"

	_, err := Assemble(source)
	if err == nil {
		t.Fatal("Expected error for invalid operation, got nil")
	}
}

func TestAssemble_LabelError(t *testing.T) {
	source := "pc <- copy pc/undefined"

	_, err := Assemble(source)
	if err == nil {
		t.Fatal("Expected error for undefined label, got nil")
	}
}

func TestFormatHex(t *testing.T) {
	code := []uint16{0x4540, 0x002A, 0x0040, 0x0000}

	hex := FormatHex(code)

	// Should format with spaces, newline every 8 words
	expected := "4540 002A 0040 0000"
	if hex != expected {
		t.Errorf("Expected %q, got %q", expected, hex)
	}
}

func TestFormatHex_MultiLine(t *testing.T) {
	// Create 16 words to test newline formatting
	code := make([]uint16, 16)
	for i := range code {
		code[i] = uint16(i)
	}

	hex := FormatHex(code)

	// Should have a newline after 8 words
	lines := strings.Split(hex, "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
}
