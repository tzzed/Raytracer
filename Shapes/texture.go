package shapes

import "math"

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

func (sc SolidColor) Value(u, v float64, p Point3) Color  {
	return sc.ColorValue
}

type CheckerTexture struct {
	Odd  Texture
	Even Texture
}

func (c CheckerTexture) Value(u, v float64, p Point3) Color  {
	sines := math.Sin(10*p.X)* math.Sin(10*p.Y) * math.Sin(10 * p.Z)
	if sines < 0 {
		return c.Odd.Value(u, v, p)
	}
	
	return c.Even.Value(u, v, p)
}