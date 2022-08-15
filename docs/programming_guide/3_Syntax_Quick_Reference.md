## uCISC Programming Guide

1. [Getting Started](1.0_Getting_Started.md)
   1. [Configuring IntelliJ](1.1_Configuring_IntelliJ.md)
   2. [Configuring VIM](1.2_Configuring_VIM.md)
   3. [Compiling uCISC Code](1.3_Compiling_uCISC_Code.md)
   4. [Using the uCISC Simulator](1.4_Simulating_uCISC.md)
   5. [Using the uCISC Soft Core](1.5_Running_uCISC_Soft_Core.md)
2. [Introduction to Programming With uCISC](2.0_Program_With_uCISC.md)
   1. [Accessing External Devices](2.1_Accessing_Devices.md)
   2. [Common Devices](2.2.0_Common_Devices.md)
      1. [GPIO](2.2.1_GPIO_Devices.md)
      2. [I2C](2.2.2_I2C_Devices.md)
      3. [UART](2.2.3_UART_Devices.md)
      4. [Video Devices](2.2.4_Video_Devices.md)
   3. [Advanced uCISC Programming Techiques](2.3_Advanced_Programming_Techniques.md)
3. [uCISC Syntax Quick Reference](3_Syntax_Quick_Reference.md)
4. [Standard Libraries](04_Syntax_Quick_Reference.md)
5. [Instruction Set Details](5_Instruction_Set_Details.md)

# Chapter 3 - uCISC Syntax Quick Reference

uCISC syntax is heavily influenced by [Mu](https://github.com/akkartik/mu/blob/main/mu.md).

### Architecture Distinctives

uCISC is different from RISC:
* Registers are usually pointers to memory
* Data is loaded from, computed and stored back to memory in a single instruction

uCISC is different from CISC:
* uCISC is a tiny instruction set
* uCISC is brutally simple whenever possible

The uCISC name is a bit (intentionally) ironic since it is really neither RISC nor CISC.
It lies on the spectrum somewhere between.

#### Comments

Line comments start with a '#' character and run all the way to the newline.
Comments can be added to the end of a line after code.

```
# This is a comment, hopefully a helpful one
```

#### Registers

The registers are:
 - r1
 - r2
 - r3
 - r4
 - r5
 - r6
 - pc
 - val
 - control/banking
 - interrupt

R1-R6 are the general purpose registers. The `pc` register is the program counter.
Val is "constant" register for directly loading immediate values. The control/banking
register configures the local/device setup for the general purpose registers.

#### Offsets

Registers point to a memory address, but accept an offset to refer to a value several
addresses offset from the register pointer.

The maximum offsets depend on the source and destination.

| Destination | Source Min | Source Max | Destination Min | Destination Max |
|-------------|------------|------------|-----------------|-----------------|
| Memory      | -2048      | 2047       | 0               | 15              |
| Register    | -32768*    | 65535*     | N/A             | N/A             |

**Register destinations:** &r1, &r2, &r3, &r4, &r5, &r6, pc, control, flags, interrupt

**Memory destinations:** r1, r2, r3, r4, r5, r6

Examples:

```
# Copy value at address r2 + 5 to address r1 + 10 (mem destination)
r1/10 <- copy r2/5

# Copy value at address r2 + 200 into r1 (register destination)
&r1 <- copy r2/200

# Copy value at address r2 - 4096 (pc register destination)
pc <- copy r2/-4096
```

Invalid statements
```
# r1 is memory destination (same as r1/0), maximum source is 2047
r1 <- copy r2/2100

# r1 is memory destination, maximum destination is 15
r1/16 <- copy r2/10

# r1 is a register destination, destination offsets are not supported
&r1/16 <- copy r2/10
```

*Technical note: since the local memory address is a 16-bit space, the offset
representation wraps around. Thus, an offset of 65535 (%FFFF) is the same as an offset
of -1. Likewise and offset of -32768 is the same as and offset of 32768. The sign
does not technically matter for a 16-bit offset in a 16-bit address space.

#### Defines

The first step in a program is usually to give registers useful names to refer to
them by. You do this by defining a new name for the register. For example:

```
# Define stack
def stack/r1
# Initialize stack address
&stack <- copy val/0

# Optionally, you can initialize at the same time
def stack/r1 <- copy val/0
```

There is no limit to the number of names a register can have:

```
def stack/r1 <- copy val/0
def objects/r1
def references/r1
def calls/r1

# The stack, objects, references and calls defines will all have the same value
# Any push/pop statements will update the reference count for all defines
```

#### Variables

After this definition you can use stack instead of r1. Once you have a definition,
you can create variables on that definition:

```
def object/r1 <- copy val/1024
# r1 points to address 1024

var object.value/0 <- copy val/100
# object.value refers to address 1024, contains the number 100

var object.answer/1 <- copy val/200
# object.value refers to address 1025, contains the number 200
```

You can push variables during definition to create new space. A push decrements
the register before storing, preserving the value of whatever address it is
currently pointing at.

```
def stack/r1 <- copy val/0
# stack points to address 0

stack <- copy val/5
# address zero has value 5

var stack.value/0 push <- copy val/100

# stack.value refers to address %FFFF, contains the number 100
# address zero still has value 5

var stack.answer/0 push <- copy val/200

# stack.answer refers to address %FFFE, contains the number 200
# address zero still has value 5
# The offset for stack.value was automatically incremented to keep it pointing
# at address %FFFF
```

If you are mapping definitions and variables to existing data, the initialization is
optional. For example, you could initialize r3 with a player pointer and define
a variable mapping for it.

```
def player/r3 <- copy stack.playerAddress
var player.energy/0
var player.shields/1
var player.torpedoes/2
var player.nextPlayer/3

# Or you can still initialize it if it is useful
var player.counter/4 <- add val/1
```

This creates a set of variables that view the memory starting at the address in r3
in a certain way. In this case, the energy value is at offset 0, shields at offset 1
and and so on. uCISC keeps track of these offsets, and you can use them any time they
are in scope. Make sure `r3` points to a valid player object in memory, however.

Conveniently, the definitions stick around if you change r3:
```
&player = player.nextPlayer

# The player variables are still valid
# For example, this still works but sets the value on the next player
player.counter <- add val/1
```

Definitions for a register don't have to be unique and don't override each other.

```
def player/r3 <- copy stack.playerAddress
var player.angle/0
var player.magnitude/1
var player.energy/2
var player.shields/3
var player.torpedoes/4
var player.nextPlayer/5

def vector/r3
var vector.angle/0
var vector.magnitude/1
# player definition is still valid.

def data_table/r3
var data_table.rows/0
var data_table.columns/1
# player, vector definitions are still valid.

# All 3 of these result in exactly the same instruction
stack.angle <- copy player.angle
stack.angle <- copy vector.angle
stack.angle <- copy data_table.rows
```

You can define as many different ways of looking at the data as you want. Just be
careful to keep track of what is what. In the example above, viewing a player
as a `data_table` is probably a bad idea.

#### Scope

Scope in uCISC is tied to code blocks surrounded by curly braces:

```
{
  # Can't use vector here, it's not defined yet

  def vector/r3
  var vector.angle/0
  var vector.magnitude/1

  # We can use vector here all we want
}

# Can't use vector here, it's hidden
```

You can think of blocks as namespaces without the names. This allows you to "hide"
variables, labels and other things in scopes and avoid name collisions across larger
programs. In general, it's best to have very few definitions in the global namespace
and group like functionality together in blocks.

#### Statements

Statements are the backbone of uCISC. Every uCISC statement (but not necessarily
`def` and `var` descriptions) have a destination (a register and offset), effect and
source (register and offset). The destination comes first. The effect is the arrow
(more on that later).

```
r3/0 <- copy r1/1
```

Of course, you'll almost always use the variables you defined instead:

```
vector.angle <- copy stack.angle
```

Of course that reads much better to humans, but notice that these are just variable
definitions the compiler is looking up the register and the offset for.

The effect controls if the statement results in a value being stored. In the example
above copy will always happen. The list of effects is:

* <-, <| (store always)
* <~, <|~ (don't store, just set the flags for the operation result)
* <0, <|0, (store if zero)
* <1, <|1, (store if !zero)
* <+, <p, <|+, <|p (store if positive)
* <n, <|n, <|- (store if negative)
* <&, <|&, <o (store if overflow)
* <#, <|# (store if error)

Side Bar on Readability:
> You can use any of the varieties listed. Some are easier/harder to type,
so it depends on how much effort you want to put in for readability and your preferences
on that. The <| versions look great in a font with ligatures (e.g. Fira Code).
I haven't settled on these and may add text only versions (e.g. neg, pos, err, etc)
in the future. Personally, I think assembly programs are significantly impacted with
poor readability, so I don't mind a bit of extra typing to keep it clean. On the other
hand without ligatures, some of these don't look great, so idk. ¯\_(ツ)_/¯

These effects always operate on the result of the last ALU operation that wasn't a
copy. A copy will never alter the flags. This is useful for testing multiple
conditions on a previous statement. For example:

```
{
    stack.myAngle <~ sub lookup.angle # <---- flags are set
    pc <0? copy pc/break              # <---- This copy doesn't modify the flags
    pc <n? copy pc/break              # <---- This copy uses the flags from the sub

    &lookup <- copy lookup.value
    pc <- copy pc/loop
}
```

#### Arithmetic and Logic

The ALU operations are:

* `copy` - Copy
* `and` - And
* `or` - Or
* `xor` - Xor
* `inv` - Invert
* `shl` - Shift left, zero extend
* `shr` - Shift right, respect signed mode
* `swap` - Swap MSB and LSB bytes
* `msb` - MSB only: A & 0xFF00
* `lsb` - LSB only: A & 0x00FF
* `add` - Add, respect signed mode
* `addc` - Add with carry in, respect signed mode
* `sub` - Subtract, respect signed mode
* `mult` - Multiply, respect signed mode, carry is zero
* `multsw` - Multiply, return MSW, respect signed mode, carry is zero

#### Functions

Functions are fairly powerful convenience tools in uCISC. They are the only statement
that maps to multiple uCISC instructions under the hood. This compromise was made to
solidify a calling convention, save a lot of duplicate typing and help organize code
effectively. Overall, it seemed worth the compromise.

Functions are defined like so:

```
fun stack.readChar(serialDevice) -> char {
    # Code goes here
    
}
```

Functions need to be declared with a stack register that is used as the call stack.
Typically, you'll have defined that as "stack". The name of this function is
"readChar" and it reads a character from a provided serial device and returns it.
Inside the function, you can reference `stack.serialDevice` for the argument and
`stack.char` for the return value.

The calling convention for uCISC is a bit unusual. The caller will reserve space on
the stack for the return variables in advance. The return address will be pushed next
and the arguments last.

A full listing of the `readChar` method might be:

```
# Reads a single character from the serial device
# Relies on register 6 being banked, but preserves
# the current content of register 6
fun stack.readChar(serialDevice) -> char {
    var stack.saveR6/0 push <- copy &r6
    def serial/r6 <- copy stack.serialDevice
    var serial.flags/1
    var serial.rx/4
    var val.readReady/512

    stack.waitForRead()
    stack.char <- copy serial.rx
    serial.rx <- copy val/0 # mark byte as read

    &r6 <- copy stack.saveR6 pop

    pc <- copy stack.return pop

    fun stack.waitForRead() {
        serial.flags <~ and val.readReady
        pc <0? copy pc/loop

        pc <- copy stack.return pop
    }
}
```

A couple things to notice:

1. Scoping rules apply. Functions are required to start a new scope.
2. Functions are just syntax sugar around scopes, so anything you can do in a scope
   block you can do in a function, including `def`, `var` and other `fun`
   declarations

The helper functions and serial definitions are restricted to be visible only within
the function scope.

#### Push and pop

Push and pop are convenient operations for specific situations, but they have
important rules about when they can be applied. Push adds a new variable to the stack
and pop removes variables from the stack. Both operate on the variable at the
specified offset.

For example:

```
pc <- copy stack.return pop
```

Remember `stack.return` is a register and offset that uCISC has kept track of. The
pop will pop to and including the offset. You will see this statement at the end of
most functions because `stack.return` is the reference to the return address and this
statement pops all of the arguments and variables created in the function right up
to and including the return address, while stuffing the return address into the
program counter.

Push works on the variable offset, but stores the result at a new stack/0 location.

For example:
```
var stack.newVar/10 <- add 5
```

This will `add 5` to the value at `stack/10`, decrement stack and store the result
at `stack/0`.

The restrictions:

1. You can only use push and pop on memory variables, never addresses (i.e. `&`)
2. You can only push the destination and pop the source
3. If you can push, you can not pop (because push and pop use the same bit in the
   instruction encoding)

Note: pc, banking and interrupt registers are not memory variable addressable.

Legal:

```
   var stack.myVar/0 push <- copy stack.other
   var stack.myVar/0 push <- copy pc/10
   var stack.myVar/0 push <- copy banking
   var stack.myVar/0 push <- copy interrupt

   pc <- copy stack.return pop
   &r3 <- copy stack.someVar pop
```

NOT legal:
```
   var stack.myVar/0 <- copy stack.other pop # stack.myVar can push, so can't pop

   var &r3/0 push <- copy pc/10 # &r3 is an address not a memory value
   
   var stack.myVar/0 push <- copy stack.other pop # Can/t push and pop
```
#### Data

Data can be added to uCISC files with hex values. The line must begin with a `%`
symbol. Typically you will want to add a label to refer to the beginning of the data 
block so you can use it in your code.

```
my_data:
% 001F 0A48 656C 6C6f 2C20 776F 726C 6421 2075 4349 5343 2069 7320 6865 7265 210A 0A00

# Refer to my data
&r2 <- copy pc/my_data
```

#### Strings

Strings are in quotes. They will be converted inline to string data (a length
followed by ascii data). Strings are packed 1 character per word and are zero
terminated.

For example:
```
# label is optional
hello_world: "\nHello, world! uCISC is here!\n\n"
```

Strings are converted to data inline, so their placement must be appropriate or you
will try to execute the string as code.

#### Continue Reading

* Back: [The Path Not Taken](../history/11_The_Path_Not_Taken.md)
* Next: [Beginning to Code](12_Beginning_To_Code.md)

