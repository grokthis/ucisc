# uCISC Go Emulator

A professional-quality Go implementation of the uCISC processor emulator.

## Features

- **Accurate Emulation**: Bit-accurate implementation of the uCISC instruction set
- **Full Device Support**: GPIO, UART, I2C, Block Memory, and Video devices
- **Comprehensive Debugging**: Breakpoints, watchpoints, step execution, and tracing
- **High Performance**: Optimized execution (>1M instructions/second)
- **Clean Architecture**: Modular, testable, and extensible design

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

# Run benchmarks
go test -bench=. ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Documentation

- [Requirements](../docs/go-emulator/REQUIREMENTS.md)
- [Architecture](../docs/go-emulator/ARCHITECTURE.md)
- [uCISC Instruction Set](../docs/programming_guide/5_Instruction_Set_Details.md)

## License

See LICENSE file in the root of the ucisc repository.
