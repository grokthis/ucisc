## Transparent Computing

* Back: [Beginnings](02_Beginnings.md)
* Next: [Chips and Hardware](04_Chips_And_Hardware.md)

What kind of computing device should I aim for? The first steps will be small and
the first runnable version is unlikely to be more than a simple microcontroller.
I won't feel like I'm done, though, until I build an actual computer. Basically,
I want to build something worthy of being used on a regular basis. I want to be able
to use things like:

#### Basic Goals

* Send and recieve email
* Display images
* Edit and organize text files and documents
* Interact with chat services
* Read and post on open social services like Mastadon
* Build applications that consume web APIs (e.g. calendar)
* Build local applications that use structured data (e.g. manage todos or finances)

*Side note:* The thing I need to pull off many of the things above that fuels my
nightmares when I think about writing the code for it is building an SSL layer.
Simple things like posting to Mastadon will require it to function, but I'm not
looking forward to that.

#### Rejected Ideas

At some level, all parts of this project are an act of reimagination. However, the
following I'm certain I don't want to look anything like they do now on modern
desktop computers if I support them at all. We are focused more on curiosity and
discovery, allowing users to fully understand how things work on these homebrew
computers. The modern web stack doesn't seem to fit in that paradigm.

*The World Wide Web*

I have no intention of building a JavaScript interpreter and supporting HTML and CSS
layout methods is wholly outside the realm of possibility when building a personal
computer like this that is comprehensible to it's users. Vast amounts of what is
available doesn't need to be accessible, of course. However, I could imagine building
some targeted apps for things I commonly use (e.g. a Wikipedia app) that consumed
the web on an API basis.

It also might be able to build some basic web like protocols that really can be
supported on a simple computer like this. This type of web will look more like the
early days of the web, except nicer looking (I hope).

*Modern Programming*

I can imagine supporting some additional high level languages on this machine. There
is a niche for things like bash scripts that I'm not sure I want to do away with.
The jury is still out on that (though it obviously won't actually be bash, just
something that fills the scripting role for the OS). However, I'm certain I don't
want to pull in ruby, python or even C/C++ into this world.

It's certainly possible to write it, but as soon as you do, you leave transparency
behind and accept a very, very large environment between you and the computer. This
essentially means I'm not targeting compatibility with anything that exists except
through external communication protocols like HTTP and simple data formats like Json.

## The Minimum Specs

*1. Proprietary software is not transparent*

The computer, as much as possible should be fully open source. The hardware design,
operating system, development toolchain and more should all be open source. Where
not possible, such as required proprietary hardware components (e.g. FPGA's), they
should be readily available and open source tooling.

*2. Visually Appealing*

I'd like it to display resonably appealing graphics. I enjoy retro-graphics looks
but I don't what that to be a turn-off to others who don't get the same nostalgia
from 8-bit graphics that I do.

It should be able to drive something close to a 1080p display at 60Hz in "desktop"
mode to feel highly responsive and easy on the eyes. The hardware doesn't need to be
able to push high end graphics at 60 FPS but should be able to render a graphical
desktop at 60 FPS. It should use modern display connectors such as HDMI or
DisplayPort at 8-bit color depth (a total of 3 bytes per pixel).

*3. Modern Memory*

If I want to drive 1080p displays, display high resolution images with 24-bits
of color and build data driven applications, it needs to be able to make use of
modern memory sizes. I imagine it will be more memory effecient that a full modern
OS and it won't run the most memory intensive application of today's computers (the
web browser) but being able to utilize several GB of memory will still be necessary.

*4. Common Disk Storage Media*

Loading programs on/off of old computers was tedious. We have modern storage
solutions that are portable and easy to use (SD cards, USB drives, etc.). The system
will need to be able to make use of them. It is NOT required to use a compatible
file system, but tooling will be required so that modern computers can load data
to/from whatever file system is used if not commonly readable like FAT.

*5. Network Connected*

You can't build applications that interact with the web without a network. So, the
computer needs to be able to connect to the internet. WiFi would be nice for
portability, but it's not a hard requirement.

## Transparent Software

Thinking back to when I was a kid and I fell in love with computers and programming,
I was always trying to peel back the layers and understand how things really worked.
I would love to build a computer that faciliated that as much as possible. Imagine
booting into a GUI based operating system. You have the clock, text editor and
perhaps a few games readily available. After playing around with the interface for
a bit, you want to know how it works, so you read the docs on how to write a basic
program.

It turns out to be really easy to get into. Before long, you are writing hello world,
maybe even getting a few windows showing on the screen. Then you say, well, I wonder
how the text editor works. So you launch the text editor, flip into debug mode and
type a key. You immediately drop into a debug window and are stepping through the
code for the text editor. It's clean, well written and really not that complicated.
You borrow a few things and add them to your hello world program. Suddenly you turn
your hello world program into a simple text editor.

Then you think, "I wonder how to load and save these files I'm making." This time,
instead of debugging the text editor, you decide to just open the executable file.
And there it is, the human readable code that was written for the text editor. You
read it through, noting several new techniques you discover, but you find the
load/save code.

Then it hits you, "wait, are all the executables like this?" It turns out they are.
Before long, you are reading the code for the browser, the game you were playing
earlier, and even the operating system itself. Then you stumble across the startup
code for the OS. You see how it's loading everything and then running the programs
for the UI startup. Now you're modifying the OS itself, adding steps to the startup,
changing the UI appearance, customizing your background, even adding a timed
background rotation.

Ideally, the computer would:

* Blur the lines between source code and executables
* Allow you to modify the system on the fly
* Recover and undo when you inevitably mess something up
* Debug any program live without only looking at hex dumps or binary bits
* Immediately start writing code that runs
* Slowly introduce how the OS and even the hardware works as you dive in
* Understand how the entire computer works by using the computer

#### Continue Reading

* Back: [Beginnings](02_Beginnings.md)
* Next: [Chips and Hardware](04_Chips_And_Hardware.md)

