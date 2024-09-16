package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Scroll struct {
	X, Y float64
}

type Game struct {
	Player  *Player
	Tilemap *Tilemap
	Scroll  Scroll
	Debug   bool
}

func (g *Game) Update() error {
	g.Player.Update(g.Tilemap)

	g.Scroll.X += (g.Player.Bb.X + g.Player.Bb.W/2) - screenWidth/2 - g.Scroll.X
	g.Scroll.Y += (g.Player.Bb.Y + g.Player.Bb.H/2) - screenHeight/2 - g.Scroll.Y

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Player.Sprite.X = 20
		g.Player.Sprite.Y = 120
		g.Player.CurrentState = "normal"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.Debug = !g.Debug
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{34, 34, 35, 255})
	g.Tilemap.Draw(screen, g.Scroll)
	g.Player.Draw(screen, g.Scroll)

	if g.Debug {
		for _, row := range g.Tilemap.Tiles {
			for _, t := range row {
				vector.StrokeRect(screen, float32(t.Bb.X-g.Scroll.X), float32(t.Bb.Y-g.Scroll.Y), float32(t.Bb.W), float32(t.Bb.H), 1, color.RGBA{255, 0, 0, 255}, false)
			}
		}
		vector.StrokeRect(screen, float32(g.Player.Bb.X-g.Scroll.X), float32(g.Player.Bb.Y-g.Scroll.Y), float32(g.Player.Bb.W), float32(g.Player.Bb.H), 1, color.RGBA{255, 0, 0, 255}, false)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	spritesheet, _, err := ebitenutil.NewImageFromFile("assets/monochrome_spritesheet.png")
	if err != nil {
		log.Fatal(err)
	}

	player := Player{
		Sprite:       Rect{20, 120, tileSize, tileSize},
		Bb:           Rect{20, 120, 14, 12},
		Oldbb:        Rect{100, 100, tileSize, tileSize},
		OffsetX:      1,
		OffsetY:      4,
		Vx:           0,
		Vy:           0,
		Image:        spritesheet.SubImage(image.Rect(0, 238, 16, 254)).(*ebiten.Image),
		Collisions:   make(map[string]bool),
		Animations:   make(map[string]*Animation),
		States:       make(map[string]State),
		CurrentState: "normal",
	}
	player.AddAnimation("idle", CreateFramesFromSpritesheetHorizontal(spritesheet, 16, 16, 1, 0, 204, 0), 10)
	player.AddAnimation("run", CreateFramesFromSpritesheetHorizontal(spritesheet, 16, 16, 3, 17, 204, 1), 10)
	player.AddAnimation("jump", CreateFramesFromSpritesheetHorizontal(spritesheet, 16, 16, 1, 68, 204, 0), 10)
	player.AddAnimation("death", CreateFramesFromSpritesheetHorizontal(spritesheet, 16, 16, 1, 102, 204, 0), 10)
	player.SetAnimation("idle")

	player.AddState("normal", &PlayerNormalState{})
	player.AddState("death", &PlayerDeathState{})

	tilemap := NewTilemap(100, 40)
	tilemap.LoadTiles(spritesheet, "level1.csv", 20, tileSize, 1)

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
