package cpu

import (
	"testing"
)

// TestCPU_AllRegisters verifies all 6 registers can be read and written
func TestCPU_AllRegisters(t *testing.T) {
	c := New(1024)

	tests := []struct {
		name  string
		reg   Register
		value uint16
		get   func() uint16
		set   func(uint16)
	}{
		{"r1", RegR1, 0x1111, c.GetR1, func(v uint16) { c.setRegister(RegR1, v) }},
		{"r2", RegR2, 0x2222, c.GetR2, func(v uint16) { c.setRegister(RegR2, v) }},
		{"r3", RegR3, 0x3333, c.GetR3, func(v uint16) { c.setRegister(RegR3, v) }},
		{"r5", RegR5, 0x5555, c.GetR5, func(v uint16) { c.setRegister(RegR5, v) }},
		{"r6", RegR6, 0x6666, c.GetR6, func(v uint16) { c.setRegister(RegR6, v) }},
		{"r7", RegR7, 0x7777, c.GetR7, func(v uint16) { c.setRegister(RegR7, v) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the register
			tt.set(tt.value)

			// Read it back
			got := tt.get()
			if got != tt.value {
				t.Errorf("%s: got 0x%04X, want 0x%04X", tt.name, got, tt.value)
			}

			// Verify other registers weren't affected
			for _, other := range tests {
				if other.name == tt.name {
					continue
				}
				if other.get() == tt.value && other.name != tt.name {
					t.Errorf("%s: writing to %s affected %s", tt.name, tt.name, other.name)
				}
			}
		})
	}
}

// TestCPU_RegisterIsolation ensures registers are truly independent
func TestCPU_RegisterIsolation(t *testing.T) {
	c := New(1024)

	// Set all registers to unique values
	c.setRegister(RegR1, 0xAA11)
	c.setRegister(RegR2, 0xBB22)
	c.setRegister(RegR3, 0xCC33)
	c.setRegister(RegR5, 0xDD55)
	c.setRegister(RegR6, 0xEE66)
	c.setRegister(RegR7, 0xFF77)

	// Verify all values are correct
	if got := c.GetR1(); got != 0xAA11 {
		t.Errorf("r1 = 0x%04X, want 0xAA11", got)
	}
	if got := c.GetR2(); got != 0xBB22 {
		t.Errorf("r2 = 0x%04X, want 0xBB22", got)
	}
	if got := c.GetR3(); got != 0xCC33 {
		t.Errorf("r3 = 0x%04X, want 0xCC33", got)
	}
	if got := c.GetR5(); got != 0xDD55 {
		t.Errorf("r5 = 0x%04X, want 0xDD55", got)
	}
	if got := c.GetR6(); got != 0xEE66 {
		t.Errorf("r6 = 0x%04X, want 0xEE66", got)
	}
	if got := c.GetR7(); got != 0xFF77 {
		t.Errorf("r7 = 0x%04X, want 0xFF77", got)
	}
}

// TestCPU_RegisterALUOperations tests ALU operations between all register pairs
func TestCPU_RegisterALUOperations(t *testing.T) {
	tests := []struct {
		name     string
		destReg  Register
		srcReg   Register
		destVal  uint16
		srcVal   uint16
		op       Operation
		expected uint16
	}{
		// Test r1 <- r2
		{"r1_copy_r2", RegR1, RegR2, 0, 100, OpCopy, 100},
		{"r1_add_r2", RegR1, RegR2, 50, 30, OpAdd, 80},
		{"r1_sub_r2", RegR1, RegR2, 100, 30, OpSub, 70},
		{"r1_and_r2", RegR1, RegR2, 0xFF00, 0x00FF, OpAnd, 0},
		{"r1_or_r2", RegR1, RegR2, 0xFF00, 0x00FF, OpOr, 0xFFFF},

		// Test r3 <- r5
		{"r3_copy_r5", RegR3, RegR5, 0, 200, OpCopy, 200},
		{"r3_add_r5", RegR3, RegR5, 150, 50, OpAdd, 200},
		{"r3_xor_r5", RegR3, RegR5, 0xAAAA, 0x5555, OpXor, 0xFFFF},

		// Test r6 <- r7
		{"r6_copy_r7", RegR6, RegR7, 0, 300, OpCopy, 300},
		{"r6_add_r7", RegR6, RegR7, 250, 50, OpAdd, 300},

		// Test r2 <- r3
		{"r2_copy_r3", RegR2, RegR3, 0, 400, OpCopy, 400},

		// Test r5 <- r6
		{"r5_copy_r6", RegR5, RegR6, 0, 500, OpCopy, 500},

		// Test r7 <- r1
		{"r7_copy_r1", RegR7, RegR1, 0, 600, OpCopy, 600},

		// Cross-bank operations (r1 <- r5, r2 <- r6, r3 <- r7)
		{"r1_add_r5", RegR1, RegR5, 100, 200, OpAdd, 300},
		{"r2_add_r6", RegR2, RegR6, 150, 150, OpAdd, 300},
		{"r3_add_r7", RegR3, RegR7, 50, 250, OpAdd, 300},

		// Reverse cross-bank (r5 <- r1, r6 <- r2, r7 <- r3)
		{"r5_sub_r1", RegR5, RegR1, 500, 200, OpSub, 300},
		{"r6_sub_r2", RegR6, RegR2, 400, 100, OpSub, 300},
		{"r7_sub_r3", RegR7, RegR3, 350, 50, OpSub, 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(1024)

			// Set up initial values
			c.setRegister(tt.destReg, tt.destVal)
			c.setRegister(tt.srcReg, tt.srcVal)

			// Build instruction: dest <- op src
			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: tt.destReg},
				Source:      Source{Type: SourceReg, Reg: tt.srcReg},
				Operation:   tt.op,
				Effect:      EffectAlways,
				Immediate:   0,
				Offset:      0,
				Increment:   false,
			}

			// Execute
			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			// Check result
			got := c.getRegister(tt.destReg)
			if got != tt.expected {
				t.Errorf("%s = 0x%04X, want 0x%04X", tt.name, got, tt.expected)
			}

			// Verify source register wasn't modified
			srcGot := c.getRegister(tt.srcReg)
			if srcGot != tt.srcVal {
				t.Errorf("source register %v was modified: got 0x%04X, want 0x%04X",
					tt.srcReg, srcGot, tt.srcVal)
			}
		})
	}
}

// TestCPU_RegisterEdgeValues tests edge case values in all registers
func TestCPU_RegisterEdgeValues(t *testing.T) {
	edgeValues := []struct {
		name  string
		value uint16
	}{
		{"zero", 0x0000},
		{"one", 0x0001},
		{"max_positive", 0x7FFF},
		{"min_negative", 0x8000},
		{"minus_one", 0xFFFF},
		{"alternating_bits_1", 0xAAAA},
		{"alternating_bits_2", 0x5555},
		{"all_ones", 0xFFFF},
	}

	registers := []struct {
		name string
		reg  Register
		get  func(*CPU) uint16
	}{
		{"r1", RegR1, (*CPU).GetR1},
		{"r2", RegR2, (*CPU).GetR2},
		{"r3", RegR3, (*CPU).GetR3},
		{"r5", RegR5, (*CPU).GetR5},
		{"r6", RegR6, (*CPU).GetR6},
		{"r7", RegR7, (*CPU).GetR7},
	}

	for _, reg := range registers {
		for _, edge := range edgeValues {
			t.Run(reg.name+"_"+edge.name, func(t *testing.T) {
				c := New(1024)
				c.setRegister(reg.reg, edge.value)

				got := reg.get(c)
				if got != edge.value {
					t.Errorf("%s with %s: got 0x%04X, want 0x%04X",
						reg.name, edge.name, got, edge.value)
				}
			})
		}
	}
}

// TestCPU_RegisterMemoryOperations tests reading from and writing to memory via each register
func TestCPU_RegisterMemoryOperations(t *testing.T) {
	registers := []struct {
		name string
		reg  Register
	}{
		{"r1", RegR1},
		{"r2", RegR2},
		{"r3", RegR3},
		{"r5", RegR5},
		{"r6", RegR6},
		{"r7", RegR7},
	}

	for _, tt := range registers {
		t.Run(tt.name+"_memory_write", func(t *testing.T) {
			c := New(1024)

			// Set register to point to address 100
			c.setRegister(tt.reg, 100)

			// Write value 0x1234 to memory at register address
			inst := &Instruction{
				Destination: Destination{Type: DestMem, Reg: tt.reg},
				Source:      Source{Type: SourceVal, Reg: RegR1},
				Operation:   OpCopy,
				Effect:      EffectAlways,
				Immediate:   0x1234,
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			// Read back from memory
			got := c.memory.ReadLocal(100)
			if got != 0x1234 {
				t.Errorf("memory[100] = 0x%04X, want 0x1234", got)
			}
		})

		t.Run(tt.name+"_memory_read", func(t *testing.T) {
			c := New(1024)

			// Set register to point to address 200
			c.setRegister(tt.reg, 200)

			// Write test value to memory
			c.memory.WriteLocal(200, 0x1234)

			// Read from memory via register into r1
			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: RegR1},
				Source:      Source{Type: SourceMem, Reg: tt.reg},
				Operation:   OpCopy,
				Effect:      EffectAlways,
				Immediate:   0,
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			// Check r1 has the value
			if got := c.GetR1(); got != 0x1234 {
				t.Errorf("r1 = 0x%04X, want 0x1234", got)
			}
		})

		t.Run(tt.name+"_memory_with_offset", func(t *testing.T) {
			c := New(1024)

			// Set register to point to base address 300
			c.setRegister(tt.reg, 300)

			// Write to memory at base + 5
			c.memory.WriteLocal(305, 0xDEAD)

			// Read from memory with offset
			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: RegR1},
				Source:      Source{Type: SourceMem, Reg: tt.reg},
				Operation:   OpCopy,
				Effect:      EffectAlways,
				Immediate:   5, // offset
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			if got := c.GetR1(); got != 0xDEAD {
				t.Errorf("r1 = 0x%04X, want 0xDEAD", got)
			}
		})
	}
}

// TestCPU_RegisterOverflow tests overflow behavior in all registers
func TestCPU_RegisterOverflow(t *testing.T) {
	registers := []struct {
		name string
		reg  Register
		get  func(*CPU) uint16
	}{
		{"r1", RegR1, (*CPU).GetR1},
		{"r2", RegR2, (*CPU).GetR2},
		{"r3", RegR3, (*CPU).GetR3},
		{"r5", RegR5, (*CPU).GetR5},
		{"r6", RegR6, (*CPU).GetR6},
		{"r7", RegR7, (*CPU).GetR7},
	}

	for _, tt := range registers {
		t.Run(tt.name+"_add_overflow", func(t *testing.T) {
			c := New(1024)

			// Set register to max value
			c.setRegister(tt.reg, 0xFFFF)

			// Add 1 (should overflow to 0)
			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: tt.reg},
				Source:      Source{Type: SourceVal, Reg: RegR1},
				Operation:   OpAdd,
				Effect:      EffectAlways,
				Immediate:   1,
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			got := tt.get(c)
			if got != 0 {
				t.Errorf("%s after overflow: got 0x%04X, want 0x0000", tt.name, got)
			}

			// Check carry flag is set
			if !c.flags.Carry() {
				t.Errorf("%s overflow: carry flag not set", tt.name)
			}
		})

		t.Run(tt.name+"_sub_underflow", func(t *testing.T) {
			c := New(1024)

			// Set register to 0
			c.setRegister(tt.reg, 0)

			// Subtract 1 (should underflow to 0xFFFF)
			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: tt.reg},
				Source:      Source{Type: SourceVal, Reg: RegR1},
				Operation:   OpSub,
				Effect:      EffectAlways,
				Immediate:   1,
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			got := tt.get(c)
			if got != 0xFFFF {
				t.Errorf("%s after underflow: got 0x%04X, want 0xFFFF", tt.name, got)
			}
		})
	}
}

// TestCPU_RegisterShiftRotate tests shift/rotate operations on all registers
func TestCPU_RegisterShiftRotate(t *testing.T) {
	registers := []struct {
		name string
		reg  Register
		get  func(*CPU) uint16
	}{
		{"r1", RegR1, (*CPU).GetR1},
		{"r2", RegR2, (*CPU).GetR2},
		{"r3", RegR3, (*CPU).GetR3},
		{"r5", RegR5, (*CPU).GetR5},
		{"r6", RegR6, (*CPU).GetR6},
		{"r7", RegR7, (*CPU).GetR7},
	}

	for _, tt := range registers {
		t.Run(tt.name+"_shift_left", func(t *testing.T) {
			c := New(1024)
			c.setRegister(tt.reg, 0x0001)

			// Shift left by 4
			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: tt.reg},
				Source:      Source{Type: SourceVal, Reg: RegR1},
				Operation:   OpShl,
				Effect:      EffectAlways,
				Immediate:   4,
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			got := tt.get(c)
			if got != 0x0010 {
				t.Errorf("%s shift left: got 0x%04X, want 0x0010", tt.name, got)
			}
		})

		t.Run(tt.name+"_shift_right", func(t *testing.T) {
			c := New(1024)
			c.setRegister(tt.reg, 0x1000)

			// Shift right by 4
			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: tt.reg},
				Source:      Source{Type: SourceVal, Reg: RegR1},
				Operation:   OpShr,
				Effect:      EffectAlways,
				Immediate:   4,
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			got := tt.get(c)
			if got != 0x0100 {
				t.Errorf("%s shift right: got 0x%04X, want 0x0100", tt.name, got)
			}
		})
	}
}

// TestCPU_RegisterMultiplication tests multiplication with all registers
func TestCPU_RegisterMultiplication(t *testing.T) {
	tests := []struct {
		name     string
		destReg  Register
		srcReg   Register
		destVal  uint16
		srcVal   uint16
		expected uint16
	}{
		{"r1_mult_r2", RegR1, RegR2, 10, 20, 200},
		{"r2_mult_r3", RegR2, RegR3, 15, 15, 225},
		{"r3_mult_r5", RegR3, RegR5, 100, 2, 200},
		{"r5_mult_r6", RegR5, RegR6, 50, 4, 200},
		{"r6_mult_r7", RegR6, RegR7, 25, 8, 200},
		{"r7_mult_r1", RegR7, RegR1, 20, 10, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(1024)
			c.setRegister(tt.destReg, tt.destVal)
			c.setRegister(tt.srcReg, tt.srcVal)

			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: tt.destReg},
				Source:      Source{Type: SourceReg, Reg: tt.srcReg},
				Operation:   OpMult,
				Effect:      EffectAlways,
				Immediate:   0,
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			got := c.getRegister(tt.destReg)
			if got != tt.expected {
				t.Errorf("%s: got %d, want %d", tt.name, got, tt.expected)
			}
		})
	}
}

// TestCPU_RegisterAsStackPointer tests using each register as a stack pointer
func TestCPU_RegisterAsStackPointer(t *testing.T) {
	registers := []struct {
		name string
		reg  Register
	}{
		{"r1", RegR1},
		{"r2", RegR2},
		{"r3", RegR3},
		{"r5", RegR5},
		{"r6", RegR6},
		{"r7", RegR7},
	}

	for _, tt := range registers {
		t.Run(tt.name+"_push", func(t *testing.T) {
			c := New(1024)

			// Set register as stack pointer (grows down)
			c.setRegister(tt.reg, 500)

			// Push value 0x5678
			inst := &Instruction{
				Destination: Destination{Type: DestMem, Reg: tt.reg},
				Source:      Source{Type: SourceVal, Reg: RegR1},
				Operation:   OpCopy,
				Effect:      EffectAlways,
				Immediate:   0x5678,
				Offset:      0,
				Increment:   true, // push
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			// Stack pointer should have decremented
			if got := c.getRegister(tt.reg); got != 499 {
				t.Errorf("%s after push: got %d, want 499", tt.name, got)
			}

			// Value should be at decremented address
			if got := c.memory.ReadLocal(499); got != 0x5678 {
				t.Errorf("memory[499] = 0x%04X, want 0x5678", got)
			}
		})

		t.Run(tt.name+"_pop", func(t *testing.T) {
			c := New(1024)

			// Set register as stack pointer
			c.setRegister(tt.reg, 400)

			// Put value on "stack"
			c.memory.WriteLocal(400, 0xCAFE)

			// Choose a different destination register (not the stack pointer)
			destReg := RegR7
			if tt.reg == RegR7 {
				destReg = RegR1
			}

			// Pop value into destination register
			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: destReg},
				Source:      Source{Type: SourceMem, Reg: tt.reg},
				Operation:   OpCopy,
				Effect:      EffectAlways,
				Immediate:   0,
				Offset:      0,
				Increment:   true, // pop
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			// Check value was read
			if got := c.getRegister(destReg); got != 0xCAFE {
				t.Errorf("%v = 0x%04X, want 0xCAFE", destReg, got)
			}

			// Stack pointer should have incremented
			if got := c.getRegister(tt.reg); got != 401 {
				t.Errorf("%s after pop: got %d, want 401", tt.name, got)
			}
		})
	}
}

// TestCPU_RegisterMultiplicationOverflow tests multiplication overflow and MSW operations
func TestCPU_RegisterMultiplicationOverflow(t *testing.T) {
	tests := []struct {
		name        string
		reg1        Register
		reg2        Register
		val1        uint16
		val2        uint16
		expectedLSW uint16
		expectedMSW uint16
		signed      bool
	}{
		// Unsigned overflow cases
		{"r1_mult_r2_overflow", RegR1, RegR2, 0xFFFF, 0xFFFF, 0x0001, 0xFFFE, false},
		{"r3_mult_r5_overflow", RegR3, RegR5, 256, 256, 0x0000, 0x0001, false},
		{"r5_mult_r6_overflow", RegR5, RegR6, 1000, 1000, 0x4240, 0x000F, false},
		{"r6_mult_r7_large", RegR6, RegR7, 0x8000, 2, 0x0000, 0x0001, false},

		// Signed overflow cases (positive * positive)
		{"r1_mult_r2_signed_pos", RegR1, RegR2, 100, 100, 10000, 0, true},
		{"r2_mult_r3_signed_overflow", RegR2, RegR3, 200, 200, 40000, 0, true},

		// Signed cases (negative * positive)
		{"r3_mult_r5_signed_neg", RegR3, RegR5, 0xFFFF, 100, 0xFF9C, 0xFFFF, true}, // -1 * 100 = -100

		// Small values (no overflow)
		{"r5_mult_r6_small", RegR5, RegR6, 10, 10, 100, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_lsw", func(t *testing.T) {
			c := New(1024)
			if tt.signed {
				c.flags = c.flags.Set(FlagSigned)
			}

			c.setRegister(tt.reg1, tt.val1)
			c.setRegister(tt.reg2, tt.val2)

			// Multiply to get LSW
			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: tt.reg1},
				Source:      Source{Type: SourceReg, Reg: tt.reg2},
				Operation:   OpMult,
				Effect:      EffectAlways,
				Immediate:   0,
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			got := c.getRegister(tt.reg1)
			if got != tt.expectedLSW {
				t.Errorf("%s LSW: got 0x%04X, want 0x%04X", tt.name, got, tt.expectedLSW)
			}
		})

		t.Run(tt.name+"_msw", func(t *testing.T) {
			c := New(1024)
			if tt.signed {
				c.flags = c.flags.Set(FlagSigned)
			}

			c.setRegister(tt.reg1, tt.val1)
			c.setRegister(tt.reg2, tt.val2)

			// Multiply to get MSW (upper 16 bits)
			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: tt.reg1},
				Source:      Source{Type: SourceReg, Reg: tt.reg2},
				Operation:   OpMultMsw,
				Effect:      EffectAlways,
				Immediate:   0,
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			got := c.getRegister(tt.reg1)
			if got != tt.expectedMSW {
				t.Errorf("%s MSW: got 0x%04X, want 0x%04X", tt.name, got, tt.expectedMSW)
			}
		})
	}
}

// TestCPU_RegisterMultiplication32Bit tests getting full 32-bit multiplication results
func TestCPU_RegisterMultiplication32Bit(t *testing.T) {
	tests := []struct {
		name       string
		val1       uint16
		val2       uint16
		lswReg     Register
		mswReg     Register
		expected32 uint32
	}{
		{"max_unsigned", 0xFFFF, 0xFFFF, RegR1, RegR2, 0xFFFE0001},
		{"256_times_256", 256, 256, RegR3, RegR5, 65536},
		{"1000_times_1000", 1000, 1000, RegR5, RegR6, 1000000},
		{"small_values", 10, 20, RegR6, RegR7, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(1024)

			// Store values in temporary registers (r7 and r3)
			c.setRegister(RegR7, tt.val1)
			c.setRegister(RegR3, tt.val2)

			// Get LSW: lswReg = r7 * r3
			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: tt.lswReg},
				Source:      Source{Type: SourceReg, Reg: RegR7},
				Operation:   OpMult,
				Effect:      EffectAlways,
				Immediate:   0,
				Offset:      0,
				Increment:   false,
			}
			inst.Source.Reg = RegR3
			inst.Destination.Reg = RegR7
			if err := c.execute(inst); err != nil {
				t.Fatalf("mult failed: %v", err)
			}
			lsw := c.GetR7()

			// Reset source values
			c.setRegister(RegR7, tt.val1)
			c.setRegister(RegR3, tt.val2)

			// Get MSW: mswReg = multmsw(r7, r3)
			inst.Operation = OpMultMsw
			if err := c.execute(inst); err != nil {
				t.Fatalf("multmsw failed: %v", err)
			}
			msw := c.GetR7()

			// Reconstruct 32-bit result
			result32 := (uint32(msw) << 16) | uint32(lsw)

			if result32 != tt.expected32 {
				t.Errorf("%s: got 0x%08X, want 0x%08X (LSW=0x%04X, MSW=0x%04X)",
					tt.name, result32, tt.expected32, lsw, msw)
			}
		})
	}
}

// TestCPU_RegisterMultiplicationSigned tests signed multiplication with all registers
func TestCPU_RegisterMultiplicationSigned(t *testing.T) {
	tests := []struct {
		name     string
		reg1     Register
		reg2     Register
		val1     int16
		val2     int16
		expected int16
	}{
		{"r1_neg_times_pos", RegR1, RegR2, -10, 20, -200},
		{"r2_pos_times_neg", RegR2, RegR3, 15, -10, -150},
		{"r3_neg_times_neg", RegR3, RegR5, -20, -5, 100},
		{"r5_large_pos", RegR5, RegR6, 100, 100, 10000},
		{"r6_large_neg", RegR6, RegR7, -100, 100, -10000},
		{"r7_minus_one", RegR7, RegR1, -1, 50, -50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(1024)

			// Set signed mode
			c.flags = c.flags.Set(FlagSigned)

			c.setRegister(tt.reg1, uint16(tt.val1))
			c.setRegister(tt.reg2, uint16(tt.val2))

			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: tt.reg1},
				Source:      Source{Type: SourceReg, Reg: tt.reg2},
				Operation:   OpMult,
				Effect:      EffectAlways,
				Immediate:   0,
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			got := int16(c.getRegister(tt.reg1))
			if got != tt.expected {
				t.Errorf("%s: got %d, want %d", tt.name, got, tt.expected)
			}
		})
	}
}

// TestCPU_RegisterMultiplicationEdgeCases tests edge cases for multiplication
func TestCPU_RegisterMultiplicationEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		reg         Register
		val1        uint16
		val2        uint16
		expectedLSW uint16
		expectedMSW uint16
	}{
		{"zero_times_anything", RegR1, 0, 12345, 0, 0},
		{"one_times_value", RegR2, 1, 0x5678, 0x5678, 0},
		{"max_times_one", RegR3, 0xFFFF, 1, 0xFFFF, 0},
		{"max_times_two", RegR5, 0xFFFF, 2, 0xFFFE, 0x0001},
		{"power_of_two", RegR6, 0x100, 0x100, 0x0000, 0x0001},
		{"half_max_squared", RegR7, 0x8000, 0x8000, 0x0000, 0x4000},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_lsw", func(t *testing.T) {
			c := New(1024)
			c.setRegister(tt.reg, tt.val1)
			c.setRegister(RegR7, tt.val2) // Use r7 as source

			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: tt.reg},
				Source:      Source{Type: SourceReg, Reg: RegR7},
				Operation:   OpMult,
				Effect:      EffectAlways,
				Immediate:   0,
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			got := c.getRegister(tt.reg)
			if got != tt.expectedLSW {
				t.Errorf("%s: got 0x%04X, want 0x%04X", tt.name, got, tt.expectedLSW)
			}
		})

		t.Run(tt.name+"_msw", func(t *testing.T) {
			c := New(1024)
			c.setRegister(tt.reg, tt.val1)
			c.setRegister(RegR7, tt.val2)

			inst := &Instruction{
				Destination: Destination{Type: DestReg, Reg: tt.reg},
				Source:      Source{Type: SourceReg, Reg: RegR7},
				Operation:   OpMultMsw,
				Effect:      EffectAlways,
				Immediate:   0,
				Offset:      0,
				Increment:   false,
			}

			if err := c.execute(inst); err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			got := c.getRegister(tt.reg)
			if got != tt.expectedMSW {
				t.Errorf("%s: got 0x%04X, want 0x%04X", tt.name, got, tt.expectedMSW)
			}
		})
	}
}

// TestCPU_RegisterCopyChain tests copying values through a chain of registers
func TestCPU_RegisterCopyChain(t *testing.T) {
	c := New(1024)

	// Chain: r1 -> r2 -> r3 -> r5 -> r6 -> r7
	testValue := uint16(0xABCD)

	// Start with r1
	c.setRegister(RegR1, testValue)

	// Copy chain
	copies := []struct {
		from Register
		to   Register
	}{
		{RegR1, RegR2},
		{RegR2, RegR3},
		{RegR3, RegR5},
		{RegR5, RegR6},
		{RegR6, RegR7},
	}

	for i, cp := range copies {
		inst := &Instruction{
			Destination: Destination{Type: DestReg, Reg: cp.to},
			Source:      Source{Type: SourceReg, Reg: cp.from},
			Operation:   OpCopy,
			Effect:      EffectAlways,
			Immediate:   0,
			Offset:      0,
			Increment:   false,
		}

		if err := c.execute(inst); err != nil {
			t.Fatalf("copy %d failed: %v", i, err)
		}

		// Verify destination has the value
		got := c.getRegister(cp.to)
		if got != testValue {
			t.Errorf("after copy %d: %v = 0x%04X, want 0x%04X", i, cp.to, got, testValue)
		}
	}

	// Verify all registers have the value
	if c.GetR1() != testValue || c.GetR2() != testValue || c.GetR3() != testValue ||
		c.GetR5() != testValue || c.GetR6() != testValue || c.GetR7() != testValue {
		t.Errorf("not all registers have value 0x%04X: r1=0x%04X r2=0x%04X r3=0x%04X r5=0x%04X r6=0x%04X r7=0x%04X",
			testValue, c.GetR1(), c.GetR2(), c.GetR3(), c.GetR5(), c.GetR6(), c.GetR7())
	}
}
