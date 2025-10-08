package asm

// DefParser parses def statements
type DefParser struct{}

func (dp *DefParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokDef
}

func (dp *DefParser) Parse(ctx ParseContext) Statement {
	// For simplicity during transition, cast to *Parser
	// TODO: Migrate to use only ParseContext methods
	p := ctx.(*Parser)
	return p.parseDef()
}
