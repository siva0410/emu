package ppu

import (
	"emu/bus"
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
)

/*
   PPU memory map
   |---------------+-------+------------------------+----------------------------------------|
   | Address range | Size  | Device                 |                                        |
   |---------------+-------+------------------------+----------------------------------------|
   | $0000-$0FFF   | $1000 | Pattern table 0        | mapped to CHR ROM                      |
   | $1000-$1FFF   | $1000 | Pattern table 1        |                                        |
   | $2000-$23FF   | $0400 | Nametable 0            | mapped to the 2kB NES internal VRAM    |
   | $2400-$27FF   | $0400 | Nametable 1            |                                        |
   | $2800-$2BFF   | $0400 | Nametable 2            |                                        |
   | $2C00-$2FFF   | $0400 | Nametable 3            |                                        |
   | $3000-$3EFF   | $0F00 | Mirrors of $2000-$2EFF | a mirror of the 2kB region             |
   | $3F00-$3F1F   | $0020 | Palette RAM indexes    | mapped to the internal palette control |
   | $3F20-$3FFF   | $00E0 | Mirrors of $3F00-$3F1F |                                        |
   |---------------+-------+------------------------+----------------------------------------|
*/
var PPU_MEM [0x4000]byte
var PPU_PTR uint32
var PPU_MEM_CHK [0x4000]bool

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
	reg.Ppuctrl = &bus.CPU_MEM[0x2000]
	reg.Ppumask = &bus.CPU_MEM[0x2001]
	reg.Ppustatus = &bus.CPU_MEM[0x2002]
	reg.Oamaddr = &bus.CPU_MEM[0x2003]
	reg.Oamdata = &bus.CPU_MEM[0x2004]
	reg.Ppuscroll = &bus.CPU_MEM[0x2005]
	reg.Ppuaddr = &bus.CPU_MEM[0x2006]
	reg.Ppudata = &bus.CPU_MEM[0x2007]
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

var dots [][]*dot

func InitPpu() {
	Ppu_reg = new(PpuRegister)
	initPpuRegisters(Ppu_reg)

	dots = makeDots()
}

func CheckPpuPtr(operand uint16) {
	switch operand {
	case 0x2007:
		PPU_MEM[PPU_PTR] = *Ppu_reg.Ppudata
		// check ppu_mem_chk
		PPU_MEM_CHK[PPU_PTR] = true
		if !GetPpuCtrl("I") {
			PPU_PTR += 0x1
		} else {
			PPU_PTR += 0x20
		}

	case 0x2006:
		PPU_PTR = PPU_PTR << 0x8
		PPU_PTR += uint32(*Ppu_reg.Ppuaddr)
		PPU_PTR &= 0xFFFF
		if ppu_addr_flag {
			ppu_addr_flag = false
		} else {
			ppu_addr_flag = true
		}

	default:

	}
}

var line int

func ExecPpu(cycle *int, window *glfw.Window) {
	// Set sprite
	if *cycle >= 341 {
		*cycle -= 341
		line++
		fmt.Println(line)
	}

	if (line+1)%8 == 0 && line < 240 {
		// set sprite
		sl := (line+1)/8 - 1
		for sw := 0; sw < 256/8; sw++ {
			sprite_num := PPU_MEM[0x2000+0x20*sl+sw]
			for l := 0; l < 8; l++ {
				for i := 0; i < 8; i++ {
					s := (PPU_MEM[0x10*int(sprite_num)+l] >> (7 - i)) & 0b1
					t := (PPU_MEM[0x08+0x10*int(sprite_num)+l] >> (7 - i)) & 0b1
					dots[sl*8+l][sw*8+i].sprite = s + t<<1
				}
			}
		}
	}

	if (line+1)%16 == 0 && line < 240 {
		// set palette
		pl := (line+1)/16 - 1
		for pw := 0; pw < 256/16; pw++ {
			for l := 0; l < 16; l++ {
				for i := 0; i < 16; i++ {
					switch {
					case i < 8 && l < 8:
						dots[pl*16+l][pw*16+i].palette = (PPU_MEM[0x2000+0x03C0+0x4*pl+int(pw/4)] >> 0) & 0b11
					case i >= 8 && l < 8:
						dots[pl*16+l][pw*16+i].palette = (PPU_MEM[0x2000+0x03C0+0x4*pl+int(pw/4)] >> 2) & 0b11
					case i < 8 && l >= 8:
						dots[pl*16+l][pw*16+i].palette = (PPU_MEM[0x2000+0x03C0+0x4*pl+int(pw/4)] >> 4) & 0b11
					case i >= 8 && l >= 8:
						dots[pl*16+l][pw*16+i].palette = (PPU_MEM[0x2000+0x03C0+0x4*pl+int(pw/4)] >> 6) & 0b11
					}
				}
			}
		}
	}

	if line == 262 {
		line = 0
		UpdatePalette()

		draw(dots)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
