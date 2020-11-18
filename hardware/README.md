# uCISC Hardware Implementation

The reference hardware implementation uses [project icestorm](http://www.clifford.at/icestorm/)
to target the [TinyFPGA Bx board](https://tinyfpga.com/). As of this writing,
the CPU has most of the functionality in the specs. It is able to do basic
animations with LEDs (See [knight.ucisc](../examples/knight.ucisc)) and compute
[fibonacci numbers](../examples/fib.ucisc). The contents of pc and r1 are sent
to output pins you can hook up to LEDs to see the results (See
[reference/top.v](reference/top.v) for the pinout - don't forget current
limiting resistors on your LEDs).

#### Basic Architecture

The CPU is a 4 stage reference implementation broken up into the following
stages:

* Stage 0: Store result, load next instruction from memory
* Stage 1: Load immediate from memory
* Stage 2: Load source value from memory, or compute source from registers
* Stage 3: Load dest value from memory, or compute dest from registers

It has all 6 general purpose registers.

Not implemented:
* Manipulating the flags register other than as a result of the ALU calculation.
* Manipulating the interrupt register
* Memory banking for devices

#### Usage

1. Follow the [project icestorm](http://www.clifford.at/icestorm/) install
instructions for the yosys toolchain.

2. Compile your uCISC code to reference/prog.hex. You can use the
[uCISC kotlin emulator](https://github.com/grokthis/ucisc-kotlin) for this. If
you follow the install instructions there, you can run something like
`ucisc -c <path-to-your-code> > reference/prog.hex` to compile your code.

3. From the reference directory, do `make clean` and `make prog` with the
TinyFPGA Bx board plugged into your computer. Make sure the bootloader on the
FPGA is running. You may need `make sudo-prog` if your user doesn't have
permissions to connect to the board.

Note: You will want to hook a reset button to pin 2 with a pull down resistor.
Setting reset to 3.3v will start the program over. The bootstrap process isn't
as clean as I would like and sometimes program start doesn't work immediately
after bootstrap and you need to reset it once. I should be able to work this
out in a future release.

#### Tests

Install [Icarus Verilog](http://iverilog.icarus.com/) to run the tests. Then, from
the reference directory, run `./run_tests` to run all tests. Note any warnings,
compile errors or test failures in the output.

#### Status

Currently, the implementation takes about 25-35% of the resources on the
Bx board, not including BRAMs of which I use all 32. I have not included
devices or division in the ALU yet. In theory there is room for 2 CPUs with
4k x 16 block memory each and a couple of simple devices if you leave out
the division opcode.

I have successfully run at a full 16MHz without division and 8MHz with, but
I may not have hit the critical long path yet with my limited testing. Time
will tell, though fib(24) does push a lot of additions and memory operations
through the processor.