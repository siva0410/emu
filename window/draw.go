package window

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

var (
	vertexPosition = []float32{
		// square
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
		0.5, 0.5, 0,
	}

	vertexInitColor = []float32{
		// color
		1, 1, 1,
	}
)

type dot struct {
	points      []float32
	colorpoints []float32
	w           int
	h           int
}

func setColor() {}

func (c *dot) setColor(colors []byte) {
	red := float32(colors[0]) / 0xFF
	green := float32(colors[1]) / 0xFF
	blue := float32(colors[2]) / 0xFF
	res := make([]float32, 0)
	res = append(res, float32(red), float32(green), float32(blue))
	c.colorpoints = res
}

// func draw(dots [][]*dot, window *glfw.Window, program uint32) {
// 	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
// 	gl.UseProgram(program)

// 	// for x := range dots {
// 	// 	for _, c := range dots[x] {
// 	// 		c.draw()
// 	// 	}
// 	// }
// 	dots[2][3].setColor([]byte{0x80, 0x00, 0x80})
// 	dots[2][3].draw()
// 	dots[0][0].setColor([]byte{0x80, 0x80, 0x00})
// 	dots[0][0].draw()
// 	dots[rows-1][columns-1].draw()

// 	glfw.PollEvents()
// 	window.SwapBuffers()
// }

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
	points := make([]float32, len(vertexPosition), len(vertexPosition))
	copy(points, vertexPosition)

	for i := 0; i < len(points); i++ {
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

	colorpoints := make([]float32, 0)
	for i := 0; i < len(points)/3; i++ {
		colorpoints = append(colorpoints, vertexInitColor...)
	}

	return &dot{
		// drawable:    makeVao(append(points, colorpoints...)),
		points:      points,
		colorpoints: colorpoints,
		h:           h,
		w:           w,
	}
}

func (c *dot) draw() {
	colorpoints := make([]float32, 0)
	for i := 0; i < len(c.points)/3; i++ {
		colorpoints = append(colorpoints, c.colorpoints...)
	}
	drawable := makeVao(append(c.points, colorpoints...))
	gl.BindVertexArray(drawable)
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
