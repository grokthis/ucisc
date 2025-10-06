package asm

import (
	"testing"

	"github.com/grokthis/ucisc-go/pkg/cpu"
)

func TestParser_SimpleInstruction(t *testing.T) {
	input := "&r1 <- copy val/42"
	p := NewParser(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	inst, ok := program.Statements[0].(*InstructionStmt)
	if !ok {
		t.Fatalf("Expected InstructionStmt, got %T", program.Statements[0])
	}

	// Check destination
	if inst.Destination.Type != OpRegister {
		t.Errorf("Expected OpRegister, got %v", inst.Destination.Type)
	}
	if inst.Destination.Register != cpu.RegR1 {
		t.Errorf("Expected RegR1, got %v", inst.Destination.Register)
	}
	if !inst.Destination.IsAddress {
		t.Error("Expected IsAddress to be true")
	}

	// Check effect
	if inst.Effect != cpu.EffectAlways {
		t.Errorf("Expected EffectAlways, got %v", inst.Effect)
	}

	// Check operation
	if inst.Operation != cpu.OpCopy {
		t.Errorf("Expected OpCopy, got %v", inst.Operation)
	}

	// Check source
	if inst.Source.Type != OpImmediate {
		t.Errorf("Expected OpImmediate, got %v", inst.Source.Type)
	}
	if inst.Source.Offset != 42 {
		t.Errorf("Expected offset 42, got %d", inst.Source.Offset)
	}
}

func TestParser_MemoryDestination(t *testing.T) {
	input := "r1/5 <- add r2/10"
	p := NewParser(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	inst := program.Statements[0].(*InstructionStmt)

	// Check destination
	if inst.Destination.Type != OpMemory {
		t.Errorf("Expected OpMemory, got %v", inst.Destination.Type)
	}
	if inst.Destination.Register != cpu.RegR1 {
		t.Errorf("Expected RegR1, got %v", inst.Destination.Register)
	}
	if inst.Destination.Offset != 5 {
		t.Errorf("Expected offset 5, got %d", inst.Destination.Offset)
	}
	if inst.Destination.IsAddress {
		t.Error("Expected IsAddress to be false")
	}

	// Check operation
	if inst.Operation != cpu.OpAdd {
		t.Errorf("Expected OpAdd, got %v", inst.Operation)
	}

	// Check source
	if inst.Source.Type != OpMemory {
		t.Errorf("Expected OpMemory, got %v", inst.Source.Type)
	}
	if inst.Source.Register != cpu.RegR2 {
		t.Errorf("Expected RegR2, got %v", inst.Source.Register)
	}
	if inst.Source.Offset != 10 {
		t.Errorf("Expected offset 10, got %d", inst.Source.Offset)
	}
}

func TestParser_ConditionalInstruction(t *testing.T) {
	input := "pc <1 copy val/10"
	p := NewParser(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	inst := program.Statements[0].(*InstructionStmt)

	if inst.Effect != cpu.EffectNotZero {
		t.Errorf("Expected EffectNotZero, got %v", inst.Effect)
	}

	if inst.Destination.Type != OpMemory {
		t.Errorf("Expected OpMemory (pc without &), got %v", inst.Destination.Type)
	}
	if inst.Destination.Register != cpu.RegPC {
		t.Errorf("Expected RegPC, got %v", inst.Destination.Register)
	}
}

func TestParser_PushModifier(t *testing.T) {
	input := "r1/0 push <- copy val/100"
	p := NewParser(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	inst := program.Statements[0].(*InstructionStmt)

	if !inst.Push {
		t.Error("Expected Push to be true")
	}

	if inst.Destination.Type != OpMemory {
		t.Errorf("Expected OpMemory, got %v", inst.Destination.Type)
	}
}

func TestParser_PopModifier(t *testing.T) {
	input := "pc <- copy r1/0 pop"
	p := NewParser(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	inst := program.Statements[0].(*InstructionStmt)

	if !inst.Pop {
		t.Error("Expected Pop to be true")
	}

	if inst.Source.Type != OpMemory {
		t.Errorf("Expected OpMemory, got %v", inst.Source.Type)
	}
}

func TestParser_Label(t *testing.T) {
	input := `loop:
    r1 <- add val/1`

	p := NewParser(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	if len(program.Statements) != 2 {
		t.Fatalf("Expected 2 statements, got %d", len(program.Statements))
	}

	label, ok := program.Statements[0].(*LabelStmt)
	if !ok {
		t.Fatalf("Expected LabelStmt, got %T", program.Statements[0])
	}

	if label.Name != "loop" {
		t.Errorf("Expected label name 'loop', got %q", label.Name)
	}

	inst, ok := program.Statements[1].(*InstructionStmt)
	if !ok {
		t.Fatalf("Expected InstructionStmt, got %T", program.Statements[1])
	}

	if inst.Operation != cpu.OpAdd {
		t.Errorf("Expected OpAdd, got %v", inst.Operation)
	}
}

func TestParser_LabelReference(t *testing.T) {
	input := "pc <1 copy pc/loop"
	p := NewParser(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	inst := program.Statements[0].(*InstructionStmt)

	if inst.Source.Type != OpPCRelative {
		t.Errorf("Expected OpPCRelative, got %v", inst.Source.Type)
	}

	if inst.Source.LabelRef != "loop" {
		t.Errorf("Expected label ref 'loop', got %q", inst.Source.LabelRef)
	}
}

func TestParser_AllOperations(t *testing.T) {
	operations := []struct {
		input string
		op    cpu.Operation
	}{
		{"&r1 <- copy val/1", cpu.OpCopy},
		{"&r1 <- and val/1", cpu.OpAnd},
		{"&r1 <- or val/1", cpu.OpOr},
		{"&r1 <- xor val/1", cpu.OpXor},
		{"&r1 <- inv val/1", cpu.OpInv},
		{"&r1 <- shl val/1", cpu.OpShl},
		{"&r1 <- shr val/1", cpu.OpShr},
		{"&r1 <- swap val/1", cpu.OpSwap},
		{"&r1 <- msb val/1", cpu.OpMsb},
		{"&r1 <- lsb val/1", cpu.OpLsb},
		{"&r1 <- add val/1", cpu.OpAdd},
		{"&r1 <- addc val/1", cpu.OpAddc},
		{"&r1 <- sub val/1", cpu.OpSub},
		{"&r1 <- mult val/1", cpu.OpMult},
		{"&r1 <- multsw val/1", cpu.OpMultMsw},
	}

	for _, tt := range operations {
		p := NewParser(tt.input)
		program, err := p.Parse()

		if err != nil {
			t.Errorf("input=%q - Parse() returned error: %v", tt.input, err)
			continue
		}

		inst := program.Statements[0].(*InstructionStmt)
		if inst.Operation != tt.op {
			t.Errorf("input=%q - Expected %v, got %v", tt.input, tt.op, inst.Operation)
		}
	}
}

func TestParser_AllEffects(t *testing.T) {
	effects := []struct {
		input  string
		effect cpu.Effect
	}{
		{"&r1 <- copy val/1", cpu.EffectAlways},
		{"&r1 <~ copy val/1", cpu.EffectFlagsOnly},
		{"&r1 <0 copy val/1", cpu.EffectZero},
		{"&r1 <1 copy val/1", cpu.EffectNotZero},
		{"&r1 <+ copy val/1", cpu.EffectPositive},
		{"&r1 <n copy val/1", cpu.EffectNegative},
		{"&r1 <& copy val/1", cpu.EffectOverflow},
		{"&r1 <# copy val/1", cpu.EffectInterrupt},
	}

	for _, tt := range effects {
		p := NewParser(tt.input)
		program, err := p.Parse()

		if err != nil {
			t.Errorf("input=%q - Parse() returned error: %v", tt.input, err)
			continue
		}

		inst := program.Statements[0].(*InstructionStmt)
		if inst.Effect != tt.effect {
			t.Errorf("input=%q - Expected %v, got %v", tt.input, tt.effect, inst.Effect)
		}
	}
}

func TestParser_NegativeNumber(t *testing.T) {
	input := "&r1 <- copy val/-5"
	p := NewParser(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	inst := program.Statements[0].(*InstructionStmt)
	if inst.Source.Offset != -5 {
		t.Errorf("Expected offset -5, got %d", inst.Source.Offset)
	}
}

func TestParser_HexNumber(t *testing.T) {
	input := "&r1 <- copy val/0xFFFF"
	p := NewParser(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	inst := program.Statements[0].(*InstructionStmt)
	if inst.Source.Offset != -1 { // 0xFFFF as int16 is -1
		t.Errorf("Expected offset -1, got %d", inst.Source.Offset)
	}
}

func TestParser_CommentsIgnored(t *testing.T) {
	input := `# This is a comment
    &r1 <- copy val/42  # inline comment
    # another comment`

	p := NewParser(input)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement (comments ignored), got %d", len(program.Statements))
	}
}
