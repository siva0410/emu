package cpu

import (
	"fmt"

	"emu/bus"
)

/* register */
type Register struct {
	A  byte
	X  byte
	Y  byte
	S  byte
	PC uint16
	P  byte
}

var reg *Register

// Init register
func initRegister() *Register {
	init_reg := new(Register)
	init_reg.A = 0x00
	init_reg.X = 0x00
	init_reg.Y = 0x00
	init_reg.S = 0x00
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
	reset_point := uint16(bus.CPU_MEM[lower_addr]) + (uint16(bus.CPU_MEM[upper_addr]) << 0x8)
	reset_reg := new(Register)
	reset_reg.A = 0x00
	reset_reg.X = 0x00
	reset_reg.Y = 0x00
	reset_reg.S = 0x00
	reset_reg.PC = reset_point

	return reset_reg
}

/*
   7654 3210
   ---- ----
   NVss DIZC
   |||| ||||
   |||| |||+- Carry
   |||| ||+-- Zero
   |||| |+--- Interrupt Disable
   |||| +---- Decimal
   ||++------ No CPU effect, see: the B flag
   |+-------- Overflow
   +--------- Negative
*/
func setStatus(flagname string, status bool) {
	if status {
		// status = true
		switch flagname {
		case "carry":
			reg.P |= 0b00000001
		case "zero":
			reg.P |= 0b00000010
		case "interrupt_disable":
			reg.P |= 0b00000100
		case "decimal":
			reg.P |= 0b00001000
		case "b_flag1":
			reg.P |= 0b00010000
		case "b_flag2":
			reg.P |= 0b00100000
		case "overflow":
			reg.P |= 0b01000000
		case "negative":
			reg.P |= 0b10000000
		}
	} else {
		// staus = false
		switch flagname {
		case "carry":
			reg.P &= 0b11111110
		case "zero":
			reg.P &= 0b11111101
		case "interrupt_disable":
			reg.P &= 0b11111011
		case "decimal":
			reg.P &= 0b11110111
		case "b_flag1":
			reg.P &= 0b11101111
		case "b_flag2":
			reg.P &= 0b11011111
		case "overflow":
			reg.P &= 0b10111111
		case "negative":
			reg.P &= 0b01111111
		}
	}

}

func getStatus(flagname string) bool {
	var status byte
	switch flagname {
	case "carry":
		status = reg.P >> 0 & 0b1
	case "zero":
		status = reg.P >> 1 & 0b1
	case "interrupt_disable":
		status = reg.P >> 2 & 0b1
	case "decimal":
		status = reg.P >> 3 & 0b1
	case "b_flag1":
		status = reg.P >> 4 & 0b1
	case "b_flag2":
		status = reg.P >> 5 & 0b1
	case "overflow":
		status = reg.P >> 6 & 0b1
	case "negative":
		status = reg.P >> 7 & 0b1
	}

	var res bool = false
	if status == 1 {
		res = true
	}

	return res
}

func setZeroFlag(num uint) {
	if num == 0 {
		setStatus("zero", true)
	} else {
		setStatus("zero", false)
	}
}

// func setCarryFlag(num uint) {
// 	if num>>7 == 1 {
// 		setStatus("negative", true)
// 	} else {
// 		setStatus("negative", false)
// 	}
// }

func setNegativeFlag(num uint) {
	if num>>7 == 1 {
		setStatus("negative", true)
	} else {
		setStatus("negative", false)
	}
}

// Fetch inst by PC
func fetchPC() byte {
	return bus.CPU_MEM[reg.PC]
}

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
		operand = uint16(fetchPC())
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
		// check ppu_addr register
		bus.CheckPpuPtr(operand)
		setZeroFlag(res)
		setNegativeFlag(res)
		reg.A = byte(operand)

	case "LDX":
		res = uint(operand)
		// check ppu_addr register
		bus.CheckPpuPtr(operand)
		setZeroFlag(res)
		setNegativeFlag(res)
		reg.X = byte(operand)

	case "LDY":
		res = uint(operand)
		// check ppu_addr register
		bus.CheckPpuPtr(operand)
		setZeroFlag(res)
		setNegativeFlag(res)
		reg.Y = byte(operand)

	case "STA":
		// check ppu_addr register
		bus.CheckPpuPtr(operand)
		bus.CPU_MEM[operand] = reg.A

	case "STX":
		// check ppu_addr register
		bus.CheckPpuPtr(operand)
		bus.CPU_MEM[operand] = reg.X

	case "STY":
		// check ppu_addr register
		bus.CheckPpuPtr(operand)
		bus.CPU_MEM[operand] = reg.Y

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
