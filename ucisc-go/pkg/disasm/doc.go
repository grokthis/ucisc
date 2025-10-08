// Package disasm provides a disassembler for the uCISC architecture.
//
// The disassembler converts uCISC machine code back into human-readable
// assembly language syntax. This is useful for debugging, program analysis,
// and understanding compiled code.
//
// # Features
//
//   - Decodes 32-bit instructions into assembly syntax
//   - Produces output compatible with the assembler
//   - Handles all instruction formats (register, memory, immediate)
//   - Formats operations, effects, and operands clearly
//
// # Example Usage
//
//	// Disassemble a single instruction
//	inst, err := cpu.Decode(0x45400064) // r1 <- copy val/100
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(disasm.FormatInstruction(inst))
//	// Output: &r1 <- copy val/100
//
//	// Disassemble a program
//	program := []uint16{0x4540, 0x0064, 0x0040, 0x0000}
//	for i := 0; i < len(program); i += 2 {
//	    word1 := program[i]
//	    word2 := program[i+1]
//	    raw := (uint32(word1) << 16) | uint32(word2)
//	    inst, _ := cpu.Decode(raw)
//	    fmt.Printf("%04X: %s\n", i, disasm.FormatInstruction(inst))
//	}
package disasm
