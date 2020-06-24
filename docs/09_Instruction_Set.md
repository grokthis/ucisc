# uCISC Instructions

uCISC only has two instructions, copy and ALU instructions. Most of the arguments and
behaviors between these instructions are shared.

# uCISC Statements

```
$stack as 1.mem
$pc as 0.reg

# Example ALU instructions
# Add value at address stack + 1 to value at address stack, push to stack
compute A.op push $stack:1 $stack

# Convert value at address stack to page boundary, store in pc
compute E.op pop $stack $pc

# Increment value at stack address
compute A.op 1.imm $stack

# Example copy instructions
# Push duplicate value on the top of the stack
copy push $stack $stack

# Jump return
copy pop $stack $pc

# Set register 2 to stack + 2 address
copy $stack:2 2.reg
```

These instructions pack in to similar instruction bits:

```
# COPY Instruction Packed Bits:
0EEDDDSS SMIIIIII

# ALU Instruction Packed Bits:
0EEDDDSS SMIIAAAA

# The MG bits are interpreted as part of the immediate I unless DDD is 1, 2 or 3
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
will look at the result flags from the previous ALU operation and decide whether to
store the result or not depending on the effect argument. Only ALU instructions set
the result flags.

* 0.eff - store if zero
* 1.eff - store if not zero
* 2.eff - store if negative
* 3.eff - store

Branches/jumps are implemented by setting the destination of the instruction as the
PC with the jump address as the source. The effect controls whether the branch is
taken or not.

*Note:* The flags register is not directly accessible. You can only figure out what
is in the register via the use of effect modifiers on instructions.

### (I) Immediate

The immediate bit width depends on the instruction and increment conditions. See the
increment and sign section for more details. The bit width and whether it is signed
determine the range of acceptible values. Normally when coding, the compiler will
verify that the requested immediate fits in the given space when combined with the
other parameters of the instruction.

Immediates are interpreted as unsigned if the source is a "mem" argument, but are
always interpreted as signed otherwise.

### (M) Incre(M)ent

* 0.inc - No increment
* 1.inc - Increment register arguments before storing the result

The effect of the increment is as follows:

* Push: If the destination is a "mem" argument, decrement the destination
* Pop: If the source is a "mem" argument, increment the source
* If *both* are mem arguments, only the destination is affected

Typically used to "push" variables to the stack by incrementing the address before
storing the result. However, you can also use it to increment through arrays and
other uses where you want to move a memory reference in small increments.

If neither the source or the destination are "mem" arguments specifying an increment
is illegal. In that case, the M bit is treated as part of the immediate.

*Note:* The increment is only persisted if the effect causes a store operation.

*Important:* For ALU instructions, the destination argument is not incremented before
being passed into the ALU. So, A + B => B can be turned into A + B => C with an
increment push specified.

As noted above, the interpretation of the M bit depends on other values:

```
# Example: COPY Instructions
# When a mem arg is used, immediate is 6-bits
0EEDDDSS SMIIIIII 

# When no mem arg is used, immediate is 7-bits
0EEDDDSS SIIIIIII 

# Example: ALU Instructions
# When a mem arg is used, immediate is 2-bits
0EEDDDSS SMIIAAAA 

# When no mem arg is used, immediate is 3-bits
0EEDDDSS SIIIAAAA 
```

### Control Register

```
# Control bits:
000000HG WWWWWWFT
```

The control register specifies various aspects of instruction behavior.

#### Si(G)ned Mode

ALU arithmetic instructions can be executed in signed or unsigned mode.

#### Address Space Flags - (F)rom and (T)o

The uCISC processor has 2 address spaces. The default address space is the local
memory address space. The other is the banked memory space. Processors can use banked
memory to access other hardware devices in the system. The flags in the control
register determine which space is read or written to. You can control these spaces
independently to read and write between local memory and banked memory.

Values:

* 0 = local memory space
* 1 = banked memory space

#### Copy (W)idth

*Warning:* This mechanism will most likely be removed or modified dramatically.

Most uCISC implementations will be able to perform read/write operations in bulk,
particularly for transactions to/from banked memory. Without higher throughput,
data transfers will be painfully slow. In order to allow better transfer, code can
set the transfer page control.

When zero, all instructions will execute normally, affecting a single source and
destination address. If the transfer width is set, COPY instructions will be affected
by the width value. W + 1 words will be transferred as if the COPY instruction had
been executed W + 1 times. At most 64 words can be copied in a single instruction.
The copy will be at least as fast as executing the copy instruction W + 1 times, but
may be significantly faster.

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

You can also trigger a halt by setting the correct bit in the control register.

On halt, processors execution stops. To restart a halted processor, another processor
or controller must manipulate the device using a control segment to initialize it.

#### (A) ALU opcode

The ALU opcode controls the arithmetic or logic operation invoked. Math functions
that take a single arg (e.g. invert, negate, byte operations, etc.) use the source
as the argument and store the result in the target. Math functions that require two
arguments, use both the source and target as the arguments and store the result in
the target.

*Note:* All non-associative, 2-arg operations use the target as the first arg and the
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
