# Common Devices

* Back: [Instruction Set](09_Instruction_Set.md)
* Next: [The Path Not Taken](11_The_Path_Not_Taken.md)

## Banked memory control

The control segment layout is as follows:

* 0x0 - Device ID - read only. Unique system wide.
* 0x1 - Local bank block (MSB) | Device type (LSB) - read only
* 0x2 - Init device ID - read only if set, writable if 0
* 0x3 - Remote bank block - writable by init device
* 0x4 - <Reserved> (MSB) | Device status (LSB)
* 0x5 - Local register with interrupt handler address (read/write)
* 0x6 to 0xF - Device type specific (see below)

The device ID can be used to uniquely identify a device system wide if it is mapped
to more than one processor. This limits systems to 2^16 - 1 unique devices since
device 0 is an invalid ID.

The bank address is a shortcut to load the bank address into a register. You can use
an ALU operation to zero the LSB (the device type) and get a reference to the bank
address. You can similarly zero the MSB and get the device type.

The interrupt handler register can be set to a local interrupt handler address for
that device. Each device must have it's own independent register and that register
must be local to this processor, even if the device is shared.

## Device Types

* 0x00 - Missing/Invalid
* 0x01 - Processor
* 0x02 - Block memory device
* 0x03 - Block I/O device
* 0x04 - Serial device
* 0x05 - Human interface device
* 0x06 - Terminal device
* 0x07 to 0xFF - TBD

## Device Status

0DEI BPRW

The device status is one of the following:

* W - Write ready (write to banked memory will succeed)
* R - Read ready (read from banked memory will succeed)
* P - Operation in progress (device is fulfilling a request)
* B - Device is busy, interactions are generally ignored
* I - Interrupt set, interrupt code must be non-zero
* E - Device is in error state
* D - Device disabled

If E or D are set, it's not safe to assume the other flags are correct. All
non-banked devices will have W and R set to 0.

## Device Control Segment Layouts

#### Processor

* 0x6 - Maximum local memory block (MSB)
* 0x7 - Control block address of next non-suppressed interrupted device (0 if none)
* 0x8 - 0.reg (PC)
* 0x9 - 1.reg
* 0xA - 2.reg
* 0xB - 3.reg
* 0xC - flags register
* 0xD - control register
* 0xE to 0xF - not used

#### Block Memory Device

Block memory devices provide access to large memory address spaces. The memory device
is split up into blocks of 256 words. You can control the block address by setting
the remote bank block in the control block. The requested block may or may not
become available immediately in the banked address space, depending on the hardware.

* 0x6 (LSW) to 0x7 (MSW) - 32-bit maximum block address, read only
* 0x8 - Extended remote bank block address, read/write
* 0x9 to 0x0F - manufacturer specific

If the block device has more than 2^16 blocks, (word at 0x7 is non-zero) then 0x8 can
be used as the most significant word of the block address. This setup is able to
address 2^32 blocks at 256 words each. This is up to 2 terabytes of data.

#### Block I/O Device

Block devices actually read/write from a buffer since the device access itself is
slow. The main difference from the processor's perspective is that a command must be
issued to initiate a flush or read from the device. This helps prevent spurious I/O
requests for blocks that don't actually need to be loaded.

* 0x6 (LSW) to 0x7 (MSW) - 32-bit maximum block address, read only
* 0x8 - Extended remote bank block address, read/write
* 0x9 - Device command, TBD
* 0xA to 0x0F - manufacturer specific

#### Serial I/O Device

Serial devices actually read/write from a buffer. To read/write from a serial device
you must set a pointer range to trigger the underlying hardware.

* 0x6 - Write buffer offset
* 0x7 - Read start inclusive
* 0x8 - Read end exclusive
* 0x9 - Write start inclusive
* 0xA - Write end exclusive
* 0xB to 0x0F - Not used

The write offset is the byte offset within the banked block where the read buffer
ends and the write buffer begins. The address is the first word in the write buffer.
Data in from the device will show up in the read buffer range. The processor
should write to the write buffer range.

The readable range is the full local address of the start and end addresses. The MSB
is not strictly necessary, but it is convenient to be able to directly read and write
the addresses and let the hardware strip the block address out.

If read start != read end, the interrupt flag will be set. If write start != write
end, the underlying device will write the data to the bust and update the write start
address. The buffers wrap within the address space of the block.

Note: the data format of the words is up to the device hardware. For example, devices
that need to write bits rather than bytes or words may encode the number of bits to
be transmitted in the word somewhere.

#### Continue Reading

* Back: [Instruction Set](09_Instruction_Set.md)
* Next: [The Path Not Taken](11_The_Path_Not_Taken.md)

