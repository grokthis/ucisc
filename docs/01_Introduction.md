# uCISC Introduction

uCISC (pronounced like micro-sisk) is a small, complex instruction set for
general purpose personal hobby computing that can scale from micro-controllers
to large server clusters using the same set of primitives understandable by
a single human head.

Processors fundamentally do one thing: move data around and transform it in
some way, usually with mathematics or bitwise operations. uCISC doubles down
on the data movement metaphor and centers the entire instruction set on it.

However, we want to build as compact an instruction set as possible, making
the primitives as few and as powerful as possible. We believe a computing
stack that fits in a single head is more comprehensible, more rewrite friendly
and leads to a more resilient society (see
[Mu: A minimal hobbyist computing stack](http://akkartik.name/post/mu-2019-1)
and
[Mu GitHub](https://github.com/akkartik/mu)).

While the instruction set is largely an invention I've been working on for a
while to scratch my hobbyist itch, somewhere along the line I came across
Mu and it helped really refine the approach to implementing the stack. My
goal is to take the single mind concept for Mu even further and push to
the hardware level.

## Design Goals

#### Maximum Micro

The instruction set should be as small as possible. Initially, the aim was
an 8-bit architecture, but in the end, that was too limiting because the limit
for immediate values was too small. Additionally, a 16-bit address space is
the minimum I could imagine to be useful.

All instructions, therefore, should be 16-bits wide including the immediate
value. Variable instruction widths are not allowed since this makes things harder
to mentally keep track of.

#### Maximum Value

The instruction set should pack as much value as possible in the limited space.
This means that clever tricks are used to get double use of instructions. Humans
are clever beings that can wield powerful tools in creative ways. It makes the
learning curve a bit higher, but the payoff is big.

For example, uCISC doesn't have jump instructions. A copy, move or transform
instruction will do just fine if the resulting value is stored in the program
counter.

#### Constant Time Peformance

All instructions should execute in a consistent time on a processor. There are
a few cases where this is not possible (e.g. page to/from main memory). This
implies consistent instruction size, roughly equivalent instruction computation
time, etc.

This also means no caching is allowed since cache misses necessarily produce
longer execution times. The local CPU memory is essentially the cache and paging
is manually handled by the program.

Some exceptions being considered:
 - Floating point arithmetic
 - Page locks. These are necessarily variable in a multi-core system.
 - Page address references when >16-bit page address
 - Back to back conditional jumps - generally speaking, implementations should
   take both branches on conditionals and throw away the incorrect one. However,
   there is a ratio of pipeline width to number of conditionals that can be
   effectively handled.

#### Easy to Implement in Hardware

The hardware should be understandable and easy to implement by single human (e.g.
on an FPGA device).

#### Concept Scalable

This means concepts should work well on a single micro-processor core to large
core counts without breaking the single brain paradigm. A single mind should
reasonably be able to understand how the program will be have on 1 core and 100
cores.

#### Uncomplicated Operating System

Operating System kernels are a nightmare. We want to make an OS as simple and
understandable as possible. This means we want to avoid at least:

*privileged/non-priviledged divide* - Some of this is inevitable if the system
ever wants to run untrusted software. However, the solution should sacrifice
flexibility for uniformity and simplicity wherever possible.

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

#### Avoid Intermediate Steps

RISC based systems copy values to registers and manipulate them there.
uCISC directly manipulates memory instead, avoiding the intermediate step
of copying to a register first and copying out after. The main benefit
of registers is speed (over memory access). Most processors use several
levels of cache to speed up memory access. uCISC uses a small amount of
memory that is as fast as the processor itself (constant time execution)
and just uses registers as addresses.
