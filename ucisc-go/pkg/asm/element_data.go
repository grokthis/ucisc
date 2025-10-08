package asm

import (
	"fmt"
	"strconv"
)

// DataParser parses data block statements
type DataParser struct{}

func (dp *DataParser) Matches(p *Parser) bool {
	return p.curTok.Type == TokPercent
}

func (dp *DataParser) Parse(p *Parser) Statement {
	dataStmt := &DataStmt{
		Line:   p.curTok.Line,
		Col:    p.curTok.Column,
		Values: []uint16{},
	}
	
	p.nextToken() // consume %
	
	// Parse hex numbers until newline or EOF
	for p.curTok.Type != TokNewline && p.curTok.Type != TokEOF {
		if p.curTok.Type == TokNumber {
			// Parse as hex
			value := p.curTok.Value
			// Remove any prefix
			if len(value) > 0 && value[0] == '%' {
				value = value[1:]
			}
			num, err := strconv.ParseInt(value, 16, 32)
			if err != nil {
				p.addError(fmt.Sprintf("invalid hex number %q: %v", p.curTok.Value, err))
				return nil
			}
			dataStmt.Values = append(dataStmt.Values, uint16(num))
			p.nextToken()
		} else if p.curTok.Type == TokIdent {
			// In data blocks, identifiers like ABCD are treated as hex
			value := p.curTok.Value
			num, err := strconv.ParseInt(value, 16, 32)
			if err != nil {
				p.addError(fmt.Sprintf("invalid hex identifier %q in data block", value))
				return nil
			}
			dataStmt.Values = append(dataStmt.Values, uint16(num))
			p.nextToken()
		} else {
			p.addError(fmt.Sprintf("expected hex number in data block, got %v", p.curTok.Type))
			return nil
		}
	}
	
	return dataStmt
}
