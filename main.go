package main

import (
	"emu/cpu"
	"emu/romloader"
)

func main() {
	// Read ROM
	path := "./ROM/sample1.nes"

	// Load ROM
	romloader.LoadRom(path)

	cpu.Exec()
	// ppu.PpuTest()
}
