package asm

import (
	"testing"

	"github.com/grokthis/ucisc-go/pkg/cpu"
)

func TestAssemble_SimpleFunctionDefinition(t *testing.T) {
	source := `
def stack/r1 <- copy val/1000

fun stack.add5(value) -> result {
	&r2 <- copy val/5
	&r2 <- add val/3
	pc <- copy pc/0
}

pc <- copy pc/0  # halt
`

	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	if len(code) == 0 {
		t.Fatal("No code generated")
	}
}

func TestAssemble_FunctionWithMultipleArgs(t *testing.T) {
	source := `
def stack/r1 <- copy val/1000

fun stack.multiply(a, b) -> result {
	&r2 <- copy stack.a
	&r3 <- copy stack.b
	&r2 <- mult &r3
	stack.result <- copy &r2
	pc <- copy stack.return pop
}

pc <- copy pc/0
`

	_, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}
}

func TestAssemble_NestedScopes(t *testing.T) {
	source := `
def stack/r1 <- copy val/1000

{
	def local/r2 <- copy val/100
	var local.temp/0 <- copy val/42

	{
		# Inner scope can access outer scope
		&r3 <- copy local.temp
	}
}

pc <- copy pc/0
`

	_, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}
}

func TestAssemble_FunctionCallConvention(t *testing.T) {
	// Test a more complete function with proper calling convention
	source := `
def stack/r1 <- copy val/1000

# Jump to main program
pc <- copy pc/program_start

# Simple function that increments a value
fun stack.increment(value) -> result {
	stack.result <- copy stack.value
	stack.result <- add val/1
	pc <- copy stack.return pop
}

# Main program
program_start:
	# Reserve space for return value
	stack/0 push <- copy val/0  # result placeholder

	# Push return address
	stack/0 push <- copy pc/after_call

	# Push argument
	stack/0 push <- copy val/42

	# Jump to function
	pc <- copy pc/stack.increment

after_call:
	# Get result (at stack/0 now)
	&r2 <- copy stack/0 pop

	# Should be 43
	pc <- copy pc/0
`

	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Run it to verify it works
	c := cpu.New(65536)
	c.LoadProgram(code, 0)
	err = c.Run()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Result should be in r2
	if c.GetR2() != 43 {
		t.Errorf("Expected r2=43, got %d", c.GetR2())
	}
}
