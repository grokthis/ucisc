package asm

// ElementParser is an interface for parsing different statement types
// Uses ParseContext to follow Interface Segregation Principle
type ElementParser interface {
	// Matches returns true if this parser can handle the current token sequence
	// It should not consume tokens or add errors
	Matches(ctx ParseContext) bool

	// Parse parses the statement and returns it, or nil on error
	// It should consume tokens and may add errors
	Parse(ctx ParseContext) Statement
}
