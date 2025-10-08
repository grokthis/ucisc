// Package cpu implements a complete emulator for the uCISC (Micro-CISC) architecture.
//
// The uCISC is a custom CISC architecture with a 16-bit data path and 32-bit instructions.
// Each instruction consists of two 16-bit words encoding the source, destination, operation,
// effect, and optional immediate value or memory offset.
//
// # Architecture Overview
//
// The CPU provides:
//   - 6 general-purpose 16-bit registers (R1-R3, R5-R7)
//   - Special registers: PC, Flags, Banking, Interrupt
//   - 64K words of addressable memory
//   - 16 operations (copy, add, sub, and, or, xor, shl, shr, addc, etc.)
//   - Flag-based conditional execution
//   - Memory banking for future device I/O support
//
// # Key Types
//
//   - CPU: Main processor state and execution engine
//   - Memory: 64K word memory system with banking support
//   - Instruction: Decoded instruction representation
//   - ALU: Arithmetic and logic operations
//
// # Example Usage
//
//	// Create a new CPU with 64K memory
//	c := cpu.New(65536)
//
//	// Load a program
//	program := []uint16{0x4540, 0x0064, 0x0040, 0x0000} // r1 <- 100; halt
//	c.LoadProgram(program, 0)
//
//	// Run until halt
//	if err := c.Run(); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Check results
//	fmt.Printf("R1 = %d\n", c.GetR1())
//	fmt.Printf("Cycles: %d\n", c.GetCycles())
package cpu
