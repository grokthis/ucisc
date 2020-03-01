Control:   FFFCCCRR IIIIIIII  -  uCISC: 7/control/ 1:mem/stack/ [0x0:imm/stack-offset/] 0:arg/push-task/

(C) Control code:

* 000 - push task address
* 001 - unlock page

* 010 - lock page against write
* 011 - lock page against read

* 100 - halt/stop? (it pains me how many potential codes I waste if I do it this way, all of the R an I bits)
* 101 - _unused_

* 110 - _unused_
* 111 - _unused_

(F) Format (op code in []):

* 00XX [0] - Copy instructions
* 10XX [2] - Transform instructions
* 111X [7] - Control instructions
* 1100 [C] - Page instructions
* 1101 [D] - Move instructions

(R) Register for source address:

* 0 - If transform instruction, PC as reg
    - Otherwise, immediate value (sign extended if S == 1)
* 1 - r1 as mem address; if copy or meta, immediate is index (stack?)
* 2 - r2 as mem address; if copy or meta, immediate is index (refA?)
* 3 - r3 as mem address; if copy or meta, immediate is index (refB?)

