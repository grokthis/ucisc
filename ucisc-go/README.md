# uCISC Go Emulator

A professional-quality Go implementation of the uCISC processor emulator.

## Features

- **Accurate Emulation**: Bit-accurate implementation of the uCISC instruction set
- **Complete Assembler**: Full assembly language with labels, all operations and addressing modes
- **High Performance**: ~28.8 MIPS average (verified with 13 comprehensive benchmarks)
- **Comprehensive Testing**: 300+ test cases covering all operations and edge cases
- **Performance Validation**: Automated benchmarks ensure >10 MIPS minimum requirement
- **Clean Architecture**: Modular, testable, and extensible design
- **Full Register Support**: All 6 registers (r1, r2, r3, r5, r6, r7) with comprehensive testing

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

| Benchmark | MIPS | Description |
|-----------|------|-------------|
| ConditionalBranching | 31.09 | Conditional execution with effect flags |
| Multiplication | 30.36 | Integer multiplication |
| BitwiseOperations | 30.19 | AND/OR/XOR/INV/SWAP operations |
| InstructionFetch | 29.93 | Raw fetch/decode overhead |
| MemoryAccess | 29.80 | Memory read/write operations |
| RegisterToRegister | 29.16 | Register-to-register copies |
| NestedLoops | 29.13 | Nested loops (10×10) |
| SimpleLoop | 28.74 | Simple counting loop |
| Fibonacci | 28.11 | Fibonacci sequence |
| ALUOperations | 27.72 | Mixed ALU operations |
| ShiftOperations | 27.46 | Left/right shifts |
| MemoryCopy | 27.31 | Memory copy operations |
| AllRegisters | 26.23 | All 6 registers |

**Average: ~28.8 MIPS** (exceeds 10 MIPS minimum by 2.9x)

## Documentation

- [Requirements](../docs/go-emulator/REQUIREMENTS.md)
- [Architecture](../docs/go-emulator/ARCHITECTURE.md)
- [uCISC Instruction Set](../docs/programming_guide/5_Instruction_Set_Details.md)

## License

See LICENSE file in the root of the ucisc repository.
