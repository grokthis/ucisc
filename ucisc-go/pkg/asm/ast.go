package asm

import (
	"github.com/grokthis/ucisc-go/pkg/cpu"
)

// Node is the base interface for all AST nodes
type Node interface {
	node()
}

// Statement represents a statement in the assembly program
type Statement interface {
	Node
	statement()
}

// Program is the root of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) node() {}

// LabelStmt represents a label declaration
type LabelStmt struct {
	Name string
	Line int
	Col  int
}

func (l *LabelStmt) node()      {}
func (l *LabelStmt) statement() {}

// InstructionStmt represents an assembly instruction
type InstructionStmt struct {
	Destination *Operand
	Effect      cpu.Effect
	Operation   cpu.Operation
	Source      *Operand
	Push        bool
	Pop         bool
	Line        int
	Col         int
}

func (i *InstructionStmt) node()      {}
func (i *InstructionStmt) statement() {}

// Operand represents an operand (source or destination)
type Operand struct {
	Type         OperandType
	Register     cpu.Register
	Offset       int16
	IsAddress    bool   // true for &r1 (register destination)
	LabelRef     string // for pc/label references
	FunctionCall *FunctionCall // for function call syntax
	Line         int
	Col          int
}

// FunctionCall represents a function call expression like stack.func(arg1, arg2)
type FunctionCall struct {
	DefName   string     // e.g., "stack"
	FuncName  string     // e.g., "double"
	Args      []*Operand // Arguments
	Line      int
	Col       int
}

// OperandType specifies the type of operand
type OperandType int

const (
	OpImmediate OperandType = iota // val/42
	OpRegister                     // r1, &r1
	OpMemory                       // r1/offset
	OpPCRelative                   // pc/offset, pc/label
)

// DefStmt represents a register definition: def name/register
// Example: def stack/r1 <- copy val/0
type DefStmt struct {
	Name         string            // The alias name (e.g., "stack")
	Register     cpu.Register      // The register being aliased
	Initializer  *InstructionStmt  // Optional initialization instruction
	Line         int
	Col          int
}

func (d *DefStmt) node()      {}
func (d *DefStmt) statement() {}

// VarStmt represents a variable declaration with offset
// Example: var stack.value/0 <- copy val/42
type VarStmt struct {
	DefName      string            // The definition name (e.g., "stack")
	VarName      string            // The variable name (e.g., "value")
	Offset       int16             // The offset
	Push         bool              // Whether to push
	Initializer  *InstructionStmt  // Optional initialization instruction
	Line         int
	Col          int
}

func (v *VarStmt) node()      {}
func (v *VarStmt) statement() {}

// FunStmt represents a function definition
// Example: fun stack.readChar(serialDevice) -> char { ... }
type FunStmt struct {
	DefName      string       // The stack definition (e.g., "stack")
	Name         string       // Function name (e.g., "readChar")
	Args         []string     // Argument names
	Returns      []string     // Return value names
	Body         []Statement  // Function body
	Line         int
	Col          int
}

func (f *FunStmt) node()      {}
func (f *FunStmt) statement() {}

// DataStmt represents a data block
// Example: % 001F 0A48 656C 6C6F
type DataStmt struct {
	Values []uint16 // The data values
	Line   int
	Col    int
}

func (d *DataStmt) node()      {}
func (d *DataStmt) statement() {}

// StringStmt represents a string literal that gets converted to data
// Example: "Hello, world!"
type StringStmt struct {
	Value string // The string content
	Line  int
	Col   int
}

func (s *StringStmt) node()      {}
func (s *StringStmt) statement() {}

// ScopeStmt represents a scope block with curly braces
// Example: { ... }
type ScopeStmt struct {
	Statements []Statement
	Line       int
	Col        int
}

func (s *ScopeStmt) node()      {}
func (s *ScopeStmt) statement() {}

// CallStmt represents a standalone function call statement
// Example: stack.double(val/2) -> result
type CallStmt struct {
	Call *FunctionCall
	Line int
	Col  int
}

func (c *CallStmt) node()      {}
func (c *CallStmt) statement() {}

// DebugStmt represents a debug device function call
// Example: debug.assertEqual(&r3, val/42)
type DebugStmt struct {
	FuncName string     // e.g., "assertEqual", "printValue", "reset"
	Args     []*Operand // Arguments
	Line     int
	Col      int
}

func (d *DebugStmt) node()      {}
func (d *DebugStmt) statement() {}
