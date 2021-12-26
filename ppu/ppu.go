package ppu

import (
	"emu/bus"
	"fmt"
)

type PpuRegister struct {
	ppuctrl   byte // mode:W
	ppumask   byte // mode:W
	ppustatus byte // mode:R
	oamaddr   byte // mode:W
	oamdata   byte // mode:R/W
	ppuscroll byte // mode:W
	ppuaddr   byte // mode:W
	ppudata   byte // mode:R/W
}

var ppu_reg *PpuRegister

func fetchPpuRegisters() *PpuRegister {
	res := new(PpuRegister)
	res.ppustatus = bus.CPU_MEM[0x2002]
	res.oamdata = bus.CPU_MEM[0x2004]
	res.ppudata = bus.CPU_MEM[0x2007]

	return res
}

func setPpuCtrlRegister(new_num byte) {
	bus.CPU_MEM[0x2000] = new_num
}

func setPpuMaskRegister(new_num byte) {
	bus.CPU_MEM[0x2001] = new_num
}

func setOamAddrRegister(new_num byte) {
	bus.CPU_MEM[0x2003] = new_num
}

func setOamDataRegister(new_num byte) {
	bus.CPU_MEM[0x2004] = new_num
}

func setPpuScrollRegister(new_num byte) {
	bus.CPU_MEM[0x2005] = new_num
}

func setPpuAddrRegister(new_num byte) {
	bus.CPU_MEM[0x2006] = new_num
}

func setPpuDataRegister(new_num byte) {
	bus.CPU_MEM[0x2007] = new_num
}

func PpuTest() {
	ppu_reg = fetchPpuRegisters()

	fmt.Println("PPU TEST!!")
	fmt.Printf("ppuctrl:%x\n", ppu_reg.ppustatus)
	fmt.Printf("ppuctrl:%x\n", ppu_reg.oamdata)
	fmt.Printf("ppuctrl:%x\n", ppu_reg.ppudata)
}
