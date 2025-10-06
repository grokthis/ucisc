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
