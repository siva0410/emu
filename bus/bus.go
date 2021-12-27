package bus

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
var PPU_PTR uint32

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
	Ppuctrl   byte // mode:W
	Ppumask   byte // mode:W
	Ppustatus byte // mode:R
	Oamaddr   byte // mode:W
	Oamdata   byte // mode:R/W
	Ppuscroll byte // mode:W
	Ppuaddr   byte // mode:W
	Ppudata   byte // mode:R/W
}

var Ppu_reg, Before_ppu_reg *PpuRegister

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
		status = Ppu_reg.Ppuctrl >> 0 & 0b1
	case "NM2":
		status = Ppu_reg.Ppuctrl >> 1 & 0b1
	case "I":
		status = Ppu_reg.Ppuctrl >> 2 & 0b1
	case "S":
		status = Ppu_reg.Ppuctrl >> 3 & 0b1
	case "B":
		status = Ppu_reg.Ppuctrl >> 4 & 0b1
	case "H":
		status = Ppu_reg.Ppuctrl >> 5 & 0b1
	case "P":
		status = Ppu_reg.Ppuctrl >> 6 & 0b1
	case "V":
		status = Ppu_reg.Ppuctrl >> 7 & 0b1
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
		status = Ppu_reg.Ppumask >> 0 & 0b1
	case "m":
		status = Ppu_reg.Ppumask >> 1 & 0b1
	case "M":
		status = Ppu_reg.Ppumask >> 2 & 0b1
	case "b":
		status = Ppu_reg.Ppumask >> 3 & 0b1
	case "s":
		status = Ppu_reg.Ppumask >> 4 & 0b1
	case "RGB_B":
		status = Ppu_reg.Ppumask >> 5 & 0b1
	case "RGB_G":
		status = Ppu_reg.Ppumask >> 6 & 0b1
	case "RGB_R":
		status = Ppu_reg.Ppumask >> 7 & 0b1
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
