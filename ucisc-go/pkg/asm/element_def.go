package asm

// DefParser parses def statements
type DefParser struct{}

func (dp *DefParser) Matches(p *Parser) bool {
	return p.curTok.Type == TokDef
}

func (dp *DefParser) Parse(p *Parser) Statement {
	return p.parseDef()
}
