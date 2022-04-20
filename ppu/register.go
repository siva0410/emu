package ppu

import "github.com/siva0410/emu/cpu"

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
	Ppuctrl   *byte // mode:W
	Ppumask   *byte // mode:W
	Ppustatus *byte // mode:R
	Oamaddr   *byte // mode:W
	Oamdata   *byte // mode:R/W
	Ppuscroll *byte // mode:W
	Ppuaddr   *byte // mode:W
	Ppudata   *byte // mode:R/W
}

var Ppu_reg *PpuRegister

// Init Ppu register
func initPpuRegisters(reg *PpuRegister) {
	reg.Ppuctrl = &cpu.CPU_MEM[0x2000]
	reg.Ppumask = &cpu.CPU_MEM[0x2001]
	reg.Ppustatus = &cpu.CPU_MEM[0x2002]
	reg.Oamaddr = &cpu.CPU_MEM[0x2003]
	reg.Oamdata = &cpu.CPU_MEM[0x2004]
	reg.Ppuscroll = &cpu.CPU_MEM[0x2005]
	reg.Ppuaddr = &cpu.CPU_MEM[0x2006]
	reg.Ppudata = &cpu.CPU_MEM[0x2007]
}

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
func GetPpuCtrl(flagname string) bool {
	var status byte
	switch flagname {
	case "NM1":
		status = *Ppu_reg.Ppuctrl >> 0 & 0b1
	case "NM2":
		status = *Ppu_reg.Ppuctrl >> 1 & 0b1
	case "I":
		status = *Ppu_reg.Ppuctrl >> 2 & 0b1
	case "S":
		status = *Ppu_reg.Ppuctrl >> 3 & 0b1
	case "B":
		status = *Ppu_reg.Ppuctrl >> 4 & 0b1
	case "H":
		status = *Ppu_reg.Ppuctrl >> 5 & 0b1
	case "P":
		status = *Ppu_reg.Ppuctrl >> 6 & 0b1
	case "V":
		status = *Ppu_reg.Ppuctrl >> 7 & 0b1
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
func GetPpuMask(flagname string) bool {
	var status byte
	switch flagname {
	case "G":
		status = *Ppu_reg.Ppumask >> 0 & 0b1
	case "m":
		status = *Ppu_reg.Ppumask >> 1 & 0b1
	case "M":
		status = *Ppu_reg.Ppumask >> 2 & 0b1
	case "b":
		status = *Ppu_reg.Ppumask >> 3 & 0b1
	case "s":
		status = *Ppu_reg.Ppumask >> 4 & 0b1
	case "RGB_B":
		status = *Ppu_reg.Ppumask >> 5 & 0b1
	case "RGB_G":
		status = *Ppu_reg.Ppumask >> 6 & 0b1
	case "RGB_R":
		status = *Ppu_reg.Ppumask >> 7 & 0b1
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

var ppu_addr_flag bool = false

func CheckPpuPtr() {
	if cpu.CPU_MEM_CHK[0x2007] {
		PPU_MEM[PPU_PTR] = *Ppu_reg.Ppudata
		// check ppu_mem_chk
		PPU_MEM_CHK[PPU_PTR] = true
		if !GetPpuCtrl("I") {
			PPU_PTR += 0x1
		} else {
			PPU_PTR += 0x20
		}
		cpu.CPU_MEM_CHK[0x2007] = false
	}

	if cpu.CPU_MEM_CHK[0x2006] {
		PPU_PTR = PPU_PTR << 0x8
		PPU_PTR += uint32(*Ppu_reg.Ppuaddr)
		PPU_PTR &= 0xFFFF
		if ppu_addr_flag {
			ppu_addr_flag = false
		} else {
			ppu_addr_flag = true
		}
		cpu.CPU_MEM_CHK[0x2006] = false
	}
}
