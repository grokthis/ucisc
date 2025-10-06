// Package disasm provides disassembly of uCISC instructions
package disasm

import (
	"fmt"
	"github.com/grokthis/ucisc-go/pkg/cpu"
)

// FormatInstruction formats an instruction as assembly text
func FormatInstruction(inst *cpu.Instruction) string {
	// Format destination
	dest := formatDestination(inst.Destination, inst.Offset, inst.Increment)

	// Format effect
	effect := formatEffect(inst.Effect)

	// Format operation
	op := formatOperation(inst.Operation)

	// Format source
	src := formatSource(inst.Source, inst.Immediate, inst.Increment)

	return fmt.Sprintf("%s %s %s %s", dest, effect, op, src)
}

func formatDestination(dest cpu.Destination, offset uint8, increment bool) string {
	push := ""
	if increment && dest.Type == cpu.DestMem {
		push = " push"
	}

	if dest.Type == cpu.DestReg {
		return fmt.Sprintf("&%s%s", formatRegister(dest.Reg), push)
	}

	// DestMem
	if offset == 0 {
		return fmt.Sprintf("%s%s", formatRegister(dest.Reg), push)
	}
	return fmt.Sprintf("%s/%d%s", formatRegister(dest.Reg), offset, push)
}

func formatSource(src cpu.Source, immediate int16, increment bool) string {
	pop := ""
	if increment && src.Type == cpu.SourceMem {
		pop = " pop"
	}

	switch src.Type {
	case cpu.SourceVal:
		return fmt.Sprintf("val/%d%s", immediate, pop)
	case cpu.SourceReg:
		if immediate == 0 {
			return fmt.Sprintf("&%s%s", formatRegister(src.Reg), pop)
		}
		return fmt.Sprintf("&%s/%d%s", formatRegister(src.Reg), immediate, pop)
	case cpu.SourceMem:
		if immediate == 0 {
			return fmt.Sprintf("%s%s", formatRegister(src.Reg), pop)
		}
		return fmt.Sprintf("%s/%d%s", formatRegister(src.Reg), immediate, pop)
	case cpu.SourcePC:
		if immediate == 0 {
			return fmt.Sprintf("pc%s", pop)
		}
		return fmt.Sprintf("pc/%d%s", immediate, pop)
	}
	return "???"
}

func formatRegister(reg cpu.Register) string {
	switch reg {
	case cpu.RegPC:
		return "pc"
	case cpu.RegR1:
		return "r1"
	case cpu.RegR2:
		return "r2"
	case cpu.RegR3:
		return "r3"
	case cpu.RegR5:
		return "r5"
	case cpu.RegR6:
		return "r6"
	case cpu.RegR7:
		return "r7"
	case cpu.RegBanking:
		return "banking"
	case cpu.RegFlags:
		return "flags"
	case cpu.RegInt:
		return "interrupt"
	}
	return "???"
}

func formatEffect(effect cpu.Effect) string {
	switch effect {
	case cpu.EffectAlways:
		return "<-"
	case cpu.EffectZero:
		return "<0"
	case cpu.EffectNotZero:
		return "<1"
	case cpu.EffectNegative:
		return "<n"
	case cpu.EffectPositive:
		return "<p"
	case cpu.EffectOverflow:
		return "<o"
	case cpu.EffectInterrupt:
		return "<#"
	case cpu.EffectFlagsOnly:
		return "<~"
	}
	return "<?>"
}

func formatOperation(op cpu.Operation) string {
	switch op {
	case cpu.OpCopy:
		return "copy"
	case cpu.OpAnd:
		return "and"
	case cpu.OpOr:
		return "or"
	case cpu.OpXor:
		return "xor"
	case cpu.OpInv:
		return "inv"
	case cpu.OpShl:
		return "shl"
	case cpu.OpShr:
		return "shr"
	case cpu.OpSwap:
		return "swap"
	case cpu.OpMsb:
		return "msb"
	case cpu.OpLsb:
		return "lsb"
	case cpu.OpAdd:
		return "add"
	case cpu.OpAddc:
		return "addc"
	case cpu.OpSub:
		return "sub"
	case cpu.OpMult:
		return "mult"
	case cpu.OpMultMsw:
		return "multmsw"
	}
	return "???"
}
