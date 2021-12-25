package main

import (
	"emu/cpu"
)

func main() {
	// Read ROM
	path := "./ROM/sample1.nes"

	cpu.Exec(path)

}
