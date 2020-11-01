package mandelbrot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMandelbrot_readPaletteFromCsv(t *testing.T) {
	g := NewDefaultParam(128, 128, "")
	g.readPaletteFromCsv("../../palette.csv")
	assert.Equal(t, Color{r: 66, g: 30, b: 15, a: 255}, g.palette[0])
	assert.Equal(t, Color{r: 25, g: 7, b: 26, a: 255}, g.palette[1])
	assert.Equal(t, Color{r: 9, g: 1, b: 47, a: 255}, g.palette[2])
	assert.Equal(t, Color{r: 4, g: 4, b: 73, a: 255}, g.palette[3])
	assert.Equal(t, Color{r: 0, g: 7, b: 100, a: 255}, g.palette[4])
	assert.Equal(t, Color{r: 12, g: 44, b: 138, a: 255}, g.palette[5])
	assert.Equal(t, Color{r: 24, g: 82, b: 177, a: 255}, g.palette[6])
	assert.Equal(t, Color{r: 57, g: 125, b: 209, a: 255}, g.palette[7])
	assert.Equal(t, Color{r: 134, g: 181, b: 229, a: 255}, g.palette[8])
	assert.Equal(t, Color{r: 211, g: 236, b: 248, a: 255}, g.palette[9])
	assert.Equal(t, Color{r: 241, g: 233, b: 191, a: 255}, g.palette[10])
	assert.Equal(t, Color{r: 248, g: 201, b: 95, a: 255}, g.palette[11])
	assert.Equal(t, Color{r: 255, g: 170, b: 0, a: 255}, g.palette[12])
	assert.Equal(t, Color{r: 204, g: 128, b: 0, a: 255}, g.palette[13])
	assert.Equal(t, Color{r: 153, g: 87, b: 0, a: 255}, g.palette[14])
	assert.Equal(t, Color{r: 106, g: 52, b: 3, a: 255}, g.palette[15])
}

func TestMandelbrot_numOfCalcUntilDivergence(t *testing.T) {
	g := NewDefaultParam(128, 128, "")
	tests := []struct {
		a    float64
		b    float64
		iter int
		want int
	}{
		{a: 0.0, b: 0.0, iter: 100, want: 100},
		{a: 0.0, b: 0.0, iter: 1000, want: 1000},
		{a: -2.0, b: -1.0, iter: 100, want: 1},
		{a: -1.4, b: -0.01, iter: 100, want: 31},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			g.iterMax = tt.iter
			if got := g.numOfCalcUntilDivergence(tt.a, tt.b); got != tt.want {
				t.Fatalf("want = %d, got  %d", tt.want, got)
			}
		})
	}
}

func TestMandelbrot_initOffscreen(t *testing.T) {
	g := NewDefaultParam(128, 128, "")
	// 初期化直後はallocate必須
	assert.Nil(t, g.offscreenImage)
	assert.Nil(t, g.workbuffer)
	// 1回呼び出して確保されることを見る
	g.initOffscreen()
	assert.NotNil(t, g.offscreenImage)
	assert.NotNil(t, g.workbuffer)
	assert.Equal(t, 128, g.offscreenImage.Bounds().Dx())
	assert.Equal(t, 128, g.offscreenImage.Bounds().Dy())
	assert.Equal(t, 128*128*4, len(g.workbuffer))
	// サイズ変更への対応
	g.width = 512
	g.height = 256
	g.initOffscreen()
	assert.NotNil(t, g.offscreenImage)
	assert.NotNil(t, g.workbuffer)
	assert.Equal(t, 512, g.offscreenImage.Bounds().Dx())
	assert.Equal(t, 256, g.offscreenImage.Bounds().Dy())
	assert.Equal(t, 512*256*4, len(g.workbuffer))
}
