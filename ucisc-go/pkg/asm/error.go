package asm

import "fmt"

// ErrorType represents the type of parse error
type ErrorType int

const (
	ErrorUnexpectedToken ErrorType = iota
	ErrorUndefinedIdentifier
	ErrorExpectedToken
	ErrorInvalidNumber
	ErrorInvalidSyntax
	ErrorDuplicateDefinition
	ErrorSymbolTable
)

// ParseError represents a structured parse error with location information
type ParseError struct {
	Type    ErrorType
	Line    int
	Column  int
	Message string
	Token   string // The problematic token, if applicable
}

// Error implements the error interface
func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d:%d: %s", e.Line, e.Column, e.Message)
}

// NewParseError creates a new parse error
func NewParseError(errType ErrorType, line, col int, message string) *ParseError {
	return &ParseError{
		Type:    errType,
		Line:    line,
		Column:  col,
		Message: message,
	}
}

// NewParseErrorWithToken creates a new parse error with token information
func NewParseErrorWithToken(errType ErrorType, line, col int, message, token string) *ParseError {
	return &ParseError{
		Type:    errType,
		Line:    line,
		Column:  col,
		Message: message,
		Token:   token,
	}
}

// ParseErrors is a collection of parse errors
type ParseErrors struct {
	Errors []*ParseError
}

// Error implements the error interface for ParseErrors
func (pe *ParseErrors) Error() string {
	if len(pe.Errors) == 0 {
		return "no errors"
	}

	result := "parse errors:\n"
	for _, err := range pe.Errors {
		result += err.Error() + "\n"
	}
	return result[:len(result)-1] // Remove trailing newline
}

// Add adds a parse error
func (pe *ParseErrors) Add(err *ParseError) {
	pe.Errors = append(pe.Errors, err)
}

// AddString adds a string error (for backward compatibility during migration)
func (pe *ParseErrors) AddString(msg string) {
	pe.Errors = append(pe.Errors, &ParseError{
		Type:    ErrorInvalidSyntax,
		Message: msg,
	})
}

// HasErrors returns true if there are any errors
func (pe *ParseErrors) HasErrors() bool {
	return len(pe.Errors) > 0
}
