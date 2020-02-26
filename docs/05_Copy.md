# uCISC Copy Instructions

Copy instructions are targeted at copying data values referenced by
registers. In almost all cases, the register is de-referenced to get
at the memory value itself.

```
0/copy/ - 2.mem/my var/ -to 1.mem/stack/ 2.imm/stack offset/ 0.extend

# Packed Instruction Bits:
0SNDDDRR IIIIIIII
```

#### Arguments

*(I) Immediate*

8-bit immediate value, extended to 16-bits according to (S)

*(S) Sign Extend*

0 - Zero extend the immediate to 16-bits
1 - Sign extend the immediate to 16-bits

*(R) Register Arg*

0.val (if source) - Value of immediate
0.reg (if target) - Value in PC, add immediate

1.mem - Value at memory location (r1 + imm)
2.mem - Value at memory location (r2 + imm)
3.mem - Value at memory location (r3 + imm)

*(D) Destination*

0.reg - Value in PC
1.mem - Value at memory location r1
2.mem - Value at memory location r2
3.mem - Value at memory location r3

4.reg - Value in flags register
5.reg - Value in r1
6.reg - Value in r2
7.reg - Value in r3

*(N) Direction*

Value is transferred from:

0.dir - source to destination
1.dir - destination to source

This argument is inferred from the combination of arguments and
source/destination postions. Valid copy statements have one argument
matching the R options and one matching the D options. The order of
those two arguments determines the value N takes. Also, the immediate
value can only modify the R position.

Therefore the following combinations are not valid:

 - Imm value following 4-7.reg
 - Two 4-7.reg arguments
 - Two imm arguments
 - Two val arguments

