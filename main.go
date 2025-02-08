package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// GREY 34, 34, 35
// BLUE 34, 34, 130

type Game struct {
	Inputs       Inputs
	Spritesheet  *ebiten.Image
	Player       *Player
	SceneManager SceneManager
}

func (g *Game) Update() error {
	g.HandleInputs()
	g.SceneManager.UpdateScenes(g.Inputs, g.Player)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.SceneManager.DrawScenes(screen, g.Player)
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	inputs := Inputs{
		"escape":     false,
		"space":      false,
		"up":         false,
		"down":       false,
		"left":       false,
		"right":      false,
		"r":          false,
		"e":          false,
		"q":          false,
		"leftClick":  false,
		"rightClick": false,
	}

	spritesheet, _, err := ebitenutil.NewImageFromFile("assets/monochrome_spritesheet.png")
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
		States:       make(map[string]PlayerState),
		CurrentState: "normal",
	}
	player.AddAnimation("idle", CreateFramesFromSpritesheetHorizontal(spritesheet, 16, 16, 1, 0, 204, 0), 10)
	player.AddAnimation("run", CreateFramesFromSpritesheetHorizontal(spritesheet, 16, 16, 3, 17, 204, 1), 10)
	player.AddAnimation("jump", CreateFramesFromSpritesheetHorizontal(spritesheet, 16, 16, 1, 68, 204, 0), 10)
	player.AddAnimation("death", CreateFramesFromSpritesheetHorizontal(spritesheet, 16, 16, 1, 102, 204, 0), 10)
	player.SetAnimation("idle")
	player.AddState("normal", &PlayerNormalState{})
	player.AddState("death", &PlayerDeathState{})

	sm := SceneManager{
		Scenes: []SceneInterface{&MainMenuScene{}},
		Camera: &Camera{},
	}

	g := &Game{
		Inputs:       inputs,
		Spritesheet:  spritesheet,
		Player:       &player,
		SceneManager: sm,
	}

	ebiten.SetWindowTitle("asdf")
	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
