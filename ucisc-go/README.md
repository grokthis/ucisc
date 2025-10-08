# uCISC Go Emulator

A professional-quality Go implementation of the uCISC processor emulator.

## Features

- **Accurate Emulation**: Bit-accurate implementation of the uCISC instruction set
- **Complete Assembler**: Full assembly language with labels, all operations and addressing modes
- **High Performance**: ~55 MIPS average (verified with 13 comprehensive benchmarks)
- **Zero Allocations**: Optimized hot path with array-based instruction cache
- **Comprehensive Testing**: 300+ test cases covering all operations and edge cases
- **Performance Validation**: Automated benchmarks ensure >50 MIPS across all tests
- **Clean Architecture**: Modular, testable, and extensible design
- **Full Register Support**: All 6 registers (r1, r2, r3, r5, r6, r7) with comprehensive testing
- **Debug Device**: Hardware-emulated debug device with assertions and register dumps

## Installation

```bash
go install github.com/grokthis/ucisc-go/cmd/ucisc@latest
```

## Quick Start

```bash
# Run a program
ucisc run examples/fib.hex

# Debug a program
ucisc debug examples/fib.hex

# Compile and run (requires ucisc compiler)
ucisc compile examples/fib.ucisc
ucisc run examples/fib.hex
```

## Building from Source

```bash
git clone https://github.com/grokthis/ucisc-go.git
cd ucisc-go
go build -o ucisc ./cmd/ucisc
```

## Architecture

The emulator is organized into several key packages:

- `pkg/cpu`: Core processor implementation (ALU, decoder, executor)
- `pkg/device`: Device framework and standard devices
- `pkg/debugger`: Debugging capabilities
- `pkg/loader`: Program loading (hex files, binary)
- `pkg/disasm`: Disassembler

See [docs/ARCHITECTURE.md](../docs/go-emulator/ARCHITECTURE.md) for detailed architecture documentation.

## Testing

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run benchmarks (reports IPS/MIPS)
go test -bench=. ./pkg/cpu -benchtime=1000x

# Run specific benchmark
go test -bench=BenchmarkCPU_SimpleLoop ./pkg/cpu -benchtime=1000x

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Performance Benchmarks

The emulator includes 13 comprehensive benchmarks that validate performance and correctness:

| Benchmark | MIPS | Time/op | Allocs/op | Description |
|-----------|------|---------|-----------|-------------|
| BitwiseOperations | 58.92 | 13.6 μs | 0 | AND/OR/XOR/INV/SWAP operations |
| Multiplication | 58.15 | 6.9 μs | 0 | Integer multiplication |
| ALUOperations | 57.90 | 13.9 μs | 0 | Mixed ALU operations |
| SimpleLoop | 55.95 | 5.4 μs | 0 | Simple counting loop |
| ConditionalBranching | 55.59 | 4.5 μs | 0 | Conditional execution with effect flags |
| ShiftOperations | 55.54 | 9.0 μs | 0 | Left/right shifts |
| InstructionFetch | 55.20 | 53.5 μs | 0 | Raw fetch/decode overhead |
| NestedLoops | 55.16 | 6.2 μs | 0 | Nested loops (10×10) |
| MemoryAccess | 53.03 | 4.8 μs | 0 | Memory read/write operations |
| AllRegisters | 52.61 | 1.2 ms | 0 | All 6 registers |
| RegisterToRegister | 52.13 | 13.5 μs | 0 | Register-to-register copies |
| MemoryCopy | 50.96 | 6.9 μs | 0 | Memory copy operations |
| Fibonacci | 50.30 | 2.4 μs | 0 | Fibonacci sequence |

**Average: ~55 MIPS** (exceeds 50 MIPS target across all benchmarks)

All benchmarks achieve zero allocations through array-based instruction caching.

## Documentation

- [Requirements](../docs/go-emulator/REQUIREMENTS.md)
- [Architecture](../docs/go-emulator/ARCHITECTURE.md)
- [uCISC Instruction Set](../docs/programming_guide/5_Instruction_Set_Details.md)

## License

See LICENSE file in the root of the ucisc repository.
