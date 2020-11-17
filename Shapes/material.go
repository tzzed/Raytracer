package shapes

import (
	"math"
)

// Material defines how a Material scatter light
type Material interface {
	Scatter(r *Ray, rec *HitRecord) (wasScattered bool, attenuation *Color, scattered *Ray)
}

// Lambertian Material (diffuse only)
type Lambertian struct {
	Albedo Texture
}

func (l Lambertian) Scatter(r *Ray, rec *HitRecord) (bool, *Color, *Ray) {
	target := rec.P.Translate(rec.Normal).Translate(RandomInUnitSphere(r.Rnd))
	scattered := &Ray{Origin: rec.P, Dir: target.Sub(rec.P), Rnd: r.Rnd}
	sc := l.Albedo.Value(rec.U, rec.V, rec.P)
	return true, &sc, scattered
	
}

// Metal Material
type Metal struct {
	Albedo Color
	Fuzz   float64
}

func (m Metal) Scatter(r *Ray, rec *HitRecord) (bool, *Color, *Ray) {
	reflected := r.Dir.Unit().Reflect(rec.Normal)
	if m.Fuzz < 1 {
		reflected = reflected.Add(RandomInUnitSphere(r.Rnd).Scale(m.Fuzz))
	}
	
	scattered := &Ray{Origin: rec.P, Dir: reflected, Rnd: r.Rnd}
	if DotProduct(scattered.Dir, rec.Normal) < 0 {
		return false, nil, nil
	}
	
	return true, &m.Albedo, scattered
}

// Dielectric Material
type Dielectric struct {
	Ri float64
}

func schlick(cosine float64, iRefIdx float64) float64 {
	r0 := (1.0 - iRefIdx) / (1.0 + iRefIdx)
	r0 = r0 * r0
	return r0 + (1.0-r0)*math.Pow(1.0-cosine, 5)
}

func (d Dielectric) Scatter(r *Ray, rec *HitRecord) (bool, *Color, *Ray) {
	var (
		outwardNormal Vec3
		niOverNt      float64
		cosine        float64
	)
	
	DotRayNormal := DotProduct(r.Dir, rec.Normal)
	if DotRayNormal > 0 {
		outwardNormal = rec.Normal.Negate()
		niOverNt = d.Ri
		cosine = DotRayNormal / r.Dir.Length()
		cosine = math.Sqrt(1.0 - d.Ri*d.Ri*(1.0-cosine*cosine))
	} else {
		outwardNormal = rec.Normal
		niOverNt = 1.0 / d.Ri
		cosine = -DotRayNormal / r.Dir.Length()
	}
	
	wasRefracted, refracted := r.Dir.Refract(outwardNormal, niOverNt)
	// refract only with some probability
	if !wasRefracted || r.Rnd.Float64() < schlick(cosine, d.Ri) {
		return true, &Color{R: 1.0, G: 1.0, B: 1.0}, &Ray{Origin: rec.P, Dir: r.Dir.Unit().Reflect(rec.Normal), Rnd: r.Rnd}
	}
	
	return true, &Color{R: 1.0, G: 1.0, B: 1.0}, &Ray{Origin: rec.P, Dir: refracted, Rnd: r.Rnd}
	
}
