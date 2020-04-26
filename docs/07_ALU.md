# uCISC ALU Instructions

ALU instructions perform arithmetic and logic operations on data.

```
20A/ALU add/ 2.mem/my var/ 1.mem/stack val/ 1.inc 0.eff/store if zero/

# Packed Instruction Bits:
1EEDDDSS SAAAAAGM
```

Note: for uCISC syntax, the ALU opcode (00-1F) is concatenated with the instruction
opcode to indicate which math operation is being performed. Therefore, the coded
opcode value is in the range 200-21F.

#### Flags Register Effects

The flags register effects apply to this instruction. See the Copy instruction
for more details.

#### Special Case: NOP

For opcode 200, you can produce a NOP instruction by making the CPU do nothing
and incrementing the PC to the next instruction with 1.inc.

```
200/ALU copy/ 0.pc 0.pc 1.inc
```

#### Arguments

*(S) Source Argument*

ALU instructions use the standard source instruction behavior. ALU instructions support
the 4.reg argument since there is no immediate value.

*(D) Destination*

ALU instructions use the standard destination instruction behavior.

*(M) Incre(M)ent*

* 0.inc - No increment
* 1.inc - Increment source and destination address registers by 2

This works on any X.mem arguments and increments the corresponding X.reg address values
by 1 on completion of this operation. This allows you to do things like chain ALU operations
to succinctly operate on numbers larger than 16-bits (in multiples of 16-bits). Another
common use case is to execute a jump and pop operation by sending the stack value to the PC
and incrementing the stack pointer.

Note: The increment only occurs if the effect causes a store operation

If both arguments are reg arguments, the increment argument has no effect. This behavior may
change in the future, so don't rely on it.

*(G) Sign*

If M == 1, interpret the sign as the increment sign.

* 0.sign - positive increment
* 1.sign - negative increment

If M == 0, interpret the sign as modifying the Effect instead. See below.

*(E) Effect*

If M == 0, G is not needed to specify the increment sign and is treated as the MSB of the
effect number, allowing you to access the extended effects.

*(A) ALU opcode*

Math functions that take a single arg (e.g. inverse, increment, decrement, negate, byte swaps)
use the source as the argument and store the result in the target. Math functions that
require two arguments, use both the source and target as the arguments and store the result in
the target.

*Note:* All 2-arg operations use the target as the first arg and the source as the second.
For example, subtraction is D - S => D (subtract S from D). Shift uses D >> S => D (shift D by
S bits) and store in D. The value in D is destroyed in the process.

For example:

```
# source is 2.mem
# target is 1.mem (the stack)
# This means 1.mem - 2.mem => 1.mem
20C/subtract/ 2.mem 1.mem/SP/
```

This allows you to repeatedly modify a number in some way that. For example if you want to shift
a number a bit at a time and take action on the bit, you can do so by running the same shift
multiple times. You can also use this to subtract a constant value from a larger number multiple
times. You can run the example instruction above repeatedly to subtract the 2.mem value each time
rather than having to reset the decrement value each time.

#### Bit operations

Flags are set based on 16-bit result value, overflow is always 0, carry register is set to 0.

* 00 - Pass through, no modification
* 01 - AND
* 02 - OR
* 03 - XOR
* 04 - INV

Flags are set based on 16-bit result value, overflow is 1 if any bits shifted out were 1's.
Carry register LS bits are shifted out values in the order they were shifted out. The
LS bit is the first shifted bit, the next LS bit is the second shifted bit, etc. When chaining
shift operations, shifted in bits are pulled LS bit to MS bit from carry register. This
allows you to chain both shift left and shift right operations.

* 05 - Shift left carry extended
* 06 - Shift right carry extended
* 07 - Shift right sign extend (does carry extended when overflow is set)

#### Arithmetic operations

Flags are set based on 16-bit result value, overflow indicates multiplication overflow. Overflow
register is set to multiplication overflow value (bits 31:16).

* 08 - Multiply Unsigned
* 09 - Multiply Signed

Flags are set based on 16-bit result value, overflow indicates multiplication overflow. Overflow
register is set to 1 on carry (not overflow).

* 0A - Add Unsigned
* 0B - Add Signed
* 0C - Subtract Unsigned
* 0D - Subtract Signed

* 0E - Increment
* 0F - Decrement

Clear flags/overflow operation
* 10 - Clear all flags (set for zero result), clear carry register

Flags are set based on 16-bit result value, overflow is 1 on carry, respects carry in
* 11 - Negate (two's complement)

#### Byte operations

Since the memory is word addressed, uCISC provides a number of byte specific
operations to improve performance and ease of byte based logic.

Flags are set based on 16-bit result value, overflow is always 0, carry register is set to 0.

* 12 - Swap MSB and LSB bytes
* 13 - MSB to LSB: (A >> 8, ignores overflow)
* 14 - Zero LSB:  (A & 0xFF00)
* 15 - Zero MSB:  (A & 0x00FF)
* 16 - Write MSB only: (A & 0xFF00) | (B & 0x00FF)
* 17 - Write LSB only: (A & 0x00FF) | (B & 0xFF00)
* 18 - Write LSB as MSB: ((A & 0x00FF) << 8) | (B & 0x00FF)
* 19 - Write MSB as LSB: ((A & 0xFF00) >> 8) | (B & 0xFF00)

Unused ALU operations

* 1A -
* 1B -
* 1C -
* 1D -
* 1E -
* 1F -

