package Shapes

import (
	"math"
)

type Rnd interface {
	Float64() float64
}

// Vec3 defines a vector in 3D space
type Vec3 struct {
	X, Y, Z float64
}

// Scale scales the vector by the value (return a new vector)
func (v Vec3) Scale(t float64) Vec3 {
	return Vec3{X: v.X * t, Y: v.Y * t, Z: v.Z * t}
}

// Mult multiplies the vector by the other one (return a new vector)
func (v Vec3) Mult(v2 Vec3) Vec3 {
	return Vec3{X: v.X * v2.X, Y: v.Y * v2.Y, Z: v.Z * v2.Z}
}

// Sub substracts the 2 vectors (return a new vector)
func (v Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{X: v.X - v2.X, Y: v.Y - v2.Y, Z: v.Z - v2.Z}
}

// Add adds the 2 vectors (return a new vector)
func (v Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{X: v.X + v2.X, Y: v.Y + v2.Y, Z: v.Z + v2.Z}
}

// Length returns the size of the vector
func (v Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec3) squared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Unit returns Normalized vectors
func (v Vec3) Unit() Vec3 {
	return v.Scale(1.0 / v.Length())
}

// Negate returns a new vector with X/Y/Z negated
func (v Vec3) Negate() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

// DotProduct returns the dot product (a scalar) of 2 vectors
func DotProduct(v1 Vec3, v2 Vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// Cross returns the cross product of 2 vectors (another vector)
func Cross(v1 Vec3, v2 Vec3) Vec3 {
	return Vec3{v1.Y*v2.Z - v1.Z*v2.Y, -(v1.X*v2.Z - v1.Z*v2.X), v1.X*v2.Y - v1.Y*v2.X}
}

// Reflect simply reflects the vector based on the Normal n
func (v Vec3) Reflect(n Vec3) Vec3 {
	return v.Sub(n.Scale(2.0 * DotProduct(v, n)))
}

// Refract returns a refracted vector (or not if there is no refraction possible)
func (v Vec3) Refract(n Vec3, niOverNt float64) (bool, Vec3) {
	uv := v.Unit()
	un := n.Unit()
	
	dt := DotProduct(uv, un)
	d := 1.0 - niOverNt*niOverNt*(1-dt*dt)
	if d <= 0 {
		return false, Vec3{}
	}
	
	refracted := uv.Sub(un.Scale(dt)).Scale(niOverNt).Sub(un.Scale(math.Sqrt(d)))
	return true, refracted
	
}

// Point3  3D point
type Point3 struct {
	X, Y, Z float64
}

// Translate translates the point to a new location (return a new point)
func (p Point3) Translate(v Vec3) Point3 {
	return Point3{p.X + v.X, p.Y + v.Y, p.Z + v.Z}
}

// Sub subtracts a point to another p which gives a vector
func (p Point3) Sub(p2 Point3) Vec3 {
	return Vec3{p.X - p2.X, p.Y - p2.Y, p.Z - p2.Z}
}

// Vec3 converts a point to a vector (centered at origin)
func (p Point3) Vec3() Vec3 {
	return Vec3{p.X, p.Y, p.Z}
}

// Ray represents a ray defined by its origin and direction
type Ray struct {
	Origin Point3
	Dir    Vec3
	Rnd    Rnd
}

// PointAt returns a new point along the ray.
// P(t) = o + dt
func (r *Ray) PointAt(t float64) Point3 {
	return r.Origin.Translate(r.Dir.Scale(t))
}

// Color defines the basic Red/Green/Blue as raw float64 values
type Color struct {
	R, G, B float64
}

// Scale scales the Color by the value (return a new Color)
func (c Color) Scale(t float64) Color {
	return Color{R: c.R * t, G: c.G * t, B: c.B * t}
}

// Mult Multiplies 2 colors together (component by component multiplication)
func (c Color) Mult(c2 Color) Color {
	return Color{R: c.R * c2.R, G: c.G * c2.G, B: c.B * c2.B}
}

// Add adds the 2 colors (return a new color)
func (c Color) Add(c2 Color) Color {
	return Color{R: c.R + c2.R, G: c.G + c2.G, B: c.B + c2.B}
}

// PixelValue converts a raw Color into a pixel value (0-255) packed into a uint32
func (c Color) PixelValue() uint32 {
	r := uint32(math.Min(255.0, c.R*255.99))
	g := uint32(math.Min(255.0, c.G*255.99))
	b := uint32(math.Min(255.0, c.B*255.99))
	
	return ((r & 0xFF) << 16) | ((g & 0xFF) << 8) | (b & 0xFF)
}

// HitRecord
type HitRecord struct {
	T      float64  // which T generated the hit
	P      Point3   // which point when hit
	Normal Vec3     // Normal at that point
	Mat    Material // the material associated to this record
}

// HitTable interface of objects that can be hit by a ray
type HitTable interface {
	Hit(r *Ray, tMin float64, tMax float64) (bool, *HitRecord)
}

// HitTableList defines a simple list of hitable
type HitTableList []HitTable

// Hit returns the one closest
func (hl HitTableList) Hit(r *Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	var res *HitRecord
	hitAnything := false
	
	closestSoFar := tMax
	
	for i := range hl {
		if hit, hr := hl[i].Hit(r, tMin, closestSoFar); hit {
			hitAnything = true
			res = hr
			closestSoFar = hr.T
		}
	}
	
	return hitAnything, res
}

func RandomInUnitSphere(rnd Rnd) Vec3 {
	for {
		p := Vec3{X: 2.0*rnd.Float64() - 1.0, Y: 2.0*rnd.Float64() - 1.0, Z: 2.0*rnd.Float64() - 1.0}
		// squared of p
		p2 := p.X*p.X + p.Y*p.Y + p.Z*p.Z
		if p2 < 1.0 {
			return p
		}
	}
}

func RandomInUnitDisk(rnd Rnd) Vec3 {
	for {
		p := Vec3{X: 2.0*rnd.Float64() - 1.0, Y: 2.0*rnd.Float64() - 1.0}
		p2 := p.X*p.X + p.Y*p.Y + p.Z*p.Z
		if p2 < 1.0 {
			return p
		}
	}
}
