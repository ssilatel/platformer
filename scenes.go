package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SceneInterface interface {
	Update(sm *SceneManager, player *Player)
	Draw(screen *ebiten.Image, sm *SceneManager, player *Player)
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

func (sm *SceneManager) UpdateScenes(player *Player) {
	sm.Scenes[len(sm.Scenes)-1].Update(sm, player)
}

func (sm *SceneManager) DrawScenes(screen *ebiten.Image, player *Player) {
	sm.Scenes[len(sm.Scenes)-1].Draw(screen, sm, player)
}

type MainMenuScene struct{}

func (s *MainMenuScene) Update(sm *SceneManager, player *Player) {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		sm.EnterScene(NewGameScene())
	}
}

func (s *MainMenuScene) Draw(screen *ebiten.Image, sm *SceneManager, player *Player) {
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

func (s *GameScene) Update(sm *SceneManager, player *Player) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		sm.ExitScene()
		//resetInputs(inputs)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		player.Sprite.X = 20
		player.Sprite.Y = 120
		player.CurrentState = "normal"
	}

	player.Update(s.Tilemap)

	sm.Camera.X += (player.Bb.X + player.Bb.W/2) - screenWidth/2 - sm.Camera.X
	sm.Camera.Y += (player.Bb.Y + player.Bb.H/2) - screenHeight/2 - sm.Camera.Y
}

func (s *GameScene) Draw(screen *ebiten.Image, sm *SceneManager, player *Player) {
	screen.Fill(color.RGBA{34, 34, 35, 255})

	s.Tilemap.Draw(screen, sm.Camera)
	player.Draw(screen, sm.Camera)

	ebitenutil.DebugPrint(screen, "Game Scene")
}
