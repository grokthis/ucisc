Page:      FFF0NQRR IIIIIIII  -  uCISC: C/page 0:arg/in/ 1:mem/stack/ 0x10:imm/local-page/

(F) Format (op code in []):
* 00XX [0] - Copy instructions
* 10XX [2] - Transform instructions
* 111X [7] - Control instructions
* 1100 [C] - Page instructions
* 1101 [D] - Move instructions

(Q) Queue:
* 0 - No
* 1 - Yes, add to task queue

(R) Register for source address:
* 0 - If transform instruction, PC as reg
    - Otherwise, immediate value (sign extended if S == 1)
* 1 - r1 as mem address; if copy or meta, immediate is index (stack?)
* 2 - r2 as mem address; if copy or meta, immediate is index (refA?)
* 3 - r3 as mem address; if copy or meta, immediate is index (refB?)

(N) Direction:
* 0 - from source to destination
* 1 - from destination to source

