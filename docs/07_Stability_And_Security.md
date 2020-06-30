## Stability and Security

* Back: [Memory And Devices](06_Memory_And_Devices.md)
* Next: [Banking Mechanics](08_Banking_Mechanics.md)

TL;DR: I'm largely ignoring security for now. A homebrew just doesn't have as much of
a need for it and I think we can work around most of the issues we will run into in
the short term. You can skip ahead unless you are curious about some potentially
interesting security related effects of the processor architecture.

I have difficulty seeing far enough in advance to know how this computer could be
exploited and I'm no security expert. However, a few things stand out to me:

* A transparent computing environment will find it more difficult to hide malicious
code, if, as I hope, users are encouraged to understand most if not all of the
software on the system.

* There are quite simply fewer failure modes. Pre-emptive computing is abolished and
processes always have CPU affinity for the core they run on. There are just fewer
ways for processes to interfere with each other through back doors because many of
those categories of back doors simply don't exist.

* I intend to avoid as many forms of running untrusted code on the machine as
possible. This amounts to trusting to good behavior of users and programmers, which
definitely makes me nervous, but I hope to generally avoid running untrusted code,
even if sandboxed. I don't want web applications.

So, more or less, at the moment it looks like security will be a bit of a different
game on a homebrew computer than a more modern computer. I wonder if it might be the
only real option if we are truly to make progress on security. Let's simply stop
doing the things that have been shown to lead to exploits that we need to patch all
the time.

Having said that, not all non-malicious code is well behaved, so we do need to make
provision for failure modes where code doesn't work. If we don't, then the computer
won't be conducive to experimentation and iteration. We need to be able to put the
computer into an infinite loop because we failed to terminate the loop properly and
still recover from that mistake.

Fortunately, we already have a possible answer to this. The inter-processor
connections we have already explored in banking are sufficient to allow us to
maintain control. For the simple case, we only need to make sure that processor 0,
the init processor is able to function. Since it is the init processor, it has the
ability to kill and restart any other processor in the system. If there is more than
one init processor all of them must run stable and safe code at all times or you may
lose control of parts of the system.

This also implies a mechanism to build sandboxes within which to run untrusted code.
In addition to the init processors, there will be some number of "device" processors
where all the devices are connected (like memory, I/O, etc.). If we maintain only
stable and trusted code on the device processors as well, we should be reasonably
able to partition which data they access on the system. Untrusted code therefore
runs only on a set of processors to which they have direct control over both the
local memory and code. Each processor is effectively, by it's very design, a sandbox
with very limited connections to either other processors or any hardware. There are
not privileged access modes.

So we are left with the following if you were to run untrusted code on this computer:

* Attacks that rely on speculative execution or processor state will only reveal the
contents of memory the malicious code already has access to: processor local memory.

* Accessing main memory and using cache failures or bit fiddling to read adjacent
bits is unlikely to work since you need to request an entire block from a memory
controller process and copy it to local memory before you can read it. Access will
be rejected at the block level long before any bytes make it anywhere close to the
malicious process.

* Stack overflows in the local processor are limited to corrupting the local memory
and process space. It is probably still possible to craft a request to exploit a bug
in the host OS and gain remote execution, but it is much more difficult with the
banked memory approach. Banked memory is not copied, it is merely accessed. It's hard
to overrun buffers in the host OS this way. Generally banked memory isn't going to
be copied outside the block it and won't ever be executed. Again, it's not impossible
to conceive of an exploit, but I hypothesize they will be much harder to generate and
easier to defend against.

In summary, I hop that even when malicious code is set free on a particular
processor, it will find itself inside a sandbox that is difficult to break out of
with a host OS that can easily kill it and reset it. Any non-malicious but buggy
code can similarily be contained and reset as needed. I wonder if such a system
might actually be more stable and secure than most of our modern system. This is all
in theory, and most of these thoughts are probably a bit premature. In the end, the
right move is probably to move forwards with experimentation rather than get hung
up on the details now. This computer is intended to be both transparent and highly
modifiable, which is likely what is needed to integrate security when it becomes
relevant.

#### Continue Reading

* Back: [Memory And Devices](06_Memory_And_Devices.md)
* Next: [Banking Mechanics](08_Banking_Mechanics.md)

