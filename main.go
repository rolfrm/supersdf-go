package main

import (
	"math"

	. "github.com/supersdf-go/engine"
	vec3 "github.com/supersdf-go/engine/vec3"
	vec4 "github.com/supersdf-go/engine/vec4"
)

type Entity struct {
	Polygons []Polygon

	Transform Mat4
}

type Game struct {
	Entities []*Entity
	time     float32
}

func (g *Game) Draw(screen Screen) {
	proj := PerspectiveMatrix(1.2, 1.0, 0.01, 1000.0)
	//proj := Mat4Identity()

	screen.SetCamera(proj, vec3.New(0, 0, -5), vec3.New(0, 1, 0), vec3.New(float32(math.Cos(float64(g.time))), 0, float32(math.Sin(float64(g.time)))))

	for _, e := range g.Entities {
		for _, p := range e.Polygons {
			screen.Draw(p, e.Transform, p.Color)

		}
	}
}

func (g *Game) Layout(width int, height int) (int, int) {
	return width, height
}

func (g *Game) Update() {

	g.time += 0.016
	if len(g.Entities) == 0 {

		p1 := Polygon{Color: vec4.New(0.0, 0.0, 1.0, 1.0)}
		p1.Load3D([]vec3.Vec3{
			vec3.New(-1, -1, 0),
			vec3.New(1, -1, 0),
			vec3.New(-1, 1, 0),
			vec3.New(1, -1, 0),
			vec3.New(-1, 1, 0),
			vec3.New(1, 1, 0),
		})
		e1 := Entity{Polygons: []Polygon{p1}, Transform: Mat4Scale(0.5, 0.5, 0.5)}

		p2 := Polygon{Color: vec4.New(1.0, 0.0, 0.0, 1.0)}
		p2.Load3D([]vec3.Vec3{
			vec3.New(-1, -1, 0),
			vec3.New(0.0, 1, 0),
			vec3.New(1, -1, 0),
		})
		e2 := Entity{Polygons: []Polygon{p2}, Transform: Mat4Scale(0.9, 0.9, 0.9)}

		g.Entities = []*Entity{&e1, &e2}
	}
	//for _, e := range g.Entities {
	//e.Transform = RotationMatrix(g.time, vec3.New(0, 0, 1))
	//}

}

func main() {
	game := Game{}
	RunApp(&game)

}
