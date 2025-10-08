package asm

// LabelParser parses label statements
type LabelParser struct{}

func (lp *LabelParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokIdent && ctx.PeekToken().Type == TokColon
}

func (lp *LabelParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	label := &LabelStmt{
		Name: curTok.Value,
		Line: curTok.Line,
		Col:  curTok.Column,
	}

	ctx.NextToken() // consume label name
	ctx.NextToken() // consume colon

	return label
}
