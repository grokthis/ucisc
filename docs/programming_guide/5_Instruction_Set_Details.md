## uCISC Programming Guide

1. [Getting Started](1.0_Getting_Started.md)
   1. [Configuring IntelliJ](1.1_Configuring_IntelliJ.md)
   2. [Configuring VIM](1.2_Configuring_VIM.md)
   3. [Compiling uCISC Code](1.3_Compiling_uCISC_Code.md)
   4. [Using the uCISC Simulator](1.4_Simulating_uCISC.md)
   5. [Using the uCISC Soft Core](1.5_Running_uCISC_Soft_Core.md)
2. [Introduction to Programming With uCISC](2.0_Program_With_uCISC.md)
   1. [Accessing External Devices](2.1_Accessing_Devices.md)
   2. [Common Devices](2.2.0_Common_Devices.md)
      1. [GPIO](2.2.1_GPIO_Devices.md)
      2. [I2C](2.2.2_I2C_Devices.md)
      3. [UART](2.2.3_UART_Devices.md)
      4. [Video Devices](2.2.4_Video_Devices.md)
   3. [Advanced uCISC Programming Techniques](2.3_Advanced_Programming_Techniques.md)
3. [uCISC Syntax Quick Reference](3_Syntax_Quick_Reference.md)
4. [Standard Libraries](4_Standard_Libraries.md)
5. [Instruction Set Details](5_Instruction_Set_Details.md) <-- you are here

** Warning ** Some parts of this doc may be slightly ot of date. Nothing too serious,
but some names/options/syntax may be slightly off. Where it conflicts with the other
docs, the other docs are right.

# uCISC Instructions

uCISC only has two instruction bit structures, they are as follows:

```
# Packed bits for non-memory destinations
SSSSDDDD MEEEAAAA IIIIIIII IIIIIIII

# Packed bits for memory destinations
SSSSDDDD MEEEAAAA OOOOIIII IIIIIIII
```

# uCISC Statements

Some examples:

```
# This is a comment

# Define a bunch of name substitutions for syntax sugar
def stack as 1.mem
def pc as 0.reg
def val as 4.val
def copy as 0.op
def add as 10.op
def r2 as 2.reg

# Example compute instructions
# Add value at address stack + 1 to value at address stack, push to stack
add stack 1.imm stack push

# Increment value at stack address
add val 1.imm stack

# Example copy instructions
# Push duplicate value on the top of the stack
copy stack stack push

# Jump return
# Prefix reg names by `&` to refer to the reg content rather than the memory content
copy stack &pc pop

# Set register 2 to stack address + 2
copy &stack 2.imm &r2
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

If the source argument is a `mem` value, then the immediate is interpreted as
unsigned, otherwise it is treated as signed and added to the source value.

*Note:* The address space referenced by a register depends on the memory banking
setup. See the chapter on Banking Mechanics for details.

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
| ------------ | ------------ | ----------------------------------------------- |
| 8.reg        | 1000 (8)     | Flags Register                                  |
| ------------ | ------------ | ----------------------------------------------- |
| 5.mem        | 1001 (9)    | The value at location (r5 + immediate)          |
| 6.mem        | 1010 (10)    | The value at location (r6 + immediate)          |
| 7.mem        | 1011 (11)    | The value at location (r7 + immediate)          |
| ------------ | ------------ | ----------------------------------------------- |
| 12.reg       | 1100 (12)    | Interrupt Register                              |
| ------------ | ------------ | ----------------------------------------------- |
| 5.reg        | 1101 (13)    | r5 + immediate                                  |
| 6.reg        | 1110 (14)    | r6 + immediate                                  |
| 7.reg        | 1111 (15)    | r7 + immediate                                  |

Note that the pattern here is intended to make it easy for hardware decoding. That's
the reason for the numbering skips.

### (D) Destination

The destination argument indicates where the data being sent. When using uCISC
instructions, you should think of registers like pointers. They should only contain
addresses under normal conditions.

Destination arguments have an identifier (the register number) and a meta type (reg,
mem or val) indicating if it is a register, memory or immediate value. Registers can
be accessed as `mem` or `reg` values. If `mem` then the register is dereferenced as a
memory address and the source value is the value at the memory location. If `reg`
then the value of the register is used directly.

If the destination is a `mem` type, the offset adds to the value before determining
the final memory lookup address.

*Note:* The address space referenced by a register depends on the memory banking
setup. See the chapter on Banking Mechanics for details.

| Argument     | Bit Encoding | Target Location                                 |
|:------------ |:------------ |:----------------------------------------------- |
| 0.reg        | 000 (0)      | The program counter                             |
| ------------ | ------------ | ----------------------------------------------- |
| 1.mem        | 001 (1)      | Address location r1                             |
| 2.mem        | 010 (2)      | Address location r2                             |
| 3.mem        | 011 (3)      | Address location r3                             |
| ------------ | ------------ | ----------------------------------------------- |
| 4.reg        | 100 (4)      | Banking Control Register                        |
| ------------ | ------------ | ----------------------------------------------- |
| 1.reg        | 101 (5)      | r1                                              |
| 2.reg        | 110 (6)      | r2                                              |
| 3.reg        | 111 (7)      | r3                                              |
| ------------ | ------------ | ----------------------------------------------- |
| 8.reg        | 1000 (8)     | Flags Register                                  |
| ------------ | ------------ | ----------------------------------------------- |
| 5.mem        | 1001 (10)    | The banked value at location (r5 + immediate)   |
| 6.mem        | 1010 (10)    | The banked value at location (r6 + immediate)   |
| 7.mem        | 1011 (10)    | The banked value at location (r7 + immediate)   |
| ------------ | ------------ | ----------------------------------------------- |
| 12.reg       | 1100 (12)    | Interrupt Register                              |
| ------------ | ------------ | ----------------------------------------------- |
| 5.reg        | 1101 (13)    | r5 + immediate                                  |
| 6.reg        | 1110 (14)    | r6 + immediate                                  |
| 7.reg        | 1111 (15)    | r7 + immediate                                  |

Note that the pattern here is intended to make it easy for hardware decoding. That's
the reason for the numbering skips.

### (E) Effect

The effect argument controls whether the operation is completed or not. The processor
will look at the result flags from the previous non-copy operation and decide whether
to store the result or not depending on the effect argument. Only non-copy instructions
set the result flags.

* 0.eff - store if zero
* 1.eff - store if !zero
* 2.eff - store if negative
* 3.eff - store if positive
* 4.eff - store
* 5.eff - store if overflow
* 6.eff - store if interrupt
* 7.eff - set flags only

Branches/jumps are implemented by setting the destination of the instruction as the
PC with the jump address as the source. The effect controls whether the branch is
taken or not.

### (I) Immediate

How the immediate is interpreted depends on the destination type. If the destination
is a memory argument (values 1-3 and 9-11) the immediate is only 12 bits wide to make
room for the offset. When signed, this allows -2048 to 2047. When unsigned, 0 to 4096.
If the destination is not a memory element, it is a full 16-bit value and can
represent values from -32768 to 65535. Since it's a full 16-bits the sign is
irrelevant to the hardware since it doesn't need to do sign extending. In other words,
adding 65535 is the same as adding -1 to the source because of the magic of 2's
complement.

### (M) Incre(M)ent

* 0.inc - No increment
* 1.inc - Increment register argument as described below

The effect of the increment is as follows:

* Push: If the destination is a `mem` argument, decrement the destination **before**
  storing the result, but after reading the value from the current address if needed
  for the op code.
* Pop: If **only** the source is a "mem" argument, increment the source *after
* copying the
  value.
* If **both** are mem arguments, only the destination is affected.
* If *neither* are mem arguments, the increment must be zero.

Typically, it is used to "push" variables to the stack by incrementing the address
before storing the result. However, you can also use it to increment through arrays
perform more efficient copy operations and more.

*Note:* The incremented value only saves to the register if the effect causes a store
operation.

*Important:* For compute instructions, the destination argument does not increment
before being passed into the ALU. So, A + B => B can be turned into A + B => C with an
increment push specified.

### Flags Register

The flags register contains the result of compute operations and other special cases.

```
# Flag Register Bits
HRB0 000S 000I OCNZ
```
You can control the operation of the processor by setting these flags:
* (H) Halt. Halts all execution.
* (R) Resume execution when interrupt occurs. Processes can set this flag and halt.
  An interrupt will cause the processor to resume execution at the interrupt address.
* (B) Blocked memory is writable while executing. In special cases, you may want to
  allow the init device to write data to memory while the processor is executing.
  Setting this flag turns that on.
* (S) Signed. Indicates the signedness of compute operations (1 is signed)

These flags indicate the result of compute operations:
* (I) Interrupt flag
* (Z) Zero flag indicates if the last compute resulted in a zero
* (N) Negative flag indicates if the last compute resulted in a negative value
* (C) Carry flag indicates if the last compute resulted in a carry
* (O) Overflow flag indicates if the last compute resulted in an overflow
  
### Control Register

The control register controls memory banking.
See [Banking Mechanics](2.1_Accessing_Devices.md) for more information.

### Special Instruction Cases

#### (H)alt Instruction

If you copy 0.reg (the program counter) to itself with an immediate offset of 0
this would normally put the processor into an infinite loop, since the next
instruction executed will be the currently executing instruction. uCISC processors
instead interpret this as a HALT instruction. HALT instructions are conditional on
the specified effect and only trigger if the effect condition matches.

```
# Copy HALT Instruction Bits
00000000 0EEE0000 00000000 00000000
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

Carry and overflow flags are 0

Format: CODE - Name and description (instruction)

* 00 - Copy (copy)
* 01 - And (and)
* 02 - Or (or)
* 03 - Xor (xor)
* 04 - Invert (inv)

* 05 - Shift left, zero extend (shl)
* 06 - Shift right, respect signed mode (shr)

#### Byte operations

Carry and overflow flags are 0

* 07 - Swap MSB and LSB bytes (swap)
* 08 - MSB only: A & 0xFF00 (msb)
* 09 - LSB only: A & 0x00FF (lsb)

#### Arithmetic operations

Carry and overflow flags are set appropriately. Overflow means the
result is too big to correctly represent in the result.

* 0A - Add, respect signed mode (add)
* 0B - Subtract, respect signed mode (sub)
* 0C - Multiply, respect signed mode, carry is zero (mult)
* 0D - Multiply, return MSW, respect signed mode, carry is zero (mult)
* 0E - Add with carry in, respect signed mode (addc)
* 0F - TBD

#### Continue Reading

* Back: [Banking Mechanics](2.1_Accessing_Devices.md)
* Next: [Common Devices](2.2.0_Common_Devices.md)

