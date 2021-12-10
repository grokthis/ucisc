# uCISC Hardware Implementation

The Micro80 is the uCISC computer designed around a ECP5 FPGA that allocates
80k of memory to the video controller. This allows up to 320x240 bitmap mode
at 16-bit color, 640x480 bitmap mode with 4-bit color and a wide variety of
character modes at 16-bit color.

The FPGA resources are allocated to four key components:

* CPU with 64k of memory
* I/O controller with up to 48k of memory
* Video controller with 80k of memory
* Memory controller with 16k of memory allocated

All I/0 devices (game ports, keyboard, UART, SD card, etc.) except video ports
are attached to the I/O controller. The video controller produces a 16-bit
color output at a specific VGA pixel frequency (depending on the resolution)
that can be hooked up to a DAC for conversion to a VGA signal. The memory
controller is connected to a large memory address space that is implementation
specific. The CPU, I/O controller and memory controller are all implemented as
uCISC processors. The video controller presents as a directly addressable
memory space to the CPU and uses PLLs to output the correct pixel clock for
the video mode selected.

This implementation requires an ECP5 to run as it needs 208k words of memory
and true dual port memory access for the processors. This code uses
[project trellis](https://github.com/YosysHQ/prjtrellis) to target the ECP5
development board.

#### Basic Architecture

The CPU is a 4 stage implementation broken up into the following stages:

* Stage 1: Store result, load the next instruction and immediate from memory
* Stage 2: Load source value from memory, or compute it from registers
* Stage 3: Load destination value from memory, or compute it from registers
* Stage 4: Compute the result

Most of the computer is implemented on the ECP5 FPGA with the exception of
the 16-bit DAC and RAM (which is supplied by an external chip).

```
#        ##############   ##############
#        #   Video    #   #            #
#        # Controller #===# 16-bit DAC #===|- VGA Port
#        #            #   #            #
#        ##############   ##############
#              ||
#              ||
#        ##############   ##############   |- Game Ports
#        #            #   #    I/O     #   |- PS/2 Keyboard
#        #    CPU     #===# Controller #===|- SD Card
#        #            #   #            #   |- UART Port
#        ##############   ##############   |- ROM
#              ||
#              ||
#        ##############   ##############
#        #   Memory   #   #            #
#        # Controller #===#    RAM     #
#        #            #   #            #
#        ##############   ##############
```

#### Usage

TODO: fix this section

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

TODO: Update this section

Currently, the implementation takes about 25-35% of the resources on the
Bx board, not including BRAMs of which I use all 32. I have not included
devices or division in the ALU yet. In theory there is room for 2 CPUs with
4k x 16 block memory each and a couple of simple devices if you leave out
the division opcode.

I have successfully run at a full 16MHz without division and 8MHz with, but
I may not have hit the critical long path yet with my limited testing. Time
will tell, though fib(24) does push a lot of additions and memory operations
through the processor.
