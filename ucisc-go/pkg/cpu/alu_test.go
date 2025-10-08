package cpu

import "testing"

func TestALU_Copy(t *testing.T) {
	alu := NewALU()

	tests := []struct {
		name   string
		value  uint16
		signed bool
		zero   bool
		neg    bool
	}{
		{"zero", 0, false, true, false},
		{"positive", 42, false, false, false},
		{"negative unsigned", 0x8000, false, false, false}, // In unsigned mode, no negative flag
		{"max unsigned", 0xFFFF, false, false, false},      // In unsigned mode, no negative flag
		{"negative signed", 0x8000, true, false, true},     // In signed mode, MSB set = negative
		{"max signed", 0xFFFF, true, false, true},          // In signed mode, MSB set = negative
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alu.SetSigned(tt.signed)
			result := alu.Execute(OpCopy, 0, tt.value, false)
			if result.Value != tt.value {
				t.Errorf("got value %d, want %d", result.Value, tt.value)
			}
			if result.Zero != tt.zero {
				t.Errorf("got zero %v, want %v", result.Zero, tt.zero)
			}
			if result.Negative != tt.neg {
				t.Errorf("got negative %v, want %v", result.Negative, tt.neg)
			}
		})
	}
}

func TestALU_Add(t *testing.T) {
	alu := NewALU()

	tests := []struct {
		name     string
		dest     uint16
		src      uint16
		signed   bool
		expected uint16
		carry    bool
		overflow bool
	}{
		{"simple add", 5, 3, false, 8, false, false},
		{"with carry", 0xFFFF, 1, false, 0, true, true},
		{"signed positive", 5, 3, true, 8, false, false},
		{"signed overflow", 0x7FFF, 1, true, 0x8000, false, true},
		{"signed negative", 0xFFFE, 1, true, 0xFFFF, false, false}, // -2 + 1 = -1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alu.SetSigned(tt.signed)
			result := alu.Execute(OpAdd, tt.dest, tt.src, false)

			if result.Value != tt.expected {
				t.Errorf("got value 0x%04X, want 0x%04X", result.Value, tt.expected)
			}
			if result.Carry != tt.carry {
				t.Errorf("got carry %v, want %v", result.Carry, tt.carry)
			}
			if result.Overflow != tt.overflow {
				t.Errorf("got overflow %v, want %v", result.Overflow, tt.overflow)
			}
		})
	}
}

func TestALU_Sub(t *testing.T) {
	alu := NewALU()

	tests := []struct {
		name     string
		dest     uint16
		src      uint16
		signed   bool
		expected uint16
		borrow   bool
	}{
		{"simple sub", 10, 3, false, 7, false},
		{"with borrow", 3, 10, false, 0xFFF9, true},
		{"signed", 5, 3, true, 2, false},
		{"zero result", 5, 5, false, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alu.SetSigned(tt.signed)
			result := alu.Execute(OpSub, tt.dest, tt.src, false)

			if result.Value != tt.expected {
				t.Errorf("got value 0x%04X, want 0x%04X", result.Value, tt.expected)
			}
			if result.Carry != tt.borrow {
				t.Errorf("got borrow %v, want %v", result.Carry, tt.borrow)
			}
		})
	}
}

func TestALU_BitwiseOps(t *testing.T) {
	alu := NewALU()

	tests := []struct {
		name     string
		op       Operation
		dest     uint16
		src      uint16
		expected uint16
	}{
		{"AND", OpAnd, 0xFF0F, 0x00FF, 0x000F},
		{"OR", OpOr, 0xFF00, 0x00FF, 0xFFFF},
		{"XOR", OpXor, 0xFF00, 0xFF0F, 0x000F},
		{"INV", OpInv, 0, 0xFF00, 0x00FF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := alu.Execute(tt.op, tt.dest, tt.src, false)
			if result.Value != tt.expected {
				t.Errorf("got 0x%04X, want 0x%04X", result.Value, tt.expected)
			}
		})
	}
}

func TestALU_Shift(t *testing.T) {
	alu := NewALU()

	tests := []struct {
		name     string
		op       Operation
		value    uint16
		shift    uint16
		signed   bool
		expected uint16
	}{
		{"SHL by 1", OpShl, 0x0001, 1, false, 0x0002},
		{"SHL by 4", OpShl, 0x000F, 4, false, 0x00F0},
		{"SHR unsigned", OpShr, 0xF000, 4, false, 0x0F00},
		{"SHR signed neg", OpShr, 0x8000, 1, true, 0xC000},
		{"SHR signed pos", OpShr, 0x4000, 1, true, 0x2000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alu.SetSigned(tt.signed)
			result := alu.Execute(tt.op, tt.value, tt.shift, false)
			if result.Value != tt.expected {
				t.Errorf("got 0x%04X, want 0x%04X", result.Value, tt.expected)
			}
		})
	}
}

func TestALU_ByteOps(t *testing.T) {
	alu := NewALU()

	tests := []struct {
		name     string
		op       Operation
		value    uint16
		expected uint16
	}{
		{"SWAP", OpSwap, 0x1234, 0x3412},
		{"MSB", OpMsb, 0x1234, 0x1200},
		{"LSB", OpLsb, 0x1234, 0x0034},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := alu.Execute(tt.op, 0, tt.value, false)
			if result.Value != tt.expected {
				t.Errorf("got 0x%04X, want 0x%04X", result.Value, tt.expected)
			}
		})
	}
}

func TestALU_Mult(t *testing.T) {
	alu := NewALU()

	tests := []struct {
		name     string
		dest     uint16
		src      uint16
		signed   bool
		expected uint16
		overflow bool
	}{
		{"simple", 5, 3, false, 15, false},
		{"overflow unsigned", 0x100, 0x100, false, 0, true},
		{"signed positive", 5, 3, true, 15, false},
		{"signed negative", 0xFFFB, 3, true, 0xFFF1, false}, // -5 * 3 = -15
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alu.SetSigned(tt.signed)
			result := alu.Execute(OpMult, tt.dest, tt.src, false)

			if result.Value != tt.expected {
				t.Errorf("got value 0x%04X, want 0x%04X", result.Value, tt.expected)
			}
			if result.Overflow != tt.overflow {
				t.Errorf("got overflow %v, want %v", result.Overflow, tt.overflow)
			}
		})
	}
}
