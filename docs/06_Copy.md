# uCISC Copy Instructions

Copy instructions are targeted at copying data values referenced by
registers or located in registers. The copy instruction does not
transform the data in any way.

```
0/copy/ 2.mem/my var/ 'to 1.mem/stack/ 2.imm/stack offset/ 0.sign

# Packed Instruction Bits:
0EEDDDSS SCIIIIII
```

#### Arguments

*(I) Immediate*

6-bit (or 7-bit with C) signed immediate value

If the destination (D) is a memory location, the immediate is 6 bits
and the value must be between -32 and 31.

If the destination (D) is a register value then C is treated as the
most significant bit of the immediate value, increasing the range
of values to -64 to 63.

*(E) Effect*

Only the normal effects (0-3) are supported for copy instructions. See
instruction behaviors for more information.

*(C) Control Flag*

The control flag behavior depends on the destination argument.

Push:

If the destination argument is a memory value, the control flag
allows you to treat the destination like a stack pointer and push
the data to the stack. If C == 1, the destination address will be
decremented before the value is stored. The source argument
will not be changed.

Immediate:

If the destination is not a memory value then C is treated as the
most significant bit of the immediate value.

*(S) Source Argument*

Copy uses the standard source instruction behavior. Copy uses 4.val
since it supports immediate values.

*(D) Destination Argument*

Copy uses the standard destination instruction behavior.

