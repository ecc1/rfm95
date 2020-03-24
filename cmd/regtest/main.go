package main

import (
	"fmt"
	"log"

	"github.com/ecc1/rfm95"
)

func main() {
	r := rfm95.Open()
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	r.Reset()
	dumpRegs(r)

	fmt.Printf("\nTesting individual writes\n")
	hw := r.Hardware()
	hw.WriteRegister(rfm95.RegSyncValue1, 0x44)
	hw.WriteRegister(rfm95.RegSyncValue2, 0x55)
	hw.WriteRegister(rfm95.RegSyncValue3, 0x66)
	readRegs(r)

	r.Reset()
	fmt.Printf("\nTesting burst writes\n")
	hw.WriteBurst(rfm95.RegSyncValue1, []byte{0x77, 0x88, 0x99})
	readRegs(r)
}

func dumpRegs(r *rfm95.Radio) {
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	fmt.Printf("\nConfiguration registers:\n")
	regs := r.ReadConfiguration().Bytes()
	resetValue := rfm95.ResetRFConfiguration.Bytes()
	for i, v := range regs {
		fmt.Printf("%02X  %02X  %08b", rfm95.RegOpMode+i, v, v)
		r := resetValue[i]
		if v == r {
			fmt.Printf("\n")
		} else {
			fmt.Printf("  **** SHOULD BE %02X  %08b\n", r, r)
		}
	}
}

func readRegs(r *rfm95.Radio) {
	hw := r.Hardware()
	x := hw.ReadRegister(rfm95.RegSyncValue1)
	y := hw.ReadRegister(rfm95.RegSyncValue2)
	z := hw.ReadRegister(rfm95.RegSyncValue3)
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	fmt.Printf("individual: %X %X %X\n", x, y, z)
	v := hw.ReadBurst(rfm95.RegSyncValue1, 3)
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	fmt.Printf("  burst:    %X %X %X\n", v[0], v[1], v[2])
}
