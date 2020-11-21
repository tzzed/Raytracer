package shapes

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

type Texture interface {
	Value(u, v float64, p Point3) Color
}

type SolidColor struct {
	ColorValue Color
}

func NewSolidColor(c Color) SolidColor {
	return SolidColor{
		ColorValue: c,
	}
}

func (sc SolidColor) Value(u, v float64, p Point3) Color {
	return sc.ColorValue
}

type CheckerTexture struct {
	Odd  Texture
	Even Texture
}

func (c CheckerTexture) Value(u, v float64, p Point3) Color {
	sines := math.Sin(10*p.X) * math.Sin(10*p.Y) * math.Sin(10*p.Z)
	if sines < 0 {
		return c.Odd.Value(u, v, p)
	}

	return c.Even.Value(u, v, p)
}

type Perlin struct {
	X, Y, Z []int
	Rnd     []float64
	count   int
}

func NewPerlin() Perlin {
	p := Perlin{count: 256}
	r := make([]float64, p.count)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < p.count; i++ {
		rnd := rand.New(rand.NewSource(rand.Int63()))
		r[i] = rnd.Float64()
	}

	p.X = p.genPerm()
	p.Y = p.genPerm()
	p.Z = p.genPerm()
	p.Rnd = r
	return p
}

func (pe Perlin) noise(p Point3) float64 {
	u := p.X - math.Floor(p.X)
	u = u * u * (3 - (2 * u))

	v := p.Y - math.Floor(p.Y)
	v = v * v * (3 - (2 * v))

	w := p.Z - math.Floor(p.Z)
	w = w * w * (3 - (2 * w))

	i := int(math.Floor(p.X))
	j := int(math.Floor(p.Y))
	k := int(math.Floor(p.Z))

	var c [2][2][2]float64

	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				c[di][dj][dk] = pe.Rnd[pe.X[(i+di)&0xFF]^pe.Y[(j+dj)&0xFF]^pe.Z[(k+dk)&0xFF]]
			}
		}
	}

	return pe.trailingInterpret(c, u, v, w)
}

func (pe Perlin) trailingInterpret(c [2][2][2]float64, u, v, w float64) float64 {
	accu := 0.0
	for i := 0.0; i < 2; i++ {
		for j := 0.0; j < 2; j++ {
			for k := 0.0; k < 2; k++ {
				accu += (i*u + 1 - i) * (1 - u) * (j*v + (1-j)*(1-v)) *
					(k*w + (1-k)*(1-w)) * c[int(i)][int(j)][int(k)]
			}
		}
	}

	return accu
}

func permute(point []int, n int) []int {
	rand.Seed(time.Now().UnixNano())
	for i := n - 1; i > 0; i-- {
		target := int(rand.Int63n(int64(i)))
		_, _ = fmt.Fprintf(os.Stderr, "target: %d\n", target)
		point[i], point[target] = point[target], point[i]
	}

	return point

}

func (pe Perlin) genPerm() []int {
	p := make([]int, pe.count)

	for i := 0; i < pe.count; i++ {
		p[i] = i
	}

	return permute(p, pe.count)
}

type NoiseTexture struct {
	Perlin Perlin
	scale  float64
}

// NewNoiseTexture returns a new NoiseTexture with the coefficient scale.
func NewNoiseTexture(s float64) NoiseTexture {
	p := NewPerlin()
	return NoiseTexture{Perlin: p, scale: s}
}

func (n NoiseTexture) Value(u, v float64, p Point3) Color {
	t := n.Perlin.noise(p)
	black := Color{1, 1, 1}.Scale(t)
	return black
}
