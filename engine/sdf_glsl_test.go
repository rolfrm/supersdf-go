package engine

import (
	"fmt"
	"testing"

	"github.com/supersdf-go/engine/sdf"
	"github.com/supersdf-go/engine/vec3"
)

func TestSdf2Glsl(t *testing.T) {

	sdf0 := sdf.Union{sdf.Sphere{Center: vec3.New(0, -1, 1), Radius: 1},
		sdf.Sphere{Center: vec3.New(1, 0, 0), Radius: 1},
		sdf.Color{
			Color: vec3.New(1, 0, 0),
			Sub:   sdf.Sphere{Center: vec3.New(1, 3, 0), Radius: 1.5}},
	}

	glsl := SDF2GLSL(sdf0)
	fmt.Printf("glsl: %v\n", glsl)
}
