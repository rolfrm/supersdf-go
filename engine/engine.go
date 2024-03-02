package engine

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	. "github.com/supersdf-go/engine/vec3"
	"github.com/supersdf-go/engine/vec4"
)

type MainContext interface {
	Update()
	Draw(screen Screen)
	Layout(width, height int) (int, int)
}

func RunApp(ctx MainContext) error {
	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}
	if err := gl.Init(); err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	shaderProgram, e := compileShaders()
	if e != nil {
		panic(e)
	}
	gl.UseProgram(shaderProgram)

	// Create Vertex Array Object (VAO)
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Create Vertex Buffer Object (VBO) for vertices
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	// Specify vertex attribute pointers
	positionAttrib := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("vp\x00")))
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 3, gl.FLOAT, false, 0, nil)

	screen := Screen{}
	screen.transformLoc = gl.GetUniformLocation(shaderProgram, gl.Str("transform\x00"))
	screen.colorLoc = gl.GetUniformLocation(shaderProgram, gl.Str("color\x00"))
	screen.positionAttrib = positionAttrib
	for !window.ShouldClose() {
		ctx.Update()
		w, h := window.GetSize()
		ctx.Layout(w, h)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		ctx.Draw(screen)

		window.SwapBuffers()
		glfw.PollEvents()
	}
	return nil
}

func compileShaders() (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	var status int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength)
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, &log[0])

		return 0, fmt.Errorf("linking program failed: %v", string(log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return shaderProgram, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])

		return 0, fmt.Errorf("compiling shader failed: %v", string(log))
	}

	return shader, nil
}

type Screen struct {
	colorLoc       int32
	positionAttrib uint32
	transformLoc   int32
}

func (s *Screen) Draw(polygon Polygon, transform Mat4, color vec4.Vec4) {
	gl.UniformMatrix4fv(s.transformLoc, 1, false, &transform[0])
	gl.Uniform4f(s.colorLoc, color.X, color.Y, color.Z, color.W)
	gl.BindVertexArray(polygon.buffer)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(polygon.count))
}

var (
	vertices = []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}

	colors = []float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
	}
	vertexShaderSource = `
		#version 410
		uniform mat4 transform;
		in vec3 vp;
		void main() {
			gl_Position = transform * vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410

		uniform vec4 color;
		out vec4 frag_color;

		void main() {
			frag_color = color;
		}
	` + "\x00"
)

type Polygon struct {
	Color  vec4.Vec4
	buffer uint32
	count  uint32
}

func (p *Polygon) Load3D(vertices []Vec3) {
	vbo := p.buffer
	if vbo == 0 {
		gl.GenBuffers(1, &vbo)
		p.buffer = vbo
	}
	p.count = uint32(len(vertices))
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	arr := make([]float32, len(vertices)*3)
	for i, v := range vertices {
		arr[i*3] = v.X
		arr[i*3+1] = v.Y
		arr[i*3+2] = v.Z
	}
	gl.BufferData(gl.ARRAY_BUFFER, 3*4*len(vertices), unsafe.Pointer(&arr[0]), gl.STATIC_DRAW)

}
