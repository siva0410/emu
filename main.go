package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Global

/*
   memory map
   Address range	Size	Device
   $0000-$07FF	$0800	2KB internal RAM
   $0800-$0FFF	$0800	Mirrors of $0000-$07FF
   $1000-$17FF	$0800
   $1800-$1FFF	$0800
   $2000-$2007	$0008	NES PPU registers
   $2008-$3FFF	$1FF8	Mirrors of $2000-2007 (repeats every 8 bytes)
   $4000-$4017	$0018	NES APU and I/O registers
   $4018-$401F	$0008	APU and I/O functionality that is normally disabled. See
   CPU Test Mode.
   $4020-$FFFF	$BFE0	Cartridge space: PRG ROM, PRG RAM, and mapper registers
   $6000-$7FFF = Battery Backed Save or Work RAM
   $8000-$FFFF = Usual ROM, commonly with Mapper Registers (see MMC1 and UxROM for example)
   $FFFA-$FFFB = NMI vector
   $FFFC-$FFFD = Reset vector
   $FFFE-$FFFF = IRQ/BRK vector
*/
var CPU_MEM [0xFFFF]byte
var PPU_MEM [0xFFFF]byte

var PRG_ROM_ADDR int = 0x8000
var CHR_ROM_ADDR int

func main() {
	// Read ROM
	path := "./ROM/sample1.nes"
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
	// CHR_ROM_SIZE := 0x2000

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
	// CHR_ROM_PAGES := int(rom[5])
	// CHR_ROM_START := HEADER_SIZE + PRG_ROM_PAGES*PRG_ROM_SIZE

	prg_rom := rom[HEADER_SIZE : HEADER_SIZE+PRG_ROM_PAGES*PRG_ROM_SIZE]
	// chr_rom := rom[CHR_ROM_START : CHR_ROM_START+CHR_ROM_PAGES*CHR_ROM_SIZE]

	copy(CPU_MEM[PRG_ROM_ADDR:], prg_rom[:])
	// copy(CPU_MEM[CHR_ROM_ADDR:], chr_rom[:])
	fmt.Println(prg_rom, CPU_MEM[PRG_ROM_ADDR:])

}
