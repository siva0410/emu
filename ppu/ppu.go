package ppu

import (
	"emu/bus"
	"fmt"
)

/*
   |-------------+---------+-----------+------------------------------------------------------------------|
   | Common Name | Address | Bits      | Notes                                                            |
   |-------------+---------+-----------+------------------------------------------------------------------|
   | PPUCTRL     | %2000   | VPHB SINN | follow table                                                     |
   | PPUMASK     | %2001   | BGRs bMmG | follow table                                                     |
   | PPUSTATUS   | %2002   | VSO- ---- | follow table                                                     |
   | OAMADDR     | %2003   | aaaa aaaa | OAM read/write address                                           |
   | OAMDATA     | %2004   | dddd dddd | OAM data read/write                                              |
   | PPUSCROLL   | %2005   | xxxx xxxx | fine scroll position (two writes: X scroll, Y scroll)            |
   | PPUADDR     | %2006   | aaaa aaaa | PPU read/write address (two writes: most/least significant byte) |
   | PPUDATA     | %2007   | dddd dddd | PPU data read/write                                              |
   | OAMDMA      | %4014   | aaaa aaaa | OAM DMA high address                                             |
   |-------------+---------+-----------+------------------------------------------------------------------|
*/
type PpuRegister struct {
	ppuctrl   byte // mode:W
	ppumask   byte // mode:W
	ppustatus byte // mode:R
	oamaddr   byte // mode:W
	oamdata   byte // mode:R/W
	ppuscroll byte // mode:W
	ppuaddr   byte // mode:W
	ppudata   byte // mode:R/W
}

var ppu_reg, before_ppu_reg *PpuRegister
var is_ppu_ptr_read bool = false

/*
   |----------------------------+-----|
   | PPUCTRL                    | bit |
   |----------------------------+-----|
   | NMI enable (V)             |   7 |
   | PPU master/slave (P)       |   6 |
   | sprite height (H)          |   5 |
   | background tile select (B) |   4 |
   | sprite tile select (S)     |   3 |
   | increment mode (I)         |   2 |
   | nametable select (NN)      | 1-0 |
   |----------------------------+-----|
*/
func getPpuCtrl(flagname string) bool {
	var status byte
	switch flagname {
	case "NM1":
		status = ppu_reg.ppuctrl >> 0 & 0b1
	case "NM2":
		status = ppu_reg.ppuctrl >> 1 & 0b1
	case "I":
		status = ppu_reg.ppuctrl >> 2 & 0b1
	case "S":
		status = ppu_reg.ppuctrl >> 3 & 0b1
	case "B":
		status = ppu_reg.ppuctrl >> 4 & 0b1
	case "H":
		status = ppu_reg.ppuctrl >> 5 & 0b1
	case "P":
		status = ppu_reg.ppuctrl >> 6 & 0b1
	case "V":
		status = ppu_reg.ppuctrl >> 7 & 0b1
	}

	var res bool = false
	if status == 1 {
		res = true
	}

	return res
}

/*
   |-----------------------------------+-----|
   | PPUMASK                           | bit |
   |-----------------------------------+-----|
   | color emphasis (R)                |   7 |
   | color emphasis (G)                |   6 |
   | color emphasis (B)                |   5 |
   | sprite enable (s)                 |   4 |
   | background enable (b)             |   3 |
   | sprite left column enable (M)     |   2 |
   | background left column enable (m) |   1 |
   | greyscale (G)                     |   0 |
   |-----------------------------------+-----|
*/
func getPpuMask(flagname string) bool {
	var status byte
	switch flagname {
	case "G":
		status = ppu_reg.ppumask >> 0 & 0b1
	case "m":
		status = ppu_reg.ppumask >> 1 & 0b1
	case "M":
		status = ppu_reg.ppumask >> 2 & 0b1
	case "b":
		status = ppu_reg.ppumask >> 3 & 0b1
	case "s":
		status = ppu_reg.ppumask >> 4 & 0b1
	case "RGB_B":
		status = ppu_reg.ppumask >> 5 & 0b1
	case "RGB_G":
		status = ppu_reg.ppumask >> 6 & 0b1
	case "RGB_R":
		status = ppu_reg.ppumask >> 7 & 0b1
	}

	var res bool = false
	if status == 1 {
		res = true
	}

	return res
}

/*
   |-----------------------------------------+-----|
   | PPUSTATUS                               | bit |
   |-----------------------------------------+-----|
   | vblank (V)                              |   7 |
   | sprite 0 hit (S)                        |   6 |
   | sprite overflow (O)                     |   5 |
   | read resets write pair for %2005/%2006  |   4 |
   |-----------------------------------------+-----|
*/

func initPpuRegisters() *PpuRegister {
	res := new(PpuRegister)
	res.ppuctrl = bus.CPU_MEM[0x2000]
	res.ppumask = bus.CPU_MEM[0x2001]
	res.ppustatus = bus.CPU_MEM[0x2002]
	res.oamaddr = bus.CPU_MEM[0x2003]
	res.oamdata = bus.CPU_MEM[0x2004]
	res.ppuscroll = bus.CPU_MEM[0x2005]
	res.ppuaddr = bus.CPU_MEM[0x2006]
	res.ppudata = bus.CPU_MEM[0x2007]

	return res
}

func fetchPpuRegisters(reg *PpuRegister) {
	reg.ppuctrl = bus.CPU_MEM[0x2000]
	reg.ppumask = bus.CPU_MEM[0x2001]
	reg.ppustatus = bus.CPU_MEM[0x2002]
	reg.oamaddr = bus.CPU_MEM[0x2003]
	reg.oamdata = bus.CPU_MEM[0x2004]
	reg.ppuscroll = bus.CPU_MEM[0x2005]
	reg.ppuaddr = bus.CPU_MEM[0x2006]
	reg.ppudata = bus.CPU_MEM[0x2007]
}

func copyPpuRegisters(before_reg, reg *PpuRegister) {
	before_reg.ppuctrl = reg.ppuctrl
	before_reg.ppumask = reg.ppumask
	before_reg.ppustatus = reg.ppustatus
	before_reg.oamaddr = reg.oamaddr
	before_reg.oamdata = reg.oamdata
	before_reg.ppuscroll = reg.ppuscroll
	before_reg.ppuaddr = reg.ppuaddr
	before_reg.ppudata = reg.ppudata
}

func InitPpu() {
	ppu_reg = initPpuRegisters()
	before_ppu_reg = initPpuRegisters()
}

func ExecPpu() {
	fetchPpuRegisters(ppu_reg)
	if before_ppu_reg.ppuaddr != ppu_reg.ppuaddr {
		bus.PPU_PTR = bus.PPU_PTR << 0x8
		bus.PPU_PTR += uint32(ppu_reg.ppuaddr)
		bus.PPU_PTR &= 0xFFFF
		fmt.Printf("MEM ppuaddr:%x\n\n", bus.PPU_PTR)
		bus.PPU_MEM[bus.PPU_PTR] = ppu_reg.ppudata
	}

	// If ppuaddr is read, ppuaddr increments
	if is_ppu_ptr_read {
		if !getPpuCtrl("I") {
			bus.PPU_PTR += 0x1
		} else {
			bus.PPU_PTR += 0x20
		}
		is_ppu_ptr_read = false
	} else {
		is_ppu_ptr_read = true
	}

	// backup ppu_reg to before_ppu_reg
	copyPpuRegisters(before_ppu_reg, ppu_reg)

	// check ppu register
	fmt.Printf("ppuctrl:%x\t", ppu_reg.ppuctrl)
	fmt.Printf("ppustatus:%x\t", ppu_reg.ppustatus)
	fmt.Printf("oamaddr:%x\t\n", ppu_reg.oamaddr)
	fmt.Printf("oamdata:%x\t", ppu_reg.oamdata)
	fmt.Printf("ppuaddr:%x\t", ppu_reg.ppuaddr)
	fmt.Printf("ppudata:%x\t\n", ppu_reg.ppudata)
	fmt.Printf("MEM ppuaddr:%x\n\n", bus.PPU_PTR)

	fmt.Printf("before_ppuctrl:%x\t", before_ppu_reg.ppuctrl)
	fmt.Printf("before_ppustatus:%x\t", before_ppu_reg.ppustatus)
	fmt.Printf("before_oamaddr:%x\t\n", before_ppu_reg.oamaddr)
	fmt.Printf("before_oamdata:%x\t", before_ppu_reg.oamdata)
	fmt.Printf("before_ppuaddr:%x\t", before_ppu_reg.ppuaddr)
	fmt.Printf("before_ppudata:%x\t\n\n", before_ppu_reg.ppudata)

	// Palette()
}
