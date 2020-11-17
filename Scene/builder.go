package scene

import (
	"Raytracer/shapes"
	"fmt"
	"math/rand"
	"os"
)

func NewBuilder(width, height, s int) (Camera, shapes.HitTable) {
	switch s {
	case 1:
		_, _ = fmt.Fprintln(os.Stdout, "Build scene One")
		return buildOne(width, height)
	case 2:
		_, _ = fmt.Fprintln(os.Stdout, "Build Lambertian scene")
		return buildWorldLambertian(width, height)
	case 3:
		_, _ = fmt.Fprintln(os.Stdout, "Build Dielectrics scene")
		return buildWorldDielectrics(width, height)
	case 4:
		_, _ = fmt.Fprintln(os.Stdout, "Build Metal scene")
		return buildWorldMetalSpheres(width, height)
	case 5:
		_, _ = fmt.Fprintln(os.Stdout, "Final scene")
		return buildWorldFinalScene(width, height)
	case 6:
		_, _ = fmt.Fprintln(os.Stdout, "Reversed sphere scene")
		return buildReflectedSpheres(width, height)
		
	default:
		_, _ = fmt.Fprintln(os.Stdout, "Build scene One")
		return buildOne(width, height)
	}
}

// buildWorldMetalSpheres
func buildWorldMetalSpheres(width, height int) (Camera, shapes.HitTableList) {
	lookFrom := shapes.Point3{Z: 3.0}
	lookAt := shapes.Point3{Z: -1.0}
	aperture := 0.0
	distToFocus := 1.0
	camera := NewCamera(lookFrom, lookAt, shapes.Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)
	
	world := []shapes.HitTable{
		shapes.Sphere{Center: shapes.Point3{Z: -1.0}, R: 0.5, Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: 0.8, G: 0.3, B: 0.3})}},
		shapes.Sphere{Center: shapes.Point3{Y: -100.5, Z: -1.0}, R: 100, Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: 0.8, G: 0.8})}},
		shapes.Sphere{Center: shapes.Point3{X: 1.0, Y: 0, Z: -1.0}, R: 0.5, Material: shapes.Metal{Albedo: shapes.Color{R: 0.8, G: 0.6, B: 0.2}, Fuzz: 1.0}},
		shapes.Sphere{Center: shapes.Point3{X: -1.0, Y: 0, Z: -1.0}, R: 0.5, Material: shapes.Metal{Albedo: shapes.Color{R: 0.8, G: 0.8, B: 0.8}, Fuzz: 0.3}},
	}
	
	return camera, shapes.HitTableList{Hits: world}
}

func buildOne(width, height int) (Camera, shapes.HitTableList) {
	lookFrom := shapes.Point3{-2.0, 2.0, 5.0}
	lookAt := shapes.Point3{Z: -1.0}
	aperture := 0.0
	distToFocus := 1.0
	camera := NewCamera(lookFrom, lookAt, shapes.Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)
	c := shapes.Point3{X: -0.5, Y: 0.5, Z: -1.0}
	y := rand.Float64() - 0.5
	z := rand.Float64()
	fmt.Fprintf(os.Stderr, "y: %.2f, z: %.2f\n", y, z)
	world := []shapes.HitTable{
		shapes.Sphere{Center: shapes.Point3{X: -1, Y: -0.5, Z: -4.6}, R: 0.5, Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: 0.1, G: 0.2, B: 0.5})}},
		shapes.Sphere{Center: shapes.Point3{Y: -100.5, Z: -1.0}, R: 100, Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: 0.5, G: 0.5, B: 0.5})}},
		
		shapes.MovingSphere{
		Center0:  c,
		Center1:  c.Translate(shapes.Vec3{Y: y}.Scale(-0.5)),
		R:        0.1,
		Tm1:      1.0,
		Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: rand.Float64() * rand.Float64(), G: rand.Float64() * rand.Float64(), B: rand.Float64() * rand.Float64()})},
		},
	}
	
	return camera, shapes.HitTableList{Hits: world}
}

// buildWorldDielectrics is the end result chapter 10
func buildWorldDielectrics(width, height int) (Camera, shapes.HitTableList) {
	lookFrom := shapes.Point3{X: -2.0, Y: 2.0, Z: 5.0}
	lookAt := shapes.Point3{Z: -1.0}
	aperture := 0.0
	distToFocus := 1.0
	camera := NewCamera(lookFrom, lookAt, shapes.Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)
	
	world := []shapes.HitTable{
		shapes.Sphere{Center: shapes.Point3{X: -0.5, Y: 0.5, Z: -1.0}, R: 0.5, Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: 0.1, G: 0.2, B: 0.5})}},
		shapes.Sphere{Center: shapes.Point3{Y: -100.5, Z: -1.0}, R: 100, Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: 0.5, G: 0.5, B: 0.5})}},
		shapes.Sphere{Center: shapes.Point3{X: 1.0, Y: 0, Z: -1.0}, R: 0.5, Material: shapes.Metal{Albedo: shapes.Color{R: 0.8, G: 0.6, B: 0.2}, Fuzz: 1.0}},
		/*Sphere{Center: Point3{X: -1.0, Y: 0, Z: -1.0}, R: 0.5, Material: Dielectric{1.5}},
		Sphere{Center: Point3{X: -1.0, Y: 0, Z: -1.0}, R: -0.45, Material: Dielectric{1.5}},*/
	}
	
	return camera, shapes.HitTableList{Hits: world}
}

// buildWorldDielectrics is the end result book
func buildWorldFinalScene(width, height int) (Camera, shapes.HitTableList) {
	var world []shapes.HitTable
	
	maxSpheres := 500
	checker := shapes.CheckerTexture{Odd: shapes.NewSolidColor(shapes.Color{R: 0.2, G: 0.3, B: 0.1}), Even: shapes.NewSolidColor(shapes.Color{0.9, 0.9,0.9})}
	world = append(world, shapes.Sphere{Center: shapes.Point3{Y: -1000.0}, R: 1000, Material: shapes.Lambertian{Albedo: checker}})
	
	for a := -11; a < 11 && len(world) < maxSpheres; a++ {
		for b := -11; b < 11 && len(world) < maxSpheres; b++ {
			chooseMaterial := rand.Float64()
			Center := shapes.Point3{X: float64(a) + 0.9*rand.Float64(), Y: 0.2, Z: float64(b) + 0.9*rand.Float64()}
			if Center.Sub(shapes.Point3{X: 4.0, Y: 0.2}).Length() > 0.9 {
				switch {
				case chooseMaterial < 0.8: // diffuse
				world = append(world,
						shapes.Sphere{
							Center:   Center,
							R:        0.2,
							Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: rand.Float64() * rand.Float64(), G: rand.Float64() * rand.Float64(), B: rand.Float64() * rand.Float64()})}})
				case chooseMaterial < 0.95: // metal
					world = append(world,
						shapes.Sphere{
							Center:   Center,
							R:        0.2,
							Material: shapes.Metal{Albedo: shapes.Color{R: 0.5 * rand.Float64(), G: 0.5 * rand.Float64(), B: 0.5 * rand.Float64()}, Fuzz: 0.5 * rand.Float64()}})
				default:
					world = append(world,
						shapes.Sphere{
							Center:   Center,
							R:        0.2,
							Material: shapes.Dielectric{Ri: 1.5}})
					
				}
			}
		}
	}
	
	world = append(world,
		shapes.Sphere{
			Center:   shapes.Point3{Y: 1},
			R:        1.0,
			Material: shapes.Dielectric{Ri: 1.5}},
		shapes.Sphere{
			Center:   shapes.Point3{X: -4, Y: 1},
			R:        1.0,
			Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: 0.4, G: 0.2, B: 0.1})}},
		shapes.Sphere{
			Center:   shapes.Point3{X: 4, Y: 1},
			R:        1.0,
			Material: shapes.Metal{Albedo: shapes.Color{R: 0.7, G: 0.6, B: 0.5}}})
	
	lookFrom := shapes.Point3{X: 13, Y: 2, Z: 3}
	lookAt := shapes.Point3{}
	aperture := 0.1
	distToFocus := 10.0
	camera := NewCamera(lookFrom, lookAt, shapes.Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)
	
	return camera, shapes.HitTableList{Hits: world}
}

// buildWorld with Lambertian
func buildWorldLambertian(width, height int) (Camera, shapes.HitTableList) {
	lookFrom := shapes.Point3{Z: 3.0}
	lookAt := shapes.Point3{Z: -1.0}
	aperture := 0.0
	distToFocus := 1.0
	camera := NewCamera(lookFrom, lookAt, shapes.Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)
	
	world := []shapes.HitTable{
		shapes.Sphere{Center: shapes.Point3{Z: -1.0}, R: 0.5, Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: 1.0})}},
		shapes.Sphere{Center: shapes.Point3{Y: -100.5, Z: -1.0}, R: 100, Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{G: 1.0})}},
	}
	
	return camera, shapes.HitTableList{Hits: world}
}

// buildWorldDielectrics is the end result book
func buildMovingWorldFinalScene(width, height int) (Camera, shapes.HitTableList) {
	var world []shapes.HitTable
	
	maxSpheres := 500
	world = append(world, shapes.Sphere{Center: shapes.Point3{Y: -1000.0}, R: 1000, Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: 0.5, G: 0.5, B: 0.5})}})
	
	for a := -11; a < 11 && len(world) < maxSpheres; a++ {
		for b := -11; b < 11 && len(world) < maxSpheres; b++ {
			chooseMaterial := rand.Float64()
			Center := shapes.Point3{X: float64(a) + 0.9*rand.Float64(), Y: 0.2, Z: float64(b) + 0.9*rand.Float64()}
			
			if Center.Sub(shapes.Point3{X: 4.0, Y: 0.2}).Length() > 0.9 {
				switch {
				case chooseMaterial < 0.8: // diffuse
					/*	world = append(world,
						shapes.Sphere{
							Center:   Center,
							R:   0.2,
							Material: shapes.Lambertian{Albedo: shapes.Color{R: rand.Float64() * rand.Float64(), G: rand.Float64() * rand.Float64(), B: rand.Float64() * rand.Float64()}}})*/
					world = append(world,
						shapes.MovingSphere{
							Center0:  Center,
							Center1:  Center.Translate(shapes.Vec3{Y: rand.Float64() - 0.5}),
							R:        0.2,
							Tm1:      1.0,
							Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: rand.Float64() * rand.Float64(), G: rand.Float64() * rand.Float64(), B: rand.Float64() * rand.Float64()})}})
				
				case chooseMaterial < 0.95: // metal
					world = append(world,
						shapes.Sphere{
							Center:   Center,
							R:        0.2,
							Material: shapes.Metal{Albedo: shapes.Color{R: 0.5 * rand.Float64(), G: 0.5 * rand.Float64(), B: 0.5 * rand.Float64()}, Fuzz: 0.5 * rand.Float64()}})
				default:
					world = append(world,
						shapes.Sphere{
							Center:   Center,
							R:        0.2,
							Material: shapes.Dielectric{Ri: 1.5}})
					
				}
			}
		}
	}
	
	world = append(world,
		shapes.Sphere{
			Center:   shapes.Point3{Y: 1},
			R:        1.0,
			Material: shapes.Dielectric{Ri: 1.5}},
		shapes.Sphere{
			Center:   shapes.Point3{X: -4, Y: 1},
			R:        1.0,
			Material: shapes.Lambertian{Albedo: shapes.NewSolidColor(shapes.Color{R: 0.4, G: 0.2, B: 0.1})}},
		shapes.Sphere{
			Center:   shapes.Point3{X: 4, Y: 1},
			R:        1.0,
			Material: shapes.Metal{Albedo: shapes.Color{R: 0.7, G: 0.6, B: 0.5}}})
	
	lookFrom := shapes.Point3{X: 13, Y: 2, Z: 3}
	lookAt := shapes.Point3{}
	aperture := 0.1
	distToFocus := 10.0
	camera := NewCamera(lookFrom, lookAt, shapes.Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)
	
	return camera, shapes.HitTableList{Hits: world}
}

func buildReflectedSpheres(width, height int) (Camera, shapes.HitTableList) {
	lookFrom := shapes.Point3{X:13, Y: 2.0, Z: 3.0}
	lookAt := shapes.Point3{}
	aperture := 0.0
	distToFocus := 1.0
	camera := NewCamera(lookFrom, lookAt, shapes.Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)
	
	checker := shapes.CheckerTexture{Odd: shapes.NewSolidColor(shapes.Color{R: 0.2, G: 0.3, B: 0.1}), Even: shapes.NewSolidColor(shapes.Color{R: 0.9, G: 0.9, B: 0.9})}

	world := []shapes.HitTable{
		shapes.Sphere{Center: shapes.Point3{Y: -10}, R: 10, Material: shapes.Lambertian{Albedo: checker}},
		shapes.Sphere{Center: shapes.Point3{Y: 10}, R: 10, Material: shapes.Lambertian{Albedo: checker}},
	}
	
	return camera, shapes.HitTableList{Hits: world}
}