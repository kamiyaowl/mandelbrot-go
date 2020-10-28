package main

import (
	"flag"
	"log"
	"mandelbrot-go/pkg/mandelbrot"

	"github.com/hajimehoshi/ebiten"
)

func main() {
	windowScale := flag.Int("scale", 1, "window scale")
	width := flag.Int("width", 640, "screen width")
	height := flag.Int("height", 640, "screen height")
	palettePath := flag.String("palette_path", "", "color palette csv filepath")
	iterMax := flag.Int("iter", 128, "calculation iteration")
	distancePerPixel := flag.Float64("dpp", 0.009155, "distance/pixel")
	centerX := flag.Float64("cx", -0.5, "center x")
	centerY := flag.Float64("cy", 0.0, "center y")
	z0x := flag.Float64("z0x", 0.0, "z0 x")
	z0y := flag.Float64("z0y", 0.0, "z0 y")
	flag.Parse()

	ebiten.SetWindowSize(*width**windowScale, *height**windowScale)
	ebiten.SetWindowTitle("mandelbrot-go")
	ebiten.SetRunnableOnUnfocused(true)
	g := mandelbrot.NewDetailParam(*width, *height, *palettePath, *iterMax, *distancePerPixel, *centerX, *centerY, *z0x, *z0y)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
