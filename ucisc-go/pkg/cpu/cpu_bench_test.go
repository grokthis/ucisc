package cpu

import (
	"testing"
	"time"
)

// MinRequiredMIPS is the minimum performance requirement in millions of instructions per second
// Note: This measures wall-clock time including Reset() overhead, so it's lower than
// the benchmark framework's reported MIPS (50+ MIPS) which excludes overhead
const MinRequiredMIPS = 15.0

// reportIPS calculates and reports instructions per second
// totalInstructions is the total number of instructions executed across all iterations
// elapsedTime is the actual time spent executing (not including benchmark overhead)
func reportIPS(b *testing.B, totalInstructions uint64, elapsedTime time.Duration) {
	ips := float64(totalInstructions) / elapsedTime.Seconds()
	mips := ips / 1e6

	b.ReportMetric(ips, "IPS")
	b.ReportMetric(mips, "MIPS")

	// Verify we meet minimum performance requirements
	if mips < MinRequiredMIPS {
		b.Errorf("Performance requirement not met: %.2f MIPS < %.2f MIPS required", mips, MinRequiredMIPS)
	}
}

// BenchmarkCPU_SimpleLoop benchmarks a simple counting loop
func BenchmarkCPU_SimpleLoop(b *testing.B) {
	// Loop that counts from 0 to 100
	// &r1 <- copy val/0
	// loop:
	//   &r1 <- add val/1
	//   &r1 <~ sub val/100 (sets zero flag when r1==100)
	//   pc <1 copy pc/loop  (jump if NOT zero)
	// pc <- copy pc/0  (halt)

	code := []uint16{
		0x4540, 0x0000, // &r1 <- copy val/0
		0x454A, 0x0001, // loop: &r1 <- add val/1
		0x457B, 0x0064, // &r1 <~ sub val/100 (test against 100)
		0x4010, 0x0002, // pc <1 copy val/2 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0 (halt)
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}

// BenchmarkCPU_MemoryAccess benchmarks memory read/write operations
func BenchmarkCPU_MemoryAccess(b *testing.B) {
	// Loop that reads and writes memory
	// &r1 <- copy val/100   (base address)
	// &r2 <- copy val/0     (counter)
	// loop:
	//   r1/0 <- copy val/42  (write to memory)
	//   &r3 <- copy r1/0     (read from memory)
	//   &r2 <- add val/1
	//   &r2 <~ sub val/50
	//   pc <1 copy pc/loop
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x0064, // &r1 <- copy val/100
		0x4640, 0x0000, // &r2 <- copy val/0
		0x4140, 0x002A, // loop: r1/0 <- copy val/42
		0x1740, 0x0000, // &r3 <- copy r1/0
		0x464A, 0x0001, // &r2 <- add val/1
		0x467B, 0x0032, // &r2 <~ sub val/50
		0x4010, 0x0004, // pc <1 copy val/4 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}

// BenchmarkCPU_ALUOperations benchmarks various ALU operations
func BenchmarkCPU_ALUOperations(b *testing.B) {
	// Loop performing multiple ALU operations
	// &r1 <- copy val/10
	// &r2 <- copy val/20
	// &r3 <- copy val/0
	// &r5 <- copy val/0 (counter)
	// loop:
	//   &r3 <- add r1
	//   &r3 <- sub r2
	//   &r3 <- and val/0xFF
	//   &r3 <- or val/0x10
	//   &r3 <- xor val/0xAA
	//   &r5 <- add val/1
	//   &r5 <~ sub val/100
	//   pc <1 copy val/loop
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x000A, // &r1 <- copy val/10
		0x4640, 0x0014, // &r2 <- copy val/20
		0x4740, 0x0000, // &r3 <- copy val/0
		0x4D40, 0x0000, // &r5 <- copy val/0
		0x174A, 0x0000, // loop: &r3 <- add r1
		0x274B, 0x0000, // &r3 <- sub r2
		0x4741, 0x00FF, // &r3 <- and val/0xFF
		0x4742, 0x0010, // &r3 <- or val/0x10
		0x4743, 0x00AA, // &r3 <- xor val/0xAA
		0x4D4A, 0x0001, // &r5 <- add val/1
		0x4D7B, 0x0064, // &r5 <~ sub val/100
		0x4010, 0x0008, // pc <1 copy val/8 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}

// BenchmarkCPU_Multiplication benchmarks multiplication operations
func BenchmarkCPU_Multiplication(b *testing.B) {
	// Loop performing multiplications
	// &r1 <- copy val/5
	// &r2 <- copy val/7
	// &r3 <- copy val/0
	// loop:
	//   &r3 <- mult r1  (r3 = r1 * r2)
	//   &r3 <- add val/1
	//   &r3 <~ sub val/100
	//   pc <1 copy pc/loop
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x0005, // &r1 <- copy val/5
		0x4640, 0x0007, // &r2 <- copy val/7
		0x4740, 0x0000, // &r3 <- copy val/0
		0x4768, 0x0000, // loop: &r3 <- mult r1 (using r2 as implicit source)
		0x474A, 0x0001, // &r3 <- add val/1
		0x477B, 0x0064, // &r3 <~ sub val/100
		0x4010, 0x0006, // pc <1 copy val/6 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}

// BenchmarkCPU_ShiftOperations benchmarks shift operations
func BenchmarkCPU_ShiftOperations(b *testing.B) {
	// Loop performing shift operations
	// &r1 <- copy val/1
	// &r2 <- copy val/0
	// loop:
	//   &r1 <- shl val/1
	//   &r1 <- shr val/1
	//   &r2 <- add val/1
	//   &r2 <~ sub val/100
	//   pc <1 copy pc/loop
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x0001, // &r1 <- copy val/1
		0x4640, 0x0000, // &r2 <- copy val/0
		0x4545, 0x0001, // loop: &r1 <- shl val/1
		0x4546, 0x0001, // &r1 <- shr val/1
		0x464A, 0x0001, // &r2 <- add val/1
		0x467B, 0x0064, // &r2 <~ sub val/100
		0x4010, 0x0004, // pc <1 copy val/4 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}

// BenchmarkCPU_ConditionalBranching benchmarks conditional execution
func BenchmarkCPU_ConditionalBranching(b *testing.B) {
	// Loop with multiple conditional branches
	// &r1 <- copy val/0
	// loop:
	//   &r1 <- add val/1
	//   &r2 <0 copy val/100  (should not execute)
	//   &r2 <1 copy val/200  (should execute)
	//   &r1 <~ sub val/50
	//   pc <1 copy pc/loop
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x0000, // &r1 <- copy val/0
		0x454A, 0x0001, // loop: &r1 <- add val/1
		0x4600, 0x0064, // &r2 <0 copy val/100
		0x4610, 0x00C8, // &r2 <1 copy val/200
		0x457B, 0x0032, // &r1 <~ sub val/50
		0x4010, 0x0002, // pc <1 copy val/2 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}

// BenchmarkCPU_PushPop benchmarks stack operations
// TODO: Fix when assembler supports push/pop syntax
/*
func BenchmarkCPU_PushPop(b *testing.B) {
	// Loop with push and pop operations
	// &r1 <- copy val/500  (stack pointer)
	// &r2 <- copy val/0    (counter)
	// loop:
	//   r1> <- copy val/42  (push)
	//   &r3 <r1 copy val/0  (pop)
	//   &r2 <- add val/1
	//   &r2 <~ sub val/50
	//   pc <1 copy pc/loop
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x01F4, // &r1 <- copy val/500
		0x4640, 0x0000, // &r2 <- copy val/0
		0x1240, 0x002A, // loop: r1> <- copy val/42
		0x4B50, 0x0000, // &r3 <r1 copy val/0
		0x464A, 0x0001, // &r2 <- add val/1
		0x467B, 0x0032, // &r2 <~ sub val/50
		0x4010, 0x0004, // pc <1 copy val/4 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}
*/

// BenchmarkCPU_RegisterToRegister benchmarks register-to-register operations
func BenchmarkCPU_RegisterToRegister(b *testing.B) {
	// Loop copying between registers
	// &r1 <- copy val/42
	// &r5 <- copy val/0
	// loop:
	//   &r2 <- copy r1
	//   &r3 <- copy r2
	//   &r5 <- copy r3
	//   &r6 <- copy r5
	//   &r7 <- copy r6
	//   &r5 <- add val/1
	//   &r5 <~ sub val/100
	//   pc <1 copy pc/loop
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x002A, // &r1 <- copy val/42
		0x4D40, 0x0000, // &r5 <- copy val/0
		0x4E40, 0x0000, // &r6 <- copy val/0
		0x1640, 0x0000, // loop: &r2 <- copy r1
		0x2740, 0x0000, // &r3 <- copy r2
		0x3D40, 0x0000, // &r5 <- copy r3
		0x9F40, 0x0000, // &r7 <- copy r5
		0x4E4A, 0x0001, // &r6 <- add val/1
		0x4E7B, 0x0064, // &r6 <~ sub val/100
		0x4010, 0x0006, // pc <1 copy val/6 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}

// BenchmarkCPU_NestedLoops benchmarks nested loop performance
func BenchmarkCPU_NestedLoops(b *testing.B) {
	// Nested loops (outer loop 10 times, inner loop 10 times)
	// &r1 <- copy val/0    (outer counter)
	// outer:
	//   &r2 <- copy val/0  (inner counter)
	//   inner:
	//     &r2 <- add val/1
	//     &r2 <~ sub val/10
	//     pc <1 copy pc/inner
	//   &r1 <- add val/1
	//   &r1 <~ sub val/10
	//   pc <1 copy pc/outer
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x0000, // &r1 <- copy val/0
		0x4640, 0x0000, // outer: &r2 <- copy val/0
		0x464A, 0x0001, // inner: &r2 <- add val/1
		0x467B, 0x000A, // &r2 <~ sub val/10
		0x4010, 0x0004, // pc <1 copy val/4 (jump to inner if not zero)
		0x454A, 0x0001, // &r1 <- add val/1
		0x457B, 0x000A, // &r1 <~ sub val/10
		0x4010, 0x0002, // pc <1 copy val/2 (jump to outer if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}

// BenchmarkCPU_MemoryCopy benchmarks copying data in memory
func BenchmarkCPU_MemoryCopy(b *testing.B) {
	// Copy 50 words from address 100 to address 200
	// &r1 <- copy val/100  (source)
	// &r2 <- copy val/200  (dest)
	// &r3 <- copy val/0    (counter)
	// loop:
	//   &r5 <- copy r1/0   (read)
	//   r2/0 <- copy r5    (write)
	//   &r1 <- add val/1
	//   &r2 <- add val/1
	//   &r3 <- add val/1
	//   &r3 <~ sub val/50
	//   pc <1 copy pc/loop
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x0064, // &r1 <- copy val/100
		0x4640, 0x00C8, // &r2 <- copy val/200
		0x4740, 0x0000, // &r3 <- copy val/0
		0x1D40, 0x0000, // loop: &r5 <- copy r1/0
		0x9240, 0x0000, // r2/0 <- copy r5
		0x454A, 0x0001, // &r1 <- add val/1
		0x464A, 0x0001, // &r2 <- add val/1
		0x474A, 0x0001, // &r3 <- add val/1
		0x477B, 0x0032, // &r3 <~ sub val/50
		0x4010, 0x0006, // pc <1 copy val/6 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Pre-populate source memory and load program once
	c := New(1024)
	c.LoadProgram(code, 0)
	for i := 0; i < 50; i++ {
		c.memory.WriteLocal(uint16(100+i), uint16(i*10))
	}

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter

		// Reset CPU state but need to restore source memory
		c.Reset()
		for j := 0; j < 50; j++ {
			c.memory.WriteLocal(uint16(100+j), uint16(j*10))
		}
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}

// BenchmarkCPU_Fibonacci benchmarks calculating Fibonacci numbers
func BenchmarkCPU_Fibonacci(b *testing.B) {
	// Calculate first 20 Fibonacci numbers
	// &r1 <- copy val/0    (fib(n-2))
	// &r2 <- copy val/1    (fib(n-1))
	// &r3 <- copy val/0    (counter)
	// loop:
	//   &r5 <- add r1      (r5 = r1 + r2)
	//   &r1 <- copy r2     (shift: r1 = r2)
	//   &r2 <- copy r5     (shift: r2 = r5)
	//   &r3 <- add val/1
	//   &r3 <~ sub val/20
	//   pc <1 copy pc/loop
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x0000, // &r1 <- copy val/0
		0x4640, 0x0001, // &r2 <- copy val/1
		0x4740, 0x0000, // &r3 <- copy val/0
		0x1D4A, 0x0000, // loop: &r5 <- add r1 (using r2)
		0x2540, 0x0000, // &r1 <- copy r2
		0x9640, 0x0000, // &r2 <- copy r5
		0x474A, 0x0001, // &r3 <- add val/1
		0x477B, 0x0014, // &r3 <~ sub val/20
		0x4010, 0x0006, // pc <1 copy val/6 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}

// BenchmarkCPU_BitwiseOperations benchmarks bitwise operations
func BenchmarkCPU_BitwiseOperations(b *testing.B) {
	// Loop performing various bitwise operations
	// &r1 <- copy val/0xAAAA
	// &r2 <- copy val/0x5555
	// &r3 <- copy val/0
	// loop:
	//   &r3 <- and r1
	//   &r3 <- or r2
	//   &r3 <- xor val/0xFFFF
	//   &r3 <- inv val/0
	//   &r3 <- swap val/0
	//   &r3 <- add val/1
	//   &r3 <~ sub val/100
	//   pc <1 copy pc/loop
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x2AAA, // &r1 <- copy val/0xAAAA (positive signed)
		0x4640, 0x5555, // &r2 <- copy val/0x5555
		0x4740, 0x0000, // &r3 <- copy val/0
		0x4D40, 0x0000, // &r5 <- copy val/0 (counter)
		0x1741, 0x0000, // loop: &r3 <- and r1 (using r2)
		0x2742, 0x0000, // &r3 <- or r2
		0x4743, 0x7FFF, // &r3 <- xor val/0x7FFF
		0x4744, 0x0000, // &r3 <- inv val/0
		0x4747, 0x0000, // &r3 <- swap val/0
		0x4D4A, 0x0001, // &r5 <- add val/1
		0x4D7B, 0x0064, // &r5 <~ sub val/100
		0x4010, 0x0008, // pc <1 copy val/8 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}

// BenchmarkCPU_InstructionFetch benchmarks raw instruction fetch and decode
func BenchmarkCPU_InstructionFetch(b *testing.B) {
	// Minimal loop to isolate fetch/decode overhead
	// &r1 <- copy val/0
	// loop:
	//   &r1 <- add val/1
	//   &r1 <~ sub val/1000
	//   pc <1 copy pc/loop
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x0000, // &r1 <- copy val/0
		0x454A, 0x0001, // loop: &r1 <- add val/1
		0x457B, 0x03E8, // &r1 <~ sub val/1000
		0x4010, 0x0002, // pc <1 copy val/2 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}

// BenchmarkCPU_AllRegisters benchmarks operations across all registers
func BenchmarkCPU_AllRegisters(b *testing.B) {
	// Loop using all 6 registers
	// &r1 <- copy val/1
	// &r2 <- copy val/2
	// &r3 <- copy val/3
	// &r5 <- copy val/5
	// &r6 <- copy val/6
	// &r7 <- copy val/0
	// loop:
	//   &r7 <- add r1
	//   &r7 <- add r2
	//   &r7 <- add r3
	//   &r7 <- add r5
	//   &r7 <- add r6
	//   &r7 <~ sub val/1000
	//   pc <1 copy pc/loop
	// pc <- copy pc/0

	code := []uint16{
		0x4540, 0x0001, // &r1 <- copy val/1
		0x4640, 0x0002, // &r2 <- copy val/2
		0x4740, 0x0003, // &r3 <- copy val/3
		0x4D40, 0x0005, // &r5 <- copy val/5
		0x4E40, 0x0006, // &r6 <- copy val/6
		0x4F40, 0x0000, // &r7 <- copy val/0
		0x1F4A, 0x0000, // loop: &r7 <- add r1
		0x2F4A, 0x0000, // &r7 <- add r2
		0x3F4A, 0x0000, // &r7 <- add r3
		0x9F4A, 0x0000, // &r7 <- add r5
		0xAF4A, 0x0000, // &r7 <- add r6
		0x4F7B, 0x03E8, // &r7 <~ sub val/1000
		0x4010, 0x000C, // pc <1 copy val/12 (jump back to loop if not zero)
		0x0040, 0x0000, // pc <- copy pc/0
	}

	// Load program once outside the loop
	c := New(1024)
	c.LoadProgram(code, 0)

	b.ResetTimer()
	start := time.Now()
	var totalInstructions uint64
	var expectedInstructions uint64
	for i := 0; i < b.N; i++ {
		startCount := c.GetInstCount()
		if startCount != 0 {
			b.Fatalf("Iteration %d: instruction count not reset (got %d, expected 0)", i, startCount)
		}

		c.Run()
		instThisIter := c.GetInstCount()

		if i == 0 {
			expectedInstructions = instThisIter
		} else if instThisIter != expectedInstructions {
			b.Fatalf("Iteration %d: instruction count mismatch (got %d, expected %d)", i, instThisIter, expectedInstructions)
		}

		totalInstructions += instThisIter
		c.Reset()
	}
	elapsed := time.Since(start)
	reportIPS(b, totalInstructions, elapsed)
}
