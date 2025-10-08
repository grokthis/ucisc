# Quick Start Guide - uCISC Go Emulator

## Prerequisites

- Go 1.21 or later
- Make (optional, but recommended)

## Installation

### Option 1: Build from Source

```bash
cd ucisc-go
make build
```

The binary will be in `bin/ucisc`

### Option 2: Install to GOPATH

```bash
cd ucisc-go
make install
```

Now `ucisc` is in your PATH

## Your First Program

### 1. Create a simple program

Create a file `test.hex`:

```
# My first uCISC program
5400 0007  # r1 <- copy val/7
544A 0005  # r1 <- add val/5
0040 0000  # HALT
```

### 2. Run it

```bash
./bin/ucisc run test.hex
```

You should see:

```
Running program: test.hex

Program completed successfully
Cycles: 3
Instructions: 3

Registers:
  PC:      0x0006
  R1:      0x000C (12)
  R2:      0x0000 (0)
  R3:      0x0000 (0)
  R4:      0x0000 (0)
  R5:      0x0000 (0)
  R6:      0x0000 (0)
  Banking: 0x0000
  Flags:   0x0000 []
```

Success! Your program added 7 + 5 = 12

## Running Examples

### Simple Addition

```bash
make run-simple
```

### Fibonacci Sequence

```bash
make run-fib
```

## Understanding the Hex Format

Each uCISC instruction is **2 words** (32 bits total), written as 4 hex digits each:

```
WWWW WWWW  # First word, second word
```

### Example Instructions

```
# Copy immediate value to register
5400 0005  # r1 <- copy val/5
           # 5400 = source(val=4), dest(r1=5), effect(always=4), op(copy=0)
           # 0005 = immediate value (5)

# Add with immediate
544A 0003  # r1 <- add val/3
           # 544A = source(val=4), dest(r1=5), effect(always=4), op(add=A)
           # 0003 = immediate value (3)

# HALT
0040 0000  # pc <- copy pc/0
           # Special case: copy PC to itself = HALT
```

### Comments

- Use `#` for comments
- Comments can be on their own line or after hex values
- Blank lines are ignored

## Writing Programs

### Basic Structure

```
# Initialize
5400 0000  # r1 <- copy val/0

# Your code here
544A 0001  # r1 <- add val/1

# Always end with HALT
0040 0000  # HALT
```

### Using Memory

```
# Set r1 to address 100
5400 0064  # r1 <- copy val/100

# Write to memory[100]
1440 002A  # r1/0 <- copy val/42

# Read from memory[100]
6410 0000  # r2 <- copy r1/0
```

### Conditional Execution

Effects determine when instructions execute:

- `<-` or `<|` (bits: 4) - Always execute
- `<0` or `<|0` (bits: 0) - Execute if zero
- `<1` or `<|1` (bits: 1) - Execute if not zero
- `<n` or `<|n` (bits: 2) - Execute if negative
- `<+` or `<|+` (bits: 3) - Execute if positive

```
# Set r1 to 0, then conditionally execute
5400 0000  # r1 <- copy val/0 (sets zero flag)
6000 000A  # r2 <0 copy val/10 (executes because zero)
6410 0014  # r2 <1 copy val/20 (doesn't execute)
```

### Loops

```
# Count from 0 to 5
5400 0000  # r1 <- copy val/0
544A 0001  # r1 <- add val/1
547B 0005  # r1 <~ sub val/5 (test, don't store)
0410 0002  # pc <1 copy val/2 (jump if not zero)
0040 0000  # HALT
```

## Instruction Reference

### Source Types
- `val` - Immediate value (4)
- `r1`, `r2`, `r3` - Register value (5,6,7)
- `r1`, `r2`, `r3` - Memory at register (1,2,3)

### Destination Types
- `r1`, `r2`, `r3` - Register (5,6,7)
- `r1`, `r2`, `r3` - Memory at register (1,2,3)
- `pc` - Program counter (0)

### Operations
- `copy` (0) - Copy
- `and` (1) - Bitwise AND
- `or` (2) - Bitwise OR
- `xor` (3) - Bitwise XOR
- `inv` (4) - Bitwise NOT
- `shl` (5) - Shift left
- `shr` (6) - Shift right
- `swap` (7) - Swap bytes
- `msb` (8) - MSB only
- `lsb` (9) - LSB only
- `add` (A) - Add
- `sub` (B) - Subtract
- `mult` (C) - Multiply (LSW)
- `multsw` (D) - Multiply (MSW)
- `addc` (E) - Add with carry

## Encoding Guide

### Instruction Format

Non-memory destination:
```
SSSS DDDD MEEE AAAA  IIII IIII IIII IIII
```

Memory destination:
```
SSSS DDDD MEEE AAAA  OOOO IIII IIII IIII
```

Where:
- `SSSS` = Source (4 bits)
- `DDDD` = Destination (4 bits)
- `M` = Increment/Push/Pop (1 bit)
- `EEE` = Effect (3 bits)
- `AAAA` = Operation (4 bits)
- `OOOO` = Offset for memory dest (4 bits)
- `I` = Immediate (16 bits or 12 bits)

### Common Patterns

**Load immediate to register:**
```
5400 VVVV  # r1 <- copy val/VVVV
6400 VVVV  # r2 <- copy val/VVVV
7400 VVVV  # r3 <- copy val/VVVV
```

**Add immediate:**
```
544A VVVV  # r1 <- add val/VVVV
644A VVVV  # r2 <- add val/VVVV
```

**Jump to address:**
```
0440 AAAA  # pc <- copy val/AAAA
```

**Conditional jump:**
```
0410 AAAA  # pc <1 copy val/AAAA (jump if not zero)
0400 AAAA  # pc <0 copy val/AAAA (jump if zero)
```

## Testing Your Code

### Run Tests

```bash
make test
```

### Run with Verbose Output

```bash
go test -v ./...
```

### Check Coverage

```bash
make coverage
```

## Common Issues

### Program doesn't halt
- Make sure to end with `0040 0000` (HALT instruction)

### Wrong result
- Check instruction encoding
- Verify operation codes (add=A, sub=B, etc.)
- Check if immediate values are correct

### Parse error
- Each instruction must be exactly 2 words (4 hex digits each)
- Check for typos in hex values
- Make sure comments start with `#`

## Next Steps

1. **Read the full documentation**:
   - [ARCHITECTURE.md](../docs/go-emulator/ARCHITECTURE.md)
   - [REQUIREMENTS.md](../docs/go-emulator/REQUIREMENTS.md)

2. **Study the examples**:
   - `examples/simple.hex` - Basic operations
   - `examples/fibonacci.hex` - Loops and conditionals

3. **Explore the instruction set**:
   - [Instruction Set Details](../docs/programming_guide/5_Instruction_Set_Details.md)

4. **Write your own programs**:
   - Start simple (addition, subtraction)
   - Try loops
   - Experiment with memory
   - Build something cool!

## Getting Help

- Check the [uCISC documentation](../docs/Overview.md)
- Look at example programs in `examples/`
- Read the test files to see how things work
- Open an issue on GitHub

## Contributing

We welcome contributions! Areas that need work:
- Device implementations (GPIO, UART, I2C)
- Debugger features
- More example programs
- Documentation improvements

---

Happy coding with uCISC! 🚀
