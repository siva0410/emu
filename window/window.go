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
                        // vColor = vec4(1.0, 1.0, 1.0, 1.0);
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

	line := 0
	for !window.ShouldClose() {
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
			sl := (line+1)/8 - 1
			for sw := 0; sw < 256/8; sw++ {
				sprite_num := ppu.PPU_MEM[0x2000+0x20*sl+sw]
				for l := 0; l < 8; l++ {
					for i := 0; i < 8; i++ {
						s := (ppu.PPU_MEM[0x10*int(sprite_num)+l] >> (7 - i)) & 0b1
						t := (ppu.PPU_MEM[0x08+0x10*int(sprite_num)+l] >> (7 - i)) & 0b1
						dots[sl*8+l][sw*8+i].sprite = s + t<<1
					}
				}
			}
		}

		if (line+1)%16 == 0 && line < 240 {
			// set palette
			pl := (line+1)/16 - 1
			for pw := 0; pw < 256/16; pw++ {
				for l := 0; l < 16; l++ {
					for i := 0; i < 16; i++ {
						switch {
						case i < 8 && l < 8:
							dots[pl*16+l][pw*16+i].palette = (ppu.PPU_MEM[0x2000+0x03C0+0x4*pl+int(pw/4)] >> 0) & 0b11
						case i >= 8 && l < 8:
							dots[pl*16+l][pw*16+i].palette = (ppu.PPU_MEM[0x2000+0x03C0+0x4*pl+int(pw/4)] >> 2) & 0b11
						case i < 8 && l >= 8:
							dots[pl*16+l][pw*16+i].palette = (ppu.PPU_MEM[0x2000+0x03C0+0x4*pl+int(pw/4)] >> 4) & 0b11
						case i >= 8 && l >= 8:
							dots[pl*16+l][pw*16+i].palette = (ppu.PPU_MEM[0x2000+0x03C0+0x4*pl+int(pw/4)] >> 6) & 0b11
						}
					}
				}
			}
		}

		if line == 262 {
			line = 0
			ppu.UpdatePalette()

			draw(dots)

			glfw.PollEvents()
			window.SwapBuffers()

		}
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
