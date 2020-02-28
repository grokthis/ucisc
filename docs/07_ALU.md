# uCISC ALU Instructions

ALU instructions perform arithmetic and logic operations on data.

```
20A/ALU add/ 2.mem/my var/ 1.mem/stack val/ 1.inc 0.eff/store if zero/

# Packed Instruction Bits:
10NDDDRR MAAAAAEE
```

Note: opcode is instruction opcode concatenated with ALU opcode. Value
is in the range 200-21F

#### Arguments

*(R) Register Arg*

0.val (if source) - The value 0x0001
0.reg (if target) - Value in PC

1.mem - Value at memory location r1
2.mem - Value at memory location r2
3.mem - Value at memory location r3

*(D) Destination*

0.reg - Value in PC
1.mem - Value at memory location r1
2.mem - Value at memory location r2
3.mem - Value at memory location r3

4.reg - Value in flags register
5.reg - Value in r1
6.reg - Value in r2
7.reg - Value in r3

*(N) Direction*

Value is transferred from:

0.dir - source to destination
1.dir - destination to source

This argument is inferred from the combination of arguments and
source/destination postions. Valid copy statements have one argument
matching the R options and one matching the D options. The order of
those two arguments determines the value N takes. Also, the immediate
value can only modify the R position.

Therefore the following combinations are not valid:

 - Imm value following 4-7.reg
 - Two 4-7.reg arguments
 - Two imm arguments
 - Two val arguments

*(M) Incre(M)ent on Modify*
0.inc - No increment
1.inc - Increment source and destination address registers by 2

This works on any X.mem arguments and increments the corresponding X.reg
address values by 2 on completion of this operation. This allows you to
chain math operations to succinctly operate on numbers larger than 16-bits
(in multiples of 16-bits) by incrementing the address. For example, 2 ADD
operations in a row with increment turned on will add a 32-bit number,
leaving the overflow and other flags set properly.

Note: The increment always applies, even if the effect fails

*(E) Effect*

0.eff - Set overflow, carry and zero result flags; store value if zero flag already set
1.eff - Set overflow, carry and zero result flags; store value if zero flag not already set
2.eff - Store value only, leave current flags unchanged
3.eff - Set overflow, carry and zero result flags; store value

*(A) ALU opcode*

Math functions that take a single arg (e.g. INV) use the source, perform
the operation and store the result in the target. Math functions that
require two arguments, use both the source and target as the arguments
and store the result in the target.

00 - NOP
01 - AND
02 - OR
03 - XOR
04 - INV
05 - Parity

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
14 - Remainder Signed
15 - Remainder Signed
  Note: if carry flag is set, uses it to extend arithmetic, requires you
  to clear the carry flag when doing multiple independent math
  operations back to back.

  Note: overflow != carry, carry lets you quickly chain similar functions
  to do math on numbers with some multiple of 16 bits. Overflow
  indicates a numeric overflow if the result was the MSB.

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

