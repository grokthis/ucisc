package cpu

import "testing"

func TestDecode_SimpleInstructions(t *testing.T) {
	tests := []struct {
		name string
		raw  uint32
		want Instruction
	}{
		{
			name: "copy val/5 to r1",
			raw:  0x45400005, // r1 <- copy val/5 (src=4, dest=5, effect=4, op=0)
			want: Instruction{
				Source:      Source{Type: SourceVal, Reg: 0},
				Destination: Destination{Type: DestReg, Reg: RegR1},
				Increment:   false,
				Effect:      EffectAlways,
				Operation:   OpCopy,
				Immediate:   5,
				Offset:      0,
			},
		},
		{
			name: "add r2 to r1",
			raw:  0x654A0000, // r1 <- add r2 (src=6(r2 reg), dest=5(r1 reg), effect=4, op=A)
			want: Instruction{
				Source:      Source{Type: SourceReg, Reg: RegR2},
				Destination: Destination{Type: DestReg, Reg: RegR1},
				Increment:   false,
				Effect:      EffectAlways,
				Operation:   OpAdd,
				Immediate:   0,
				Offset:      0,
			},
		},
		{
			name: "copy r1 mem to r2 mem (push)",
			raw:  0x12C00000, // r2 push <- copy r1 (src=1(r1 mem), dest=2(r2 mem), M=1, effect=4, op=0)
			want: Instruction{
				Source:      Source{Type: SourceMem, Reg: RegR1},
				Destination: Destination{Type: DestMem, Reg: RegR2},
				Increment:   true,
				Effect:      EffectAlways,
				Operation:   OpCopy,
				Immediate:   0,
				Offset:      0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.raw)
			if err != nil {
				t.Fatalf("Decode() error = %v", err)
			}

			if got.Source.Type != tt.want.Source.Type {
				t.Errorf("Source.Type = %v, want %v", got.Source.Type, tt.want.Source.Type)
			}
			if got.Source.Reg != tt.want.Source.Reg {
				t.Errorf("Source.Reg = %v, want %v", got.Source.Reg, tt.want.Source.Reg)
			}
			if got.Destination.Type != tt.want.Destination.Type {
				t.Errorf("Destination.Type = %v, want %v", got.Destination.Type, tt.want.Destination.Type)
			}
			if got.Destination.Reg != tt.want.Destination.Reg {
				t.Errorf("Destination.Reg = %v, want %v", got.Destination.Reg, tt.want.Destination.Reg)
			}
			if got.Increment != tt.want.Increment {
				t.Errorf("Increment = %v, want %v", got.Increment, tt.want.Increment)
			}
			if got.Effect != tt.want.Effect {
				t.Errorf("Effect = %v, want %v", got.Effect, tt.want.Effect)
			}
			if got.Operation != tt.want.Operation {
				t.Errorf("Operation = %v, want %v", got.Operation, tt.want.Operation)
			}
			if got.Immediate != tt.want.Immediate {
				t.Errorf("Immediate = %d, want %d", got.Immediate, tt.want.Immediate)
			}
			if got.Offset != tt.want.Offset {
				t.Errorf("Offset = %d, want %d", got.Offset, tt.want.Offset)
			}
		})
	}
}

func TestDecode_MemoryDestination(t *testing.T) {
	// Memory destination: r1/5 <- add val/10
	// Format: SSSSDDDD MEEEAAAA OOOOIIII IIIIIIII
	// Source = val (0100), Dest = r1 mem (0001), Inc = 0, Effect = always (100), Op = add (1010)
	// Offset = 5, Immediate = 10
	raw := uint32(0x414A500A) // 0100 0001 0 100 1010 0101 000000001010

	inst, err := Decode(raw)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if inst.Destination.Type != DestMem {
		t.Errorf("Destination.Type = %v, want DestMem", inst.Destination.Type)
	}
	if inst.Offset != 5 {
		t.Errorf("Offset = %d, want 5", inst.Offset)
	}
	if inst.Immediate != 10 {
		t.Errorf("Immediate = %d, want 10", inst.Immediate)
	}
}

func TestDecode_ConditionalEffects(t *testing.T) {
	tests := []struct {
		name   string
		effect Effect
		bits   uint16
	}{
		{"zero", EffectZero, 0x0000},
		{"not zero", EffectNotZero, 0x0010},
		{"negative", EffectNegative, 0x0020},
		{"positive", EffectPositive, 0x0030},
		{"always", EffectAlways, 0x0040},
		{"overflow", EffectOverflow, 0x0050},
		{"interrupt", EffectInterrupt, 0x0060},
		{"flags only", EffectFlagsOnly, 0x0070},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create instruction: r1 <- copy val/0 with specific effect
			raw := uint32(0x5400|tt.bits)<<16 | 0x0000

			inst, err := Decode(raw)
			if err != nil {
				t.Fatalf("Decode() error = %v", err)
			}

			if inst.Effect != tt.effect {
				t.Errorf("Effect = %v, want %v", inst.Effect, tt.effect)
			}
		})
	}
}

func TestDecode_HaltInstruction(t *testing.T) {
	// HALT: copy pc to pc with immediate 0
	raw := uint32(0x00400000)

	inst, err := Decode(raw)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if !inst.IsHalt() {
		t.Error("Expected HALT instruction")
	}
}

func TestDecode_SignedImmediate(t *testing.T) {
	// Test negative immediate with memory destination (12-bit signed)
	// r1/0 <- add val/-5
	// -5 in 12-bit two's complement: 0xFFB (1111 1111 1011)
	raw := uint32(0x414A0FFB)

	inst, err := Decode(raw)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if inst.Immediate != -5 {
		t.Errorf("Immediate = %d, want -5", inst.Immediate)
	}
}
