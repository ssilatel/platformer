package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// GREY 34, 34, 35
// BLUE 34, 34, 130

type Scroll struct {
	X, Y float64
}

type Game struct {
	Player        *Player
	Tilemap       *Tilemap
	Scroll        Scroll
	Colours       map[string]color.RGBA
	CurrentColour string

	Debug       bool
	MouseX      int
	MouseY      int
	ClickedTile *Tile
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
		g.ChangeColour()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		g.Debug = !g.Debug
	}

	return nil
}

func (g *Game) ChangeColour() {
	if g.CurrentColour == "grey" {
		g.CurrentColour = "blue"
	} else {
		g.CurrentColour = "grey"
	}

	if g.CurrentColour != "grey" {
		for _, row := range g.Tilemap.Tiles {
			for _, t := range row {
				if t.Colour == g.CurrentColour {
					t.Collidable = false
				} else {
					t.Collidable = true
				}
			}
		}
	} else {
		for _, row := range g.Tilemap.Tiles {
			for _, t := range row {
				if t.Colour != "grey" {
					t.Collidable = true
				}
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.Colours[g.CurrentColour])
	g.Tilemap.Draw(screen, g.Scroll)
	g.Player.Draw(screen, g.Scroll)

	if g.Debug {
		for _, t := range g.Tilemap.TilesVisible(g.Scroll) {
			vector.StrokeRect(screen, float32(t.Bb.X-g.Scroll.X), float32(t.Bb.Y-g.Scroll.Y), float32(t.Bb.W), float32(t.Bb.H), 1, color.RGBA{255, 0, 0, 255}, false)
		}
		vector.StrokeRect(screen, float32(g.Player.Bb.X-g.Scroll.X), float32(g.Player.Bb.Y-g.Scroll.Y), float32(g.Player.Bb.W), float32(g.Player.Bb.H), 1, color.RGBA{255, 0, 0, 255}, false)

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			g.MouseX = int((float64(mx) + g.Scroll.X) / tileSize)
			g.MouseY = int((float64(my) + g.Scroll.Y) / tileSize)
			g.ClickedTile = g.Tilemap.Tiles[g.MouseY][g.MouseX]
		}
		if g.ClickedTile.Image != nil {
			vector.StrokeRect(screen, float32(g.ClickedTile.Bb.X-g.Scroll.X), float32(g.ClickedTile.Bb.Y-g.Scroll.Y), float32(g.ClickedTile.Bb.W), float32(g.ClickedTile.Bb.H), 1, color.RGBA{0, 0, 255, 255}, false)
			tilePosX := fmt.Sprintf("%f", g.ClickedTile.Bb.X)
			tilePosY := fmt.Sprintf("%f", g.ClickedTile.Bb.Y)
			ebitenutil.DebugPrintAt(screen, tilePosX, int(g.ClickedTile.Bb.X-g.Scroll.X), int(g.ClickedTile.Bb.Y-20-g.Scroll.Y))
			ebitenutil.DebugPrintAt(screen, tilePosY, int(g.ClickedTile.Bb.X-g.Scroll.X), int(g.ClickedTile.Bb.Y-10-g.Scroll.Y))
		}
		posX := fmt.Sprintf("%f", g.Player.Sprite.X)
		posY := fmt.Sprintf("%f", g.Player.Sprite.Y)
		ebitenutil.DebugPrintAt(screen, posX, int(g.Player.Sprite.X-g.Scroll.X), int(g.Player.Sprite.Y-20-g.Scroll.Y))
		ebitenutil.DebugPrintAt(screen, posY, int(g.Player.Sprite.X-g.Scroll.X), int(g.Player.Sprite.Y-10-g.Scroll.Y))
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
	spritesheetBlue, _, err := ebitenutil.NewImageFromFile("assets/monochrome_spritesheet_blue.png")
	if err != nil {
		log.Fatal(err)
	}

	player := Player{
		Sprite:       Rect{20, 120, tileSize, tileSize},
		Bb:           Rect{20, 120, 12, 12},
		Oldbb:        Rect{100, 100, tileSize, tileSize},
		OffsetX:      2,
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
	tilemap.LoadTiles(spritesheet, "data/level1_tile_layer.csv", 20, tileSize, 1, "grey")
	tilemap.LoadTiles(spritesheetBlue, "data/level1_blue_layer.csv", 20, tileSize, 1, "blue")

	g := &Game{
		Player:  &player,
		Tilemap: tilemap,
		Colours: map[string]color.RGBA{
			"grey": {34, 34, 35, 255},
			"blue": {34, 34, 130, 255},
		},
		CurrentColour: "grey",
	}

	ebiten.SetWindowTitle("qwer")
	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)

	if err = ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
