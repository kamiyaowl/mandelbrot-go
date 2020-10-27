package mandelbrot

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/cmplx"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// Paletteで使用する色情報
type Color struct {
	r byte
	g byte
	b byte
	a byte
}

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
	// 描画色, nilのままであればデフォルトの設定で描画
	palette []Color
}

// CSVファイルから色情報を読み込みます
func (g *Game) readPaletteFromCsv(palettePath string) {
	// read from file
	file, err := os.Open(palettePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	isHeaderSkipped := false
	p := []Color{}
	r := csv.NewReader(file)
	for {
		// record: [4]string{r,g,b,a}
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		// other error
		if err != nil {
			panic(err)
		}
		// skip header
		if !isHeaderSkipped {
			isHeaderSkipped = true
			continue
		}
		// 要素数が満たなければSkip
		if len(record) < 4 {
			continue
		}
		// add color
		var r int
		var g int
		var b int
		var a int
		if r, err = strconv.Atoi(record[0]); err != nil {
			panic(err)
		}
		if g, err = strconv.Atoi(record[1]); err != nil {
			panic(err)
		}
		if b, err = strconv.Atoi(record[2]); err != nil {
			panic(err)
		}
		if a, err = strconv.Atoi(record[3]); err != nil {
			panic(err)
		}

		c := Color{
			r: byte(r),
			g: byte(g),
			b: byte(b),
			a: byte(a),
		}
		p = append(p, c)
	}

	// 要素が存在していれば新規に置き換える
	if len(p) > 0 {
		g.palette = p
	}
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

// 描画色を取得します, Paletteが事前に読み込まれていなければデフォルト設定(収束:白, 発散: 黒)になります
func (g *Game) getColor(n int) *Color {
	// default color
	if g.palette == nil {
		if n == g.iterMax {
			return &Color{
				r: 255,
				g: 255,
				b: 255,
				a: 255,
			}
		} else {
			return &Color{
				r: 0,
				g: 0,
				b: 0,
				a: 0,
			}
		}
	}

	// from palette array
	index := n % len(g.palette)
	return &g.palette[index]
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
			// 色を決定して載せておく
			c := g.getColor(n)
			g.workbuffer[ptr+0] = c.r
			g.workbuffer[ptr+1] = c.g
			g.workbuffer[ptr+2] = c.b
			g.workbuffer[ptr+3] = c.a

		}
	}

	// Imageの中身を更新
	g.offscreenImage.ReplacePixels(g.workbuffer)
}

// 描画範囲以外を初期値で埋めた構造体を返します
func NewDefaultParam(width int, height int, palettePath string) *Game {
	g := Game{
		width:            width,
		height:           height,
		iterMax:          1024,
		centerX:          -0.5,
		centerY:          0.0,
		distancePerPixel: 0.009155,
		z0x:              0,
		z0y:              0,
		isParamChanged:   true,
		offscreenImage:   nil, // drawOffscreen -> initOffscreen で初期化
		workbuffer:       nil, // drawOffscreen -> initOffscreen で初期化
		palette:          nil, // readPaletteFromCsvで作成するか、nilを継続
	}
	// paletteの読み込み。rgbaのCSVになっているのでpaletteに配列で展開
	if palettePath != "" {
		g.readPaletteFromCsv(palettePath)
	}
	return &g
}

// tickごとに呼び出されます
func (g *Game) Update(screen *ebiten.Image) error {
	// クリックした位置に移動, 連続では反映しない
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// window内に収まっていることを確認してから、実座標に変換する
		cx, cy := ebiten.CursorPosition()
		if cx >= 0 && cy >= 0 && cx < g.width && cy < g.height {
			// スクリーン座標で、中央からのの距離を求める
			sx := cx - (g.width / 2)
			sy := cy - (g.height / 2)
			// 実座標系に変換して、現在地点に足し合わせる
			g.centerX += float64(sx) * g.distancePerPixel
			g.centerY += float64(sy) * g.distancePerPixel

			g.isParamChanged = true
		}
	}
	// 拡大, 縮小
	_, wy := ebiten.Wheel()
	if wy != 0 {
		if wy < 0 {
			g.distancePerPixel *= 1.25
		} else {
			g.distancePerPixel *= 0.75
		}
		g.isParamChanged = true
	}
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
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nC: (%0.8f, %0.8f)\niter: %d\nd/p: %f\n", ebiten.CurrentFPS(), g.centerX, g.centerY, g.iterMax, g.distancePerPixel))
}

// screen size取得時に呼び出されます
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}
