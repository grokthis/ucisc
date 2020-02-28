# uCISC Syntax

uCISC syntax is heavily derived from [Mu](https://github.com/akkartik/mu). However,
I made a few changes to make the syntax a bit more visually readable in my opinion.

## Instructions

The basic instruction format is:

```
  D/move/ 1.reg/stack address/ 'to 2.reg/return address/ 0x2.imm/offset/
```

#### Codes and Values

Every code or argument has a type. The compiler will validate that each instruction only
includes valid types and raise a compile error otherwise. Each non-comment position takes
the following format:

```
  <code>.<type>
```

The code is the op code, argument or value in case insensitive hexadecimal. The type is a
valid type abbreviation (see below). You can include the "0x" prefix if desired to clearly
denote the hexadecimal values. Decimal values can be specified by adding a 'd' at the end
of the number (for example `-19d` or `24d`).

#### Position Within Statements

In uCISC statements, position is mostly ignored, except for the following:

1. The first position must be the instruction code. The type specification can be left off
  of the instruction code if desired.

2. Source and destination arguments must be in the order of data movement direction (e.g.
  `<from> <to>`). The compiler will be able to correctly infer the copy direction, source
  and destination details based on the position and types given. This cleans up the
  statements quite a bit while maintaining clarity and readability.

3. When instructions use the immediate value to offset a register address to reference
   the actual address, the immediate value must be located in the next position after
   the address it is modifying (the source or destination).

#### Position Types

Source and destination types:
 - `reg` - register, the content of the register is being referenced
 - `mem` - memory location, the content of the memory referenced by the register is being referenced
 - `val` - immediate value, the content of the immediate value is being referenced

Argument types:
 - `imm` - 8-bit immediate value, if label resolves to the label offset in the page
 - `disp` - immediate value, if label and will resolve to difference between the current instruction and the label. Must result in a valid imm value for the given instruction.
 - `aluc` - ALU op code, values between 00 and 1F
 - `eff` - Effect code, values between 0 and 4
 - `inc` - Word increment, value of 0 or 1
 - `sign` - Sign extension mode, value of 0 or 1
 - `dir` - Paging direction, value of 0 or 1

#### Comments

There are two types of comments:
 - word comments begin with a "'" character and terminate at the next
   whitespace (e.g. `'comment <arg>`).
 - inline comments begin with a '/' character and terminate at the next
   '/' character (e.g.  `<arg>/comment/ <arg>`)
 - line comments are lines where the first non-whitespace character is '#'.
   Line comments cause the entire line to be ignored.

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

