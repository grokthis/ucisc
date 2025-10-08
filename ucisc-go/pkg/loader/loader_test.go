package loader

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoadHex_Simple(t *testing.T) {
	input := `4540 0064
0040 0000`

	reader := strings.NewReader(input)
	program, err := LoadHex(reader)
	if err != nil {
		t.Fatalf("LoadHex failed: %v", err)
	}

	expected := []uint16{0x4540, 0x0064, 0x0040, 0x0000}
	if len(program) != len(expected) {
		t.Fatalf("Expected %d words, got %d", len(expected), len(program))
	}

	for i, word := range expected {
		if program[i] != word {
			t.Errorf("Word %d: expected 0x%04X, got 0x%04X", i, word, program[i])
		}
	}
}

func TestLoadHex_WithComments(t *testing.T) {
	input := `# This is a program
4540 0064  # r1 <- copy val/100
0040 0000  # halt`

	reader := strings.NewReader(input)
	program, err := LoadHex(reader)
	if err != nil {
		t.Fatalf("LoadHex failed: %v", err)
	}

	expected := []uint16{0x4540, 0x0064, 0x0040, 0x0000}
	if len(program) != len(expected) {
		t.Fatalf("Expected %d words, got %d", len(expected), len(program))
	}

	for i, word := range expected {
		if program[i] != word {
			t.Errorf("Word %d: expected 0x%04X, got 0x%04X", i, word, program[i])
		}
	}
}

func TestLoadHex_BlankLines(t *testing.T) {
	input := `
4540 0064

0040 0000

`

	reader := strings.NewReader(input)
	program, err := LoadHex(reader)
	if err != nil {
		t.Fatalf("LoadHex failed: %v", err)
	}

	expected := []uint16{0x4540, 0x0064, 0x0040, 0x0000}
	if len(program) != len(expected) {
		t.Fatalf("Expected %d words, got %d", len(expected), len(program))
	}
}

func TestLoadHex_SingleWord(t *testing.T) {
	input := "FFFF"

	reader := strings.NewReader(input)
	program, err := LoadHex(reader)
	if err != nil {
		t.Fatalf("LoadHex failed: %v", err)
	}

	if len(program) != 1 {
		t.Fatalf("Expected 1 word, got %d", len(program))
	}

	if program[0] != 0xFFFF {
		t.Errorf("Expected 0xFFFF, got 0x%04X", program[0])
	}
}

func TestLoadHex_InvalidHex(t *testing.T) {
	input := "ZZZZ"

	reader := strings.NewReader(input)
	_, err := LoadHex(reader)
	if err == nil {
		t.Fatal("Expected error for invalid hex, got nil")
	}

	if !strings.Contains(err.Error(), "invalid hex value") {
		t.Errorf("Expected 'invalid hex value' error, got: %v", err)
	}
}

func TestLoadHex_Empty(t *testing.T) {
	input := ""

	reader := strings.NewReader(input)
	_, err := LoadHex(reader)
	if err == nil {
		t.Fatal("Expected error for empty input, got nil")
	}

	if !strings.Contains(err.Error(), "no program data found") {
		t.Errorf("Expected 'no program data found' error, got: %v", err)
	}
}

func TestLoadHex_OnlyComments(t *testing.T) {
	input := `# Just comments
# Nothing else`

	reader := strings.NewReader(input)
	_, err := LoadHex(reader)
	if err == nil {
		t.Fatal("Expected error for comments only, got nil")
	}

	if !strings.Contains(err.Error(), "no program data found") {
		t.Errorf("Expected 'no program data found' error, got: %v", err)
	}
}

func TestLoadHex_MultipleWordsPerLine(t *testing.T) {
	input := "4540 0064 0040 0000"

	reader := strings.NewReader(input)
	program, err := LoadHex(reader)
	if err != nil {
		t.Fatalf("LoadHex failed: %v", err)
	}

	expected := []uint16{0x4540, 0x0064, 0x0040, 0x0000}
	if len(program) != len(expected) {
		t.Fatalf("Expected %d words, got %d", len(expected), len(program))
	}

	for i, word := range expected {
		if program[i] != word {
			t.Errorf("Word %d: expected 0x%04X, got 0x%04X", i, word, program[i])
		}
	}
}

func TestLoadHex_LineNumbers(t *testing.T) {
	input := `4540 0064
ZZZZ
0040 0000`

	reader := strings.NewReader(input)
	_, err := LoadHex(reader)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("Expected error to include 'line 2', got: %v", err)
	}
}

func TestLoadBinary_Simple(t *testing.T) {
	// Little-endian 16-bit words: 0x4540, 0x0064, 0x0040, 0x0000
	data := []byte{
		0x40, 0x45, // 0x4540
		0x64, 0x00, // 0x0064
		0x40, 0x00, // 0x0040
		0x00, 0x00, // 0x0000
	}

	reader := bytes.NewReader(data)
	program, err := LoadBinary(reader)
	if err != nil {
		t.Fatalf("LoadBinary failed: %v", err)
	}

	expected := []uint16{0x4540, 0x0064, 0x0040, 0x0000}
	if len(program) != len(expected) {
		t.Fatalf("Expected %d words, got %d", len(expected), len(program))
	}

	for i, word := range expected {
		if program[i] != word {
			t.Errorf("Word %d: expected 0x%04X, got 0x%04X", i, word, program[i])
		}
	}
}

func TestLoadBinary_Empty(t *testing.T) {
	data := []byte{}

	reader := bytes.NewReader(data)
	program, err := LoadBinary(reader)
	if err != nil {
		t.Fatalf("LoadBinary failed: %v", err)
	}

	if len(program) != 0 {
		t.Errorf("Expected empty program, got %d words", len(program))
	}
}

func TestLoadBinary_IncompleteWord(t *testing.T) {
	// Only 1 byte (incomplete 16-bit word)
	data := []byte{0x40}

	reader := bytes.NewReader(data)
	_, err := LoadBinary(reader)
	if err == nil {
		t.Fatal("Expected error for incomplete word, got nil")
	}

	// Can be either "incomplete word" or "unexpected EOF"
	errMsg := err.Error()
	if !strings.Contains(errMsg, "incomplete word") && !strings.Contains(errMsg, "unexpected EOF") {
		t.Errorf("Expected 'incomplete word' or 'unexpected EOF' error, got: %v", err)
	}
}

func TestLoadBinary_SingleWord(t *testing.T) {
	data := []byte{0xFF, 0xFF} // 0xFFFF

	reader := bytes.NewReader(data)
	program, err := LoadBinary(reader)
	if err != nil {
		t.Fatalf("LoadBinary failed: %v", err)
	}

	if len(program) != 1 {
		t.Fatalf("Expected 1 word, got %d", len(program))
	}

	if program[0] != 0xFFFF {
		t.Errorf("Expected 0xFFFF, got 0x%04X", program[0])
	}
}
