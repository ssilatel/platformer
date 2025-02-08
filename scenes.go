package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SceneInterface interface {
	Update(sm *SceneManager, inputs map[string]bool, player *Player)
	Draw(screen *ebiten.Image, sm *SceneManager)
}

type SceneManager struct {
	Scenes []SceneInterface
	Camera *Camera
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
	sm.Scenes[len(sm.Scenes)-1].Draw(screen, sm)
}

type MainMenuScene struct{}

func (s *MainMenuScene) Update(sm *SceneManager, inputs map[string]bool, player *Player) {
	if inputs["space"] {
		sm.EnterScene(NewGameScene())
	}
	resetInputs(inputs)
}

func (s *MainMenuScene) Draw(screen *ebiten.Image, sm *SceneManager) {
	screen.Fill(color.RGBA{34, 34, 35, 255})
	ebitenutil.DebugPrintAt(screen, "Main Menu", screenWidth/2-30, screenHeight/2-20)
	ebitenutil.DebugPrintAt(screen, "Press SPACE to start", screenWidth/2-60, screenHeight/2)
}

type GameScene struct {
	Tilemap *Tilemap
}

func NewGameScene() *GameScene {
	spritesheet, _, err := ebitenutil.NewImageFromFile("assets/monochrome_spritesheet.png")
	if err != nil {
		log.Fatal(err)
	}
	t := NewTilemap(100, 40)
	t.LoadTiles(spritesheet, "data/level1.csv", 20, tileSize, 1, "grey")
	return &GameScene{
		Tilemap: t,
	}
}

func (s *GameScene) Update(sm *SceneManager, inputs map[string]bool, player *Player) {
	if inputs["left"] {
		sm.Camera.X -= 1
	}
	if inputs["right"] {
		sm.Camera.X += 1
	}

	if inputs["up"] {
		sm.Camera.Y -= 1
	}
	if inputs["down"] {
		sm.Camera.Y += 1
	}

	if inputs["escape"] {
		sm.ExitScene()
	}
}

func (s *GameScene) Draw(screen *ebiten.Image, sm *SceneManager) {
	screen.Fill(color.RGBA{34, 34, 35, 255})
	s.Tilemap.Draw(screen, *sm.Camera)
	ebitenutil.DebugPrint(screen, "Game Scene")
}
