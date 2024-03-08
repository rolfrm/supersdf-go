// Distance field calculations
// todo: Should distance fields have hashes
// or should I calculate hashes as I go?

package sdf

import (
	"hash/fnv"
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
	Hash(hash.Hash)
}

func Optimize(s *Sdf) Sdf {
	return *s
}

func SphereIntersects(sdf Sdf, sphere *Sphere) bool {
	d0 := sdf.Distance(sphere.Center)
	return d0 <= sphere.Radius
}

func GenericIntersects(sdf Sdf, sdf2 Sdf) bool {
	return true
}

func isInfinity(sdf Sdf) bool {
	_, ok := sdf.(Infinity)
	return ok
}

func OptimizeIntersect(sdf Sdf, intersect Sdf) Sdf {

	switch obj := (sdf).(type) {
	case Sphere:
		if SphereIntersects(intersect, &obj) {
			return obj
		}
		return Infinity{}
	case Color:
		inner := OptimizeIntersect(obj.Sub, intersect)
		if isInfinity(inner) {
			return Infinity{}
		}
		if obj.Sub == inner {
			return obj
		}
		return Color{
			Color: obj.Color,
			Sub:   inner,
		}
	case Union:
		result := Union{}
		for _, v := range obj {
			sub := OptimizeIntersect(v, intersect)
			if !isInfinity(sub) {
				result = append(result, sub)
			}
		}
		if len(result) == 0 {
			return Infinity{}
		}
		if len(result) == 1 {
			return result[0]
		}
		return result
	case Cube:
		sbounds := obj.SphereBounds()
		if SphereIntersects(intersect, &sbounds) {
			if GenericIntersects(sdf, &obj) {
				return obj
			}
		}
		return Infinity{}
	}

	return sdf
}

func CompareSdfs(a Sdf, b Sdf) bool {
	h64 := fnv.New64()
	a.Hash(h64)
	asum := h64.Sum64()
	h64.Reset()
	b.Hash(h64)
	bsum := h64.Sum64()
	return asum == bsum
}

type Sphere struct {
	Center vec3.Vec3
	Radius float32
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
	return p.Subtract(s.Center).Length() - s.Radius
}

var sphereSalt []byte = []byte{1, 2, 3, 4}

func (s Sphere) Hash(hasher hash.Hash) {
	HashVec3(s.Center, hasher)
	HashFloat32(s.Radius, hasher)
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

func (c Cube) Hash(h hash.Hash) {
	HashVec3(c.center, h)
	HashVec3(c.halfSize, h)
	h.Write(sphereSalt)
}

func (c *Cube) SphereBounds() Sphere {
	return Sphere{Center: c.center, Radius: max(c.halfSize.X, c.halfSize.Y, c.halfSize.Z)}
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
	for _, sdf := range s {
		sdf.Hash(h)
	}
}

type Infinity struct {
}

func (s Infinity) Distance(p vec3.Vec3) float32 {
	return infinity
}

func (s Infinity) Hash(h hash.Hash) {
	h.Write(sphereSalt)
}

type Color struct {
	Color vec3.Vec3
	Sub   Sdf
}

func (s Color) Distance(p vec3.Vec3) float32 {
	return infinity
}

func (s Color) Hash(h hash.Hash) {
	h.Write(sphereSalt)
	HashVec3(s.Color, h)
	s.Sub.Hash(h)
}
