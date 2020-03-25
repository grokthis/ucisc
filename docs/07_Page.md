# uCISC Page Instructions

Page instructions control the paging operations to/from local
and main memory.

NOTE: These instructions are a work in progress and have not
been implemented in the simulator yet. They are highly likely
to change rapidly.

#### Local Memory Space

All uCISC processors must execute programs in their local memory
space. Each byte pair in the memory space has one address in the
2-byte address space. Therefore, uCISC processors can address up
to 128kb of local memory. It is not required for the processor to
have that much memory actually available, however. Each page is
512 bytes, 256 2-byte values. The addressable memory in a uCISC
processor must be a multiple of one page.

Some form of main memory is necessary for the processors to load
data and execute programs as the processors are only allowed to
communicate with main memory by paging to/from main memory pages.
There are few limits placed on how the main memory can be mapped,
but it should be done in such a way that the software can discover
the layout (or use a known layout) depending on the application.

#### Main Memory Layout

Since uCISC processors use 2-byte values, they will naturally be
able to address 2^16 pages of 512 bytes (32MiB). However, it's
possible that some implementations may support larger page address
spaces. Any multiple of 2-byte address length can be supported in
theory, but doing so would require software support since it
would change the size of values being stored for memory page
addresses.

One possible memory layout might be the following:

```
0x0000 - 0x000F : Boot ROM
0x0010 - 0x0076 : 102 device mapped pages
0x0077 - 0x0077 : Video device configuration
0x0078 - 0x2FFF : Video memory, up to 1080p x 24bit
0x3000 - 0x30FF : Message queue storage (2^16 messages * 2 bytes)
0x3100 - 0xFFFF : 12544 pages of general access memory (6.125 MB)
```

This layout contains typical elements in a memory space available
to a processor, but the details may vary based on the design needs
of the hardware and the use case. Over time, conventions may be
developed, such as using device 0 as the main board. Typically,
the main board would include the memory mapping layout in a known
format so that the software could detect what resources were
available for use.

#### Devices

Every device or virtual device should be mapped to a device page
in main memory. This includes things like processor cores, I/O
devices, the main system board, etc. Depending on the device,
the mapped page may only accept reads or writes to certain parts
of the mapped page. Reads or writes outside of these zones must
be ignored and not cause errors.

#### Page Visibility & Locking

Each page has 3 settings that control access to the page by a
given processor:

* Assigned PID: this is the PID that "owns" the memory page
* Locked: the processor device ID with a lock on the page
* Content: code or data

By default, all pages are assigned to PID 0 (the root process) and
private (only visible to PID 0). This means that while other
processes can address the page, any interactions with those pages
will fail and set the error flag. That is, any page request, lock,
unlock or anything else should fail in constant time with the
error flag set.

Every processor has a read-only (unless PID == 0) register that
specifies the PID for that process. The processor also has a
device ID that can't be changed while the hardware is powered on.

A process can access a page for another PID when all of these
conditions apply:

* The page has been set to public
* The page is not locked against reads
* The process is trying to read the page

A process can access a page for it's own PID when any of these
conditions apply:

* The page is unlocked
* The current processor owns the lock
* Another processor owns the write lock, but not the read lock

Locks are reentrant to the same processor. When aquiring a read
or write lock and a lock of the other type is already aquired,
the lock type is changed to the requested type. This allows you
to upgrade or downgrade the lock as needed.

Locks should resolve within a single clock cycle when there is no
contention for the locking hardware. However, when under
contention, they are not guaranteed to resolve in a constant
number of cycles. The hardware must not deadlock under any case
and can choose to abort the lock change and return an error to the
processor. The processor must always check the error flag to know
if the lock was successfully acquired or released.

Notes:

Local memory access is specifically designed to simplify
programming for processors that implement this instruction set. It
eliminates the need for contention in the process memory.

Page locks were implemented as a way for processes to run multiple
threads at the same time and coordinate memory access between
them. The goal is to make this as simple as possible while still
providing the necessary capability to coordinate highly parallel
processes. However, this is by far the most complex portion of
the hardware implementation and the details are subject to change
if the hardware proves difficult to understand.

The choice was specifically made to avoid memory virtualization
where each process had a virtual memory space that an MMU unit
mapped to real memory. This would add a lot of complexity and
require a lot of functionality to allow the OS to control the
mapping. It seems that in keeping with single mind principal such
virtualized and invisible memory mappings should be avoided.
However, this opens some security concerns since rogue software
will be able to directly probe the address space for
meta-information about other processes. By asking the OS for
additional pages, it could also map out how much memory is being
used and potentially infer information about memory movements
based on timing information.

The security concerns definitely require more research and thought
to determine if exposing the memory map like this is too
dangerous. Some lines of research on mitigation include:

* Assigning untrusted applications to certain memory ranges they
  are then unable to escalate out of, even if memory is available.
* Constant time responses for rejected access so you can't tell if
  the page doesn't exist, that you don't have access or if it is
  owned by another process.
* Moving devices an other mapped memory to the high address ranges
  and leaving general memory in the low range so you can't tell
  the reserved device memory range.
* Limiting access to public pages in some way. Currently public
  pages are used to transfer data between processes. However,
  there are no process specific permissions for these pages. Any
  process that requests the page can access the page.

#### Instruction Format

```
6/page/ 1.mem/sp local/ 1.mem/sp remote/ 2.imm/sp offset/ 0.dir/page out/ 3.lock/release/

# Packed Instruction Bits:
110NLLPP PKIIIIII
```

*(K) Lock*

Processors will be executing
When performing the copy operation, the instruction can aquire or
release the lock on the page. All the normal lock behaviors apply
(see the control instructions for more information).

* 0 - Release lock (or leave unlocked)
* 1 - Acquire lock (or leave locked)

*(N) Direction*

* 0 - Page out - copy local page to main memory
* 1 - Page in - copy main memory to the local page

*(P) Page in Main Memory*

0. 0.????

1. 1.mem - Value at memory location (r1 + imm)
2. 2.mem - Value at memory location (r2 + imm)
3. 3.mem - Value at memory location (r3 + imm)

4. 4.val - Immediate

5. 1.xmem - 4-byte value at memory location (r1 + imm)
6. 2.xmem - 4-byte value at memory location (r2 + imm)
7. 3.xmem - 4-byte value at memory location (r3 + imm)

Note: xmem arguments are little endian with respect to the word order.
That is location (r1 + imm) is the least significant word and
(r1 + imm + 1) is the most significant word.

*(L) Local Page Identifier*

If the identifier is interpreted as local, then the page identifer
is computed as ((value & 0xFF00) >> 8) where value is the result
of the identifier code.

* 0.blank - Copy all zeros

* 1.mem - Value at memory location r1
* 2.mem - Value at memory location r2
* 3.mem - Value at memory location r3

