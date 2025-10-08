package asm

import "github.com/grokthis/ucisc-go/pkg/cpu"

// DebugParser parses debug device function calls
// Syntax: debug.functionName(args...)
type DebugParser struct{}

func (dp *DebugParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokDebug && ctx.PeekToken().Type == TokDot
}

func (dp *DebugParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	ctx.NextToken() // consume 'debug'
	ctx.NextToken() // consume '.'

	// Expect function name
	if ctx.CurrentToken().Type != TokIdent {
		ctx.AddError("expected function name after 'debug.'")
		return nil
	}

	funcName := ctx.CurrentToken().Value
	funcTok := ctx.CurrentToken()
	ctx.NextToken() // consume function name

	// Expect '('
	if ctx.CurrentToken().Type != TokLParen {
		ctx.AddError("expected '(' after debug function name")
		return nil
	}
	ctx.NextToken() // consume '('

	// Parse arguments
	args := []*Operand{}
	for ctx.CurrentToken().Type != TokRParen && ctx.CurrentToken().Type != TokEOF {
		// Parse operand (source-style)
		source, _ := ctx.ParseSource()
		if source == nil {
			// Error already added by ParseSource
			return nil
		}
		args = append(args, source)

		// Check for comma or closing paren
		if ctx.CurrentToken().Type == TokComma {
			ctx.NextToken() // consume comma
		} else if ctx.CurrentToken().Type != TokRParen {
			ctx.AddError("expected ',' or ')' in debug function arguments")
			return nil
		}
	}

	if ctx.CurrentToken().Type != TokRParen {
		ctx.AddError("expected ')' to close debug function call")
		return nil
	}
	ctx.NextToken() // consume ')'

	return &DebugStmt{
		FuncName: funcName,
		Args:     args,
		Line:     curTok.Line,
		Col:      funcTok.Column,
	}
}

// Debug function name to command code mapping
var debugCommands = map[string]uint16{
	// Hardware-compatible assertions
	"assertEqual":       cpu.CmdAssertEq,
	"assertNotEqual":    cpu.CmdAssertNe,
	"assertGreater":     cpu.CmdAssertGt,
	"assertLess":        cpu.CmdAssertLt,
	"assertGreaterOrEqual": cpu.CmdAssertGte,
	"assertLessOrEqual": cpu.CmdAssertLte,
	"assertGreaterSigned": cpu.CmdAssertGtSigned,
	"assertLessSigned":  cpu.CmdAssertLtSigned,
	"assertZero":        cpu.CmdAssertZero,
	"assertNotZero":     cpu.CmdAssertNotZero,
	"assertHasFlags":    cpu.CmdAssertHasFlags,
	"assertNoFlags":     cpu.CmdAssertNoFlags,

	// Emulator-only commands
	"breakpoint":     cpu.CmdBreakpoint,
	"printValue":     cpu.CmdPrintValue,
	"printString":    cpu.CmdPrintString,
	"dumpRegisters":  cpu.CmdDumpRegisters,
	"dumpMemory":     cpu.CmdDumpMemory,
	"traceOn":        cpu.CmdTraceOn,
	"traceOff":       cpu.CmdTraceOff,
	"logMessage":     cpu.CmdLogMessage,

	// Control
	"reset": 0, // Special - handled by CtrlReset
}

// Expected argument counts for each command
var debugArgCounts = map[string]int{
	"assertEqual":       2,
	"assertNotEqual":    2,
	"assertGreater":     2,
	"assertLess":        2,
	"assertGreaterOrEqual": 2,
	"assertLessOrEqual": 2,
	"assertGreaterSigned": 2,
	"assertLessSigned":  2,
	"assertZero":        1,
	"assertNotZero":     1,
	"assertHasFlags":    1,
	"assertNoFlags":     1,
	"breakpoint":        0,
	"printValue":        1,
	"printString":       1,
	"dumpRegisters":     0,
	"dumpMemory":        2,
	"traceOn":           0,
	"traceOff":          0,
	"logMessage":        1,
	"reset":             0,
}
