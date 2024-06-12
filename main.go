package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"
	"slices"

	"github.com/falanger/hexgrid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	windowWidth  = 1280
	windowHeight = 720
)

var (
	//go:embed tiles.png
	Tiles_png  []byte
	tilesImage *ebiten.Image
)

type Game struct {
	viewport viewport
	layout   hexgrid.Layout
	grid     []Tile
}

type viewport struct {
	x int
	y int
}

type Tile struct {
	hex hexgrid.Hex
	typ int
}

func (g *Game) init() {
	g.layout = hexgrid.Layout{
		Origin:      hexgrid.Point{X: 0, Y: 0},
		Size:        hexgrid.Point{X: 112 / 2, Y: 112 / 2},
		Orientation: hexgrid.OrientationFlat,
	}

	img, _, err := image.Decode(bytes.NewReader(Tiles_png))
	if err != nil {
		log.Fatal(err)
	}

	tilesImage = ebiten.NewImageFromImage(img)

	// g.grid = make([]Tile, 25)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			g.grid = append(g.grid, Tile{typ: 0, hex: hexgrid.NewHex(i, j)})
		}
	}

}

func (v *viewport) handleMove() {
	presedKeys := inpututil.PressedKeys()
	if slices.Contains(presedKeys, ebiten.KeyLeft) {
		v.x += 5
	}

	if slices.Contains(presedKeys, ebiten.KeyRight) {
		v.x -= 5
	}

	if slices.Contains(presedKeys, ebiten.KeyUp) {
		v.y += 5
	}

	if slices.Contains(presedKeys, ebiten.KeyDown) {
		v.y -= 5
	}
}

func (g *Game) Update() error {
	g.viewport.handleMove()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &colorm.DrawImageOptions{}

	var c colorm.ColorM

	c.ChangeHSV(-1.87, 1, 1)

	for i, t := range g.grid {
		point := hexgrid.HexToPixel(g.layout, t.hex)

		op.GeoM.Reset()
		// op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(point.X, point.Y)
		op.GeoM.Translate(float64(g.viewport.x), float64(g.viewport.y))
		x := (112 * (i % 3))
		colorm.DrawImage(screen, tilesImage.SubImage(image.Rect(x, 0, x+112, 97)).(*ebiten.Image), c, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	game := Game{}

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Hello, World!")

	game.init()

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
