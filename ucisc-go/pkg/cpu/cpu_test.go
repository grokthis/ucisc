package cpu

import "testing"

func TestCPU_BasicExecution(t *testing.T) {
	cpu := New(1024)

	// Program: r1 <- copy val/5
	// src=val(4), dest=r1 reg(5), effect=always(4), op=copy(0)
	program := []uint16{
		0x4540, 0x0005, // r1 <- copy val/5
	}

	err := cpu.LoadProgram(program, 0)
	if err != nil {
		t.Fatalf("LoadProgram() error = %v", err)
	}

	err = cpu.Step()
	if err != nil {
		t.Fatalf("Step() error = %v", err)
	}

	if cpu.GetR1() != 5 {
		t.Errorf("R1 = %d, want 5", cpu.GetR1())
	}

	if cpu.GetPC() != 2 {
		t.Errorf("PC = %d, want 2", cpu.GetPC())
	}
}

func TestCPU_AddOperation(t *testing.T) {
	cpu := New(1024)

	// Program:
	//   r1 <- copy val/10
	//   r1 <- add val/5
	program := []uint16{
		0x4540, 0x000A, // r1 <- copy val/10
		0x454A, 0x0005, // r1 <- add val/5
	}

	cpu.LoadProgram(program, 0)

	cpu.Step() // Execute first instruction
	if cpu.GetR1() != 10 {
		t.Errorf("After first instruction: R1 = %d, want 10", cpu.GetR1())
	}

	cpu.Step() // Execute second instruction
	if cpu.GetR1() != 15 {
		t.Errorf("After second instruction: R1 = %d, want 15", cpu.GetR1())
	}

	// Check flags (15 is non-zero, positive)
	if cpu.GetFlags().Zero() {
		t.Error("Zero flag should not be set")
	}
	if cpu.GetFlags().Negative() {
		t.Error("Negative flag should not be set")
	}
}

func TestCPU_MemoryOperations(t *testing.T) {
	cpu := New(1024)

	// Program:
	//   r1 <- copy val/100  (set r1 to address 100)
	//   r1/0 <- copy val/42 (write 42 to address 100)
	//   r2 <- copy r1/0     (read from address 100 into r2)
	program := []uint16{
		0x4540, 0x0064, // r1 <- copy val/100
		0x4140, 0x002A, // r1 mem/0 <- copy val/42
		0x1640, 0x0000, // r2 <- copy r1 mem/0 (src=1(r1 mem), dest=6(r2 reg))
	}

	cpu.LoadProgram(program, 0)

	cpu.Step() // r1 = 100
	cpu.Step() // mem[100] = 42
	cpu.Step() // r2 = mem[100]

	if cpu.GetR2() != 42 {
		t.Errorf("R2 = %d, want 42", cpu.GetR2())
	}

	// Verify memory
	if cpu.memory.ReadLocal(100) != 42 {
		t.Errorf("mem[100] = %d, want 42", cpu.memory.ReadLocal(100))
	}
}

func TestCPU_ConditionalExecution(t *testing.T) {
	cpu := New(1024)

	// Program:
	//   r1 <- copy val/0     (set r1 to 0, sets zero flag via copy - NO! copy doesn't set flags)
	//   r1 <~ sub val/0      (test zero, set flags)
	//   r2 <0 copy val/10    (should execute: store 10 in r2)
	//   r2 <1 copy val/20    (should NOT execute: zero flag is set)
	program := []uint16{
		0x4540, 0x0000, // r1 <- copy val/0
		0x457B, 0x0000, // r1 <~ sub val/0 (effect=flags only(7), op=sub(B))
		0x4600, 0x000A, // r2 <0 copy val/10 (effect=zero(0))
		0x4610, 0x0014, // r2 <1 copy val/20 (effect=not zero(1))
	}

	cpu.LoadProgram(program, 0)

	cpu.Step() // r1 = 0
	cpu.Step() // r1 <~ sub val/0, sets zero flag
	cpu.Step() // r2 = 10 (because zero flag is set)
	cpu.Step() // r2 unchanged (because zero flag is still set)

	if cpu.GetR2() != 10 {
		t.Errorf("R2 = %d, want 10", cpu.GetR2())
	}
}

func TestCPU_Jump(t *testing.T) {
	cpu := New(1024)

	// Program:
	//   r1 <- copy val/5
	//   pc <- copy val/10    (jump to address 10)
	//   r1 <- copy val/99    (should be skipped)
	//   ... padding ...
	//   [addr 10] r2 <- copy val/42
	program := []uint16{
		0x4540, 0x0005, // [0] r1 <- copy val/5
		0x4040, 0x000A, // [2] pc <- copy val/10 (src=val(4), dest=pc(0))
		0x4540, 0x0063, // [4] r1 <- copy val/99 (skipped)
		0x0000, 0x0000, // [6] padding
		0x0000, 0x0000, // [8] padding
		0x4640, 0x002A, // [10] r2 <- copy val/42
	}

	cpu.LoadProgram(program, 0)

	cpu.Step() // r1 = 5
	cpu.Step() // pc = 10
	cpu.Step() // r2 = 42

	if cpu.GetR1() != 5 {
		t.Errorf("R1 = %d, want 5", cpu.GetR1())
	}
	if cpu.GetR2() != 42 {
		t.Errorf("R2 = %d, want 42", cpu.GetR2())
	}
}

func TestCPU_PushPop(t *testing.T) {
	cpu := New(1024)

	// Program:
	//   r1 <- copy val/1000       (stack pointer)
	//   r1/0 push <- copy val/42  (push 42 onto stack)
	//   r2 <- copy r1/0 pop       (pop into r2)
	program := []uint16{
		0x4540, 0x03E8, // r1 <- copy val/1000
		0x41C0, 0x002A, // r1 mem/0 push <- copy val/42 (M=1, dest=r1 mem, effect=always(4))
		0x16C0, 0x0000, // r2 <- copy r1 mem/0 pop (M=1, src=r1 mem(1), dest=r2 reg(6), effect=always(4))
	}

	cpu.LoadProgram(program, 0)

	cpu.Step() // r1 = 1000

	initialR1 := cpu.GetR1()
	cpu.Step() // push: r1 = 999, mem[999] = 42

	if cpu.GetR1() != initialR1-1 {
		t.Errorf("After push: R1 = %d, want %d", cpu.GetR1(), initialR1-1)
	}
	if cpu.memory.ReadLocal(999) != 42 {
		t.Errorf("mem[999] = %d, want 42", cpu.memory.ReadLocal(999))
	}

	cpu.Step() // pop: r2 = mem[999], r1 = 1000

	if cpu.GetR2() != 42 {
		t.Errorf("R2 = %d, want 42", cpu.GetR2())
	}
	if cpu.GetR1() != initialR1 {
		t.Errorf("After pop: R1 = %d, want %d", cpu.GetR1(), initialR1)
	}
}

func TestCPU_Halt(t *testing.T) {
	cpu := New(1024)

	// Program: HALT instruction
	program := []uint16{
		0x0040, 0x0000, // pc <- copy pc/0 (HALT)
		0x4540, 0x0063, // r1 <- copy val/99 (should not execute)
	}

	cpu.LoadProgram(program, 0)

	err := cpu.Step()
	if err != nil {
		t.Fatalf("Step() error = %v", err)
	}

	if !cpu.IsHalted() {
		t.Error("CPU should be halted")
	}

	// Try to step again
	err = cpu.Step()
	if err == nil {
		t.Error("Expected error when stepping halted CPU")
	}

	// r1 should still be 0
	if cpu.GetR1() != 0 {
		t.Errorf("R1 = %d, want 0 (instruction after HALT should not execute)", cpu.GetR1())
	}
}

func TestCPU_Flags(t *testing.T) {
	cpu := New(1024)

	// Test zero flag - need to use an ALU op to set flags
	program := []uint16{
		0x454B, 0x0000, // r1 <- sub val/0 (0 - 0 = 0, sets zero)
	}
	cpu.LoadProgram(program, 0)
	cpu.Step()

	if !cpu.GetFlags().Zero() {
		t.Error("Zero flag should be set")
	}

	// Test negative flag
	cpu.Reset()
	program = []uint16{
		0x454A, 0xFFFF, // r1 <- add val/-1 (0 + -1 = -1)
	}
	cpu.LoadProgram(program, 0)
	cpu.Step()

	if !cpu.GetFlags().Negative() {
		t.Error("Negative flag should be set")
	}

	// Test carry flag
	cpu.Reset()
	program = []uint16{
		0x4540, 0xFFFF, // r1 <- copy val/65535
		0x454A, 0x0001, // r1 <- add val/1 (overflow)
	}
	cpu.LoadProgram(program, 0)
	cpu.Step()
	cpu.Step()

	if !cpu.GetFlags().Carry() {
		t.Error("Carry flag should be set")
	}
}

func TestCPU_Run(t *testing.T) {
	cpu := New(1024)

	// Program that counts to 5 then halts
	//   r1 <- copy val/0
	// loop:
	//   r1 <- add val/1
	//   r1 <~ sub val/5      (test if r1 == 5)
	//   pc <1 copy val/2     (if not zero, jump to loop)
	//   halt
	program := []uint16{
		0x4540, 0x0000, // [0] r1 <- copy val/0
		0x454A, 0x0001, // [2] r1 <- add val/1
		0x457B, 0x0005, // [4] r1 <~ sub val/5
		0x4010, 0x0002, // [6] pc <1 copy val/2 (src=val(4), dest=pc(0), effect=not-zero(1))
		0x0040, 0x0000, // [8] halt
	}

	cpu.LoadProgram(program, 0)

	err := cpu.Run()
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if cpu.GetR1() != 5 {
		t.Errorf("R1 = %d, want 5", cpu.GetR1())
	}

	if !cpu.IsHalted() {
		t.Error("CPU should be halted")
	}
}

func TestCPU_Callbacks(t *testing.T) {
	cpu := New(1024)

	beforeCalled := false
	afterCalled := false

	cpu.SetBeforeExec(func(c *CPU, addr uint16) bool {
		beforeCalled = true
		return true // Allow execution
	})

	cpu.SetAfterExec(func(c *CPU, inst *Instruction) {
		afterCalled = true
	})

	program := []uint16{
		0x4540, 0x0005, // r1 <- copy val/5
	}

	cpu.LoadProgram(program, 0)
	cpu.Step()

	if !beforeCalled {
		t.Error("Before callback was not called")
	}
	if !afterCalled {
		t.Error("After callback was not called")
	}
}

func TestCPU_BreakpointCallback(t *testing.T) {
	cpu := New(1024)

	breakpointHit := false

	cpu.SetBeforeExec(func(c *CPU, addr uint16) bool {
		if addr == 2 {
			breakpointHit = true
			return false // Don't execute
		}
		return true
	})

	program := []uint16{
		0x4540, 0x0001, // [0] r1 <- copy val/1
		0x4540, 0x0002, // [2] r1 <- copy val/2 (breakpoint)
		0x4540, 0x0003, // [4] r1 <- copy val/3
	}

	cpu.LoadProgram(program, 0)

	cpu.Step() // Execute first instruction
	cpu.Step() // Hit breakpoint, don't execute

	if !breakpointHit {
		t.Error("Breakpoint was not hit")
	}

	if cpu.GetR1() != 1 {
		t.Errorf("R1 = %d, want 1 (breakpoint instruction should not execute)", cpu.GetR1())
	}

	if cpu.GetPC() != 2 {
		t.Errorf("PC = %d, want 2", cpu.GetPC())
	}
}
