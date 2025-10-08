package asm

// FunParser parses fun statements
type FunParser struct{}

func (fp *FunParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokFun
}

func (fp *FunParser) Parse(ctx ParseContext) Statement {
	// For simplicity during transition, cast to *Parser
	// TODO: Migrate to use only ParseContext methods
	p := ctx.(*Parser)
	return p.parseFun()
}
