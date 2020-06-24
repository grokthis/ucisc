# Micro CISC (uCISC)

uCISC is an opinionated micro instruction set, attempting to build a computer that
is understandable by a single human from the CPU micro-architecture to the high
level programming language.

Quick link to the docs: [Introduction](/docs/01_Introduction.md)

This is the main repo for uCISC development. This project will hold all the uCISC
code developed for the instruction set, compiler, core libraries and OS. In addition
to this project you will need a VM (or yet-to-be-built hardware that runs the ISA).

VMs:

* [ucisc-ruby](https://github.com/grokthis/ucisc-ruby) - A prototype compiler and VM
written in ruby.

## Getting Started

### Get a VM

Prerequisites:
 * Ruby 2.5+
 * Bundler 2

Install the ucisc command by:

```
$ git clone https://github.com/grokthis/ucisc-ruby
$ cd ucisc-ruby
$ bundle install
$ bundle exec rake install
```

You can now do something like:

```
ucisc <file.ucisc>
```

Or you can just run directly in the repo if you don't want to install the gem:

```
$ cd ucisc-ruby
$ exe/ucisc examples/fib.ucisc
```

See the [ucisc-ruby](https://github.com/grokthis/ucisc-ruby#usage)
documentation for more details on how the compiler and VM work.

### Helpful Extras

You will find helpful extras like syntax highlighting in the extras folder. At the
moment we have a vim syntax setup. Feel free to send config files for other editors
my way. 

### Examples

The examples folder contains standalone examples of ucisc code that you can run
using the VM.

### Documentation

You'll find helpful documentation in the docs directory. The documentation is
constantly improving. Start at the [Introduction](/docs/01_Introduction.md) for more
information.

