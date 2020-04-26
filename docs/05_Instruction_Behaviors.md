# uCISC Common Instruction Behaviors

uCISC only has two instructions, Copy and ALU instructions. Most of the
arguments and behaviors between these instructions are shared.

### Common Arguments

#### (S) Source Argument

The source argument indicates where the data is coming from. When using
uCISC instructions, you should think of registers like pointers. They
should only contain addresses under normal conditions.

Source arguments have an identifier (the register number) and a meta type
(reg, mem or val) indicating if it is a register, memory or immediate value.
Registers can be accessed as `mem` or `reg` values. If `mem` then the register
is dereferenced as a memory address and the source value is the value at the
memory location. If `reg` then the value of the register is used directly.

If the instruction has an immediate value (i.e. it is a copy instruction) then
it will add the signed immediate value to the register before interpreting the
value.

*Note:* Copy and ALU instructions treat argument "4" differently as noted below.
Copy instructions treat the source as `4.val` and use the immediate value in
the instruction as the source. ALU instructions use `4.reg` instead and use the
flags register as the source (ALU instructions don't have an immediate).

*Note:* The address space the memory arguments refer to is modified by the
address space indicators in the flags register. See the details in the flags
register section below.

*Note:* Operations may perform the specified operation multiple times depending
on the word width specified in the flags register. Each additional operation will
be repeated at subsequent addresses. See the details in the flags register section.

| Argument     | Bit Encoding | Source Value                                    |
|:------------ |:------------ |:----------------------------------------------- |
| 0.reg        | 000 (0)      | The address in the PC for this instruction      |
| ------------ | ------------ | ----------------------------------------------- |
| 1.mem        | 001 (1)      | The value at location (r1 + immediate)          |
| 2.mem        | 010 (2)      | The value at location (r2 + immediate)          |
| 3.mem        | 011 (3)      | The value at location (r3 + immediate)          |
| ------------ | ------------ | ----------------------------------------------- |
| 4.val (Copy) | 100 (4)      | The signed immediate value from the instruction |
| 4.reg (ALU)  | 100 (4)      | The current value in the carry register         |
| ------------ | ------------ | ----------------------------------------------- |
| 1.reg        | 101 (5)      | r1 + immediate                                  |
| 2.reg        | 110 (6)      | r2 + immediate                                  |
| 3.reg        | 111 (7)      | r3 + immediate                                  |

Notice that the least significant bits always indicate the register (1-3). This
makes these instructions fairly easy to decode in hardware.

#### (D) Destination

The destination argument indicates where the data being sent. When using
uCISC instructions, you should think of registers like pointers. They
should only contain addresses under normal conditions.

Destination arguments have an identifier (the register number) and a meta type
(reg, mem or val) indicating if it is a register, memory or immediate value.
Registers can be accessed as `mem` or `reg` values. If `mem` then the register
is dereferenced as a memory address and the source value is the value at the
memory location. If `reg` then the value of the register is used directly.

Destination arguments are *never* modified by immediate values.

*Note:* The address space the memory arguments refer to is modified by the
address space indicators in the flags register. See the details in the flags
register section below.

*Note:* Operations may perform the specified operation multiple times depending
on the word width specified in the flags register. Each additional operation will
be repeated at subsequent addresses. See the details in the flags register section.

| Argument     | Bit Encoding | Target Location       |
|:------------ |:------------ |:--------------------- |
| 0.reg        | 000 (0)      | The program counter   |
| ------------ | ------------ | --------------------- |
| 1.mem        | 001 (1)      | Address location r1   |
| 2.mem        | 010 (2)      | Address location r2   |
| 3.mem        | 011 (3)      | Address location r3   |
| ------------ | ------------ | --------------------- |
| 4.reg        | 100 (4)      | The flags register    |
| ------------ | ------------ | --------------------- |
| 1.reg        | 101 (5)      | r1                    |
| 2.reg        | 110 (6)      | r2                    |
| 3.reg        | 111 (7)      | r3                    |

Notice that the least significant bits always indicate the register (1-3). This
makes these instructions fairly easy to decode in hardware.

#### (E) Effects

The effect argument controls whether the operation is completed or not. The
processor will look at the flags in the flags register and decide whether to
store the result or not depending on the effect argument.

Effects 0-3 are always available. Effects 4-7 are available on ALU instructions
under certain conditions. See the ALU instruction for more details.

Normal effects:

* 0.eff - store if zero
* 1.eff - store if not zero
* 2.eff - store if positive
* 3.eff - store

Extended effects:

* 4.eff - store if not negative
* 5.eff - store if negative
* 6.eff - store if not overflow
* 7.eff - do not store

Branches/jumps are implemented by setting the destination of the instruction as the
PC with the jump address as the source. The effect controls whether the branch is
taken or not.

#### Flags Register

The flags register controls the access behavior of instructions. The upper
byte of the flags word contains the ALU output flags. The lower byte is used
to control the address space of the source and destination bits. The uCISC
processor has 2 address spaces. The default address space is the local memory
address space. The other is the device space. The device space is used to
control processor function, read processor metadata, the hardware clock and
any devices attached to the processor (shared or directly attached).

The device space may have access to hardware accelerator functions, communication
buses, peripherals and more.

```
# Flags bits:
NZCOSFTR RRRRRRRR
```
ALU Result Flags

* N = negative result flag
* Z = zero result flag
* C = carry flag
* O = overflow flag

*Address Space Flags*

Values:

* 0 = local mem space
* 1 = device space

If the flag is set to 1, the address space refers to device space. If the
device space address does not exist, it falls through to memory space. If the
flag is set to 0, or the address falls through, then the address refers to
memory space, respecting banked memory configuration. If there is not memory
at the address, it falls through to the void. The void always reads 0 and
ignores writes. See the memory system description for more details.

* F = From address space; controls the address space of RRR
* T = To address space; controls the address space of DDD
* R = Repetition factor
* S = Sticky value of transfer width. If 0, the transfer width is reset to 1
  after the next instruction. If 1, the width is retained.

The repetition factor R is interpreted as an 9-bit unsigned integer that
indicates the number of times an operation will be repeated. When set, the
instruction will behave exactly as if the processor had executed the exact same
instruction R times. By default, the value is 1, meaning the instruction happens
exactly once.

Some examples:

```
# Setup: Set the transfer width to 8
0/copy/ 4.val 8.imm 4.reg

# Given that setup, the following instructions behave as follows:

# Copy zero for 8 words, 1.mem is unchanged, so this instruction writes 0.imm
# to the same address 8 times. Could be useful when writing to certain I/O
# devices, or could be used as a fixed cycle delay.
0/copy/ 4.val 0.imm 1.mem

# Copy zero for 8 words, 1.mem is incremented each time because of the `1.push`
# argument. Could be useful when initializing blocks of memory or arrays.
0/copy/ 4.val 0.imm 1.mem 1.push

# Mem replicate for 8 words from one source address, 1.mem is incremented each time
# because of the `1.push` argument. Effectively replicates a value across the
# destionation memory. Could be useful when initializing blocks of memory or arrays
# to non-zero values.
0/copy/ 2.mem 0.imm 1.mem 1.push

# Memcopy for 8 words, 1.mem and 2.mem are both incremented after each copy because
# the increment flag is set. Effectively replicates a block of memory from one
# location to another (a mem copy).
200/alu copy/ 2.mem 0.imm 1.mem 1.inc

# Add two values for 8 words, 1.mem and 2.mem are both incremented after each add
# operation because the increment flag is set. The carry flag is set and carried
# between each add. You can essentially use this to add arbitrarily large numbers.
20A/alu add/ 2.mem 0.imm 1.mem 1.inc
```

The main purpose of the transfer width is performance optimization of data transfers
between memory locations and processors (see below). However, you can get other
interesting effects by creative use of the memory pointers. For example, it's possible
to do n-factorial using n-bytes of memory in two passes:

```
# Setup: N is on the stack (1.mem), load into r2
0/load N into r2/ 1.mem 2.reg
# create stack space for math
20C/alu sub N to stack/ 2.reg 1.reg

# Init width with N
217/alu write MSB/ 2.reg 4.reg
# Create
20F/alu decrement N times on stack/ 1.mem 1.mem 1.inc 0.sign

# reset stack pointer to beginning of the array of numbers
20C/alu sub N to stack/ 2.reg 1.reg
# Reset instruction width to N (we didn't set the sticky bit)
217/alu write MSB/ 2.reg 4.reg
# Multiple N(N-1) N times
208/alu mult N times/ 1.mem 1.mem 1.inc 0.sign

# N! is now next location after 1.meme. Copy N! over N on stack
0/copy/ 1.mem 1.imm 1.mem
```

Notice this is a bit more compact that the more standard implementation of factorial
in the examples (8 instructions instead of 13). Further, the example algorithm executes
~200 instructions to compute the answer to `factorial(16)`. This one will execute 38,
between 5 and 6 times faster and uses less memory than the recursive algorithm. There
may be additional performance benefits achieved by using this super-scaler technique.
See the performance discussion below.

#### Performance

uCISC processors are designed to provide predictable performance as much as possible
while allowing. Any processor using uCISC hardware must execute instructions in constant
time (1 instruction per N clock cycles, called a tick) regardless of any effects, jumps
or other conditions. The one exception is intentionally blocking operations when
reading/writing from devices or other processors.

##### Debug Mode

If the repetition value is set to 0, the instruction will behave exactly as if the
instruction had not been executed for a single tick, including any and all side effects
such as setting the flags register, except that:

* The instruction takes one tick to execute
* The repetition sticky flag is respected

If the sticky flag is 0, execution will continue at this instruction after one tick. If
the sticky flag is 1, the processor will pause indefinitely and enable the debug interface.
If there is no debug interface present (or JTAG based equivalent), the processor will
freeze indefinitely until reset by it's parent. The mechanics of the debug interface are
not specified.

##### Super Scaler Performance

if the repitition value is greater than 1, super scaler mode is enabled. If performance
scaling is supported by the hardware, it is allowed to execute the repeated instruction
as fast as possible as long as the result is equivalent (other than clock cycles) to
repeating the instruction the number of times the repetition value indicates.

Some uses:

* Faster data transfer to parent processors. Using repetition, hardware may be able to
  copy larger data widths in a single clock cycle.
* Super scaler ALU operations to perform array math more effectively.
* Super scaler device interactions. If the device is an FPU interface, for example, you
  may be able to achieve high bandwidth floating point math using repetition.
* Manipulating higher bit width numbers. If you want to manipulate 32-bit or 64-bit numbers
  you may be able to do that by setting a higher repetition width (2 for 32-bit, 4 for 64-bit)
  For example, multiplying larger width numbers can use the
  [Karatsuba algorithm](https://en.wikipedia.org/wiki/Karatsuba_algorithm) to turn 32-bit
  multiplication into 16-bit multiplication and addition.

The exact acceleration behavior is hardware dependent and you'll need to understand the
specifics of the hardware you are using. Note that if the hardware doesn't support super
scaler operations, it will always fall back to brute force repetition.

#### Special Case: Halt Instruction

If you copy 0.reg (the program counter) to itself with an immediate
offset of 0 this would normally put the processor into an infinite loop,
since the next instruction executed will be the currently executing
instruction.

uCISC processors instead interpret this as a HALT instruction. HALT
instructions are conditional on the specified effect and only happen
if the effect condition matches.

```
# Copy HALT Instruction Bits
0EE00000 00000000

# ALU HALT Instruction Bits, note that sign G modifies which effect is used
1EE00000 000000G0
```

On halt, processors will reset themselves. The following actions are
taken:

* All execution is halted.
* The comm memory bank for the first 4k block is mapped to the parent
  (at ID 0) and enabled.
* All other memory banks are disabled (but the mapping is preserved).
* If the parent attempts to read or write memory, the processor will
  read or write the indicated address.

In this state, the parent process can read the contents of the first
block, give it new code to execute and start execution.
