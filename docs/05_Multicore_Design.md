## Multicore Design

* Back: [Chips and Hardware](04_Chips_And_Hardware.md)
* Next: [Memory and Devices](06_Memory_And_Devices.md)

When looking at the design goals and available FPGA's on an open source stack, we
need to be able to run an estimated 20 16-bit CPU cores on a ECP5 84k FPGA to
achieve what I am aiming for.

#### Avoiding Complexity in Multicore Processors

This presents a challenge. Multi-threaded CPU architectures have some inherent
complexities that I would like to avoid. To keep everything transparent and
understandable the following decisions come naturally:

*Complete Processor Affinity*

Modern processors and operating systems are built around
[Task Preemption](https://en.wikipedia.org/wiki/Preemption_%28computing%29). Well,
really, they have been since we started running multiple processes on a single
machine. I would like to avoid this for uCISC, however, as preemption adds a lot
of complexity, context switching and human brain space when dealing with it. Since
we are trying to fit the whole system in a single mind, adding this complexity is
best avoided.

Basically, this means that all of our executing code will have
[Processor Affinity](https://en.wikipedia.org/wiki/Processor_affinity).

*No Global Address Space*

A 16-bit address is not enough to address a decent amount of memory. We could get
around this by having either extended address instructions or some fancy logic in a
[memory management unit](https://en.wikipedia.org/wiki/Memory_management_unit). We
would also have to deal with shared memory, process permissions and build a complex
OS just to run simple code.

I spent quite a bit of time trying to figure how how to implement some sort of
memory management scheme and it just adds so much complexity. We avoid almost all of
this complexity by only supporting "local memory" which is only addressable to the
local processor and under it's full control. We will need a way for processors to
communicate, however. More on that later.

*No Preemptive Interrupts*

For the same reasons we want to avoid process preemption, we want to avoid
preemption for any reason or we won't be able to ditch that complexity. Interrupts
have the added complexity of needing to switch in/out of privileged mode to execute
and then you need to write very, very hard to understand and carefully crafted code
to not mess it up and avoid interruption while you are handling interrupts among
other things.

Therefore we will only support interrupt handling via polling mechanisms. Note that
this dovetails nicely with processor affinity.

*No Privilege Mode*

A computer, designed to be transparent to the user, with complete processor affinity
and only local memory really doesn't need a privileged execution mode. There is
an open question of how to avoid a rogue process taking over the whole machine, which
we will deal with in [Safety and Security](07_Safety_And_Security.md) but without
privileged access. This greatly simplifies how the processor works as 100% of the
processor is availble to you and under your control 100% of the time.

#### Running Multiple Cores

So, having done away with all these complex processor modes, we still need to be
able to run multiple processes on different cores to be able to do things like
render a screen. Also, we will need a way to launch processes on peer cores without
special processor modes in a way that a human can easily mentally keep track of.

uCISC approaches this problem through
[memory bank switching](https://en.wikipedia.org/wiki/Bank_switching). With banking
you can set change which device or slot belongs to a given memory address. In this
case, we say that address 0x1000 may normally refer to a local memory word. However,
with banking mode turned on, it could refer instead to an address in another CPU
address space. I could turn on banking for 0x1000 and copy some of my CPU's data
into the address space of another CPU.

To make this work without privileged modes or giving processes direct memory access
to other processes running on other cores, we want to simplify this access. The
following design choices will allow us to keep multicore processes simple while
giving users the ability to controll what happens in their system.

*As a process:*

* I control my own memory space at all times. I can *optionally* designate a block
of memory that I will allow other processes to write into.
* I control which other processor is allowed to write into that space (only one at
a time).
* If a processor I have set to have write access to a block of memory halts, that
chain is automatically broken. I have to re-establish it if I want to use it.

That is, for a processor to communicate, both must want to talk to each other and
directly setup the connection.

*As a halted processor:*

1. I will erase all local memory before resuming execution
2. I will only communicate with a designated "init" processor that is hardwired
3. I will allow the init processor to write code to my local memory
4. I will allow the init processor to start code execution

With these basic rules you can design a CPU structure that will allow the OS to
maintain system control at all times while still executing code that may be less
stable.

As a simple example:

* Processor 0 is the boot processor with the lower banks mapped to ROM
* Processor 0 is the "init" processor for processors 1-20
* Processor 0 can restart any failed processors or serve requests from other
processors to start a new process.

This is more or less what I intend to build for the uCISC computer. I imagine that
the render engine will estabilish several processes running on different cores that
each render a sub-section of the screen. Other cores will be left to other user
processes.

To make this work well, I think processor 0 needs the following primitives:

* Launch new process
* Allow a process currently executing to voluntarily halt, to be restarted on some
event (say, network activity or keyboard input)
* Monitor devices for interupts and schedule work on available cores
* Restart dead processes
* Kill locked processes (e.g. I accidentlly wrote an infinite loop)

If we give processor 0 permission to force halt another processor, we can keep the
machine from being locked up by unstable or malicious code. You can read up on the
[Banking Mechanics](08_Banking_Mechanics.md) to see how this works. It's fully
achieved via well defined structured memory banking. While doing so, it avoids some
really nasty hardware complexities and programming complexities.

Read on about how we can generalize this approach for accessing large amounts of
memory and just about any hardware device you can plug into a computer.

#### Continue Reading

* Back: [Chips and Hardware](04_Chips_And_Hardware.md)
* Next: [Memory and Devices](06_Memory_And_Devices.md)

