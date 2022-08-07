# Common Devices

* Back: [Instruction Set](09_Instruction_Set.md)
* Next: [The Path Not Taken](11_The_Path_Not_Taken.md)

## Standard Device Control

The control segment layout is as follows:

* 0x0 - Device ID - read only. Unique system wide.
* 0x1 - Flags (MSB) | Device type (LSB) - read only
* 0x2 to 0xF - Device type specific (see below)

The device ID can be used to uniquely identify a device system wide if it is mapped
to more than one processor. This limits systems to 2^16 - 2 unique devices since
device 0 is an invalid ID.

## Device Types

* 0x00 - Missing/Invalid
* 0x01 - Processor
* 0x02 - Block memory device
* 0x03 - Block I/O device
* 0x04 - Serial device (including UART)
* 0x05 - Human interface device
* 0x06 - Terminal device
* 0x07 - Raster graphics device
* 0x08 - GPIO device
* 0x09 - I2C device
* 0x0A - SPI device
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

#### GPIO Device

GPIO Devices drive a set of GPIO pins. Each pin can be set to read or write.
When reading from the device, you always read the current values. For input pins
this is the value at the pin. For outputs, this is the current value of the output.
Each GPIO device can support up to 64 pins.

Device specific flags (most flags are read only):

Packed flag bits:
```
PPPP PPIB 0000 1000
```

* P - 6 bit number indicating the maximum GPIO index. 0-63 for 1-64 pins.
* I - Indicates if this is an I/O expander device
* B - Busy indicator. Indicates if I/O state is being refreshed
      (always 0 unless I/O expander)

Note: if the device is an I/O expander, a refresh can be triggered by
writing a 1 to the busy indicator. It will automatically clearn when done.

Control block addresses:

* 0x2 - Signal device address (e.g. I2C device address)
* 0x3 - Reserved
* 0x4 to 0x7 - Pin configuration for each pin (1 = out, 0 = in)
* 0x8 to 0xB - Pin values for output enabled pins
* 0xC to 0xF - Current input value for all pins (current pin value for outputs)

#### I2C Device

I2C devices communicate with I2C components on the bus.

Device specific flags (most flags are read-only):

Packed flag bits:
```
00CC 00NA 0000 1001
```

* A - (R/W) address mode (0 = control, 1 = block)
* N - (R/W) enables NACK, rather than ACK for last read byte
* C - Indicates the error code for the last transmission. 0 indicates
      no error.

Control block addresses:

* 0x2 - I2C bus device address (i.e. I2C device address)
        The address is in the lower byte of the word. The 7-bit
        I2C address is in the most significan 7 positions of the byte
        and has the R/W bit at the end. You can change the R/W mode
        by flipping the LSB in the register.
* 0x3 - The remaining bytes in the current I2C interaction. If zero,
        the message is complete and the device is not busy. To initiate
        a communication, write a positive integer to this register and
        the device will either send or receive that number of bytes.
        If the device is in control mode, at most 24 bytes may be
        transferred at once. If in block mode, at most 512 bytes may be
        transferred at once.
* 0x4 to 0xF - Data buffer for control mode (up to 24 bytes). Bytes are
        written or read to/from the least significant byte in the word first.

Note: Not all I2C devices will allow block mode. In addition specific I2C
devices may limit the maximum number of bytes available in control mode. See
the specific documentation of the device in question.

#### Continue Reading

* Back: [Instruction Set](09_Instruction_Set.md)
* Next: [The Path Not Taken](11_The_Path_Not_Taken.md)

