package asm

// ScopeParser parses scope block statements
type ScopeParser struct{}

func (sp *ScopeParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokLBrace
}

func (sp *ScopeParser) Parse(ctx ParseContext) Statement {
	// For simplicity during transition, cast to *Parser
	// TODO: Migrate to use only ParseContext methods
	p := ctx.(*Parser)
	return p.parseScope()
}
