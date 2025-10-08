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
	// For simplicity during transition, cast to *Parser
	// TODO: Migrate to use only ParseContext methods
	p := ctx.(*Parser)
	callStmt := &CallStmt{
		Line: p.curTok.Line,
		Col:  p.curTok.Column,
	}
	
	// Parse: stack.funcName(args) -> result
	name := p.curTok.Value
	p.nextToken() // consume identifier
	
	p.nextToken() // consume dot
	
	name = name + "." + p.curTok.Value
	p.nextToken() // consume function name
	
	// Now at (
	fc := p.parseFunctionCall(name)
	if fc == nil {
		return nil
	}
	callStmt.Call = fc
	
	// Handle -> result part
	if p.curTok.Type == TokArrow {
		p.nextToken() // consume ->
		if p.curTok.Type == TokIdent {
			resultName := p.curTok.Value
			// Define stack.result pointing to offset 0 (where return value is)
			parts := strings.Split(name, ".")
			if len(parts) == 2 {
				defName := parts[0]
				if err := p.symtab.DefineVar(defName, resultName, 0, callStmt.Line, callStmt.Col); err != nil {
					p.addError(err.Error())
					return nil
				}
			}
			p.nextToken() // consume result name
		}
	}
	
	return callStmt
}
