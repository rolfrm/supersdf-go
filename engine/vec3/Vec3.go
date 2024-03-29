package vec3

import "math"

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

func New(x, y, z float32) Vec3 {
	return Vec3{
		X: x, Y: y, Z: z,
	}
}

func Add(v1, v2 Vec3) Vec3 {
	return Vec3{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z}
}

// Subtract returns the difference between two vectors.
func (v1 Vec3) Subtract(v2 Vec3) Vec3 {
	return Vec3{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
}

// MultiplyScalar multiplies a vector by a scalar.
func (v Vec3) MultiplyScalar(scalar float32) Vec3 {
	return Vec3{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (v Vec3) Abs() Vec3 {
	return New(float32(math.Abs(float64(v.X))), float32(math.Abs(float64(v.Y))), float32(math.Abs(float64(v.Z))))
}

// DotProduct calculates the dot product of two vectors.
func (v1 Vec3) DotProduct(v2 Vec3) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// CrossProduct calculates the cross product of two vectors.
func (v1 Vec3) CrossProduct(v2 Vec3) Vec3 {
	return Vec3{
		v1.Y*v2.Z - v1.Z*v2.Y,
		v1.Z*v2.X - v1.X*v2.Z,
		v1.X*v2.Y - v1.Y*v2.X,
	}
}

// Length calculates the magnitude of a vector.
func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Normalize normalizes the vector to have a magnitude of 1.
func (v Vec3) Normalize() Vec3 {
	magnitude := v.Length()
	if magnitude == 0 {
		return Vec3{} // Avoid division by zero
	}
	return v.MultiplyScalar(1 / magnitude)
}
