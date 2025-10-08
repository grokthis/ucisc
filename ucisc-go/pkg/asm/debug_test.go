package asm

import (
	"testing"

	"github.com/grokthis/ucisc-go/pkg/cpu"
)

func TestDebugSyntax(t *testing.T) {
	tests := []struct {
		name string
		code string
		want int // Expected number of words
	}{
		{
			name: "assertEqual",
			code: `
				r3 <- copy val/42
				debug.assertEqual(&r3, val/42)
			`,
			want: 10, // r3 <- copy (2) + cmd (2) + arg1 (2) + arg2 (2) + exec (2) = 10
		},
		{
			name: "assertZero",
			code: `
				r1 <- copy val/0
				debug.assertZero(&r1)
			`,
			want: 8, // r1 <- copy (2) + cmd (2) + arg1 (2) + exec (2) = 8
		},
		{
			name: "reset",
			code: `debug.reset()`,
			want: 2, // Just write -2 to &debug
		},
		{
			name: "printValue",
			code: `
				r2 <- copy val/100
				debug.printValue(&r2)
			`,
			want: 8,
		},
		{
			name: "dumpRegisters",
			code: `debug.dumpRegisters()`,
			want: 4, // cmd (2) + exec (2)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := Assemble(tt.code)
			if err != nil {
				t.Fatalf("Assembly failed: %v", err)
			}

			if len(code) != tt.want {
				t.Errorf("Expected %d words, got %d", tt.want, len(code))
			}
		})
	}
}

func TestDebugExecution(t *testing.T) {
	source := `
		r3 <- copy val/42
		debug.assertEqual(&r3, val/42)
		pc <- copy &pc  # halt
	`

	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assembly failed: %v", err)
	}

	// Run the code
	c := cpu.New(1024)
	c.LoadProgram(code, 0)
	if err := c.Run(); err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// Check debug device - should have no failures
	debug := c.GetDebugDevice()
	if debug.HasFailures() {
		t.Errorf("Unexpected failures: %v", debug.GetFailures())
	}
}

func TestDebugExecutionFailure(t *testing.T) {
	source := `
		r3 <- copy val/42
		debug.assertEqual(&r3, val/43)  # This should fail
		pc <- copy &pc  # halt
	`

	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assembly failed: %v", err)
	}

	// Run the code
	c := cpu.New(1024)
	c.LoadProgram(code, 0)
	if err := c.Run(); err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// Check debug device - should have 1 failure
	debug := c.GetDebugDevice()
	if !debug.HasFailures() {
		t.Error("Expected assertion failure")
	}
	if debug.Read(0x4) != 1 { // Failure count
		t.Errorf("Expected 1 failure, got %d", debug.Read(0x4))
	}
}

func TestDebugInvalidFunction(t *testing.T) {
	source := `debug.invalidFunction()`

	_, err := Assemble(source)
	if err == nil {
		t.Fatal("Expected error for invalid debug function")
	}
}

func TestDebugWrongArgCount(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{"assertEqual_too_few", `debug.assertEqual(&r1)`},
		{"assertEqual_too_many", `debug.assertEqual(&r1, val/1, val/2)`},
		{"assertZero_too_many", `debug.assertZero(&r1, val/1)`},
		{"reset_with_args", `debug.reset(val/1)`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Assemble(tt.code)
			if err == nil {
				t.Error("Expected error for wrong argument count")
			}
		})
	}
}
