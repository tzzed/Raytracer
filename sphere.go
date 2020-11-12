package main

import (
	"math"
)

type Sphere struct {
	center   Point3
	radius   float64
	material Material
}

func (s Sphere) hit(r *Ray, tMin, tMax float64) (bool, *HitRecord) {
	oc := r.Origin.Sub(s.center)                // O-C
	a := DotProduct(r.Dir, r.Dir)               // d.d = d2
	b := DotProduct(oc, r.Dir)                  //  (O-C).d
	c := DotProduct(oc, oc) - s.radius*s.radius // oc.oc - r2 = od2 - r2
	d := b*b - a*c
	if d <= 0 {
		return false, nil
	}
	
	discriminantSquareRoot := math.Sqrt(d)
	temp := (-b - discriminantSquareRoot) / a
	if temp < tMax && temp > tMin {
		
		hitPoint := r.PointAt(temp)
		hr := HitRecord{
			t:      temp,
			p:      hitPoint,
			normal: hitPoint.Sub(s.center).Scale(1 / s.radius),
			mat:    s.material,
		}
		return true, &hr
	}
	
	temp = (-b + discriminantSquareRoot) / a
	if temp < tMax && temp > tMin {
		
		hitPoint := r.PointAt(temp)
		hr := HitRecord{
			t:      temp,
			p:      hitPoint,
			normal: hitPoint.Sub(s.center).Scale(1 / s.radius),
			mat:    s.material,
		}
		return true, &hr
	}
	
	return false, nil
}
