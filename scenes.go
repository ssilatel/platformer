package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	t.LoadTiles(spritesheet, "data/level1_tile_layer.csv", 20, tileSize, 1, "black")
	t.LoadTiles(spritesheet, "data/level1_red_layer.csv", 20, tileSize, 1, "red")
	t.LoadTiles(spritesheet, "data/level1_green_layer.csv", 20, tileSize, 1, "green")
	t.LoadTiles(spritesheet, "data/level1_blue_layer.csv", 20, tileSize, 1, "blue")
	return &GameScene{
		Tilemap: t,
	}
}

func (s *GameScene) Update(sm *SceneManager, player *Player) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		sm.ExitScene()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		player.Sprite.X = 20
		player.Sprite.Y = 120
		player.CurrentState = "normal"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		player.Sprite.X = 1220
		player.Sprite.Y = 190
	}

	player.Update(s.Tilemap)

	sm.Camera.X += (player.Bb.X + player.Bb.W/2) - screenWidth/2 - sm.Camera.X
	sm.Camera.Y += (player.Bb.Y + player.Bb.H/2) - screenHeight/2 - sm.Camera.Y
}

func (s *GameScene) Draw(screen *ebiten.Image, sm *SceneManager, player *Player) {
	screen.Fill(color.RGBA{34, 34, 35, 255})

	if ebiten.IsKeyPressed(ebiten.KeyJ) && player.HasRed {
		vector.DrawFilledCircle(screen, float32((player.Sprite.X+player.Sprite.W/2)-sm.Camera.X), float32((player.Sprite.Y+player.Sprite.H/2)-sm.Camera.Y), float32(20), color.RGBA{255, 0, 0, 255}, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) && player.HasGreen {
		vector.DrawFilledCircle(screen, float32((player.Sprite.X+player.Sprite.W/2)-sm.Camera.X), float32((player.Sprite.Y+player.Sprite.H/2)-sm.Camera.Y), float32(20), color.RGBA{0, 255, 0, 255}, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) && player.HasBlue {
		vector.DrawFilledCircle(screen, float32((player.Sprite.X+player.Sprite.W/2)-sm.Camera.X), float32((player.Sprite.Y+player.Sprite.H/2)-sm.Camera.Y), float32(20), color.RGBA{0, 0, 255, 255}, false)
	}
	s.Tilemap.Draw(screen, sm.Camera)
	player.Draw(screen, sm.Camera)

	ebitenutil.DebugPrint(screen, "Game Scene")
}
