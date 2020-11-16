package shapes

import (
	"math"
)

type Sphere struct {
	Center   Point3
	R        float64
	Material Material
}

func (s Sphere) Hit(r *Ray, tMin, tMax float64) (bool, *HitRecord) {
	oc := r.Origin.Sub(s.Center)      // O-C
	a := DotProduct(r.Dir, r.Dir)     // d.d = d2
	b := DotProduct(oc, r.Dir)        //  (O-C).d
	c := DotProduct(oc, oc) - s.R*s.R // oc.oc - r2 = od2 - r2
	d := b*b - a*c
	if d <= 0 {
		return false, nil
	}
	
	discriminantSquareRoot := math.Sqrt(d)
	temp := (-b - discriminantSquareRoot) / a
	if temp < tMax && temp > tMin {
		
		hitPoint := r.PointAt(temp)
		hr := HitRecord{
			T:      temp,
			P:      hitPoint,
			Normal: hitPoint.Sub(s.Center).Scale(1 / s.R),
			Mat:    s.Material,
		}
		return true, &hr
	}
	
	temp = (-b + discriminantSquareRoot) / a
	if temp < tMax && temp > tMin {
		
		hitPoint := r.PointAt(temp)
		hr := HitRecord{
			T:      temp,
			P:      hitPoint,
			Normal: hitPoint.Sub(s.Center).Scale(1 / s.R),
			Mat:    s.Material,
		}
		return true, &hr
	}
	
	return false, nil
}


type MovingSphere struct {
	Center0, Center1   Point3
	R        float64
	tm0, Tm1 float64
	Material Material
}

func (ms MovingSphere) center(t float64)  Point3{
	tm := (t-ms.tm0) / ms.Tm1
	return ms.Center0.Translate(ms.Center1.Sub(ms.Center0).Scale(tm))
}

func (ms MovingSphere) Hit(r *Ray, tMin, tMax float64) (bool, *HitRecord) {
	center := ms.center(r.Rnd.Float64())
	oc := r.Origin.Sub(center)      // O-C
	a := DotProduct(r.Dir, r.Dir)     // d.d = d2
	b := DotProduct(oc, r.Dir)        //  (O-C).d
	c := DotProduct(oc, oc) - ms.R*ms.R // oc.oc - r2 = od2 - r2
	d := b*b - a*c
	if d <= 0 {
		return false, nil
	}
	
	discriminantSquareRoot := math.Sqrt(d)
	temp := (-b - discriminantSquareRoot) / a
	if temp < tMax && temp > tMin {
		
		hitPoint := r.PointAt(temp)
		hr := HitRecord{
			T:      temp,
			P:      hitPoint,
			Normal: hitPoint.Sub(center).Scale(1 / ms.R),
			Mat:    ms.Material,
		}
		return true, &hr
	}
	
	temp = (-b + discriminantSquareRoot) / a
	if temp < tMax && temp > tMin {
		
		hitPoint := r.PointAt(temp)
		hr := HitRecord{
			T:      temp,
			P:      hitPoint,
			Normal: hitPoint.Sub(center).Scale(1 / ms.R),
			Mat:    ms.Material,
		}
		return true, &hr
	}
	
	return false, nil
}
