# uCISC Instructions

* Back: [Banking Mechanics](08_Banking_Mechanics.md)
* Next: [Common Devices](10_Common_Devices.md)

uCISC only has two instructions, copy and compute instructions. Most of the arguments and
behaviors between these instructions are shared.

# uCISC Statements

Some examples:

```
$stack as 1.mem
$pc as 0.reg

# Example compute instructions
# Add value at address stack + 1 to value at address stack, push to stack
compute A.op $stack 1.imm $stack push

# Convert value at address stack to page boundary, store in pc
compute E.op $stack $pc pop

# Increment value at stack address
compute A.op 1.imm $stack

# Example copy instructions
# Push duplicate value on the top of the stack
copy $stack $stack push

# Jump return
copy $stack $pc pop

# Set register 2 to stack + 2 address
copy $stack 2.imm 2.reg
```

Compute and copy instructions pack in to similar instruction bit structures:

```
# COPY Instruction Packed Bits:
0EEDDDSS SMIIIIII

# Compute Instruction Packed Bits:
0EEDDDSS SMIIAAAA

# The M bit is interpreted as part of the immediate I unless DDD is 1, 2 or 3
# See the Increment section below for details.
```

### (S) Source Argument

The source argument indicates where the data is coming from. When using uCISC
instructions, you should think of registers like pointers. They should only contain
addresses under normal conditions.

Source arguments have an identifier (the register number) and a meta type (reg, mem
or val) indicating if it is a register, memory or immediate value. Registers can be
accessed as `mem` or `reg` values. If `mem` then the register is dereferenced as a
memory address and the source value is the value at the memory location. If `reg`
then the value of the register is used directly.

If the instruction has an immediate value (i.e. it is a copy instruction) then it
will add the signed immediate value to the register before interpreting the value.

*Note:* The address space the memory arguments refer to is modified by the
address space indicators in the flags register. See the details in the flags
register section below.

| Argument     | Bit Encoding | Source Value                                    |
|:------------ |:------------ |:----------------------------------------------- |
| 0.reg        | 000 (0)      | The address in the PC for this instruction      |
| ------------ | ------------ | ----------------------------------------------- |
| 1.mem        | 001 (1)      | The value at location (r1 + immediate)          |
| 2.mem        | 010 (2)      | The value at location (r2 + immediate)          |
| 3.mem        | 011 (3)      | The value at location (r3 + immediate)          |
| ------------ | ------------ | ----------------------------------------------- |
| 4.val        | 100 (4)      | The signed immediate value from the instruction |
| ------------ | ------------ | ----------------------------------------------- |
| 1.reg        | 101 (5)      | r1 + immediate                                  |
| 2.reg        | 110 (6)      | r2 + immediate                                  |
| 3.reg        | 111 (7)      | r3 + immediate                                  |

Notice that the least significant bits always indicate the register (1-3). This
makes these instructions fairly easy to decode in hardware.

### (D) Destination

The destination argument indicates where the data being sent. When using uCISC
instructions, you should think of registers like pointers. They should only contain
addresses under normal conditions.

Destination arguments have an identifier (the register number) and a meta type (reg,
mem or val) indicating if it is a register, memory or immediate value. Registers can
be accessed as `mem` or `reg` values. If `mem` then the register is dereferenced as a
memory address and the source value is the value at the memory location. If `reg`
then the value of the register is used directly.

Destination arguments are *never* modified by immediate values.

*Note:* The address space the memory arguments refer to is modified by the
address space indicators in the flags register. See the details in the flags
register section below.

| Argument     | Bit Encoding | Target Location       |
|:------------ |:------------ |:--------------------- |
| 0.reg        | 000 (0)      | The program counter   |
| ------------ | ------------ | --------------------- |
| 1.mem        | 001 (1)      | Address location r1   |
| 2.mem        | 010 (2)      | Address location r2   |
| 3.mem        | 011 (3)      | Address location r3   |
| ------------ | ------------ | --------------------- |
| 4.reg        | 100 (4)      | The control register  |
| ------------ | ------------ | --------------------- |
| 1.reg        | 101 (5)      | r1                    |
| 2.reg        | 110 (6)      | r2                    |
| 3.reg        | 111 (7)      | r3                    |

Notice that the least significant bits always indicate the register (1-3). This
makes these instructions fairly easy to decode in hardware.

### (E) Effect

The effect argument controls whether the operation is completed or not. The processor
will look at the result flags from the previous compute operation and decide whether
to store the result or not depending on the effect argument. Only compute
instructions set the result flags.

* 0.eff - store if zero
* 1.eff - store if not zero
* 2.eff - store if negative
* 3.eff - store

Branches/jumps are implemented by setting the destination of the instruction as the
PC with the jump address as the source. The effect controls whether the branch is
taken or not.

*Note:* The flags register is not directly accessible. You can figure out what is in
the register via the use of effect modifiers on instructions. If you really need the
bits directly, you can always use the processor control segment to read the register
value directly.

### (I) Immediate

How the immediate is interpreted is depending on the source and destination types.
This is done to support as many useful modes of refering to values as possible at the
cost of adding complexity to the decoding hardware.

Conditions:

* If neither source or destination is a "mem" reference, the increment bit is treated
  as part of the immediate.

In this case, the increment bit would have no meaning. Using it as part of the
immediate means we get a signed 7-bits for copy instructions. This means we get a
value between -64 and 63 inclusive for copy instructions (-4 to 3 for compute). The
primary use cases for this are jumps within a page or loading immediate values into
registers. Note, pages are 64 words precisely because this leaves the entire page
accessible via jumps from anywhere else in the page (except jumping from the very
first instruction in a page to the very last).

* If the source is a "mem" argument, the immediate is treated as unsigned. It is
  treated as signed otherwise.

When the source is a "mem" argument, you are almost certainly wanting to address
some value at an offset from an object pointer, stack reference or array. Negative
values rarely make sense in this case and using unsigned values here gives a larger
offset range. For copy instructions, we get an offset between 0 and 63, for compute
instructions we get 0 to 3. This let's us address values close to the top of the
stack and avoid fancy register work unless we are computing values further up the
stack.

* Special case: For copy instructions only, when both the source and destination
  arguments are "mem" references, the immediate is split between the two arguments.

When writing real code, I discovered that it is incredible useful to be able to copy
values up the stack. This is often done when copying return values to the right place
on the stack, for example. Splitting the 6-bit immediate into two 3-bit immediates
is very useful in practice for using fewer instructions and writing code that is
easier to understand. In this case, both the source and destination can reference
values offset by 0 to 7 from the register address.

### (M) Incre(M)ent

* 0.inc - No increment
* 1.inc - Increment register argument as described below

The effect of the increment is as follows:

* Push: If the destination is a "mem" argument, decrement the destination *before*
  storing the result, but after reading the compute argument.
* Pop: If the source is a "mem" argument, increment the source *after* copying the
  value.
* If *both* are mem arguments, only the destination is affected.
* If *neither* are mem arguments, the increment is actually part of the immediate

Typically used to "push" variables to the stack by incrementing the address before
storing the result. However, you can also use it to increment through arrays and
other uses where you want to move a memory reference in word increments.

*Note:* The increment is only persisted if the effect causes a store operation.

*Important:* For compute instructions, the destination argument is not incremented before
being passed into the ALU. So, A + B => B can be turned into A + B => C with an
increment push specified.

### Examples

As noted above, the interpretation of the M bit depends on other values:

```
# Example: COPY Instructions
# When a mem arg is used, immediate is 6-bits
0EEDDDSS SMIIIIII 

# When no mem arg is used, immediate is 7-bits
0EEDDDSS SIIIIIII 

# Example: Compute Instructions
# When a mem arg is used, immediate is 2-bits
1EEDDDSS SMIIAAAA 

# When no mem arg is used, immediate is 3-bits
1EEDDDSS SIIIAAAA 
```

### Flags Register

The flags register contains the result of compute operations and other special cases.
Note that the flags register is not directly readable by instructions. You can use
the effect flags to infer it or use the banked control slot for the current processor
to manipulate it.

```
# Flag Register Bits
H000000S 0000 OCNZ
```

* (H) Halt. This flag is set when the processor is halted.
* (S) Signed. Indicates the signedness of compute operations (1 is signed)
* (Z) Zero flag is set if the last compute resulted in a zero
* (N) Negative flag is set if the last compute resulted in a negative value
* (C) Carry flag is set if the last compute resulted in a carry
* (O) Overflow flag is set if the last compute resulted in an overflow

### Control Register

The control register controls memory banking.
See [Banking Mechancs](08_Banking_Mechanics.md) for more information.

### Special Instruction Cases

#### (H)alt Instruction

If you copy 0.reg (the program counter) to itself with an immediate offset of 0
this would normally put the processor into an infinite loop, since the next
instruction executed will be the currently executing instruction. uCISC processors
instead interpret this as a HALT instruction. HALT instructions are conditional on
the specified effect and only trigger if the effect condition matches.

```
# Copy HALT Instruction Bits
0EE00000 00000000
```

You can also trigger a halt by setting the correct bit in the flags register.

On halt, processors execution stops. To restart a halted processor, init processor
must manipulate the device using a control segment to initialize it.

#### (A) Compute Opcode

The compute opcode controls the arithmetic or logic operation invoked. Math functions
that take a single arg (e.g. invert, negate, byte operations, etc.) use the source
as the argument and store the result in the target. Math functions that require two
arguments, use both the source and target as the arguments and store the result in
the target.

*Important:* All 2-arg compute operations use the target as the first arg and the
source as the second. For example, subtraction is D - S => D (subtract S from D).
Shift uses D >> S => D (shift D by S bits) and store in D. The value in D is
destroyed in the process (unless a push is used).

For example:

```
# source is 2.mem
# This means 1.mem - 2.mem => 1.mem
compute B.op 2.mem $stack
```

This allows you to repeatedly modify a number in some way that. For example if you
want to shift a number a bit at a time and take action on the bit, you can do so by
running the same shift multiple times. You can also use this to subtract a constant
value from a larger number multiple times. You can run the example instruction above
repeatedly to subtract the 2.mem value each time.

*Signedness:* The processor can be put into signed or unsigned mode. Arithmetic and
right shift operations respect this mode.

#### Bit operations

Overflow register is set to 0, carry and overflow flags are 0

* 00 - INV
* 01 - AND
* 02 - OR
* 03 - XOR
* 04 - Negate (2's compliment)

Shift operations put the bits shifted out in the least significant position of the
overflow register. Oveflow and carry flags are set to 1 if any non-zero bit is
shifted out, 0 otherwise.

* 05 - Shift left, zero extend
* 06 - Shift right, respect signed mode

#### Byte operations

Overflow register is set to 0, carry and overflow flags are 0

* 07 - Swap MSB and LSB bytes
* 08 - Zero LSB:  (A & 0xFF00)
* 09 - Zero MSB:  (A & 0x00FF)

#### Arithmetic operations

Arithmetic operations always set the overflow register, resulting in a 32-bit number
if you consider the overflow. Carry and overflow flags are set appropriately if there
is a carry or overflow of the 16-bit width.

* 0A - Add, respect signed mode
* 0B - Subtract, respect signed mode
* 0C - Multiply, respect signed mode
* 0D - Divide, respect signed mode, overflow is remainder

#### Extra operations

Overflow register is set to 0, carry and overflow flags are 0

* 0E - Extract page boundary from address (a page is 64 words)
* 0F - Add source + overflow register => destination

#### Continue Reading

* Back: [Banking Mechanics](08_Banking_Mechanics.md)
* Next: [Common Devices](10_Common_Devices.md)

