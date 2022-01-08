package window

import (
	"emu/cpu"
	"emu/ppu"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	height = 960
	width  = 1024

	rows    = 240
	columns = 256

	vertexShaderSource = `
		#version 410 core
		layout(location = 0) in vec3 vp;
                layout(location = 1) in vec3 vc;
                out vec4 vColor;

		void main() {
			gl_Position = vec4(vp, 1.0);
                        vColor = vec4(vc, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410 core
                in vec4 vColor;
		out vec4 frag_color;
		void main() {
			frag_color = vColor;
		}
	` + "\x00"
)

func Window() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	dots := makeDots()

	var cycle *int
	cycle = new(int)

	i := 0
	line := 0
	for !window.ShouldClose() {
		// draw(dots, window, program)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		// ----------------------------------------------
		// Exec CPU and PPU
		// PPU clock = 3*CPU clock
		fmt.Printf("#cycle: %d\n", *cycle)
		cpu.ExecCpu(cycle)
		// ppu.ExecPpu(cycle)

		if *cycle >= 341 {
			*cycle -= 341
			line++
			fmt.Println(line)
		}

		if (line+1)%8 == 0 && line < 240 {
			// set sprite
			for ws := 0; ws < columns/8; ws++ {
				sprite_num := ppu.PPU_MEM[0x2000+0x20*((line+1)/8-1)+ws]
				for l := 8 * ((line+1)/8 - 1); l < line+1; l++ {
					for i := 8 * ws; i < 8*(ws+1); i++ {
						if sprite_num != 0 {
							fmt.Println("--------------------------------")
							fmt.Println(ws, l, i, sprite_num, 0x10*sprite_num, ppu.PPU_MEM[0x10*uint16(0x10*sprite_num)])
						}
						s := (ppu.PPU_MEM[0x10*int(sprite_num)+l-8*((line+1)/8-1)] >> (7 - (i - 8*ws))) & 0b1
						t := ((ppu.PPU_MEM[0x08+0x10*int(sprite_num)+l-8*((line+1)/8-1)] >> (7 - (i - 8*ws))) & 0b1) << 1
						dots[l][i].sprite = s + t
					}
				}
			}
		}

		if (line+1)%16 == 0 && line < 240 {
			// set palette

			for ps := 0; ps < columns/16; ps++ {
				for l := 16 * ((line+1)/16 - 1); l < line+1; l++ {
					for i := 16 * ps; i < 16*(ps+1); i++ {
						switch {
						case i < 8*(ps+1) && l < (line+1)/2:
							dots[l][i].palette = (ppu.PPU_MEM[0x2000+0x03C0+0x4*((line+1)/16-1)+int(ps/4)] >> 0) & 0b11
						case i >= 8*(ps+1) && l < (line+1)/2:
							dots[l][i].palette = (ppu.PPU_MEM[0x2000+0x03C0+0x4*((line+1)/16-1)+int(ps/4)] >> 2) & 0b11
						case i < 8*(ps+1) && l >= (line+1)/2:
							dots[l][i].palette = (ppu.PPU_MEM[0x2000+0x03C0+0x4*((line+1)/16-1)+int(ps/4)] >> 4) & 0b11
						case i >= 8*(ps+1) && l >= (line+1)/2:
							dots[l][i].palette = (ppu.PPU_MEM[0x2000+0x03C0+0x4*((line+1)/16-1)+int(ps/4)] >> 6) & 0b11
						}
					}
				}
			}
		}

		if line == 262 {
			line = 0
			ppu.UpdatePalette()
			for x := range dots {
				for _, dot := range dots[x] {
					dot.setColor(ppu.Palettes[dot.palette][dot.sprite][:])
					// dot.draw()
					if dot.sprite != 0 {
						// fmt.Println("--------------------------------")
						// fmt.Println(dot.sprite)
						fmt.Println(dot.sprite, dot.palette, dot.sprite)
						dot.draw()
					}
				}
			}
			// for x := range dots {
			// 	for _, c := range dots[x] {
			// 		c.draw()
			// 	}
			// }
			// dots[2][3].setColor([]byte{0x80, byte(i) & 0xFF, 0x80})
			// dots[2][3].draw()

			// dots[0][0].setColor([]byte{0x80, 0x80, 0x00})
			// dots[0][0].draw()
			// dots[rows-1][columns-1].draw()
			glfw.PollEvents()
			window.SwapBuffers()

		}

		// for x := range dots {
		// 	for _, c := range dots[x] {
		// 		fmt.Println(c.sprite, c.palette)
		// 	}
		// }
		// // for x := range dots {
		// // 	for _, c := range dots[x] {
		// // 		c.draw()
		// // 	}
		// // }
		// dots[2][3].setColor([]byte{0x80, byte(i) & 0xFF, 0x80})
		// dots[2][3].draw()

		// dots[0][0].setColor([]byte{0x80, 0x80, 0x00})
		// dots[0][0].draw()
		// dots[rows-1][columns-1].draw()

		// glfw.PollEvents()
		// window.SwapBuffers()
		i++
	}
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "NES EMU", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
