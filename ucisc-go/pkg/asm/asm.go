// Package asm provides a uCISC assembler for compiling assembly source to machine code
package asm

import (
	"fmt"
	"os"
	"strings"
)

// Assemble compiles uCISC assembly source code into machine code
func Assemble(source string) ([]uint16, error) {
	// Parse the source code
	parser := NewParser(source)
	program, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	// Generate machine code
	codegen := NewCodeGenerator(program)
	code, err := codegen.Generate()
	if err != nil {
		return nil, err
	}

	return code, nil
}

// AssembleFile compiles a .ucisc file into machine code
func AssembleFile(filename string) ([]uint16, error) {
	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", filename, err)
	}

	// Assemble the source
	return Assemble(string(data))
}

// AssembleWithComments compiles uCISC assembly and returns code with source line comments
func AssembleWithComments(source string) ([]uint16, []string, error) {
	// Parse the source code
	parser := NewParser(source)
	program, err := parser.Parse()
	if err != nil {
		return nil, nil, err
	}

	// Generate machine code with comments
	codegen := NewCodeGenerator(program)
	code, err := codegen.Generate()
	if err != nil {
		return nil, nil, err
	}

	// Extract source lines for comments
	sourceLines := strings.Split(source, "\n")
	comments := make([]string, len(code)/2) // One comment per instruction (2 words each)

	// Map instructions to source lines
	instrIdx := 0
	for _, stmt := range program.Statements {
		if instr, ok := stmt.(*InstructionStmt); ok {
			if instrIdx < len(comments) && instr.Line > 0 && instr.Line <= len(sourceLines) {
				// Get the source line and trim whitespace
				line := strings.TrimSpace(sourceLines[instr.Line-1])
				// Remove any trailing comments
				if idx := strings.Index(line, "#"); idx >= 0 {
					line = strings.TrimSpace(line[:idx])
				}
				comments[instrIdx] = line
			}
			instrIdx++
		}
	}

	return code, comments, nil
}

// AssembleToHex compiles assembly source and formats it as hex output
// Output format is compatible with the hex loader
func AssembleToHex(source string) (string, error) {
	code, err := Assemble(source)
	if err != nil {
		return "", err
	}

	return FormatHex(code), nil
}

// FormatHex formats machine code as hex string
func FormatHex(code []uint16) string {
	var sb strings.Builder

	for i, word := range code {
		// Write hex value
		sb.WriteString(fmt.Sprintf("%04X", word))

		// Add space between words, newline every 8 words
		if i < len(code)-1 {
			if (i+1)%8 == 0 {
				sb.WriteString("\n")
			} else {
				sb.WriteString(" ")
			}
		}
	}

	return sb.String()
}

// =============================================================================
// Constants
// =============================================================================

// Instruction encoding constants
const (
	// InstructionSizeWords is the size of each uCISC instruction in 16-bit words
	// All instructions are encoded as 32 bits (2 words)
	InstructionSizeWords = 2
)

// Data encoding constants
const (
	// StringNullTerminator is the number of words added for string null terminator
	StringNullTerminator = 1
)
