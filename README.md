NOTE: I am in the middle of refactoring the project structure to support multiple
CPU versions with different purposes. Also, I am moving the build system to
gradle. Documentation coming soon.

# Micro CISC (uCISC)

uCISC is an opinionated micro instruction set, attempting to build a computer that
is understandable by a single human from the CPU micro-architecture to the high
level programming language.

Quick link to the docs: [Overview](/docs/Overview.md)

This is the main repo for uCISC development. This project will hold all the uCISC
code developed for the instruction set, compiler, core libraries and OS. In addition
to this project you will need an emulator or FPGA hardware to run the code.

Virtual Execution:

* [ucisc-kotlin](https://github.com/grokthis/ucisc-kotlin) - An emulator and debugger
* OUTDATED - [ucisc-ruby](https://github.com/grokthis/ucisc-ruby) - A prototype compiler and VM
written in ruby.

## Getting Started

### Install the Emulater

See the [ucisc-kotlin](https://github.com/grokthis/ucisc-kotlin) readme for
instructions on how to install the emulator. It is recommended to install it such
that you have `ucisc` as a command on the terminal.

### Compiling and Running uCISC

To compile uCISC code:

```
ucisc -c examples/fib.ucisc
```

To run uCISC code on the emulator

```
ucisc examples/knight.ucisc
```

If you need the UART port hooked up, you'll need to do a little more
work. First, you'll need to create a virtual 2-way serial port. On linux you
can do that with `socat`

```
> sudo apt install socat
> socat -d -d pty,raw,echo=0 pty,raw,echo=0
```
This will print out something like:

```
2020/11/30 19:06:28 socat[9758] N PTY is /dev/pts/4
2020/11/30 19:06:28 socat[9758] N PTY is /dev/pts/5
2020/11/30 19:06:28 socat[9758] N starting data transfer loop with FDs [5,5] and [7,7]
```

The two PTY devices are the two ends of your serial port. Hook one end to
your serial application (I recommend GTKTerm) and the other to the uCISC
emulator:

```
ucisc -r=/dev/pts/4 -t=/dev/pts/4 examples/space.ucisc
```

### Running the Hardware

You can find details about the reference implementation on the hardware in its
own [README](hardware/README.md).

### Helpful Extras

You will find helpful extras like syntax highlighting in the extras folder. At the
moment we have a vim syntax setup. Feel free to send config files for other editors
my way. 

### Examples

The examples folder contains standalone examples of ucisc code that you can run
using the VM.

### Documentation

You'll find helpful documentation in the docs directory. The documentation is
constantly improving. Start at the [Introduction](/docs/Overview.md) for more
information.

