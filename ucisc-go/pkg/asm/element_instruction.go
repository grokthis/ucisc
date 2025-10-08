package asm

// InstructionParser parses instruction statements (catch-all)
type InstructionParser struct{}

func (ip *InstructionParser) Matches(ctx ParseContext) bool {
	// Always matches - this should be last in the list
	return true
}

func (ip *InstructionParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	inst := &InstructionStmt{
		Line: curTok.Line,
		Col:  curTok.Column,
	}

	// Parse destination
	dest, push := ctx.ParseDestination()
	if dest == nil {
		return nil
	}
	inst.Destination = dest
	inst.Push = push

	// Parse effect
	effect := ctx.ParseEffect()
	inst.Effect = effect

	// Parse operation
	op := ctx.ParseOperation()
	inst.Operation = op

	// Parse source
	source, pop := ctx.ParseSource()
	if source == nil {
		return nil
	}
	inst.Source = source
	inst.Pop = pop

	return inst
}
