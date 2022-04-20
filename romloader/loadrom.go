package romloader

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/siva0410/emu/bus"
	"github.com/siva0410/emu/ppu"
)

// Read ROM and load to CPU/PPU memory
func LoadRom(path string) {
	f, err := os.Open(path)

	if err != nil {
		fmt.Printf("Cannot read %s", path)
	}

	defer f.Close()

	rom, err := ioutil.ReadAll(f)

	/*
	   Header (16 bytes)
	   PRG ROM data (16384 * x bytes)
	   CHR ROM data, if present (8192 * y bytes)
	*/
	HEADER_SIZE := 0x0010
	PRG_ROM_SIZE := 0x4000
	CHR_ROM_SIZE := 0x2000

	/*
	   Header
	   0-3:   Constant $4E $45 $53 $1A ("NES" followed by MS-DOS end-of-file)
	   4:     Size of PRG ROM in 16 KB units
	   5:     Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
	   6:     Flags 6 - Mapper, mirroring, battery, trainer
	   7:     Flags 7 - Mapper, VS/Playchoice, NES 2.0
	   8:     Flags 8 - PRG-RAM size (rarely used extension)
	   9:     Flags 9 - TV system (rarely used extension)
	   10:    Flags 10 - TV system, PRG-RAM presence (unofficial, rarely used extension)
	   11-15: Unused padding (should be filled with zero, but some rippers put their name across bytes 7-15)
	*/
	PRG_ROM_PAGES := int(rom[4])
	CHR_ROM_PAGES := int(rom[5])
	CHR_ROM_START := HEADER_SIZE + PRG_ROM_PAGES*PRG_ROM_SIZE

	prg_rom := rom[HEADER_SIZE : HEADER_SIZE+PRG_ROM_PAGES*PRG_ROM_SIZE]
	chr_rom := rom[CHR_ROM_START : CHR_ROM_START+CHR_ROM_PAGES*CHR_ROM_SIZE]

	copy(bus.CPU_MEM[bus.PRG_ROM_ADDR:], prg_rom[:])
	copy(ppu.PPU_MEM[bus.CHR_ROM_ADDR:], chr_rom[:])
}
