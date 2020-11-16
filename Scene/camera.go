package scene

import (
	"Raytracer/shapes"
	"math"
	"math/rand"
	"time"
)

type Camera struct {
	origin     shapes.Point3
	llc        shapes.Point3
	horizontal shapes.Vec3
	vertical   shapes.Vec3
	u, v       shapes.Vec3
	lensRadius float64
	rnd        shapes.Rnd
}

// NewCamera computes the parameters necessary for the Camera...
func NewCamera(lookFrom shapes.Point3, lookAt shapes.Point3, vup shapes.Vec3, vfov float64, aspect float64, aperture float64, focusDist float64) Camera {
	// vfov deg to radian.
	theta := vfov * (math.Pi / 180.0)
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight
	
	origin := lookFrom
	// Normalize lf - lat
	w := lookFrom.Sub(lookAt).Unit()
	u := shapes.Cross(vup, w).Unit()
	v := shapes.Cross(w, u)
	
	lowerLeftCorner := origin.Translate(u.Scale(-(halfWidth * focusDist))).Translate(v.Scale(-(halfHeight * focusDist))).Translate(w.Scale(-focusDist))
	horizontal := u.Scale(2 * halfWidth * focusDist)
	vertical := v.Scale(2 * halfHeight * focusDist)
	
	return Camera{origin, lowerLeftCorner, horizontal, vertical, u, v, aperture / 2.0, rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (c Camera) ray(rnd shapes.Rnd, u, v float64) *shapes.Ray {
	d := c.llc.Translate(c.horizontal.Scale(u)).Translate(c.vertical.Scale(v)).Sub(c.origin)
	
	if c.lensRadius <= 0 {
		return &shapes.Ray{Origin: c.origin, Dir: d, Rnd: rnd}
	}
	
	rd := shapes.RandomInUnitDisk(rnd).Scale(c.lensRadius)
	offset := c.u.Scale(rd.X).Add(c.v.Scale(rd.Y))
	
	return &shapes.Ray{Origin: c.origin.Translate(offset), Dir: d.Sub(offset), Rnd: rnd}
}
