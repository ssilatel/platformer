package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SceneInterface interface {
	Update(sm *SceneManager, inputs map[string]bool, player *Player)
	Draw(screen *ebiten.Image)
}

type SceneManager struct {
	Scenes []SceneInterface
}

func (sm *SceneManager) EnterScene(newScene SceneInterface) {
	sm.Scenes = append(sm.Scenes, newScene)
}

func (sm *SceneManager) ExitScene() {
	sm.Scenes = sm.Scenes[:len(sm.Scenes)-1]
}

func (sm *SceneManager) UpdateScenes(inputs map[string]bool, player *Player) {
	sm.Scenes[len(sm.Scenes)-1].Update(sm, inputs, player)
}

func (sm *SceneManager) DrawScenes(screen *ebiten.Image) {
	sm.Scenes[len(sm.Scenes)-1].Draw(screen)
}

type MainMenuScene struct{}

func (s *MainMenuScene) Update(sm *SceneManager, inputs map[string]bool, player *Player) {
	if inputs["space"] {
		sm.EnterScene(&GameScene{})
	}
	resetInputs(inputs)
}

func (s *MainMenuScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{34, 34, 35, 255})
	ebitenutil.DebugPrintAt(screen, "Main Menu", screenWidth/2-30, screenHeight/2-20)
	ebitenutil.DebugPrintAt(screen, "Press SPACE to start", screenWidth/2-60, screenHeight/2)
}

type GameScene struct{}

func (s *GameScene) Update(sm *SceneManager, inputs map[string]bool, player *Player) {
	if inputs["space"] {
		sm.ExitScene()
	}
	resetInputs(inputs)
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{100, 200, 100, 255})
	ebitenutil.DebugPrint(screen, "Game Scene")
}
