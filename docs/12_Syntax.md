# uCISC Syntax

uCISC syntax is heavily influenced by [Mu](https://github.com/akkartik/mu).

If you use vim, there is syntax highlighting available:
[extras/ucisc.vim](/extras/ucisc.vim). 

## Instructions

The basic instruction format is:

```
  # Copy an address off the stack to a the PC (effectively a jump instruction)
  0/copy/ 1.mem/stack variable/ 2.imm/stack offset/ 0.reg/PC register/
```

#### Codes and Values

Every code or argument has a type. The compiler will validate that each instruction only
includes valid types and raise a compile error otherwise. The format for arguments is:

```
  <code>.<type>

  # Register and memory examples
  1.mem # memory value at the address in r1
  1.reg # register value in r1
  4.imm # immediate value as an argument
  0.reg # PC register
  4.reg # PC or flags register depending on context

  # Other argument types
  3.eff # the conditional effect of the instruction
  1.inc # enables increment on an instruction
  1.push # decrements destination register before storing (for pushing to the stack)
  1.sign # the sign of an argument (effect depends on context, usually affects inc direction)

  # Labels are interpreted as immediate values
  start.imm # Absolute address reference of the start label
  start.disp # Relative reference from the current instruction to the start label 
  
```

The code is the op code, argument or value in case insensitive hexadecimal. The type is a
valid type abbreviation (see below). You can include the "0x" prefix if desired to clearly
denote the hexadecimal values. Decimal values can be specified by adding a 'd' at the end
of the number (for example `-19d` or `24d`).

Recommended: It's a good idea to include the "0x" prefix on immediate values. You'll get
less confused when 12.imm is actually understood as 18 in decimal.

#### Position Within Statements

In uCISC statements, position is mostly ignored, except for the following:

1. The first position must be the instruction code. The type specification can be left off
  of the instruction code if desired.

2. Source and destination arguments must be in the order of data movement direction (e.g.
  `<from> <to>`).

3. The immediate value must follow the source argument, not the destination argument.

#### Position Types

Source and destination types:
 - `reg` - register, the content of the register is being referenced
 - `mem` - memory location, the content of the memory referenced by the register is being referenced
 - `val` - immediate value, the content of the immediate value is being referenced

Valid types for source and destination depend on the instruction. Generally, `val` is
only valid as a source, `mem` is almost always valid as either and `reg` depends on
the details of the instruction.

Argument types:
 - `imm` - 6-bit or 7-bit signed immediate value, depending on context
 - `disp` - immediate value, if label and will resolve to difference between the current instruction and the label. Must result in a valid imm value for the given instruction.
 - `eff` - Effect code, values between 0 and 4
 - `inc` - Word increment, value of 0 or 1
 - `sign` - Sign extension mode, value of 0 or 1

#### Comments

There are three types of comments:
 - word comments begin with a "'" character and terminate at the next
   whitespace (e.g. `'comment <arg>`).
 - inline comments begin with a '/' character and terminate at the next
   '/' character (e.g.  `<arg>/comment/ <arg>`) on the same line
 - line comments start with a '#' character and run all the way to the
   newline.

Before processing a statement, the compiler will strip all comments from the line. Be
careful about inline comments that cause args to collide when removed (e.g.
`<arg>/comment/<arg>`). You will want to include a whitespace on one side of the inline
comment.

The following two statements are identical:

```
  # Readable statement:
  D/move/ 1.reg/stack address/ 'to 2.reg/return address/ 0x2.imm/offset/

  # Without inline or word comments:
  D 1.reg 2.reg 0x2.imm
```

