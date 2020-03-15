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
4. [Processor Architecture](04_Processor_Architecture.md) - a high level overview
   of the processor architecture. Helpful background information on how instructions
   are executed in a real context.
5. [Copy Instructions](05_Copy.md) - a detailed description of the copy instruction.
6. [ALU Instructions](06_ALU.md) - a detailed description of the ALU instruction.
7. [Page Instructions](07_Page.md) - a detailed description of memory paging
   instructions.
