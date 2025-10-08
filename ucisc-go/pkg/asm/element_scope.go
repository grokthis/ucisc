package asm

// ScopeParser parses scope block statements
type ScopeParser struct{}

func (sp *ScopeParser) Matches(p *Parser) bool {
	return p.curTok.Type == TokLBrace
}

func (sp *ScopeParser) Parse(p *Parser) Statement {
	return p.parseScope()
}
