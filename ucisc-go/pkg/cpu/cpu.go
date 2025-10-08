package cpu

import "fmt"

// BeforeExecFunc is called before each instruction execution
// Return false to skip execution (for breakpoints)
type BeforeExecFunc func(cpu *CPU, addr uint16) bool

// AfterExecFunc is called after each instruction execution
type AfterExecFunc func(cpu *CPU, inst *Instruction)

// CPU represents the uCISC processor
type CPU struct {
	// Registers (matches ISA: r1, r2, r3, r5, r6, r7)
	r1      uint16
	r2      uint16
	r3      uint16
	r5      uint16
	r6      uint16
	r7      uint16
	pc      uint16
	banking uint16
	flags   Flags
	intReg  uint16

	// Components
	memory *Memory
	alu    *ALU

	// Execution state
	halted bool
	cycles uint64
	instCount uint64

	// Debug/safety options
	strictBounds bool // If true, error on out-of-bounds instead of wrapping

	// Callbacks for debugging
	beforeExec BeforeExecFunc
	afterExec  AfterExecFunc
}

// New creates a new CPU with the specified memory size
func New(memSize int) *CPU {
	return &CPU{
		memory:       NewMemory(memSize),
		alu:          NewALU(),
		halted:       false,
		cycles:       0,
		instCount:    0,
		strictBounds: true, // Default to strict bounds checking for development
	}
}

// SetStrictBounds enables/disables strict bounds checking
// When false, PC and memory accesses wrap around (hardware-accurate)
// When true, out-of-bounds accesses cause errors (safer for development)
func (c *CPU) SetStrictBounds(strict bool) {
	c.strictBounds = strict
}

// Reset resets the CPU to initial state
func (c *CPU) Reset() {
	c.r1 = 0
	c.r2 = 0
	c.r3 = 0
	c.r5 = 0
	c.r6 = 0
	c.r7 = 0
	c.pc = 0
	c.banking = 0
	c.flags = 0
	c.intReg = 0
	c.halted = false
	c.cycles = 0
	c.instCount = 0
}

// Step executes a single instruction
func (c *CPU) Step() error {
	if c.halted {
		return fmt.Errorf("CPU is halted")
	}

	// Call before-exec callback (for breakpoints)
	if c.beforeExec != nil {
		if !c.beforeExec(c, c.pc) {
			return nil // Breakpoint hit, don't execute
		}
	}

	// Check PC bounds if strict mode enabled
	memSize := c.memory.Size()
	if memSize == 0 {
		return fmt.Errorf("cannot execute: memory size is 0")
	}

	if c.strictBounds {
		// In strict mode, check that we can fetch a complete 2-word instruction
		if int(c.pc) >= memSize || int(c.pc)+1 >= memSize {
			return fmt.Errorf("PC 0x%04X out of bounds (memory size: %d words)", c.pc, memSize)
		}
	}

	// Fetch instruction (2 words)
	var word1, word2 uint16
	if c.strictBounds || int(c.pc)+1 < memSize {
		// Direct access when in bounds or strict mode already checked
		word1 = c.memory.ReadLocal(c.pc)
		word2 = c.memory.ReadLocal(c.pc + 1)
	} else {
		// Wrap around in non-strict mode (only when memSize fits in uint16)
		if memSize <= 65536 {
			word1 = c.memory.ReadLocal(c.pc % uint16(memSize))
			word2 = c.memory.ReadLocal((c.pc + 1) % uint16(memSize))
		} else {
			// Should never happen since NewMemory caps at 65536
			word1 = c.memory.ReadLocal(c.pc)
			word2 = c.memory.ReadLocal(c.pc + 1)
		}
	}
	raw := (uint32(word1) << 16) | uint32(word2)

	// Decode instruction
	inst, err := Decode(raw)
	if err != nil {
		return fmt.Errorf("decode error at 0x%04X: %w", c.pc, err)
	}

	// Check for HALT instruction
	if inst.IsHalt() {
		c.halted = true
		c.pc += InstructionSizeWords
		return nil
	}

	// Execute instruction
	if err := c.execute(inst); err != nil {
		return fmt.Errorf("execution error at 0x%04X: %w", c.pc, err)
	}

	// Update cycles (each instruction takes at least 1 cycle)
	c.cycles++
	c.instCount++

	// Call after-exec callback
	if c.afterExec != nil {
		c.afterExec(c, inst)
	}

	return nil
}

// Run executes instructions until halt or error
func (c *CPU) Run() error {
	for !c.halted {
		if err := c.Step(); err != nil {
			return err
		}
	}
	return nil
}

// execute executes a decoded instruction
func (c *CPU) execute(inst *Instruction) error {
	// Update ALU signed mode from flags
	c.alu.SetSigned(c.flags.Signed())

	// Read source value
	srcValue := c.readSource(inst.Source, inst.Immediate)

	// Read destination value (for 2-arg operations)
	destValue := c.readDestination(inst.Destination, inst.Offset)

	// Handle increment (pop) on source if needed
	if inst.Increment && inst.Source.Type == SourceMem && inst.Destination.Type != DestMem {
		// Pop: increment source register after read
		// Pop amount is immediate + bonus increment (you always get 1 for free)
		popAmount := uint16(inst.Immediate) + PopBonusIncrement
		c.setRegister(inst.Source.Reg, c.getRegister(inst.Source.Reg)+popAmount)
	}

	// Perform ALU operation
	result := c.alu.Execute(inst.Operation, destValue, srcValue, c.flags.Carry())

	// Check effect condition and store if satisfied
	// Note: conditions check PREVIOUS flags (from last ALU op), not current result
	stored := false
	if c.shouldStore(inst.Effect) {
		c.writeDestination(inst.Destination, inst.Offset, result.Value, inst.Increment)
		stored = true
	}

	// Update flags (only for non-copy operations and when effect allows it)
	if inst.Operation != OpCopy && inst.Effect != EffectFlagsOnly {
		c.updateFlags(result)
	} else if inst.Effect == EffectFlagsOnly {
		// Flags-only effect: update flags but don't store
		c.updateFlags(result)
	}

	// Update PC (unless PC was the destination register AND we stored to it)
	pcIsDestination := inst.Destination.Type == DestReg && inst.Destination.Reg == RegPC
	if !pcIsDestination || !stored {
		// Only increment if: PC is NOT the destination OR we didn't store
		c.pc += InstructionSizeWords
	}

	return nil
}

// readSource reads the source value
func (c *CPU) readSource(src Source, immediate int16) uint16 {
	switch src.Type {
	case SourceVal:
		// Immediate value
		return uint16(immediate)
	case SourceReg:
		// Register value + immediate
		regValue := c.getRegister(src.Reg)
		return uint16(int32(regValue) + int32(immediate))
	case SourceMem:
		// Memory at register + immediate
		regValue := c.getRegister(src.Reg)
		addr := uint16(int32(regValue) + int32(immediate))
		return c.memory.Read(addr, src.Reg, c.banking)
	case SourcePC:
		// PC value + immediate
		return uint16(int32(c.pc) + int32(immediate))
	default:
		return 0
	}
}

// readDestination reads the destination value
func (c *CPU) readDestination(dest Destination, offset uint8) uint16 {
	if dest.Type == DestReg {
		return c.getRegister(dest.Reg)
	}

	// DestMem
	regValue := c.getRegister(dest.Reg)
	addr := regValue + uint16(offset)
	return c.memory.Read(addr, dest.Reg, c.banking)
}

// writeDestination writes the result to the destination
func (c *CPU) writeDestination(dest Destination, offset uint8, value uint16, increment bool) {
	if dest.Type == DestReg {
		c.setRegister(dest.Reg, value)
		return
	}

	// DestMem
	regValue := c.getRegister(dest.Reg)

	if increment {
		// Push: decrement register before write
		regValue = regValue - PushDecrement
		c.setRegister(dest.Reg, regValue)
	}

	addr := regValue + uint16(offset)
	c.memory.Write(addr, dest.Reg, c.banking, value)
}

// shouldStore checks if the effect condition is satisfied based on CPU flags
func (c *CPU) shouldStore(effect Effect) bool {
	switch effect {
	case EffectAlways:
		return true
	case EffectZero:
		return c.flags.Zero()
	case EffectNotZero:
		return !c.flags.Zero()
	case EffectNegative:
		return c.flags.Negative()
	case EffectPositive:
		return !c.flags.Negative()
	case EffectOverflow:
		return c.flags.Overflow()
	case EffectInterrupt:
		return c.flags.Interrupt()
	case EffectFlagsOnly:
		return false // Don't store, just set flags
	default:
		return false
	}
}

// updateFlags updates the processor flags based on ALU result
func (c *CPU) updateFlags(result Result) {
	if result.Zero {
		c.flags = c.flags.Set(FlagZero)
	} else {
		c.flags = c.flags.Clear(FlagZero)
	}

	if result.Negative {
		c.flags = c.flags.Set(FlagNegative)
	} else {
		c.flags = c.flags.Clear(FlagNegative)
	}

	if result.Carry {
		c.flags = c.flags.Set(FlagCarry)
	} else {
		c.flags = c.flags.Clear(FlagCarry)
	}

	if result.Overflow {
		c.flags = c.flags.Set(FlagOverflow)
	} else {
		c.flags = c.flags.Clear(FlagOverflow)
	}
}

// getRegister gets a register value
func (c *CPU) getRegister(reg Register) uint16 {
	switch reg {
	case RegPC:
		return c.pc
	case RegR1:
		return c.r1
	case RegR2:
		return c.r2
	case RegR3:
		return c.r3
	case RegBanking:
		return c.banking
	case RegR5:
		return c.r5
	case RegR6:
		return c.r6
	case RegR7:
		return c.r7
	case RegFlags:
		return uint16(c.flags)
	case RegInt:
		return c.intReg
	default:
		return 0
	}
}

// setRegister sets a register value
func (c *CPU) setRegister(reg Register, value uint16) {
	switch reg {
	case RegPC:
		c.pc = value
	case RegR1:
		c.r1 = value
	case RegR2:
		c.r2 = value
	case RegR3:
		c.r3 = value
	case RegBanking:
		c.banking = value
	case RegR5:
		c.r5 = value
	case RegR6:
		c.r6 = value
	case RegR7:
		c.r7 = value
	case RegFlags:
		c.flags = Flags(value)
		// Update ALU signed mode when flags change
		c.alu.SetSigned(c.flags.Signed())
		// Check halt flag
		if c.flags.Halted() {
			c.halted = true
		}
	case RegInt:
		c.intReg = value
	}
}

// Public accessor methods for registers
func (c *CPU) GetPC() uint16        { return c.pc }
func (c *CPU) GetR1() uint16        { return c.r1 }
func (c *CPU) GetR2() uint16        { return c.r2 }
func (c *CPU) GetR3() uint16        { return c.r3 }
func (c *CPU) GetR5() uint16        { return c.r5 }
func (c *CPU) GetR6() uint16        { return c.r6 }
func (c *CPU) GetR7() uint16        { return c.r7 }
func (c *CPU) GetBanking() uint16   { return c.banking }
func (c *CPU) GetFlags() Flags      { return c.flags }
func (c *CPU) GetInterrupt() uint16 { return c.intReg }
func (c *CPU) GetCycles() uint64    { return c.cycles }
func (c *CPU) GetInstCount() uint64 { return c.instCount }
func (c *CPU) IsHalted() bool       { return c.halted }

// SetPC sets the program counter
func (c *CPU) SetPC(value uint16) { c.pc = value }

// LoadProgram loads a program into memory
func (c *CPU) LoadProgram(data []uint16, offset uint16) error {
	return c.memory.LoadProgram(data, offset)
}

// SetBeforeExec sets the before-execution callback
func (c *CPU) SetBeforeExec(fn BeforeExecFunc) {
	c.beforeExec = fn
}

// SetAfterExec sets the after-execution callback
func (c *CPU) SetAfterExec(fn AfterExecFunc) {
	c.afterExec = fn
}

// Memory returns the memory system (for debugger access)
func (c *CPU) Memory() *Memory {
	return c.memory
}
