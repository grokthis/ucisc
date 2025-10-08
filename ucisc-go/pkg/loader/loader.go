package loader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Load loads a program from a file, auto-detecting format
func Load(filename string) ([]uint16, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Try to detect format based on extension or content
	if strings.HasSuffix(filename, ".hex") {
		return LoadHex(file)
	}

	return LoadHex(file) // Default to hex format
}

// LoadHex loads a hex file with optional comments
// Format supports:
//   - Simple: space-separated hex values
//   - Commented: hex values followed by # and assembly/comments
//   - Blank lines and full-line comments are ignored
//
// Example:
//   5400 0000  # r1 <- copy val/0
//   544A 0001  # r1 <- add val/1
func LoadHex(r io.Reader) ([]uint16, error) {
	var program []uint16
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Trim whitespace
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip full-line comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Split on # to separate hex from comments
		parts := strings.SplitN(line, "#", 2)
		hexPart := strings.TrimSpace(parts[0])

		// Skip if only comment
		if hexPart == "" {
			continue
		}

		// Parse hex values (space-separated)
		hexValues := strings.Fields(hexPart)
		for _, hexStr := range hexValues {
			// Parse as 16-bit hex value
			value, err := strconv.ParseUint(hexStr, 16, 16)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid hex value '%s': %w", lineNum, hexStr, err)
			}
			program = append(program, uint16(value))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(program) == 0 {
		return nil, fmt.Errorf("no program data found")
	}

	return program, nil
}

// LoadBinary loads a raw binary file (16-bit words, little-endian)
func LoadBinary(r io.Reader) ([]uint16, error) {
	var program []uint16
	buf := make([]byte, 2)

	for {
		n, err := io.ReadFull(r, buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if n != 2 {
			return nil, fmt.Errorf("incomplete word at end of file")
		}

		// Little-endian 16-bit word
		word := uint16(buf[0]) | (uint16(buf[1]) << 8)
		program = append(program, word)
	}

	return program, nil
}
