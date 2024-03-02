package sdf

import (
	"math"

	"hash"

	vec3 "github.com/supersdf-go/engine/vec3"
)

var infinity float32 = float32(math.MaxFloat32)

type Hashable interface {
	Hash(hash.Hash)
}

type Sdf interface {
	Distance(p vec3.Vec3) float32
}

func Optimize(s *Sdf) Sdf {
	return *s
}

func SphereIntersects(sdf Sdf, sphere *Sphere) bool {
	d0 := sdf.Distance(sphere.center)
	return d0 <= sphere.radius
}

func OptimizeIntersect(sdf *Sdf, intersect *Sdf) Sdf {
	return *sdf
}

type Sphere struct {
	center vec3.Vec3
	radius float32
}

func HashFloat32(f float32, hasher hash.Hash) {
	var buffer [4]byte
	v := math.Float32bits(f)
	buffer[0] = byte(v)
	buffer[1] = byte(v >> 8)
	buffer[2] = byte(v >> 16)
	buffer[3] = byte(v >> 24)
	hasher.Write(buffer[:])
}

func HashVec3(v vec3.Vec3, hasher hash.Hash) {
	HashFloat32(v.X, hasher)
	HashFloat32(v.Y, hasher)
	HashFloat32(v.Z, hasher)
}

func (s Sphere) Distance(p vec3.Vec3) float32 {
	return p.Subtract(s.center).Length() - s.radius
}

var sphereSalt []byte = []byte{1, 2, 3, 4}

func (s Sphere) Hash(hasher hash.Hash64) {
	HashVec3(s.center, hasher)
	HashFloat32(s.radius, hasher)
	hasher.Write(sphereSalt)
}

type Cube struct {
	center   vec3.Vec3
	halfSize vec3.Vec3
}

func (c Cube) Distance(p vec3.Vec3) float32 {
	d := p.Subtract(c.center).Abs().Subtract(c.halfSize)

	// Calculate the distance to the surface of the cube
	return float32(math.Sqrt(math.Max(0.0, float64(d.DotProduct(d)))))
}

type Union []Sdf

func (s Union) Distance(p vec3.Vec3) float32 {
	d := infinity
	for _, sdf := range s {
		dist := sdf.Distance(p)
		d = min(dist, d)
	}
	return d
}

func (s Union) Hash(h hash.Hash) {
	h.Write(sphereSalt)
	//for _, sdf := range s {
	//sdf.Hash(h)
	//}
}

type Infinity struct {
}

func (s Infinity) Distance(p vec3.Vec3) float32 {
	return infinity
}
