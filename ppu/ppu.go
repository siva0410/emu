package ppu

import (
	"emu/bus"
	"fmt"
)

func initPpuRegisters() *bus.PpuRegister {
	res := new(bus.PpuRegister)
	res.Ppuctrl = bus.CPU_MEM[0x2000]
	res.Ppumask = bus.CPU_MEM[0x2001]
	res.Ppustatus = bus.CPU_MEM[0x2002]
	res.Oamaddr = bus.CPU_MEM[0x2003]
	res.Oamdata = bus.CPU_MEM[0x2004]
	res.Ppuscroll = bus.CPU_MEM[0x2005]
	res.Ppuaddr = bus.CPU_MEM[0x2006]
	res.Ppudata = bus.CPU_MEM[0x2007]

	return res
}

func fetchPpuRegisters(reg *bus.PpuRegister) {
	reg.Ppuctrl = bus.CPU_MEM[0x2000]
	reg.Ppumask = bus.CPU_MEM[0x2001]
	reg.Ppustatus = bus.CPU_MEM[0x2002]
	reg.Oamaddr = bus.CPU_MEM[0x2003]
	reg.Oamdata = bus.CPU_MEM[0x2004]
	reg.Ppuscroll = bus.CPU_MEM[0x2005]
	reg.Ppuaddr = bus.CPU_MEM[0x2006]
	reg.Ppudata = bus.CPU_MEM[0x2007]
}

func copyPpuRegisters(before_reg, reg *bus.PpuRegister) {
	before_reg.Ppuctrl = reg.Ppuctrl
	before_reg.Ppumask = reg.Ppumask
	before_reg.Ppustatus = reg.Ppustatus
	before_reg.Oamaddr = reg.Oamaddr
	before_reg.Oamdata = reg.Oamdata
	before_reg.Ppuscroll = reg.Ppuscroll
	before_reg.Ppuaddr = reg.Ppuaddr
	before_reg.Ppudata = reg.Ppudata
}

func InitPpu() {
	bus.Ppu_reg = initPpuRegisters()
	bus.Before_ppu_reg = initPpuRegisters()
}

func ExecPpu() {
	fetchPpuRegisters(bus.Ppu_reg)
	if bus.Before_ppu_reg.Ppuaddr != bus.Ppu_reg.Ppuaddr {
		bus.PPU_PTR = bus.PPU_PTR << 0x8
		bus.PPU_PTR += uint32(bus.Ppu_reg.Ppuaddr)
		bus.PPU_PTR &= 0xFFFF
		fmt.Printf("MEM ppuaddr:%x\n\n", bus.PPU_PTR)
		bus.PPU_MEM[bus.PPU_PTR] = bus.Ppu_reg.Ppudata
	}

	// backup ppu_reg to before_ppu_reg
	copyPpuRegisters(bus.Before_ppu_reg, bus.Ppu_reg)

	// check ppu register
	fmt.Printf("ppuctrl:%x\t", bus.Ppu_reg.Ppuctrl)
	fmt.Printf("ppustatus:%x\t", bus.Ppu_reg.Ppustatus)
	fmt.Printf("oamaddr:%x\t\n", bus.Ppu_reg.Oamaddr)
	fmt.Printf("oamdata:%x\t", bus.Ppu_reg.Oamdata)
	fmt.Printf("ppuaddr:%x\t", bus.Ppu_reg.Ppuaddr)
	fmt.Printf("ppudata:%x\t\n", bus.Ppu_reg.Ppudata)
	fmt.Printf("MEM ppuaddr:%x\n\n", bus.PPU_PTR)

	fmt.Printf("before_ppuctrl:%x\t", bus.Before_ppu_reg.Ppuctrl)
	fmt.Printf("before_ppustatus:%x\t", bus.Before_ppu_reg.Ppustatus)
	fmt.Printf("before_oamaddr:%x\t\n", bus.Before_ppu_reg.Oamaddr)
	fmt.Printf("before_oamdata:%x\t", bus.Before_ppu_reg.Oamdata)
	fmt.Printf("before_ppuaddr:%x\t", bus.Before_ppu_reg.Ppuaddr)
	fmt.Printf("before_ppudata:%x\t\n\n", bus.Before_ppu_reg.Ppudata)

	// Palette()
}
