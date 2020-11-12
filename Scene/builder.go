package Scene

import (
	"Raytracer/Shapes"
	"fmt"
	"math/rand"
	"os"
)

func NewBuilder(width, height, s int) (Camera, Shapes.HitTable) {
	switch s {
	case 1:
		_, _ = fmt.Fprintln(os.Stdout, "Build Scene One")
		return buildOne(width, height)
	case 2:
		_, _ = fmt.Fprintln(os.Stdout, "Build Lambertian Scene")
		return buildWorldLambertian(width, height)
	case 3:
		_, _ = fmt.Fprintln(os.Stdout, "Build Dielectrics Scene")
		return buildWorldDielectrics(width, height)
	case 4:
		_, _ = fmt.Fprintln(os.Stdout, "Build Metal Scene")
		return buildWorldMetalSpheres(width, height)
	case 5:
		_, _ = fmt.Fprintln(os.Stdout, "Final Scene")
		return buildWorldFinalScene(width, height)
	default:
		_, _ = fmt.Fprintln(os.Stdout, "Build Scene One")
		return buildOne(width, height)
	}
}

// buildWorldMetalSpheres is the end result chapter 8
func buildWorldMetalSpheres(width, height int) (Camera, Shapes.HitTableList) {
	lookFrom := Shapes.Point3{Z: 3.0}
	lookAt := Shapes.Point3{Z: -1.0}
	aperture := 0.0
	distToFocus := 1.0
	camera := NewCamera(lookFrom, lookAt, Shapes.Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)
	
	world := Shapes.HitTableList{
		Shapes.Sphere{Center: Shapes.Point3{Z: -1.0}, R: 0.5, Material: Shapes.Lambertian{Albedo: Shapes.Color{R: 0.8, G: 0.3, B: 0.3}}},
		Shapes.Sphere{Center: Shapes.Point3{Y: -100.5, Z: -1.0}, R: 100, Material: Shapes.Lambertian{Albedo: Shapes.Color{R: 0.8, G: 0.8}}},
		Shapes.Sphere{Center: Shapes.Point3{X: 1.0, Y: 0, Z: -1.0}, R: 0.5, Material: Shapes.Metal{Albedo: Shapes.Color{R: 0.8, G: 0.6, B: 0.2}, Fuzz: 1.0}},
		Shapes.Sphere{Center: Shapes.Point3{X: -1.0, Y: 0, Z: -1.0}, R: 0.5, Material: Shapes.Metal{Albedo: Shapes.Color{R: 0.8, G: 0.8, B: 0.8}, Fuzz: 0.3}},
	}
	
	return camera, world
}

func buildOne(width, height int) (Camera, Shapes.HitTableList) {
	lookFrom := Shapes.Point3{-2.0, 2.0, 5.0}
	lookAt := Shapes.Point3{Z: -1.0}
	aperture := 0.0
	distToFocus := 1.0
	camera := NewCamera(lookFrom, lookAt, Shapes.Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)
	
	world := Shapes.HitTableList{
		Shapes.Sphere{Center: Shapes.Point3{X: -0.5, Y: 0.5, Z: -1.0}, R: 0.5, Material: Shapes.Lambertian{Albedo: Shapes.Color{R: 0.1, G: 0.2, B: 0.5}}},
		Shapes.Sphere{Center: Shapes.Point3{Y: -100.5, Z: -1.0}, R: 100, Material: Shapes.Lambertian{Albedo: Shapes.Color{R: 0.5, G: 0.5, B: 0.5}}},
		/*Sphere{Center: Point3{X: -1.0, Y: 0, Z: -1.0}, R: 0.5, Material: Dielectric{1.5}},
		Sphere{Center: Point3{X: -1.0, Y: 0, Z: -1.0}, R: -0.45, Material: Dielectric{1.5}},*/
	}
	
	return camera, world
}

// buildWorldDielectrics is the end result chapter 10
func buildWorldDielectrics(width, height int) (Camera, Shapes.HitTableList) {
	lookFrom := Shapes.Point3{X: -2.0, Y: 2.0, Z: 5.0}
	lookAt := Shapes.Point3{Z: -1.0}
	aperture := 0.0
	distToFocus := 1.0
	camera := NewCamera(lookFrom, lookAt, Shapes.Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)
	
	world := Shapes.HitTableList{
		Shapes.Sphere{Center: Shapes.Point3{X: -0.5, Y: 0.5, Z: -1.0}, R: 0.5, Material: Shapes.Lambertian{Albedo: Shapes.Color{R: 0.1, G: 0.2, B: 0.5}}},
		Shapes.Sphere{Center: Shapes.Point3{Y: -100.5, Z: -1.0}, R: 100, Material: Shapes.Lambertian{Albedo: Shapes.Color{R: 0.5, G: 0.5, B: 0.5}}},
		Shapes.Sphere{Center: Shapes.Point3{X: 1.0, Y: 0, Z: -1.0}, R: 0.5, Material: Shapes.Metal{Albedo: Shapes.Color{R: 0.8, G: 0.6, B: 0.2}, Fuzz: 1.0}},
		/*Sphere{Center: Point3{X: -1.0, Y: 0, Z: -1.0}, R: 0.5, Material: Dielectric{1.5}},
		Sphere{Center: Point3{X: -1.0, Y: 0, Z: -1.0}, R: -0.45, Material: Dielectric{1.5}},*/
	}
	
	return camera, world
}

// buildWorldDielectrics is the end result book
func buildWorldFinalScene(width, height int) (Camera, Shapes.HitTableList) {
	var world Shapes.HitTableList
	
	maxSpheres := 500
	world = append(world, Shapes.Sphere{Center: Shapes.Point3{Y: -1000.0}, R: 1000, Material: Shapes.Lambertian{Albedo: Shapes.Color{R: 0.5, G: 0.5, B: 0.5}}})
	
	for a := -11; a < 11 && len(world) < maxSpheres; a++ {
		for b := -11; b < 11 && len(world) < maxSpheres; b++ {
			chooseMaterial := rand.Float64()
			Center := Shapes.Point3{X: float64(a) + 0.9*rand.Float64(), Y: 0.2, Z: float64(b) + 0.9*rand.Float64()}
			
			if Center.Sub(Shapes.Point3{X: 4.0, Y: 0.2}).Length() > 0.9 {
				switch {
				case chooseMaterial < 0.8: // diffuse
					world = append(world,
						Shapes.Sphere{
							Center:   Center,
							R:   0.2,
							Material: Shapes.Lambertian{Albedo: Shapes.Color{R: rand.Float64() * rand.Float64(), G: rand.Float64() * rand.Float64(), B: rand.Float64() * rand.Float64()}}})
				case chooseMaterial < 0.95: // metal
					world = append(world,
						Shapes.Sphere{
							Center:   Center,
							R:   0.2,
							Material: Shapes.Metal{Albedo: Shapes.Color{R: 0.5 * (1 + rand.Float64()), G: 0.5 * (1 + rand.Float64()), B: 0.5 * (1 + rand.Float64())}, Fuzz: 0.5 * rand.Float64()}})
				default:
					world = append(world,
						Shapes.Sphere{
							Center:   Center,
							R:   0.2,
							Material: Shapes.Dielectric{Ri: 1.5}})
					
				}
			}
		}
	}
	
	world = append(world,
		Shapes.Sphere{
			Center:   Shapes.Point3{Y: 1},
			R:        1.0,
			Material: Shapes.Dielectric{Ri: 1.5}},
		Shapes.Sphere{
			Center:   Shapes.Point3{X: -4, Y: 1},
			R:        1.0,
			Material: Shapes.Lambertian{Albedo: Shapes.Color{R: 0.4, G: 0.2, B: 0.1}}},
		Shapes.Sphere{
			Center:   Shapes.Point3{X: 4, Y: 1},
			R:        1.0,
			Material: Shapes.Metal{Albedo: Shapes.Color{R: 0.7, G: 0.6, B: 0.5}}})
	
	lookFrom := Shapes.Point3{X: 13, Y: 2, Z: 3}
	lookAt := Shapes.Point3{}
	aperture := 0.1
	distToFocus := 10.0
	camera := NewCamera(lookFrom, lookAt, Shapes.Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)
	
	return camera, world
}

// buildWorld is the end result chapter 7
func buildWorldLambertian(width, height int) (Camera, Shapes.HitTableList) {
	lookFrom := Shapes.Point3{Z: 3.0}
	lookAt := Shapes.Point3{Z: -1.0}
	aperture := 0.0
	distToFocus := 1.0
	camera := NewCamera(lookFrom, lookAt, Shapes.Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)
	
	world := Shapes.HitTableList{
		Shapes.Sphere{Center: Shapes.Point3{Z: -1.0}, R: 0.5, Material: Shapes.Lambertian{Albedo: Shapes.Color{R: 1.0}}},
		Shapes.Sphere{Center: Shapes.Point3{Y: -100.5, Z: -1.0}, R: 100, Material: Shapes.Lambertian{Albedo: Shapes.Color{G: 1.0}}},
	}
	
	return camera, world
}

