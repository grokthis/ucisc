## Multicore Design

* Back: [Multicore Design](05_Multicore_Design.md)
* Next: [Stability And Security](07_Stability_And_Security.md)

Multiple processors are not enough. Our target design calls for 20 processors with a
theoretical maximum of 64k words per processor, though the FPGA only has enough for
10k words per processor. 200k of memory is far short of the multiple GB of memory
we are likely to need to meaningfully hold data that can be displayed on a 1080p
monitor, so we need some scheme to deal with that.

#### Banked Memory

As it turns out, the memory banking mechanics that we use to communicate between
processors can be generalized to memory as well. If you have the ability to swap
which memory space is accessed for a given local address, we can use the same
approach to swap in blocks of main memory.

To make this work, the processor needs a control bank to specify the the address for
the main memory. The basic idea is to write control information to a well known
location that triggers an address lookup in the memory system. That block is then
mapped to a memory bank and the processor can read from that memory directly or copy
it to local memory.

This isn't quite enough, however, as we will need to access memory from multiple
processors. Which processor will controll the address and read the data? It is
possible to build a locking mechanism to deal with the contention, but this is
complicated to implement and difficult to predict how it will perform. A better way
will be to directly attach memory to processors so they have sole control over the
specified memory bank. Other processors can request pages from the memory linked
processor, which will serve requests in order and copy the data to the requesting
processor via banked mechanics.

#### Memory Performance

In order to drive our displays, we need to be able to write 200 MW/s (million words
per second) to video memory. A single 100MHz processor can drive a throughput of
around 25 MW/s, suggesting that we need to have 8 processor cores dedicated to
serving memory requests in such a system. We will also want a few cores to serve
requests for main memory to store application data and other non-video memory.

This suggests the following system architecture:

* Processor 0 - bootstrap processor, init processor + OS call handler
  * 8 memory cores driving independent video memory pipelines
  * 2-4 memory cores driving independent memory lanes
  * 10 application cores running the user processes on the computer

To write code for this system, I will basically need to write a memory controller
and run several copies of it. This controller can range from a simple implementation
serving requests to other processors to more complex memory space mapping for
processes, basically taking on some of the functions a MMU would perform. You could
implement process address translation if you wanted, but I don't currently see a
need for this until I want to run code that I don't trust. I have no plans to do
that at this time.

Because of the need for copying memory blocks around the system, I have wanted to
figure out a way to implement a batch memory copy operation in the instruction set.
So far, my attempts haven't worked out. I need something that is easy to implement
in hardware and seamless to use in software. As I have settled on a page size (64
words) it seems more likely to fit somewhere if I focus on speeding up page transfers
rather than variable sized transfers.

If there was a page copy operation, it would save significant processing time since
instead of 4 instructions per word, the processor could transfer one word per clock
cycle in chunks of 64 words. This would reduce the video memory CPU cores by a
factor of 4. We could devote probably 16 cores to executing user applications. Some
implementations could use the memory clock cycle and multi-word transfers to speed
this up even more, giving these memory controller processors much higher throughput.

#### Devices

This banking system can be expanded to include all devices, not just processors and
main memory. Any device can be mapped into block of memory plus the ability to issue
commands. Even serial I/O devices like USB can utilize a buffer to write these
commands to the bus.

The general solution is for processors, memory, I/O devices, etc. to all use banked
memory to access those devices. Each device needs a "control" space to read
configuration information and write to control registers. If a device needs either
block data or streamed data, it will need to be mapped to a banked block.

The banked memory layout will be as follows:

* 0x0000 - 0x000F - Control segment for the current processor
* 0x0010 - 0x001F - Control segment for the "init" device
* 0x0020 - 0x00FF - 14 non-banked control segments (no corresponding banked block)
* 0x0100 - 0x0FFF - 240 banked device control segment (16 words each)
* 0x1000 - 0xFFFF - 240 banked memory spaces (256 words each)

This means that each processor can control up to 14 non-banked devices and 240 banked
devices. To simplify system configuration, each device should be statically linked
to a certain processor at a specific bank position. If a single physical device can
support interactions with multiple processors at once, it will need to map to
multiple logical devices.

For more information on the specifics, see the
[Banking Mechanics](08_Banking_Mechanics.md) documentation. It goes into much more
details on how the address space is laid out and how to access it. Also, the
[Hello World](14_Hello_World.md) example uses banked memory to access a terminal
device and print a string.

#### Continue Reading

* Back: [Multicore Design](05_Multicore_Design.md)
* Next: [Stability And Security](07_Stability_And_Security.md)

