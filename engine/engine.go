package engine

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/supersdf-go/engine/vec3"
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

	// Specify vertex attribute pointers
	positionAttrib := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("vp\x00")))
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 3, gl.FLOAT, false, 0, nil)

	screen := Screen{cameraTransform: Mat4Identity()}
	screen.modelViewLoc = gl.GetUniformLocation(shaderProgram, gl.Str("modelView\x00"))
	screen.colorLoc = gl.GetUniformLocation(shaderProgram, gl.Str("color\x00"))
	screen.cameraPositionLoc = gl.GetUniformLocation(shaderProgram, gl.Str("cameraPosition\x00"))
	screen.modelTransformLoc = gl.GetUniformLocation(shaderProgram, gl.Str("modelTransform\x00"))
	screen.positionAttrib = positionAttrib
	for !window.ShouldClose() {
		ctx.Update()
		w, h := window.GetSize()
		ctx.Layout(w, h)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(shaderProgram)
		gl.Uniform3f(screen.cameraPositionLoc, screen.cameraPosition.X, screen.cameraPosition.Y, screen.cameraPosition.Z)
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
	colorLoc          int32
	positionAttrib    uint32
	modelViewLoc      int32
	cameraPositionLoc int32
	modelTransformLoc int32
	cameraPosition    vec3.Vec3
	cameraTransform   Mat4
}

func (s *Screen) SetCamera(viewTransform Mat4, cameraPosition Vec3, cameraUp Vec3, cameraRight Vec3) {
	s.cameraPosition = cameraPosition
	var camRotation = RotationMatrix2(cameraUp, cameraRight)
	var camTranslation = Mat4Translation(-cameraPosition.X, -cameraPosition.Y, -cameraPosition.Z)
	s.cameraTransform = viewTransform.Multiply(camRotation.Multiply(camTranslation))

}

func (s *Screen) Draw(polygon Polygon, modelTransform Mat4, color vec4.Vec4) {
	modelView := s.cameraTransform.Multiply(modelTransform)
	gl.UniformMatrix4fv(s.modelViewLoc, 1, false, &modelView[0])
	gl.Uniform4f(s.colorLoc, color.X, color.Y, color.Z, color.W)
	gl.BindVertexArray(polygon.buffer)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(polygon.count))
	gl.BindVertexArray(0)
}

var (
	vertexShaderSource = `
		#version 410
		uniform mat4 modelView;
		uniform mat4 modelTransform;
		uniform vec3 cameraPosition;
		in vec3 vp;
		out vec3 wp;
		out vec3 eye_dir;
		void main() {
			gl_Position = modelView * vec4(vp, 1.0);
			wp = (modelTransform * vec4(vp, 1.0)).xyz;
			eye_dir = normalize(wp - eye_dir);
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
	vao    uint32
	buffer uint32
	count  uint32
}

func (p *Polygon) Load3D(vertices []Vec3) {

	vbo := p.buffer
	if vbo == 0 {
		var vao uint32
		gl.GenVertexArrays(1, &vao)

		gl.GenBuffers(1, &vbo)
		p.buffer = vbo
		p.vao = vao
		gl.BindVertexArray(p.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, p.buffer)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
		gl.EnableVertexAttribArray(0)
	}
	gl.BindVertexArray(p.vao)
	p.count = uint32(len(vertices))
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	arr := make([]float32, len(vertices)*3)
	for i, v := range vertices {
		arr[i*3] = v.X
		arr[i*3+1] = v.Y
		arr[i*3+2] = v.Z
	}
	gl.BufferData(gl.ARRAY_BUFFER, 3*4*len(vertices), unsafe.Pointer(&arr[0]), gl.STATIC_DRAW)
	gl.BindVertexArray(0)

}
