package cpu

import (
	"fmt"

	"emu/bus"
	"emu/ppu"
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
		tmp = uint16((fetchPC() + reg.Y) & 0xFF)
		reg.PC++
		operand = (tmp + uint16(fetchPC())<<0x8)
		reg.PC++

	case "INDABS":
		tmp = uint16(fetchPC())
		reg.PC++
		tmp = tmp + uint16(fetchPC())<<0x8
		reg.PC++
		operand = (uint16(bus.CPU_MEM[tmp])&0xFF + uint16(fetchPC())<<0x8)
		reg.PC++

	default:

	}
	return operand
}

func execOpecode(opecode byte) {
	operand := getOperand(inst_arr[opecode].mode)
	var res uint
	switch inst_arr[opecode].name {
	case "LDA":
		res = uint(operand)
		setZeroFlag(res)
		setNegativeFlag(res)
		reg.A = bus.CPU_MEM[operand]
		// check ppu_addr register
		ppu.CheckPpuPtr(operand)

	case "LDX":
		res = uint(operand)
		setZeroFlag(res)
		setNegativeFlag(res)
		reg.X = bus.CPU_MEM[operand]
		// check ppu_addr register
		ppu.CheckPpuPtr(operand)

	case "LDY":
		res = uint(operand)
		setZeroFlag(res)
		setNegativeFlag(res)
		reg.Y = bus.CPU_MEM[operand]
		// check ppu_addr register
		ppu.CheckPpuPtr(operand)

	case "STA":
		bus.CPU_MEM[operand] = reg.A
		// check ppu_addr register
		ppu.CheckPpuPtr(operand)

	case "STX":
		bus.CPU_MEM[operand] = reg.X
		// check ppu_addr register
		ppu.CheckPpuPtr(operand)

	case "STY":
		bus.CPU_MEM[operand] = reg.Y
		// check ppu_addr register
		ppu.CheckPpuPtr(operand)

	case "TAX":
		res = uint(reg.A)
		setZeroFlag(res)
		setNegativeFlag(res)
		reg.X = reg.A

	case "TAY":
		res = uint(reg.A)
		setZeroFlag(res)
		setNegativeFlag(res)
		reg.Y = reg.A

	case "TXA":
		res = uint(reg.X)
		setZeroFlag(res)
		setNegativeFlag(res)
		reg.A = reg.X

	case "TYA":
		res = uint(reg.Y)
		setZeroFlag(res)
		setNegativeFlag(res)
		reg.A = reg.Y

	case "TXS":
		res = uint(reg.X)
		setZeroFlag(res)
		setNegativeFlag(res)
		reg.S = reg.X

	case "TSX":
		res = uint(reg.S)
		setZeroFlag(res)
		setNegativeFlag(res)
		reg.X = reg.S

	// case "PHP":
	// case "PLP":
	// case "PHA":
	// case "PLA":
	// case "ADC":
	// 	reg.A = reg.A + bus.CPU_MEM[operand] + getStatus("carry")
	// case "SBC":
	// case "CPX":
	// case "CPY":
	// case "CMP":
	// case "AND":
	// 	reg.A = reg.A & bus.CPU_MEM[operand]
	// case "EOR":
	// case "ORA":
	// case "BIT":
	// case "ASL":
	// case "LSR":
	// case "ROL":
	// case "ROR":
	case "INX":
		reg.X++
		res = uint(reg.X)
		setZeroFlag(res)
		setNegativeFlag(res)

	case "INY":
		reg.Y++
		res = uint(reg.Y)
		setZeroFlag(res)
		setNegativeFlag(res)

	case "INC":
		bus.CPU_MEM[operand]++
		res = uint(bus.CPU_MEM[operand])
		setZeroFlag(res)
		setNegativeFlag(res)

	case "DEX":
		reg.X--
		res = uint(reg.X)
		setZeroFlag(res)
		setNegativeFlag(res)

	case "DEY":
		reg.Y--
		res = uint(reg.Y)
		setZeroFlag(res)
		setNegativeFlag(res)

	case "DEC":
		bus.CPU_MEM[operand]--
		res = uint(bus.CPU_MEM[operand])
		setZeroFlag(res)
		setNegativeFlag(res)

	// case "CLC":
	// case "CLI":
	// case "CLV":
	// case "CLD":
	// case "SEC":
	case "SEI":
		setStatus("interrupt_disable", false)
	// case "SED":
	// case "NOP":
	case "BRK":
		if !getStatus("interrupt_disable") {
			// interrupt handler
		}
	// case "JSR":
	case "JMP":
		reg.PC = operand
	// case "RTI":
	// case "RTS":
	// case "BPL":
	// case "BMI":
	// case "BVC":
	// case "BVS":
	// case "BCC":
	// case "BCS":
	case "BNE":
		if !getStatus("zero") {
			reg.PC = operand
		}
	// case "BEQ":
	default:
		fmt.Println("NOT IMPL INST:", inst_arr[opecode].name, operand)
	}

	fmt.Printf("NUM:0x%x\tOP:%s\tMODE:%s\tOPERAND:0x%x\n", opecode, inst_arr[opecode].name, inst_arr[opecode].mode, operand)
	fmt.Printf("A:0x%x\tX:0x%x\tY:0x%x\tZERO:%v\t\n", reg.A, reg.X, reg.Y, getStatus("zero"))
	// check ppu register
	fmt.Printf("ppuctrl:%x\t", ppu.Ppu_reg.Ppuctrl)
	fmt.Printf("ppustatus:%x\t", ppu.Ppu_reg.Ppustatus)
	fmt.Printf("oamaddr:%x\t\n", ppu.Ppu_reg.Oamaddr)
	fmt.Printf("oamdata:%x\t", ppu.Ppu_reg.Oamdata)
	fmt.Printf("ppuaddr:%x\t", ppu.Ppu_reg.Ppuaddr)
	fmt.Printf("ppudata:%x\t\n", ppu.Ppu_reg.Ppudata)
	fmt.Printf("MEM ppuaddr:%x\n\n", ppu.PPU_PTR)

}

// Init CPU
func InitCpu() {
	// init register
	reg = initRegister()
	reg = resetRegister()

	// init inst list
	initInstList()
	setInstList()
}

// Execute loaded ROM
func ExecCpu() {
	// Execute ROM
	opecode := fetchPC()
	reg.PC++

	execOpecode(opecode)
}
