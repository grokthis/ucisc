package asm

import (
	"strings"
)

// CallParser parses standalone function call statements
type CallParser struct{}

func (cp *CallParser) Matches(p *Parser) bool {
	// Pattern: IDENT DOT IDENT LPAREN
	// Parser will restore state after this, so we can consume tokens freely

	// Check: IDENT
	if p.curTok.Type != TokIdent {
		return false
	}
	p.nextToken()

	// Check: DOT
	if p.curTok.Type != TokDot {
		return false
	}
	p.nextToken()

	// Check: IDENT
	if p.curTok.Type != TokIdent {
		return false
	}
	p.nextToken()

	// Check: LPAREN
	return p.curTok.Type == TokLParen
}

func (cp *CallParser) Parse(p *Parser) Statement {
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
