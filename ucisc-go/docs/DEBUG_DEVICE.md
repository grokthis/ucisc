# Debug Device Specification

## Overview

The debug device provides testing and debugging facilities for uCISC programs. It is accessed by writing to the `&debug` register (which is actually the flags register - writes are intercepted since flags is read-only).

**Device ID**: 13 (0x000D)
**Device Address**: 0x00D0 - 0x00DF (control registers only, no block memory)
**Device Type**: 0x0D (Debug device)

## Access Method

```assembly
# Write to &debug register (flags register on write)
# Note: Can write from immediate values or registers
&debug <- copy val/0x0001   # Stage command (immediate value)
&debug <- copy &r3          # Stage arg1 from register
&debug <- copy val/42       # Stage arg2 (immediate value)
&debug <- copy val/-1       # Execute staged command (0xFFFF)
```

## Protocol

1. **Stage command**: Write command code to `&debug`
2. **Stage arguments**: Write arguments to `&debug` (if command needs them)
3. **Execute or Reset**: Write special value to `&debug`:
   - `-1` (0xFFFF) = Execute staged command
   - `-2` (0xFFFE) = Reset all state (clear failures, staged values)
4. **Read status**: Read device control registers to check results

Special control values:
- `-1` (0xFFFF): Execute - triggers execution of staged command with staged arguments
- `-2` (0xFFFE): Reset - clears all failures, staged command, and staged arguments

## Device Control Layout (0x00D0 - 0x00DF)

```
Address  | R/W | Description
---------|-----|----------------------------------------------------------
0x00D0   | R   | Device ID (0x000D)
0x00D1   | R   | Flags (MSB, bits 8-15) | Device Type (LSB, 0x0D = Debug)
0x00D2   | R/W | Behavior flags (see below)
0x00D3   | R   | Status (0=pass/idle, 1=fail, 2=breakpoint)
0x00D4   | R   | Failure count (increments on each assertion failure)
0x00D5   | R   | Last failed command code
0x00D6   | R   | Last failed arg1 value (for debugging)
0x00D7   | R   | Last failed arg2 value (for debugging)
0x00D8   | R/W | Reserved
0x00D9   | R/W | Reserved
0x00DA   | R/W | Reserved
0x00DB   | R/W | Reserved
0x00DC   | R/W | Reserved
0x00DD   | R/W | Reserved
0x00DE   | R/W | Reserved
0x00DF   | R/W | Reserved
```

## Behavior Flags (0x00D2)

```
Bit  | Description
-----|------------------------------------------------------------
0    | HALT_ON_FAIL - Set CPU halt flag when assertion fails
1    | PRINT_ON_FAIL - Print failure message (emulator only)
2-15 | Reserved
```

## Commands

### Hardware-Compatible Commands

These commands work on both real hardware and emulator. They set device status flags and can halt the CPU based on behavior flags.

```
Code   | Command              | Args | Description
-------|----------------------|------|----------------------------------------
0x0001 | ASSERT_EQ            | 2    | Assert arg1 == arg2
0x0002 | ASSERT_NE            | 2    | Assert arg1 != arg2
0x0003 | ASSERT_GT            | 2    | Assert arg1 > arg2 (unsigned)
0x0004 | ASSERT_LT            | 2    | Assert arg1 < arg2 (unsigned)
0x0005 | ASSERT_GTE           | 2    | Assert arg1 >= arg2 (unsigned)
0x0006 | ASSERT_LTE           | 2    | Assert arg1 <= arg2 (unsigned)
0x0007 | ASSERT_GT_SIGNED     | 2    | Assert arg1 > arg2 (signed)
0x0008 | ASSERT_LT_SIGNED     | 2    | Assert arg1 < arg2 (signed)
0x0009 | ASSERT_ZERO          | 1    | Assert arg1 == 0
0x000A | ASSERT_NOT_ZERO      | 1    | Assert arg1 != 0
0x000B | ASSERT_HAS_FLAGS     | 1    | Assert (CPU flags & arg1) != 0
0x000C | ASSERT_NO_FLAGS      | 1    | Assert (CPU flags & arg1) == 0
0x000D | Reserved             | -    | -
0x000E | Reserved             | -    | -
0x000F | Reserved             | -    | -
```

### Emulator-Only Commands

These commands are only functional in the emulator and are ignored on real hardware. Useful for debugging and testing.

```
Code   | Command              | Args | Description
-------|----------------------|------|----------------------------------------
0x0100 | BREAKPOINT           | 0    | Trigger breakpoint (halt + dump state)
0x0101 | PRINT_VALUE          | 1    | Print arg1 to console (decimal and hex)
0x0102 | PRINT_STRING         | 1    | Print null-terminated string at arg1
0x0103 | DUMP_REGISTERS       | 0    | Dump all CPU registers to console
0x0104 | DUMP_MEMORY          | 2    | Dump memory from arg1, length arg2
0x0105 | TRACE_ON             | 0    | Enable instruction trace logging
0x0106 | TRACE_OFF            | 0    | Disable instruction trace logging
0x0107 | LOG_MESSAGE          | 1    | Log message ID arg1 (app-defined)
```

## Usage Examples

### Basic Assertions

```assembly
# Assert that r3 == 42
&debug <- copy val/0x0001    # Stage ASSERT_EQ command
&debug <- copy &r3           # Stage arg1 = r3
&debug <- copy val/42        # Stage arg2 = 42
&debug <- copy val/-1        # Execute (0xFFFF)

# Assert that r1 != 0
&debug <- copy val/0x000A    # Stage ASSERT_NOT_ZERO
&debug <- copy &r1           # Stage arg1 = r1
&debug <- copy val/-1        # Execute (0xFFFF)

# Assert zero flag is set
&debug <- copy val/0x000B    # Stage ASSERT_HAS_FLAGS
&debug <- copy val/0x0001    # Stage arg1 = FlagZero
&debug <- copy val/-1        # Execute (0xFFFF)

# Reset all state (failures and staged values)
&debug <- copy val/-2        # Reset (0xFFFE)
```

### Emulator Commands

```assembly
# Print a value
&debug <- copy val/0x0101    # Stage PRINT_VALUE
&debug <- copy &r3           # Stage arg1 = r3 value
&debug <- copy val/-1        # Execute (0xFFFF)

# Dump registers
&debug <- copy val/0x0103    # Stage DUMP_REGISTERS
&debug <- copy val/-1        # Execute (0xFFFF)

# Breakpoint
&debug <- copy val/0x0100    # Stage BREAKPOINT
&debug <- copy val/-1        # Execute (0xFFFF)
```

### Configure Behavior (Requires Banking)

To read/write device control registers, you need to enable banking:

```assembly
# Save banking state
r5 <- copy &banking
r7 push <- copy &r5

# Enable R1 for device memory
banking <- copy val/0b000010

# Configure to halt on assertion failure
r1 <- copy val/0x00D2        # Behavior flags address (moved earlier)
r6 <- copy r1                # Read current flags
r6 <- copy val/0b0001        # Set HALT_ON_FAIL
r1 <- copy &r6               # Write back

# Restore banking
r5 <- copy r7 pop
banking <- copy &r5
```

### Assembler Syntactic Sugar

The assembler will provide syntactic sugar for common operations:

```assembly
# Function-style syntax (assembler expands to writes + execute)
debug.assertEqual(&r3, val/42)
# Expands to:
#   &debug <- copy val/0x0001
#   &debug <- copy &r3
#   &debug <- copy val/42
#   &debug <- copy val/-1

debug.assertNotZero(&r1)
# Expands to:
#   &debug <- copy val/0x000A
#   &debug <- copy &r1
#   &debug <- copy val/-1

debug.reset()
# Expands to:
#   &debug <- copy val/-2

# Emulator functions
debug.printValue(&r3)
debug.dumpRegisters()
debug.breakpoint()
```

## Implementation Notes

### CPU Implementation

1. **Flags Register Write Interception**: When a write to the flags register is detected, route it to the debug device instead
2. **Device State**: Debug device maintains:
   - Staged command code
   - Staged arg1
   - Staged arg2
   - Failure history
3. **Special Values**:
   - `-1` (0xFFFF): Execute staged command with staged arguments
   - `-2` (0xFFFE): Reset all state (clear failures, staged values)

### Hardware Implementation (Verilog/VHDL)

```verilog
// Pseudo-code for hardware implementation
always @(posedge clk) begin
    if (flags_write_enable) begin
        if (flags_write_data == 16'hFFFF) begin
            // Execute staged command (-1)
            case (staged_cmd)
                16'h0001: status <= (staged_arg1 == staged_arg2) ? 0 : 1;
                16'h0002: status <= (staged_arg1 != staged_arg2) ? 0 : 1;
                // ... other commands
            endcase

            if (status == 1 && behavior_flags[0]) begin
                cpu_halt <= 1;
            end
        end else if (flags_write_data == 16'hFFFE) begin
            // Reset all state (-2)
            staged_cmd <= 0;
            staged_arg1 <= 0;
            staged_arg2 <= 0;
            status <= 0;
            fail_count <= 0;
        end else begin
            // Stage value
            if (staged_cmd == 0) begin
                staged_cmd <= flags_write_data;
            end else if (staged_arg1 == 0) begin
                staged_arg1 <= flags_write_data;
            end else begin
                staged_arg2 <= flags_write_data;
            end
        end
    end
end
```

### Assembler Implementation

1. **Lexer**: Recognize `debug` as a keyword and `.` for method calls
2. **Parser**: Parse `debug.functionName(args...)` syntax
3. **Code Generator**: Map function names to command codes and expand to instruction sequence
4. **Validation**: Check argument count matches command requirements

## Benefits

1. **Zero ISA Impact**: Uses existing write-to-flags mechanism
2. **Hardware Simple**: Just stage values and execute on trigger
3. **Explicit Execution**: No arg-counting logic needed
4. **Hardware Compatible**: Core assertions work on real hardware
5. **Flexible**: Easy to add new commands
6. **Clean Syntax**: Function-style calls via assembler sugar
7. **Discoverable**: Device control registers document state
8. **No Banking for Writes**: Direct `&debug` register writes for commands/args
9. **Banking Only for Config**: Only need banking to read/write device control registers
