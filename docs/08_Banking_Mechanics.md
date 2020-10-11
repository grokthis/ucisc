# Banking Mechanics

* Back: [Stability And Security](07_Stability_And_Security.md)
* Next: [Instruction Set](09_Instruction_Set.md)

uCISC processors communicate with other processors and devices through a modified
memory [bank switching](https://en.wikipedia.org/wiki/Bank_switching) mechanism.

### Local Memory

The processor has up to 64k words of 16-bits each for local memory. A page of memory
is 64 words (128 bytes) and a block of memory is 4k words (8k bytes).

### Banked Memory

Each processor can have up to 254 devices mapped into it's banked memory spaces.
At design time, a computer should be architected for a specific use case.

Each system can have up to 2^16 - 1 devices (device ID 0 is used in a special case
and is therefore invalid as a device ID). Everything in the computer that is not the
processor itself is considered a device (even other processors are devices). The
device IDs are globally unique, but are intended only to help identify shared devices
in a complex setup. In general a processor should not have a single device mapped to
more than one bank as it will be confusing to write code for.

Normally, in banked memory systems, individual address spaces swap between
different devices or memory spaces. For uCISC processors and devices, the registers
swap between a local memory space, and a global banked memory space. In other words
you control banking on the registers, not the memory blocks. Time will tell if this
is more intuitive to work with. By convention, r1-r3 always refer to local memory
while r5-r7 refer to banked memory. You can modify the behavior by writing to r4,
but most programs should avoid it if possible.

Once banked, any reads or writes to the mapped memory block will actually read or
write from the specified device. Note: Instructions always load from local memory,
regardless of the bank status.

### Banking Rules

1. You can only access a halted processor, it's memory is inaccessible while running
2. You can, however, read and write to the init register and flags register, allowing
   the init process to halt a runaway process
3. Only the init device can write to a device, which is controlled by the init
   register. If the init device is 0, any device can "claim" it.
4. Read/writes to banked memory are blocking until the request can be served, unless
   you are not the init device, in which case your writes have no effect and reads
   always return 0.

### Accessing Device Space

The 16-bit control register turns on or off banking for each of the 6 general
purpose registers. Each register is controled by the 0-indexed bit that matches it`s
name. R1 is the second bit, R5 is the sixth, and so on. Bits 0, 4 and 8 and higher
are not currently used.

Note: instructions are *always* loaded from local memory. That is, all code must be
loaded into local memory before execution is possible. You can not execute code
directly from a device.

#### Hardware Specification

It is not necessary for all cores in a processor to have the same devices available.

##### Figure 1 - Processor Architecture

Use the following image to help understand the memory layout for uCISC processors.
```
#        #############
#        #           #
#        # Processor #============================================
#        #   Core    #                ||      Memory Space      ||
#        #           #         Banked ||                        || Local
#        #############   ##############################  ####################
#                        # 0x0000 .. 0x0FFF => Control#  # 0x0000 .. 0x0FFF #
#                        #----------------------------#  #------------------#
#        #############   # 0x1000 .. 0x2FFF => 0xNNNN #  # 0x1000 .. 0x2FFF #
#        #           #   #----------------------------#  #------------------#
#        # Device 1  #===# 0x2000 .. 0x3FFF => 0xNNNN #  # 0x2000 .. 0x3FFF #
#        #           #   #----------------------------#  #------------------#
#        #############   #            ...             #  #       ...        #
#                        #----------------------------#  #------------------#
#        #############   # 0xD000 .. 0xDFFF => 0xNNNN #  # 0xD000 .. 0xDFFF #
#        #           #   #----------------------------#  #------------------#
#        # Device 2  #===# 0xE000 .. 0xEFFF => 0xNNNN #  # 0xE000 .. 0xEFFF #
#        #           #   #----------------------------#  #------------------#
#        #############   # 0xF000 .. 0xFFFF => 0xNNNN #  # 0xFFFF .. 0xFFFF #
#                        ##############################  ####################
#        #############   ||                      64k             Up to 64k
#        #           #   ||
#        # Device 3  #===||
#        #           #   ||
#        #############   ||
#                        ||
#             ...        ||
#                        ||
#        #############   ||
#        #           #   ||
#        # Device N  #===||
#        #           #
#        #############
```

#### Banked Memory Layout

The first segment of banked memory contains "control segments" which provide
information about the banked devices and let the processor control them.

* 0x0000 - 0x000F - Current processor device control segment
* 0x0010 - 0x00FF - 15 device control segments, no corresponding memory block
* 0x0100 - 0x0FFF - 240 device control segment, memory blocks starting at 0x1000
* 0x1000 - 0xFFFF - 240 banked memory blocks (256 words each)

Each bank maps to one and only one hardware device. The first control segment at
0x0000 is always the device control block for the current processor. The segment at
0x0010 is always "init" device.

#### Continue Reading

* Back: [Stability And Security](07_Stability_And_Security.md)
* Next: [Instruction Set](09_Instruction_Set.md)

