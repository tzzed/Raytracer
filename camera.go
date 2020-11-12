package main

import (
	"math"
)

type Camera struct {
	origin     Point3
	llc        Point3
	horizontal Vec3
	vertical   Vec3
	u, v       Vec3
	lensRadius float64
}

// NewCamera computes the parameters necessary for the camera...
//	vfov is expressed in degrees (not radians)
func NewCamera(lookFrom Point3, lookAt Point3, vup Vec3, vfov float64, aspect float64, aperture float64, focusDist float64) Camera {
	theta := vfov * (math.Pi / 180.0)
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight
	
	origin := lookFrom
	w := lookFrom.Sub(lookAt).Unit()
	u := Cross(vup, w).Unit()
	v := Cross(w, u)
	
	lowerLeftCorner := origin.Translate(u.Scale(-(halfWidth * focusDist))).Translate(v.Scale(-(halfHeight * focusDist))).Translate(w.Scale(-focusDist))
	horizontal := u.Scale(2 * halfWidth * focusDist)
	vertical := v.Scale(2 * halfHeight * focusDist)
	
	return Camera{origin, lowerLeftCorner, horizontal, vertical, u, v, aperture / 2.0}
}

func (c Camera) ray(rnd Rnd, u, v float64) *Ray {
	d := c.llc.Translate(c.horizontal.Scale(u)).Translate(c.vertical.Scale(v)).Sub(c.origin)
	
	if c.lensRadius <= 0 {
		return &Ray{c.origin, d, rnd}
	}
	
	rd := randomInUnitDisk(rnd).Scale(c.lensRadius)
	offset := c.u.Scale(rd.X).Add(c.v.Scale(rd.Y))
	
	return &Ray{c.origin.Translate(offset), d.Sub(offset), rnd}
}
