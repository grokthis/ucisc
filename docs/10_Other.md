# Instruction Set Architecture

Memory Layout:
Processes can address 2^16 by default (or 2^64 pages by setting the page
address width) address. Each page is 256 bytes, and therefore a total of
2^24 bytes (16Mb) is addressable by default (2^72 bytes under max config).

Page address ranges:

```
0000 0000 - 0000 0001 : 256 bytes, device 0 (main board) configuration
                      : Includes info on number of configurable devices (say FFFE)
                      : Includes task queue size (say FFFF)
0000 0000 - 0000 FFFF : ~16m, Device configuration range, exact range depends on
                      : number of configurable devices, 1 page per device (65,535 in this case)
0001 0000 - 0001 FFFF : 16m, Task queue, exact range depends on task queue size (e.g. 64k)
0002 0000 - FFFF FFFF : Main memory, this takes up the remainder of the memory space
```

Notes:
 - At startup, the BIOS should be mapped over 0000 0000 and following. Once BIOS init is complete
   main board configuration can unmap the BIOS.
 - The address ranges are determined by the contents of the first memory page. The details
   of that memory layout and which pieces are r/w are device specific, but it should at a
   minumum include the number of configurable devices (though not all may be available) and
   the task queue size (in pages).
 - Main board configuration and device configuration are intended to be memory mapped to the
   connected devices. That is, they are not likely actually RAM and the contents and r/w details
   are device specific.
 - All available devices must have a non-zero byte in the first position of the page. A device
   may return zero as the first byte to indicate it is disabled or otherwise unavailable. This
   generally should not be used as a "busy" indicator.
 - The device pages are allowed to have completely different latency characteristics. The task
   queue range is a shared, high performance resource. Details vary by implementation and needs
   but conceptually it can be thought of as performing similar to a shared cache.
 - Main memory will typically be DRAM speeds, though details are implementation specific.
 - Portions of main memory may be mapped to devices, this should be reflected in the device
   configuration as appropriate. The specifics of how to do this are out of the scope of the ISA.

The Task Queue:

The task queue is a LIFO queue that stores 32-bit page addresses. Each page address references
the memory page for the execution entry point for the task. A processor will copy the content
of that memory location into it's local memory and begin execution at the first address. The
first page will need to fetch any additional pages needed from main memory. The following
gaurantees are made:

 - A task address will be handed to at most one processor
 - The task push and pop operations are atomic
 - A task can be queued multiple times

Interrupts:

Interrupts are implemented as tasks. Each device that can send interrupts needs to be initialized
with a task address. Upon interrupt, the device will push the task address into the task queue and
a processor will pick it up.
