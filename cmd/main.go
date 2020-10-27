package main

import (
	"log"
	"mandelbrot-go/pkg/mandelbrot"

	"github.com/hajimehoshi/ebiten"
)

func main() {
	// TODO: commandline argsから取得?
	width := 400
	height := 300
	windowScale := 2
	palettePath := "palette.csv"

	ebiten.SetWindowSize(width*windowScale, height*windowScale)
	ebiten.SetWindowTitle("mandelbrot-go")
	ebiten.SetRunnableOnUnfocused(true)
	g := mandelbrot.NewDefaultParam(width, height, palettePath)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
