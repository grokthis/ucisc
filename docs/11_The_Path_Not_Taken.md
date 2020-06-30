## The Path Not Taken

* Back: [Common Devices](10_Common_Devices.md)
* Next: [Syntax](12_Syntax.md)

With all the stops and starts on this project, there are things that I tried that
either didn't work or ended up not fitting in the architecture. They were cut
but are documented here for posterity.

This document is a work in progress and most of it is just text dumped here because
it contains relevant info. Needs work.

TODO: To document:

### 8-bit ISA

TODO: I tried 8-bits first. Describe why it didn't work.

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


### Variable width memory copy

TODO: I tried adding multi-word copy for performance. Hasn't worked out yet.

### Chained math functions

TODO: I tried allowing for arithmetic chaining. Describe why it was removed.

### Memory paging instructions

TODO: Banking worked better than paging instructions

### Clocks timers, etc.

TODO: I removed clocks and timers as a native part of the spec. Describe why.

### "Parent" processors

TODO: I had a more hierarchical processor structure in earlier iterations.

---------

TODO: The text below needs to be mined for design decisions that didn't work out
and and added above, or if the decision stuck, moved to another doc.

#### RISC Architecture

TODO: Polish this. The early iterations were more RISC like.

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

##### Constant Time Peformance

TODO: This led to good things, but not sure it should be a hard requirement

All instructions should execute in a consistent time on a processor. There are
a few cases where this is not possible (e.g. page to/from main memory), but
these should be rare. This implies consistent instruction size, roughly equivalent
instruction computation time, etc. A lot of effort has been made to make the
instructions have roughly the same flavor even across different types.

Achieving constant time performance also typically means there are very few forms
of caching allowed since cache misses necessarily produce longer execution times.

### Concept Scalable

TODO: I abandoned the non-homebrew segment

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

#### Continue Reading

* Back: [Common Devices](10_Common_Devices.md)
* Next: [Syntax](12_Syntax.md)

