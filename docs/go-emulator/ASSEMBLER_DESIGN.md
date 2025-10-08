# uCISC Go Assembler - Design Document

## 1. Executive Summary

This document outlines the design for a professional-quality Go implementation of a uCISC assembler. The assembler will compile uCISC assembly language into binary hex format suitable for execution on the uCISC emulator.

## 2. Design Philosophy

The uCISC assembler will be implemented in phases, starting with core functionality and progressively adding advanced features. This approach ensures:
- Early delivery of a working assembler
- Incremental complexity management
- Thorough testing at each phase
- Clear milestone markers

## 3. Assembly Language Overview

uCISC supports a human-friendly assembly syntax with the following key features:

### 3.1 Basic Syntax
```ucisc
# Comments start with #
destination <- operation source
```

### 3.2 Registers
- General purpose: r1, r2, r3, r4, r5, r6
- Special: pc, val (immediate), control/banking, interrupt, flags

### 3.3 Addressing Modes
- **Register destination** (address): `&r1` - stores to register
- **Memory destination** (value): `r1`, `r1/5` - stores to memory at register+offset
- **Immediate source**: `val/42` - constant value
- **Register source**: `r2/10` - value from register+offset
- **Memory source**: `r2/10` - value from memory at register+offset
- **PC-relative**: `pc/10` - PC+offset

### 3.4 Effects (Conditional Execution)
- `<-`, `<|` - Store always
- `<~`, `<|~` - Set flags only (don't store)
- `<0`, `<|0` - Store if zero
- `<1`, `<|1` - Store if not zero
- `<+`, `<p`, `<|+`, `<|p` - Store if positive
- `<n`, `<|-`, `<|n` - Store if negative
- `<&`, `<o`, `<|&`, `<|o` - Store if overflow
- `<#`, `<|#` - Store if interrupt

### 3.5 Operations
- **Logical**: copy, and, or, xor, inv
- **Shift**: shl, shr, swap
- **Bit manipulation**: msb, lsb
- **Arithmetic**: add, addc, sub, mult, multsw

### 3.6 Push/Pop
- `push` modifier on destination: decrements register before store
- `pop` modifier on source: increments register after load

### 3.7 Advanced Features
- **Labels**: `loop:` - mark code location
- **Definitions**: `def stack/r1 <- copy val/0` - name registers
- **Variables**: `var stack.value/0 <- copy val/100` - offset tracking
- **Functions**: `fun stack.readChar(arg) -> result { ... }` - calling convention
- **Scopes**: `{ ... }` - namespace isolation
- **Strings**: `"Hello"` - inline string data
- **Data**: `% 0x1234 0x5678` - raw hex data

## 4. Assembler Architecture

### 4.1 Pipeline Overview
```
Source Code
    ↓
[Lexer] → Tokens
    ↓
[Parser] → AST (Abstract Syntax Tree)
    ↓
[Symbol Resolution] → Resolved AST
    ↓
[Code Generator] → Machine Code
    ↓
Hex Output
```

### 4.2 Component Design

#### 4.2.1 Lexer (Tokenizer)
**Purpose**: Convert source text into stream of tokens

**Input**: Raw source code string
**Output**: `[]Token`

**Token Types**:
```go
type TokenType int

const (
    // Literals
    TokIdentifier    // label, def_name, var_name
    TokNumber        // 42, 0x1234, %FFFF
    TokString        // "hello"

    // Keywords
    TokDef           // def
    TokVar           // var
    TokFun           // fun
    TokPush          // push
    TokPop           // pop
    TokVal           // val

    // Registers
    TokR1, TokR2, TokR3, TokR4, TokR5, TokR6
    TokPC, TokControl, TokBanking, TokInterrupt, TokFlags

    // Effects
    TokEffectAlways  // <-, <|
    TokEffectFlags   // <~, <|~
    TokEffectZero    // <0, <|0
    TokEffectNotZero // <1, <|1
    // ... etc

    // Operations
    TokCopy, TokAnd, TokOr, TokXor, TokInv
    TokShl, TokShr, TokSwap
    TokMsb, TokLsb
    TokAdd, TokAddc, TokSub, TokMult, TokMultsw

    // Punctuation
    TokSlash         // /
    TokAmpersand     // &
    TokColon         // :
    TokLBrace        // {
    TokRBrace        // }
    TokLParen        // (
    TokRParen        // )
    TokComma         // ,
    TokArrow         // ->
    TokPercent       // %

    // Special
    TokNewline
    TokEOF
    TokComment
)

type Token struct {
    Type   TokenType
    Value  string
    Line   int
    Column int
}
```

#### 4.2.2 Parser
**Purpose**: Convert token stream into Abstract Syntax Tree

**Input**: `[]Token`
**Output**: `*Program` (AST root)

**AST Node Types**:
```go
type Program struct {
    Statements []Statement
}

type Statement interface {
    statement()
}

// Instruction statement (the core assembly instruction)
type InstructionStmt struct {
    Destination Operand
    Effect      Effect
    Operation   Operation
    Source      Operand
    Push        bool
    Pop         bool
}

// Label statement
type LabelStmt struct {
    Name string
}

// Definition statement
type DefStmt struct {
    Name     string
    Register Register
    Init     *InstructionStmt // optional initialization
}

// Variable statement
type VarStmt struct {
    Name     string
    Base     string  // base definition name
    Offset   int
    Init     *InstructionStmt // optional initialization
}

// Function statement
type FunStmt struct {
    Name       string
    Stack      string  // stack register name
    Args       []string
    Returns    []string
    Body       *BlockStmt
}

// Block statement
type BlockStmt struct {
    Statements []Statement
}

// Data statement
type DataStmt struct {
    Values []uint16
}

// String statement
type StringStmt struct {
    Value string
}

type Operand struct {
    Type     OperandType  // Register, Memory, Immediate, PCRelative
    Register Register
    Offset   int16
    IsAddr   bool  // true for &r1 (register destination)
}
```

#### 4.2.3 Symbol Table
**Purpose**: Track definitions, variables, labels, and scopes

```go
type SymbolTable struct {
    scopes []*Scope
}

type Scope struct {
    parent   *Scope
    symbols  map[string]*Symbol
    labels   map[string]int  // label -> address
}

type Symbol struct {
    Type     SymbolType  // Definition, Variable, Function, Label
    Name     string
    Register Register    // for defs and vars
    Offset   int16       // for vars
    Address  int         // for labels
}
```

#### 4.2.4 Code Generator
**Purpose**: Convert resolved AST into machine code

**Input**: Resolved AST + Symbol Table
**Output**: `[]uint16` (machine code)

**Process**:
1. Traverse AST in order
2. For each instruction:
   - Resolve operands (lookup symbols, calculate offsets)
   - Encode instruction fields
   - Generate 32-bit instruction (2 words)
3. For data/strings:
   - Convert to word values
4. Output hex format

## 5. Implementation Phases

### Phase 1: Core Assembler (MVP) - ✅ COMPLETE
**Goal**: Assemble basic instructions without advanced features

**Features**:
- ✅ Lexer for basic tokens
- ✅ Parser for simple instructions
- ✅ Register and immediate operands
- ✅ All ALU operations
- ✅ All effect codes
- ✅ Push/pop modifiers (with correct immediate+1 behavior)
- ✅ Labels (forward and backward references)
- ✅ PC-relative addressing with correct offset calculation
- ✅ Comments
- ✅ Hex output

**Example Programs**:
```ucisc
# Fibonacci (simple version)
    &r1 <- copy val/0     # fib(0) = 0
    &r2 <- copy val/1     # fib(1) = 1
    &r3 <- copy val/6     # counter

loop:
    r1 <- add r2          # r1 = r1 + r2
    &r4 <- copy r1        # temp = r1
    &r1 <- copy r2        # r1 = r2
    &r2 <- copy r4        # r2 = temp
    r3 <- sub val/1       # r3--
    pc <1 copy pc/loop    # if r3 != 0, goto loop

    pc <- copy pc/0       # halt
```

### Phase 2: Definitions and Variables - ✅ COMPLETE
**Goal**: Add symbolic names for registers and offsets

**Features**:
- ✅ def statements (def name/register <- init)
- ✅ var statements (var name.field/offset <- init)
- ✅ Symbol table with scope management
- ✅ Variable offset tracking
- ✅ Qualified name resolution (stack.value)

**Example Programs**:
```ucisc
def stack/r1 <- copy val/1000
var stack.temp/0
var stack.counter/1 <- copy val/10

loop:
    stack.counter <- sub val/1
    pc <1 copy pc/loop

    pc <- copy pc/0  # halt
```

### Phase 3: Functions - ✅ COMPLETE
**Goal**: Add function calling convention and scopes

**Features**:
- ✅ fun statements (fun stack.name(args) -> returns { body })
- ✅ Automatic stack frame management
- ✅ Argument/return handling with correct offsets
- ✅ Function calls with stack setup (stack.func(args) -> result)
- ✅ Scope blocks ({ ... })
- ✅ Qualified function labels (defName.funcName)

### Phase 4: Data and Strings - ✅ COMPLETE
**Goal**: Add inline data support

**Features**:
- ✅ String literals ("text")
- ✅ Hex data blocks (% FFFF 1234)
- ✅ PC-relative label references for data

### Phase 5: Advanced Features - ✅ COMPLETE
**Goal**: Complete feature parity and extensibility

**Features**:
- ✅ Element-based parser architecture
- ✅ Mark/restore lexer for non-destructive lookahead
- ✅ Error recovery (skip to next line on parse error)
- ✅ Detailed error messages with line/column information
- ✅ Two-pass code generation for forward references
- ✅ Standalone function call statements

## 6. Package Structure

```
pkg/
├── asm/
│   ├── asm.go         # Public API, Assemble() function
│   ├── lexer.go       # Tokenization
│   ├── lexer_test.go
│   ├── parser.go      # AST generation
│   ├── parser_test.go
│   ├── ast.go         # AST node definitions
│   ├── symbol.go      # Symbol table
│   ├── codegen.go     # Code generation
│   ├── codegen_test.go
│   └── errors.go      # Error types and formatting
```

## 7. Error Handling

The assembler will provide high-quality error messages:

```
error: undefined symbol 'loop'
  --> fibonacci.ucisc:15:12
   |
15 |     pc <1 copy pc/loop
   |                   ^^^^ not found in current scope
   |
   = help: did you mean 'Loop' (defined at line 8)?
```

Error categories:
- Syntax errors (unexpected token, missing operand, etc.)
- Semantic errors (undefined symbol, type mismatch, etc.)
- Encoding errors (immediate value out of range, etc.)

## 8. Testing Strategy

### 8.1 Unit Tests
- Lexer: Test tokenization of all constructs
- Parser: Test AST generation for all statement types
- Symbol table: Test scope resolution, symbol lookup
- Code generator: Test instruction encoding

### 8.2 Integration Tests
- End-to-end assembly of example programs
- Verify assembled output matches expected hex
- Cross-validate with emulator execution

### 8.3 Golden Tests
- Assemble existing .ucisc examples
- Compare output with ucisc-kotlin compiler (if compatible)

## 9. API Design

### 9.1 Basic API
```go
package asm

// Assemble compiles uCISC assembly source to machine code
func Assemble(source string) ([]uint16, error)

// AssembleFile compiles a .ucisc file to machine code
func AssembleFile(filename string) ([]uint16, error)

// AssembleToHex compiles and formats as hex output
func AssembleToHex(source string) (string, error)
```

### 9.2 Advanced API
```go
// Options for assembler behavior
type Options struct {
    DebugInfo    bool   // Include debug symbols
    Optimize     bool   // Enable optimizations
    StartAddress uint16 // Load address
}

// AssembleWithOptions provides fine-grained control
func AssembleWithOptions(source string, opts Options) (*Result, error)

type Result struct {
    Code       []uint16
    Symbols    *SymbolTable
    DebugInfo  *DebugInfo
    Errors     []Error
    Warnings   []Warning
}
```

## 10. Success Criteria

### Phase 1 (Core Assembler) - ✅ COMPLETE
- ✅ Can assemble basic instructions with all operations
- ✅ Supports labels and forward references
- ✅ Generates correct machine code validated by emulator
- ✅ Has comprehensive test coverage (>80%)
- ✅ Produces helpful error messages
- ✅ CLI integration (`ucisc compile program.ucisc`)

### Phase 2 (Definitions/Variables) - ✅ COMPLETE
- ✅ Symbol table with scope management
- ✅ def and var statements working
- ✅ Qualified name resolution (def.var)
- ✅ Optional initializers on def/var

### Phase 3 (Functions) - ✅ COMPLETE
- ✅ Function declarations with args/returns
- ✅ Automatic stack frame layout
- ✅ Function calls with automatic setup
- ✅ Scope blocks working correctly

### Phase 4 (Data/Strings) - ✅ COMPLETE
- ✅ String literals compile to words
- ✅ Hex data blocks inline in code
- ✅ PC-relative references work correctly

### Phase 5 (Advanced) - ✅ COMPLETE
- ✅ Element-based parser architecture
- ✅ Extensible for future syntax additions
- ✅ Mark/restore lookahead working
- ✅ All tests passing

## 11. Implementation Notes

### Critical Fixes Applied
1. **PC-Relative Addressing**: Changed from SourcePC to SourceReg with RegPC, calculate relative offsets
2. **Pop Behavior**: CPU pops by immediate+1 ("you always get 1 for free")
3. **Register Destinations**: Both `pc` and `&pc` treated identically (16-bit immediate range)
4. **Halt Detection**: Recognizes both `pc <- copy pc` and `pc <- copy pc/0`
5. **Address Tracking**: pass2 maintains currentAddr through recursion for correct PC-relative encoding
6. **Stack Frame Layout**: args at offset 0, return address at args+1, results at higher offsets

### Architecture Decisions
1. **Element Parsers**: Each statement type has independent Matches()/Parse() methods
2. **Mark/Restore**: Lexer supports save/restore for non-destructive lookahead
3. **Two-Pass Generation**: Pass1 calculates addresses, Pass2 generates code
4. **Symbol Table**: Scoped with definitions and variables in separate namespaces

---

**Document Version**: 1.1
**Last Updated**: 2025-10-07
**Status**: All Phases Complete - Production Ready
