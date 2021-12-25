package cpu

import (
	"fmt"
	"io/ioutil"
	"os"
)

/*
   CPU memory map
   |---------------+-------+---------------------------------------------------------|
   | Address range | Size  | Device                                                  |
   |---------------+-------+---------------------------------------------------------|
   | $0000-$07FF   | $0800 | 2KB internal RAM                                        |
   | $0800-$0FFF   | $0800 | Mirrors of $0000-$07FF                                  |
   | $1000-$17FF   | $0800 |                                                         |
   | $1800-$1FFF   | $0800 |                                                         |
   | $2000-$2007   | $0008 | NES PPU registers                                       |
   | $2008-$3FFF   | $1FF8 | Mirrors of $2000-2007 (repeats every 8 bytes)           |
   | $4000-$4017   | $0018 | NES APU and I/O registers                               |
   | $4018-$401F   | $0008 | APU and I/O functionality that is normally disabled     |
   | $4020-$FFFF   | $BFE0 | Cartridge space: PRG ROM, PRG RAM, and mapper registers |
   | $6000-$7FFF   |       | Battery Backed Save or Work RAM                         |
   | $8000-$FFFF   |       | Usual ROM, commonly with Mapper Registers               |
   | $FFFA-$FFFB   |       | NMI vector                                              |
   | $FFFC-$FFFD   |       | Reset vector                                            |
   | $FFFE-$FFFF   |       | IRQ/BRK vector                                          |
   |---------------+-------+---------------------------------------------------------|

*/
var CPU_MEM [0xFFFF]byte
var PRG_ROM_ADDR uint16 = 0x8000

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
var PPU_MEM [0x3FFF]byte
var CHR_ROM_ADDR uint16 = 0x0000

/* register */
type Register struct {
	A  byte
	X  byte
	Y  byte
	SP uint16
	PC uint16
	P  byte
}

// Init register
func initRegister() *Register {
	init_reg := new(Register)
	init_reg.A = 0x00
	init_reg.X = 0x00
	init_reg.Y = 0x00
	init_reg.SP = 0x01FD
	init_reg.PC = 0x0000

	return init_reg
}

// Reset register
/*
   |-----------+------------+------------|
   | Interrupt | Lower Addr | Upper Addr |
   |-----------+------------+------------|
   | RESET     |     0xFFFC |     0xFFFD |
   |-----------+------------+------------|
*/
func resetRegister() *Register {
	lower_addr := 0xFFFC
	upper_addr := 0xFFFD
	reset_point := uint16(CPU_MEM[lower_addr]) + (uint16(CPU_MEM[upper_addr]) << 0x8)
	fmt.Printf("%x, %x\n", CPU_MEM[lower_addr], CPU_MEM[upper_addr])
	reset_reg := new(Register)
	reset_reg.A = 0x00
	reset_reg.X = 0x00
	reset_reg.Y = 0x00
	reset_reg.SP = 0x01FD
	reset_reg.PC = reset_point

	return reset_reg
}

// Fetch inst by PC
func fetchPC(reg Register) byte {
	// fmt.Printf("fetch:%x \n", CPU_MEM[reg.PC])
	// inc PC
	// reg.PC++
	return CPU_MEM[reg.PC]
}

// Execute loaded ROM
func Exec(path string) {
	loadRom(path)
	// Init and reset register
	reg := initRegister()
	reg = resetRegister()

	opcode := fetchPC(*reg)
	fmt.Printf("reset:%x \n", opcode)

}

// Read ROM and load to CPU/PPU memory
func loadRom(path string) {
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

	copy(CPU_MEM[PRG_ROM_ADDR:], prg_rom[:])
	copy(PPU_MEM[CHR_ROM_ADDR:], chr_rom[:])
}
