package vec4

import "math"

type Vec4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

func New(x, y, z, w float32) Vec4 {
	return Vec4{
		X: x, Y: y, Z: z, W: w,
	}
}
func (v1 *Vec4) Add(v2 Vec4) Vec4 {
	return Vec4{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z, v1.W + v2.W}
}

// Subtract returns the difference between two vectors.
func (v1 *Vec4) Subtract(v2 Vec4) Vec4 {
	return Vec4{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z, v1.W - v2.W}
}

// MultiplyScalar multiplies a vector by a scalar.
func (v *Vec4) MultiplyScalar(scalar float32) Vec4 {
	return Vec4{v.X * scalar, v.Y * scalar, v.Z * scalar, v.W * scalar}
}

// DotProduct calculates the dot product of two vectors.
func (v1 *Vec4) DotProduct(v2 Vec4) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z + v1.W*v2.W
}

// Length calculates the magnitude of a vector.
func (v Vec4) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Normalize normalizes the vector to have a magnitude of 1.
func (v Vec4) Normalize() Vec4 {
	magnitude := v.Length()
	if magnitude == 0 {
		return Vec4{} // Avoid division by zero
	}
	return v.MultiplyScalar(1 / magnitude)
}
