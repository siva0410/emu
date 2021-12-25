package cpu

import (
	"fmt"
)

var tmp struct {
	A  byte
	X  byte
	Y  byte
	SP uint16
	PC uint16
	P  byte
}

func Cpu() {
	fmt.Print(tmp)
}
