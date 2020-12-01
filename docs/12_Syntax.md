# uCISC Syntax

* Back: [The Path Not Taken](11_The_Path_Not_Taken.md)
* Next: [Beginning to Code](12_Beginning_To_Code.md)

uCISC syntax is heavily influenced by [Mu](https://github.com/akkartik/mu/blob/main/mu.md).

If you use vim, there is syntax highlighting available:
[extras/ucisc.vim](/extras/ucisc.vim). 

### Fundamentals

uCISC is intended to almost exclusively use memory to store values rather than
registers. This is the most important and fundamental difference from most other
assembly languages. Rather than loading values into registers, operating on them
and putting them back into memory, uCISC operates on data in memory and puts it
back into memory.

The most important register is the stack register. It points to the call stack
which is usually where you will store variables and data in common use. The stack
register can be any one of the general purpose registers you decide to use, but
it should be consistent throughout your application.

#### Comments

Line comments start with a '#' character and run all the way to the newline. Before
processing a statement, the compiler will strip all comments from the line. Comments
are used in code examples throughout this documentation to provide helpful tips.

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
 - banking
 - interrupt

R1-R6 are the general purpose registers, the others are special in some way. PC is
the program counter, pointing to the current instruction. Val is a special register
that always has the value of 0 (this is important for immediates). The banking
register controls the banking mechanics of the general purpose registers.

Since we are pointing at memory, registers have offsets. For example, `r1/2`
indicates the value in memory at address r1 + 2. Since the `val` register is always
zero, it is useful for constants, such as `val/42`.

The first step in a program is usually to give registers useful names to refer to
them by. You do this by defining a new name for the register. The most common is
the stack. For example:

```
def stack/r1

# Optionally, you can initialize r1 at the same time
# This is what you will usually see
def stack/r1 <- copy val/0
```

The initialization is not required, but usually good practice. After this definition
you can use stack instead of r1. Once you have a definition, you can create variables
on that definition:

```
var stack.value/0 push <- copy val/5
var stack.answer/0 push <- copy val/42
```

Important: Notice the "push" at the end of this variable declaration. The causes
uCISC processor to push the value to the top of the stack, effectively creating
space for the variable. If you push a variable, you must initialize it with a
statement (more on statements later).

If you are mapping definitions and variables to existing data, the statement is
optional. For example, you could initialize r3 with a player pointer and define
a variable mapping for it.

```
def player/r3 <- copy stack.playerAddress
var player.energy/0
var player.shields/1
var player.torpedoes/2
var player.nextPlayer/3

# Or you can still use a statement if it is useful
var player.counter/4 <- add val/1
```

This creates a set of variables that view the memory starting at r3 a certain way.
In this case, the energy value is at offset 0, shields at offset 1 and and so on.
uCISC keeps track of these and you can use them any time they are in scope.
Make sure `r3` points to a valid player object in memory, however. Conveniently,
the definitiions stick around if you change r3:

```
&player = player.nextPlayer

# The player variables are still valid
# For example, this still works but on the new player
player.counter <- add val/1
```

Definitions don't have to be unique and don't override each other.

```
def vector/r3
var vector.angle/0
var vector.magnitude/1
```

I now have another definition based on `r3` that I can use. I can use either vector
or player any time they are in scope, as long as I make sure `r3` points to a valid
address with the data I expect to be there. It is just much simpler to remember
vector.angle rather than r3/0 in this context and know what that refers to.
At the moment, you can not define variables independent of a register, but that
ability might be added in the future.

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

- `<-` - always store
- `<~` - never store, just set the flags as if you had after computing the result
- `<0?` - store if zero
- `<!?` - store if not zero
- `<n?` - store if negative
- `<p?` - store if positive
- `<o?` - store on overflow
- `<i?` - store if processor is currently interrupted

These effects always operate on the result of the last ALU operation that wasn't a
copy. A copy will never alter the flags. This is useful for testing multiple
conditions on a previous statement. For example:

```
{
    stack.myAngle <~ sub lookup.angle
    pc <0? copy pc/break # <---- This copy doesn't modify the flags
    pc <n? copy pc/break # <---- This copy acts on the flags from the subtract

    &lookup <- copy &lookup.value
    pc <- copy pc/loop
}
```

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
* `sub` - Subtract, respect signed mode
* `mult` - Multiply, respect signed mode, carry is zero
* `multu` - Multiply, return MSW, respect signed mode, carry is zero
* `addc` - Add with carry in, respect signed mode

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
and the arguments last. By convention, registers 2 and 3 are volatile and can be
overwritten by functions at will without saving and restoring state. Registers 4-6
should be saved and restored by functions. R1 is typically the stack register.

Given those conventions, a full listing of the `readChar` method might be:

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

* Back: [The Path Not Taken](11_The_Path_Not_Taken.md)
* Next: [Beginning to Code](12_Beginning_To_Code.md)

