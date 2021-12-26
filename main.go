package main

import (
	"emu/cpu"
	"emu/ppu"
)

func main() {
	// Read ROM
	path := "./ROM/sample1.nes"

	cpu.Exec(path)
	ppu.PpuTest()

}
