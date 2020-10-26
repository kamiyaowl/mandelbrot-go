package main

import (
	"log"
	"mandelbrot-go/pkg/mandelbrot"

	"github.com/hajimehoshi/ebiten"
)

func main() {
	// TODO: commandline argsから取得?
	width := 800
	height := 600
	windowScale := 2

	ebiten.SetWindowSize(width*windowScale, height*windowScale)
	ebiten.SetWindowTitle("mandelbrot-go")
	g := mandelbrot.NewDefaultParam(width, height)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
