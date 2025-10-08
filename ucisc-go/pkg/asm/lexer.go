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

// LexerMark represents a saved lexer state
type LexerMark struct {
	pos     int
	nextPos int
	ch      rune
	line    int
	column  int
}

// Mark saves the current lexer state
func (l *Lexer) Mark() LexerMark {
	return LexerMark{
		pos:     l.pos,
		nextPos: l.nextPos,
		ch:      l.ch,
		line:    l.line,
		column:  l.column,
	}
}

// Restore restores a previously saved lexer state
func (l *Lexer) Restore(mark LexerMark) {
	l.pos = mark.pos
	l.nextPos = mark.nextPos
	l.ch = mark.ch
	l.line = mark.line
	l.column = mark.column
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
	case '{':
		tok.Type = TokLBrace
		tok.Value = string(l.ch)
		l.readChar()
	case '}':
		tok.Type = TokRBrace
		tok.Value = string(l.ch)
		l.readChar()
	case '(':
		tok.Type = TokLParen
		tok.Value = string(l.ch)
		l.readChar()
	case ')':
		tok.Type = TokRParen
		tok.Value = string(l.ch)
		l.readChar()
	case ',':
		tok.Type = TokComma
		tok.Value = string(l.ch)
		l.readChar()
	case '=':
		tok.Type = TokEquals
		tok.Value = string(l.ch)
		l.readChar()
	case '.':
		tok.Type = TokDot
		tok.Value = string(l.ch)
		l.readChar()
	case '"':
		tok.Type = TokString
		tok.Value = l.readString()
	case '<':
		// Effect tokens: <-, <~, <0, <1, <+, <n, <&, <#, and <| variants
		tok = l.readEffect()
	case '%':
		// Could be % for data blocks or %FFFF for hex numbers
		if isHexDigit(l.peekChar()) {
			tok.Type = TokNumber
			tok.Value = l.readNumber()
		} else {
			tok.Type = TokPercent
			tok.Value = string(l.ch)
			l.readChar()
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		tok.Type = TokNumber
		tok.Value = l.readNumber()
	case '-':
		// Could be negative number or -> arrow
		if l.peekChar() == '>' {
			tok.Type = TokArrow
			tok.Value = "->"
			l.readChar()
			l.readChar()
		} else if isDigit(l.peekChar()) {
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
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
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

	// Read hex digits (could be decimal or hex depending on context)
	for isHexDigit(l.ch) {
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

// readString reads a quoted string literal
func (l *Lexer) readString() string {
	l.readChar() // consume opening quote
	start := l.pos

	for l.ch != '"' && l.ch != 0 && l.ch != '\n' {
		// Handle escape sequences
		if l.ch == '\\' {
			l.readChar() // consume backslash
			if l.ch != 0 && l.ch != '\n' {
				l.readChar() // consume escaped character
			}
		} else {
			l.readChar()
		}
	}

	value := l.input[start:l.pos]

	if l.ch == '"' {
		l.readChar() // consume closing quote
	}

	return value
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
		tok.Type = TokEffectAlways
		if pipe {
			tok.Value = "<|-"
		} else {
			tok.Value = "<-"
		}
	case '~':
		l.readChar()
		tok.Type = TokEffectFlags
		if pipe {
			tok.Value = "<|~"
		} else {
			tok.Value = "<~"
		}
	case '0':
		l.readChar()
		tok.Type = TokEffectZero
		if pipe {
			tok.Value = "<|0"
		} else {
			tok.Value = "<0"
		}
		// Check for optional ? suffix
		if l.ch == '?' {
			tok.Value += "?"
			l.readChar()
		}
	case '1':
		l.readChar()
		tok.Type = TokEffectNotZero
		if pipe {
			tok.Value = "<|1"
		} else {
			tok.Value = "<1"
		}
		// Check for optional ? suffix
		if l.ch == '?' {
			tok.Value += "?"
			l.readChar()
		}
	case '+', 'p':
		ch := l.ch
		l.readChar()
		tok.Type = TokEffectPositive
		if pipe {
			tok.Value = fmt.Sprintf("<|%c", ch)
		} else {
			tok.Value = fmt.Sprintf("<%c", ch)
		}
		// Check for optional ? suffix
		if l.ch == '?' {
			tok.Value += "?"
			l.readChar()
		}
	case 'n':
		l.readChar()
		tok.Type = TokEffectNegative
		if pipe {
			tok.Value = "<|n"
		} else {
			tok.Value = "<n"
		}
		// Check for optional ? suffix
		if l.ch == '?' {
			tok.Value += "?"
			l.readChar()
		}
	case '&', 'o':
		ch := l.ch
		l.readChar()
		tok.Type = TokEffectOverflow
		if pipe {
			tok.Value = fmt.Sprintf("<|%c", ch)
		} else {
			tok.Value = fmt.Sprintf("<%c", ch)
		}
		// Check for optional ? suffix
		if l.ch == '?' {
			tok.Value += "?"
			l.readChar()
		}
	case '#':
		l.readChar()
		tok.Type = TokEffectInterrupt
		if pipe {
			tok.Value = "<|#"
		} else {
			tok.Value = "<#"
		}
		// Check for optional ? suffix
		if l.ch == '?' {
			tok.Value += "?"
			l.readChar()
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
