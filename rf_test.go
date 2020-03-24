package rfm95

import (
	"bytes"
	"testing"
)

func TestFrequency(t *testing.T) {
	cases := []struct {
		f       uint32
		b       []byte
		fApprox uint32 // 0 => equal to f
	}{
		{315000000, []byte{0x4E, 0xC0, 0x00}, 0},
		{434000000, []byte{0x6C, 0x80, 0x00}, 0},
		{868000000, []byte{0xD9, 0x00, 0x00}, 0},
		{915000000, []byte{0xE4, 0xC0, 0x00}, 0},
		// some that can't be represented exactly:
		{916300000, []byte{0xE5, 0x13, 0x33}, 916299987},
		{916600000, []byte{0xE5, 0x26, 0x66}, 916599975},
	}
	for _, c := range cases {
		b := frequencyToRegisters(c.f)
		if !bytes.Equal(b, c.b) {
			t.Errorf("frequencyToRegisters(%d) == % X, want % X", c.f, b, c.b)
		}
		f := registersToFrequency(c.b)
		if c.fApprox != 0 {
			if f != c.fApprox {
				t.Errorf("registersToFrequency(% X) == %d, want %d", c.b, f, c.fApprox)
			}
		} else {
			if f != c.f {
				t.Errorf("registersToFrequency(% X) == %d, want %d", c.b, f, c.f)
			}
		}
	}
}

// See data sheet Table 18 Bit Rate Examples
func TestBitrate(t *testing.T) {
	cases := []struct {
		br       uint32
		b        []byte
		brApprox uint32 // 0 => equal to br
	}{
		{1200, []byte{0x68, 0x2B}, 0},
		{2400, []byte{0x34, 0x15}, 0},
		{25000, []byte{0x05, 0x00}, 0},
		{50000, []byte{0x02, 0x80}, 0},
		// some that can't be represented exactly:
		{16384, []byte{0x07, 0xA1}, 16385},
		{19200, []byte{0x06, 0x83}, 19196},
		{38400, []byte{0x03, 0x41}, 38415},
		{150000, []byte{0x00, 0xD5}, 150235},
	}
	for _, c := range cases {
		b := bitrateToRegisters(c.br)
		if !bytes.Equal(b, c.b) {
			t.Errorf("bitrateToRegisters(%d) == % X, want % X", c.br, b, c.b)
		}
		f := registersToBitrate(c.b)
		if c.brApprox != 0 {
			if f != c.brApprox {
				t.Errorf("registersToBitrate(% X) == %d, want %d", c.b, f, c.brApprox)
			}
		} else {
			if f != c.br {
				t.Errorf("registersToBitrate(% X) == %d, want %d", c.b, f, c.br)
			}
		}
	}
}

// See data sheet Table 38 Available RxBw Settings
func TestChannelBW(t *testing.T) {
	cases := []struct {
		bw       uint32
		r        byte
		bwApprox uint32 // 0 => equal to bw
	}{
		{12500, RxBwMant20 | 5<<RxBwExpShift, 0},
		{25000, RxBwMant20 | 4<<RxBwExpShift, 0},
		{166666, RxBwMant24 | 1<<RxBwExpShift, 0},
		{200000, RxBwMant20 | 1<<RxBwExpShift, 0},
		{250000, RxBwMant16 | 1<<RxBwExpShift, 0},
		// some that can't be represented exactly:
		{0, RxBwMant24 | 7<<RxBwExpShift, 2604},
		{1000, RxBwMant24 | 7<<RxBwExpShift, 2604},
		{48000, RxBwMant20 | 3<<RxBwExpShift, 50000},
		{112500, RxBwMant20 | 2<<RxBwExpShift, 100000},
		{150000, RxBwMant24 | 1<<RxBwExpShift, 166666},
		{300000, RxBwMant16 | 1<<RxBwExpShift, 250000},
	}
	for _, c := range cases {
		r := channelBWToRegister(c.bw)
		if r != c.r {
			t.Errorf("channelBWToRegister(%d) == %02X, want %02X", c.bw, r, c.r)
		}
		bw := registerToChannelBW(c.r)
		if c.bwApprox != 0 {
			if bw != c.bwApprox {
				t.Errorf("registerToChannelBW(%02X) == %d, want %d", c.r, bw, c.bwApprox)
			}
		} else {
			if bw != c.bw {
				t.Errorf("registerToChannelBW(%02X) == %d, want %d", c.r, bw, c.bw)
			}
		}
	}
}
