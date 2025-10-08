package asm

// VarParser parses var statements
type VarParser struct{}

func (vp *VarParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokVar
}

func (vp *VarParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	varStmt := &VarStmt{
		Line: curTok.Line,
		Col:  curTok.Column,
	}

	ctx.NextToken() // consume 'var'

	// Parse qualified name: def.varname
	curTok = ctx.CurrentToken()
	if curTok.Type != TokIdent {
		ctx.AddError("expected identifier after 'var'")
		return nil
	}

	varStmt.DefName = curTok.Value
	ctx.NextToken()

	// Expect dot
	if ctx.CurrentToken().Type != TokDot {
		ctx.AddError("expected '.' after def name in var declaration")
		return nil
	}
	ctx.NextToken() // consume '.'

	// Parse variable name
	curTok = ctx.CurrentToken()
	if curTok.Type != TokIdent {
		ctx.AddError("expected variable name after '.'")
		return nil
	}
	varStmt.VarName = curTok.Value
	ctx.NextToken()

	// Expect /
	if ctx.CurrentToken().Type != TokSlash {
		ctx.AddError("expected '/' after var name")
		return nil
	}
	ctx.NextToken()

	// Parse offset
	offset, ok := ctx.ParseNumber()
	if !ok {
		return nil
	}
	varStmt.Offset = int16(offset)

	// Register the variable in symbol table
	if err := ctx.Symbols().DefineVar(varStmt.DefName, varStmt.VarName, varStmt.Offset, varStmt.Line, varStmt.Col); err != nil {
		ctx.AddError(err.Error())
		return nil
	}

	// Check for push modifier
	if ctx.CurrentToken().Type == TokPush {
		varStmt.Push = true
		ctx.NextToken()
	}

	// Check for optional initializer
	if ctx.CurrentToken().Type == TokEffectAlways || ctx.CurrentToken().Type == TokEquals {
		if ctx.CurrentToken().Type == TokEquals {
			ctx.NextToken() // consume =
		}

		// Look up the register for this def
		def, ok := ctx.Symbols().LookupDef(varStmt.DefName)
		if !ok {
			ctx.AddError("undefined definition \"" + varStmt.DefName + "\"")
			return nil
		}

		// Parse instruction with def.varname as destination
		curTok = ctx.CurrentToken()
		inst := &InstructionStmt{
			Line: curTok.Line,
			Col:  curTok.Column,
			Destination: &Operand{
				Type:     OpMemory,
				Register: def.Register,
				Offset:   varStmt.Offset,
			},
			Push: varStmt.Push,
		}

		// Parse effect
		inst.Effect = ctx.ParseEffect()

		// Parse operation
		inst.Operation = ctx.ParseOperation()

		// Parse source
		source, pop := ctx.ParseSource()
		if source == nil {
			return nil
		}
		inst.Source = source
		inst.Pop = pop

		varStmt.Initializer = inst
	}

	return varStmt
}
