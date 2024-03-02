package main

import (
	"fmt"

	. "github.com/supersdf-go/engine"
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
	for _, e := range g.Entities {
		for _, p := range e.Polygons {
			screen.Draw(p, e.Transform, vec4.New(1.0, 0.5, 0.5, 1.0))

		}
	}
}

func (g *Game) Layout(width int, height int) (int, int) {
	return width, height
}

func (g *Game) Update() {
	g.time += 0.016
	if len(g.Entities) == 0 {
		fmt.Println("New")
		p1 := Polygon{}
		p1.Load3D([]Vec3{
			NewVec3(-0.5, -0.5, 0.0),
			NewVec3(0.5, -0.5, 0.0),
			NewVec3(0.0, 0.5, 0.0),
		})
		e1 := Entity{Polygons: []Polygon{p1}, Transform: Mat4Translation(0.25, 0.0, 0.0)}

		g.Entities = []*Entity{&e1}
	}
	for _, e := range g.Entities {
		e.Transform = RotationMatrix(g.time, NewVec3(0, 0, 1))
	}

}

func main() {
	game := Game{}
	RunApp(&game)

}
