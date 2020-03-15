## Design Goals

With a few exceptions, these design goals have mostly been around since the
beginning of my hobby to build a processor from scratch. However, as the endeavor
became more serious and I learned more about retro computing, computer
architectures and the current state of hobby electronics, I have refined them a
great deal.

The examples below are given in uCISC syntax since it is more readable than
the binary instruction format. See the [syntax chapter](/03_Syntax.md) for more
details. If you want to read through here first, just not that the examples
make verbose use of inline comments which are between slashes:

```
# This is a line comment
... /this is an inline comment/ ...
# Comments are removed before compliation and have no bearing on the result
```

#### Maximum Micro

The instruction set should be as small as possible. Initially, the aim was
an 8-bit architecture, but in the end, that was too limiting. I tried really
hard for a long time to stay at 8-bits, but a number of issues cropped up.

* The space for immediate values was too small, 16-bit uCISC only has space
  for 6 signed bits as it is, so you can imagine trying to fit it in 8-bits.
* I didn't want to compromise on constant space and time (see below) and allow
  variable byte width instructions.
* The 16-bit address space is the minimum I could imagine to be useful for even
  a hobby computer (256 addresses in an 8-bit space seemed small to me). In the
  end 16-bit address space and 16-bit instructions has huge hardware benefits
  since addresses and instructions can fit in the same registers, bus width, etc.

All instructions, therefore, should be 16-bits wide including the immediate
value. Variable instruction widths are not allowed.

#### Memory Based System

uCISC comes firmly down on the side of CISC and is opinionated that memory is
the primary thing being manipulated. Registers are used only because it makes
the instructions easier to implement in hardware in many respects. RISC
architectures aim to "reduce" the amount of work a processor does on each
instructions. The
[Wikipedia article on RISC](https://en.wikipedia.org/wiki/Reduced_instruction_set_computer)
puts it like this:

> A common misunderstanding of the phrase "reduced instruction set computer" is
> the mistaken idea that instructions are simply eliminated, resulting in a smaller
> set of instructions. In fact, over the years, RISC instruction sets have grown in
> size, and today many of them have a larger set of instructions than many CISC CPUs.
> Some RISC processors such as the PowerPC have instruction sets as large as the CISC
> IBM System/370, for example; conversely, the DEC PDP-8—clearly a CISC CPU because
> many of its instructions involve multiple memory accesses—has only 8 basic
> instructions and a few extended instructions.
>
> The term "reduced" in that phrase was intended to describe the fact that the amount
> of work any single instruction accomplishes is reduced—at most a single data memory
> cycle—compared to the "complex instructions" of CISC CPUs that may require dozens of
> data memory cycles in order to execute a single instruction. In particular, RISC
> processors typically have separate instructions for I/O and data processing.
>
> The term load/store architecture is sometimes preferred.

In other words, RISC based systems copy values to registers and manipulate them there,
whereas CISC usually more directly manipulates memory instead. uCISC typically does
up to four memory accesses per instruction:

* Read the instruction from the PC
* Read the value at the source memory location
* Read the value at the target memory location
* Write the result value to the target memory location

Processors that implement the uCISC instruction set are required to do all of these in
constant time execution (see below) with only a few exceptions. The example
implementations do all of this in a single clock cycle. For details on how this works
read the memory design documentation (coming soon).

#### Maximum Value

The instruction set should pack as much value as possible in the limited space.
This means that clever tricks are used to get double use of instructions. Humans
are clever beings that can wield powerful tools in creative ways. It makes the
learning curve a bit higher, but the payoff is big.

For example, uCISC doesn't have jump instructions. A copy or ALU instruction will
do just fine if the resulting value is stored in the program counter. You can do
quite a lot by just varying the arguments.

```
# Do a relative jump by adding a displacement value to the program counter
0/copy/ 0.reg/from pc/ Entry.disp/displacement to Entry label/ 0.reg/into pc/

# Absolute jumps are also possible
0/copy/ 4.imm/from immediate/ Entry.imm/absolute offset to Entry/ 0.reg/into pc/
```

Another example of packing more value into less space is the "halt" instruction.
Since copying the value to the PC to itself with no offset will result in an
infinite loop, this particular combindation is treated as the "halt" instruction.
As an added benifit of the instruction set setup, you actually get conditional
halts by manipulating the effect argument.

```
# All of these instructions look at the flags set by the last ALU operation.
# So you do some math, comparisons, etc. and check the result.

# halt instruction (3.eff/always store/ is inferred here)
0/copy/ 0.reg/from pc/ 0.imm/zero offset/ 0.reg/into pc/

# halt if zero
0/copy/ 0.reg/from pc/ 0.imm/zero offset/ 0.reg/into pc/ 0.eff/if zero/

# halt if NOT zero
0/copy/ 0.reg/from pc/ 0.imm/zero offset/ 0.reg/into pc/ 1.eff/if not zero/

# halt if positive
0/copy/ 0.reg/from pc/ 0.imm/zero offset/ 0.reg/into pc/ 2.eff/if positive/
```

#### Constant Time Peformance

All instructions should execute in a consistent time on a processor. There are
a few cases where this is not possible (e.g. page to/from main memory), but
these should be rare. This implies consistent instruction size, roughly equivalent
instruction computation time, etc. A lot of effort has been made to make the
instructions have roughly the same flavor even across different types.

For example, ALU and Copy instructions feel roughly the same, even if there is
some variance in bit widths or semantics on a few arguments based on the context
of a copy verses an ALU operion. Notably, ALU operations have no immediate value
due to space constraints.

```
# Jump relative instruction
0/copy/ 0.reg/from pc/ 6.imm 0.reg/into pc/

# Jump relative with ALU instead
# First, store 6 into a register (r3 in this cas)
0/copy/ 4.imm 6.imm/immediate value 6/ 7.reg/into r3/
# Then, add r3 to PC
20A/add/ 7.reg/r3/ 0.reg/to pc and store in pc/

# Jump relative with ALU, this time with value off the stack
# May be useful for longer jumps that won't fit in an immediate
20A/add/ 1.mem/value on stack/ 0.reg/to pc and store in pc/
```

Achieving constant time performance also typically means there are very few forms
of caching allowed since cache misses necessarily produce longer execution times.

Some exceptions being considered:
 - Memory paging. Steps should be taken to ensure minimal variability in access,
   but accessing any shared resource will involve contention.
 - Device access, which is implemented through memory paging mechanics.
 - Page locks. Locking/unlocking pages involves accessing a shared resource.

Note: It is possible for a system to have page addresses that are any multiple of
16-bits. In this case, software will have to be recompiled as stack and memory
offsets will be different, but it doesn't fundamentally change the instruction set.
Regardless, the hardware must still conform to the constant time execution, which
implies a wider memory lookahead in the processor internals.

#### Easy to Implement in Hardware

Easy to implement is such a vague measure. More or less this means it should be
implementable on an FPGA by a hobbyist. Since I am a electronics hobbyist and
I'm targeting the ECP5 family of FPGA devices, I seem unlikely to violate this
constraint.

Further, you can find a digital logic design for this on CircuitVerse. Almost
anything that can be implemented in that very simple digital logic simulator
should conform to this goal. Checkout the
[uCISC project](https://circuitverse.org/users/6119/projects/58784) for more
details.

![uCISC Hardware Simulation](/images/uCISC_circuitverse.png)

Note, however, that easy to implement is intended to capture the fact that it
should be easy for a human to have an accurate mental model of how the hardware
works and does not mean that it is simple in the same way "RISC" architectures
intend. As the design depends on dual port memory, that necessarily adds
challenge at the hardware level.

#### Concept Scalable

This means concepts should work well on a single micro-processor core to large
core counts without breaking the single brain paradigm. A single mind should
reasonably be able to understand how the program will be have on 1 core and 100
cores. Significant portions of this goal are achieved by the local memory model
of uCISC processors. However, more work is being done on the paging and task
execution model with the goal of being able to reason about it across large
multi-core systems.

One tantalizing side effect of the model as it has been realized so far is that
it seems uCISC programs will not tend to be limited by single core performance
in the same way that typical modern programs are.

Once you break the barrier of 128k of locally available memory, you will tend
to be forced into a functionale style of programming where you are copying
chunks of memory to the local processor, computing the result and copying the
result out. Because you are forced to break your work into chunks that will
fit in 128k (along with the code needed to process the data), you'll naturally
be able to farm that work out to multiple cores and take advantage of all the
available processing power for normally single threaded work.

For example, let's say you are processing an image of some kind. This could be
resizing it or preparing it for screen rendering or constructing an updated
view of your application GUI. These processes are normally single threaded on
modern processors and operating systems. It's actuallly a lot of work to split
them up because of the hardware model and shared memory access.

For a uCISC based system, you can split the work into grids and render portions
of the image separately. Each portion can be executed on a separate processor
core since fundamentally the logic is the same if you do it on a single core or
multiple cores: copy part of the data to the local CPU, process it, copy out out.

#### Uncomplicated Operating System

Operating System kernels are a nightmare. I still remember taking my operating
systems class in college where I got the first taste of how complicated it was
to properly handle interrupts. What happens if you get interrupted while in the
middle of the interrupt? How do you call priviledged code and have the processor
constantly switch into priviledged execution mode? It's a nightmare to get right
even in the simplest cases.

uCISC aims to make an OS as simple and understandable as possible. This means we
want to avoid at least:

*privileged/non-priviledged divide* - Some of this is inevitable if the system
ever wants to run untrusted software. However, uCISC should minimize the
boundary and eliminate as much context switching as possible.

*interrupts* - Interrupt handling is nightmarishly complicated to get right:
 - Which processor/core should handle it
 - It's stopping process execution mid-stream and must carefully preserve state
 - It has to turn off interrupts while handling them to avoid interrupting itself
 - Interrupts can get lost if the processor is ignoring them
 - You need to handle interrupt starvation

*complicated device communication* - there are lots of various protocols for
moving data between the processor and other devices. The hardware is broad
and varied depending on performance, power, reliability and complexity. We need
the same instruction set to scale between I/O pins on a microcontroller and
potentially multi-lane PCIe communication across multiple device boards.

#### Security

While there is no all-encompasing solution, uCISC aims to design away as many
security flaws as possible. Admittedly, I'm not very experienced in desigining
processor security paradigms, but it seems that there are a few categories of
things to attempt to limit exposure to:

* Timing attacks:
  * All instructions must execute in constant time except for paging operations
  * No caching is allowed unless cache misses are executed in constant time
  * All memory access is local, uncontended and CPU specific.

* Process isolation:
  * The process can't access memory on other processes since all memory is local.
  * Shared memory access via paging, which is secured against cross processes and
  can only be referenced a page at a time (512 bytes).

* Memory safety:
  * Only a single thread can be running on a give core at a given moment
  * It has full access to the memory, and no access to other core's memory
  * Bits can't change unless you changed them
  * Data is copied (with optional locks) from main memory so you are forced to use
  a functional style of programming between threads.
  * Each 4k block of memory can be flagged as code or data, the process can't write
  to code blocks.

Some security concerns I'm still looking into:

* It would be nice to eliminate stack based attacks and prevent functions from
  overwriting stack data they aren't supposed to.
* Memory page probing. Since no virtualization of the memory space is allowed
  it means rogue code can probe the memory space and infer things about the
  system based on what it can/cannot access. This could be likely used to
  fingerprint a system or worse, identify programs that are running by memory
  behaviors.
* What, if any, machine virtualization should be supported? Is there a way to unify
  some VM and/or process virtualization mechanics that would be intuitive and support
  hardware level sandboxing of processes?

