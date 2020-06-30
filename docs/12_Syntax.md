# uCISC Syntax

* Back: [The Path Not Taken](11_The_Path_Not_Taken.md)
* Next: [Beginning to Code](12_Beginning_To_Code.md)

uCISC syntax is heavily influenced by [Mu](https://github.com/akkartik/mu).

If you use vim, there is syntax highlighting available:
[extras/ucisc.vim](/extras/ucisc.vim). 

### Fundamentals

The basic instruction format is:

```
copy|compute [ALU.op] <source> <destination> [effect] [push|pop]
```

uCISC only has 2 instructions, copy and compute. Copy instructions copy a value from
the source to the destination. Compute instructions transform the data in some way,
as specified by the ALU op code.

Each statement in uCISC translates to a single machine instruction and fits in a
single memory word (16-bits). The compiler will validate that all the limitations are
fulfilled (e.g. the immediate value is within bounds).

Arguments consist of a numerical code and a type. The numerical code is generally the
value of the bits that get packed into the instruction. The only exception is "1.mem"
(encoded as 1) vs. "1.reg" (encoded as 5). This exception is made because they both
use the same register (register 1).

```
  # Register and memory examples
  1.mem # memory value at the address in r1
  1.reg # register value in r1
  4.imm # immediate value as an argument
  0.reg # PC register
  4.reg # control register

  # Other argument types
  3.eff # the conditional effect of the instruction
  1.push # decrements destination register before storing (for pushing to the stack)
  1.pop # increments the source register after execution (for popping from the stack)
  0xA.op # Add op code for compute instructions

  # Labels are interpreted as immediate values
  start.imm # Absolute address reference of the start label
  start.disp # Relative reference from the current instruction to the start label 
```

Numbers are interpreted as decimal values unless prefixed by `0x`.

#### Position Within Statements

The arguments must follow the position rules for uCISC statements:

1. The first position must be the instruction (copy or compute)
2. The source must be first, destination second
3. The immediate value must follow the argument it affects (source or destination)

#### Position Types

Source and destination types:
 - `reg` - register, the content of the register is being referenced
 - `mem` - memory location, the content of the memory referenced by the register is being referenced
 - `val` - immediate value, the content of the immediate value is being referenced

*Note:* as a shortcut, a source argument that only specifies an immediate value
implies 4.val as the actual argument.

Other arguments:
 - `imm` - Immediate value, bit-width and signed-ness depends on context. Can be
   applied to a label to resolve to the absolute address of a label.
 - `disp` - immediate value, if label and will resolve to difference between the
   current instruction and the label. Must result in a valid imm value for the given
   instruction.
 - `eff` - Effect code, values between 0 and 3, defaults to 3.eff
 - `push`, `pop` - Increment argument, value of 1 is implied

#### Comments

There are two types of comments:
 - inline comments begin with a '/' character and terminate at the next
   '/' character (e.g.  `<arg>/comment/ <arg>`) on the same line
 - line comments start with a '#' character and run all the way to the
   newline.

Before processing a statement, the compiler will strip all comments from the line. Be
careful about inline comments that cause args to collide when removed (e.g.
`<arg>/comment/<arg>`). You will want to include a whitespace on one side of the
inline comment.

The following two statements are identical:

```
  # Readable statement:
  copy 1.reg/stack address/ 2.imm/offset/ 2.reg/other comment/

  # Without inline or word comments:
  copy 1.reg 2.imm 2.reg
```

#### Data

Data can be added to uCISC files with hex values. The line must begin with a `%`
symbol and there must be a multiple of 4 hex characters on the line. Typically you
will want to add a label to refer to the beginning of the data block so you can use
it in your code.

```
my_data:
% 001F 0A48 656C 6C6f 2C20 776F 726C 6421 2075 4349 5343 2069 7320 6865 7265 210A 0A00

# Make 2.reg point to my_data
copy my_data.imm 2.reg
```

#### Strings

Strings are in quotes. They will be converted inline to string data (a length
followed by ascii data). Strings are packed 2 chars per word. The last word will end
with 0x00 if there are an odd number of characters.

For example, the following are equivalent:
```
# String form
hello_world:
"\nHello, world! uCISC is here!\n\n"

# Data equivalent
hello_world:
% 001F 0A48 656C 6C6f 2C20 776F 726C 6421 2075 4349 5343 2069 7320 6865 7265 210A 0A00
```

Strings are converted to data inline, so their placement must be appropriate for the
context.

#### Variables

uCISC has a simple variable syntax to allow easier to understand code. Variables
begin with an `&` or a `$`. The `$` variant refers to "mem" arguments while the `&`
variant refers to "reg" arguments. That is, `$my_val` will refer to the value of
my_val and `&my_val` will refer to the address of my_val. When you declare either
flavor, the other is automatically created as well if appropriate.

```
# The "as" keyword does simple substitution
$if_zero as 0.eff
# You can now use $if_zero instead of 0.eff in statements

# You can create them without an instruction
$stack as 1.mem
# &stack now evaluates to 1.reg, no instruction was generated

# You can create them as a result of an instruction
$my_data = copy $stack $stack push

# Add clarity by refering to PC by name
&pc as 0.reg

# These are invalid because the type is wrong:
$pc as 0.reg # Syntax error, not a mem register
&my_data = copy $stack $stack push # Syntax error, statement destination is mem
$my_ref = copy $stack 2.reg # Syntax error, statement destination is reg
```

Instead of specifying "as" you can use "=". These arguments are indexed variables,
that keep track of changing offsets as you modify their registers. This is best seen by examples:

```
$stack as 1.mem

# push new value to stack
copy 4.imm $stack push
# $stack now refers to a different memory location since it was the target of a push

# modifying offsets
$val = copy 4.imm $stack push
# &stack is affected as above
# &val is equivalent to &stack

copy 1.imm $stack push
# &stack is affected as above
# &val is now &stack 1.imm

# This now does the right thing
copy $val 2.reg
# equivalent to "copy $stack 1.imm 2.reg"
# equivalent to "copy 1.mem 1.imm 2.reg"
```

This is very important for easily writing code because you would otherwise find
yourself counting stack offset all the time. I find they are easy to get wrong and
lead to bugs where I am off by one on the stack reference because I miscounted the
number of pushes I did. It allows you to add a push without manually recomputing the
later offsets.

This creates a useful syntax for popping multiple values off the stack at once:
```
# Everything push after $val is popped:
copy &val &stack

# Everything push after $val is popped, including $val itself
copy &val 1.imm &stack # actual immediate value is 1.imm + the immediate of &val
```

#### Functions

Functions are the only syntax sugar that resolves to multiple instructions. The
format is as follows:

```
1.mem[1] <= function_name([<arg>[, <arg>[, ...]]])
```

The function name must resolve to a label somewhere. Arguments are a comma separated
list of valid "source" arguments for instructions. The "mem" value to the left of the
arrow refers to the stack register where arguments are pushed. The value in the
brackets indicates the number of words pushed to the stack by the function.

Function calls are compatible with variables:

```
$stack from 1.mem

# You can init variables already on the stack
$arg1, $arg2 from $stack
# $arg1 is the same as $stack, $arg2 is the $stack register offset by 1

# Or you can directly reference the offset
$my_data = $stack[2]

# You can leave off the brackets if nothing is returned
$stack <= my_function($my_var)

# You can assign return values to new variables
$result = $stack[1] <= my_function($my_var)

# Or with multiple values returned, assign them in one go
$start, $finish = $stack[2] <= my_function($my_var)
```

When multiple values are returned, the first one refers to the top of the stack, the
second is offset by 1, and so on. Remember, variables are just mem arguments with an
immediate offset that is kept track of by the compiler.

Function calls get translated into multiple statements:

1. The stack pointer is decremented by the return word count. Called functions are
  expected to put return values in these slots.
2. The return address is pushed to the stack
3. The arguments are pushed to the stack one by one
4. The PC register is loaded with the jump address

For example:

```
$stack <= do_something
# becomes
copy 0.reg 2.imm $stack # return address
copy do_something.disp 0.reg # jump

# and
$stack[5] <= something_else(3.imm)
# becomes
copy &stack 5.imm &stack # make room for returns
copy 0.reg 3.imm $stack # return address
copy 3.imm $stack # push arg
copy something_else.disp 0.reg # jump
```

I have found this calling convention to mostly result in the cleanest code in uCISC.
Before returning, the called function is responsible for popping all values from the
stack except the return arguments.

Function declaration is:

```
function_name: [$arg, $return, $result = $stack]
```

Note that the variables on the same line as the function are syntax sugar. It's just
as valid to put that declaration on the next line. If you declare your functions in
this manner, the compiler can do argument matching and ensure that the signature of
called functions matches the callee declaration, at least in terms of argument and
result values. No type checking is performed on the content of those arguments.

### Control blocks

uCISC provides a label convention using curly braces ('{' and '}'). These are used
with special labels to allow easy clear control flow.

```
{
  copy break.disp 0.eff

  # Some logic here
  copy loop.disp 1.eff
}
```

The `break.disp` refers to the end brace and `loop.disp` refers to the open brace.
Curly braces can be nested and the `break` and `loop` labels refer to the inner most
pair.

*Note:* Braces also modify variable scope. Variables first declared inside a brace
will pass out of scope at the end of that brace section. Assigning a new register
and/or index to a variable inside a sub-scope affects the variable for the root
most scope the variable exists in. That is, variables are global to all child scopes
in the context they are created in.

#### Continue Reading

* Back: [The Path Not Taken](11_The_Path_Not_Taken.md)
* Next: [Beginning to Code](12_Beginning_To_Code.md)

