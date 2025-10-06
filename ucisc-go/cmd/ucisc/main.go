package main

import (
	"fmt"
	"os"

	"github.com/grokthis/ucisc-go/pkg/cpu"
	"github.com/grokthis/ucisc-go/pkg/loader"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: ucisc <command> [args]")
		fmt.Fprintln(os.Stderr, "\nCommands:")
		fmt.Fprintln(os.Stderr, "  run <file>    Run a program")
		fmt.Fprintln(os.Stderr, "  debug <file>  Debug a program")
		fmt.Fprintln(os.Stderr, "  version       Show version")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: ucisc run <file>")
			os.Exit(1)
		}
		runProgram(os.Args[2])
	case "debug":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: ucisc debug <file>")
			os.Exit(1)
		}
		debugProgram(os.Args[2])
	case "version":
		fmt.Println("ucisc-go version 0.1.0")
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func runProgram(filename string) {
	// Load program
	program, err := loader.Load(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading program: %v\n", err)
		os.Exit(1)
	}

	// Create CPU
	c := cpu.New(65536)

	// Load program into memory
	if err := c.LoadProgram(program, 0); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading program into memory: %v\n", err)
		os.Exit(1)
	}

	// Run program
	fmt.Printf("Running program: %s\n", filename)
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}

	// Print final state
	fmt.Println("\nProgram completed successfully")
	fmt.Printf("Cycles: %d\n", c.GetCycles())
	fmt.Printf("Instructions: %d\n", c.GetInstCount())
	printRegisters(c)
}

func debugProgram(filename string) {
	fmt.Fprintf(os.Stderr, "Debug mode not yet implemented\n")
	os.Exit(1)
}


func printRegisters(c *cpu.CPU) {
	fmt.Println("\nRegisters:")
	fmt.Printf("  PC:      0x%04X\n", c.GetPC())
	fmt.Printf("  R1:      0x%04X (%d)\n", c.GetR1(), c.GetR1())
	fmt.Printf("  R2:      0x%04X (%d)\n", c.GetR2(), c.GetR2())
	fmt.Printf("  R3:      0x%04X (%d)\n", c.GetR3(), c.GetR3())
	fmt.Printf("  R4:      0x%04X (%d)\n", c.GetR4(), c.GetR4())
	fmt.Printf("  R5:      0x%04X (%d)\n", c.GetR5(), c.GetR5())
	fmt.Printf("  R6:      0x%04X (%d)\n", c.GetR6(), c.GetR6())
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
