package cpu

// Memory represents the processor's memory system
type Memory struct {
	local   []uint16       // Local processor memory (up to 64K words)
	size    int            // Size of local memory
	devices *DeviceManager // Device manager for banked memory access
}

// NewMemory creates a new memory system
func NewMemory(size int) *Memory {
	if size > 65536 {
		size = 65536
	}
	if size < 0 {
		size = 0
	}

	return &Memory{
		local:   make([]uint16, size),
		size:    size,
		devices: NewDeviceManager(),
	}
}

// ReadLocal reads from local memory (no banking)
func (m *Memory) ReadLocal(addr uint16) uint16 {
	if int(addr) >= m.size {
		return 0 // Out of bounds returns 0
	}
	return m.local[addr]
}

// WriteLocal writes to local memory (no banking)
func (m *Memory) WriteLocal(addr uint16, value uint16) {
	if int(addr) < m.size {
		m.local[addr] = value
	}
}

// Read reads from memory, handling banking
func (m *Memory) Read(addr uint16, reg Register, banking uint16) uint16 {
	// Check if register is banked
	if m.isBanked(reg, banking) {
		// Route to device memory
		return m.devices.Read(addr)
	}

	return m.ReadLocal(addr)
}

// Write writes to memory, handling banking
func (m *Memory) Write(addr uint16, reg Register, banking uint16, value uint16) {
	// Check if register is banked
	if m.isBanked(reg, banking) {
		// Route to device memory
		m.devices.Write(addr, value)
		return
	}

	m.WriteLocal(addr, value)
}

// isBanked checks if a register is banked
func (m *Memory) isBanked(reg Register, banking uint16) bool {
	// Banking bits: bit N corresponds to register N
	// R1 = bit 1, R2 = bit 2, R3 = bit 3, R5 = bit 5, R6 = bit 6
	switch reg {
	case RegR1:
		return banking&(1<<1) != 0
	case RegR2:
		return banking&(1<<2) != 0
	case RegR3:
		return banking&(1<<3) != 0
	case RegR5:
		return banking&(1<<5) != 0
	case RegR6:
		return banking&(1<<6) != 0
	default:
		return false
	}
}

// LoadProgram loads a program into memory at the specified offset
func (m *Memory) LoadProgram(data []uint16, offset uint16) error {
	for i, word := range data {
		addr := offset + uint16(i)
		if int(addr) >= m.size {
			break
		}
		m.local[addr] = word
	}
	return nil
}

// Size returns the size of local memory
func (m *Memory) Size() int {
	return m.size
}

// Dump returns a copy of memory contents for debugging
func (m *Memory) Dump(start, end uint16) []uint16 {
	if int(start) >= m.size {
		return nil
	}
	if int(end) > m.size {
		end = uint16(m.size)
	}
	if end <= start {
		return nil
	}

	result := make([]uint16, end-start)
	copy(result, m.local[start:end])
	return result
}

// AttachDevice attaches a device to the memory system
func (m *Memory) AttachDevice(deviceNum uint8, device Device) {
	m.devices.AttachDevice(deviceNum, device)
}

// GetDeviceManager returns the device manager
func (m *Memory) GetDeviceManager() *DeviceManager {
	return m.devices
}
