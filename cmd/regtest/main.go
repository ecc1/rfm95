package main

import (
	"bytes"
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
	checkRegs(r)

	hw := r.Hardware()
	data := []byte{0x44, 0x55, 0x66}
	hw.WriteRegister(rfm95.RegSyncValue1, data[0])
	hw.WriteRegister(rfm95.RegSyncValue2, data[1])
	hw.WriteRegister(rfm95.RegSyncValue3, data[2])
	readRegs(r, "single", data)

	r.Reset()
	data = []byte{0x77, 0x88, 0x99}
	hw.WriteBurst(rfm95.RegSyncValue1, data)
	readRegs(r, "burst", data)
}

func checkRegs(r *rfm95.Radio) {
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	resetValue := rfm95.ResetConfiguration()
	regs0 := r.ReadConfiguration(false)
	regs1 := r.ReadConfiguration(true)
	if len(regs0) != len(resetValue) {
		log.Fatal("%d individual registers, expected %d", len(regs0), len(resetValue))
	}
	if len(regs1) != len(resetValue) {
		log.Fatal("%d burst-mode registers, expected %d", len(regs1), len(resetValue))
	}
	mismatches := 0
	for i, v := range regs0 {
		if regs1[i] != v {
			fmt.Printf("%02X  %02X  %08b (single) != %02X  %08b (burst)\n", i, v, v, regs1[i], regs1[i])
			mismatches++
		}
	}
	if mismatches == 0 {
		fmt.Printf("Burst-mode read is working correctly.\n")
	} else {
		fmt.Printf("WARNING: burst read did not match %d of %d single reads\n", mismatches, len(regs0))
	}
	fmt.Println("Configuration registers:")
	for i, v := range regs1 {
		fmt.Printf("%02X  %02X  %08b", i, v, v)
		r := resetValue[i]
		if v == r {
			fmt.Printf("\n")
		} else {
			fmt.Printf("  **** SHOULD BE %02X  %08b\n", r, r)
		}
	}
}

func readRegs(r *rfm95.Radio, kind string, data []byte) {
	fmt.Printf("\nTesting %s writes\n", kind)
	fmt.Printf("source: % X\n", data)
	hw := r.Hardware()
	x := hw.ReadRegister(rfm95.RegSyncValue1)
	y := hw.ReadRegister(rfm95.RegSyncValue2)
	z := hw.ReadRegister(rfm95.RegSyncValue3)
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	fmt.Printf("single: %02X %02X %02X\n", x, y, z)
	if x != data[0] || y != data[1] || z != data[2] {
		fmt.Printf("ERROR: single reads did not match %s writes\n", kind)
	}
	v := hw.ReadBurst(rfm95.RegSyncValue1, 3)
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	fmt.Printf(" burst: %02X %02X %02X\n", v[0], v[1], v[2])
	if !bytes.Equal(v, data) {
		fmt.Printf("ERROR: burst reads did not match %s writes\n", kind)
	}
}
