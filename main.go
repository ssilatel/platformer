package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// GREY 34, 34, 35
// BLUE 34, 34, 130

// /// TEMPORARY
type Player struct{}

///// TEMPORARY

type Game struct {
	Inputs       Inputs
	SceneManager SceneManager
	Player       *Player
}

func (g *Game) Update() error {
	g.HandleInputs()
	g.SceneManager.UpdateScenes(g.Inputs, g.Player)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.SceneManager.DrawScenes(screen)
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

	sm := SceneManager{
		Scenes: []SceneInterface{&MainMenuScene{}},
	}

	g := &Game{
		Inputs:       inputs,
		SceneManager: sm,
		Player:       &Player{},
	}

	ebiten.SetWindowTitle("asdf")
	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
