package sdf

import (
	"fmt"
	"testing"

	"github.com/supersdf-go/engine/vec3"
)

func abs(a float32) float32 {
	if a < 0.0 {
		return -a
	}
	return a
}
func TestUnion(t *testing.T) {
	sdf := Union{}

	d := sdf.Distance(vec3.New(1, 1, 1))
	if d != infinity {
		t.Error("Expected infinity, got:", d)
	}

	sdf = Union{
		Sphere{
			Center: vec3.New(0, 0, 0),
			Radius: 1.0,
		}, Sphere{
			Center: vec3.New(1, 0, 0),
			Radius: 1.0,
		},
		Infinity{},
	}

	d1 := sdf.Distance(vec3.New(0, 1.1, 0))
	if abs(d1-0.1) > 0.01 {
		t.Error("Expected 0, got:", d1-0.1)
	}
	d2 := sdf.Distance(vec3.New(1, 1.1, 0))
	if abs(d2-0.1) > 0.01 {
		t.Error("Expected 0, got:", d2-0.1)
	}

	d3 := sdf.Distance(vec3.New(0, 0, 0))
	if abs(d3 - -1) > 0.01 {
		t.Error("Expected 0, got:", d3-1)
	}

	fmt.Printf("%v %v\n", d1, d2)
}

func TestIntersect(t *testing.T) {
	sdf := Union{
		Sphere{
			Center: vec3.New(0, 0, 0),
			Radius: 1.0,
		}, Sphere{
			Center: vec3.New(1, 0, 0),
			Radius: 1.0,
		},
		Infinity{},
	}
	i1 := SphereIntersects(sdf, &Sphere{Center: vec3.New(-1, 0, 0), Radius: 1.1})
	if !i1 {
		t.Error("Expected intersection")
	}
	i2 := SphereIntersects(sdf, &Sphere{Center: vec3.New(1, 0, 0), Radius: 1.1})
	if !i2 {
		t.Error("Expected intersection")
	}
	i3 := SphereIntersects(sdf, &Sphere{Center: vec3.New(3.2, 0, 0), Radius: 1.1})
	if i3 {
		t.Error("Did not expect intersection")
	}
}

func TestOptimizeIntersect(t *testing.T) {
	sdf := Union{
		Sphere{
			Center: vec3.New(0, 0, 0),
			Radius: 1.0,
		}, Sphere{
			Center: vec3.New(1, 0, 0),
			Radius: 1.0,
		}, Cube{
			center:   vec3.New(0.5, 0.5, 0.5),
			halfSize: vec3.New(0.1, 0.2, 0.3),
		},
	}

	testcases := []struct {
		intersect      Sdf
		expectedResult Sdf
	}{
		{
			intersect:      Sphere{Center: vec3.New(-5, 0, 0), Radius: 1},
			expectedResult: Infinity{},
		},
		{
			intersect: Sphere{Center: vec3.New(-1.8, 0, 0), Radius: 1},
			expectedResult: Sphere{
				Center: vec3.New(0, 0, 0),
				Radius: 1.0,
			},
		},
		{
			intersect:      Sphere{Center: vec3.New(0.0, 0, 0), Radius: 1},
			expectedResult: sdf,
		},
	}
	for _, item := range testcases {
		result := OptimizeIntersect(sdf, item.intersect)
		fmt.Printf("a: %v \n", result)
		fmt.Printf("b: %v \n", item.expectedResult)
		if !CompareSdfs(item.expectedResult, result) {
			t.Error("Objects are not equal!")
		}
	}

}
