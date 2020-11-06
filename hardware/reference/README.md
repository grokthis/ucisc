The reference design is implemented for an TinyFPGA Bx board. You'll need to
install project icestorm to get it to build and run. It should be easily
adaptable to other projects (e.g. the ECP5) in the yosys family.

## Getting Started
Step 1: Install [icestorm](http://www.clifford.at/icestorm/) for your computer:

```sh
sudo apt-get install build-essential clang bison flex libreadline-dev \
                     gawk tcl-dev libffi-dev git mercurial graphviz   \
                     xdot pkg-config python python3 libftdi-dev

mkdir icestorm-build
cd icestorm-build

git clone https://github.com/cliffordwolf/icestorm.git icestorm
cd icestorm
make -j$(nproc)
sudo make install
cd ..


git clone https://github.com/cseed/arachne-pnr.git arachne-pnr
cd arachne-pnr
make -j$(nproc)
sudo make install
cd ..

git clone https://github.com/cliffordwolf/yosys.git yosys
cd yosys
make -j$(nproc)
sudo make install
cd ..

pip install --user tinyprog
```

Step 2: Clone the project and build it

```shell
git clone git@github.com:grokthis/ucisc.git
cd ucisc/hardware/reference
make
```

Program the TinyFPGA B-series board with the bitstream:
```shell
make prog
```

## Testing

I am working on building out a pretty extensive set of tests. You can
run all of the tests by doing:

```
cd ucisc/hardware/reference
./run_tests
```