package asm

import (
	"strings"
	"testing"
)

// TestStructuredErrors verifies structured error types are created correctly
func TestStructuredErrors(t *testing.T) {
	source := `
def stack INVALID
`
	p := NewParser(source)
	_, err := p.Parse()

	if err == nil {
		t.Fatal("Expected parse error for missing /")
	}

	// Check that it's a ParseErrors type
	parseErrs, ok := err.(*ParseErrors)
	if !ok {
		t.Fatalf("Expected *ParseErrors, got %T", err)
	}

	if !parseErrs.HasErrors() {
		t.Fatal("Expected errors to be present")
	}

	if len(parseErrs.Errors) == 0 {
		t.Fatal("Expected at least one error")
	}

	// Verify error has location information
	firstErr := parseErrs.Errors[0]
	if firstErr.Line == 0 {
		t.Error("Expected error to have line number")
	}
	if firstErr.Column == 0 {
		t.Error("Expected error to have column number")
	}
}

// TestErrorFormatting verifies error messages are formatted correctly
func TestErrorFormatting(t *testing.T) {
	source := `
def stack MISSING_SLASH
`
	_, err := Assemble(source)

	if err == nil {
		t.Fatal("Expected parse error")
	}

	errMsg := err.Error()

	// Should contain line number
	if !strings.Contains(errMsg, "line 2:") {
		t.Errorf("Error message should contain line number: %s", errMsg)
	}
}

// TestMultipleErrors verifies multiple errors are collected
func TestMultipleErrors(t *testing.T) {
	source := `
def stack MISSING1
def other MISSING2
`
	p := NewParser(source)
	_, err := p.Parse()

	if err == nil {
		t.Fatal("Expected parse errors")
	}

	parseErrs, ok := err.(*ParseErrors)
	if !ok {
		t.Fatalf("Expected *ParseErrors, got %T", err)
	}

	// Should have collected both errors
	if len(parseErrs.Errors) < 2 {
		t.Errorf("Expected at least 2 errors, got %d", len(parseErrs.Errors))
	}
}

// TestLegacyErrorsMethod verifies backward compatibility
func TestLegacyErrorsMethod(t *testing.T) {
	source := `
def stack INVALID
`
	p := NewParser(source)
	p.Parse()

	// Legacy Errors() method should still work
	legacyErrs := p.Errors()
	if len(legacyErrs) == 0 {
		t.Error("Expected legacy Errors() method to return errors")
	}

	// Should be strings with line:col format
	if len(legacyErrs) > 0 && !strings.Contains(legacyErrs[0], "line") {
		t.Errorf("Expected legacy error format with 'line', got: %s", legacyErrs[0])
	}
}

// TestGetErrors verifies GetErrors returns structured errors
func TestGetErrors(t *testing.T) {
	source := `
def stack INVALID
`
	p := NewParser(source)
	p.Parse()

	// GetErrors() should return *ParseErrors
	errs := p.GetErrors()
	if errs == nil {
		t.Fatal("Expected GetErrors() to return non-nil")
	}

	if !errs.HasErrors() {
		t.Error("Expected GetErrors() to have errors")
	}
}

// TestErrorTypes verifies different error types can be distinguished
func TestErrorTypes(t *testing.T) {
	// Create errors of different types
	err1 := NewParseError(ErrorUnexpectedToken, 1, 1, "unexpected token")
	err2 := NewParseError(ErrorUndefinedIdentifier, 2, 5, "undefined identifier")
	err3 := NewParseError(ErrorInvalidNumber, 3, 10, "invalid number")

	if err1.Type != ErrorUnexpectedToken {
		t.Error("Expected ErrorUnexpectedToken")
	}
	if err2.Type != ErrorUndefinedIdentifier {
		t.Error("Expected ErrorUndefinedIdentifier")
	}
	if err3.Type != ErrorInvalidNumber {
		t.Error("Expected ErrorInvalidNumber")
	}

	// Verify error messages contain location
	if !strings.Contains(err1.Error(), "line 1:1") {
		t.Errorf("Error should contain location: %s", err1.Error())
	}
}
