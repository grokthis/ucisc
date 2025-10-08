package asm

import "testing"

// TestInstructionAlignment_OddData verifies that instructions are aligned to 2-word boundaries
// when preceded by an odd number of data words
func TestInstructionAlignment_OddData(t *testing.T) {
	source := `
% 0001 0002 0003
&r1 <- copy val/42
`
	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Verify structure:
	// Word 0-2: Data (0001, 0002, 0003)
	// Word 3: Padding (0000)
	// Word 4-5: Instruction
	if len(code) != 6 {
		t.Fatalf("Expected 6 words (3 data + 1 padding + 2 instruction), got %d", len(code))
	}

	if code[0] != 0x0001 || code[1] != 0x0002 || code[2] != 0x0003 {
		t.Errorf("Data words incorrect: %04X %04X %04X", code[0], code[1], code[2])
	}

	if code[3] != 0x0000 {
		t.Errorf("Expected padding word at index 3, got %04X", code[3])
	}

	// Words 4-5 should be the instruction (exact encoding doesn't matter, just that it exists)
	if code[4] == 0x0000 && code[5] == 0x0000 {
		t.Error("Instruction appears to be missing or all zeros")
	}
}

// TestInstructionAlignment_EvenData verifies no padding is added when data is even
func TestInstructionAlignment_EvenData(t *testing.T) {
	source := `
% 0001 0002
&r1 <- copy val/42
`
	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Verify structure:
	// Word 0-1: Data (0001, 0002)
	// Word 2-3: Instruction (no padding needed)
	if len(code) != 4 {
		t.Fatalf("Expected 4 words (2 data + 2 instruction), got %d", len(code))
	}

	if code[0] != 0x0001 || code[1] != 0x0002 {
		t.Errorf("Data words incorrect: %04X %04X", code[0], code[1])
	}

	// Words 2-3 should be the instruction immediately (no padding)
	if code[2] == 0x0000 && code[3] == 0x0000 {
		t.Error("Instruction appears to be missing or all zeros")
	}
}

// TestInstructionAlignment_OddString verifies alignment with odd-length strings
func TestInstructionAlignment_OddString(t *testing.T) {
	source := `
"ABC"
&r1 <- copy val/42
`
	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Verify structure:
	// Word 0-2: String "ABC" (3 chars)
	// Word 3: Null terminator
	// Word 4: Padding (total 4 words is even, so no padding needed)
	// Word 4-5: Instruction
	// Actually "ABC" + null = 4 words (even), so no padding
	expectedLen := 4 + 2 // 4 for string+null, 2 for instruction
	if len(code) != expectedLen {
		t.Fatalf("Expected %d words, got %d", expectedLen, len(code))
	}

	// Last two words should be instruction
	if code[4] == 0x0000 && code[5] == 0x0000 {
		t.Error("Instruction appears to be missing or all zeros")
	}
}

// TestInstructionAlignment_DefInitializer verifies alignment for def with initializer
func TestInstructionAlignment_DefInitializer(t *testing.T) {
	source := `
% 0001
def stack/r1 <- copy val/100
`
	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Verify structure:
	// Word 0: Data (0001)
	// Word 1: Padding
	// Word 2-3: Def initializer instruction
	if len(code) != 4 {
		t.Fatalf("Expected 4 words (1 data + 1 padding + 2 instruction), got %d", len(code))
	}

	if code[0] != 0x0001 {
		t.Errorf("Data word incorrect: %04X", code[0])
	}

	if code[1] != 0x0000 {
		t.Errorf("Expected padding word at index 1, got %04X", code[1])
	}
}

// TestInstructionAlignment_VarInitializer verifies alignment for var with initializer
func TestInstructionAlignment_VarInitializer(t *testing.T) {
	source := `
def obj/r2
% 0001
var obj.x/0 <- copy val/100
`
	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Verify structure:
	// Word 0: Data (0001)
	// Word 1: Padding
	// Word 2-3: Var initializer instruction
	if len(code) != 4 {
		t.Fatalf("Expected 4 words (1 data + 1 padding + 2 instruction), got %d", len(code))
	}

	if code[0] != 0x0001 {
		t.Errorf("Data word incorrect: %04X", code[0])
	}

	if code[1] != 0x0000 {
		t.Errorf("Expected padding word at index 1, got %04X", code[1])
	}
}

// TestInstructionAlignment_Function verifies alignment for functions
func TestInstructionAlignment_Function(t *testing.T) {
	source := `
def stack/r1
% 0001
fun stack.test() {
	&r1 <- copy val/42
}
`
	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Verify structure:
	// Word 0: Data (0001)
	// Word 1: Padding
	// Word 2-3: Function body instruction
	if len(code) != 4 {
		t.Fatalf("Expected 4 words (1 data + 1 padding + 2 instruction), got %d", len(code))
	}

	if code[0] != 0x0001 {
		t.Errorf("Data word incorrect: %04X", code[0])
	}

	if code[1] != 0x0000 {
		t.Errorf("Expected padding word at index 1, got %04X", code[1])
	}
}

// TestInstructionAlignment_MultipleOddBlocks verifies alignment with multiple odd data blocks
func TestInstructionAlignment_MultipleOddBlocks(t *testing.T) {
	source := `
% 0001
&r1 <- copy val/1
% 0002
&r2 <- copy val/2
`
	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Verify structure:
	// Word 0: Data (0001)
	// Word 1: Padding
	// Word 2-3: First instruction
	// Word 4: Data (0002)
	// Word 5: Padding
	// Word 6-7: Second instruction
	expectedLen := 8
	if len(code) != expectedLen {
		t.Fatalf("Expected %d words, got %d", expectedLen, len(code))
	}

	if code[0] != 0x0001 {
		t.Errorf("First data word incorrect: %04X", code[0])
	}
	if code[1] != 0x0000 {
		t.Errorf("First padding word incorrect: %04X", code[1])
	}
	if code[4] != 0x0002 {
		t.Errorf("Second data word incorrect: %04X", code[4])
	}
	if code[5] != 0x0000 {
		t.Errorf("Second padding word incorrect: %04X", code[5])
	}
}
