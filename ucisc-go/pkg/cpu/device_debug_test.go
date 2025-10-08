package cpu

import (
	"testing"
)

func TestDebugDevice_Basic(t *testing.T) {
	cpu := New(1024)
	debug := cpu.GetDebugDevice()

	// Test device ID and type
	if debug.ID() != DebugDeviceID {
		t.Errorf("Expected device ID 0x%04X, got 0x%04X", DebugDeviceID, debug.ID())
	}
	if debug.Type() != DebugDeviceType {
		t.Errorf("Expected device type 0x%02X, got 0x%02X", DebugDeviceType, debug.Type())
	}
}

func TestDebugDevice_AssertEqual(t *testing.T) {
	cpu := New(1024)
	debug := cpu.GetDebugDevice()

	// Stage ASSERT_EQ command
	debug.WriteDebug(CmdAssertEq)
	debug.WriteDebug(42)
	debug.WriteDebug(42)
	debug.WriteDebug(CtrlExecute)

	// Should pass
	if debug.HasFailures() {
		t.Errorf("Assert equal should pass: %v", debug.GetFailures())
	}

	// Stage ASSERT_EQ command with different values
	debug.WriteDebug(CmdAssertEq)
	debug.WriteDebug(42)
	debug.WriteDebug(43)
	debug.WriteDebug(CtrlExecute)

	// Should fail
	if !debug.HasFailures() {
		t.Error("Assert equal should fail with different values")
	}
	if debug.Read(0x4) != 1 { // Failure count
		t.Errorf("Expected failure count 1, got %d", debug.Read(0x4))
	}
}

func TestDebugDevice_Reset(t *testing.T) {
	cpu := New(1024)
	debug := cpu.GetDebugDevice()

	// Create a failure
	debug.WriteDebug(CmdAssertEq)
	debug.WriteDebug(1)
	debug.WriteDebug(2)
	debug.WriteDebug(CtrlExecute)

	if !debug.HasFailures() {
		t.Fatal("Should have failure")
	}

	// Reset
	debug.WriteDebug(CtrlReset)

	// Check state is cleared
	if debug.HasFailures() {
		t.Error("Failures should be cleared after reset")
	}
	if debug.Read(0x3) != 0 { // Status
		t.Error("Status should be cleared")
	}
	if debug.Read(0x4) != 0 { // Failure count
		t.Error("Failure count should be cleared")
	}
}

func TestDebugDevice_ThroughCPU(t *testing.T) {
	cpu := New(1024)
	debug := cpu.GetDebugDevice()

	// Set r3 to 42
	cpu.setRegister(RegR3, 42)

	// Write to &debug through setRegister (simulates &debug <- copy ...)
	cpu.setRegister(RegFlags, CmdAssertEq)  // Stage command
	cpu.setRegister(RegFlags, cpu.GetR3())  // Stage arg1 (r3 value = 42)
	cpu.setRegister(RegFlags, 42)           // Stage arg2 (expected = 42)
	cpu.setRegister(RegFlags, CtrlExecute)  // Execute

	// Should pass
	if debug.HasFailures() {
		t.Errorf("Assert should pass: %v", debug.GetFailures())
	}
}

func TestDebugDevice_AllAssertions(t *testing.T) {
	tests := []struct {
		name    string
		cmd     uint16
		arg1    uint16
		arg2    uint16
		wantFail bool
	}{
		{"eq_pass", CmdAssertEq, 5, 5, false},
		{"eq_fail", CmdAssertEq, 5, 6, true},
		{"ne_pass", CmdAssertNe, 5, 6, false},
		{"ne_fail", CmdAssertNe, 5, 5, true},
		{"gt_pass", CmdAssertGt, 10, 5, false},
		{"gt_fail", CmdAssertGt, 5, 10, true},
		{"lt_pass", CmdAssertLt, 5, 10, false},
		{"lt_fail", CmdAssertLt, 10, 5, true},
		{"gte_pass_gt", CmdAssertGte, 10, 5, false},
		{"gte_pass_eq", CmdAssertGte, 5, 5, false},
		{"gte_fail", CmdAssertGte, 5, 10, true},
		{"lte_pass_lt", CmdAssertLte, 5, 10, false},
		{"lte_pass_eq", CmdAssertLte, 5, 5, false},
		{"lte_fail", CmdAssertLte, 10, 5, true},
		{"zero_pass", CmdAssertZero, 0, 0, false},
		{"zero_fail", CmdAssertZero, 1, 0, true},
		{"not_zero_pass", CmdAssertNotZero, 1, 0, false},
		{"not_zero_fail", CmdAssertNotZero, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := New(1024)
			debug := cpu.GetDebugDevice()

			debug.WriteDebug(tt.cmd)
			debug.WriteDebug(tt.arg1)
			if tt.cmd != CmdAssertZero && tt.cmd != CmdAssertNotZero {
				debug.WriteDebug(tt.arg2)
			}
			debug.WriteDebug(CtrlExecute)

			if tt.wantFail && !debug.HasFailures() {
				t.Errorf("Expected failure but passed")
			}
			if !tt.wantFail && debug.HasFailures() {
				t.Errorf("Expected pass but failed: %v", debug.GetFailures())
			}
		})
	}
}

func TestDebugDevice_SignedAssertions(t *testing.T) {
	cpu := New(1024)
	debug := cpu.GetDebugDevice()

	// -5 (0xFFFB) > -10 (0xFFF6) as signed
	debug.WriteDebug(CmdAssertGtSigned)
	debug.WriteDebug(0xFFFB) // -5
	debug.WriteDebug(0xFFF6) // -10
	debug.WriteDebug(CtrlExecute)

	if debug.HasFailures() {
		t.Errorf("Signed assertion failed: %v", debug.GetFailures())
	}

	// Reset
	debug.WriteDebug(CtrlReset)

	// -10 (0xFFF6) < -5 (0xFFFB) as signed
	debug.WriteDebug(CmdAssertLtSigned)
	debug.WriteDebug(0xFFF6) // -10
	debug.WriteDebug(0xFFFB) // -5
	debug.WriteDebug(CtrlExecute)

	if debug.HasFailures() {
		t.Errorf("Signed assertion failed: %v", debug.GetFailures())
	}
}

func TestDebugDevice_FlagAssertions(t *testing.T) {
	cpu := New(1024)
	debug := cpu.GetDebugDevice()

	// Set zero flag
	cpu.flags = cpu.flags.Set(FlagZero)

	// Assert has zero flag
	debug.WriteDebug(CmdAssertHasFlags)
	debug.WriteDebug(uint16(FlagZero))
	debug.WriteDebug(CtrlExecute)

	if debug.HasFailures() {
		t.Errorf("Should have zero flag: %v", debug.GetFailures())
	}

	// Reset
	debug.WriteDebug(CtrlReset)

	// Assert doesn't have carry flag
	debug.WriteDebug(CmdAssertNoFlags)
	debug.WriteDebug(uint16(FlagCarry))
	debug.WriteDebug(CtrlExecute)

	if debug.HasFailures() {
		t.Errorf("Should not have carry flag: %v", debug.GetFailures())
	}
}

func TestDebugDevice_BehaviorFlags(t *testing.T) {
	cpu := New(1024)
	debug := cpu.GetDebugDevice()

	// Set halt on fail
	debug.Write(0x2, BehaviorHaltOnFail)

	// Create a failure
	debug.WriteDebug(CmdAssertEq)
	debug.WriteDebug(1)
	debug.WriteDebug(2)
	debug.WriteDebug(CtrlExecute)

	// CPU should be halted
	if !cpu.flags.Halted() {
		t.Error("CPU should be halted on assertion failure")
	}
}
