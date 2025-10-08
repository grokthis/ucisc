package disasm

import (
	"testing"

	"github.com/grokthis/ucisc-go/pkg/cpu"
)

func TestFormatInstruction_Copy(t *testing.T) {
	// r1 <- copy val/100
	inst := &cpu.Instruction{
		Source: cpu.Source{
			Type: cpu.SourceVal,
			Reg:  cpu.RegR1,
		},
		Destination: cpu.Destination{
			Type: cpu.DestReg,
			Reg:  cpu.RegR1,
		},
		Operation: cpu.OpCopy,
		Effect:    cpu.EffectAlways,
		Immediate: 100,
	}

	result := FormatInstruction(inst)
	expected := "&r1 <- copy val/100"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatInstruction_Add(t *testing.T) {
	// r1 <- add &r2
	inst := &cpu.Instruction{
		Source: cpu.Source{
			Type: cpu.SourceReg,
			Reg:  cpu.RegR2,
		},
		Destination: cpu.Destination{
			Type: cpu.DestReg,
			Reg:  cpu.RegR1,
		},
		Operation: cpu.OpAdd,
		Effect:    cpu.EffectAlways,
	}

	result := FormatInstruction(inst)
	expected := "&r1 <- add &r2"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatInstruction_MemoryDest(t *testing.T) {
	// r1/5 <- copy val/42
	inst := &cpu.Instruction{
		Source: cpu.Source{
			Type: cpu.SourceVal,
		},
		Destination: cpu.Destination{
			Type: cpu.DestMem,
			Reg:  cpu.RegR1,
		},
		Operation: cpu.OpCopy,
		Effect:    cpu.EffectAlways,
		Immediate: 42,
		Offset:    5,
	}

	result := FormatInstruction(inst)
	expected := "r1/5 <- copy val/42"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatInstruction_MemorySource(t *testing.T) {
	// r2 <- copy r1/3
	inst := &cpu.Instruction{
		Source: cpu.Source{
			Type: cpu.SourceMem,
			Reg:  cpu.RegR1,
		},
		Destination: cpu.Destination{
			Type: cpu.DestReg,
			Reg:  cpu.RegR2,
		},
		Operation: cpu.OpCopy,
		Effect:    cpu.EffectAlways,
		Immediate: 3,
	}

	result := FormatInstruction(inst)
	expected := "&r2 <- copy r1/3"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatInstruction_Push(t *testing.T) {
	// r7 push <- copy val/5
	inst := &cpu.Instruction{
		Source: cpu.Source{
			Type: cpu.SourceVal,
		},
		Destination: cpu.Destination{
			Type: cpu.DestMem,
			Reg:  cpu.RegR7,
		},
		Operation: cpu.OpCopy,
		Effect:    cpu.EffectAlways,
		Immediate: 5,
		Increment: true,
	}

	result := FormatInstruction(inst)
	expected := "r7 push <- copy val/5"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatInstruction_Pop(t *testing.T) {
	// r1 <- copy r7 pop
	inst := &cpu.Instruction{
		Source: cpu.Source{
			Type: cpu.SourceMem,
			Reg:  cpu.RegR7,
		},
		Destination: cpu.Destination{
			Type: cpu.DestReg,
			Reg:  cpu.RegR1,
		},
		Operation: cpu.OpCopy,
		Effect:    cpu.EffectAlways,
		Increment: true,
	}

	result := FormatInstruction(inst)
	expected := "&r1 <- copy r7 pop"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatInstruction_Conditional(t *testing.T) {
	tests := []struct {
		name     string
		effect   cpu.Effect
		expected string
	}{
		{"zero", cpu.EffectZero, "&r1 <0 copy val/0"},
		{"not_zero", cpu.EffectNotZero, "&r1 <1 copy val/0"},
		{"negative", cpu.EffectNegative, "&r1 <n copy val/0"},
		{"positive", cpu.EffectPositive, "&r1 <p copy val/0"},
		{"overflow", cpu.EffectOverflow, "&r1 <o copy val/0"},
		{"interrupt", cpu.EffectInterrupt, "&r1 <# copy val/0"},
		{"flags_only", cpu.EffectFlagsOnly, "&r1 <~ copy val/0"},
		{"always", cpu.EffectAlways, "&r1 <- copy val/0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &cpu.Instruction{
				Source: cpu.Source{
					Type: cpu.SourceVal,
				},
				Destination: cpu.Destination{
					Type: cpu.DestReg,
					Reg:  cpu.RegR1,
				},
				Operation: cpu.OpCopy,
				Effect:    tt.effect,
			}

			result := FormatInstruction(inst)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFormatInstruction_AllOperations(t *testing.T) {
	tests := []struct {
		name     string
		op       cpu.Operation
		opStr    string
	}{
		{"copy", cpu.OpCopy, "copy"},
		{"add", cpu.OpAdd, "add"},
		{"sub", cpu.OpSub, "sub"},
		{"and", cpu.OpAnd, "and"},
		{"or", cpu.OpOr, "or"},
		{"xor", cpu.OpXor, "xor"},
		{"inv", cpu.OpInv, "inv"},
		{"shl", cpu.OpShl, "shl"},
		{"shr", cpu.OpShr, "shr"},
		{"swap", cpu.OpSwap, "swap"},
		{"msb", cpu.OpMsb, "msb"},
		{"lsb", cpu.OpLsb, "lsb"},
		{"mult", cpu.OpMult, "mult"},
		{"multmsw", cpu.OpMultMsw, "multmsw"},
		{"addc", cpu.OpAddc, "addc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &cpu.Instruction{
				Source: cpu.Source{
					Type: cpu.SourceVal,
				},
				Destination: cpu.Destination{
					Type: cpu.DestReg,
					Reg:  cpu.RegR1,
				},
				Operation: tt.op,
				Effect:    cpu.EffectAlways,
			}

			result := FormatInstruction(inst)
			expected := "&r1 <- " + tt.opStr + " val/0"
			if result != expected {
				t.Errorf("Expected %q, got %q", expected, result)
			}
		})
	}
}

func TestFormatInstruction_AllRegisters(t *testing.T) {
	tests := []struct {
		name   string
		reg    cpu.Register
		regStr string
	}{
		{"pc", cpu.RegPC, "pc"},
		{"r1", cpu.RegR1, "r1"},
		{"r2", cpu.RegR2, "r2"},
		{"r3", cpu.RegR3, "r3"},
		{"r5", cpu.RegR5, "r5"},
		{"r6", cpu.RegR6, "r6"},
		{"r7", cpu.RegR7, "r7"},
		{"banking", cpu.RegBanking, "banking"},
		{"flags", cpu.RegFlags, "flags"},
		{"interrupt", cpu.RegInt, "interrupt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &cpu.Instruction{
				Source: cpu.Source{
					Type: cpu.SourceReg,
					Reg:  tt.reg,
				},
				Destination: cpu.Destination{
					Type: cpu.DestReg,
					Reg:  cpu.RegR1,
				},
				Operation: cpu.OpCopy,
				Effect:    cpu.EffectAlways,
			}

			result := FormatInstruction(inst)
			expected := "&r1 <- copy &" + tt.regStr
			if result != expected {
				t.Errorf("Expected %q, got %q", expected, result)
			}
		})
	}
}

func TestFormatInstruction_NegativeImmediate(t *testing.T) {
	// r1 <- copy val/-5
	inst := &cpu.Instruction{
		Source: cpu.Source{
			Type: cpu.SourceVal,
		},
		Destination: cpu.Destination{
			Type: cpu.DestReg,
			Reg:  cpu.RegR1,
		},
		Operation: cpu.OpCopy,
		Effect:    cpu.EffectAlways,
		Immediate: -5,
	}

	result := FormatInstruction(inst)
	expected := "&r1 <- copy val/-5"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatInstruction_PCRelative(t *testing.T) {
	// pc <- copy &pc/10
	inst := &cpu.Instruction{
		Source: cpu.Source{
			Type: cpu.SourceReg,
			Reg:  cpu.RegPC,
		},
		Destination: cpu.Destination{
			Type: cpu.DestReg,
			Reg:  cpu.RegPC,
		},
		Operation: cpu.OpCopy,
		Effect:    cpu.EffectAlways,
		Immediate: 10,
	}

	result := FormatInstruction(inst)
	expected := "&pc <- copy &pc/10"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatInstruction_Halt(t *testing.T) {
	// pc <- copy &pc (halt instruction)
	inst := &cpu.Instruction{
		Source: cpu.Source{
			Type: cpu.SourceReg,
			Reg:  cpu.RegPC,
		},
		Destination: cpu.Destination{
			Type: cpu.DestReg,
			Reg:  cpu.RegPC,
		},
		Operation: cpu.OpCopy,
		Effect:    cpu.EffectAlways,
		Immediate: 0,
	}

	result := FormatInstruction(inst)
	expected := "&pc <- copy &pc"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
