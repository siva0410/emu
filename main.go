package main

import (
	"emu/cpu"
	"emu/ppu"
	"emu/romloader"
	"fmt"
)

func main() {
	// Read ROM
	path := "./ROM/helloworld/helloworld.nes"
	// path := "./ROM/tkshoot/SHOOT.nes"
	romloader.LoadRom(path)

	// Init CPU and PPU
	cpu.InitCpu()
	ppu.InitPpu()

	// window.Window()

	for i := 0; i < 200; i++ {
		// Exec CPU and PPU
		fmt.Printf("%d:\n", i)
		cpu.ExecCpu()
		// PPU clock = 3*CPU clock
		for j := 0; j < 3; j++ {
			ppu.ExecPpu()
		}
	}
}
