package window

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	square = []float32{
		// square
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
		0.5, 0.5, 0,

		// color
		1, 0, 1,
	}

	vertexcolor = []float32{
		1, 1, 0,
	}
)

type dot struct {
	drawable uint32
	colors   uint32
	w        int
	h        int
}

func draw(dots [][]*dot, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	for x := range dots {
		for _, c := range dots[x] {
			c.draw()
		}
	}
	dots[2][3].draw()
	dots[0][0].draw()
	dots[rows-1][columns-1].draw()

	glfw.PollEvents()
	window.SwapBuffers()
}

func makeDots() [][]*dot {
	dots := make([][]*dot, rows, rows)
	for h := 0; h < rows; h++ {
		for w := 0; w < columns; w++ {
			c := newDot(h, w)
			dots[h] = append(dots[h], c)
		}
	}

	return dots
}

func newDot(h, w int) *dot {
	points := make([]float32, len(square), len(square))
	copy(points, square)

	for i := 0; i < 4*3; i++ {
		var position float32
		var size float32
		switch i % 3 {
		case 0:
			size = 1.0 / float32(columns)
			position = float32(w) * size
		case 1:
			size = 1.0 / float32(rows)
			position = float32(rows-1-h) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (position * 2) - 1
		} else {
			points[i] = ((position + size) * 2) - 1
		}
	}

	return &dot{
		drawable: makeVao(points),

		h: h,
		w: w,
	}
}

func (c *dot) draw() {
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, gl.PtrOffset(3*4*4))

	return vao
}
