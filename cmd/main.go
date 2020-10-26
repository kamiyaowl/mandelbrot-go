package main

import (
	"log"
	"mandelbrot-go/pkg/mandelbrot"

	"github.com/hajimehoshi/ebiten"
)

func main() {
	// ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("mandelbrot-go")
	if err := ebiten.RunGame(&mandelbrot.Game{}); err != nil {
		log.Fatal(err)
	}
}
