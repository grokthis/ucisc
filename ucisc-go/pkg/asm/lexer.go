package asm

import (
	"fmt"
	"unicode"
)

// Lexer tokenizes uCISC assembly source code
type Lexer struct {
	input   string
	pos     int  // current position in input
	nextPos int  // next reading position
	ch      rune // current character
	line    int  // current line number (1-based)
	column  int  // current column number (1-based)
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar() // Initialize first character
	return l
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.nextPos >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = rune(l.input[l.nextPos])
	}
	l.pos = l.nextPos
	l.nextPos++
	l.column++
}

// peekChar looks at the next character without advancing
func (l *Lexer) peekChar() rune {
	if l.nextPos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.nextPos])
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	tok := Token{
		Line:   l.line,
		Column: l.column,
	}

	switch l.ch {
	case 0:
		tok.Type = TokEOF
	case '\n':
		tok.Type = TokNewline
		l.line++
		l.column = 0
		l.readChar()
	case '#':
		tok.Type = TokComment
		tok.Value = l.readComment()
	case '/':
		tok.Type = TokSlash
		tok.Value = string(l.ch)
		l.readChar()
	case '&':
		tok.Type = TokAmpersand
		tok.Value = string(l.ch)
		l.readChar()
	case ':':
		tok.Type = TokColon
		tok.Value = string(l.ch)
		l.readChar()
	case '?':
		tok.Type = TokQuestion
		tok.Value = string(l.ch)
		l.readChar()
	case '<':
		// Effect tokens: <-, <~, <0, <1, <+, <n, <&, <#, and <| variants
		tok = l.readEffect()
	case '%', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		tok.Type = TokNumber
		tok.Value = l.readNumber()
	case '-':
		// Negative number
		if isDigit(l.peekChar()) {
			tok.Type = TokNumber
			tok.Value = l.readNumber()
		} else {
			tok.Type = TokIllegal
			tok.Value = string(l.ch)
			l.readChar()
		}
	default:
		if isLetter(l.ch) {
			tok.Value = l.readIdentifier()
			tok.Type = LookupIdent(tok.Value)
			return tok
		}
		tok.Type = TokIllegal
		tok.Value = string(l.ch)
		l.readChar()
	}

	return tok
}

// skipWhitespace skips spaces and tabs (but not newlines)
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier (label, variable, etc.)
func (l *Lexer) readIdentifier() string {
	start := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '.' || l.ch == '-' {
		l.readChar()
	}
	return l.input[start:l.pos]
}

// readNumber reads a number (decimal, hex, or % notation)
func (l *Lexer) readNumber() string {
	start := l.pos

	// Handle negative sign
	if l.ch == '-' {
		l.readChar()
	}

	// Handle % prefix (hex in uCISC format)
	if l.ch == '%' {
		l.readChar()
		// Read hex digits
		for isHexDigit(l.ch) {
			l.readChar()
		}
		return l.input[start:l.pos]
	}

	// Handle 0x prefix
	if l.ch == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		l.readChar() // consume '0'
		l.readChar() // consume 'x'
		for isHexDigit(l.ch) {
			l.readChar()
		}
		return l.input[start:l.pos]
	}

	// Read decimal number
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[start:l.pos]
}

// readComment reads a comment from # to end of line
func (l *Lexer) readComment() string {
	start := l.pos
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return l.input[start:l.pos]
}

// readEffect reads an effect token (<-, <~, <0, etc.)
func (l *Lexer) readEffect() Token {
	tok := Token{
		Line:   l.line,
		Column: l.column,
	}

	if l.ch != '<' {
		tok.Type = TokIllegal
		tok.Value = string(l.ch)
		l.readChar()
		return tok
	}

	l.readChar() // consume '<'

	// Check for <| variant
	pipe := false
	if l.ch == '|' {
		pipe = true
		l.readChar()
	}

	// Now check the effect character
	switch l.ch {
	case '-':
		l.readChar()
		if pipe {
			tok.Type = TokEffectAlways
			tok.Value = "<|-"
		} else {
			tok.Type = TokEffectAlways
			tok.Value = "<-"
		}
	case '~':
		l.readChar()
		if pipe {
			tok.Type = TokEffectFlags
			tok.Value = "<|~"
		} else {
			tok.Type = TokEffectFlags
			tok.Value = "<~"
		}
	case '0':
		l.readChar()
		if pipe {
			tok.Type = TokEffectZero
			tok.Value = "<|0"
		} else {
			tok.Type = TokEffectZero
			tok.Value = "<0"
		}
	case '1':
		l.readChar()
		if pipe {
			tok.Type = TokEffectNotZero
			tok.Value = "<|1"
		} else {
			tok.Type = TokEffectNotZero
			tok.Value = "<1"
		}
	case '+', 'p':
		ch := l.ch
		l.readChar()
		if pipe {
			tok.Type = TokEffectPositive
			tok.Value = fmt.Sprintf("<|%c", ch)
		} else {
			tok.Type = TokEffectPositive
			tok.Value = fmt.Sprintf("<%c", ch)
		}
	case 'n':
		l.readChar()
		if pipe {
			tok.Type = TokEffectNegative
			tok.Value = "<|n"
		} else {
			tok.Type = TokEffectNegative
			tok.Value = "<n"
		}
	case '&', 'o':
		ch := l.ch
		l.readChar()
		if pipe {
			tok.Type = TokEffectOverflow
			tok.Value = fmt.Sprintf("<|%c", ch)
		} else {
			tok.Type = TokEffectOverflow
			tok.Value = fmt.Sprintf("<%c", ch)
		}
	case '#':
		l.readChar()
		if pipe {
			tok.Type = TokEffectInterrupt
			tok.Value = "<|#"
		} else {
			tok.Type = TokEffectInterrupt
			tok.Value = "<#"
		}
	default:
		// Just '<' or '<|' without valid effect
		tok.Type = TokIllegal
		if pipe {
			tok.Value = "<|"
		} else {
			tok.Value = "<"
		}
	}

	return tok
}

// Helper functions

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isHexDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}
