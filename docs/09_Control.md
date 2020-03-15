# uCISC Control Instructions

NOTE: These instructions are a work in progress. You should
consider these pre-alpha. There are lots of issues and there is
almost zero chance this is anything like what the final result
will be.

```
7/control/

# Packed Instruction Bits:
111CCCCX XXDDDRRR
```

(C) Control code:

Need RRR:

* 0x0 - Try lock page against write, RRR is page
* 0x1 - Wait for lock page against write, RRR is page
* 0x2 - Try lock page against read, RRR is page
* 0x3 - Wait for lock page against read, RRR is page
* 0x4 - Unlock page for read and write, RRR is page

* 0x5 - Enable page for public read
* 0x6 - Disable page for public read

* 0x7 - Mark page as readonly, can't be undone without reset
* 0x8 - Reset page (PID 0 only)

RRR and DDD
* 0x9 - Queue message RRR to pid DDD
* 0xA - Cede page RRR to pid DDD
* 0xB - Execute page RRR as pid DDD (PID 0 only)

* 0xC - Kill PID from RRR (PID 0 only)
* 0xD - Unused
* 0xE - NOP
* 0xF - Unused

Virtualization controls, memory mapping and other possibilities:
 * message passing to other pids
 * control page lookup mapping for other pids
 * rate limit pids (allowed CPU count, task count, etc)

(R) Register for source address:

* 0.reg - from PID register

* 1.mem - Value at memory location (r1 + imm)
* 2.mem - Value at memory location (r2 + imm)
* 3.mem - Value at memory location (r3 + imm)

* 4.reg - from message register

* 5.reg - Value in r1 + imm
* 6.reg - Value in r2 + imm
* 7.reg - Value in r3 + imm


(D) Destination
* 0.reg - to PID register (PID 0 only)

* 1.mem - Value at memory location (r1 + imm)
* 2.mem - Value at memory location (r2 + imm)
* 3.mem - Value at memory location (r3 + imm)

* 4.msg - to message callback

* 5.reg - Value in r1
* 6.reg - Value in r2
* 7.reg - Value in r3

