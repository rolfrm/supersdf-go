package engine

import (
	"math"
	"testing"
)

func TestMat4Identity(t *testing.T) {
	// Test the Mat4Identity function.
	identity := Mat4Identity()
	expected := Mat4{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}

	for i := 0; i < 16; i++ {
		if identity[i] != expected[i] {
			t.Errorf("Mat4Identity: Expected %f, got %f", expected[i], identity[i])
		}
	}
}

func TestMat4Translation(t *testing.T) {
	// Test the Mat4Translation function.
	x, y, z := float32(1.0), float32(2.0), float32(3.0)
	translation := Mat4Translation(x, y, z)

	// Check the translation matrix.
	expected := Mat4{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		x, y, z, 1.0,
	}

	for i := 0; i < 16; i++ {
		if translation[i] != expected[i] {
			t.Errorf("Mat4Translation: Expected %f, got %f", expected[i], translation[i])
		}
	}
}

func TestMat4Multiply(t *testing.T) {
	// Test the Mat4Multiply function.
	a := Mat4{
		1.0, 2.0, 3.0, 4.0,
		5.0, 6.0, 7.0, 8.0,
		9.0, 10.0, 11.0, 12.0,
		13.0, 14.0, 15.0, 16.0,
	}

	b := Mat4{
		2.0, 0.0, 0.0, 0.0,
		0.0, 2.0, 0.0, 0.0,
		0.0, 0.0, 2.0, 0.0,
		0.0, 0.0, 0.0, 2.0,
	}

	result := a.Multiply(b)

	// Check the result matrix.
	expected := Mat4{
		2.0, 4.0, 6.0, 8.0,
		10.0, 12.0, 14.0, 16.0,
		18.0, 20.0, 22.0, 24.0,
		26.0, 28.0, 30.0, 32.0,
	}

	for i := 0; i < 16; i++ {
		if result[i] != expected[i] {
			t.Errorf("Mat4Multiply: Expected %f, got %f", expected[i], result[i])
		}
	}
}

func TestMat4Rotate(t *testing.T) {
	testCases := []struct {
		angle        float32
		axis         Vec3
		testVector   Vec3
		expectVector Vec3
	}{
		{
			angle:        90,
			axis:         Vec3{0, 0, 1},
			testVector:   Vec3{1, 0, 0},
			expectVector: Vec3{0, 1, 0},
		}, {
			angle:        -90,
			axis:         Vec3{0, 0, 1},
			testVector:   Vec3{1, 0, 0},
			expectVector: Vec3{0, -1, 0},
		},
		{
			angle:        90,
			axis:         Vec3{0, 1, 0},
			testVector:   Vec3{1, 0, 0},
			expectVector: Vec3{0, 0, -1},
		},
		{
			angle:        90,
			axis:         Vec3{1, 0, 0},
			testVector:   Vec3{1, 0, 0},
			expectVector: Vec3{1, 0, 0},
		},
		{
			angle:        90,
			axis:         Vec3{1, 0, 0},
			testVector:   Vec3{0, 1, 0},
			expectVector: Vec3{0, 0, 1},
		},
		{
			angle:        45,
			axis:         Vec3{0, 1, 0},
			testVector:   Vec3{1, 0, 0},
			expectVector: Vec3{0.70710677, 0, -0.70710677},
		},
	}
	for i, c := range testCases {
		r := RotationMatrix(c.angle/360.0*2.0*math.Pi, c.axis)
		out := r.Apply(c.testVector)
		if out.Subtract(c.expectVector).Length() > 0.01 {
			t.Errorf("Unexpected rotation output for case %v", i)
		}
		//fmt.Printf(" %v  %v\n", c.expectVector, out)

	}
}
