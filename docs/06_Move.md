# uCISC Move Instructions

Move instructions are targeted at moving information between registers.

```
D/move/ 0.reg/flags/ 0.imm/clear all/ 'to 0.reg/flags/

# Packed Instruction Bits:
1101DDRR IIIIIIII
```

#### Arguments

*(I) Immediate*

8-bit signed immediate value

The stored immediate value is left shifted once loaded.
All even values between -256 and 254 are supported. This
allows more range when using the immediate to offset memory
addresses since loading off the 2-byte boundary is not allowed.

*(R) Register Arg*

* 0.reg - Value in ((PC & 0xFF00) | immediate)
* 1.reg - Value in r1 + imm
* 2.reg - Value in r2 + imm
* 3.reg - Value in r3 + imm

*(D) Destination*

* 0.reg - Value in PC
* 1.reg - Value in r1
* 2.reg - Value in r2
* 3.reg - Value in r3

*Direction*

The direction is inferred from the combination of arguments and
source/destination postions. The immediate value must be given
after the source argument.

