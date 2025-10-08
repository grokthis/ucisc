package asm

// StringParser parses string literal statements  
type StringParser struct{}

func (sp *StringParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokString
}

func (sp *StringParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	stringStmt := &StringStmt{
		Value: curTok.Value,
		Line:  curTok.Line,
		Col:   curTok.Column,
	}

	ctx.NextToken() // consume string

	return stringStmt
}
