package asm

// VarParser parses var statements
type VarParser struct{}

func (vp *VarParser) Matches(p *Parser) bool {
	return p.curTok.Type == TokVar
}

func (vp *VarParser) Parse(p *Parser) Statement {
	return p.parseVar()
}
