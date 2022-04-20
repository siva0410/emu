package ppu

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/siva0410/emu/casette"
)

var dots [][]*dot

var line int

func InitPpu() {
	// Read Rom
	copy(PPU_MEM[CHR_ROM_ADDR:], casette.Chr_rom[:])

	Ppu_reg = new(PpuRegister)
	initPpuRegisters(Ppu_reg)

	dots = makeDots()
	line = 0
}

func ExecPpu(cycle *int, window *glfw.Window) {
	// Update ppu register
	CheckPpuPtr()

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
			palette_num := PPU_MEM[0x2000+0x03C0+0x4*pl+int(pw/4)]
			for l := 0; l < 16; l++ {
				for i := 0; i < 16; i++ {
					switch {
					case i < 8 && l < 8:
						dots[pl*16+l][pw*16+i].palette = (palette_num >> 0) & 0b11
					case i >= 8 && l < 8:
						dots[pl*16+l][pw*16+i].palette = (palette_num >> 2) & 0b11
					case i < 8 && l >= 8:
						dots[pl*16+l][pw*16+i].palette = (palette_num >> 4) & 0b11
					case i >= 8 && l >= 8:
						dots[pl*16+l][pw*16+i].palette = (palette_num >> 6) & 0b11
					}
				}
			}
		}
	}

	if line == 262 {
		line = 0
		updatePalette()

		draw(dots)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
