package cpu

import (
	"fmt"

	"github.com/siva0410/emu/casette"
)

/*
   |---------------------+--------------|
   | Addressing mode     | Abbreviation |
   |---------------------+--------------|
   | zeroPage            | ZERO         |
   | relative            | REL          |
   | implied             | IMPL         |
   | absolute            | ABS          |
   | accumulator         | ACCUM        |
   | immediate           | IMM          |
   | zeroPageX           | ZEROX        |
   | zeroPageY           | ZEROY        |
   | absoluteX           | ABSX         |
   | absoluteY           | ABSY         |
   | preIndexedIndirect  | INDX         |
   | postIndexedIndirect | INDY         |
   | indirectAbsolute    | INDABS       |
   |---------------------+--------------|
*/
func getOperand(mode string) uint16 {
	var operand uint16
	var tmp uint16
	switch mode {
	case "IMPL", "ACCUM":

	case "IMM":
		operand = reg.PC
		reg.PC++

	case "ZERO":
		operand = uint16(fetchPC() & 0xFF)
		reg.PC++

	case "ZEROX":
		operand = uint16((fetchPC() + reg.X) & 0xFF)
		reg.PC++

	case "ZEROY":
		operand = uint16((fetchPC() + reg.Y) & 0xFF)
		reg.PC++

	case "ABS":
		tmp = uint16(fetchPC())
		reg.PC++
		operand = tmp + uint16(fetchPC())<<0x8
		reg.PC++

	case "ABSX":
		tmp = uint16(fetchPC())
		reg.PC++
		operand = (tmp + uint16(fetchPC())<<0x8) + uint16(reg.X)
		reg.PC++

	case "ABSY":
		tmp = uint16(fetchPC())
		reg.PC++
		operand = (tmp + uint16(fetchPC())<<0x8) + uint16(reg.Y)
		reg.PC++

	case "REL":
		tmp = uint16(fetchPC())
		reg.PC++
		if (tmp >> 7 & 1) == 1 {
			operand = uint16(reg.PC - (^tmp+0b1)&0xFF)
		} else {
			operand = uint16(reg.PC + tmp)
		}

	case "INDX":
		tmp = uint16((fetchPC() + reg.X) & 0xFF)
		reg.PC++
		operand = (tmp + uint16(fetchPC())<<0x8)
		reg.PC++

	case "INDY":
		tmp = uint16(fetchPC())
		reg.PC++
		operand = (tmp + uint16(fetchPC())<<0x8) + uint16(reg.Y)
		reg.PC++

	case "INDABS":
		tmp = uint16(fetchPC())
		reg.PC++
		tmp = tmp + uint16(fetchPC())<<0x8
		reg.PC++
		operand = (uint16(CPU_MEM[tmp]) + uint16(fetchPC())<<0x8)
		reg.PC++

	default:

	}
	return operand
}

func execOpecode(opecode byte) int {
	operand := getOperand(inst_arr[opecode].mode)
	var res byte
	switch inst_arr[opecode].name {
	case "LDA":
		reg.A = CPU_MEM[operand]
		setZeroFlag(reg.A)
		setNegativeFlag(reg.A)

	case "LDX":
		reg.X = CPU_MEM[operand]
		setZeroFlag(reg.X)
		setNegativeFlag(reg.X)

	case "LDY":
		reg.Y = CPU_MEM[operand]
		setZeroFlag(reg.Y)
		setNegativeFlag(reg.Y)

	case "STA":
		CPU_MEM[operand] = reg.A

	case "STX":
		CPU_MEM[operand] = reg.X

	case "STY":
		CPU_MEM[operand] = reg.Y

	case "TAX":
		reg.X = reg.A
		setZeroFlag(reg.X)
		setNegativeFlag(reg.X)

	case "TAY":
		reg.Y = reg.A
		setZeroFlag(reg.Y)
		setNegativeFlag(reg.Y)

	case "TXA":
		reg.A = reg.X
		setZeroFlag(reg.A)
		setNegativeFlag(reg.A)

	case "TYA":
		reg.A = reg.Y
		setZeroFlag(reg.A)
		setNegativeFlag(reg.A)

	case "TXS":
		reg.S = reg.X

	case "TSX":
		reg.X = reg.S
		setZeroFlag(reg.X)
		setNegativeFlag(reg.X)

	case "PHP":
		SP := 0x0100 + uint16(reg.S)&0xFF
		CPU_MEM[SP] = reg.P
		reg.S--

	case "PLP":
		reg.S++
		SP := 0x0100 + uint16(reg.S)&0xFF
		reg.P = CPU_MEM[SP]

	case "PHA":
		SP := 0x0100 + uint16(reg.S)&0xFF
		CPU_MEM[SP] = reg.A
		reg.S--

	case "PLA":
		reg.S++
		SP := 0x0100 + uint16(reg.S)&0xFF
		reg.A = CPU_MEM[SP]
		setZeroFlag(reg.A)
		setNegativeFlag(reg.A)

	case "ADC":
		overflow := (reg.A >> 6) & 0b1
		reg.A = reg.A + CPU_MEM[operand]
		if getStatus("carry") {
			reg.A += 1
		}
		if overflow != (reg.A>>6)&0b1 {
			setStatus("overflow", true)
		}
		setZeroFlag(reg.A)
		setNegativeFlag(reg.A)

	case "SBC":
		overflow := (reg.A >> 6) & 0b1
		reg.A = reg.A - CPU_MEM[operand]
		if !getStatus("carry") {
			reg.A -= 1
		}
		if overflow != (reg.A>>6)&0b1 {
			setStatus("overflow", true)
		}
		setZeroFlag(reg.A)
		setNegativeFlag(reg.A)

	case "CPX":
		res = reg.X - CPU_MEM[operand]
		setCarryFlag(res)
		setZeroFlag(res)
		setNegativeFlag(res)

	case "CPY":
		res = reg.Y - CPU_MEM[operand]
		setCarryFlag(res)
		setZeroFlag(res)
		setNegativeFlag(res)

	case "CMP":
		res = reg.A - CPU_MEM[operand]
		setCarryFlag(res)
		setZeroFlag(res)
		setNegativeFlag(res)

	case "AND":
		reg.A = reg.A & CPU_MEM[operand]
		setZeroFlag(reg.A)
		setNegativeFlag(reg.A)

	case "EOR":
		reg.A = reg.A ^ CPU_MEM[operand]
		setZeroFlag(reg.A)
		setNegativeFlag(reg.A)

	case "ORA":
		reg.A = reg.A | CPU_MEM[operand]
		setZeroFlag(reg.A)
		setNegativeFlag(reg.A)

	case "BIT":
		if CPU_MEM[operand]&reg.A == 0 {
			setStatus("zero", true)
		} else {
			setStatus("zero", false)
		}
		if CPU_MEM[operand]>>6 == 1 {
			setStatus("overflow", true)
		} else {
			setStatus("overflow", false)
		}
		setNegativeFlag(CPU_MEM[operand])

	case "ASL":
		if inst_arr[opecode].mode == "ACCUM" {
			if reg.A>>7 == 1 {
				setStatus("carry", true)
			} else {
				setStatus("carry", false)
			}
			reg.A = reg.A << 1
			setZeroFlag(reg.A)
			setNegativeFlag(reg.A)
		} else {
			if CPU_MEM[operand]>>7 == 1 {
				setStatus("carry", true)
			} else {
				setStatus("carry", false)
			}
			CPU_MEM[operand] = CPU_MEM[operand] << 1
			setZeroFlag(CPU_MEM[operand])
			setNegativeFlag(CPU_MEM[operand])
		}

	case "LSR":
		if inst_arr[opecode].mode == "ACCUM" {
			if reg.A&0b1 == 1 {
				setStatus("carry", true)
			} else {
				setStatus("carry", false)
			}
			reg.A = reg.A >> 1
			setZeroFlag(reg.A)
			setNegativeFlag(reg.A)
		} else {
			if CPU_MEM[operand]&0b1 == 1 {
				setStatus("carry", true)
			} else {
				setStatus("carry", false)
			}
			CPU_MEM[operand] = CPU_MEM[operand] >> 1
			setZeroFlag(CPU_MEM[operand])
			setNegativeFlag(CPU_MEM[operand])
		}

	case "ROL":
		if inst_arr[opecode].mode == "ACCUM" {
			if reg.A>>7 == 1 {
				setStatus("carry", true)
			} else {
				setStatus("carry", false)
			}
			reg.A = reg.A << 1
			if getStatus("carry") {
				reg.A += 1
			}
			setZeroFlag(reg.A)
			setNegativeFlag(reg.A)
		} else {
			if CPU_MEM[operand]>>7 == 1 {
				setStatus("carry", true)
			} else {
				setStatus("carry", false)
			}
			CPU_MEM[operand] = CPU_MEM[operand] << 1
			if getStatus("carry") {
				CPU_MEM[operand] += 1
			}
			setZeroFlag(CPU_MEM[operand])
			setNegativeFlag(CPU_MEM[operand])
		}

	case "ROR":
		if inst_arr[opecode].mode == "ACCUM" {
			if reg.A&0b1 == 1 {
				setStatus("carry", true)
			} else {
				setStatus("carry", false)
			}
			reg.A = reg.A >> 1
			if getStatus("carry") {
				reg.A += 1 << 7
			}
			setZeroFlag(reg.A)
			setNegativeFlag(reg.A)
		} else {
			if CPU_MEM[operand]&0b1 == 1 {
				setStatus("carry", true)
			} else {
				setStatus("carry", false)
			}
			CPU_MEM[operand] = CPU_MEM[operand] >> 1
			if getStatus("carry") {
				CPU_MEM[operand] += 1 << 7
			}
			setZeroFlag(CPU_MEM[operand])
			setNegativeFlag(CPU_MEM[operand])
		}

	case "INX":
		reg.X++
		setZeroFlag(reg.X)
		setNegativeFlag(reg.X)

	case "INY":
		reg.Y++
		setZeroFlag(reg.Y)
		setNegativeFlag(reg.Y)

	case "INC":
		CPU_MEM[operand]++
		setZeroFlag(CPU_MEM[operand])
		setNegativeFlag(CPU_MEM[operand])

	case "DEX":
		reg.X--
		setZeroFlag(reg.X)
		setNegativeFlag(reg.X)

	case "DEY":
		reg.Y--
		setZeroFlag(reg.Y)
		setNegativeFlag(reg.Y)

	case "DEC":
		CPU_MEM[operand]--
		setZeroFlag(CPU_MEM[operand])
		setNegativeFlag(CPU_MEM[operand])

	case "CLC":
		setStatus("carry", false)

	case "SEC":
		setStatus("carry", true)

	case "CLI":
		setStatus("interrupt_disable", false)

	case "SEI":
		setStatus("interrupt_disable", true)

	case "CLD":
		setStatus("decimal", false)

	case "SED":
		setStatus("decimal", true)

	case "CLV":
		setStatus("overflow", false)

	case "NOP":
		/* nop */

	case "BRK":
		if !getStatus("interrupt_disable") {
			setStatus("b_flag1", true)
			SP := 0x0100 + uint16(reg.S)&0xFF
			CPU_MEM[SP] = byte((reg.PC >> 8) & 0xFF)
			reg.S--
			SP = 0x0100 + uint16(reg.S)&0xFF
			CPU_MEM[SP] = byte(reg.PC & 0xFF)
			reg.S--
			SP = 0x0100 + uint16(reg.S)&0xFF
			CPU_MEM[SP] = reg.P
			reg.S--
			setStatus("interrupt_disable", true)
			reg.PC = uint16(CPU_MEM[0xFFFE]) & 0xFF
			reg.PC += uint16(CPU_MEM[0xFFFF]) << 8
		}

	case "RTI":
		reg.S++
		SP := 0x0100 + uint16(reg.S)&0xFF
		reg.P = CPU_MEM[SP]
		reg.S++
		SP = 0x0100 + uint16(reg.S)&0xFF
		reg.PC = uint16(CPU_MEM[SP]) & 0xFF
		reg.S++
		SP = 0x0100 + uint16(reg.S)&0xFF
		reg.PC += uint16(CPU_MEM[SP]) << 8

	case "JMP":
		reg.PC = operand

	case "JSR":
		SP := 0x0100 + uint16(reg.S)&0xFF
		CPU_MEM[SP] = byte((reg.PC >> 8) & 0xFF)
		reg.S--
		SP = 0x0100 + uint16(reg.S)&0xFF
		CPU_MEM[SP] = byte(reg.PC & 0xFF)
		reg.S--
		reg.PC = operand

	case "RTS":
		reg.S++
		SP := 0x0100 + uint16(reg.S)&0xFF
		reg.PC = uint16(CPU_MEM[SP]) & 0xFF
		reg.S++
		SP = 0x0100 + uint16(reg.S)&0xFF
		reg.PC += uint16(CPU_MEM[SP]) << 8

	case "BPL":
		if !getStatus("negative") {
			reg.PC = operand
		}

	case "BMI":
		if getStatus("negative") {
			reg.PC = operand
		}

	case "BVC":
		if !getStatus("overflow") {
			reg.PC = operand
		}

	case "BVS":
		if getStatus("overflow") {
			reg.PC = operand
		}

	case "BCC":
		if !getStatus("carry") {
			reg.PC = operand
		}

	case "BCS":
		if getStatus("carry") {
			reg.PC = operand
		}

	case "BNE":
		if !getStatus("zero") {
			reg.PC = operand
		}

	case "BEQ":
		if getStatus("zero") {
			reg.PC = operand
		}

	default:
		fmt.Println("NOT IMPL INST:", inst_arr[opecode].name, operand)
	}

	// fmt.Printf("NUM:0x%x\tOP:%s\tMODE:%s\tOPERAND:0x%x\n", opecode, inst_arr[opecode].name, inst_arr[opecode].mode, operand)
	// fmt.Printf("A:0x%x\tX:0x%x\tY:0x%x\tZERO:%v\t\n", reg.A, reg.X, reg.Y, getStatus("zero"))
	// // check ppu register
	// fmt.Printf("ppuctrl:%x\t", Ppu_reg.Ppuctrl)
	// fmt.Printf("ppustatus:%x\t", Ppu_reg.Ppustatus)
	// fmt.Printf("oamaddr:%x\t\n", Ppu_reg.Oamaddr)
	// fmt.Printf("oamdata:%x\t", Ppu_reg.Oamdata)
	// fmt.Printf("ppuaddr:%x\t", Ppu_reg.Ppuaddr)
	// fmt.Printf("ppudata:%x\t\n", Ppu_reg.Ppudata)
	// fmt.Printf("MEM ppuaddr:%x\n\n", PPU_PTR)

	CPU_MEM_CHK[operand] = true

	return inst_arr[opecode].cycle
}

// Init CPU
func InitCpu() {
	// Read Rom
	copy(CPU_MEM[PRG_ROM_ADDR:], casette.Prg_rom[:])

	// init register
	reg = initRegister()
	reg = resetRegister()

	// init inst list
	initInstList()
	setInstList()
}

// Execute loaded ROM
func ExecCpu(cycle *int) {
	// Execute ROM
	opecode := fetchPC()
	reg.PC++

	*cycle += execOpecode(opecode) * 3
}
