package asm

import (
	"testing"
)

// TestCustomElementParsers verifies custom parser ordering works
func TestCustomElementParsers(t *testing.T) {
	source := `&r1 <- copy val/42`

	// Create parser with default parsers
	p1 := NewParser(source)
	prog1, err := p1.Parse()
	if err != nil {
		t.Fatalf("Default parser failed: %v", err)
	}
	if len(prog1.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog1.Statements))
	}

	// Create parser with custom parser list (same as default)
	customParsers := DefaultElementParsers()
	p2 := NewParserWithElementParsers(source, customParsers)
	prog2, err := p2.Parse()
	if err != nil {
		t.Fatalf("Custom parser failed: %v", err)
	}
	if len(prog2.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog2.Statements))
	}
}

// TestDefaultElementParsersOrder verifies the default parser order
func TestDefaultElementParsersOrder(t *testing.T) {
	parsers := DefaultElementParsers()

	// Verify we have all expected parsers
	expectedCount := 10 // Label, Def, Var, Fun, Debug, Call, Data, String, Scope, Instruction
	if len(parsers) != expectedCount {
		t.Errorf("Expected %d default parsers, got %d", expectedCount, len(parsers))
	}

	// Verify order - InstructionParser should be last (most general)
	if len(parsers) > 0 {
		if _, ok := parsers[len(parsers)-1].(*InstructionParser); !ok {
			t.Error("Expected InstructionParser to be last in default order")
		}
	}

	// Verify LabelParser is first
	if len(parsers) > 0 {
		if _, ok := parsers[0].(*LabelParser); !ok {
			t.Error("Expected LabelParser to be first in default order")
		}
	}
}

// TestMinimalParserSet demonstrates using a minimal parser set
func TestMinimalParserSet(t *testing.T) {
	// Only allow instructions and labels
	minimalParsers := []ElementParser{
		&LabelParser{},
		&InstructionParser{},
	}

	source := `
label1:
&r1 <- copy val/42
pc <- copy pc/0
`

	p := NewParserWithElementParsers(source, minimalParsers)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("Minimal parser failed: %v", err)
	}

	// Should have 3 statements: label, instruction, instruction
	if len(prog.Statements) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(prog.Statements))
	}
}

// TestParserRejectsUnsupportedSyntax verifies minimal parser rejects unsupported syntax
func TestParserRejectsUnsupportedSyntax(t *testing.T) {
	// Parser without DefParser
	limitedParsers := []ElementParser{
		&LabelParser{},
		&InstructionParser{},
	}

	source := `def stack/r1`

	p := NewParserWithElementParsers(source, limitedParsers)
	_, err := p.Parse()

	// Should fail because DefParser is not in the list
	// InstructionParser will try to parse it and fail
	if err == nil {
		t.Error("Expected parser to reject 'def' statement without DefParser")
	}
}
