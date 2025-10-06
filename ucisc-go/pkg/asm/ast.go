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
	Type       OperandType
	Register   cpu.Register
	Offset     int16
	IsAddress  bool   // true for &r1 (register destination)
	LabelRef   string // for pc/label references
	Line       int
	Col        int
}

// OperandType specifies the type of operand
type OperandType int

const (
	OpImmediate OperandType = iota // val/42
	OpRegister                     // r1, &r1
	OpMemory                       // r1/offset
	OpPCRelative                   // pc/offset, pc/label
)
