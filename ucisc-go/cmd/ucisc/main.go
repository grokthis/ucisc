package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/grokthis/ucisc-go/pkg/asm"
	"github.com/grokthis/ucisc-go/pkg/cpu"
	"github.com/grokthis/ucisc-go/pkg/disasm"
	"github.com/grokthis/ucisc-go/pkg/loader"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: ucisc <command> [args]")
		fmt.Fprintln(os.Stderr, "\nCommands:")
		fmt.Fprintln(os.Stderr, "  run <file>       Run a program (.hex or .ucisc)")
		fmt.Fprintln(os.Stderr, "  compile <file>   Compile assembly to hex")
		fmt.Fprintln(os.Stderr, "  version          Show version")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: ucisc run <file> [-trace [max_instructions]]")
			os.Exit(1)
		}
		trace := false
		maxInstructions := 100
		if len(os.Args) >= 4 && os.Args[3] == "-trace" {
			trace = true
			if len(os.Args) >= 5 {
				if n, err := strconv.Atoi(os.Args[4]); err == nil && n > 0 {
					maxInstructions = n
				}
			}
		}
		runProgram(os.Args[2], trace, maxInstructions)
	case "compile":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: ucisc compile <file.ucisc> [-o output.hex]")
			os.Exit(1)
		}
		outputFile := ""
		if len(os.Args) >= 5 && os.Args[3] == "-o" {
			outputFile = os.Args[4]
		}
		compileProgram(os.Args[2], outputFile)
	case "version":
		fmt.Println("ucisc-go version 0.2.0")
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// loadProgramFile loads a program from a file (either .ucisc or .hex)
func loadProgramFile(filename string) ([]uint16, error) {
	if strings.HasSuffix(filename, ".ucisc") {
		return compileToMemory(filename)
	}
	return loader.Load(filename)
}

// execProgram executes a program and returns the final CPU state
func execProgram(program []uint16, trace bool, maxInstructions int) (*cpu.CPU, error) {
	// Create CPU
	c := cpu.New(65536)

	// Check program size before loading
	if len(program) > 65536 {
		return nil, fmt.Errorf("error loading program into memory: program too large (%d words, max 65536)", len(program))
	}

	// Load program into memory
	if err := c.LoadProgram(program, 0); err != nil {
		return nil, fmt.Errorf("error loading program into memory: %w", err)
	}

	// Set up trace callback if requested
	if trace {
		instrCount := 0
		c.SetBeforeExec(func(c *cpu.CPU, addr uint16) bool {
			instrCount++
			if instrCount > maxInstructions {
				fmt.Printf("\n[Trace limit reached: %d instructions]\n", maxInstructions)
				return false // Stop execution
			}

			// Read and decode the instruction at current PC
			word1 := c.Memory().ReadLocal(addr)
			word2 := c.Memory().ReadLocal(addr + 1)
			instrRaw := (uint32(word1) << 16) | uint32(word2)

			// Decode and disassemble
			inst, err := cpu.Decode(instrRaw)
			asmText := "???"
			if err == nil {
				asmText = disasm.FormatInstruction(inst)
			}

			fmt.Printf("[%04X %04X] %-30s  PC=%04X R1=%04X R2=%04X R3=%04X R5=%04X R6=%04X R7=%04X F=%04X\n",
				word1, word2, asmText, addr, c.GetR1(), c.GetR2(), c.GetR3(), c.GetR5(), c.GetR6(), c.GetR7(), uint16(c.GetFlags()))
			return true // Continue execution
		})
	}

	// Run program
	if err := c.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %w", err)
	}

	return c, nil
}

func runProgram(filename string, trace bool, maxInstructions int) {
	program, err := loadProgramFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading program: %v\n", err)
		os.Exit(1)
	}

	if !trace {
		fmt.Printf("Running program: %s\n", filename)
	}

	c, err := execProgram(program, trace, maxInstructions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Print final state
	fmt.Println("\nProgram completed successfully")
	fmt.Printf("Cycles: %d\n", c.GetCycles())
	fmt.Printf("Instructions: %d\n", c.GetInstCount())
	printRegisters(c)
}

func compileToMemory(filename string) ([]uint16, error) {
	// Assemble the file
	code, err := asm.AssembleFile(filename)
	if err != nil {
		return nil, fmt.Errorf("assembly failed: %w", err)
	}

	return code, nil
}

func compileProgram(inputFile, outputFile string) {
	// Read source file
	source, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading source: %v\n", err)
		os.Exit(1)
	}

	// Compile with comments
	code, comments, err := asm.AssembleWithComments(string(source))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Determine output filename
	if outputFile == "" {
		outputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".hex"
	}

	// Write hex file
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Write hex words (2 per instruction) with comments
	for i := 0; i < len(code); i += 2 {
		if i+1 < len(code) {
			fmt.Fprintf(f, "%04X %04X", code[i], code[i+1])
		} else {
			fmt.Fprintf(f, "%04X", code[i])
		}

		// Add source comment if available
		instrIdx := i / 2
		if instrIdx < len(comments) && comments[instrIdx] != "" {
			fmt.Fprintf(f, "  # %s", comments[instrIdx])
		}
		fmt.Fprintf(f, "\n")
	}

	fmt.Printf("Compiled %s -> %s (%d words, %d instructions)\n", inputFile, outputFile, len(code), len(code)/2)
}

func printRegisters(c *cpu.CPU) {
	fmt.Println("\nRegisters:")
	fmt.Printf("  PC:      0x%04X\n", c.GetPC())
	fmt.Printf("  R1:      0x%04X (%d)\n", c.GetR1(), c.GetR1())
	fmt.Printf("  R2:      0x%04X (%d)\n", c.GetR2(), c.GetR2())
	fmt.Printf("  R3:      0x%04X (%d)\n", c.GetR3(), c.GetR3())
	fmt.Printf("  R5:      0x%04X (%d)\n", c.GetR5(), c.GetR5())
	fmt.Printf("  R6:      0x%04X (%d)\n", c.GetR6(), c.GetR6())
	fmt.Printf("  R7:      0x%04X (%d)\n", c.GetR7(), c.GetR7())
	fmt.Printf("  Banking: 0x%04X\n", c.GetBanking())

	flags := c.GetFlags()
	fmt.Printf("  Flags:   0x%04X [", uint16(flags))
	if flags.Zero() {
		fmt.Print("Z")
	}
	if flags.Negative() {
		fmt.Print("N")
	}
	if flags.Carry() {
		fmt.Print("C")
	}
	if flags.Overflow() {
		fmt.Print("O")
	}
	if flags.Interrupt() {
		fmt.Print("I")
	}
	if flags.Signed() {
		fmt.Print("S")
	}
	if flags.Halted() {
		fmt.Print("H")
	}
	fmt.Println("]")
}
