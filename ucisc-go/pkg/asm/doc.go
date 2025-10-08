// Package asm implements an assembler for the uCISC architecture.
//
// The assembler processes uCISC assembly language source code and produces
// executable machine code in the form of 16-bit word arrays. It supports
// a rich assembly language with labels, constants, variables, functions,
// and macro-like syntax sugar.
//
// # Assembly Language Features
//
//   - Instructions: Full uCISC instruction set with source/dest operands
//   - Labels: Named locations in code for jumps and references
//   - Constants: Named immediate values (.def directive)
//   - Variables: Named memory locations (.var directive)
//   - Functions: Scoped code blocks with local labels (.fun directive)
//   - Syntax sugar: call/ret for function calls, simplified addressing
//   - String literals: Automatic conversion to word data
//
// # Pipeline
//
// The assembler uses a multi-stage pipeline:
//   1. Lexer: Tokenizes source text into tokens
//   2. Parser: Builds an Abstract Syntax Tree (AST)
//   3. Code Generator: Two-pass code generation (addresses, then code)
//
// # Example Usage
//
//	source := `
//	    .def MAX 100
//	    r1 <- copy val/MAX
//	    r2 <- add &r1/val/1
//	    pc <- copy pc  # halt
//	`
//
//	code, err := asm.Assemble(source)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// code is []uint16 ready to load into CPU
//
// # File Assembly
//
//	code, err := asm.AssembleFile("program.ucisc")
//	if err != nil {
//	    log.Fatal(err)
//	}
package asm
