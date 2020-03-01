# uCISC Copy Instructions

Copy instructions are targeted at copying data values referenced by
registers. In almost all cases, the register is de-referenced to get
at the memory value itself.

```
0/copy/ 2.mem/my var/ 'to 1.mem/stack/ 2.imm/stack offset/ 0.sign

# Packed Instruction Bits:
0EEDDDRR RMIIIIII
```

#### Arguments

*(I) Immediate*

6-bit signed immediate value

If source is 4.val, a value between -32 and 31 is supported.
If source is anything else, the value is left shifted once
and all even values between -64 and 62 are supported. This
allows more range when using the immediate to offset memory
addresses since loading off the 2-byte boundary is not allowed.

*(E) Effect*

0.eff - store if zero
1.eff - store if not zero
2.eff - store if positive
3.eff - store

*(M) Incre(M)ent on Modify*

0.inc - No increment
1.inc - Decrement source and destination mem registers by 2 BEFORE storing

This works on any X.mem arguments and increments the corresponding X.reg
address values by 2 before storing the result of this operation. The increment
happens only on X.mem arguments. The increment will only happen if the
effect results in a stored value.

*(R) Register Arg*

0.reg - Value in PC, add immediate

1.mem - Value at memory location (r1 + imm)
2.mem - Value at memory location (r2 + imm)
3.mem - Value at memory location (r3 + imm)

4.val - Value of immediate

5.reg - Value in r1 + imm
6.reg - Value in r2 + imm
7.reg - Value in r3 + imm

*(D) Destination*

0.reg - Value in PC

1.mem - Value at memory location r1
2.mem - Value at memory location r2
3.mem - Value at memory location r3

4.reg - Value in flags

5.reg - Value in r1
6.reg - Value in r2
7.reg - Value in r3


