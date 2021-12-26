package main

import (
	"emu/cpu"
	"emu/ppu"
	"emu/romloader"
)

func main() {
	// Read ROM
	path := "./ROM/sample1.nes"
	romloader.LoadRom(path)

	// Exec CPU and PPU
	cpu.Exec()
	ppu.PpuTest()
}
