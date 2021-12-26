package ppu

import (
	"fmt"

	"emu/bus"
)

type PpuRegister struct {
	a int
}

// func fetchPpuRegisters() {

// }

func PpuTest() {
	fmt.Println("PPU TEST!!")
	fmt.Println(bus.PPU_MEM)
}
