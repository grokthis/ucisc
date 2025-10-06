# uCISC Go Emulator - Architecture Design Document

## 1. Overview

This document describes the high-level architecture of the uCISC Go emulator. The design emphasizes modularity, testability, and extensibility while maintaining high performance and accuracy.

## 2. Design Principles

### 2.1 Core Principles

1. **Separation of Concerns**: Each package has a single, well-defined responsibility
2. **Interface-Driven Design**: Components interact through well-defined interfaces
3. **Testability**: All components can be tested in isolation
4. **Performance**: Hot paths are optimized, with minimal allocations
5. **Idiomatic Go**: Follow Go best practices and conventions
6. **Zero-Cost Abstractions**: Interfaces should compile away or have negligible overhead

### 2.2 SOLID Principles

- **Single Responsibility**: Each type has one reason to change
- **Open/Closed**: Open for extension (new devices), closed for modification
- **Liskov Substitution**: All device implementations are interchangeable
- **Interface Segregation**: Small, focused interfaces
- **Dependency Inversion**: Depend on abstractions, not concretions

## 3. System Architecture

### 3.1 High-Level Component Diagram

```
┌─────────────────────────────────────────────────────────┐
│                     CLI Application                      │
│                    (cmd/ucisc)                          │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                   Debugger / Controller                  │
│                   (pkg/debugger)                        │
│  • Breakpoints  • Stepping  • Inspection  • Tracing     │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                      CPU Core                           │
│                     (pkg/cpu)                           │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌──────────┐  │
│  │Registers│  │  ALU    │  │ Decoder │  │ Executor │  │
│  └─────────┘  └─────────┘  └─────────┘  └──────────┘  │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                    Memory System                         │
│                   (pkg/cpu/memory)                      │
│  ┌──────────────┐           ┌─────────────────────┐    │
│  │Local Memory  │           │   Device Memory     │    │
│  │   (64K)      │           │  (Control + Block)  │    │
│  └──────────────┘           └─────────────────────┘    │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                    Device Manager                        │
│                   (pkg/device)                          │
│  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐    │
│  │ GPIO │  │ UART │  │ I2C  │  │Block │  │Video │    │
│  └──────┘  └──────┘  └──────┘  └──────┘  └──────┘    │
└─────────────────────────────────────────────────────────┘

Supporting Components:
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│   Loader     │  │ Disassembler │  │   Utilities  │
│ (pkg/loader) │  │(pkg/disasm)  │  │ (pkg/util)   │
└──────────────┘  └──────────────┘  └──────────────┘
```

### 3.2 Data Flow

#### 3.2.1 Execution Flow
```
1. Load Program (Loader) → Local Memory
2. Fetch Instruction (CPU) → from PC address
3. Decode Instruction (Decoder) → instruction structure
4. Execute (Executor):
   - Read source value (Memory System)
   - Perform ALU operation
   - Evaluate condition (Flags)
   - Write destination (Memory System)
   - Update PC
5. Repeat or Halt
```

#### 3.2.2 Memory Access Flow
```
CPU wants to read/write at address X using register R:
1. Check banking register for R
2. If banked:
   - Extract device ID from address
   - Route to Device Manager
   - Device handles read/write
3. If not banked:
   - Access local memory directly
4. Return value to CPU
```

#### 3.2.3 Debug Flow
```
1. User sets breakpoint → Debugger
2. CPU executes → check breakpoint before each instruction
3. Hit breakpoint → pause execution
4. User inspects/steps → Debugger controls CPU
5. User continues → CPU resumes execution
```

## 4. Core Components

### 4.1 CPU Package (`pkg/cpu`)

#### 4.1.1 CPU Core (`cpu.go`)

```go
package cpu

// CPU represents the uCISC processor
type CPU struct {
    // Registers
    r1, r2, r3, r4, r5, r6 uint16  // General purpose
    pc                      uint16  // Program counter
    banking                 uint16  // Banking control
    flags                   uint16  // Processor flags
    interrupt               uint16  // Interrupt register

    // Memory
    memory *Memory

    // Execution state
    halted bool
    cycles uint64

    // Callbacks for debugging
    beforeExec BeforeExecFunc
    afterExec  AfterExecFunc
}

// BeforeExecFunc is called before each instruction
type BeforeExecFunc func(cpu *CPU, addr uint16) (shouldExecute bool)

// AfterExecFunc is called after each instruction
type AfterExecFunc func(cpu *CPU, instruction uint32)

// Core methods
func New(memSize int) *CPU
func (c *CPU) Reset()
func (c *CPU) Step() error           // Execute one instruction
func (c *CPU) Run() error            // Run until halt
func (c *CPU) LoadProgram(data []uint16, offset uint16) error

// Register access
func (c *CPU) GetRegister(reg Register) uint16
func (c *CPU) SetRegister(reg Register, value uint16)
func (c *CPU) GetPC() uint16
func (c *CPU) GetFlags() Flags
```

**Responsibilities**:
- Manage register state
- Orchestrate fetch-decode-execute cycle
- Track execution statistics (cycles, instructions)
- Provide hooks for debugging

**Design Notes**:
- Registers stored as uint16 for direct manipulation
- Callbacks allow debugger to intercept execution
- No direct device knowledge (delegates to Memory)

#### 4.1.2 ALU (`alu.go`)

```go
package cpu

// ALU performs arithmetic and logic operations
type ALU struct {
    signed bool  // Current signed/unsigned mode
}

// Operation represents an ALU operation code
type Operation uint8

const (
    OpCopy Operation = iota
    OpAnd
    OpOr
    OpXor
    OpInv
    OpShl
    OpShr
    OpSwap
    OpMsb
    OpLsb
    OpAdd
    OpAddc
    OpSub
    OpMult
    OpMultsw
)

// Result contains the operation result and flags
type Result struct {
    Value    uint16
    Zero     bool
    Negative bool
    Carry    bool
    Overflow bool
}

func (a *ALU) Execute(op Operation, dest, src uint16, carryIn bool) Result
```

**Responsibilities**:
- Perform all arithmetic and logic operations
- Calculate flags (Z, N, C, O)
- Handle signed/unsigned modes

**Design Notes**:
- Pure functions (no state mutation except signed flag)
- Returns structured Result for clarity
- Optimized for common operations (copy, add, sub)

#### 4.1.3 Instruction Decoder (`decoder.go`)

```go
package cpu

// Instruction represents a decoded uCISC instruction
type Instruction struct {
    Raw         uint32
    Source      Source
    Destination Destination
    Increment   bool
    Effect      Effect
    Operation   Operation
    Immediate   int16  // Signed immediate
    Offset      uint8  // For memory destinations
}

// Source represents instruction source
type Source struct {
    Type   SourceType
    Reg    Register
    Signed bool
}

type SourceType uint8

const (
    SourceReg SourceType = iota  // Register value
    SourceMem                     // Memory at register
    SourceVal                     // Immediate value
)

// Destination represents instruction destination
type Destination struct {
    Type DestType
    Reg  Register
}

type DestType uint8

const (
    DestReg DestType = iota  // Register
    DestMem                   // Memory at register
)

// Effect represents conditional execution
type Effect uint8

const (
    EffectZero Effect = iota
    EffectNotZero
    EffectNegative
    EffectPositive
    EffectAlways
    EffectOverflow
    EffectInterrupt
    EffectFlagsOnly
)

func Decode(raw uint32) (*Instruction, error)
```

**Responsibilities**:
- Decode 32-bit instruction into structured format
- Validate instruction encoding
- Extract all fields (source, dest, op, immediate, etc.)

**Design Notes**:
- Stateless decoder (pure function)
- Returns rich instruction structure
- Separate types for source/destination clarity

#### 4.1.4 Executor (`executor.go`)

```go
package cpu

// Executor handles instruction execution
type Executor struct {
    cpu *CPU
    alu *ALU
}

func NewExecutor(cpu *CPU) *Executor

// Execute runs a decoded instruction
func (e *Executor) Execute(inst *Instruction) error {
    // 1. Read source value
    srcValue := e.readSource(inst.Source, inst.Immediate)

    // 2. Read destination value (for 2-arg operations)
    destValue := e.readDestination(inst.Destination, inst.Offset)

    // 3. Perform ALU operation
    result := e.alu.Execute(inst.Operation, destValue, srcValue, e.cpu.flags.Carry())

    // 4. Check effect condition
    if e.shouldStore(inst.Effect, result) {
        e.writeDestination(inst.Destination, inst.Offset, result.Value, inst.Increment)
        e.cpu.updateFlags(result)
    }

    // 5. Update PC (unless it was the destination)
    if inst.Destination.Reg != RegPC {
        e.cpu.pc += 2  // Each instruction is 2 words
    }

    return nil
}

func (e *Executor) readSource(src Source, immediate int16) uint16
func (e *Executor) readDestination(dest Destination, offset uint8) uint16
func (e *Executor) writeDestination(dest Destination, offset uint8, value uint16, increment bool)
func (e *Executor) shouldStore(effect Effect, result Result) bool
```

**Responsibilities**:
- Coordinate instruction execution
- Handle source/destination value resolution
- Apply increment/decrement logic
- Evaluate conditional effects
- Update processor state

**Design Notes**:
- Encapsulates complex execution logic
- Handles banking through Memory abstraction
- Clear separation between read/compute/write phases

### 4.2 Memory System (`pkg/cpu/memory.go`)

```go
package cpu

// Memory manages local and device memory
type Memory struct {
    local  []uint16         // Local memory (up to 64K words)
    device *device.Manager  // Device manager
}

func NewMemory(size int, deviceMgr *device.Manager) *Memory

// Read reads from memory, handling banking
func (m *Memory) Read(addr uint16, reg Register, banking uint16) uint16

// Write writes to memory, handling banking
func (m *Memory) Write(addr uint16, reg Register, banking uint16, value uint16)

// Direct local memory access (for PC, loading programs)
func (m *Memory) ReadLocal(addr uint16) uint16
func (m *Memory) WriteLocal(addr uint16, value uint16)

// Banking helpers
func (m *Memory) isBanked(reg Register, banking uint16) bool
func (m *Memory) routeToDevice(addr uint16, isRead bool) (deviceID uint8, offset uint16)
```

**Responsibilities**:
- Manage local memory array
- Implement banking logic
- Route device memory accesses to Device Manager
- Provide direct access for program loading

**Design Notes**:
- Encapsulates all banking complexity
- Simple interface for CPU (doesn't need to know about banking)
- Local memory is just a []uint16 for performance

### 4.3 Device System (`pkg/device`)

#### 4.3.1 Device Interface (`device.go`)

```go
package device

// Device represents any uCISC device
type Device interface {
    // ID returns the device ID (0-255)
    ID() uint8

    // Type returns the device type code
    Type() uint16

    // ReadControl reads from device control memory (0-15)
    ReadControl(offset uint8) uint16

    // WriteControl writes to device control memory (0-15)
    WriteControl(offset uint8, value uint16)

    // ReadBlock reads from device block memory (0-255)
    // Returns 0 if device doesn't support block memory
    ReadBlock(offset uint8) uint16

    // WriteBlock writes to device block memory (0-255)
    // Ignored if device doesn't support block memory
    WriteBlock(offset uint8, value uint16)

    // Reset resets the device to initial state
    Reset()

    // Tick allows device to update (called periodically)
    Tick()
}

// Manager manages all devices
type Manager struct {
    devices [256]Device
}

func NewManager() *Manager

func (m *Manager) Register(dev Device) error
func (m *Manager) Unregister(id uint8)
func (m *Manager) Get(id uint8) Device
func (m *Manager) ReadControl(deviceID uint8, offset uint8) uint16
func (m *Manager) WriteControl(deviceID uint8, offset uint8, value uint16)
func (m *Manager) ReadBlock(deviceID uint8, offset uint8) uint16
func (m *Manager) WriteBlock(deviceID uint8, offset uint8, value uint16)
func (m *Manager) TickAll()
```

**Design Notes**:
- Simple interface allows any device implementation
- Manager provides device registry
- Tick() allows devices to have autonomous behavior
- ID embedded in interface ensures device knows its position

#### 4.3.2 Base Device (`base.go`)

```go
package device

// BaseDevice provides common device functionality
type BaseDevice struct {
    id      uint8
    devType uint16
    control [16]uint16
    block   [256]uint16
    hasBlock bool
}

func NewBaseDevice(id uint8, devType uint16, hasBlock bool) *BaseDevice

func (b *BaseDevice) ID() uint8 { return b.id }
func (b *BaseDevice) Type() uint16 { return b.devType }

func (b *BaseDevice) ReadControl(offset uint8) uint16 {
    if offset == 0 {
        return b.id  // Device ID
    }
    if offset == 1 {
        return b.devType  // Device type
    }
    return b.control[offset]
}

func (b *BaseDevice) WriteControl(offset uint8, value uint16) {
    if offset < 2 {
        return  // ID and type are read-only
    }
    b.control[offset] = value
}

func (b *BaseDevice) ReadBlock(offset uint8) uint16 {
    if !b.hasBlock { return 0 }
    return b.block[offset]
}

func (b *BaseDevice) WriteBlock(offset uint8, value uint16) {
    if !b.hasBlock { return }
    b.block[offset] = value
}

func (b *BaseDevice) Reset() {
    b.control = [16]uint16{}
    b.block = [256]uint16{}
}

func (b *BaseDevice) Tick() {}
```

**Design Notes**:
- Provides default implementation for common functionality
- Devices can embed BaseDevice and override as needed
- Handles read-only ID and type fields automatically

#### 4.3.3 Example Device: GPIO (`gpio.go`)

```go
package device

const DeviceTypeGPIO uint16 = 0x0001

// GPIO represents a GPIO device
type GPIO struct {
    *BaseDevice
    pins uint16  // Pin states
    direction uint16  // Pin directions (0=in, 1=out)
}

func NewGPIO(id uint8) *GPIO {
    return &GPIO{
        BaseDevice: NewBaseDevice(id, DeviceTypeGPIO, false),
    }
}

func (g *GPIO) ReadControl(offset uint8) uint16 {
    switch offset {
    case 2:  // Pin states
        return g.pins
    case 3:  // Pin directions
        return g.direction
    default:
        return g.BaseDevice.ReadControl(offset)
    }
}

func (g *GPIO) WriteControl(offset uint8, value uint16) {
    switch offset {
    case 2:  // Write to pins (output only)
        g.pins = (g.pins & ^g.direction) | (value & g.direction)
    case 3:  // Set direction
        g.direction = value
    default:
        g.BaseDevice.WriteControl(offset, value)
    }
}

func (g *GPIO) SetPin(pin uint8, value bool) {
    // External input (for testing/simulation)
    mask := uint16(1) << pin
    if value {
        g.pins |= mask
    } else {
        g.pins &= ^mask
    }
}
```

### 4.4 Debugger (`pkg/debugger`)

```go
package debugger

// Debugger provides debugging capabilities
type Debugger struct {
    cpu        *cpu.CPU
    breakpoints map[uint16]*Breakpoint
    watchpoints map[uint16]*Watchpoint
    traceEnabled bool
    traceWriter io.Writer
}

// Breakpoint represents an execution breakpoint
type Breakpoint struct {
    Address  uint16
    Enabled  bool
    Condition BreakCondition  // Optional condition
}

// BreakCondition is evaluated when breakpoint is hit
type BreakCondition func(cpu *cpu.CPU) bool

// Watchpoint represents a memory watchpoint
type Watchpoint struct {
    Address uint16
    OnRead  bool
    OnWrite bool
    Enabled bool
}

func New(cpu *cpu.CPU) *Debugger

// Breakpoint management
func (d *Debugger) AddBreakpoint(addr uint16) *Breakpoint
func (d *Debugger) RemoveBreakpoint(addr uint16)
func (d *Debugger) AddConditionalBreakpoint(addr uint16, cond BreakCondition) *Breakpoint

// Watchpoint management
func (d *Debugger) AddWatchpoint(addr uint16, onRead, onWrite bool) *Watchpoint
func (d *Debugger) RemoveWatchpoint(addr uint16)

// Execution control
func (d *Debugger) Step() error                    // Single step
func (d *Debugger) StepOver() error               // Step over function calls
func (d *Debugger) Continue() error               // Continue execution
func (d *Debugger) RunTo(addr uint16) error      // Run to address
func (d *Debugger) Reset()                        // Reset CPU

// Inspection
func (d *Debugger) DumpRegisters() string
func (d *Debugger) DumpMemory(start, end uint16) string
func (d *Debugger) Disassemble(start uint16, count int) string

// Tracing
func (d *Debugger) EnableTrace(w io.Writer)
func (d *Debugger) DisableTrace()
```

**Design Notes**:
- Uses CPU callbacks for breakpoint detection
- Conditional breakpoints allow flexible debugging
- Watchpoints hook into memory read/write
- Trace output configurable (file, stdout, etc.)

### 4.5 Loader (`pkg/loader`)

```go
package loader

// Loader loads programs into memory
type Loader struct{}

// Load loads a program from various formats
func Load(filename string) ([]uint16, error)

// LoadIntelHex loads Intel HEX format
func LoadIntelHex(r io.Reader) ([]uint16, error)

// LoadBinary loads raw binary (16-bit words)
func LoadBinary(r io.Reader) ([]uint16, error)

// LoadHex loads custom hex format (space-separated hex)
func LoadHex(r io.Reader) ([]uint16, error)
```

**Design Notes**:
- Support multiple formats
- Returns []uint16 for loading into memory
- Format auto-detection based on file extension/content

### 4.6 Disassembler (`pkg/disasm`)

```go
package disasm

// Disassembler converts instructions to assembly
type Disassembler struct {
    symbols map[uint16]string  // Optional symbol table
}

func New() *Disassembler

// LoadSymbols loads a symbol table from file
func (d *Disassembler) LoadSymbols(filename string) error

// Disassemble converts a single instruction
func (d *Disassembler) Disassemble(addr uint16, inst uint32) string

// DisassembleRange disassembles a range of memory
func (d *Disassembler) DisassembleRange(memory []uint16, start, count int) string

// Format example: "0x0100: r1/5 <- add r2/3"
```

**Design Notes**:
- Optional symbol support for readability
- Human-readable output matching uCISC syntax
- Can work standalone or with debugger

## 5. Key Algorithms

### 5.1 Instruction Execution

```
function Execute(instruction):
    decoded = Decode(instruction)

    // Read source
    if decoded.source.type == SourceVal:
        srcValue = decoded.immediate
    else if decoded.source.type == SourceReg:
        srcValue = GetRegister(decoded.source.reg) + decoded.immediate
    else:  // SourceMem
        addr = GetRegister(decoded.source.reg) + decoded.immediate
        srcValue = Memory.Read(addr, decoded.source.reg, banking)
        if decoded.increment and decoded.dest.type != DestMem:
            // Pop: increment source after read
            SetRegister(decoded.source.reg, GetRegister(decoded.source.reg) + 1)

    // Read destination (for 2-arg ops)
    if decoded.dest.type == DestReg:
        destValue = GetRegister(decoded.dest.reg)
    else:  // DestMem
        baseAddr = GetRegister(decoded.dest.reg)
        if decoded.increment:
            // Push: decrement before write
            baseAddr = baseAddr - 1
            SetRegister(decoded.dest.reg, baseAddr)
        addr = baseAddr + decoded.offset
        destValue = Memory.Read(addr, decoded.dest.reg, banking)

    // Execute ALU operation
    result = ALU.Execute(decoded.operation, destValue, srcValue)

    // Check effect condition
    if ShouldStore(decoded.effect, result, flags):
        if decoded.dest.type == DestReg:
            SetRegister(decoded.dest.reg, result.value)
        else:
            addr = GetRegister(decoded.dest.reg) + decoded.offset
            Memory.Write(addr, decoded.dest.reg, banking, result.value)

        if decoded.operation != OpCopy:
            UpdateFlags(result)

    // Update PC (unless PC was destination)
    if decoded.dest.reg != RegPC:
        pc = pc + 2
```

### 5.2 Banking Resolution

```
function Memory.Read(addr, reg, banking):
    if IsBanked(reg, banking):
        deviceID = addr >> 4  // Top 8 bits for control, different for block
        if addr < 0x1000:
            // Control memory: deviceID = addr / 16
            offset = addr & 0x0F
            return DeviceManager.ReadControl(deviceID, offset)
        else:
            // Block memory: deviceID = (addr - 0x1000) / 256 + 16
            deviceID = ((addr - 0x1000) >> 8) + 16
            offset = addr & 0xFF
            return DeviceManager.ReadBlock(deviceID, offset)
    else:
        return localMemory[addr]
```

### 5.3 Breakpoint Detection

```
function CheckBreakpoint(addr):
    bp = breakpoints[addr]
    if bp == nil or not bp.enabled:
        return false

    if bp.condition != nil:
        return bp.condition(cpu)

    return true

// In CPU.Step():
    if beforeExec != nil and not beforeExec(cpu, pc):
        return  // Breakpoint hit, don't execute

    Execute(currentInstruction)

    if afterExec != nil:
        afterExec(cpu, currentInstruction)
```

## 6. Performance Considerations

### 6.1 Hot Path Optimizations

**Instruction Execution** (executed millions of times):
- Avoid allocations in execute path
- Use array lookups instead of map lookups where possible
- Inline small functions
- Pre-decode instruction formats for common patterns

**Memory Access**:
- Local memory is simple array access
- Banking check is bit test (very fast)
- Device access uses array lookup, not map

**ALU Operations**:
- Use lookup tables for flag calculations where beneficial
- Optimize common operations (copy, add, sub)
- Signed/unsigned handled with type conversions, not conditionals

### 6.2 Memory Layout

```go
// Prefer this (better cache locality):
type CPU struct {
    // Hot fields first
    pc      uint16
    flags   uint16
    r1, r2, r3 uint16  // Most commonly used

    // Less frequently accessed
    r4, r5, r6 uint16
    banking    uint16
    interrupt  uint16

    // Pointers last (may be nil during init)
    memory    *Memory
    callbacks *Callbacks
}
```

### 6.3 Profiling Points

Key areas to profile:
1. Instruction decode (should be < 5% of execution time)
2. ALU operations (should be < 10% of execution time)
3. Memory access (should be < 20% of execution time)
4. Banking logic (should be < 5% of execution time)
5. Device access (acceptable to be slower, not in hot path)

## 7. Error Handling

### 7.1 Error Types

```go
package cpu

type ErrorType uint8

const (
    ErrInvalidInstruction ErrorType = iota
    ErrInvalidMemoryAccess
    ErrInvalidRegister
    ErrInvalidDevice
    ErrHalted
)

type CPUError struct {
    Type ErrorType
    Addr uint16
    Msg  string
}

func (e *CPUError) Error() string
```

### 7.2 Error Strategy

- **Fatal Errors**: Invalid instruction, memory corruption → stop execution
- **Recoverable Errors**: Device errors, missing device → return error value (0x0000)
- **Debug Errors**: Breakpoint hit → special return value, not an error

## 8. Testing Strategy

### 8.1 Unit Tests

Each component tested in isolation:
- ALU: Test each operation with edge cases
- Decoder: Test all instruction formats
- Memory: Test banking logic
- Devices: Test control/block memory

### 8.2 Integration Tests

```go
// Example integration test
func TestSimpleProgram(t *testing.T) {
    cpu := cpu.New(65536)

    // Load program: r1 <- 5; r1 <- r1 + 3
    program := []uint16{
        0x5400, 0x0005,  // r1 <- copy val/5
        0x5A55, 0x0003,  // r1 <- add val/3
    }
    cpu.LoadProgram(program, 0)

    cpu.Step()  // Execute first instruction
    assert.Equal(t, uint16(5), cpu.GetRegister(cpu.R1))

    cpu.Step()  // Execute second instruction
    assert.Equal(t, uint16(8), cpu.GetRegister(cpu.R1))
}
```

### 8.3 Validation Tests

Run all programs in `examples/` and compare:
- Final register states
- Final memory states
- Execution traces
- Against ucisc-kotlin emulator

## 9. Build and Deployment

### 9.1 Build Configuration

```makefile
# Makefile
.PHONY: build test bench lint

build:
	go build -o bin/ucisc ./cmd/ucisc

build-release:
	go build -ldflags="-s -w" -o bin/ucisc ./cmd/ucisc

test:
	go test -v -race ./...

bench:
	go test -bench=. -benchmem ./...

lint:
	golangci-lint run

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
```

### 9.2 CI/CD Pipeline

```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.21', '1.22']
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go }}
      - run: go test -v -race ./...
      - run: go test -bench=. ./...

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: golangci/golangci-lint-action@v3
```

## 10. Future Extensibility

### 10.1 Plugin System (Future)

```go
// Device plugins loaded at runtime
type DevicePlugin interface {
    Device
    Name() string
    Version() string
}

func LoadPlugin(path string) (DevicePlugin, error)
```

### 10.2 JIT Compilation (Future)

```go
// JIT compiler for hot code paths
type JITCompiler struct {
    cpu *CPU
    cache map[uint16]CompiledCode
}

type CompiledCode func(*CPU) error

func (j *JITCompiler) Compile(addr uint16) CompiledCode
```

### 10.3 Multi-Processor (Future)

```go
// System with multiple CPUs
type System struct {
    cpus []*CPU
    bus  *Bus
}

func (s *System) RunParallel()
```

## 11. Open Design Questions

1. **Interrupt Handling**: How should interrupts be triggered and handled?
   - Option A: Timer-based (every N cycles)
   - Option B: Event-based (device signals)
   - **Recommendation**: Both, configurable

2. **Device Timing**: Should devices run in lockstep or asynchronously?
   - Option A: Tick all devices every CPU cycle (accurate but slow)
   - Option B: Tick devices periodically (fast but less accurate)
   - **Recommendation**: Option B with configurable tick rate

3. **Memory Model**: Should we support memory-mapped I/O beyond devices?
   - Option A: Strict device model only
   - Option B: Allow custom memory regions
   - **Recommendation**: Option A for v1.0, Option B later

4. **Concurrency**: Should the emulator support concurrent execution?
   - Option A: Single-threaded only (simple)
   - Option B: Concurrent with mutex protection (complex but flexible)
   - **Recommendation**: Option A for v1.0, with mutex in Memory for future

---

**Document Version**: 1.0
**Last Updated**: 2025-10-05
**Author**: Principal Engineer Review
**Status**: Draft for Review
