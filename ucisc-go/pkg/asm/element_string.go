package asm

// StringParser parses string literal statements  
type StringParser struct{}

func (sp *StringParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokString
}

func (sp *StringParser) Parse(ctx ParseContext) Statement {
	// For simplicity during transition, cast to *Parser
	// TODO: Migrate to use only ParseContext methods
	p := ctx.(*Parser)
	stringStmt := &StringStmt{
		Value: p.curTok.Value,
		Line:  p.curTok.Line,
		Col:   p.curTok.Column,
	}
	
	p.nextToken() // consume string
	
	return stringStmt
}
