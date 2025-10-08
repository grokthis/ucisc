package asm

import (
	"fmt"
	"strconv"
)

// DataParser parses data block statements
type DataParser struct{}

func (dp *DataParser) Matches(ctx ParseContext) bool {
	return ctx.CurrentToken().Type == TokPercent
}

func (dp *DataParser) Parse(ctx ParseContext) Statement {
	curTok := ctx.CurrentToken()
	dataStmt := &DataStmt{
		Line:   curTok.Line,
		Col:    curTok.Column,
		Values: []uint16{},
	}

	ctx.NextToken() // consume %

	// Parse hex numbers until newline or EOF
	for ctx.CurrentToken().Type != TokNewline && ctx.CurrentToken().Type != TokEOF {
		curTok = ctx.CurrentToken()
		if curTok.Type == TokNumber {
			// Parse as hex
			value := curTok.Value
			// Remove any prefix
			if len(value) > 0 && value[0] == '%' {
				value = value[1:]
			}
			num, err := strconv.ParseInt(value, 16, 32)
			if err != nil {
				ctx.AddError(fmt.Sprintf("invalid hex number %q: %v", curTok.Value, err))
				return nil
			}
			dataStmt.Values = append(dataStmt.Values, uint16(num))
			ctx.NextToken()
		} else if curTok.Type == TokIdent {
			// In data blocks, identifiers like ABCD are treated as hex
			value := curTok.Value
			num, err := strconv.ParseInt(value, 16, 32)
			if err != nil {
				ctx.AddError(fmt.Sprintf("invalid hex identifier %q in data block", value))
				return nil
			}
			dataStmt.Values = append(dataStmt.Values, uint16(num))
			ctx.NextToken()
		} else {
			ctx.AddError(fmt.Sprintf("expected hex number in data block, got %v", curTok.Type))
			return nil
		}
	}

	return dataStmt
}
