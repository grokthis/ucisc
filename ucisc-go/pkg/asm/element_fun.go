package asm

// FunParser parses fun statements
type FunParser struct{}

func (fp *FunParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokFun
}

func (fp *FunParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	funStmt := &FunStmt{
		Line:    curTok.Line,
		Col:     curTok.Column,
		Args:    []string{},
		Returns: []string{},
		Body:    []Statement{},
	}

	ctx.NextToken() // consume 'fun'

	// Parse qualified name: stack.functionName
	curTok = ctx.CurrentToken()
	if curTok.Type != TokIdent {
		ctx.AddError("expected identifier after 'fun'")
		return nil
	}

	funStmt.DefName = curTok.Value
	ctx.NextToken()

	// Expect dot
	if ctx.CurrentToken().Type != TokDot {
		ctx.AddError("expected '.' after stack name in function declaration")
		return nil
	}
	ctx.NextToken() // consume '.'

	// Parse function name
	curTok = ctx.CurrentToken()
	if curTok.Type != TokIdent {
		ctx.AddError("expected function name after '.'")
		return nil
	}
	funStmt.Name = curTok.Value
	ctx.NextToken()

	// Parse arguments (arg1, arg2, ...)
	if ctx.CurrentToken().Type == TokLParen {
		ctx.NextToken() // consume '('

		for ctx.CurrentToken().Type != TokRParen && ctx.CurrentToken().Type != TokEOF {
			if ctx.CurrentToken().Type != TokIdent {
				ctx.AddError("expected argument name")
				return nil
			}
			funStmt.Args = append(funStmt.Args, ctx.CurrentToken().Value)
			ctx.NextToken()

			if ctx.CurrentToken().Type == TokComma {
				ctx.NextToken()
			}
		}

		if ctx.CurrentToken().Type != TokRParen {
			ctx.AddError("expected ')' after arguments")
			return nil
		}
		ctx.NextToken() // consume ')'
	}

	// Parse return values (-> ret1, ret2)
	if ctx.CurrentToken().Type == TokArrow {
		ctx.NextToken() // consume '->'

		for ctx.CurrentToken().Type == TokIdent {
			funStmt.Returns = append(funStmt.Returns, ctx.CurrentToken().Value)
			ctx.NextToken()

			if ctx.CurrentToken().Type == TokComma {
				ctx.NextToken()
			} else {
				break
			}
		}
	}

	// Parse body { ... }
	if ctx.CurrentToken().Type != TokLBrace {
		ctx.AddError("expected '{' to start function body")
		return nil
	}

	// Before parsing the body, we need to set up the function's stack frame
	// Push a new scope for the function
	ctx.NextToken() // consume '{'
	ctx.Symbols().PushScope()

	// According to uCISC calling convention:
	// Stack layout (from stack pointer upward):
	//   offset 0: arguments (pushed last, so at top of stack)
	//   offset 1: return address
	//   offset 2+: return values (reserved by caller, pushed first)

	// Define variables for arguments, return address, and return values
	offset := int16(0)

	// Arguments are at the top (offset 0+)
	for _, argName := range funStmt.Args {
		if err := ctx.Symbols().DefineVar(funStmt.DefName, argName, offset, funStmt.Line, funStmt.Col); err != nil {
			ctx.AddError(err.Error())
			return nil
		}
		offset++
	}

	// Return address is next
	if err := ctx.Symbols().DefineVar(funStmt.DefName, "return", offset, funStmt.Line, funStmt.Col); err != nil {
		ctx.AddError(err.Error())
		return nil
	}
	offset++

	// Return values are at the bottom (highest offsets)
	for _, retName := range funStmt.Returns {
		if err := ctx.Symbols().DefineVar(funStmt.DefName, retName, offset, funStmt.Line, funStmt.Col); err != nil {
			ctx.AddError(err.Error())
			return nil
		}
		offset++
	}

	// Parse statements until '}'
	for ctx.CurrentToken().Type != TokRBrace && ctx.CurrentToken().Type != TokEOF {
		// Skip empty lines and comments
		if ctx.CurrentToken().Type == TokNewline || ctx.CurrentToken().Type == TokComment {
			ctx.NextToken()
			continue
		}

		stmt := ctx.ParseStatement()
		if stmt != nil {
			funStmt.Body = append(funStmt.Body, stmt)
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
		ctx.AddError("expected '}' to close function body")
		return nil
	}
	ctx.NextToken() // consume '}'

	// Pop the function scope
	ctx.Symbols().PopScope()

	return funStmt
}
