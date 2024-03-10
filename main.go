package main

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	remotevm "github.com/rolfrm/remotevm"
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
	location vec3.Vec3
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

	screen.SetCamera(proj, g.location, vec3.New(0, 1, 0), vec3.New(1, 0, 0))
	g.fb.Bind()
	gl.Viewport(0, 0, 64, 64)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	for _, e := range g.Entities {
		e.Draw(screen, Mat4Identity())
	}

	g.fb.Unbind()
	screen.SetCamera(Mat4Identity(), vec3.New(0, 0, 0), vec3.New(0, 1, 0), vec3.New(1, 0, 0))

	gl.Viewport(0, 0, int32(screen.ScreenWidth), int32(screen.ScreenHeight))

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

		e1 := Node{polygon: &p1, transform: Mat4Scale(2, 2, 2)}
		e2 := Node{polygon: &p1, transform: Mat4Scale(2, 2, 2).Multiply(Mat4Translation(2, 0, 0))}

		/*p2 := Polygon{Color: vec4.New(1.0, 0.0, 0.0, 1.0)}
		p2.Load3D([]vec3.Vec3{
			vec3.New(-1, -1, 0),
			vec3.New(0.0, 1, 0),
			vec3.New(1, -1, 0),
		})
		e2 := Node{polygon: &p2, transform: Mat4Scale(0.9, 0.9, 0.9)}
		*/
		g.Entities = []*Node{&e1, &e2}
		fb, e := NewFramebuffer(64, 64)
		if e != nil {
			panic(e)
		}
		g.fb = fb
	}
}

func main() {
	game := Game{}

	commands := []remotevm.Command{
		remotevm.Command{
			Name:      "load-location",
			Arguments: []remotevm.Type{remotevm.Type_F64, remotevm.Type_F64, remotevm.Type_F64},
			Func: func(x, y, z float64) {
				game.location = vec3.New(float32(x), float32(y), float32(z))
				fmt.Printf("Loaded location: %v\n ", game.location)
			},
		},
	}

	file, err := os.Open("./save.bin")
	if err == nil {
		remotevm.EvalStream(commands, file, io.Discard)
	} else {
		panic(err.Error())
	}

	RunApp(&game)
	file, err = os.OpenFile("./save.bin", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModeAppend|os.ModePerm)
	if err != nil {
		panic(err)
	}
	codeStr := remotevm.CodeStream{Stream: file}

	codeStr.Write(remotevm.Op_Ld, float64(game.location.X),
		remotevm.Op_Ld, float64(game.location.Y),
		remotevm.Op_Ld, float64(game.location.Z),
		remotevm.Op_Call, byte(0))
	s, e := file.Stat()
	fmt.Printf("ending... %v %v", s.Size(), e)
	file.Close()

}
