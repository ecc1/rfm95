package rfm95

// Configuration for Raspberry Pi Zero W with Adafruit RFM95W bonnet:
// https://www.adafruit.com/product/4074
// https://learn.adafruit.com/adafruit-radio-bonnets/pinouts

const (
	spiDevice    = "/dev/spidev0.1"
	spiSpeed     = 6000000 // Hz
	interruptPin = 24      // GPIO for receive interrupts (DIO2)
	resetPin     = 25      // GPIO for hardware reset
)
