package asm

import "fmt"

// TokenType represents the type of a lexical token
type TokenType int

const (
	// Special tokens
	TokEOF TokenType = iota
	TokIllegal
	TokComment

	// Literals
	TokIdent  // label, var_name
	TokNumber // 42, 0x1234, %FFFF

	// Keywords
	TokVal   // val (immediate value pseudo-register)
	TokDef   // def (register definition)
	TokVar   // var (variable offset)
	TokFun   // fun (function)
	TokDebug // debug (debug device)

	// Registers
	TokR1
	TokR2
	TokR3
	TokR5
	TokR6
	TokR7
	TokPC
	TokControl
	TokBanking
	TokFlags
	TokInterrupt

	// Operations (ALU)
	TokCopy
	TokAnd
	TokOr
	TokXor
	TokInv
	TokShl
	TokShr
	TokSwap
	TokMsb
	TokLsb
	TokAdd
	TokAddc
	TokSub
	TokMult
	TokMultsw

	// Effects (conditionals)
	TokEffectAlways    // <-, <|
	TokEffectFlags     // <~, <|~
	TokEffectZero      // <0, <|0
	TokEffectNotZero   // <1, <|1
	TokEffectPositive  // <+, <p, <|+, <|p
	TokEffectNegative  // <n, <-, <|n, <|-
	TokEffectOverflow  // <&, <o, <|&, <|o
	TokEffectInterrupt // <#, <|#

	// Modifiers
	TokPush
	TokPop

	// Punctuation
	TokSlash      // /
	TokAmpersand  // &
	TokColon      // :
	TokNewline    // \n
	TokQuestion   // ?
	TokLBrace     // {
	TokRBrace     // }
	TokLParen     // (
	TokRParen     // )
	TokComma      // ,
	TokArrow      // ->
	TokPercent    // % (for data blocks)
	TokString     // "string literal"
	TokEquals     // =
	TokDot        // .
)

var tokenNames = map[TokenType]string{
	TokEOF:     "EOF",
	TokIllegal: "ILLEGAL",
	TokComment: "COMMENT",
	TokIdent:   "IDENT",
	TokNumber:  "NUMBER",

	TokVal:   "val",
	TokDef:   "def",
	TokVar:   "var",
	TokFun:   "fun",
	TokDebug: "debug",

	TokR1:        "r1",
	TokR2:        "r2",
	TokR3:        "r3",
	TokR5:        "r5",
	TokR6:        "r6",
	TokR7:        "r7",
	TokPC:        "pc",
	TokControl:   "control",
	TokBanking:   "banking",
	TokFlags:     "flags",
	TokInterrupt: "interrupt",

	TokCopy:    "copy",
	TokAnd:     "and",
	TokOr:      "or",
	TokXor:     "xor",
	TokInv:     "inv",
	TokShl:     "shl",
	TokShr:     "shr",
	TokSwap:    "swap",
	TokMsb:     "msb",
	TokLsb:     "lsb",
	TokAdd:     "add",
	TokAddc:    "addc",
	TokSub:     "sub",
	TokMult:    "mult",
	TokMultsw:  "multsw",

	TokEffectAlways:    "EFFECT_ALWAYS",
	TokEffectFlags:     "EFFECT_FLAGS",
	TokEffectZero:      "EFFECT_ZERO",
	TokEffectNotZero:   "EFFECT_NOT_ZERO",
	TokEffectPositive:  "EFFECT_POSITIVE",
	TokEffectNegative:  "EFFECT_NEGATIVE",
	TokEffectOverflow:  "EFFECT_OVERFLOW",
	TokEffectInterrupt: "EFFECT_INTERRUPT",

	TokPush: "push",
	TokPop:  "pop",

	TokSlash:      "/",
	TokAmpersand:  "&",
	TokColon:      ":",
	TokNewline:    "NEWLINE",
	TokQuestion:   "?",
	TokLBrace:     "{",
	TokRBrace:     "}",
	TokLParen:     "(",
	TokRParen:     ")",
	TokComma:      ",",
	TokArrow:      "->",
	TokPercent:    "%",
	TokString:     "STRING",
	TokEquals:     "=",
	TokDot:        ".",
}

// Token represents a lexical token
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// String returns a string representation of the token
func (t Token) String() string {
	name, ok := tokenNames[t.Type]
	if !ok {
		name = fmt.Sprintf("TokenType(%d)", t.Type)
	}
	if t.Value != "" && t.Type != TokNewline {
		return fmt.Sprintf("%s(%s)@%d:%d", name, t.Value, t.Line, t.Column)
	}
	return fmt.Sprintf("%s@%d:%d", name, t.Line, t.Column)
}

// Keywords maps keyword strings to token types
var keywords = map[string]TokenType{
	"val":   TokVal,
	"def":   TokDef,
	"var":   TokVar,
	"fun":   TokFun,
	"debug": TokDebug,

	"r1": TokR1,
	"r2": TokR2,
	"r3": TokR3,
	"r5": TokR5,
	"r6": TokR6,
	"r7": TokR7,
	"pc": TokPC,

	"control":   TokControl,
	"banking":   TokBanking,
	"flags":     TokFlags,
	"interrupt": TokInterrupt,

	"copy":    TokCopy,
	"and":     TokAnd,
	"or":      TokOr,
	"xor":     TokXor,
	"inv":     TokInv,
	"shl":     TokShl,
	"shr":     TokShr,
	"swap":    TokSwap,
	"msb":     TokMsb,
	"lsb":     TokLsb,
	"add":     TokAdd,
	"addc":    TokAddc,
	"sub":     TokSub,
	"mult":    TokMult,
	"multsw":  TokMultsw,

	"push": TokPush,
	"pop":  TokPop,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TokIdent
}
