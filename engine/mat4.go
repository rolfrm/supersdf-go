package engine

import (
	"fmt"
	"math"

	"github.com/supersdf-go/engine/vec3"
	. "github.com/supersdf-go/engine/vec3"
)

// column-major mat4.
type Mat4 [16]float32

func Mat4Identity() Mat4 {
	return Mat4{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0}
}

func (a *Mat4) Get(r, c int) float32 {
	return a[r+c*4]
}

func (a *Mat4) Set(r, c int, v float32) {
	a[r+c*4] = v
}

func Mat4Translation(x, y, z float32) Mat4 {
	v := Mat4Identity()
	v.Set(0, 3, x)
	v.Set(1, 3, y)
	v.Set(2, 3, z)
	return v
}

func Mat4Scale(x, y, z float32) Mat4 {
	v := Mat4Identity()
	v.Set(0, 0, x)
	v.Set(1, 1, y)
	v.Set(2, 2, z)
	return v
}

func (a *Mat4) Multiply(b Mat4) Mat4 {
	result := Mat4{}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			var sum float32
			for k := 0; k < 4; k++ {
				sum += a.Get(i, k) * b.Get(k, j)
			}
			result.Set(i, j, sum)
		}
	}

	return result
}

func (m *Mat4) Apply(v Vec3) Vec3 {
	result := vec3.New(
		m[0]*v.X+m[1]*v.Y+m[2]*v.Z+m[3],
		m[4]*v.X+m[5]*v.Y+m[6]*v.Z+m[7],
		m[8]*v.X+m[9]*v.Y+m[10]*v.Z+m[11],
	)
	return result
}

func (m *Mat4) ApplyN(v []Vec3) []Vec3 {
	out := make([]Vec3, len(v))
	for i, v := range v {
		out[i] = m.Apply(v)
	}
	return out
}

func (m *Mat4) ToString() string {
	return fmt.Sprintf("%v %v %v %v\n%v %v %v %v\n%v %v %v %v\n%v %v %v %v\n",
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15])
}

func (m Mat4) Transpose() Mat4 {
	return Mat4{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
}

func PerspectiveMatrix(fov, aspect, near, far float32) Mat4 {
	// Convert field of view to radians
	fovRad := fov

	f := float32(1.0 / math.Tan(float64(fovRad)/2.0))
	znear := near
	zfar := far
	aspectInv := 1.0 / aspect

	return Mat4{
		f * aspectInv, 0, 0, 0,
		0, f, 0, 0,
		0, 0, -(zfar + znear) / (znear - zfar), -1,
		0, 0, (-2 * zfar * znear) / (znear - zfar), 0,
	}
}

func OrthographicMatrix(left, right, bottom, top, near, far float32) Mat4 {
	return Mat4{
		2 / (right - left), 0, 0, -(right + left) / (right - left),
		0, 2 / (top - bottom), 0, -(top + bottom) / (top - bottom),
		0, 0, -2 / (far - near), -(far + near) / (far - near),
		0, 0, 0, 1,
	}
}

// RotationMatrix generates a 4x4 rotation matrix for the specified angle and axis.
func RotationMatrix0(angle float32, axis Vec3) Mat4 {
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))
	oneMinusCos := 1.0 - cos

	xx := axis.X * axis.X
	yy := axis.Y * axis.Y
	zz := axis.Z * axis.Z
	xy := axis.X * axis.Y
	xz := axis.X * axis.Z
	yz := axis.Y * axis.Z

	return Mat4{
		cos + xx*oneMinusCos,
		axis.Y*sin + xy*oneMinusCos,
		-axis.Y*sin + xz*oneMinusCos,
		0,
		-axis.X*sin + xy*oneMinusCos,
		cos + yy*oneMinusCos,
		axis.X*sin + yz*oneMinusCos,
		0,
		axis.Y*sin + xz*oneMinusCos,
		-axis.X*sin + yz*oneMinusCos,
		cos + zz*oneMinusCos,
		0,
		0, 0, 0, 1,
	}.Transpose()
}

func RotationMatrix(angle float32, axis Vec3) Mat4 {
	rad := angle
	cosA := float32(math.Cos(float64(rad)))
	sinA := float32(math.Sin(float64(rad)))
	invCosA := 1 - cosA

	return Mat4{
		cosA + axis.X*axis.X*invCosA, axis.X*axis.Y*invCosA - axis.Z*sinA, axis.X*axis.Z*invCosA + axis.Y*sinA, 0,
		axis.Y*axis.X*invCosA + axis.Z*sinA, cosA + axis.Y*axis.Y*invCosA, axis.Y*axis.Z*invCosA - axis.X*sinA, 0,
		axis.Z*axis.X*invCosA - axis.Y*sinA, axis.Z*axis.Y*invCosA + axis.X*sinA, cosA + axis.Z*axis.Z*invCosA, 0,
		0, 0, 0, 1,
	}
}

func RotationMatrix2(up, right Vec3) Mat4 {

	// Calculate the forward vector using the cross product of up and right
	up = up.Normalize()
	right = right.Normalize()
	forward := right.CrossProduct(up).Normalize()

	// Assign the vectors to the rotation matrix
	matrix := Mat4{
		right.X, right.Y, right.Z, 0,
		up.X, up.Y, up.Z, 0,
		forward.X, forward.Y, forward.Z, 0,
		0, 0, 0, 1,
	}.Transpose() // column major

	return matrix
}
