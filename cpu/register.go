package cpu

import (
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
