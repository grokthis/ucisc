package asm

// DefParser parses def statements
type DefParser struct{}

func (dp *DefParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokDef
}

func (dp *DefParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	defStmt := &DefStmt{
		Line: curTok.Line,
		Col:  curTok.Column,
	}

	ctx.NextToken() // consume 'def'

	// Parse name/register
	curTok = ctx.CurrentToken()
	if curTok.Type != TokIdent {
		ctx.AddError("expected identifier after 'def'")
		return nil
	}
	defStmt.Name = curTok.Value
	ctx.NextToken()

	// Expect /
	if ctx.CurrentToken().Type != TokSlash {
		ctx.AddError("expected '/' after def name")
		return nil
	}
	ctx.NextToken()

	// Parse register
	reg, ok := ctx.ParseRegister()
	if !ok {
		ctx.AddError("expected register after '/'")
		return nil
	}
	defStmt.Register = reg

	// Register the definition in symbol table
	if err := ctx.Symbols().Define(defStmt.Name, defStmt.Register, defStmt.Line, defStmt.Col); err != nil {
		ctx.AddError(err.Error())
		return nil
	}

	// Check for optional initializer
	if ctx.CurrentToken().Type == TokEffectAlways || ctx.CurrentToken().Type == TokEquals {
		if ctx.CurrentToken().Type == TokEquals {
			// Convert = to <- for consistency
			ctx.NextToken() // consume =
		}
		// Parse the instruction
		defStmt.Initializer = ctx.ParseInstructionFromDef(defStmt.Name, defStmt.Register)
	}

	return defStmt
}
