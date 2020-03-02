# Micro CISC (uCISC)

uCISC is an opinionated micro instruction set, attempting to
build a computer that is understandable by a single human
from the CPU micro-architecture to the high level programming
language.

This is the main repo for uCISC development. This project will
hold all the uCISC code developed for the instruction set,
compiler, core libraries and OS. In addition to this project
you will need a VM (or yet-to-be-built hardware that runs the
ISA).

VMs:

* [ucisc-ruby](https://github.com/grokthis/ucisc-ruby) - A
  prototype compiler and VM written in ruby.
* ucisc-subx - A much faster, more powerful VM written in
  subx. *Coming Soon-ish*

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

See the [ucisc-ruby](https://github.com/grokthis/ucisc-ruby#usage)
documentation for more details on how the compiler and VM work.

### Helpful Extras

You will find helpful extras like syntax highlighting in the
extras folder. At the moment we have a vim syntax setup. Feel
free to send config files for other editors my way. 

### Examples

The examples folder contains standalone examples of ucisc code
that you can run using the VM.

### Documentation

You'll find helpful documentation in the docs directory. Here
is an index of what you'll find there:

* Instruction Set Architecture
  * [01 - Introduction](/docs/01_Introduction.md)
  * [02 - uCISC Syntax](/docs/02_Syntax.md)
  * [05 - Copy Instructions](/docs/05_Copy.md)
  * [04 - Move Instructions](/docs/06_Move.md)
  * [07 - ALU Instructions](/docs/07_ALU.md)
  * [08 - Page Address Instructions](/docs/08_Page.md)
  * [09 - Control Instructions](/docs/09_Control.md)
  * [10 - Other Information](/docs/10_Other.md)

## Background Information

This is part of a larger project I'm working on to build a
scalable computer system understandable by a single human
mind. Here is some background on the various pieces of the
project and it's history:

#### Level 1 - The micro-architecture

Coming soon. I have designed and implemented several versions
of the CPU at the logic level on previous instruction sets,
but once I got things working I found it almost impossible
to write code well with those instructions. It was time well
spent, however, since I gained an eye for how easy or hard
something is to implement in hardware.

I will come back to this level once the ISA is solid.

#### Level 2 - The instruction set architecture

This is where I'm spending a good bit of my time now.
Originally, I was aiming for a 8-bit RISC like ISA that was
powerful (as far as 8-bit goes) and easy to understand.
Eventually I moved to 16-bits and a CISC architecture. In the
end, I think it will prove way more effective at accomplishing
my goals.

#### Level 3 - uCISC the programmer friendly language based on SubX

Around the time I was moving toward 16-bit CISC, I ran across
[a blog post on Mu](http://akkartik.name/post/mu-2019-1). The
rest, as they say, is history.

Micro CISC was even named by Kartik (Unashamedly CISC or uCISC)
which natrually became MicroCISC and a reference to Mu.

#### Level 4 - MuCISC, a memory-safe language based on Mu

Not started. uCISC and MuCISC are "spiritual forks" of SubX
and Mu since different architectures lead to slightly
different tradeoffs. You can see the beginnings of Mu's logic
in [another blog post](http://akkartik.name/post/mu-2019-2).

#### Level 5 - TBD, an expressive high-level language

By this point, you should get the idea that Levels 3 - 5
borrow heavily from Kartik's exploration of single mind
computing and leaky abstractions. I don't have much info
on this level yet.
