package testbench

import (
	"fmt"
	"time"

	"github.com/grokthis/ucisc-go/pkg/asm"
	"github.com/grokthis/ucisc-go/pkg/cpu"
)

// Executor defines the interface for running uCISC code
// This allows testing on emulator or real hardware
type Executor interface {
	// Load loads a program into memory
	Load(code []uint16) error

	// Run executes the loaded program until halt or timeout
	Run(timeout time.Duration) error

	// GetDebugDevice returns the debug device for assertion checking
	// Returns nil if not supported (e.g., on real hardware without debug device)
	GetDebugDevice() DebugDevice

	// Reset resets the execution environment
	Reset() error
}

// DebugDevice provides access to debug device state
type DebugDevice interface {
	// HasFailures returns true if any assertions failed
	HasFailures() bool

	// GetFailures returns all assertion failure messages
	GetFailures() []string

	// GetFailureCount returns the number of failures
	GetFailureCount() int
}

// EmulatorExecutor implements Executor for the CPU emulator
type EmulatorExecutor struct {
	cpu *cpu.CPU
}

// NewEmulatorExecutor creates a new emulator-based executor
func NewEmulatorExecutor(memSize int) *EmulatorExecutor {
	return &EmulatorExecutor{
		cpu: cpu.New(memSize),
	}
}

func (e *EmulatorExecutor) Load(code []uint16) error {
	return e.cpu.LoadProgram(code, 0)
}

func (e *EmulatorExecutor) Run(timeout time.Duration) error {
	// Create a channel for execution result
	done := make(chan error, 1)

	// Run in goroutine with timeout
	go func() {
		done <- e.cpu.Run()
	}()

	// Wait for completion or timeout
	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("execution timeout after %v", timeout)
	}
}

func (e *EmulatorExecutor) GetDebugDevice() DebugDevice {
	return e.cpu.GetDebugDevice()
}

func (e *EmulatorExecutor) Reset() error {
	e.cpu.Reset()
	return nil
}

// Test represents a single test case
type Test struct {
	Name        string
	Description string
	Source      string // uCISC assembly source code
	Timeout     time.Duration
}

// TestResult contains the results of running a test
type TestResult struct {
	Test     *Test
	Passed   bool
	Duration time.Duration
	Error    error
	Failures []string
}

// TestSuite is a collection of tests
type TestSuite struct {
	Name     string
	Tests    []*Test
	Executor Executor
}

// NewTestSuite creates a new test suite
func NewTestSuite(name string, executor Executor) *TestSuite {
	return &TestSuite{
		Name:     name,
		Tests:    make([]*Test, 0),
		Executor: executor,
	}
}

// AddTest adds a test to the suite
func (ts *TestSuite) AddTest(test *Test) {
	// Set default timeout if not specified
	if test.Timeout == 0 {
		test.Timeout = 5 * time.Second
	}
	ts.Tests = append(ts.Tests, test)
}

// Run runs all tests in the suite
func (ts *TestSuite) Run() *TestSuiteResult {
	result := &TestSuiteResult{
		SuiteName: ts.Name,
		Results:   make([]*TestResult, 0),
		StartTime: time.Now(),
	}

	for _, test := range ts.Tests {
		testResult := ts.runTest(test)
		result.Results = append(result.Results, testResult)

		if testResult.Passed {
			result.Passed++
		} else {
			result.Failed++
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result
}

// runTest runs a single test
func (ts *TestSuite) runTest(test *Test) *TestResult {
	startTime := time.Now()
	result := &TestResult{
		Test: test,
	}

	// Reset executor before each test
	if err := ts.Executor.Reset(); err != nil {
		result.Error = fmt.Errorf("reset failed: %w", err)
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Assemble the test code
	code, err := asm.Assemble(test.Source)
	if err != nil {
		result.Error = fmt.Errorf("assembly failed: %w", err)
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Load the program
	if err := ts.Executor.Load(code); err != nil {
		result.Error = fmt.Errorf("load failed: %w", err)
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Run the program
	if err := ts.Executor.Run(test.Timeout); err != nil {
		result.Error = fmt.Errorf("execution failed: %w", err)
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Check debug device for assertion failures
	debugDevice := ts.Executor.GetDebugDevice()
	if debugDevice != nil && debugDevice.HasFailures() {
		result.Passed = false
		result.Failures = debugDevice.GetFailures()
	} else {
		result.Passed = true
	}

	result.Duration = time.Since(startTime)
	return result
}

// TestSuiteResult contains the results of running a test suite
type TestSuiteResult struct {
	SuiteName string
	Results   []*TestResult
	Passed    int
	Failed    int
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

// Print prints the test results in a readable format
func (r *TestSuiteResult) Print() {
	fmt.Printf("\n=== Test Suite: %s ===\n", r.SuiteName)
	fmt.Printf("Started: %s\n", r.StartTime.Format("15:04:05"))
	fmt.Printf("Duration: %v\n\n", r.Duration)

	for i, testResult := range r.Results {
		status := "✓ PASS"
		if !testResult.Passed {
			status = "✗ FAIL"
		}

		fmt.Printf("[%d] %s - %s (%v)\n", i+1, status, testResult.Test.Name, testResult.Duration)

		if testResult.Test.Description != "" {
			fmt.Printf("    %s\n", testResult.Test.Description)
		}

		if testResult.Error != nil {
			fmt.Printf("    Error: %v\n", testResult.Error)
		}

		if len(testResult.Failures) > 0 {
			fmt.Printf("    Assertion Failures:\n")
			for _, failure := range testResult.Failures {
				fmt.Printf("      - %s\n", failure)
			}
		}

		fmt.Println()
	}

	fmt.Printf("=== Summary ===\n")
	fmt.Printf("Total:  %d\n", len(r.Results))
	fmt.Printf("Passed: %d\n", r.Passed)
	fmt.Printf("Failed: %d\n", r.Failed)

	if r.Failed == 0 {
		fmt.Printf("\n✓ All tests passed!\n\n")
	} else {
		fmt.Printf("\n✗ %d test(s) failed\n\n", r.Failed)
	}
}

// HasFailures returns true if any tests failed
func (r *TestSuiteResult) HasFailures() bool {
	return r.Failed > 0
}
