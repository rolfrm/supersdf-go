package main

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	. "github.com/supersdf-go/engine"
	"github.com/supersdf-go/engine/vec2"
	vec3 "github.com/supersdf-go/engine/vec3"
	vec4 "github.com/supersdf-go/engine/vec4"
)

type Node struct {
	polygon   *Polygon
	subnodes  []Node
	transform Mat4
}

type Game struct {
	Entities []*Node
	time     float32
	fb       *Framebuffer
	square   Polygon
}

func (n *Node) Draw(screen *Screen, tform Mat4) {
	tform = tform.Multiply(n.transform)
	if n.polygon != nil {
		screen.Draw(*n.polygon, tform, n.polygon.Color)

	}
	if n.subnodes != nil {
		for _, node := range n.subnodes {
			node.Draw(screen, tform)
		}
	}
}

func (g *Game) Draw(screen *Screen) {
	proj := PerspectiveMatrix(1.2, 1.0, 1, 1000.0)
	//g.time = 0.0
	//screen.SetCamera(proj, vec3.New(0, 0, 4), vec3.New(0, 1, 0), vec3.New(float32(math.Cos(float64(g.time))), 0, float32(math.Sin(float64(g.time)))))
	screen.SetCamera(proj, vec3.New(float32(math.Sin(float64(g.time)))*4.0, float32(math.Sin(float64(g.time*0.5)))*4.0, float32(math.Cos(float64(g.time)))+7), vec3.New(0, 1, 0), vec3.New(1, 0, 0))
	g.fb.Bind()
	gl.Viewport(0, 0, 64, 64)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	for _, e := range g.Entities {
		e.Draw(screen, Mat4Identity())
	}

	g.fb.Unbind()
	screen.SetCamera(Mat4Identity(), vec3.New(0, 0, 0), vec3.New(0, 1, 0), vec3.New(1, 0, 0))
	//screen.SetCamera(proj, vec3.New(float32(math.Sin(float64(g.time)))*1.0, float32(math.Sin(float64(g.time*0.5)))*1.0, float32(math.Cos(float64(g.time)))+5), vec3.New(0, 1, 0), vec3.New(1, 0, 0))
	gl.Viewport(0, 0, 1024, 1024)
	screen.DrawTextured(g.square, Mat4Identity(), g.fb.Texture)
}

func (g *Game) Layout(width int, height int) (int, int) {
	return width, height
}

func (g *Game) Update() {

	g.time += 0.016
	if len(g.Entities) == 0 {

		p0 := Polygon{Color: vec4.New(1.0, 1.0, 1.0, 1.0)}
		p0.Load3DUv([]vec3.Vec3{
			vec3.New(-1, -1, 0),
			vec3.New(1, -1, 0),
			vec3.New(-1, 1, 0),
			vec3.New(1, -1, 0),
			vec3.New(1, 1, 0),
			vec3.New(-1, 1, 0),
		}, []vec2.Vec2{
			vec2.New(0, 0),
			vec2.New(1, 0),
			vec2.New(0, 1),
			vec2.New(1, 0),
			vec2.New(1, 1),
			vec2.New(0, 1),
		})
		g.square = p0

		p1 := Polygon{Color: vec4.New(1.0, 1.0, 1.0, 1.0)}
		points := []vec3.Vec3{
			vec3.New(-1, -1, 1),
			vec3.New(1, -1, 1),
			vec3.New(-1, 1, 1),
			vec3.New(1, -1, 1),
			vec3.New(1, 1, 1),
			vec3.New(-1, 1, 1)}
		outPoints := []vec3.Vec3{}
		for i := 0; i < 4; i++ {
			rx := RotationMatrix(math.Pi/2*float32(i), vec3.New(1, 0, 0))
			outPoints = append(outPoints, rx.ApplyN(points)...)
		}
		for i := 1; i < 4; i += 2 {
			rx := RotationMatrix(math.Pi/2*float32(i), vec3.New(0, 1, 0))
			outPoints = append(outPoints, rx.ApplyN(points)...)
		}

		p1.Load3D(outPoints)

		e1 := Node{polygon: &p1, transform: Mat4Scale(3, 3, 3)}
		//e2 := Node{polygon: &p1, transform: Mat4Translation(3, -1, 0)}

		/*p2 := Polygon{Color: vec4.New(1.0, 0.0, 0.0, 1.0)}
		p2.Load3D([]vec3.Vec3{
			vec3.New(-1, -1, 0),
			vec3.New(0.0, 1, 0),
			vec3.New(1, -1, 0),
		})
		e2 := Node{polygon: &p2, transform: Mat4Scale(0.9, 0.9, 0.9)}
		*/
		g.Entities = []*Node{&e1}
		fb, e := NewFramebuffer(64, 64)
		if e != nil {
			panic(e)
		}
		g.fb = fb
	}
}

func main() {
	game := Game{}
	RunApp(&game)

}
