package mandelbrot

import (
	"fmt"
	"math/cmplx"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Mandelbrotを表示するゲーム構造体
type Game struct {
	// 描画幅
	width int
	// 描画高さ
	height int
	// 収束判定時の反復計算回数
	iterMax int
	// 描画中心のx座標
	centerX float64
	// 描画中心のy座標
	centerY float64
	// pixel間の距離
	distancePerPixel float64
	// Z_0項のX値、ジュリア集合を描画するなら0以外に設定する
	z0x float64
	// Z_0項のY値、ジュリア集合を描画するなら0以外に設定する
	z0y float64
	// パラメータの変更があり、offscreen bufferの更新が必要
	isParamChanged bool
	// オフスクリーンバッファ
	offscreenImage *ebiten.Image
	// offscreenImage計算用のworkbuffer
	workbuffer []byte
}

// 描画範囲以外を初期値で埋めた構造体を返します
func NewDefaultParam(width, height int) *Game {
	g := Game{
		width:            width,
		height:           height,
		iterMax:          256,
		centerX:          0.0,
		centerY:          0.0,
		distancePerPixel: 0.005,
		z0x:              0,
		z0y:              0,
		isParamChanged:   true,
		offscreenImage:   nil, // drawOffscreen -> initOffscreen で初期化
		workbuffer:       nil, // drawOffscreen -> initOffscreen で初期化
	}
	return &g
}

// 発散するまでにかかったループ数を返します。iterMaxの値で頭打ちされます
// C     = a + bi
// Z_n+1 = (Z_n)^2 + C
func (g *Game) numOfCalcUntilDivergence(a, b float64) int {
	var c complex128 = complex(a, b)
	var z complex128 = complex(g.z0x, g.z0y)
	for i := 0; i < g.iterMax; i++ {
		// 発散する場合は現在のループ数を返す
		if cmplx.Abs(z) >= 2 {
			return i
		}
		// 次のZを計算
		z = cmplx.Pow(z, 2) + c
	}
	// 発散しなかった
	return g.iterMax
}

// offscreenImageにマンデルブロ集合を描画します
func (g *Game) initOffscreen() {
	// オフスクリーンバッファを取得していない、もしくはサイズが変更された場合は再生成する
	isNeedAllocate := (g.offscreenImage == nil)
	if !isNeedAllocate {
		// 未初期化で触れないため
		rect := g.offscreenImage.Bounds()
		isNeedAllocate = (rect.Dx() != g.width) || (rect.Dy() != g.height)
	}
	if isNeedAllocate {
		g.offscreenImage, _ = ebiten.NewImage(g.width, g.height, ebiten.FilterDefault)
		g.workbuffer = make([]byte, g.width*g.height*4) // 4 = RGBA
	}
}

// offscreenImageにマンデルブロ集合を描画します
func (g *Game) drawOffscreen() {
	// オフスクリーンバッファを取得していない、もしくはサイズが変更された場合は再生成する
	g.initOffscreen()
	// 左上を計算
	var x0 float64 = g.centerX - ((float64(g.width) / 2.0) * g.distancePerPixel)
	var y0 float64 = g.centerY - ((float64(g.height) / 2.0) * g.distancePerPixel)
	// 左上からの差分を計算して処理する
	for j := 0; j < g.height; j++ {
		for i := 0; i < g.width; i++ {
			// 発散までの回数を計算
			x := x0 + float64(i)*g.distancePerPixel
			y := y0 + float64(j)*g.distancePerPixel
			n := g.numOfCalcUntilDivergence(x, y)

			// 対象のpixel
			ptr := ((j * g.width) + i) * 4 // 4 = RGBA
			// TODO: 色をおしゃれに
			if n == g.iterMax {
				g.workbuffer[ptr+0] = 0xff // R
				g.workbuffer[ptr+1] = 0xff // G
				g.workbuffer[ptr+2] = 0xff // B
			} else {
				g.workbuffer[ptr+0] = 0 // R
				g.workbuffer[ptr+1] = 0 // G
				g.workbuffer[ptr+2] = 0 // B
			}
			g.workbuffer[ptr+3] = 0xff // A

		}
	}

	// Imageの中身を更新
	g.offscreenImage.ReplacePixels(g.workbuffer)
}

// tickごとに呼び出されます
func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}

// 描画時に呼び出されます
func (g *Game) Draw(screen *ebiten.Image) {
	// パラメータ変更していたら画像更新
	if g.isParamChanged {
		g.drawOffscreen()
		g.isParamChanged = false
	}
	// 画像を描画
	screen.DrawImage(g.offscreenImage, nil)

	// debug print
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f C: (%0.8f, %0.8f) iter: %d distance/pixel: %f", ebiten.CurrentFPS(), g.centerX, g.centerY, g.iterMax, g.distancePerPixel))
}

// screen size取得時に呼び出されます
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}
