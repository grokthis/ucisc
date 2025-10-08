package asm

import (
	"testing"
)

func TestAssemble_DataBlocks(t *testing.T) {
	source := `
	pc <- copy pc/data  # skip over data
data:
% 0042 00FF 1234 ABCD
program_start:
	&r1 <- copy val/0x42
	pc <- copy pc/0
`

	code, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}

	// Should have: jump (2) + data (4) + copy (2) + halt (2) = 10 words
	expectedLen := 10
	if len(code) != expectedLen {
		t.Errorf("Expected %d words, got %d", expectedLen, len(code))
		t.Logf("Code: %#v", code)
	}

	// Verify data is present (at words 2-5)
	if len(code) >= 6 {
		if code[2] != 0x0042 || code[3] != 0x00FF || code[4] != 0x1234 || code[5] != 0xABCD {
			t.Errorf("Data not properly embedded: got %04X %04X %04X %04X",
				code[2], code[3], code[4], code[5])
		}
	}
}

func TestAssemble_StringLiterals(t *testing.T) {
	source := `
	pc <- copy pc/hello  # skip string
hello: "Hello, world!"
program_start:
	&r1 <- copy val/0
	pc <- copy pc/0
`

	_, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}
}

func TestAssemble_DefSimple(t *testing.T) {
	// Simplest def usage - just an alias
	source := `
def stack/r1
&stack <- copy val/1000
&r2 <- copy stack/0
pc <- copy pc/0
`

	_, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}
}

func TestAssemble_DefWithInitializer(t *testing.T) {
	source := `
def stack/r1 <- copy val/1000
&r2 <- copy stack/0
pc <- copy pc/0
`

	_, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}
}

func TestAssemble_VarSimple(t *testing.T) {
	source := `
def object/r2 <- copy val/2000

var object.x/0 <- copy val/100
var object.y/1 <- copy val/200
var object.z/2 <- copy val/300

&r1 <- copy object.x
&r1 <- add object.y

pc <- copy pc/0
`

	_, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}
}

func TestAssemble_AlternativeEffectSyntax(t *testing.T) {
	source := `
&r1 <- copy val/0
&r1 <~ sub val/0        # Set flags
&r2 <0? copy val/1      # With ? suffix
&r3 <1? copy val/2
&r5 <n? copy val/3
&r6 <|+ copy val/4      # With | prefix
pc <- copy pc/0
`

	_, err := Assemble(source)
	if err != nil {
		t.Fatalf("Assemble failed: %v", err)
	}
}
