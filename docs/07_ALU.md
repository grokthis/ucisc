# uCISC ALU Instructions

ALU instructions perform arithmetic and logic operations on data.

```
20A/ALU add/ 2.mem/my var/ 1.mem/stack val/ 1.inc 0.eff/store if zero/

# Packed Instruction Bits:
10SDDRRR MAAAAAEE
```

Note: opcode is instruction opcode concatenated with ALU opcode. Value
is in the range 200-21F

#### Arguments

*(R) Register Arg*

0.reg - Value in PC

1.mem - Value at memory location r1
2.mem - Value at memory location r2
3.mem - Value at memory location r3

4.reg - Value in flags register
5.reg - Value in r1
6.reg - Value in r2
7.reg - Value in r3

*(D) Destination*

0.reg - Value in PC

1.mem - Value at memory location r1
2.mem - Value at memory location r2
3.mem - Value at memory location r3

*(M) Incre(M)ent on Modify*
0.inc - No increment
1.inc - Increment source and destination address registers by 2

This works on any X.mem arguments and increments the corresponding X.reg
address values by 2 on completion of this operation. This allows you to
chain math operations to succinctly operate on numbers larger than 16-bits
(in multiples of 16-bits) by incrementing the address. For example, 2 ADD
operations in a row with increment turned on will add a 32-bit number,
leaving the overflow and other flags set properly.

Note: The increment only occurs if the effect causes a store operation

*(S) Sign*

If M == 1, interpret the sign as the increment sign.

0.sign - positive increment
1.sign - negative increment

If M == 0, interpret the sign as modifying the Effect instead. See below.

*(E) Effect*

If M == 1 and/or S == 0

0.eff - store if zero
1.eff - store if not zero
2.eff - store if positive
3.eff - store

If M == 0 and S == 1

4.eff - store if not negative
5.eff - store if negative
6.eff - store if not overflow
7.eff - do not store

Note: 4-7 are stored as 0-3 in the instruction code. However, in uCISC to avoid ambiguity
and allow compile time checking, 4-7 are used. If 4-7 don't fit given the values of M and S
the compiler can throw an error rather than blindly trust you know what you are doing.

*(A) ALU opcode*

Math functions that take a single arg (e.g. INV) use the source, perform
the operation and store the result in the target. Math functions that
require two arguments, use both the source and target as the arguments
and store the result in the target.

Argument order: All single argument logice uses the source as the
input and the target as the output (e.g. NOP, INV, etc.).
All multi-arg operations use the target as the first arg and the
source as the second.

For example:

```
# source is 2.mem
# target is 1.mem (the stack)
# This means 1.mem - 2.mem => 1.mem
20C/subtract/ 2.mem/increment/ 1.mem/stack val/ 0.inc 3.eff/store/
```

This is important for repeated operation where you want to do
things like decrement a value repeatedly. You can run the example
instruction above repeatedly to subtract the 2.mem value each time
rather than having to reset the decrement value each time. The
same thing works for all multi-arg ops.

00 - NOP
01 - AND
02 - OR
03 - XOR
04 - INV
05 - Parity?

06 - Shift Left
07 - Shift Right Zero Extend
08 - Shift Right Sign Extend

09 - Swap MSB and LSB bytes

0A - Add Unsigned
0B - Add Signed
0C - Subtract Unsigned
0D - Subtract Signed

0E - Multiply Unsigned
0F - Multiply Signed
10 - Divide Unsigned
11 - Divide Signed
12 - Mod Unsigned
13 - Mod Signed
  Note: if carry flag is set, uses it to extend arithmetic, requires you
  to clear the carry flag when doing multiple independent math
  operations back to back.

  Note: overflow != carry, carry lets you quickly chain similar functions
  to do math on numbers with some multiple of 16 bits. Overflow
  indicates a numeric overflow if the result was the MSB.

14 -
15 -
16 - Add floating point
17 - Subtract floating point
18 - Multiply floating point
19 - Divide floating point
1A - Exponent from floating point?
1B - Whole part from floating point?
1C - Fractional part from floating point?

1D - <unused>
1E - <unused>
1F - <unused>

