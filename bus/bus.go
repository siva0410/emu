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
var CHR_ROM_ADDR uint16 = 0x0000
