# Common Devices

* Back: [Instruction Set](09_Instruction_Set.md)
* Next: [The Path Not Taken](11_The_Path_Not_Taken.md)

## Standard Device Control

The control segment layout is as follows:

* 0x0 - Device ID - read only. Unique system wide.
* 0x1 - Flags (MSB) | Device type (LSB) - read only
* 0x2 to 0xF - Device type specific (see below)

The device ID can be used to uniquely identify a device system wide if it is mapped
to more than one processor. This limits systems to 2^16 - 1 unique devices since
device 0 is an invalid ID.

## Device Types

* 0x00 - Missing/Invalid
* 0x01 - Processor
* 0x02 - Block memory device
* 0x03 - Block I/O device
* 0x04 - Serial device
* 0x05 - Human interface device
* 0x06 - Terminal device
* 0x07 - Raster graphics device
* 0x08 to 0xFF - TBD

## Device Flags

DDDE XPRW TTTT TTTT

T is the device type

The device status is one of the following:

* W - Write ready (data writes will succeed)
* R - Read ready (data reads will succeed)
* P - Operation in progress (device is fulfilling a request), may not be busy
* X - Device is in error state
* E - Device is enabled. If zero (disabled), flags must all be 0
* D - Device specific flags, must conform to zero if disabled rule.

You can quickly check for enabled devices by checking if the MSB is non-zero.

## Device Control Segment Layouts

#### Serial I/O Device

Generally, you can think of serial devices as UART. There are both read and write
buffers which the hardware translates into serial I/O over UART lines.

Device specific flags:

PSWE XPRW 0000 0100

* W - Write ready (Tx writes will succeed)
* R - Read ready (Rx reads will succeed)
* P - Operation in progress (device is fulfilling a request), may not be busy
* X - Device is in error state
* E - Device is enabled. If zero (disabled), flags must all be 0
* W - Bit width, 0 for 8 bits, 1 for 9 bits
* S - Stop bits, 0 for 1 bit, 1 for 2 bits
* P - Parity bit, 0 for no parity bit, 1 for parity bit if supported

Note: the parity bit flag will only set when supported for the device.

* 0x2 - Clock divider
* 0x3 - Tx buffer write
* 0x4 - Rx buffer read
* 0x5 to 0x0F - Not used

Tx buffer and Rx buffer are only responsive when the write and read
device flags are set respectively.

The clock divider determines the baud rate depending on the clock signal from the
processor. The Tx/Rx clock will count the number of clock cycles in the divider
value as the time between data bits for the Tx/Rx signals. For example, for a 16MHz
clock, setting this value to 1667 will give you approximately 9600 baud.

#### Continue Reading

* Back: [Instruction Set](09_Instruction_Set.md)
* Next: [The Path Not Taken](11_The_Path_Not_Taken.md)

