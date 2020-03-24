package rfm95

import (
	"fmt"
	"log"
	"unsafe"
)

const (
	bitrate   = 16384  // baud
	channelBW = 100000 // Hz
)

// Bytes returns the RFConfiguration as a byte slice.
func (config *RFConfiguration) Bytes() []byte {
	return (*[RegVersion - RegOpMode + 1]byte)(unsafe.Pointer(config))[:]
}

// ReadConfiguration reads the current RFConfiguration from the radio.
func (r *Radio) ReadConfiguration() *RFConfiguration {
	if r.Error() != nil {
		return nil
	}
	regs := r.hw.ReadBurst(RegOpMode, RegVersion-RegOpMode+1)
	return (*RFConfiguration)(unsafe.Pointer(&regs[0]))
}

// WriteConfiguration writes the given RFConfiguration to the radio.
func (r *Radio) WriteConfiguration(config *RFConfiguration) {
	r.hw.WriteBurst(RegOpMode, config.Bytes())
}

// InitRF initializes the radio to communicate with
// a Medtronic insulin pump at the given frequency.
func (r *Radio) InitRF(frequency uint32) {
	rf := DefaultRFConfiguration

	rf.RegOpMode = FskOokMode | ModulationTypeOOK | SleepMode

	// Use 2^(5+1) = 64 samples for RSSI.
	rf.RegRssiConfig = 5

	// Make sure enough preamble bytes are sent.
	rf.RegPreambleMsb = 0x00
	rf.RegPreambleLsb = 0x18

	// Use 4 bytes for Sync word.
	rf.RegSyncConfig = SyncOn | 3<<SyncSizeShift

	// Sync word.
	rf.RegSyncValue1 = 0xFF
	rf.RegSyncValue2 = 0x00
	rf.RegSyncValue3 = 0xFF
	rf.RegSyncValue4 = 0x00

	// Must be in Sleep mode first before changing to FSK/OOK mode.
	r.hw.WriteRegister(RegOpMode, FskOokMode|ModulationTypeOOK|SleepMode)

	r.WriteConfiguration(&rf)
	r.SetFrequency(frequency)
	r.SetBitrate(bitrate)
	r.SetChannelBW(channelBW)
}

// Frequency returns the radio's current frequency, in Hertz.
func (r *Radio) Frequency() uint32 {
	return registersToFrequency(r.hw.ReadBurst(RegFrfMsb, 3))
}

func registersToFrequency(frf []byte) uint32 {
	f := uint32(frf[0])<<16 + uint32(frf[1])<<8 + uint32(frf[2])
	return uint32(uint64(f) * FXOSC >> 19)
}

// SetFrequency sets the radio to the given frequency, in Hertz.
func (r *Radio) SetFrequency(freq uint32) {
	r.hw.WriteBurst(RegFrfMsb, frequencyToRegisters(freq))
}

func frequencyToRegisters(freq uint32) []byte {
	f := (uint64(freq)<<19 + FXOSC/2) / FXOSC
	return []byte{byte(f >> 16), byte(f >> 8), byte(f)}
}

// ReadRSSI returns the radio's RSSI, in dBm.
func (r *Radio) ReadRSSI() int {
	rssi := r.hw.ReadRegister(RegRssiValue)
	return -int(rssi) / 2
}

// Bitrate returns the radio's bit rate, in bps.
func (r *Radio) Bitrate() uint32 {
	return registersToBitrate(r.hw.ReadBurst(RegBitrateMsb, 2))
}

// See data sheet section 4.2.1.
func registersToBitrate(br []byte) uint32 {
	d := uint32(br[0])<<8 + uint32(br[1])
	return (FXOSC + d/2) / d
}

// SetBitrate sets the radio's bit rate to the given rate, in bps.
func (r *Radio) SetBitrate(br uint32) {
	r.hw.WriteBurst(RegBitrateMsb, bitrateToRegisters(br))
}

func bitrateToRegisters(br uint32) []byte {
	b := (FXOSC + br/2) / br
	return []byte{byte(b >> 8), byte(b)}
}

// ReadModulationType returns the radio's modulation type.
func (r *Radio) ReadModulationType() byte {
	return r.hw.ReadRegister(RegOpMode) & ModulationTypeMask
}

// ChannelBW returns the radio's channel bandwidth, in Hertz.
func (r *Radio) ChannelBW() uint32 {
	return registerToChannelBW(r.hw.ReadRegister(RegRxBw))
}

func registerToChannelBW(bw byte) uint32 {
	mant := 0
	switch bw & RxBwMantMask {
	case RxBwMant16:
		mant = 16
	case RxBwMant20:
		mant = 20
	case RxBwMant24:
		mant = 24
	default:
		log.Panicf("unknown RX bandwidth mantissa (%X)", bw&RxBwMantMask)
	}
	e := bw & RxBwExpMask
	return uint32(FXOSC) / (uint32(mant) << (e + 2))
}

// SetChannelBW sets the radio's channel bandwidth to the given value, in Hertz.
func (r *Radio) SetChannelBW(bw uint32) {
	r.hw.WriteRegister(RegRxBw, channelBWToRegister(bw))
}

// Channel BW = FXOSC / (RxBwMant * 2^(RxBwExp + 2)).
func channelBWToRegister(bw uint32) byte {
	bb := uint32(2604) // lowest possible channel bandwidth
	rr := byte(RxBwMant24 | 7<<RxBwExpShift)
	if bw < bb {
		return rr
	}
	// RxBwExp value of 0 is reserved.
	for i := 0; i < 7; i++ {
		e := byte(7 - i)
		for j := 0; j < 3; j++ {
			m := byte((6 - j) * 4)
			b := uint32(FXOSC) / (uint32(m) << (e + 2))
			r := byte(2-j)<<RxBwMantShift | e<<RxBwExpShift
			if b >= bw {
				if b-bw < bw-bb {
					return r
				}
				return rr
			}
			bb = b
			rr = r
		}
	}
	return rr
}

func (r *Radio) mode() byte {
	return r.hw.ReadRegister(RegOpMode) & ModeMask
}

func (r *Radio) setMode(mode uint8) {
	r.SetError(nil)
	cur := r.hw.ReadRegister(RegOpMode)
	if cur&ModeMask == mode {
		return
	}
	if verbose {
		log.Printf("change from %s to %s", stateName(cur&ModeMask), stateName(mode))
	}
	r.hw.WriteRegister(RegOpMode, cur&^ModeMask|mode)
	for r.Error() == nil {
		s := r.mode()
		if s == mode && r.modeReady() {
			break
		}
		if verbose {
			log.Printf("  %s", stateName(s))
		}
	}
}

func (r *Radio) modeReady() bool {
	return r.hw.ReadRegister(RegIrqFlags1)&ModeReady != 0
}

// Sleep puts the radio into sleep mode.
func (r *Radio) Sleep() {
	r.setMode(SleepMode)
}

func stateName(mode uint8) string {
	switch mode {
	case SleepMode:
		return "Sleep"
	case StandbyMode:
		return "Standby"
	case FreqSynthModeTX:
		return "TX Frequency Synthesizer"
	case TransmitterMode:
		return "Transmitter"
	case FreqSynthModeRX:
		return "RX Frequency Synthesizer"
	case ReceiverMode:
		return "Receiver"
	default:
		return fmt.Sprintf("Unknown Mode (%X)", mode)
	}
}

// State returns the radio's current state as a string.
func (r *Radio) State() string {
	return stateName(r.mode())
}
