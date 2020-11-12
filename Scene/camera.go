package Scene

import (
	"Raytracer/Shapes"
	"math"
)

type Camera struct {
	origin     Shapes.Point3
	llc        Shapes.Point3
	horizontal Shapes.Vec3
	vertical   Shapes.Vec3
	u, v       Shapes.Vec3
	lensRadius float64
}

// NewCamera computes the parameters necessary for the Camera...
func NewCamera(lookFrom Shapes.Point3, lookAt Shapes.Point3, vup Shapes.Vec3, vfov float64, aspect float64, aperture float64, focusDist float64) Camera {
	// vfov deg to radian.
	theta := vfov * (math.Pi / 180.0)
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight
	
	origin := lookFrom
	// Normalize lf - lat
	w := lookFrom.Sub(lookAt).Unit()
	u := Shapes.Cross(vup, w).Unit()
	v := Shapes.Cross(w, u)
	
	lowerLeftCorner := origin.Translate(u.Scale(-(halfWidth * focusDist))).Translate(v.Scale(-(halfHeight * focusDist))).Translate(w.Scale(-focusDist))
	horizontal := u.Scale(2 * halfWidth * focusDist)
	vertical := v.Scale(2 * halfHeight * focusDist)
	
	return Camera{origin, lowerLeftCorner, horizontal, vertical, u, v, aperture / 2.0}
}

func (c Camera) ray(rnd Shapes.Rnd, u, v float64) *Shapes.Ray {
	d := c.llc.Translate(c.horizontal.Scale(u)).Translate(c.vertical.Scale(v)).Sub(c.origin)
	
	if c.lensRadius <= 0 {
		return &Shapes.Ray{Origin: c.origin, Dir: d, Rnd: rnd}
	}
	
	rd := Shapes.RandomInUnitDisk(rnd).Scale(c.lensRadius)
	offset := c.u.Scale(rd.X).Add(c.v.Scale(rd.Y))
	
	return &Shapes.Ray{Origin: c.origin.Translate(offset), Dir: d.Sub(offset), Rnd: rnd}
}
