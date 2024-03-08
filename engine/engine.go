package engine

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/supersdf-go/engine/vec2"
	"github.com/supersdf-go/engine/vec3"
	. "github.com/supersdf-go/engine/vec3"
	"github.com/supersdf-go/engine/vec4"
)

type MainContext interface {
	Update()
	Draw(screen *Screen)
	Layout(width, height int) (int, int)
}

type ShaderProgram struct {
	program, position, uv                            uint32
	modelView, cameraPosition, model, color, texture int32
}

func (s *ShaderProgram) Activate() {
	gl.UseProgram(s.program)
}

func NewShaderProgram(program uint32) ShaderProgram {
	fmt.Printf("program: %v\n", program)
	return ShaderProgram{
		program:        program,
		position:       uint32(gl.GetAttribLocation(program, gl.Str("vp\x00"))),
		model:          gl.GetUniformLocation(program, gl.Str("model\x00")),
		modelView:      gl.GetUniformLocation(program, gl.Str("modelView\x00")),
		cameraPosition: gl.GetUniformLocation(program, gl.Str("cameraPosition\x00")),
		color:          gl.GetUniformLocation(program, gl.Str("color\x00")),
		texture:        gl.GetUniformLocation(program, gl.Str("tex1\x00")),
		uv:             uint32(gl.GetAttribLocation(program, gl.Str("uv\x00"))),
	}
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
	window, err := glfw.CreateWindow(512, 512, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}
	if err := gl.Init(); err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	gl.Enable(gl.CULL_FACE)
	shaderProgram, e := compileShaders(vertexShaderSource, fragmentShaderSource)
	if e != nil {
		panic(e)
	}

	shaderProgram2, e := compileShaders(vertexShader2Source, fragmentShader2Source)
	if e != nil {
		panic(e)
	}

	s1 := NewShaderProgram(shaderProgram)

	screen := Screen{cameraTransform: Mat4Identity()}
	screen.s.program = 100000
	screen.s1 = s1
	screen.s2 = NewShaderProgram(shaderProgram2)
	screen.UseProgram(s1)
	for !window.ShouldClose() {
		ctx.Update()
		w, h := window.GetSize()
		ctx.Layout(w, h)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		ctx.Draw(&screen)

		window.SwapBuffers()
		glfw.PollEvents()
	}
	return nil
}

func compileShaders(vertexSource string, fragmentSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentSource, gl.FRAGMENT_SHADER)
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
	cameraPosition  vec3.Vec3
	cameraTransform Mat4
	s               ShaderProgram
	s1              ShaderProgram
	s2              ShaderProgram
}

func (s *Screen) SetCamera(viewTransform Mat4, cameraPosition Vec3, cameraUp Vec3, cameraRight Vec3) {
	s.cameraPosition = cameraPosition
	var camRotation = RotationMatrix2(cameraUp, cameraRight)
	var camTranslation = Mat4Translation(-cameraPosition.X, -cameraPosition.Y, -cameraPosition.Z)
	s.cameraTransform = viewTransform.Multiply(camRotation.Multiply(camTranslation))
}

func (s *Screen) UseProgram(newShader ShaderProgram) {
	if s.s.program != newShader.program {
		s.s = newShader
		gl.UseProgram(newShader.program)
		fmt.Printf("use program %v\n", newShader.program)
	}
}

func (s *Screen) Draw(polygon Polygon, modelTransform Mat4, color vec4.Vec4) {
	s.UseProgram(s.s1)
	modelView := s.cameraTransform.Multiply(modelTransform)
	fmt.Printf("s: %v %v\n", modelView, polygon.buffer2)
	gl.UniformMatrix4fv(s.s.modelView, 1, false, &modelView[0])
	gl.UniformMatrix4fv(s.s.model, 1, false, &modelTransform[0])
	gl.Uniform3f(s.s.cameraPosition, s.cameraPosition.X, s.cameraPosition.Y, s.cameraPosition.Z)

	gl.Uniform4f(s.s.color, color.X, color.Y, color.Z, color.W)
	gl.BindVertexArray(polygon.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(polygon.count))
	gl.BindVertexArray(0)
}

func (s *Screen) DrawTextured(polygon Polygon, modelTransform Mat4, texture uint32) {
	s.UseProgram(s.s2)
	modelView := s.cameraTransform.Multiply(modelTransform)
	gl.UniformMatrix4fv(s.s.modelView, 1, false, &modelView[0])
	gl.UniformMatrix4fv(s.s.model, 1, false, &modelTransform[0])
	gl.Uniform3f(s.s.cameraPosition, s.cameraPosition.X, s.cameraPosition.Y, s.cameraPosition.Z)
	gl.Uniform1i(s.s.texture, 0)

	// Activate texture unit 0
	gl.ActiveTexture(gl.TEXTURE0)

	// Bind the texture to texture unit 0
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.BindVertexArray(polygon.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(polygon.count))
	gl.BindVertexArray(0)
}

var (
	vertexShaderSource = `
		#version 410
		uniform mat4 modelView;
		uniform mat4 model;
		
		in vec3 vp;
		out vec3 wp;
		out vec3 eye_dir;
		void main() {
			gl_Position = modelView * vec4(vp, 1.0);
			wp = (model * vec4(vp, 1.0)).xyz;
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410

		uniform vec4 color;
		out vec4 frag_color;

		uniform vec3 cameraPosition;
		in vec3 wp;
		in vec3 eye_dir;

		float sph(vec3 p, vec3 c, float r){
			return length(p - c) -r;
		}

		float sdf(vec3 p){
			return min(sph(p, vec3(0,-11,0),10),min(sph(p, vec3(-0.4, 0, -3), 0.5), min(sph(p, vec3(0.4, 0, -1), 0.5), sph(p, vec3(0.5, 0, -7), 0.5))));
		}
		vec4 colorsdf(vec3 p){
			float sph1 = sph(p, vec3(-0.4, 0, -3), 0.5);
			float sph2 =  sph(p, vec3(0.4, 0, -1), 0.5); 
			float sph3 = sph(p, vec3(0.5, 0, -7), 0.5);
			float sph4 = sph(p, vec3(0,-11,0),10);
			if(sph4 < 0.2) return vec4(1,1,1,1);
			if(sph1 <= sph2 && sph1 <= sph3) return vec4(1,0,0,1);
			if(sph2 <= sph1 && sph2 <= sph3) return vec4(0,1,0,1);
			return vec4(0,0,1,1);
		}


		void main() {
			vec3 loc = wp;
			vec3 dir = normalize(wp - cameraPosition);
			for(int i =0; i <20	 ;i++){
				float d = sdf(loc);
				loc = loc + d * 1.2 * dir;
			}

			if (sdf(loc) < 0.1) {
				frag_color = colorsdf(loc);
				
			}else{
				discard;//frag_color = vec4(0.1,0.1,0.1,1);
			}
		}
	` + "\x00"

	vertexShader2Source = `
		#version 410
		uniform mat4 modelView;
		uniform mat4 model;
		
		in vec3 vp;
		in vec2 uv;
		out vec3 eye_dir;
		out vec2 uv2;
		void main() {
			uv2 = uv;
			gl_Position = modelView * vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShader2Source = `
		#version 410

		uniform vec4 color;
		out vec4 frag_color;
		uniform sampler2D tex1;
		uniform vec3 cameraPosition;
		in vec2 uv2;

		void main() {
			frag_color = texture(tex1, uv2);
		}
	` + "\x00"
)

type Polygon struct {
	Color   vec4.Vec4
	vao     uint32
	buffer  uint32
	buffer2 uint32
	count   uint32
}

func (p *Polygon) Load3D(vertices []Vec3) {
	p.Load3DUv(vertices, nil)
}
func (p *Polygon) Load3DUv(vertices []Vec3, uvs []vec2.Vec2) {

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
		if uvs != nil {
			gl.GenBuffers(1, &p.buffer2)
			gl.BindBuffer(gl.ARRAY_BUFFER, p.buffer2)
			gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
			gl.EnableVertexAttribArray(1)

		}

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
	if uvs != nil {
		gl.BindBuffer(gl.ARRAY_BUFFER, p.buffer2)

		arr := make([]float32, len(uvs)*2)
		for i, v := range uvs {
			arr[i*2] = v.X
			arr[i*2+1] = v.Y
		}
		gl.BufferData(gl.ARRAY_BUFFER, 2*4*len(uvs), unsafe.Pointer(&arr[0]), gl.STATIC_DRAW)

	}
	gl.BindVertexArray(0)

}
