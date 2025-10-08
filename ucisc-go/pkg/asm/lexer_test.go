package asm

import (
	"testing"
)

func TestLexer_BasicTokens(t *testing.T) {
	input := "/ & : ?"
	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{TokSlash, "/"},
		{TokAmpersand, "&"},
		{TokColon, ":"},
		{TokQuestion, "?"},
		{TokEOF, ""},
	}

	l := NewLexer(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, tt.expectedType, tok.Type)
		}
		if tok.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - wrong value. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}

func TestLexer_Registers(t *testing.T) {
	input := "r1 r2 r3 r5 r6 r7 pc control banking flags interrupt"
	tests := []TokenType{
		TokR1, TokR2, TokR3, TokR5, TokR6, TokR7,
		TokPC, TokControl, TokBanking, TokFlags, TokInterrupt,
		TokEOF,
	}

	l := NewLexer(input)
	for i, expectedType := range tests {
		tok := l.NextToken()
		if tok.Type != expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, expectedType, tok.Type)
		}
	}
}

func TestLexer_Operations(t *testing.T) {
	input := "copy and or xor inv shl shr swap msb lsb add addc sub mult multsw"
	tests := []TokenType{
		TokCopy, TokAnd, TokOr, TokXor, TokInv,
		TokShl, TokShr, TokSwap, TokMsb, TokLsb,
		TokAdd, TokAddc, TokSub, TokMult, TokMultsw,
		TokEOF,
	}

	l := NewLexer(input)
	for i, expectedType := range tests {
		tok := l.NextToken()
		if tok.Type != expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, expectedType, tok.Type)
		}
	}
}

func TestLexer_Effects(t *testing.T) {
	tests := []struct {
		input        string
		expectedType TokenType
	}{
		{"<-", TokEffectAlways},
		{"<|", TokIllegal}, // Just <| without effect char is illegal
		{"<|-", TokEffectAlways},
		{"<~", TokEffectFlags},
		{"<|~", TokEffectFlags},
		{"<0", TokEffectZero},
		{"<|0", TokEffectZero},
		{"<1", TokEffectNotZero},
		{"<|1", TokEffectNotZero},
		{"<+", TokEffectPositive},
		{"<p", TokEffectPositive},
		{"<|+", TokEffectPositive},
		{"<|p", TokEffectPositive},
		{"<n", TokEffectNegative},
		{"<|n", TokEffectNegative},
		{"<&", TokEffectOverflow},
		{"<o", TokEffectOverflow},
		{"<|&", TokEffectOverflow},
		{"<|o", TokEffectOverflow},
		{"<#", TokEffectInterrupt},
		{"<|#", TokEffectInterrupt},
	}

	for _, tt := range tests {
		l := NewLexer(tt.input)
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("input=%q - wrong token type. expected=%v, got=%v",
				tt.input, tt.expectedType, tok.Type)
		}
	}
}

func TestLexer_Numbers(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue string
	}{
		{"0", "0"},
		{"42", "42"},
		{"1234", "1234"},
		{"-5", "-5"},
		{"-100", "-100"},
		{"0x1234", "0x1234"},
		{"0xFFFF", "0xFFFF"},
		{"0xabcd", "0xabcd"},
		{"%1234", "%1234"},
		{"%FFFF", "%FFFF"},
	}

	for _, tt := range tests {
		l := NewLexer(tt.input)
		tok := l.NextToken()
		if tok.Type != TokNumber {
			t.Errorf("input=%q - wrong token type. expected=TokNumber, got=%v",
				tt.input, tok.Type)
		}
		if tok.Value != tt.expectedValue {
			t.Errorf("input=%q - wrong value. expected=%q, got=%q",
				tt.input, tt.expectedValue, tok.Value)
		}
	}
}

func TestLexer_Identifiers(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue string
	}{
		{"loop", "loop"},
		{"my_label", "my_label"},
		{"stack", "stack"},      // Dots are separate tokens now
		{"fib_n", "fib_n"},      // Hyphens are separate tokens now
		{"Label123", "Label123"},
	}

	for _, tt := range tests {
		l := NewLexer(tt.input)
		tok := l.NextToken()
		if tok.Type != TokIdent {
			t.Errorf("input=%q - wrong token type. expected=TokIdent, got=%v",
				tt.input, tok.Type)
		}
		if tok.Value != tt.expectedValue {
			t.Errorf("input=%q - wrong value. expected=%q, got=%q",
				tt.input, tt.expectedValue, tok.Value)
		}
	}
}

func TestLexer_Comments(t *testing.T) {
	input := `# This is a comment
r1 # inline comment
# another comment`

	l := NewLexer(input)

	tests := []TokenType{
		TokComment,
		TokNewline,
		TokR1,
		TokComment,
		TokNewline,
		TokComment,
		TokEOF,
	}

	for i, expectedType := range tests {
		tok := l.NextToken()
		if tok.Type != expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, expectedType, tok.Type)
		}
	}
}

func TestLexer_SimpleInstruction(t *testing.T) {
	input := "&r1 <- copy val/42"

	l := NewLexer(input)

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{TokAmpersand, "&"},
		{TokR1, "r1"},
		{TokEffectAlways, "<-"},
		{TokCopy, "copy"},
		{TokVal, "val"},
		{TokSlash, "/"},
		{TokNumber, "42"},
		{TokEOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, tt.expectedType, tok.Type)
		}
		if tt.expectedValue != "" && tok.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - wrong value. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}

func TestLexer_InstructionWithLabel(t *testing.T) {
	input := `loop:
    r1 <- add r2/10
    pc <1 copy pc/loop`

	l := NewLexer(input)

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{TokIdent, "loop"},
		{TokColon, ":"},
		{TokNewline, ""},
		{TokR1, "r1"},
		{TokEffectAlways, "<-"},
		{TokAdd, "add"},
		{TokR2, "r2"},
		{TokSlash, "/"},
		{TokNumber, "10"},
		{TokNewline, ""},
		{TokPC, "pc"},
		{TokEffectNotZero, "<1"},
		{TokCopy, "copy"},
		{TokPC, "pc"},
		{TokSlash, "/"},
		{TokIdent, "loop"},
		{TokEOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v (value=%q)",
				i, tt.expectedType, tok.Type, tok.Value)
		}
		if tt.expectedValue != "" && tok.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - wrong value. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}

func TestLexer_PushPop(t *testing.T) {
	input := "r1/0 push <- copy val/100\npc <- copy r1/0 pop"

	l := NewLexer(input)

	tests := []TokenType{
		TokR1, TokSlash, TokNumber,
		TokPush,
		TokEffectAlways,
		TokCopy,
		TokVal, TokSlash, TokNumber,
		TokNewline,
		TokPC,
		TokEffectAlways,
		TokCopy,
		TokR1, TokSlash, TokNumber,
		TokPop,
		TokEOF,
	}

	for i, expectedType := range tests {
		tok := l.NextToken()
		if tok.Type != expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%v, got=%v",
				i, expectedType, tok.Type)
		}
	}
}

func TestLexer_LineColumn(t *testing.T) {
	input := `r1
r2`

	l := NewLexer(input)

	tok := l.NextToken() // r1
	if tok.Line != 1 || tok.Column != 1 {
		t.Errorf("r1 position wrong. expected line=1, col=1, got line=%d, col=%d",
			tok.Line, tok.Column)
	}

	tok = l.NextToken() // newline
	if tok.Line != 1 {
		t.Errorf("newline line wrong. expected=1, got=%d", tok.Line)
	}

	tok = l.NextToken() // r2
	if tok.Line != 2 || tok.Column != 1 {
		t.Errorf("r2 position wrong. expected line=2, col=1, got line=%d, col=%d",
			tok.Line, tok.Column)
	}
}
