package rfm95

import (
	"bytes"
	"time"

	"github.com/ecc1/gpio"
	"github.com/ecc1/radio"
)

const (
	hwVersion = 0x0102
)

type hwFlavor struct{}

// SPIDevice returns the pathname of the radio's SPI device.
func (hwFlavor) SPIDevice() string {
	return spiDevice
}

// Speed returns the radio's SPI speed.
func (hwFlavor) Speed() int {
	return spiSpeed
}

// CustomCS returns the GPIO pin number to use as a custom chip-select for the radio.
func (hwFlavor) CustomCS() int {
	return customCS
}

// InterruptPin returns the GPIO pin number to use for receive interrupts.
func (hwFlavor) InterruptPin() int {
	return interruptPin
}

// ReadSingleAddress returns the (identity) encoding of an address for SPI read operations.
func (hwFlavor) ReadSingleAddress(addr byte) byte {
	return addr
}

// ReadBurstAddress returns the (identity) encoding of an address for SPI burst-read operations.
func (hwFlavor) ReadBurstAddress(addr byte) byte {
	return addr
}

// WriteSingleAddress returns the encoding of an address for SPI write operations.
func (hwFlavor) WriteSingleAddress(addr byte) byte {
	return SPIWriteMode | addr
}

// WriteBurstAddress returns the encoding of an address for SPI burst-write operations.
func (hwFlavor) WriteBurstAddress(addr byte) byte {
	return SPIWriteMode | addr
}

// Radio represents an open radio device.
type Radio struct {
	hw            *radio.Hardware
	resetPin      gpio.OutputPin
	receiveBuffer bytes.Buffer
	err           error
}

// Open opens the radio device.
func Open() *Radio {
	r := &Radio{hw: radio.Open(hwFlavor{})}
	v := r.Version()
	if r.Error() != nil {
		return r
	}
	if v != hwVersion {
		r.hw.Close()
		r.SetError(radio.HardwareVersionError{Actual: v, Expected: hwVersion})
		return r
	}
	r.resetPin, r.err = gpio.Output(resetPin, false, false)
	if r.Error() != nil {
		r.hw.Close()
	}
	return r
}

// Close closes the radio device.
func (r *Radio) Close() {
	r.setMode(SleepMode)
	r.hw.Close()
}

// Name returns the radio's name.
func (r *Radio) Name() string {
	return "RFM95W"
}

// Device returns the pathname of the radio's device.
func (*Radio) Device() string {
	return spiDevice
}

// Version returns the radio's hardware version.
func (r *Radio) Version() uint16 {
	v := r.hw.ReadRegister(RegVersion)
	return uint16(v>>4)<<8 | uint16(v&0xF)
}

// Reset resets the radio device.  See section 7.2.2 of data sheet.
func (r *Radio) Reset() {
	if r.Error() != nil {
		return
	}
	r.sequencerStop()
	r.setMode(SleepMode)
	_ = r.resetPin.Write(true)
	time.Sleep(100 * time.Microsecond)
	r.err = r.resetPin.Write(false)
	time.Sleep(5 * time.Millisecond)
}

// Init initializes the radio device.
func (r *Radio) Init(frequency uint32) {
	r.Reset()
	r.InitRF(frequency)
	r.setMode(SleepMode)
}

// Error returns the error state of the radio device.
func (r *Radio) Error() error {
	err := r.hw.Error()
	if err != nil {
		return err
	}
	return r.err
}

// SetError sets the error state of the radio device.
func (r *Radio) SetError(err error) {
	r.hw.SetError(err)
	r.err = err
}

// Hardware returns the radio's hardware information.
func (r *Radio) Hardware() *radio.Hardware {
	return r.hw
}
