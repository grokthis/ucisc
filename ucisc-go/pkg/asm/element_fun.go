package asm

// FunParser parses fun statements
type FunParser struct{}

func (fp *FunParser) Matches(p *Parser) bool {
	return p.curTok.Type == TokFun
}

func (fp *FunParser) Parse(p *Parser) Statement {
	return p.parseFun()
}
