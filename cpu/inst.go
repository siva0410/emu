package cpu

var cycle_list = [0x100]int{
	/*0x00*/ 7, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
	/*0x10*/ 2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	/*0x20*/ 6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
	/*0x30*/ 2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	/*0x40*/ 6, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
	/*0x50*/ 2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	/*0x60*/ 6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
	/*0x70*/ 2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	/*0x80*/ 2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	/*0x90*/ 2, 6, 2, 6, 4, 4, 4, 4, 2, 4, 2, 5, 5, 4, 5, 5,
	/*0xA0*/ 2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	/*0xB0*/ 2, 5, 2, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
	/*0xC0*/ 2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	/*0xD0*/ 2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	/*0xE0*/ 2, 6, 3, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	/*0xF0*/ 2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
}

type InstList struct {
	name  string
	mode  string
	cycle int
}

var inst_arr [0x100]InstList

func initInstList() {
	for i := 0; i < 0x100; i++ {
		inst_arr[i] = InstList{"NOP", "NOP", cycle_list[i]}
	}
}

func setInstList() {
	inst_arr[0xA9] = InstList{"LDA", "IMM", cycle_list[0xA9]}
	inst_arr[0xA5] = InstList{"LDA", "ZERO", cycle_list[0xA5]}
	inst_arr[0xAD] = InstList{"LDA", "ABS", cycle_list[0xAD]}
	inst_arr[0xB5] = InstList{"LDA", "ZEROX", cycle_list[0xB5]}
	inst_arr[0xBD] = InstList{"LDA", "ABSX", cycle_list[0xBD]}
	inst_arr[0xB9] = InstList{"LDA", "ABSY", cycle_list[0xB9]}
	inst_arr[0xA1] = InstList{"LDA", "INDX", cycle_list[0xA1]}
	inst_arr[0xB1] = InstList{"LDA", "INDY", cycle_list[0xB1]}

	inst_arr[0xA2] = InstList{"LDX", "IMM", cycle_list[0xA2]}
	inst_arr[0xA6] = InstList{"LDX", "ZERO", cycle_list[0xA6]}
	inst_arr[0xAE] = InstList{"LDX", "ABS", cycle_list[0xAE]}
	inst_arr[0xB6] = InstList{"LDX", "ZEROY", cycle_list[0xB6]}
	inst_arr[0xBE] = InstList{"LDX", "ABSY", cycle_list[0xBE]}

	inst_arr[0xA0] = InstList{"LDY", "IMM", cycle_list[0xA0]}
	inst_arr[0xA4] = InstList{"LDY", "ZERO", cycle_list[0xA4]}
	inst_arr[0xAC] = InstList{"LDY", "ABS", cycle_list[0xAC]}
	inst_arr[0xB4] = InstList{"LDY", "ZEROX", cycle_list[0xB4]}
	inst_arr[0xBC] = InstList{"LDY", "ABSY", cycle_list[0xBC]}

	inst_arr[0x85] = InstList{"STA", "ZERO", cycle_list[0x85]}
	inst_arr[0x8D] = InstList{"STA", "ABS", cycle_list[0x8D]}
	inst_arr[0x95] = InstList{"STA", "ZEROX", cycle_list[0x95]}
	inst_arr[0x9D] = InstList{"STA", "ABSX", cycle_list[0x9D]}
	inst_arr[0x99] = InstList{"STA", "ABSY", cycle_list[0x99]}
	inst_arr[0x81] = InstList{"STA", "INDX", cycle_list[0x81]}
	inst_arr[0x91] = InstList{"STA", "INDY", cycle_list[0x91]}

	inst_arr[0x86] = InstList{"STX", "ZERO", cycle_list[0x86]}
	inst_arr[0x8E] = InstList{"STX", "ABS", cycle_list[0x8E]}
	inst_arr[0x96] = InstList{"STX", "ZEROY", cycle_list[0x96]}

	inst_arr[0x84] = InstList{"STY", "ZERO", cycle_list[0x84]}
	inst_arr[0x8C] = InstList{"STY", "ABS", cycle_list[0x8C]}
	inst_arr[0x94] = InstList{"STY", "ZEROX", cycle_list[0x94]}

	inst_arr[0x8A] = InstList{"TXA", "IMPL", cycle_list[0x8A]}

	inst_arr[0x98] = InstList{"TYA", "IMPL", cycle_list[0x98]}

	inst_arr[0x9A] = InstList{"TXS", "IMPL", cycle_list[0x9A]}

	inst_arr[0xBA] = InstList{"TSX", "IMPL", cycle_list[0xBA]}

	inst_arr[0xA8] = InstList{"TAY", "IMPL", cycle_list[0xA8]}

	inst_arr[0xAA] = InstList{"TAX", "IMPL", cycle_list[0xAA]}

	inst_arr[0x08] = InstList{"PHP", "IMPL", cycle_list[0x08]}

	inst_arr[0x28] = InstList{"PLP", "IMPL", cycle_list[0x28]}

	inst_arr[0x48] = InstList{"PHA", "IMPL", cycle_list[0x48]}

	inst_arr[0x68] = InstList{"PLA", "IMPL", cycle_list[0x68]}

	inst_arr[0x69] = InstList{"ADC", "IMM", cycle_list[0x69]}
	inst_arr[0x65] = InstList{"ADC", "ZERO", cycle_list[0x65]}
	inst_arr[0x6D] = InstList{"ADC", "ABS", cycle_list[0x6D]}
	inst_arr[0x75] = InstList{"ADC", "ZEROX", cycle_list[0x75]}
	inst_arr[0x7D] = InstList{"ADC", "ABSX", cycle_list[0x7D]}
	inst_arr[0x79] = InstList{"ADC", "ABSY", cycle_list[0x79]}
	inst_arr[0x61] = InstList{"ADC", "INDX", cycle_list[0x61]}
	inst_arr[0x71] = InstList{"ADC", "INDY", cycle_list[0x71]}

	inst_arr[0xE9] = InstList{"SBC", "IMM", cycle_list[0xE9]}
	inst_arr[0xE5] = InstList{"SBC", "ZERO", cycle_list[0xE5]}
	inst_arr[0xED] = InstList{"SBC", "ABS", cycle_list[0xED]}
	inst_arr[0xF5] = InstList{"SBC", "ZEROX", cycle_list[0xF5]}
	inst_arr[0xFD] = InstList{"SBC", "ABSX", cycle_list[0xFD]}
	inst_arr[0xF9] = InstList{"SBC", "ABSY", cycle_list[0xF9]}
	inst_arr[0xE1] = InstList{"SBC", "INDX", cycle_list[0xE1]}
	inst_arr[0xF1] = InstList{"SBC", "INDY", cycle_list[0xF1]}

	inst_arr[0xE0] = InstList{"CPX", "IMM", cycle_list[0xE0]}
	inst_arr[0xE4] = InstList{"CPX", "ZERO", cycle_list[0xE4]}
	inst_arr[0xEC] = InstList{"CPX", "ABS", cycle_list[0xEC]}

	inst_arr[0xC0] = InstList{"CPY", "IMM", cycle_list[0xC0]}
	inst_arr[0xC4] = InstList{"CPY", "ZERO", cycle_list[0xC4]}
	inst_arr[0xCC] = InstList{"CPY", "ABS", cycle_list[0xCC]}

	inst_arr[0xC9] = InstList{"CMP", "IMM", cycle_list[0xC9]}
	inst_arr[0xC5] = InstList{"CMP", "ZERO", cycle_list[0xC5]}
	inst_arr[0xCD] = InstList{"CMP", "ABS", cycle_list[0xCD]}
	inst_arr[0xD5] = InstList{"CMP", "ZEROX", cycle_list[0xD5]}
	inst_arr[0xDD] = InstList{"CMP", "ABSX", cycle_list[0xDD]}
	inst_arr[0xD9] = InstList{"CMP", "ABSY", cycle_list[0xD9]}
	inst_arr[0xC1] = InstList{"CMP", "INDX", cycle_list[0xC1]}
	inst_arr[0xD1] = InstList{"CMP", "INDY", cycle_list[0xD1]}

	inst_arr[0x29] = InstList{"AND", "IMM", cycle_list[0x29]}
	inst_arr[0x25] = InstList{"AND", "ZERO", cycle_list[0x25]}
	inst_arr[0x2D] = InstList{"AND", "ABS", cycle_list[0x2D]}
	inst_arr[0x35] = InstList{"AND", "ZEROX", cycle_list[0x35]}
	inst_arr[0x3D] = InstList{"AND", "ABSX", cycle_list[0x3D]}
	inst_arr[0x39] = InstList{"AND", "ABSY", cycle_list[0x39]}
	inst_arr[0x21] = InstList{"AND", "INDX", cycle_list[0x21]}
	inst_arr[0x31] = InstList{"AND", "INDY", cycle_list[0x31]}

	inst_arr[0x49] = InstList{"EOR", "IMM", cycle_list[0x49]}
	inst_arr[0x45] = InstList{"EOR", "ZERO", cycle_list[0x45]}
	inst_arr[0x4D] = InstList{"EOR", "ABS", cycle_list[0x4D]}
	inst_arr[0x55] = InstList{"EOR", "ZEROX", cycle_list[0x55]}
	inst_arr[0x5D] = InstList{"EOR", "ABSX", cycle_list[0x5D]}
	inst_arr[0x59] = InstList{"EOR", "ABSY", cycle_list[0x59]}
	inst_arr[0x41] = InstList{"EOR", "INDX", cycle_list[0x41]}
	inst_arr[0x51] = InstList{"EOR", "INDY", cycle_list[0x51]}

	inst_arr[0x09] = InstList{"ORA", "IMM", cycle_list[0x09]}
	inst_arr[0x05] = InstList{"ORA", "ZERO", cycle_list[0x05]}
	inst_arr[0x0D] = InstList{"ORA", "ABS", cycle_list[0x0D]}
	inst_arr[0x15] = InstList{"ORA", "ZEROX", cycle_list[0x15]}
	inst_arr[0x1D] = InstList{"ORA", "ABSX", cycle_list[0x1D]}
	inst_arr[0x19] = InstList{"ORA", "ABSY", cycle_list[0x19]}
	inst_arr[0x01] = InstList{"ORA", "INDX", cycle_list[0x01]}
	inst_arr[0x11] = InstList{"ORA", "INDY", cycle_list[0x11]}

	inst_arr[0x24] = InstList{"BIT", "ZERO", cycle_list[0x24]}
	inst_arr[0x2C] = InstList{"BIT", "ABS", cycle_list[0x2C]}

	inst_arr[0x0A] = InstList{"ASL", "ACCUM", cycle_list[0x0A]}
	inst_arr[0x06] = InstList{"ASL", "ZERO", cycle_list[0x06]}
	inst_arr[0x0E] = InstList{"ASL", "ABS", cycle_list[0x0E]}
	inst_arr[0x16] = InstList{"ASL", "ZEROX", cycle_list[0x16]}
	inst_arr[0x1E] = InstList{"ASL", "ABSX", cycle_list[0x1E]}

	inst_arr[0x4A] = InstList{"LSR", "ACCUM", cycle_list[0x4A]}
	inst_arr[0x46] = InstList{"LSR", "ZERO", cycle_list[0x46]}
	inst_arr[0x4E] = InstList{"LSR", "ABS", cycle_list[0x4E]}
	inst_arr[0x56] = InstList{"LSR", "ZEROX", cycle_list[0x56]}
	inst_arr[0x5E] = InstList{"LSR", "ABSX", cycle_list[0x5E]}

	inst_arr[0x2A] = InstList{"ROL", "ACCUM", cycle_list[0x2A]}
	inst_arr[0x26] = InstList{"ROL", "ZERO", cycle_list[0x26]}
	inst_arr[0x2E] = InstList{"ROL", "ABS", cycle_list[0x2E]}
	inst_arr[0x36] = InstList{"ROL", "ZEROX", cycle_list[0x36]}
	inst_arr[0x3E] = InstList{"ROL", "ABSX", cycle_list[0x3E]}

	inst_arr[0x6A] = InstList{"ROR", "ACCUM", cycle_list[0x6A]}
	inst_arr[0x66] = InstList{"ROR", "ZERO", cycle_list[0x66]}
	inst_arr[0x6E] = InstList{"ROR", "ABS", cycle_list[0x6E]}
	inst_arr[0x76] = InstList{"ROR", "ZEROX", cycle_list[0x76]}
	inst_arr[0x7E] = InstList{"ROR", "ABSX", cycle_list[0x7E]}

	inst_arr[0xE8] = InstList{"INX", "IMPL", cycle_list[0xE8]}

	inst_arr[0xC8] = InstList{"INY", "IMPL", cycle_list[0xC8]}

	inst_arr[0xE6] = InstList{"INC", "ZERO", cycle_list[0xE6]}
	inst_arr[0xEE] = InstList{"INC", "ABS", cycle_list[0xEE]}
	inst_arr[0xF6] = InstList{"INC", "ZEROX", cycle_list[0xF6]}
	inst_arr[0xFE] = InstList{"INC", "ABSX", cycle_list[0xFE]}

	inst_arr[0xCA] = InstList{"DEX", "IMPL", cycle_list[0xCA]}

	inst_arr[0x88] = InstList{"DEY", "IMPL", cycle_list[0x88]}

	inst_arr[0xC6] = InstList{"DEC", "ZERO", cycle_list[0xC6]}
	inst_arr[0xCE] = InstList{"DEC", "ABS", cycle_list[0xCE]}
	inst_arr[0xD6] = InstList{"DEC", "ZEROX", cycle_list[0xD6]}
	inst_arr[0xDE] = InstList{"DEC", "ABSX", cycle_list[0xDE]}

	inst_arr[0x18] = InstList{"CLC", "IMPL", cycle_list[0x18]}
	inst_arr[0x58] = InstList{"CLI", "IMPL", cycle_list[0x58]}
	inst_arr[0xB8] = InstList{"CLV", "IMPL", cycle_list[0xB8]}
	inst_arr[0xD8] = InstList{"CLD", "IMPL", cycle_list[0xD8]}

	inst_arr[0x38] = InstList{"SEC", "IMPL", cycle_list[0x38]}
	inst_arr[0x78] = InstList{"SEI", "IMPL", cycle_list[0x78]}
	inst_arr[0xF8] = InstList{"SED", "IMPL", cycle_list[0xF8]}

	inst_arr[0xEA] = InstList{"NOP", "IMPL", cycle_list[0xEA]}

	inst_arr[0x00] = InstList{"BRK", "IMPL", cycle_list[0x00]}

	inst_arr[0x20] = InstList{"JSR", "ABS", cycle_list[0x20]}
	inst_arr[0x4C] = InstList{"JMP", "ABS", cycle_list[0x4C]}
	inst_arr[0x6C] = InstList{"SEC", "INDABS", cycle_list[0x6C]}

	inst_arr[0x40] = InstList{"RTI", "IMPL", cycle_list[0x40]}
	inst_arr[0x60] = InstList{"RTS", "IMPL", cycle_list[0x60]}

	inst_arr[0x10] = InstList{"BPL", "REL", cycle_list[0x10]}
	inst_arr[0x30] = InstList{"BMI", "REL", cycle_list[0x30]}
	inst_arr[0x50] = InstList{"BVC", "REL", cycle_list[0x50]}
	inst_arr[0x70] = InstList{"BVS", "REL", cycle_list[0x70]}
	inst_arr[0x90] = InstList{"BCC", "REL", cycle_list[0x90]}
	inst_arr[0xB0] = InstList{"BCS", "REL", cycle_list[0xB0]}
	inst_arr[0xD0] = InstList{"BNE", "REL", cycle_list[0xD0]}
	inst_arr[0xF0] = InstList{"BEQ", "REL", cycle_list[0xF0]}

}
