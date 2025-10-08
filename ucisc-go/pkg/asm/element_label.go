package asm

// LabelParser parses label statements
type LabelParser struct{}

func (lp *LabelParser) Matches(p *Parser) bool {
	return p.curTok.Type == TokIdent && p.peekTok.Type == TokColon
}

func (lp *LabelParser) Parse(p *Parser) Statement {
	label := &LabelStmt{
		Name: p.curTok.Value,
		Line: p.curTok.Line,
		Col:  p.curTok.Column,
	}
	
	p.nextToken() // consume label name
	p.nextToken() // consume colon
	
	return label
}
