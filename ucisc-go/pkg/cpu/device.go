package cpu

// Device represents a hardware device attached to the CPU
type Device interface {
	// ID returns the unique device ID (0-65535)
	ID() uint16

	// Type returns the device type code (0-255)
	Type() uint8

	// Read reads from device control memory (offset 0-15)
	Read(offset uint16) uint16

	// Write writes to device control memory (offset 0-15)
	Write(offset uint16, value uint16)

	// ReadBlock reads from device block memory (offset 0-255)
	// Returns 0 if device doesn't support block memory
	ReadBlock(offset uint16) uint16

	// WriteBlock writes to device block memory (offset 0-255)
	// Ignored if device doesn't support block memory
	WriteBlock(offset uint16, value uint16)
}

// DeviceManager manages devices attached to the CPU
type DeviceManager struct {
	devices [256]Device // Up to 256 devices (indexed by device number)
}

// NewDeviceManager creates a new device manager
func NewDeviceManager() *DeviceManager {
	return &DeviceManager{}
}

// AttachDevice attaches a device at the specified device number (0-255)
func (dm *DeviceManager) AttachDevice(deviceNum uint8, device Device) {
	dm.devices[deviceNum] = device
}

// Read reads from device memory
// addr is the full device memory address (0x0000-0xFFFF)
func (dm *DeviceManager) Read(addr uint16) uint16 {
	if addr < 0x1000 {
		// Device control space (0x0000-0x0FFF)
		// Each device has 16 words: device N at N*16
		deviceNum := uint8(addr >> 4)
		offset := addr & 0x0F

		device := dm.devices[deviceNum]
		if device == nil {
			return 0 // Unmapped device returns 0
		}

		return device.Read(offset)
	} else {
		// Device block memory (0x1000-0xFFFF)
		// Device 16+ only: device N at 0x1000 + (N-16)*256
		deviceNum := uint8(((addr - 0x1000) >> 8) + 16)
		offset := uint8(addr & 0xFF)

		if deviceNum < 16 {
			return 0 // First 16 devices have no block memory
		}

		device := dm.devices[deviceNum]
		if device == nil {
			return 0 // Unmapped device returns 0
		}

		return device.ReadBlock(uint16(offset))
	}
}

// Write writes to device memory
// addr is the full device memory address (0x0000-0xFFFF)
func (dm *DeviceManager) Write(addr uint16, value uint16) {
	if addr < 0x1000 {
		// Device control space (0x0000-0x0FFF)
		deviceNum := uint8(addr >> 4)
		offset := addr & 0x0F

		device := dm.devices[deviceNum]
		if device == nil {
			return // Unmapped device, ignore write
		}

		device.Write(offset, value)
	} else {
		// Device block memory (0x1000-0xFFFF)
		deviceNum := uint8(((addr - 0x1000) >> 8) + 16)
		offset := uint8(addr & 0xFF)

		if deviceNum < 16 {
			return // First 16 devices have no block memory
		}

		device := dm.devices[deviceNum]
		if device == nil {
			return // Unmapped device, ignore write
		}

		device.WriteBlock(uint16(offset), value)
	}
}

// GetDevice returns the device at the specified device number
func (dm *DeviceManager) GetDevice(deviceNum uint8) Device {
	return dm.devices[deviceNum]
}
