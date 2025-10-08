package asm

// VarParser parses var statements
type VarParser struct{}

func (vp *VarParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokVar
}

func (vp *VarParser) Parse(ctx ParseContext) Statement {
	// For simplicity during transition, cast to *Parser
	// TODO: Migrate to use only ParseContext methods
	p := ctx.(*Parser)
	return p.parseVar()
}
