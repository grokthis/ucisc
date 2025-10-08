package asm

// InstructionParser parses instruction statements (catch-all)
type InstructionParser struct{}

func (ip *InstructionParser) Matches(p *Parser) bool {
	// Always matches - this should be last in the list
	return true
}

func (ip *InstructionParser) Parse(p *Parser) Statement {
	inst := &InstructionStmt{
		Line: p.curTok.Line,
		Col:  p.curTok.Column,
	}
	
	// Parse destination
	dest, push := p.parseDestination()
	if dest == nil {
		return nil
	}
	inst.Destination = dest
	inst.Push = push
	
	// Parse effect
	effect := p.parseEffect()
	inst.Effect = effect
	
	// Parse operation
	op := p.parseOperation()
	inst.Operation = op
	
	// Parse source
	source, pop := p.parseSource()
	if source == nil {
		return nil
	}
	inst.Source = source
	inst.Pop = pop
	
	return inst
}
