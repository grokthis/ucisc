# uCISC Hardware Implementation

This is the very begining of the hardware implementation.
Currently, I'm targeting an ECP5 FPGA using the evaluation board,
with an 85k LUT device and a single instruction per cycle. You
can find the implementation under `ecp5_single_IPC`.

#### Current Status

1. I'm using the SymbiFlow project trellis workflow:
   https://github.com/SymbiFlow/prjtrellis
2. I have implemented the copy and ALU instructions and have
   some partial implementation of paging.
3. So far the verilog compiles and produces rougly the correct
   output:
```
=== processor_core ===

   Number of wires:              12904
   Number of wire bits:          30546
   Number of public wires:       12904
   Number of public wire bits:   30546
   Number of memories:               0
   Number of memory bits:            0
   Number of processes:              0
   Number of cells:              12809
     CCU2C                         119
     DP16KD                        128
     L6MUX21                      1751
     LUT4                         7498
     MULT18X18D                      4
     PFUMX                        3191
     TRELLIS_FF                    118
```

#### Implementation Notes

1. I need true dual port memory modules to implement my design
   since I need to read from two simultaneous addresses and write
   to one of those. However, the yosys + project trellis tool
   chain can't currently configure this on the ECP5 even though
   the sysMEM blocks support it. It can fake it by doubling
   the number of sysMEM blocks used (128 DP16KD instead of 64).
   This also incurs a cost of several thousand cells for the
   additional routing logic.

2. I have commented out the division and mod implementations in
   the ALU since they add 24k LUTs to implement.

3. I'm under no illusions that this will work the first time. But
   I have no easy way to get results out of the processor yet.
   I've decided to implment JTAG to read/write the contents of
   any memory, register or instruction decoding.

#### Next Steps

1. Implement JTAG interface so I can probe the registers
   and memory contents of the CPU and step through.
2. Write ruby implementation of JTAG debugger.
3. Load the fib and factorial examples and use JTAG to debug
   through them.

#### Compiling, Developing and Running

1. Install yosys, nextpnr and prtrellis according to the project
   trellis README: https://github.com/SymbiFlow/prjtrellis - note
   I could not get it to install on my Mac due to some library
   linking error, I had to use an ubuntu VM to get it to install
   properly.
2. `cd hardware/ecp5_single_IPC`
3. `make ucisc`
