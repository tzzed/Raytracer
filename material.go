package main

import (
	"math"
)

/***********************
 * Material
 ************************/
// Material defines how a material scatter light
type Material interface {
	scatter(r *Ray, rec *HitRecord) (wasScattered bool, attenuation *Color, scattered *Ray)
}

/***********************
 * Lambertian material (diffuse only)
 ************************/
type Lambertian struct {
	albedo Color
}

func (mat Lambertian) scatter(r *Ray, rec *HitRecord) (bool, *Color, *Ray) {
	target := rec.p.Translate(rec.normal).Translate(randomInUnitSphere(r.rnd))
	scattered := &Ray{rec.p, target.Sub(rec.p), r.rnd}
	return true, &mat.albedo, scattered
	
}

/***********************
 * Metal material
 ************************/
type Metal struct {
	albedo Color
	fuzz   float64
}

func NewMetal(a Color, f float64) Metal {
	return Metal{albedo: a, fuzz: f}
}

func (mat Metal) scatter(r *Ray, rec *HitRecord) (bool, *Color, *Ray) {
	reflected := r.Dir.Unit().Reflect(rec.normal)
	if mat.fuzz < 1 {
		reflected = reflected.Add(randomInUnitSphere(r.rnd).Scale(mat.fuzz))
	}
	
	scattered := &Ray{rec.p, reflected, r.rnd}
	if DotProduct(scattered.Dir, rec.normal) < 0 {
		return false, nil, nil
	}
	
	return true, &mat.albedo, scattered
}

/***********************
 * Dielectric material (glass)
 ************************/
type Dielectric struct {
	refIdx float64
}

func NewDielectric(ref float64) Dielectric {
	return Dielectric{refIdx: ref}
}

func schlick(cosine float64, iRefIdx float64) float64 {
	r0 := (1.0 - iRefIdx) / (1.0 + iRefIdx)
	r0 = r0 * r0
	return r0 + (1.0-r0)*math.Pow(1.0-cosine, 5)
}

func (die Dielectric) scatter(r *Ray, rec *HitRecord) (bool, *Color, *Ray) {
	var (
		outwardNormal Vec3
		niOverNt      float64
		cosine        float64
	)
	
	DotRayNormal := DotProduct(r.Dir, rec.normal)
	if DotRayNormal > 0 {
		outwardNormal = rec.normal.Negate()
		niOverNt = die.refIdx
		cosine = DotRayNormal / r.Dir.Length()
		cosine = math.Sqrt(1.0 - die.refIdx*die.refIdx*(1.0-cosine*cosine))
	} else {
		outwardNormal = rec.normal
		niOverNt = 1.0 / die.refIdx
		cosine = -DotRayNormal / r.Dir.Length()
	}
	
	wasRefracted, refracted := r.Dir.Refract(outwardNormal, niOverNt)
	
	// refract only with some probability
	if !wasRefracted || r.rnd.Float64() < schlick(cosine, die.refIdx) {
		return true,  &Color{1.0, 1.0, 1.0}, &Ray{Origin: rec.p, Dir: r.Dir.Unit().Reflect(rec.normal), rnd: r.rnd}
	}
	
	return true,  &Color{1.0, 1.0, 1.0}, &Ray{rec.p, refracted, r.rnd}
	
}
