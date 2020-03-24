package rfm95

// Configuration for Intel Edison.

const (
	spiDevice    = "/dev/spidev5.1"
	spiSpeed     = 6000000 // Hz
	interruptPin = 46      // GPIO for receive interrupts
	resetPin     = 47      // GPIO for hardware reset
)
