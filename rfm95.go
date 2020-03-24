package rfm95

// https://www.hoperf.com/data/upload/portal/20190801/RFM96W-V2.0.pdf

const (
	// FXOSC is the radio's oscillator frequency in Hertz.
	FXOSC = 32000000

	// SPIWriteMode is used to encode register addresses for SPI writes.
	SPIWriteMode = 1 << 7
)

// FIFO
const (
	RegFifo = 0x00 // FIFO read/write access
)

// Registers for Common settings
const (
	RegOpMode     = 0x01 // Operating mode & LoRa / FSK selection
	RegBitrateMsb = 0x02 // Bit Rate setting, Most Significant Bits
	RegBitrateLsb = 0x03 // Bit Rate setting, Least Significant Bits
	RegFdevMsb    = 0x04 // Frequency Deviation setting, Most Significant Bits
	RegFdevLsb    = 0x05 // Frequency Deviation setting, Least Significant Bits
	RegFrfMsb     = 0x06 // RF Carrier Frequency, Most Significant Bits
	RegFrfMid     = 0x07 // RF Carrier Frequency, Intermediate Bits
	RegFrfLsb     = 0x08 // RF Carrier Frequency, Least Significant Bits
)

// Registers for the Transmitter
const (
	RegPaConfig = 0x09 // PA selection and Output Power control
	RegPaRamp   = 0x0A // Control of the PA ramp time, low phase noise PLL
	RegOcp      = 0x0B // Over Current Protection control
)

// Registers for the Receiver
const (
	RegLna            = 0x0C // LNA settings
	RegRxConfig       = 0x0D // AFC, AGC, ctrl
	RegRssiConfig     = 0x0E // RSSI
	RegRssiCollision  = 0x0F // RSSI Collision detector
	RegRssiThresh     = 0x10 // RSSI Threshold control
	RegRssiValue      = 0x11 // RSSI value in dBm
	RegRxBw           = 0x12 // Channel Filter BW Control
	RegAfcBw          = 0x13 // AFC Channel Filter BW
	RegOokPeak        = 0x14 // OOK demodulator
	RegOokFix         = 0x15 // Threshold of the OOK demodulator
	RegOokAvg         = 0x16 // Average of the OOK demodulator
	RegAfcFei         = 0x1A // AFC and FEI control
	RegAfcMsb         = 0x1B // Frequency correction value of the AFC, MSB
	RegAfcLsb         = 0x1C // Frequency correction value of the AFC, LSB
	RegFeiMsb         = 0x1D // Value of the calculated frequency error, MSB
	RegFeiLsb         = 0x1E // Value of the calculated frequency error, LSB
	RegPreambleDetect = 0x1F // Settings of the Preamble Detector
	RegRxTimeout1     = 0x20 // Timeout duration between Rx request and RSSI detection
	RegRxTimeout2     = 0x21 // Timeout duration between RSSI detection and PayloadReady
	RegRxTimeout3     = 0x22 // Timeout duration between RSSI detection and SyncAddress
	RegRxDelay        = 0x23 // Delay between Rx cycles
)

// RC Oscillator registers
const (
	RegOsc = 0x24 // RC Oscillators Settings, CLK-OUT frequency
)

// Packet Handling registers
const (
	RegPreambleMsb   = 0x25 // Preamble length, MSB
	RegPreambleLsb   = 0x26 // Preamble length, LSB
	RegSyncConfig    = 0x27 // Sync Word Recognition control
	RegSyncValue1    = 0x28 // Sync Word bytes 1 through 8
	RegSyncValue2    = 0x29
	RegSyncValue3    = 0x2A
	RegSyncValue4    = 0x2B
	RegSyncValue5    = 0x2C
	RegSyncValue6    = 0x2D
	RegSyncValue7    = 0x2E
	RegSyncValue8    = 0x2F
	RegPacketConfig1 = 0x30 // Packet mode settings
	RegPacketConfig2 = 0x31 // Packet mode settings
	RegPayloadLength = 0x32 // Payload length setting
	RegNodeAdrs      = 0x33 // Node address
	RegBroadcastAdrs = 0x34 // Broadcast address
	RegFifoThresh    = 0x35 // Fifo threshold, Tx start condition
)

// Sequencer registers
const (
	RegSeqConfig1 = 0x36 // Top level Sequencer settings
	RegSeqConfig2 = 0x37 // Top level Sequencer settings
	RegTimerResol = 0x38 // Timer 1 and 2 resolution control
	RegTimer1Coef = 0x39 // Timer 1 setting
	RegTimer2Coef = 0x3A // Timer 2 setting
)

// Service registers
const (
	RegImageCal = 0x3B // Image calibration engine control
	RegTemp     = 0x3C // Temperature Sensor value
	RegLowBat   = 0x3D // Low Battery Indicator Settings
)

// Status registers
const (
	RegIrqFlags1 = 0x3E // Status register: PLL Lock state, Timeout, RSSI
	RegIrqFlags2 = 0x3F // Status register: FIFO handling flags, Low Battery
)

// IO control registers
const (
	RegDioMapping1 = 0x40 // Mapping of pins DIO0 to DIO3
	RegDioMapping2 = 0x41 // Mapping of pins DIO4 and DIO5, ClkOut frequency
)

// Version register
const (
	RegVersion = 0x42 // Hope RF ID relating the silicon revision
)

// Additional registers
const (
	RegPllHop      = 0x44 // Control the fast frequency hopping mode
	RegTcxo        = 0x4B // TCXO or XTAL input setting
	RegPaDac       = 0x4D // Higher power settings of the PA
	RegFormerTemp  = 0x5B // Stored temperature during the former IQ Calibration
	RegBitRateFrac = 0x5D // Fractional part in the Bit Rate division ratio
	RegAgcRef      = 0x61 // Adjustment of the AGC thresholds
	RegAgcThresh1  = 0x62
	RegAgcThresh2  = 0x63
	RegAgcThresh3  = 0x64
	RegPll         = 0x70 // Control of the PLL bandwidth
)

// RFConfiguration represents (most of) the radio's configuration registers.
type RFConfiguration struct {
	// Omit RegFifo to avoid reading or writing it with this struct.
	RegOpMode         byte // Operating mode & LoRa / FSK selection
	RegBitrateMsb     byte // Bit Rate setting, Most Significant Bits
	RegBitrateLsb     byte // Bit Rate setting, Least Significant Bits
	RegFdevMsb        byte // Frequency Deviation setting, Most Significant Bits
	RegFdevLsb        byte // Frequency Deviation setting, Least Significant Bits
	RegFrfMsb         byte // RF Carrier Frequency, Most Significant Bits
	RegFrfMid         byte // RF Carrier Frequency, Intermediate Bits
	RegFrfLsb         byte // RF Carrier Frequency, Least Significant Bits
	RegPaConfig       byte // PA selection and Output Power control
	RegPaRamp         byte // Control of the PA ramp time, low phase noise PLL
	RegOcp            byte // Over Current Protection control
	RegLna            byte // LNA settings
	RegRxConfig       byte // AFC, AGC, ctrl
	RegRssiConfig     byte // RSSI-related settings
	RegRssiCollision  byte // RSSI Collision detector
	RegRssiThresh     byte // RSSI Threshold control
	RegRssiValue      byte // RSSI value in dBm
	RegRxBw           byte // Channel Filter BW Control
	RegAfcBw          byte // AFC Channel Filter BW
	RegOokPeak        byte // OOK demodulator selection and control in peak mode
	RegOokFix         byte // Fixed threshold control of the OOK demodulator
	RegOokAvg         byte // Average threshold control of the OOK demodulator
	reserved17        byte
	reserved18        byte
	reserved19        byte
	RegAfcFei         byte // AFC and FEI control
	RegAfcMsb         byte // Frequency correction value of the AFC (MSB)
	RegAfcLsb         byte // Frequency correction value of the AFC (LSB)
	RegFeiMsb         byte // MSB of the calculated frequency error
	RegFeiLsb         byte // LSB of the calculated frequency error
	RegPreambleDetect byte // Settings of the Preamble Detector
	RegRxTimeout1     byte // Timeout duration between Rx request and RSSI detection
	RegRxTimeout2     byte // Timeout duration between RSSI detection and PayloadReady
	RegRxTimeout3     byte // Timeout duration between RSSI detection and SyncAddress
	RegRxDelay        byte // Delay between Rx cycles
	RegOsc            byte // RF Oscillators Settings, CLK-OUT frequency
	RegPreambleMsb    byte // Preamble length, MSB
	RegPreambleLsb    byte // Preamble length, LSB
	RegSyncConfig     byte // Sync Word Recognition control
	RegSyncValue1     byte // Sync Word bytes, 1 through 8
	RegSyncValue2     byte
	RegSyncValue3     byte
	RegSyncValue4     byte
	RegSyncValue5     byte
	RegSyncValue6     byte
	RegSyncValue7     byte
	RegSyncValue8     byte
	RegPacketConfig1  byte // Packet mode settings
	RegPacketConfig2  byte // Packet mode settings
	RegPayloadLength  byte // Payload length setting
	RegNodeAdrs       byte // Node address
	RegBroadcastAdrs  byte // Broadcast address
	RegFifoThresh     byte // Fifo threshold, Tx start condition
	RegSeqConfig1     byte // Top level Sequencer settings
	RegSeqConfig2     byte // Top level Sequencer settings
	RegTimerResol     byte // Timer 1 and 2 resolution control
	RegTimer1Coef     byte // Timer 1 setting
	RegTimer2Coef     byte // Timer 2 setting
	RegImageCal       byte // Image calibration engine control
	RegTemp           byte // Temperature Sensor value
	RegLowBat         byte // Low Battery Indicator Settings
	RegIrqFlags1      byte // Status register: PLL Lock state, Timeout, RSSI
	RegIrqFlags2      byte // Status register: FIFO handling flags, Low Battery
	RegDioMapping1    byte // Mapping of pins DIO0 to DIO3
	RegDioMapping2    byte // Mapping of pins DIO4 and DIO5, ClkOut frequency
	RegVersion        byte // Hope RF ID relating the silicon revision
}

// ResetRFConfiguration contains the register values after reset,
// according to data sheet section 6.
var ResetRFConfiguration = RFConfiguration{
	RegOpMode:         0x01,
	RegBitrateMsb:     0x1A,
	RegBitrateLsb:     0x0B,
	RegFdevMsb:        0x00,
	RegFdevLsb:        0x52,
	RegFrfMsb:         0x6C,
	RegFrfMid:         0x80,
	RegFrfLsb:         0x00,
	RegPaConfig:       0x4F,
	RegPaRamp:         0x09,
	RegOcp:            0x2B,
	RegLna:            0x20,
	RegRxConfig:       0x0E,
	RegRssiConfig:     0x02,
	RegRssiCollision:  0x0A,
	RegRssiThresh:     0xFF,
	RegRssiValue:      0x00,
	RegRxBw:           0x15,
	RegAfcBw:          0x0B,
	RegOokPeak:        0x28,
	RegOokFix:         0x0C,
	RegOokAvg:         0x12,
	reserved17:        0x47,
	reserved18:        0x32,
	reserved19:        0x3E,
	RegAfcFei:         0x00,
	RegAfcMsb:         0x00,
	RegAfcLsb:         0x00,
	RegFeiMsb:         0x00,
	RegFeiLsb:         0x00,
	RegPreambleDetect: 0x40,
	RegRxTimeout1:     0x00,
	RegRxTimeout2:     0x00,
	RegRxTimeout3:     0x00,
	RegRxDelay:        0x00,
	RegOsc:            0x07,
	RegPreambleMsb:    0x00,
	RegPreambleLsb:    0x03,
	RegSyncConfig:     0x93,
	RegSyncValue1:     0x55,
	RegSyncValue2:     0x55,
	RegSyncValue3:     0x55,
	RegSyncValue4:     0x55,
	RegSyncValue5:     0x55,
	RegSyncValue6:     0x55,
	RegSyncValue7:     0x55,
	RegSyncValue8:     0x55,
	RegPacketConfig1:  0x90,
	RegPacketConfig2:  0x40,
	RegPayloadLength:  0x40,
	RegNodeAdrs:       0x00,
	RegBroadcastAdrs:  0x00,
	RegFifoThresh:     0x0F,
	RegSeqConfig1:     0x00,
	RegSeqConfig2:     0x00,
	RegTimerResol:     0x00,
	RegTimer1Coef:     0xF5,
	RegTimer2Coef:     0x20,
	RegImageCal:       0x82,
	RegTemp:           0x00,
	RegLowBat:         0x02,
	RegIrqFlags1:      0x80,
	RegIrqFlags2:      0x40,
	RegDioMapping1:    0x00,
	RegDioMapping2:    0x00,
	RegVersion:        0x12,
}

// DefaultRFConfiguration contains the default (FSK) values,
// according to data sheet section 6.
var DefaultRFConfiguration = RFConfiguration{
	RegOpMode:         0x01,
	RegBitrateMsb:     0x1A,
	RegBitrateLsb:     0x0B,
	RegFdevMsb:        0x00,
	RegFdevLsb:        0x52,
	RegFrfMsb:         0x6C,
	RegFrfMid:         0x80,
	RegFrfLsb:         0x00,
	RegPaConfig:       0x4F,
	RegPaRamp:         0x09,
	RegOcp:            0x2B,
	RegLna:            0x20,
	RegRxConfig:       0x08,
	RegRssiConfig:     0x02,
	RegRssiCollision:  0x0A,
	RegRssiThresh:     0xFF,
	RegRssiValue:      0x00,
	RegRxBw:           0x15,
	RegAfcBw:          0x0B,
	RegOokPeak:        0x28,
	RegOokFix:         0x0C,
	RegOokAvg:         0x12,
	reserved17:        0x47,
	reserved18:        0x32,
	reserved19:        0x3E,
	RegAfcFei:         0x00,
	RegAfcMsb:         0x00,
	RegAfcLsb:         0x00,
	RegFeiMsb:         0x00,
	RegFeiLsb:         0x00,
	RegPreambleDetect: 0x40,
	RegRxTimeout1:     0x00,
	RegRxTimeout2:     0x00,
	RegRxTimeout3:     0x00,
	RegRxDelay:        0x00,
	RegOsc:            0x05,
	RegPreambleMsb:    0x00,
	RegPreambleLsb:    0x03,
	RegSyncConfig:     0x93,
	RegSyncValue1:     0x01,
	RegSyncValue2:     0x01,
	RegSyncValue3:     0x01,
	RegSyncValue4:     0x01,
	RegSyncValue5:     0x01,
	RegSyncValue6:     0x01,
	RegSyncValue7:     0x01,
	RegSyncValue8:     0x01,
	RegPacketConfig1:  0x90,
	RegPacketConfig2:  0x40,
	RegPayloadLength:  0x40,
	RegNodeAdrs:       0x00,
	RegBroadcastAdrs:  0x00,
	RegFifoThresh:     0x1F,
	RegSeqConfig1:     0x00,
	RegSeqConfig2:     0x00,
	RegTimerResol:     0x00,
	RegTimer1Coef:     0x12,
	RegTimer2Coef:     0x20,
	RegImageCal:       0x02,
	RegTemp:           0x00,
	RegLowBat:         0x02,
	RegIrqFlags1:      0x80,
	RegIrqFlags2:      0x40,
	RegDioMapping1:    0x00,
	RegDioMapping2:    0x00,
	RegVersion:        0x12,
}

// RegOpMode
const (
	FskOokMode = 0 << 7
	LoRaMode   = 1 << 7

	ModulationTypeMask = 3 << 5
	ModulationTypeFSK  = 0 << 5
	ModulationTypeOOK  = 1 << 5

	ModeMask        = 7
	SleepMode       = 0
	StandbyMode     = 1
	FreqSynthModeTX = 2
	TransmitterMode = 3
	FreqSynthModeRX = 4
	ReceiverMode    = 5
)

// RegLna
const (
	LnaGainMax      = 1 << 5
	LnaGainMax_6dB  = 2 << 5
	LnaGainMax_12dB = 3 << 5
	LnaGainMax_24dB = 4 << 5
	LnaGainMax_36dB = 5 << 5
	LnaGainMax_48dB = 6 << 5
)

// RegRxConfig
const (
	AfcAutoOn         = 1 << 4
	AgcAutoOn         = 1 << 3
	RxTriggerPreamble = 6 << 0
	RxTriggerRSSI     = 1 << 0
)

// RegRxBw
const (
	RxBwMantShift = 3
	RxBwMantMask  = 3 << 3
	RxBwMant16    = 0 << 3
	RxBwMant20    = 1 << 3
	RxBwMant24    = 2 << 3
	RxBwExpShift  = 0
	RxBwExpMask   = 7 << 0
)

// RegSyncConfig
const (
	SyncOn        = 1 << 4
	SyncSizeShift = 0
)

// RegPacketConfig1
const (
	FixedLength           = 0 << 7
	VariableLength        = 1 << 7
	DcFreeShift           = 5
	CrcOn                 = 1 << 4
	CrcOff                = 0 << 4
	CrcAutoClearOff       = 1 << 3
	AddressFilteringShift = 1
)

// RegPacketConfig2
const (
	PacketMode           = 1 << 6
	PayloadLengthMSBMask = 7
)

// RegFifoThresh
const (
	TxStartCondition = 1 << 7
)

// RegSeqConfig1
const (
	SequencerStart           = 1 << 7
	SequencerStop            = 1 << 6
	IdleModeSleep            = 1 << 5
	FromStartToTXOnFifoLevel = 3 << 3
)

// RegIrqFlags1
const (
	ModeReady        = 1 << 7
	RxReady          = 1 << 6
	TxReady          = 1 << 5
	PllLock          = 1 << 4
	Rssi             = 1 << 3
	Timeout          = 1 << 2
	PreambleDetect   = 1 << 1
	SyncAddressMatch = 1 << 0
)

// RegIrqFlags2
const (
	FifoFull     = 1 << 7
	FifoEmpty    = 1 << 6
	FifoLevel    = 1 << 5
	FifoOverrun  = 1 << 4
	PacketSent   = 1 << 3
	PayloadReady = 1 << 2
	CrcOk        = 1 << 1
	LowBat       = 1 << 0
)

// RegDioMapping1
const (
	Dio0MappingShift = 6
	Dio1MappingShift = 4
	Dio2MappingShift = 2
	Dio3MappingShift = 0
)

// RegDioMapping2
const (
	Dio4MappingShift  = 6
	Dio5MappingShift  = 4
	MapPreambleDetect = 1 << 0
	MapRssi           = 0 << 0
)
