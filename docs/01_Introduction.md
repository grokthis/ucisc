# uCISC Introduction

uCISC (pronounced like micro-sisk or myoo-sisk) is a small, complex instruction
set for general purpose personal hobby computing that can scale from
micro-controllers to large server clusters using the same set of primitives
understandable by a single human mind.

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

## Chapters

1. [Introduction](01_Introduction.md) - this introduction
2. [Design Goals](02_Design_Goals.md) - a description of the uCISC design goals.
   Helpful to understand the design decisions and ultimate intent of the architecture.
3. [Syntax](03_Syntax.md) - introduction to the uCISC code syntax. Breaks down how
   to read the simple uCISC code.
4. [System](04_System.md) - a description of the processor hardware setup and the system
   context it runs in, including device space and memory banking.
5. [Instruction Behaviors](05_Instruction_Behaviors.md) - common behaviors used by all
   uCISC instructions.
6. [Copy Instructions](06_Copy.md) - a detailed description of the copy instructions.
7. [ALU Instructions](07_ALU.md) - a detailed description of the ALU instructions.
