package asm

// ScopeParser parses scope block statements
type ScopeParser struct{}

func (sp *ScopeParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokLBrace
}

func (sp *ScopeParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	scopeStmt := &ScopeStmt{
		Line:       curTok.Line,
		Col:        curTok.Column,
		Statements: []Statement{},
	}

	ctx.NextToken() // consume '{'

	// Push a new scope
	ctx.Symbols().PushScope()

	// Parse statements until '}'
	for ctx.CurrentToken().Type != TokRBrace && ctx.CurrentToken().Type != TokEOF {
		// Skip empty lines and comments
		if ctx.CurrentToken().Type == TokNewline || ctx.CurrentToken().Type == TokComment {
			ctx.NextToken()
			continue
		}

		stmt := ctx.ParseStatement()
		if stmt != nil {
			scopeStmt.Statements = append(scopeStmt.Statements, stmt)
		} else {
			// Parse error occurred, skip to next line to avoid infinite loop
			for ctx.CurrentToken().Type != TokNewline && ctx.CurrentToken().Type != TokEOF {
				ctx.NextToken()
			}
		}

		// Consume trailing comments and newlines
		for ctx.CurrentToken().Type == TokComment || ctx.CurrentToken().Type == TokNewline {
			ctx.NextToken()
		}
	}

	if ctx.CurrentToken().Type != TokRBrace {
		ctx.AddError("expected '}' to close scope")
		return nil
	}
	ctx.NextToken() // consume '}'

	// Pop the scope
	ctx.Symbols().PopScope()

	return scopeStmt
}
