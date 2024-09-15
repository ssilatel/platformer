package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	Player  *Player
	Tilemap *Tilemap
}

func (g *Game) Update() error {
	g.Player.Update(g.Tilemap)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{34, 34, 35, 255})
	g.Tilemap.Draw(screen)
	g.Player.Draw(screen)
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth * 2, screenHeight * 2
}

func main() {
	spritesheet, _, err := ebitenutil.NewImageFromFile("assets/monochrome_spritesheet.png")
	if err != nil {
		log.Fatal(err)
	}

	player := Player{
		Sprite:     Rect{100, 100, tileSize, tileSize},
		Bb:         Rect{100, 100, 14, 12},
		Oldbb:      Rect{100, 100, tileSize, tileSize},
		OffsetX:    1,
		OffsetY:    4,
		Vx:         0,
		Vy:         0,
		Image:      spritesheet.SubImage(image.Rect(0, 238, 16, 254)).(*ebiten.Image),
		Collisions: make(map[string]bool),
	}

	tilemap := NewTilemap(40, 30)
	tilemap.LoadTiles(spritesheet, "level.csv", 20, tileSize, 1)

	g := &Game{
		Player:  &player,
		Tilemap: tilemap,
	}

	ebiten.SetWindowTitle("qwer")
	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)

	if err = ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
