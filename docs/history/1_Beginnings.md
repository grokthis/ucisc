## uCISC History and Motivation

1. [Beginnings](1_Beginnings.md)
2. [Transparent Computing](2_Transparent_Computing.md)
3. [Rough Hardware Outline](3_Rough_Hardware_Outline.md)

## Beginnings

Ever since I got started on my [VTech PreComputer 1000](https://www.youtube.com/watch?v=_Wm8OtLWJys&t=1s)
I have wanted to build my own computer. In the early 90s, I got into electronics,
experimenting with 555 timers, some simple digital gates and things like that. This
only made me want to do it more, but it always felt out of reach. The 8080 processors
I was reading about were mystifying-ly complex.

College rolled around, and you can bet I took some digital courses. I did some MIPS
assembly and introductory digital design where we built and simulated a simple CPU.
It always left me wanting more.

A few years ago, I decided that I should just do it. I came across the open source
FPGA stack built around yosys, and plenty of YouTube videos on the subject.

As I started working on it, I started trying to design an 8-bit instruction set,
mostly because I wanted to manage the complexity and have something I could actually
build in hardware at some point. It's possible to build an 8-bit computer from logic
gates with some time, money and effort. The most prominent example of this is
[Ben Eater](https://www.youtube.com/playlist?list=PLowKtXNTBypGqImE405J2565dvjafglHU).

However, as I worked through the details and went through the inevitable false starts
trying to work through my own instruction set design, I was inspired by the
Commodore 64, the best-selling computer of all time. The 8-Bit Guy has a great series
on [Commodore history](https://www.youtube.com/playlist?list=PLfABUWdDse7Y6LLPlfsHKcvBCgqaudzVY)
on YouTube. Computers of that era were incredibly accessible to users. The machine
was at your disposal and with some effort you could learn how to control the
registers that affected the machines' behavior. The 8-Bit Guy also did a video
describing his [ideal computer](https://www.youtube.com/watch?v=ayh0qebfD2g), and he
describes the essence of the greatness of computers of that era:

> Today when I think back to the computers that I was most fond of, it wasn't the
> Amiga. despite being much more powerful, the Amiga's operating system put a layer
> between the hardware and the end user.
>
> I still have a fondness for the 8-bit computers, and I don't have a particular
> favorite. I love writing code, and I'm just as happy to code on a VIC-20 a C64 or
> a Plus4 as long as it has that same closeness to the hardware.
>
> The 8-Bit Guy

Then I remembered the computer that I first learned on. I learned to program basic on
a 1 line LCD readout and I LOVED it! Nostalgia Nerd has a great video about the
[VTech computer toy line up](https://www.youtube.com/watch?v=9F4it_DH6ps) that brings
back those memories in all their annoying PC speaker tonal glory.

This is what I want to recapture in the uCISC homebrew computers I build, that
close-to-the-hardware, transparent computer that rewards curiosity and discovery and
leads you to a deeper understanding of how computers work.
