## Chips and Hardware

Creating a custom CPU architecture means we have to one of a few things:

* Emulation/Virtualization
* Custom silicon
* Verilog + FPGA

Emulation and virtualization is necessary for development. However, I don't want to
use emulation on the real hardware. There are some amazing things in the world that
are virtualized and/or emulated. [RetroPie](https://retropie.org.uk/) for example
uses emulation to run old games on new hardware. The entire cloud is an exercise in
virtualization. You can do amazing things with it. However, you are still left with
dealing with layers of software and hardware and that violates the core tenet of
building a transparent computer.

With emulation off the table, it would be awesome if we could build custom silicon
chips for this. Hey, I can dream can't I? Those cost hundreds of thousands of
dollars, even on the cheap. So that's not really possible. The first step of doing
that anyway is building the hardware in an FPGA, which brings us to the final option.

#### Let's use an FPGA!

FPGA's have exploded in recent years and there are tons of options.
[SymbiFlow](https://symbiflow.github.io/) is a fully open source tool chain that has
varying support for a variety of FPGA's, most of which are built by Lattice
Semiconductor. They seem to be open source friendly even though they sell a very
expensive proprietary toolchain for their FPGA's.

The best supported FPGA families so far are Project IceStorm (iCE40 FPGA's) and
Project Trellis (ECP5 FPGA's). Both iCE40 and ECP5 are Lattice products, but the ECP5
family is a bit larger with more features, though both families are on the small end
of the market. According to Mouser, the ECP5 price ranges from roughly 5.5 USD (12k
LUT) to 65 USD (84k LUT) depending on the speed rating and I/O pin count.

In order to figure out how to build a computer out of them, we need to know the basic
specs of the ECP5. They come in the following varieties (with key hardware specs):

| Device                | LFE5U12 | LFE5U25 | LFE5U45 | LFE5U85 |
|:--------------------- |:------- |:------- |:------- |:------- |
| LUTs (K)              | 12      | 24      | 44      | 84      |
| sysMEM Blocks (18 Kb) | 32      | 56      | 108     | 208     |
| 18 X 18 Multipliers   | 28      | 28      | 72      | 156     |

The FPGA's have different price points for different speeds. There are complex
timing tables for different elements, but most fall in the 200MHz to 400MHz range
for the ECP5. The development board I have, has a 200MHz clock, so we will use that
as our base clock frequency.

The bread and butter here is the number of lookup tables (LUTs). This is how FPGA's
encode logic, typically through 4-input LUTs. You can code what the output of the LUT
should be for any combination of the 4 inputs. This means that by combining the
inputs and outputs of multiple LUTs, you can build any logic circuit you want,
limited only by the LUT count required.

Some things that affect LUT count:

* Total logic described. More logic means more gates.
* The complexity of the operation (e.g. addition vs multiplication)
* Bit width. Adding 16-bit numbers takes more gates than 8-bit numbers. This effect
can be roughly linear (e.g. for an adder) but can sometimes be exponential (e.g. for
multiplication) in the number of gates required.
* The speed of the operation. Faster means more gates computing in parallel, often
adding specialized duplicate logic for speed (e.g. fast look-ahead adders).
* The ability of the developer to describe the logic minimally based on how complex
the hardware needs to be.
* The ability of the toolchain to find ways of reducing the LUTs needed.
* The ability of the toolchain to lay out the LUTs and connections optimally.

Some of these optimization problems are really hard. I've seen small tweaks to the
input verilog have vast effects on what the toolchain produced. Consequently, it's
impossible to know how many LUT's will be needed for the final design. However, we
can conclude a few things early.

#### CPU Design Limitations

*16-bit CPU*

Smaller is better in terms of the number of gates used. Using 32-bits means vastly
more LUTs and fewer cores and functionality per ECP5. From experience, I was unable
to build an 8-bit instruction set that meets what I feel are the minimums, but a
16-bit ISA proved doable.

Using 16-bits works well on the ECP5 FPGA. Since everything is 16-bits, including
memory addresses and address space. Using word boundaries like this means every
memory location is automatically the right width for an instruction. It also means
we are limited to 64k words for the CPU. Due to how the instruction set shakes out
at 16-bits, which I won't go into just yet, the natural block size is 256 words. This
holds 4x 64 word pages. This mean each sysMEM FPGA block will hold 4 CPU memory
blocks. This gives us a natural sliding scale for memory and core combinations
depending on where our LUT counts and core counts fall out.

For most metrics, the resource counts scale roughly linearly for various device
sizes. The table below roughly outlines where the max ratios fall out for the
84k device. The device blocks specify how many devices can be supported that need
block mem support on the CPU (256 words each).

| Core count | LUTs/core | Words/CPU | Blocks/core | Leftover Blocks |
|:---------------------- |:--------- |:----------- |:--------------- |
| 3          | 28k       | 64k       | 256         | 64              |
| 6          | 14k       | 32k       | 128         | 64              |
| 16         | 5.2k      | 13k       | 52          | 0               |
| 52         | 1.6k      | 4k        | 16          | 0               |

The numbers don't have to line up that perfectly, of course, and we can shift some
blocks to device blocks even if there are no leftovers. It's also possible to
configure a variable number of blocks per core.

I have some 80% complete cores of older specs that clock in at roughly 1.7k LUTs. I
would expect the LUT count for the final spec to be in the 2.5k range. That means we
are closer to 32 cores on an 84k device (4-5 on a 12k)

*Minimize ALU Operations*

Some early hardware descriptions had a lot of ALU operations. The ALU ended up being
massive in terms of the number of LUTs required. In the end, I have limited myself
to 16 carefully chosen operations that pack the most punch. The ALU went through many
iterations, but the final set is heavily influenced by what you actually need to
program the computer effectively.

The operations below seek to strike a balance between LUT counts, expressiveness of
the uCISC ALU options and requiring functions for more complex mathematics and
using up bits in the 16-bit instruction for ALU op codes.

* Copy (1) - For moving data around

* Bit operations (6) - AND, OR, XOR, 2's complement inverse, shift left, shift right -
these are relatively inexpensive and useful. These are pretty standard in ALU's for
good reason. Bit masking operations are common and useful at the hardware level.

* Byte operations (3) - swap MSB and LSB, zero LSB and zero MSB. These are relatively
inexpensive and fall out of the decision to use memory addressed on a 16-bit word
boundary rather than a byte boundary. Word boundaries have all kinds of simplifying
effects in a 16-bit CPU architecture, but we still need to be able to manipulate
byte streams. These 3 can be arbitrarily combined to quickly shuffle bytes around.

* Arithmetic operations (5) - Add, add with carry, subtract, multiply, multiply
significant word. Add and subtract are relatively cheap and important. Multiply is
cheap thanks to being provided by the ECP5 as discrete elements (18 x 18 multipliers
are conveniently more than the 16 x 16 multipliers I need).

* Reserved for the future (1) - I'm just leaving the last slot empty for now. None
of the options for the remaining opcodes were valuable enough to give this up.


Division and floating point operations are missing. They are just too expensive and
will need to be handled by coprocessing units if I really need them in hardware.
