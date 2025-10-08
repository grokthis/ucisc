package asm

import (
	"strings"
)

// CallParser parses standalone function call statements
type CallParser struct{}

func (cp *CallParser) Matches(ctx ParseContext) bool {
	// Pattern: IDENT DOT IDENT LPAREN
	// Parser will restore state after this, so we can consume tokens freely

	// Check: IDENT
	if ctx.CurrentToken().Type != TokIdent {
		return false
	}
	ctx.NextToken()

	// Check: DOT
	if ctx.CurrentToken().Type != TokDot {
		return false
	}
	ctx.NextToken()

	// Check: IDENT
	if ctx.CurrentToken().Type != TokIdent {
		return false
	}
	ctx.NextToken()

	// Check: LPAREN
	return ctx.CurrentToken().Type == TokLParen
}

func (cp *CallParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	callStmt := &CallStmt{
		Line: curTok.Line,
		Col:  curTok.Column,
	}

	// Parse: stack.funcName(args) -> result
	name := curTok.Value
	ctx.NextToken() // consume identifier

	ctx.NextToken() // consume dot

	name = name + "." + ctx.CurrentToken().Value
	ctx.NextToken() // consume function name

	// Now at (
	fc := ctx.ParseFunctionCall(name)
	if fc == nil {
		return nil
	}
	callStmt.Call = fc

	// Handle -> result part
	if ctx.CurrentToken().Type == TokArrow {
		ctx.NextToken() // consume ->
		if ctx.CurrentToken().Type == TokIdent {
			resultName := ctx.CurrentToken().Value
			// Define stack.result pointing to offset 0 (where return value is)
			parts := strings.Split(name, ".")
			if len(parts) == 2 {
				defName := parts[0]
				if err := ctx.Symbols().DefineVar(defName, resultName, 0, callStmt.Line, callStmt.Col); err != nil {
					ctx.AddError(err.Error())
					return nil
				}
			}
			ctx.NextToken() // consume result name
		}
	}

	return callStmt
}
