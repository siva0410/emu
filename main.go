package main

import (
	"emu/cpu"
	"emu/ppu"
	"emu/romloader"
	"fmt"
)

func printMem() {
	// Print PPU_MEM
	fmt.Printf("ADDR| ")
	for i := 0; i < 0x20; i++ {
		if i == 0x10 {
			fmt.Printf("| ")
		}
		fmt.Printf("%02x ", i)
	}
	fmt.Printf("\n")

	fmt.Printf("----+-")
	for i := 0; i < 0x20; i++ {
		if i == 0x10 {
			fmt.Printf("+-")
		}
		fmt.Printf("---")
	}

	for i, m := range ppu.PPU_MEM {
		if i == 0x1000 || i == 0x2000 || i == 0x23c0 || i == 0x2400 || i == 0x27c0 || i == 0x2800 || i == 0x2bc0 || i == 0x2c00 || i == 0x2fc0 || i == 0x3000 || i == 0x3f00 || i == 0x3f20 {
			fmt.Printf("\n----+-")
			for i := 0; i < 0x20; i++ {
				if i == 0x10 {
					fmt.Printf("+-")
				}
				fmt.Printf("---")
			}
		}
		if i%0x20 == 0 {
			fmt.Printf("\n")
			fmt.Printf("%04x| ", i)
		}
		if i%0x20 != 0 && i%0x10 == 0 {
			fmt.Printf("| ")
		}

		if ppu.PPU_MEM_CHK[i] {
			fmt.Printf("\x1b[31m")
		}
		fmt.Printf("%02x ", m)
		if ppu.PPU_MEM_CHK[i] {
			fmt.Printf("\x1b[0m")
		}

	}
	fmt.Printf("\n")
}

func main() {
	// Read ROM
	path := "./ROM/helloworld/helloworld.nes"
	// path := "./ROM/tkshoot/SHOOT.nes"
	romloader.LoadRom(path)

	// Init CPU and PPU
	cpu.InitCpu()
	ppu.InitPpu()

	// Create window
	// window.Window()

	var cycle *int
	cycle = new(int)
	for i := 0; i < 200; i++ {
		// Exec CPU and PPU
		// PPU clock = 3*CPU clock
		fmt.Printf("#%d: cycle: %d\n", i, *cycle)
		cpu.ExecCpu(cycle)
		for j := 0; j < 3; j++ {
			ppu.ExecPpu(cycle)
		}
	}

	printMem()
}
