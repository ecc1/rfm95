package rfm95

import (
	"fmt"
	"log"
)

const (
	bitrate   = 16384  // baud
	channelBW = 100000 // Hz
)

// ReadConfiguration reads the current register configuration from the radio,
// using either burst-mode or individual SPI reads.
func (r *Radio) ReadConfiguration(useBurst bool) []byte {
	if r.Error() != nil {
		return nil
	}
	n := len(resetConfiguration)
	config := make([]byte, n)
	start := config[ConfigurationStart:]
	if useBurst {
		copy(start, r.hw.ReadBurst(ConfigurationStart, n-ConfigurationStart))
		return config
	}
	for i := range start {
		start[i] = r.hw.ReadRegister(uint8(ConfigurationStart + i))
	}
	return config
}

// WriteConfiguration writes the given register configuration to the radio,
// using either burst-mode or individual SPI writes.
func (r *Radio) WriteConfiguration(config []byte, useBurst bool) {
	n := len(resetConfiguration)
	if len(config) != n {
		log.Panicf("WriteConfiguration: config length = %d, expected %d", len(config), n)
		return
	}
	start := config[ConfigurationStart:]
	if useBurst {
		r.hw.WriteBurst(ConfigurationStart, start)
		return
	}
	for i, v := range start {
		r.hw.WriteRegister(uint8(ConfigurationStart+i), v)
	}
}

// InitRF initializes the radio to communicate with
// a Medtronic insulin pump at the given frequency.
func (r *Radio) InitRF(frequency uint32) {
	// Must be in Sleep mode first before changing to FSK/OOK mode.
	r.setMode(SleepMode)
	rf := DefaultConfiguration()
	rf[RegOpMode] = FskOokMode | ModulationTypeOOK | SleepMode
	// Interrupt on DIO2 when Sync word is seen.
	rf[RegDioMapping1] = 3 << Dio2MappingShift
	// Use 2^(5+1) = 64 samples for RSSI.
	rf[RegRssiConfig] = 5
	// Make sure enough preamble bytes are sent.
	rf[RegPreambleMsb] = 0x00
	rf[RegPreambleLsb] = 0x18
	// Use 4 bytes for Sync word.
	rf[RegSyncConfig] = SyncOn | 3<<SyncSizeShift
	// Sync word.
	rf[RegSyncValue1] = 0xFF
	rf[RegSyncValue2] = 0x00
	rf[RegSyncValue3] = 0xFF
	rf[RegSyncValue4] = 0x00
	// Use unlimited length packet format (data sheet section 4.2.13.2).
	rf[RegPacketConfig1] = FixedLength
	rf[RegPayloadLength] = 0
	rf[RegPacketConfig2] = PacketMode | 0
	// Must use PA_BOOST pin on Adafruit RFM95 bonnet. This sets Pout = 3 dBm.
	rf[RegPaConfig] = PaBoost | 1<<OutputPowerShift
	rf[RegPaRamp] = ModulationShapingNarrow | PaRamp100μs
	r.WriteConfiguration(rf, true)
	r.SetFrequency(frequency)
	r.SetBitrate(bitrate)
	r.SetChannelBW(channelBW)
	// RegPaDac is not in the DefaultConfiguration range, so set it individually.
	r.hw.WriteRegister(RegPaDac, PaDacDefault)
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
	r.hw.WriteRegister(RegOpMode, cur&^ModeMask|mode)
	if debug {
		log.Printf("change from %s to %s", stateName(cur&ModeMask), stateName(mode))
	}
	for r.Error() == nil {
		s := r.mode()
		if s == mode {
			break
		}
		if debug {
			log.Printf("  %s", stateName(s))
		}
	}
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
