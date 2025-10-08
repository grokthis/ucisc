package cpu

import (
	"fmt"
	"strings"
)

// Debug device constants
const (
	DebugDeviceID     uint16 = 0x000D
	DebugDeviceType   uint8  = 0x0D
	DebugDeviceNumber uint8  = 13
)

// Debug commands - Hardware compatible
const (
	CmdAssertEq         uint16 = 0x0001
	CmdAssertNe         uint16 = 0x0002
	CmdAssertGt         uint16 = 0x0003
	CmdAssertLt         uint16 = 0x0004
	CmdAssertGte        uint16 = 0x0005
	CmdAssertLte        uint16 = 0x0006
	CmdAssertGtSigned   uint16 = 0x0007
	CmdAssertLtSigned   uint16 = 0x0008
	CmdAssertZero       uint16 = 0x0009
	CmdAssertNotZero    uint16 = 0x000A
	CmdAssertHasFlags   uint16 = 0x000B
	CmdAssertNoFlags    uint16 = 0x000C
)

// Debug commands - Emulator only
const (
	CmdBreakpoint    uint16 = 0x00D0
	CmdPrintValue    uint16 = 0x00D1
	CmdPrintString   uint16 = 0x00D2
	CmdDumpRegisters uint16 = 0x00D3
	CmdDumpMemory    uint16 = 0x00D4
	CmdTraceOn       uint16 = 0x00D5
	CmdTraceOff      uint16 = 0x00D6
	CmdLogMessage    uint16 = 0x00D7
)

// Special control values
const (
	CtrlExecute uint16 = 0xFFFF // -1
	CtrlReset   uint16 = 0xFFFE // -2
)

// Behavior flags
const (
	BehaviorHaltOnFail  uint16 = 1 << 0
	BehaviorPrintOnFail uint16 = 1 << 1
)

// DebugDevice implements a debug device for testing and debugging
type DebugDevice struct {
	cpu *CPU // Reference to CPU for state access

	// Device control registers
	behaviorFlags uint16
	status        uint16
	failCount     uint16
	lastFailedCmd uint16
	lastFailedArg1 uint16
	lastFailedArg2 uint16

	// Staged values
	stagedCmd     uint16
	stagedArg1    uint16
	stagedArg2    uint16
	stagedArgCount int // Track how many args have been staged

	// Emulator state
	failures []string
	traceEnabled bool
}

// NewDebugDevice creates a new debug device
func NewDebugDevice(cpu *CPU) *DebugDevice {
	return &DebugDevice{
		cpu:      cpu,
		failures: make([]string, 0),
	}
}

// ID returns the device ID
func (d *DebugDevice) ID() uint16 {
	return DebugDeviceID
}

// Type returns the device type
func (d *DebugDevice) Type() uint8 {
	return DebugDeviceType
}

// Read reads from device control registers
func (d *DebugDevice) Read(offset uint16) uint16 {
	switch offset {
	case 0x0: // Device ID
		return d.ID()
	case 0x1: // Flags (MSB) | Device Type (LSB)
		return uint16(d.Type())
	case 0x2: // Behavior flags
		return d.behaviorFlags
	case 0x3: // Status
		return d.status
	case 0x4: // Failure count
		return d.failCount
	case 0x5: // Last failed command
		return d.lastFailedCmd
	case 0x6: // Last failed arg1
		return d.lastFailedArg1
	case 0x7: // Last failed arg2
		return d.lastFailedArg2
	default:
		return 0
	}
}

// Write writes to device control registers
func (d *DebugDevice) Write(offset uint16, value uint16) {
	switch offset {
	case 0x2: // Behavior flags
		d.behaviorFlags = value
	default:
		// Other registers are read-only
	}
}

// ReadBlock reads from device block memory (not used)
func (d *DebugDevice) ReadBlock(offset uint16) uint16 {
	return 0
}

// WriteBlock writes to device block memory (not used)
func (d *DebugDevice) WriteBlock(offset uint16, value uint16) {
	// Not used
}

// WriteDebug handles writes to the &debug register (intercepted flags writes)
func (d *DebugDevice) WriteDebug(value uint16) {
	switch value {
	case CtrlExecute:
		d.execute()
	case CtrlReset:
		d.reset()
	default:
		d.stage(value)
	}
}

// stage stages a value (command or argument)
func (d *DebugDevice) stage(value uint16) {
	if d.stagedArgCount == 0 && d.stagedCmd == 0 {
		// First write is the command
		d.stagedCmd = value
	} else if d.stagedArgCount == 0 {
		// First argument
		d.stagedArg1 = value
		d.stagedArgCount = 1
	} else if d.stagedArgCount == 1 {
		// Second argument
		d.stagedArg2 = value
		d.stagedArgCount = 2
	}
	// Ignore additional writes beyond 2 arguments
}

// execute executes the staged command
func (d *DebugDevice) execute() {
	cmd := d.stagedCmd
	arg1 := d.stagedArg1
	arg2 := d.stagedArg2

	// Clear staged values after execution
	defer func() {
		d.stagedCmd = 0
		d.stagedArg1 = 0
		d.stagedArg2 = 0
		d.stagedArgCount = 0
	}()

	// Execute command
	switch cmd {
	// Hardware-compatible assertions
	case CmdAssertEq:
		d.assertEqual(arg1, arg2)
	case CmdAssertNe:
		d.assertNotEqual(arg1, arg2)
	case CmdAssertGt:
		d.assertGreater(arg1, arg2)
	case CmdAssertLt:
		d.assertLess(arg1, arg2)
	case CmdAssertGte:
		d.assertGreaterOrEqual(arg1, arg2)
	case CmdAssertLte:
		d.assertLessOrEqual(arg1, arg2)
	case CmdAssertGtSigned:
		d.assertGreaterSigned(arg1, arg2)
	case CmdAssertLtSigned:
		d.assertLessSigned(arg1, arg2)
	case CmdAssertZero:
		d.assertZero(arg1)
	case CmdAssertNotZero:
		d.assertNotZero(arg1)
	case CmdAssertHasFlags:
		d.assertHasFlags(Flags(arg1))
	case CmdAssertNoFlags:
		d.assertNoFlags(Flags(arg1))

	// Emulator-only commands
	case CmdBreakpoint:
		d.breakpoint()
	case CmdPrintValue:
		d.printValue(arg1)
	case CmdPrintString:
		d.printString(arg1)
	case CmdDumpRegisters:
		d.dumpRegisters()
	case CmdDumpMemory:
		d.dumpMemory(arg1, arg2)
	case CmdTraceOn:
		d.traceEnabled = true
	case CmdTraceOff:
		d.traceEnabled = false
	case CmdLogMessage:
		d.logMessage(arg1)
	}
}

// reset resets all debug state
func (d *DebugDevice) reset() {
	d.stagedCmd = 0
	d.stagedArg1 = 0
	d.stagedArg2 = 0
	d.stagedArgCount = 0
	d.status = 0
	d.failCount = 0
	d.lastFailedCmd = 0
	d.lastFailedArg1 = 0
	d.lastFailedArg2 = 0
	d.failures = make([]string, 0)
}

// assertEqual asserts that two values are equal
func (d *DebugDevice) assertEqual(actual, expected uint16) {
	if actual != expected {
		d.recordFailure(CmdAssertEq, actual, expected,
			fmt.Sprintf("Assert Equal failed: expected %d (0x%04X), got %d (0x%04X) at PC=%d",
				expected, expected, actual, actual, d.cpu.pc))
	}
}

// assertNotEqual asserts that two values are not equal
func (d *DebugDevice) assertNotEqual(actual, expected uint16) {
	if actual == expected {
		d.recordFailure(CmdAssertNe, actual, expected,
			fmt.Sprintf("Assert Not Equal failed: values are equal: %d (0x%04X) at PC=%d",
				actual, actual, d.cpu.pc))
	}
}

// assertGreater asserts that actual > expected (unsigned)
func (d *DebugDevice) assertGreater(actual, expected uint16) {
	if actual <= expected {
		d.recordFailure(CmdAssertGt, actual, expected,
			fmt.Sprintf("Assert Greater failed: %d <= %d at PC=%d",
				actual, expected, d.cpu.pc))
	}
}

// assertLess asserts that actual < expected (unsigned)
func (d *DebugDevice) assertLess(actual, expected uint16) {
	if actual >= expected {
		d.recordFailure(CmdAssertLt, actual, expected,
			fmt.Sprintf("Assert Less failed: %d >= %d at PC=%d",
				actual, expected, d.cpu.pc))
	}
}

// assertGreaterOrEqual asserts that actual >= expected (unsigned)
func (d *DebugDevice) assertGreaterOrEqual(actual, expected uint16) {
	if actual < expected {
		d.recordFailure(CmdAssertGte, actual, expected,
			fmt.Sprintf("Assert Greater or Equal failed: %d < %d at PC=%d",
				actual, expected, d.cpu.pc))
	}
}

// assertLessOrEqual asserts that actual <= expected (unsigned)
func (d *DebugDevice) assertLessOrEqual(actual, expected uint16) {
	if actual > expected {
		d.recordFailure(CmdAssertLte, actual, expected,
			fmt.Sprintf("Assert Less or Equal failed: %d > %d at PC=%d",
				actual, expected, d.cpu.pc))
	}
}

// assertGreaterSigned asserts that actual > expected (signed)
func (d *DebugDevice) assertGreaterSigned(actual, expected uint16) {
	if int16(actual) <= int16(expected) {
		d.recordFailure(CmdAssertGtSigned, actual, expected,
			fmt.Sprintf("Assert Greater Signed failed: %d <= %d at PC=%d",
				int16(actual), int16(expected), d.cpu.pc))
	}
}

// assertLessSigned asserts that actual < expected (signed)
func (d *DebugDevice) assertLessSigned(actual, expected uint16) {
	if int16(actual) >= int16(expected) {
		d.recordFailure(CmdAssertLtSigned, actual, expected,
			fmt.Sprintf("Assert Less Signed failed: %d >= %d at PC=%d",
				int16(actual), int16(expected), d.cpu.pc))
	}
}

// assertZero asserts that value == 0
func (d *DebugDevice) assertZero(value uint16) {
	if value != 0 {
		d.recordFailure(CmdAssertZero, value, 0,
			fmt.Sprintf("Assert Zero failed: value is %d (0x%04X) at PC=%d",
				value, value, d.cpu.pc))
	}
}

// assertNotZero asserts that value != 0
func (d *DebugDevice) assertNotZero(value uint16) {
	if value == 0 {
		d.recordFailure(CmdAssertNotZero, value, 0,
			fmt.Sprintf("Assert Not Zero failed: value is 0 at PC=%d", d.cpu.pc))
	}
}

// assertHasFlags asserts that CPU flags contain the specified flags
func (d *DebugDevice) assertHasFlags(flags Flags) {
	if !d.cpu.flags.IsSet(flags) {
		d.recordFailure(CmdAssertHasFlags, uint16(flags), uint16(d.cpu.flags),
			fmt.Sprintf("Assert Has Flags failed: expected flags 0x%04X, CPU flags are 0x%04X at PC=%d",
				flags, d.cpu.flags, d.cpu.pc))
	}
}

// assertNoFlags asserts that CPU flags don't contain the specified flags
func (d *DebugDevice) assertNoFlags(flags Flags) {
	if d.cpu.flags.IsSet(flags) {
		d.recordFailure(CmdAssertNoFlags, uint16(flags), uint16(d.cpu.flags),
			fmt.Sprintf("Assert No Flags failed: flags 0x%04X are set, CPU flags are 0x%04X at PC=%d",
				flags, d.cpu.flags, d.cpu.pc))
	}
}

// recordFailure records an assertion failure
func (d *DebugDevice) recordFailure(cmd uint16, arg1, arg2 uint16, msg string) {
	d.status = 1 // Fail
	d.failCount++
	d.lastFailedCmd = cmd
	d.lastFailedArg1 = arg1
	d.lastFailedArg2 = arg2
	d.failures = append(d.failures, msg)

	if d.behaviorFlags&BehaviorPrintOnFail != 0 {
		fmt.Println(msg)
	}

	if d.behaviorFlags&BehaviorHaltOnFail != 0 {
		d.cpu.flags = d.cpu.flags.Set(FlagHalt)
	}
}

// breakpoint triggers a breakpoint
func (d *DebugDevice) breakpoint() {
	fmt.Printf("\n*** BREAKPOINT at PC=%d ***\n", d.cpu.pc)
	d.dumpRegisters()
	d.status = 2 // Breakpoint
}

// printValue prints a value to console
func (d *DebugDevice) printValue(value uint16) {
	fmt.Printf("[DEBUG] Value at PC=%d: %d (0x%04X)\n", d.cpu.pc, value, value)
}

// printString prints a null-terminated string from memory
func (d *DebugDevice) printString(addr uint16) {
	var sb strings.Builder
	for {
		word := d.cpu.memory.ReadLocal(addr)
		if word == 0 {
			break
		}
		// Two characters per word (MSB, LSB)
		if msb := byte(word >> 8); msb != 0 {
			sb.WriteByte(msb)
		}
		if lsb := byte(word & 0xFF); lsb != 0 {
			sb.WriteByte(lsb)
		}
		addr++
	}
	fmt.Printf("[DEBUG] String at PC=%d: %s\n", d.cpu.pc, sb.String())
}

// dumpRegisters dumps all CPU registers
func (d *DebugDevice) dumpRegisters() {
	fmt.Printf("\n=== Register Dump at PC=%d ===\n", d.cpu.pc)
	fmt.Printf("PC:       0x%04X (%d)\n", d.cpu.pc, d.cpu.pc)
	fmt.Printf("R1:       0x%04X (%d)\n", d.cpu.r1, d.cpu.r1)
	fmt.Printf("R2:       0x%04X (%d)\n", d.cpu.r2, d.cpu.r2)
	fmt.Printf("R3:       0x%04X (%d)\n", d.cpu.r3, d.cpu.r3)
	fmt.Printf("R5:       0x%04X (%d)\n", d.cpu.r5, d.cpu.r5)
	fmt.Printf("R6:       0x%04X (%d)\n", d.cpu.r6, d.cpu.r6)
	fmt.Printf("R7:       0x%04X (%d)\n", d.cpu.r7, d.cpu.r7)
	fmt.Printf("Banking:  0x%04X\n", d.cpu.banking)
	fmt.Printf("Flags:    0x%04X ", d.cpu.flags)

	var flagStrs []string
	if d.cpu.flags.Zero() {
		flagStrs = append(flagStrs, "Z")
	}
	if d.cpu.flags.Negative() {
		flagStrs = append(flagStrs, "N")
	}
	if d.cpu.flags.Carry() {
		flagStrs = append(flagStrs, "C")
	}
	if d.cpu.flags.Overflow() {
		flagStrs = append(flagStrs, "O")
	}
	if d.cpu.flags.Interrupt() {
		flagStrs = append(flagStrs, "I")
	}
	if d.cpu.flags.Signed() {
		flagStrs = append(flagStrs, "S")
	}
	if d.cpu.flags.Halted() {
		flagStrs = append(flagStrs, "H")
	}

	if len(flagStrs) > 0 {
		fmt.Printf("[%s]", strings.Join(flagStrs, ","))
	}
	fmt.Printf("\n")
	fmt.Printf("Interrupt: 0x%04X\n", d.cpu.interrupt)
	fmt.Printf("Cycles:    %d\n", d.cpu.cycles)
	fmt.Printf("==============================\n\n")
}

// dumpMemory dumps a range of memory
func (d *DebugDevice) dumpMemory(start, length uint16) {
	fmt.Printf("\n=== Memory Dump at PC=%d ===\n", d.cpu.pc)
	fmt.Printf("Address range: 0x%04X - 0x%04X\n", start, start+length-1)
	for i := uint16(0); i < length; i++ {
		addr := start + i
		value := d.cpu.memory.ReadLocal(addr)
		if i%8 == 0 {
			fmt.Printf("\n0x%04X: ", addr)
		}
		fmt.Printf("%04X ", value)
	}
	fmt.Printf("\n==============================\n\n")
}

// logMessage logs a message ID
func (d *DebugDevice) logMessage(msgID uint16) {
	fmt.Printf("[DEBUG] Message ID %d (0x%04X) at PC=%d\n", msgID, msgID, d.cpu.pc)
}

// GetFailures returns all assertion failures
func (d *DebugDevice) GetFailures() []string {
	return d.failures
}

// HasFailures returns true if there are any failures
func (d *DebugDevice) HasFailures() bool {
	return len(d.failures) > 0
}

// GetFailureCount returns the number of failures
func (d *DebugDevice) GetFailureCount() int {
	return int(d.failCount)
}

// IsTraceEnabled returns true if tracing is enabled
func (d *DebugDevice) IsTraceEnabled() bool {
	return d.traceEnabled
}

// IsStaging returns true if a command is currently being staged (waiting for arguments)
func (d *DebugDevice) IsStaging() bool {
	return d.stagedCmd != 0
}
