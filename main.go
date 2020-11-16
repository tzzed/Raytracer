package main

import (
	"Raytracer/scene"
	"flag"
	"fmt"
	"image"
	clr "image/color"
	"image/png"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	
	"github.com/veandco/go-sdl2/sdl"
)

// RaysPerPixelList is used on the command line (flag) to define the number of rays per pixel per phase (hence a list)
// Example: ray-tracing -r 1 -r 10
type RaysPerPixelList []int

func (r *RaysPerPixelList) String() string {
	return fmt.Sprint(*r)
}

func (r *RaysPerPixelList) Set(value string) error {
	for _, e := range strings.Split(value, ",") {
		i, err := strconv.Atoi(e)
		if err != nil {
			return err
		}
		*r = append(*r, i)
	}
	return nil
}

// Options defines all the command line options available (all have a default value)
type Options struct {
	Width        int
	Height       int
	RaysPerPixel RaysPerPixelList
	Output       string
	Seed         int64
	CPU          int
	Scene        int
}

// saveImage saves the image (if requested) to a file in png format
func saveImage(pixels scene.Pixels, options Options) error {
	if options.Output == "" {
		return nil
	}
	
	f, err := os.OpenFile(options.Output, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	
	img := image.NewNRGBA(image.Rect(0, 0, options.Width, options.Height))
	
	k := 0
	for y := 0; y < options.Height; y++ {
		for x := 0; x < options.Width; x++ {
			img.Set(x, y, clr.NRGBA{
				R: uint8(pixels[k] >> 16 & 0xFF),
				G: uint8(pixels[k] >> 8 & 0xFF),
				B: uint8(pixels[k] & 0xFF),
				A: 255,
			})
			k++
		}
	}
	
	if err := png.Encode(f, img); err != nil {
		return err
	}
	
	fmt.Printf("Image saved to %v\n", options.Output)
	return nil
	
}

// main parses the options, set up the Window/Screen, builds the world and renders the scene.
// As the scene gets rendered the screen gets refreshed regularly to show progress. When the image is fully
// rendered, it saves it to a file (if the output option is set)
func main() {
	options := Options{}
	
	flag.IntVar(&options.Width, "w", 800, "width in pixel")
	flag.IntVar(&options.Height, "h", 400, "height in pixel")
	flag.IntVar(&options.CPU, "cpu", runtime.NumCPU(), "number of CPU to use (default to number of CPU available)")
	flag.Int64Var(&options.Seed, "seed", 2017, "seed for random number generator")
	flag.Var(&options.RaysPerPixel, "r", "comma separated list (or multiple) rays per pixel")
	flag.IntVar(&options.Scene, "scene", 1, "choose a scene to build")
	flag.StringVar(&options.Output, "o", "", "path to file for saving (do not save if not defined)")
	
	flag.Parse()
	
	if len(options.RaysPerPixel) == 0 {
		options.RaysPerPixel = []int{1, 99}
	}
	
	// initializes the random number generator (since the scene has random spheres... to be reproducible)
	rand.Seed(options.Seed)
	
	// initializes SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	
	// create (and show) window
	window, err := sdl.CreateWindow("Ray Tracer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(options.Width), int32(options.Height), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	
	// retrieves the screen
	screen, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	
	// clear the screen (otherwise there is garbage...)
	err = screen.FillRect(&sdl.Rect{W: int32(options.Width), H: int32(options.Height)}, 0x00000000)
	if err != nil {
		panic(err)
	}
	
	// Camera & World
	c, w := scene.NewBuilder(options.Width, options.Height, options.Scene)
	
	scene := scene.NewScene(options.Width, options.Height, options.RaysPerPixel, c, w)
	pixels, completed := scene.Render(options.CPU)
	
	// update the surface to show it
	if err := window.UpdateSurface(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
	}
	
	updateDisplay := true
	
	// poll for quit event
	for running := true; running; {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}
		
		// wait a few between iterations
		sdl.Delay(16)
		
		if updateDisplay {
			
			// if an error occurs it panic.
			if err := scene.Display(window, screen, pixels); err != nil {
				panic(err)
			}
			
			// check (non blocking thanks to select) that the image is completely rendered
			select {
			case <-completed:
				updateDisplay = false
				fmt.Println("Render complete.")
				if err := saveImage(pixels, options); err != nil {
					fmt.Printf("Error while saving the image [%v]\n", err)
				}
			default:
				break
			}
			
		}
	}
}
