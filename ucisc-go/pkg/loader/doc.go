// Package loader provides utilities for loading uCISC (Micro-CISC) programs from various file formats.
//
// The loader supports multiple program formats for the uCISC architecture,
// making it easy to load pre-assembled programs into memory for execution.
//
// # Supported Formats
//
//   - Hex format: Space-separated 16-bit hex values with optional comments
//   - Binary format: Raw 16-bit words in little-endian format
//
// # Hex Format
//
// The hex format is human-readable and supports comments for documentation:
//
//	5400 0064  # r1 <- copy val/100
//	5640 0000  # r2 <- copy &r1
//	0040 0000  # halt
//
// Each line can contain:
//   - Hex values (4 hex digits per 16-bit word)
//   - Comments starting with #
//   - Blank lines (ignored)
//
// # Example Usage
//
//	// Load from file (auto-detects .hex extension)
//	program, err := loader.Load("program.hex")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Load hex format explicitly
//	file, _ := os.Open("program.hex")
//	program, err := loader.LoadHex(file)
//
//	// Load binary format
//	file, _ := os.Open("program.bin")
//	program, err := loader.LoadBinary(file)
package loader
