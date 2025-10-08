package asm

// StringParser parses string literal statements  
type StringParser struct{}

func (sp *StringParser) Matches(p *Parser) bool {
	return p.curTok.Type == TokString
}

func (sp *StringParser) Parse(p *Parser) Statement {
	stringStmt := &StringStmt{
		Value: p.curTok.Value,
		Line:  p.curTok.Line,
		Col:   p.curTok.Column,
	}
	
	p.nextToken() // consume string
	
	return stringStmt
}
