package asm

import (
	"testing"
)

func TestCodeGen_SimpleInstruction(t *testing.T) {
	input := "&r1 <- copy val/42"
	p := NewParser(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	gen := NewCodeGenerator(program)
	code, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	if len(code) != 2 {
		t.Fatalf("Expected 2 words, got %d", len(code))
	}

	// The encoding should be: src=val(4), dest=r1(5), effect=always(4), op=copy(0)
	// Word1: 0x4540 (0100 0101 0100 0000)
	//   src=4 (0100), dest=5 (0101), M=0, effect=4 (100), op=0 (0000)
	// Word2: 0x002A (42)
	expectedWord1 := uint16(0x4540)
	expectedWord2 := uint16(42)

	if code[0] != expectedWord1 {
		t.Errorf("Word1: expected 0x%04X, got 0x%04X", expectedWord1, code[0])
	}
	if code[1] != expectedWord2 {
		t.Errorf("Word2: expected 0x%04X, got 0x%04X", expectedWord2, code[1])
	}
}

func TestCodeGen_MemoryDestination(t *testing.T) {
	input := "r1/5 <- add r2/10"
	p := NewParser(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	gen := NewCodeGenerator(program)
	code, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	if len(code) != 2 {
		t.Fatalf("Expected 2 words, got %d", len(code))
	}

	// src=r2 mem(2), dest=r1 mem(1), effect=always(4), op=add(A)
	// Word1: 0x214A
	// Word2: offset(5) in top 4 bits, immediate(10) in bottom 12 bits = 0x500A
	expectedWord1 := uint16(0x214A)
	expectedWord2 := uint16(0x500A)

	if code[0] != expectedWord1 {
		t.Errorf("Word1: expected 0x%04X, got 0x%04X", expectedWord1, code[0])
	}
	if code[1] != expectedWord2 {
		t.Errorf("Word2: expected 0x%04X, got 0x%04X", expectedWord2, code[1])
	}
}

func TestCodeGen_ConditionalJump(t *testing.T) {
	input := "pc <1 copy val/10"
	p := NewParser(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	gen := NewCodeGenerator(program)
	code, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	// src=val(4), dest=pc mem(0), effect=not-zero(1), op=copy(0)
	// Word1: 0x4010
	expectedWord1 := uint16(0x4010)
	expectedWord2 := uint16(10)

	if code[0] != expectedWord1 {
		t.Errorf("Word1: expected 0x%04X, got 0x%04X", expectedWord1, code[0])
	}
	if code[1] != expectedWord2 {
		t.Errorf("Word2: expected 0x%04X, got 0x%04X", expectedWord2, code[1])
	}
}

func TestCodeGen_WithLabel(t *testing.T) {
	input := `
loop:
    &r1 <- add val/1
    pc <1 copy pc/loop
`
	p := NewParser(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	gen := NewCodeGenerator(program)
	code, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	// Should have 2 instructions = 4 words
	if len(code) != 4 {
		t.Fatalf("Expected 4 words, got %d", len(code))
	}

	// Second instruction should jump to address 0 (where loop is)
	// src=pc(0), dest=pc(0), effect=not-zero(1), op=copy(0)
	// Word1: 0x0010
	// Word2: 0x0000 (loop address)
	if code[2] != 0x0010 {
		t.Errorf("Jump instruction word1: expected 0x0010, got 0x%04X", code[2])
	}
	if code[3] != 0x0000 {
		t.Errorf("Jump instruction word2: expected 0x0000 (loop address), got 0x%04X", code[3])
	}
}

func TestCodeGen_ForwardLabel(t *testing.T) {
	input := `
    pc <- copy pc/end
end:
`
	p := NewParser(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	gen := NewCodeGenerator(program)
	code, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	// Jump should resolve to address 2 (after first instruction)
	if code[1] != 0x0002 {
		t.Errorf("Forward label: expected 0x0002, got 0x%04X", code[1])
	}
}

func TestCodeGen_Push(t *testing.T) {
	input := "r1/0 push <- copy val/100"
	p := NewParser(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	gen := NewCodeGenerator(program)
	code, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	// src=val(4), dest=r1 mem(1), M=1 (push), effect=always(4), op=copy(0)
	// Word1: 0x41C0 (M bit set)
	expectedWord1 := uint16(0x41C0)

	if code[0] != expectedWord1 {
		t.Errorf("Push instruction: expected 0x%04X, got 0x%04X", expectedWord1, code[0])
	}
}

func TestCodeGen_Pop(t *testing.T) {
	input := "pc <- copy r1/0 pop"
	p := NewParser(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	gen := NewCodeGenerator(program)
	code, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	// src=r1 mem(1), dest=pc mem(0), M=1 (pop), effect=always(4), op=copy(0)
	// Word1: 0x10C0 (M bit set)
	expectedWord1 := uint16(0x10C0)

	if code[0] != expectedWord1 {
		t.Errorf("Pop instruction: expected 0x%04X, got 0x%04X", expectedWord1, code[0])
	}
}

func TestCodeGen_Halt(t *testing.T) {
	input := "pc <- copy pc/0"
	p := NewParser(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	gen := NewCodeGenerator(program)
	code, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	// HALT: src=pc(0), dest=pc(0), effect=always(4), op=copy(0)
	// Word1: 0x0040
	// Word2: 0x0000
	if code[0] != 0x0040 {
		t.Errorf("HALT word1: expected 0x0040, got 0x%04X", code[0])
	}
	if code[1] != 0x0000 {
		t.Errorf("HALT word2: expected 0x0000, got 0x%04X", code[1])
	}
}

func TestCodeGen_UndefinedLabel(t *testing.T) {
	input := "pc <- copy pc/missing"
	p := NewParser(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	gen := NewCodeGenerator(program)
	_, err = gen.Generate()

	if err == nil {
		t.Fatal("Expected error for undefined label, got nil")
	}
}

func TestCodeGen_DuplicateLabel(t *testing.T) {
	input := `
loop:
loop:
`
	p := NewParser(input)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	gen := NewCodeGenerator(program)
	_, err = gen.Generate()

	if err == nil {
		t.Fatal("Expected error for duplicate label, got nil")
	}
}
