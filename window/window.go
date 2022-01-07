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
	for !window.ShouldClose() {
		// draw(dots, window, program)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		// ----------------------------------------------
		// Exec CPU and PPU
		// PPU clock = 3*CPU clock
		fmt.Printf("#%d: cycle: %d\n", i, *cycle)
		cpu.ExecCpu(cycle)
		for j := 0; j < 3; j++ {
			ppu.ExecPpu(cycle)
		}

		// for x := range dots {
		// 	for _, c := range dots[x] {
		// 		c.draw()
		// 	}
		// }
		dots[2][3].setColor([]byte{0x80, byte(i) & 0xFF, 0x80})
		dots[2][3].draw()
		dots[0][0].setColor([]byte{0x80, 0x80, 0x00})
		dots[0][0].draw()
		dots[rows-1][columns-1].draw()

		glfw.PollEvents()
		window.SwapBuffers()
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
