# uCISC System and Hardware

* Back: [Safety and Security](07_Safety_And_Security.md)
* Next: [Instruction Set](09_Instruction_Set.md)

uCISC is designed to prefer lots of small micro processor cores over larger more
powerful cores. In general, uCISC perfers to provide the software with predictable
transparent access to the hardware and let the software decide how to use this and
utilize the capabilities optimally.

For example, modern processors make heavy used of pipelining and superscaler
architectures to transform a sequential set of instructions into a parallel matrix
of execution units. The processors do advanced fetch, decoding and planning logic to
attempt to maximize the overall throughput of the processor. Pipelines aren't
terribly difficult to understand, but superlinear scaling is highly opaque to most
programming work. Because of these design constraings, processors perform a lot of
speculative execution that is hard to unwind during security faults. Further, this
leads to huge transistor counts. uCISC prefers to use additional transitors to add
more cores that can be software controlled effectively.

### Processor Architecture

#### Memory

*Local Memory*

The processor has up to 64k words of 16-bits each for local memory. A page of memory
is 64 words (128 bytes) and a block of memory is 4k words (8k bytes).

*Banked Memory*

The processor uses [bank switching](https://en.wikipedia.org/wiki/Bank_switching) to
provide access to other processors and devices. This works by mapping a specific
block to a device ID. There are 2^16 possible device IDs. Anything that the processor
can interact with must have a device ID, including main block memory, block devices,
serial devices, etc.

Once mapped, any reads or writes to the mapped memory block will actually read or
write from the specified device. The behavior and control block is device specific.
Note that the memory bank isn't actually directly connected to the local memory
of another processor or device. Instead, successful reads and writes depend on the
device or processor responding with the requested data.

Note, the specific hardware layout and access controls for each processor to each
device is left to the specific hardware implementation. Processors may be arranged
hierarchically or as peers or the hardware may allow access controls to be
customized. 

Writing to banked memory only works if the device on the other side is reading,
and reading only works if the other device is writing. Reads and writes block until
the other device responds accordingly.

```
# Processor A has 0x1000 banked to ID = B, 1.reg is SP, 2.reg is 0x1000
# Write to parent, blocks until parent reads it
0/copy/ 1.mem 2.mem/parent/
... other instructions

# Processor B has 0x4000 banked to ID = A, 1.reg is SP, 2.reg is 0x4000
# Push value from child onto stack, blocks until child writes
0/copy/ 2.mem 1.mem 1.push
... other instructions
```

This means a few things:

* The code on the various processors need to coordinate well or you can deadlock
* Neither the reader or writer know the real address accessed in the other processor
* It is possible, for processors to "write through" to other devices in one clock
  cycle. The intermediate processors just need to be writing to/from different banked
  memory spaces.

*Accessing Device Space*

The 2 least significant bits of the flags register control source and destination
between banked memory (1) and local memory (0). Reading/writing to banked memory may
actually be accessing local memory depending on the address accessed and the relevant
control word (see Banked Memory Control below).

You can set the bits for "from" and "to" independently to control the read and write
behavior using the flags register:

* 0x0 - read and write to local memory
* 0x1 - read from local memory and write to banked block
* 0x2 - read from banked block and write to local memory
* 0x3 - read and write to banked block

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
#                        # 0x0000 .. 0x0FFF => 0xNNNN #  # 0x0000 .. 0x0FFF #
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

* 0x0000 - 0x000F - Current processor device control segment
* 0x0010 - 0x001F - "Init" device control segment
* 0x0020 - 0x00FF - 14 non-banked control segments
* 0x0100 - 0x0FFF - 240 banked device control segment
* 0x1000 - 0xFFFF - 240 banked memory spaces (256 words each)

Each control segment configures the mapping for a 256 word block of memory. The upper
8 bits of the control segment determine which block is mapped. For example, the
control segment at 0x0070 to 0x007F controls the mapping of block 0x07NN.

Each bank maps to one and only one hardware device. The first control segment at
0x0000 is always the device control block for the current processor. The segment at
0x0010 is always "init" device. The init device may be another processor or some sort
of hardware controller. A halted processor will always be attempting to copy code
from it's init device.

If a device may is mapped to more than one memory bank, that means that it can handle
multiple interactions simultaneously. For example, some I/O devices may be able to
manage multiple in-flight requests at once. This same device can be mapped to
multiple banks to support this.

Non-banked devices don't have an associated memory bank. These may be attached to
small control devices that don't need banked memory to interact with. The 16 word
control space is sufficient to interact with these devices. For example, one such
device might be GPIO pins that can be written to or read from.

##### Banked memory control

The control word layout is as follows:

* 0x0 - Device ID, read only. Unique system wide.
* 0x1 - Bank address (MSB) | Device type (LSB)
* 0x2 - Bus access device ID
* 0x3 - Local mapped block (MSB) | Device status (LSB)
* 0x4 - Local register with interrupt handler address (read/write)
* 0x5 to 0xF - Device type specific (see below)

The device ID can be used to uniquely identify a device system wide if it is mapped
to more than one processor. This limits systems to 2^16 unique devices.

The bank address is a shortcut to load the bank address into a register. You can use
an ALU operation to zero the LSB (the device type) and get a reference to the bank
address. You can similarly zero the MSB and get the device type.

The interrupt code is device specific. However, any non-zero interrupt code must
indicate an interrupt is set. The possible status codes are listed below. Finer
grained status codes can be made available in the device specific area.

The interrupt handler register can be set to a local interrupt handler address for
that device. Each device must have it's own independent register and that register
must be local to this processor, even if the device is shared. If the register is
set to 0x0000 (the default) interrupts for this device are suppressed for the current
processor. If you don't need the handler address but want this device to trigger
interrupts, set to any non-zero value.

*Device types:*

* 0x00 - Missing
* 0x01 - Processor
* 0x02 - Block memory device
* 0x03 - Extended block memory device
* 0x04 - Block I/O device
* 0x05 - Extended block I/O device
* 0x06 - Serial device
* 0x07 - Human interface device
* 0x08 - Terminal device
* 0x09 to 0xFF - TBD

*Device status:*

0DEI BPRW

The device status is one of the following:

* W - Write ready (write to banked memory will succeed)
* R - Read ready (read from banked memory will succeed)
* P - Operation in progress (device is fulfilling a request)
* B - Device is busy, interactions are generally ignored
* I - Interrupt set, interrupt code must be non-zero
* E - Device is in error state
* D - Device disabled

If E or D are set, it's generally not safe to assume the other flags are correct. All
non-banked devices will have W and R set to 0.

##### Device control segment layouts

*Processor*

* 0x5 - Maximum local memory block (MSB)
* 0x6 - Control block address of next non-suppressed interrupted device (0 if none)
* 0x7 - 0.reg (PC)
* 0x8 - 1.reg
* 0x9 - 2.reg
* 0xA - 3.reg
* 0xB - flags register
* 0xC - control register
* 0xD to 0xF - inter-processor message

*Block memory device*

Block memory devices provide access to large memory address spaces. The memory device
is split up into blocks of 256 words. You can control the block address by setting
the value in the control block. The requested block may or may not become available
immediately in the banked address space, depending on the hardware.

* 0x4 (LSB) to 0x5 (MSB) - 32-bit block count, read only
* 0x6 (LSB) to 0x7 (MSB) - 32-bit block address, read/write
* 0x9 to 0x0F - manufacturer specific controls

*Block I/O device*

Block devices actually read/write from a buffer. A flush command must be issued to
trigger the write. A read command must be issued to read from the block device.

* 0x03 - Block size, 4k (0x1000) maximum, read only
* 0x04 (LSB) to 0x03 (MSB) - 32-bit block count, read only
* 0x05 (LSB) to 0x05 (MSB) - 32-bit block address, read/write
* 0x06 - Device command, TBD
* 0x07 to 0x0F - Unused

*Serial I/O device*

Serial devices actually read/write from a buffer. A write/flush command
must be issued to trigger the write. A read command must be issued to read from
the block device.

* 0x03 - Bit buffer size (32k maximum), read/write
* 0x04 - Device command, TBD
* 0x05 to 0x0F - Unused

#### Continue Reading

* Back: [Safety and Security](07_Safety_And_Security.md)
* Next: [Instruction Set](09_Instruction_Set.md)

