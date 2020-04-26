# uCISC System and Hardware

uCISC is designed to prefer lots of small micro processor cores over larger more powerful
cores. In general, uCISC perfers to provide the software with predictable transparent access
to the hardware and let the software decide how to use this and utilize the capabilities
optimally.

For example, modern processors make heavy used of pipelining and superscaler architectures to
transform a sequential set of instructions into a parallel matrix of execution units. The
processors do advanced fetch, decoding and planning logic to attempt to maximize the overall
throughput of the processor. Pipelines aren't terribly difficult to understand, but superlinear
scaling is highly opaque to most programming work. Because of these design constraings,
processors perform a lot of speculative execution that is hard to unwind during security faults.
Further, this leads to huge transistor counts. uCISC prefers to use additional transitors to add
more cores that can be software controlled effectively.

### Processor Architecture

The basic processor structure is a hierarchical structure where each processor has a single
parent (except the root processor) and up to 255 children processors. In order to access the
processors, they must be "banked" in 4k blocks.

#### Memory

*Local Memory*

The processor has up to 64k words of 16-bits each for local memory. A page of memory is 256
words (512 bytes) and a block of memory is 4k words (8k bytes). The highest available memory
address is accessible in the device address space at address 0x1000.

*Banked Memory*

The uCISC processor uses [bank switching](https://en.wikipedia.org/wiki/Bank_switching) to
provide access to other processors. This works by mapping a specific block to a processor ID.
Once mapped, any reads or writes to this memory block will actually read or write from the
specified processor. The memory bank isn't actually directly connected to the local memory
of another processor, but it is dependent on the other processor reading or writing the
requested data.

Banked memory example:

```
# Setup: We have configured the 4k block at 0x1000 to point to processor ID 0 (our parent)
# You do this by setting specific values in device space. We assume that is already done here.
0/copy 16 into 2.reg/ 4.val 10.imm 2.reg
212/swap MSB and LSB/ 2.reg 2.reg
# 2.reg now points to 0x1000

# Copy the value on the stack to our parent processor
0/copy/ 1.mem/SP/ 2.mem

# We can copy multiple words by using repeated instructions (see Instruction Behaviors)
0/copy 0x10 to flags register and set repeat/ 4.val 10.imm 4.reg
# Copying multiple words takes fewer instructions and may have a data width speed boost.
# See Instruction Behaviors for more details on performance.
0/copy 16 words to parent/ 1.mem 2.mem

# Read return value from parent (1 word), push to stack
# Copying multiple would work as well
0/copy/ 2.mem/parent/ 1.mem 1.push
```

*Accessing Device Space*

Bits 9 and 10 of the flags register control source and destination between devices space (1)
and memory space (0). Note that the lowest 4k block of device space is missing, so the address
range from 0x0000 to 0x0FFF always addresses memory space whatever the bits are set to. When
accessing memory space, the normal banking behavior applies.

Device space example:

```
0/copy/             4.val 2.imm 2.reg
# Set 2.reg to 0x0200, this has a 1 in the "From" flag position for device space
212/swap MSB, LSB/  2.reg 2.reg
# Configure the flags register
0/copy/             2.reg 4.reg

# Set 2.reg to point to 0x1000
0/copy 16 into 2.reg/ 4.val 10.imm 2.reg
212/swap MSB and LSB/ 2.reg 2.reg
# 2.reg now points to 0x1000

# Let's assume 1.reg points to some address <= 0x0FFF
# Since this address is less than 4k, it always points to memory space
0/copy maximum local memory address/ 2.mem 1.mem
# We now have the maximum memory address in 1.mem

# You could use this to init the stack pointer by storing it in the register instead
0/copy maximum local memory address/ 2.mem 1.reg
# Increment the stack pointer so that the first "push" decrements to the last available address
20E/increment stack pointer/ 1.reg 1.reg

# Use memory space again
0/copy 0 into flags reg, turning off device space/ 4.val 0.imm 4.reg/flags reg/
```

You can set the bits for "from" and "to" independently to control the read and write behavior.

Flags register:

* 0x0000 - read and write to memory space
* 0x0100 - read from memory space and write to device space
* 0x0200 - read from device space and write to memory space
* 0x0300 - read and write to device space

#### Hardware Specification

It is not necessary for all cores in a processor to have the same resources available. In fact,
it is expected that they won't. See the section on multi-core structure below. Each processor will
need to be queried via device space to see what it has available, how much memory, how many devices
and so on.

##### Figure 1 - Processor Architecture
```
#                             Up to 60k
#                             ####################
#                             # 0x1000 .. 0x2FFF #
#        #############        #------------------#
#        #           #        # 0x2000 .. 0x3FFF #
#        # Processor #====||  #------------------#
#        #   Core    #====||  #        ...       #
#        #           #    ||  #------------------#
#        #############    ||  # 0xD000 .. 0xDFFF #
#        Parent (ID = 0)  ||  #------------------#
#                         ||  # 0xE000 .. 0xFFFF #
#        #############    ||  ####################
#        #           #    ||               || Device Space
#        # Processor #================================================ 0x0000 (the void)
#        #   Core    #    ||          ||      Memory Space      ||
#        #           #    ||   Banked ||                        || Local
#        #############   ############################  ####################
#                        # 0x0000 .. 0x0FFF => 0xNN #  # 0x0000 .. 0x0FFF #
#                        #--------------------------#  #------------------#
#        Child 1         # 0x1000 .. 0x2FFF => 0xNN #  # 0x1000 .. 0x2FFF #
#        #############   #--------------------------#  #------------------#
#        #           #   # 0x2000 .. 0x3FFF => 0xNN #  # 0x2000 .. 0x3FFF #
#        # Processor #===#--------------------------#  #------------------#
#        #   Core    #===#           ...            #  #           ...    #
#        #           #   #--------------------------#  #------------------#
#        #############   # 0xD000 .. 0xDFFF => 0xNN #  # 0xD000 .. 0xDFFF #
#                        #--------------------------#  #------------------#
#        Child 2         # 0xE000 .. 0xEFFF => 0xNN #  # 0xE000 .. 0xEFFF #
#        #############   #--------------------------#  #------------------#
#        #           #   # 0xFFFF .. 0xFFFF => 0xNN #  # 0xFFFF .. 0xFFFF #
#        # Processor #===############################  ####################
#        #   Core    #===|| ||                   64k             Up to 64k
#        #           #      ||
#        #############      ||
#             ...           ||
#        Child N (N <= 256) ||
#        #############      ||
#        #           #      ||
#        # Processor #======||
#        #   Core    #======||
#        #           #   
#        #############   
#
#
```

#### Device Space

Device space contains metadata about the processor and allows the processor to control it's
memory configuration. It also allows access to any attached devices. A device can pretty much
be anything, including processor extensions like floating point units or I/O access. The same
device can be shared between more than one processor as appropriate.

The minimum device space for a processor is the first eight bytes (0x1000 - 0x1007). This is
valid for a processor with local memory and no children.

##### Device Space Layout

* 0x000 - 0xFFF - No device space for the first 4k words, always falls through to memory space

##### Processor configuration

* 0x1000 - Maximum local memory address (real memory, not banked)
* 0x1001 - MSB: Parent connection word width, LSB: Child connection word width
* 0x1002 - Child processor count (up to 255)
* 0x1003 - Interprocessor connection width in words (MSB: parent, LSB: children)
* 0x1004 - Processor manufacturer (0 for unspecified)
* 0x1005 - Processor device identifier
* 0x1006 - Processor features flags
* 0x1007 - MSB: Timer count (up to 208); LSB: Device count (up to 238)

##### Banked memory control

* 0x1008 ... 0x100F - 16 bytes (2 per word) controlling comm memory bank blocks
  Each byte assigns a 4k block of memory to read/write from a child processor
  or parent processor.
* 0x1010 - Bit flags to enable/disable banked comm memory for each 4k block
* 0x1011 - Write child processor ID here to force halt child processor.
* 0x1012 - Write child processor ID here to reset it.
* 0x1013 - Write child processor ID here to force jump zero. If the processor is
  halted, this unhalts the processor and it resumes normal execution from 0x0000.
* 0x1014 - Reserved
* 0x1015 - Reserved
* 0x1016 - Reserved
* 0x1017 - Reserved

##### Interrupt control

* 0x1018 - 0x101E - 238 bit flags indicating which devices have active interrupts set
  The last 2 bits are always zero.
* 0x101F - Interrupt index - reading this address blocks until an interrupt occurs
  * 0 - 237 is device interrupt ID
  * 238 - 511 timer and custom interrupts
  * 512 to 767 processor write interrupts

##### Clock control

* 0x1020 - Flags indicating if time address is functional (one bit per address from
  0x1020 to 0x102F)
* 0x1021 - Nanoseconds since epoch LSW (65.535 ms)
* 0x1022 - Nanoseconds since epoch (~71.6 minutes)
* 0x1023 - Unix clock milliseconds since epoch (~8.9 years)
* 0x1024 - Unix clock milliseconds since epoch MSW (~584.5 thousand years)
* 0x1025 - Timer 1 (milliseconds remaining) (-1 to disable)
* 0x1026 - Timer 2
* 0x1027 - Timer 3
* 0x1028 - Timer 4
* 0x1029 - Timer 5
* 0x102A - Timer 6
* 0x102B - Timer 7
* 0x102C - Timer 8
* 0x102D - Timer 9
* 0x102E - Timer 10
* 0x102F - Timer 11

##### Custom control space

* 0x1030 ... 0x10FF - Timers or custom extensions

##### Device Control

* 0x11XX ... 0x12XX - 238 device IDs; 16b'0 if device block is missing, 16'b1 general IO; it exists but has no device ID
  The last 18 addresses of this space always return 16b'0.
* 0x12XX ... 0xFFXX - 238 device blocks, details are hardware specific. Could be anything from FPU to SERDES to large memory I/O banks to GPIO pins
  Devices can be, but do not have to be specific to this processor. They could be shared between a few processors. Consult the specific hardware
  documentation for more details.

#### Multi-core Structure

Every processor except the root node has 1 parent processor. A processor can have up to 255, each
of which also have 255 children and so on. More likely, you'll have fewer branch nodes, carefully
chosen to support a certain amount of data bandwidth. Processors can have up to 238 devices attached
to them. Devices can be shared among groups of processors if desired. Typically, however, core
components will be attached to the root nodes, certainly the memory and system data bus.

For example, one configuration might be:
* Level 1: 1 node
* Level 2: 4 nodes
* Level 3: 8 nodes per level 2 parent (32 total)

Level 1 controls the entire system and has the core OS primitives running there. Levels 1 and 2 all
have the system hardware attached to their 238 devices. If you have more than 238, you can split them
between cores. There are therefore 5 OS nodes responsible for memory access, process control,
scheduling, handling interrupts and all the other OS type things. They are also responsible for acting
as a memory controller if you need one.

Level 3 are all processor nodes that typically run user space processes, but could also run kernel level
tasks if needed. They do not have direct access to system I/O, memory or other hardware. This means that
code running on the 5 branch processors effectively has priviledged access, even though there is no such
processor state. The goal of the OS should be to never allow remote code execution on these nodes.

#### Hardware Configurations

There are many other configurations dending on your design goals. For example:

* Attach FPU modules to all processors to accelerate floating point math.
  You could attach 10's or 100's to each core to provide vast amounts of floating point compute capacity.
* Attach graphics memory to a group of nodes. Viola! Graphics card.
* Increase the number of branch nodes relative to leaf nodes for higher data
  throughput. You could support hundreds of channels for memory access.
* Create super compute clusters by connecting a bunch of these together with a parent tree.
* Use a single core for an embeded device. You could fit 3-4 cores on the smallest ECP5 12k LUT FPGA.

