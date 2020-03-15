# uCISC Copy Instructions

Copy instructions are targeted at copying data values referenced by
registers. In almost all cases, the register is de-referenced to get
at the memory value itself.

```
0/copy/ 2.mem/my var/ 'to 1.mem/stack/ 2.imm/stack offset/ 0.sign

# Packed Instruction Bits:
0EEDDDRR RCIIIIII
```

#### Halt Instructions

If you copy 0.reg to itself with an immediate offset of 0 this would
normally put the processor into an infinite loop. Instead this is
interpreted as a HALT instruction.

Note that the effect still applies and it is only considered a halt
if the effect results in a store operation. This allows you to do
conditional halts based on the flag results or use 3.eff to halt
unconditionally.

0EE00000 00000000

#### Arguments

*(I) Immediate*

6-bit signed immediate value

If source is 4.val, a value between -32 and 31 is supported.
If source is anything else, the value is left shifted once
and all even values between -64 and 62 are supported. This
allows more range when using the immediate to offset memory
addresses since loading off the 2-byte boundary is not allowed.

*(E) Effect*

* 0.eff - store if zero
* 1.eff - store if not zero
* 2.eff - store if positive
* 3.eff - store

*(C) Control Flag*

The control flag has a few interesting effects depending on the destination
argument of the instruction.

Push:

If the destination argument is a memory value, the control flag
allows you to treat the destination like a stack pointer and push
the data to the stack. If C == 1, the destination address will be
decremented by 2 before the value is stored. The source argument
will not be changed.

Immediate:

If the destination is not a memory value then the bit is treated as the
most significant bit of the immediate value, increasing the range
of values to -64 to 63.

*(R) Register Arg*

* 0.reg - Value in PC, add immediate

* 1.mem - Value at memory location (r1 + imm)
* 2.mem - Value at memory location (r2 + imm)
* 3.mem - Value at memory location (r3 + imm)

* 4.val - Value of immediate

* 5.reg - Value in r1 + imm
* 6.reg - Value in r2 + imm
* 7.reg - Value in r3 + imm

*(D) Destination*

* 0.reg - Value in PC

* 1.mem - Value at memory location r1
* 2.mem - Value at memory location r2
* 3.mem - Value at memory location r3

* 4.reg - Value in flags

* 5.reg - Value in r1
* 6.reg - Value in r2
* 7.reg - Value in r3


