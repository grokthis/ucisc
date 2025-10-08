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

// elementParsers is the ordered list of element parsers
// More specific parsers should come first
var elementParsers []ElementParser

// RegisterElementParser registers an element parser
func RegisterElementParser(ep ElementParser) {
	elementParsers = append(elementParsers, ep)
}

// InitElementParsers initializes the element parsers in order
func InitElementParsers() {
	elementParsers = nil
	
	// Order matters - more specific first
	RegisterElementParser(&LabelParser{})
	RegisterElementParser(&DefParser{})
	RegisterElementParser(&VarParser{})
	RegisterElementParser(&FunParser{})
	RegisterElementParser(&CallParser{})
	RegisterElementParser(&DataParser{})
	RegisterElementParser(&StringParser{})
	RegisterElementParser(&ScopeParser{})
	RegisterElementParser(&InstructionParser{}) // Most general, comes last
}
