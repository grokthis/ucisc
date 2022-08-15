## Transparent Computing

The plan is to build a computer worthy of being used on a regular basis. I want to\
be able to use things like:

* Send and receive email
* Display images
* Edit and organize text files and documents
* Interact with simple chat services
* Build applications that consume web APIs (e.g. calendar)
* Build local applications that use structured data (e.g. manage todos or finances)

Our modern computers *hide* all the inner workings of those functions from us. I want
to be able to peel back the layers and see the inner workings behind it.

## The Minimum Specs

*1. Proprietary software is not transparent*

The computer, as much as possible should be fully open source. The hardware design,
operating system, development toolchain and more should all be open source. Where
not possible, such as required proprietary hardware components (e.g. FPGA's), they
should be readily available and open source tooling.

*2. Visually Appealing*

I'd like it to display reasonably appealing graphics. I enjoy retro-graphics looks,
but I don't want that to be a turn-off to others who don't get the same nostalgia
from 8-bit graphics that I do.

It should be able to drive something close to a 1080p display at 60Hz in "desktop"
mode to feel highly responsive and easy on the eyes. The hardware doesn't need to be
able to push high-end graphics at 60 FPS but should be able to render a graphical
desktop at 60 FPS.

*3. Modern Memory*

If I want to drive 1080p displays and display high resolution images memory is going
to be a challenge. A directly addressable 64k for a 16-bit computer won't get me there.

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

Thinking back to when I was a kid when I fell in love with computers and programming,
I was always trying to peel back the layers and understand how things really worked.
I would love to build a computer that facilitated that as much as possible. Imagine
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
And there it is, the human-readable code that was written for the text editor. You
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

Of course, this is all a bit idealistic... we will see where I actually end up.
