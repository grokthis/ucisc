# uCISC ALU Instructions

ALU instructions perform arithmetic and logic operations on data.

```
20A/ALU add/ 2.mem/my var/ 1.mem/stack val/ 1.inc 0.eff/store if zero/

# Packed Instruction Bits:
10SMDDRR RAAAAAEE
```

Note: opcode is instruction opcode (2) concatenated with ALU opcode (00-1F).
Therefore, the value is in the range 200-21F.

#### Debug instruction

For opcode 200, with source 0.reg to dest 0.reg, 0.inc and an effect that stores
a result will copy the PC to itself without modification resulting in an infinite
loop. Instead of this, processors may interpret this instruction as a breakpoint.

Note, the same instruction with a 1.inc will increment the PC by one before
storing. This is effectively a NOP instruction with the side effect of setting
the flags register based on the result.

#### Arguments

*(R) Register Arg*

* 0.reg - Value in PC

* 1.mem - Value at memory location r1
* 2.mem - Value at memory location r2
* 3.mem - Value at memory location r3

* 4.reg - Value in flags register
* 5.reg - Value in r1
* 6.reg - Value in r2
* 7.reg - Value in r3

*(D) Destination*

* 0.reg - Value in PC

* 1.mem - Value at memory location r1
* 2.mem - Value at memory location r2
* 3.mem - Value at memory location r3

*(M) Incre(M)ent on Modify*

* 0.inc - No increment
* 1.inc - Increment source and destination address registers by 2

This works on any X.mem arguments and increments the corresponding X.reg
address values by 1 on completion of this operation. This allows you to
chain math operations to succinctly operate on numbers larger than 16-bits
(in multiples of 16-bits) by incrementing the address. For example, 2 ADD
operations in a row with increment turned on will add a 32-bit number,
leaving the overflow and other flags set properly.

Note: The increment only occurs if the effect causes a store operation

If both arguments are reg arguments, instead of having no effect, the
instruction will treat the value of the source register as having been
incremented or decremented by 1 depending on the sign argument.

*(S) Sign*

If M == 1, interpret the sign as the increment sign.

* 0.sign - positive increment
* 1.sign - negative increment

If M == 0, interpret the sign as modifying the Effect instead. See below.

*(E) Effect*

If M == 1 and/or S == 0

* 0.eff - store if zero
* 1.eff - store if not zero
* 2.eff - store if positive
* 3.eff - store

If M == 0 and S == 1

* 4.eff - store if not negative
* 5.eff - store if negative
* 6.eff - store if not overflow
* 7.eff - do not store

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

#### Bit operations

* 00 - Copy
* 01 - AND
* 02 - OR
* 03 - XOR
* 04 - INV

* 05 - Shift Left
* 06 - Shift Right Zero Extend
* 07 - Shift Right Sign Extend

#### Arithmetic operations

* 08 - Mod Unsigned
* 09 - Mod Signed

* 0A - Add Unsigned
* 0B - Add Signed
* 0C - Subtract Unsigned
* 0D - Subtract Signed

* 0E - Multiply Unsigned
* 0F - Multiply Signed
* 10 - Divide Unsigned
* 11 - Divide Signed

  Note: if carry flag is set, uses it to extend arithmetic, requires you
  to clear the carry flag when doing multiple independent math
  operations back to back.

  Note: overflow != carry, carry lets you quickly chain similar functions
  to do math on numbers with some multiple of 16 bits. Overflow
  indicates a numeric overflow if the result was the MSB.

#### Byte operations

Since the memory is word addressed, uCISC provides a number of byte specific
operations to improve performance and ease of byte based logic.

* 12 - Swap MSB and LSB bytes
* 13 - MSB to LSB: (A >> 8)
* 14 - Zero LSB:  (A & 0xFF00)
* 15 - Zero MSB:  (A & 0x00FF)
* 16 - Write MSB only: (A & 0xFF00) | (B & 0x00FF)
* 17 - Write LSB only: (A & 0x00FF) | (B & 0xFF00)
* 18 - Write LSB as MSB: ((A & 0x00FF) << 8) | (B & 0x00FF)
* 19 - Write MSB as LSB: ((A & 0xFF00) >> 8) | (B & 0xFF00)

#### Floating point math

WARNING: I took a stab at working through how this would work.
More than likely, it is wrong and I missed something important. Floating
point math is hard.

Floating point math should use IEEE 754 format modified for 16-bit blocks.
That is, 16-bit (half precision) values follow the IEEE 754 exactly. For
32-bit (single precision) values and larger, the format has been modified
to support chaining for larger values. Notably, the upper 12 bits of the
flags register are used to control the floating point process.

PPSSSCCC CCCCFFFF

PP - The precision indicator 0 (half precision) to 3 (quadruple precision)
SS - The next calculation step 0 (just starting) to 7 (step 8 in quad prec.)
CCCCCCC - The exponent differential, if needed. Once the exponent calculation
  informs the rest of the calculation. The bit shift is stored in here.
FFFF - The normal flag bits (negative | zero | carry | overflow).
  The error flag is set by setting negative and zero to 1.

To setup a calculation, you must first set the PP field in the flags register
properly.

The actual floating point number is stored in sign, exponent, significand. The
word order is little endian. That is, each word contains the least significant
bits remaining for the precision value. For, example, with a quadruple
precision number, the first 16-bits will be the sign bit and 15 exponent bits.
The last 3 words will be the 112 significand bits, least significant word first.
This format will allow the floating point calculations to proceed incrementally
for each successive word, minimizing revisiting earlier words.

During the calculations, overflow will remain 0, carry will pass through to each
successive step. After the final calculation step, the carry will be set. If
the carry flag is set, the operation must be executed one more time to update
the exponent word. The register address offset is automatically calculated
based on the precision indicator assuming the registers are now pointing to the
last word in the result value (i.e. you should not increment the last calculation
step). Thus, the proper final instruction will specify the "on carry" effect, but
it should be safe to execute regardless.

It is not required that hardward support longer floating point calculations. It
is up to the hardware implementers on how many precision values to implement.
It is required that the hardware support all precisions up to the maximum they
do support (i.e. you can't support double precision without also supporting single
precision). It is generally encouraged, but not required that hardware at least
support the half precision operation. This means that the floating point
calculations will take (use inline rather than loops for maximum performance):

* 1 cycle for half precision
* 3 cycles for single precision
* 5 cycles for double precision
* 9 cycles for quadruple precision

Hardware may accelerate larger width floating point operations as long as a given
width executes in constant time under all circumstances and the register states
afterwards are exactly as if it had been executed in multiple steps. That is the
increment, flags and memory values must be identicle. The acceleration can only
be performed when the (TBD) control flag is turned on to allow it since the
program will need to structure the application code differently.

* 1A - Add floating point
* 1B - Subtract floating point
* 1C - Multiply floating point
* 1D - Divide floating point
* 1E - Exponent from floating point (always
* 1F - Significand from floating point

