# Banking Mechanics

* Back: [Stability And Security](07_Stability_And_Security.md)
* Next: [Instruction Set](09_Instruction_Set.md)

uCISC processors communicate with other processors and devices through a memory
[bank switching](https://en.wikipedia.org/wiki/Bank_switching).

*Local Memory*

The processor has up to 64k words of 16-bits each for local memory. A page of memory
is 64 words (128 bytes) and a block of memory is 4k words (8k bytes).

*Banked Memory*

Each processor can have up to 254 devices mapped into it's banked memory spaces.
At design time, a computer is architected for a specific use case.

* "Init" processors are responsible for bootstrapping other processors
* "Device" processors have statically mapped devices attached to them
* "User" processors have no devices attached

Each system can have up to 2^16 - 1 devices (device ID 0 is invalid). Every processor
memory bank, I/O connection and more is a device. The device IDs are independent from
the bank connections. A processor may not have a single device mapped to more than
one bank. A processor can turn on or off each bank independently.

Once banked, any reads or writes to the mapped memory block will actually read or
write from the specified device. Note: Instructions are always loaded from local
memory, regardless of the bank status.

Some rules:

1. You can only bank a halted processor, it's memory is inaccessible while running
2. Read/writes to banked memory are blocking until the request can be served
3. Normally, you can only bank non-disabled processors
4. A disabled processor always allows its init processor

Take a system with 4 processors (ID 1, 2, 3 and 4) as an example. Let's assume that
each processor is banked to address 0xN000 for every other processor, where N is the
processor ID. So, processor 2 banked memory is accessible to processor 1 at 0x2000.
Let's also assume that the control register is at 0x00N2 (which will be explained in
a moment).

For the case where processor 1 wants to split off some work into processor 2:

1. Processor 1 writes it's own ID as the init device to processor 2
2. Processor 1 verifies processor 2 has the init device set to its own ID
3. Processor 1 is now free to read/write from any register or block on processor 2
4. When ready, processor 1 clears the halt flag for processor 1 and execution resumes
at the next instruction (which processor 1 should have set correctly).
5. Processor 1 can periodically check for halt on processor 2 and read out the
results when it is done.

If processor 2's device is not 0, it wont' be writable. Step 1 will silently fail.
Processor 1 will need to read back the set ID to make sure it worked in step 2.
Once the init device is set to it's own device ID, no other processor can make any
changes to it and processor 1 can be assured its changes will be effective.

*Accessing Device Space*

The 16-bit control register turns on or off banking for 4k segments of memory. The
LSB is for 0x0000-0x0FFF and the MSB is for 0xF000-0xFFFF. When set to 1, reads and
writes to the segment are to/from banked memory instead of local memory.

Note: instructions are *always* loaded from local memory. That is, all code must be
loaded into local memory before execution is possible.

#### Hardware Specification

It is not necessary for all cores in a processor to have the same devices available.

##### Figure 1 - Processor Architecture
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
* 0x0010 - 0x001F - "Init" device control segment
* 0x0020 - 0x00FF - 14 non-banked device control segments
* 0x0100 - 0x0FFF - 240 banked device control segment
* 0x1000 - 0xFFFF - 240 banked memory spaces (256 words each)

Each control segment configures the mapping for a 256 word block of memory. The upper
8 bits of the control segment determine which block is mapped. For example, the
control segment at 0x0070 to 0x007F controls the mapping of block 0x07NN. For
convenience, one of the bytes of the control segment is the upper byte of the memory
segment. This allows you to load that address directly into a register rather than
compute it yourself.

Each bank maps to one and only one hardware device. The first control segment at
0x0000 is always the device control block for the current processor. The segment at
0x0010 is always "init" device.

Non-banked devices don't have an associated memory bank. These may be attached to
small control devices that don't need banked memory to interact with. The 16 word
control space is sufficient to interact with these devices. For example, one such
device might be GPIO pins that can be written to or read from.

#### Continue Reading

* Back: [Stability And Security](07_Stability_And_Security.md)
* Next: [Instruction Set](09_Instruction_Set.md)

