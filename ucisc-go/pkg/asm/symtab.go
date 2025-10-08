package asm

import (
	"fmt"

	"github.com/grokthis/ucisc-go/pkg/cpu"
)

// SymbolTable manages register definitions and variables during parsing
type SymbolTable struct {
	scopes []*Scope // Stack of scopes
}

// Scope represents a single scope level
type Scope struct {
	definitions map[string]*Definition  // name -> register mapping
	variables   map[string]*Variable    // qualified name -> variable info
	parent      *Scope
}

// Definition represents a register alias (def statement)
type Definition struct {
	Name     string
	Register cpu.Register
	Line     int
	Col      int
}

// Variable represents a named offset on a register (var statement)
type Variable struct {
	QualifiedName string // e.g., "stack.value"
	DefName       string // e.g., "stack"
	VarName       string // e.g., "value"
	Offset        int16
	Line          int
	Col           int
}

// NewSymbolTable creates a new symbol table
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		scopes: []*Scope{newScope(nil)}, // Start with global scope
	}
}

// newScope creates a new scope
func newScope(parent *Scope) *Scope {
	return &Scope{
		definitions: make(map[string]*Definition),
		variables:   make(map[string]*Variable),
		parent:      parent,
	}
}

// PushScope creates a new nested scope
func (st *SymbolTable) PushScope() {
	st.scopes = append(st.scopes, newScope(st.currentScope()))
}

// PopScope removes the current scope
func (st *SymbolTable) PopScope() error {
	if len(st.scopes) <= 1 {
		return fmt.Errorf("cannot pop global scope")
	}
	st.scopes = st.scopes[:len(st.scopes)-1]
	return nil
}

// currentScope returns the current scope
func (st *SymbolTable) currentScope() *Scope {
	if len(st.scopes) == 0 {
		return nil
	}
	return st.scopes[len(st.scopes)-1]
}

// Define adds a new register definition to the current scope
func (st *SymbolTable) Define(name string, reg cpu.Register, line, col int) error {
	scope := st.currentScope()

	// Allow multiple definitions of the same register (as per spec)
	// But track them for reference
	def := &Definition{
		Name:     name,
		Register: reg,
		Line:     line,
		Col:      col,
	}

	scope.definitions[name] = def
	return nil
}

// DefineVar adds a new variable to the current scope
func (st *SymbolTable) DefineVar(defName, varName string, offset int16, line, col int) error {
	scope := st.currentScope()

	qualifiedName := fmt.Sprintf("%s.%s", defName, varName)

	// Check if variable already exists in current scope
	if _, exists := scope.variables[qualifiedName]; exists {
		return fmt.Errorf("variable %q already defined in current scope", qualifiedName)
	}

	v := &Variable{
		QualifiedName: qualifiedName,
		DefName:       defName,
		VarName:       varName,
		Offset:        offset,
		Line:          line,
		Col:           col,
	}

	scope.variables[qualifiedName] = v
	return nil
}

// LookupDef looks up a register definition by name
func (st *SymbolTable) LookupDef(name string) (*Definition, bool) {
	// Search from current scope up to global
	for i := len(st.scopes) - 1; i >= 0; i-- {
		if def, ok := st.scopes[i].definitions[name]; ok {
			return def, true
		}
	}
	return nil, false
}

// LookupVar looks up a variable by qualified name
func (st *SymbolTable) LookupVar(qualifiedName string) (*Variable, bool) {
	// Search from current scope up to global
	for i := len(st.scopes) - 1; i >= 0; i-- {
		if v, ok := st.scopes[i].variables[qualifiedName]; ok {
			return v, true
		}
	}
	return nil, false
}

// ResolveOperand resolves an operand that might use definitions/variables
// Returns the resolved register and offset
func (st *SymbolTable) ResolveOperand(name string) (cpu.Register, int16, bool) {
	// Check if it's a qualified name (e.g., "stack.value")
	if v, ok := st.LookupVar(name); ok {
		// Need to look up the definition to get the register
		if def, ok := st.LookupDef(v.DefName); ok {
			return def.Register, v.Offset, true
		}
	}

	// Check if it's a definition name (e.g., "stack")
	if def, ok := st.LookupDef(name); ok {
		return def.Register, 0, true
	}

	return 0, 0, false
}
