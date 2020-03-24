package rfm95

import (
	"log"
	"time"
)

const (
	verbose       = false
	maxPacketSize = 110
	fifoSize      = 66

	// The fifoThreshold value should allow a maximum-sized packet to be
	// written in two bursts, but be large enough to avoid fifo underflow.
	fifoThreshold = 20

	// Approximate time for one byte to be transmitted, based on the data rate.
	byteDuration = time.Millisecond
)

func init() {
	if verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	}
}

// Send transmits the given packet.
func (r *Radio) Send(data []byte) {
	if r.Error() != nil {
		return
	}
	count := len(data)
	if count > maxPacketSize {
		log.Panicf("attempting to send %d-byte packet", count)
	}
	if verbose {
		log.Printf("sending %d-byte packet in %s state", count, r.State())
	}
	// Terminate packet with zero byte.
	count++
	packet := make([]byte, count)
	copy(packet, data)
	// Change to Standby mode in case an earlier receive timeout left the radio in Receive mode.
	r.setMode(StandbyMode)
	r.clearFIFO()
	// Automatically enter Transmit state on FifoLevel interrupt.
	r.hw.WriteRegister(RegFifoThresh, TxStartCondition)
	r.hw.WriteRegister(RegSeqConfig1, SequencerStart|IdleModeSleep|FromStartToTXOnFifoLevel)
	// Specify fixed length packet format (including final zero byte)
	// so PacketSent interrupt will terminate Transmit state.
	r.hw.WriteRegister(RegPacketConfig1, FixedLength)
	lengthMSB := uint8(count>>8) & PayloadLengthMSBMask
	r.hw.WriteRegister(RegPacketConfig2, PacketMode|lengthMSB)
	r.hw.WriteRegister(RegPayloadLength, uint8(count))
	r.transmit(packet)
	r.setMode(SleepMode)
}

func (r *Radio) transmit(data []byte) {
	n := len(data)
	if n > fifoSize {
		n = fifoSize
	}
	if verbose {
		log.Printf("writing %d bytes to TX FIFO\n", n)
	}
	r.hw.WriteBurst(RegFifo, data[:n])
	data = data[n:]
	for r.Error() == nil {
		if len(data) == 0 {
			break
		}
		r.waitForFifoNonEmpty()
		r.hw.WriteRegister(RegFifo, data[0])
		data = data[1:]
	}
	// Wait for automatic return to standby mode when FIFO is empty.
	for r.Error() == nil {
		s := r.mode()
		if s == StandbyMode {
			break
		}
		if verbose || s != TransmitterMode {
			log.Printf("waiting for TX to finish in %s state", stateName(s))
		}
		time.Sleep(byteDuration)
	}
	r.sequencerStop()
	r.setMode(SleepMode)
}

func (r *Radio) waitForFifoNonEmpty() {
	for r.Error() == nil {
		if !r.fifoFull() {
			return
		}
	}
}
func (r *Radio) sequencerStop() {
	r.hw.WriteRegister(RegSeqConfig1, SequencerStop)
}

func (r *Radio) fifoEmpty() bool {
	return r.hw.ReadRegister(RegIrqFlags2)&FifoEmpty != 0
}

func (r *Radio) fifoFull() bool {
	return r.hw.ReadRegister(RegIrqFlags2)&FifoFull != 0
}

func (r *Radio) clearFIFO() {
	r.hw.WriteRegister(RegIrqFlags2, FifoOverrun)
}

// Receive listens with the given timeout for an incoming packet.
// It returns the packet and the associated RSSI.
func (r *Radio) Receive(timeout time.Duration) ([]byte, int) {
	if r.Error() != nil {
		return nil, 0
	}
	// Use unlimited length packet format (data sheet section 4.2.13.2).
	r.hw.WriteRegister(RegPacketConfig1, FixedLength)
	r.hw.WriteRegister(RegPayloadLength, 0)
	r.setMode(ReceiverMode)
	defer r.setMode(SleepMode)
	if verbose {
		log.Printf("waiting for interrupt in %s state", r.State())
	}
	r.hw.AwaitInterrupt(timeout)
	rssi := r.ReadRSSI()
	for r.Error() == nil {
		if r.fifoEmpty() {
			if timeout <= 0 {
				break
			}
			time.Sleep(byteDuration)
			timeout -= byteDuration
			continue
		}
		c := r.hw.ReadRegister(RegFifo)
		if r.Error() != nil {
			break
		}
		if c == 0 {
			// End of packet.
			return r.finishRX(rssi)
		}
		r.err = r.receiveBuffer.WriteByte(c)
	}
	return nil, rssi
}

func (r *Radio) finishRX(rssi int) ([]byte, int) {
	r.setMode(SleepMode)
	r.clearFIFO()
	size := r.receiveBuffer.Len()
	if size == 0 {
		return nil, rssi
	}
	p := make([]byte, size)
	_, r.err = r.receiveBuffer.Read(p)
	if r.Error() != nil {
		return nil, rssi
	}
	r.receiveBuffer.Reset()
	// Remove spurious final byte consisting of just one or two high bits.
	b := p[len(p)-1]
	if b == 0x80 || b == 0xC0 {
		log.Printf("end-of-packet glitch %X with RSSI %d", b, rssi)
		p = p[:len(p)-1]
	}
	if verbose {
		log.Printf("received %d-byte packet in %s state", size, r.State())
	}
	return p, rssi
}

// SendAndReceive transmits the given packet,
// then listens with the given timeout for an incoming packet.
// It returns the packet and the associated RSSI.
// (This could be further optimized by using an Automode to go directly
// from TX to RX, rather than returning to standby in between.)
func (r *Radio) SendAndReceive(data []byte, timeout time.Duration) ([]byte, int) {
	r.Send(data)
	if r.Error() != nil {
		return nil, 0
	}
	return r.Receive(timeout)
}
