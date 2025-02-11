package main

import (
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// GREY 34, 34, 35
// BLUE 34, 34, 130

var audioContext *audio.Context
var single *audio.Player
var grassSound *audio.Player

type Game struct {
	Spritesheet  *ebiten.Image
	Player       *Player
	SceneManager SceneManager
}

func (g *Game) Update() error {
	g.SceneManager.UpdateScenes(g.Player)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.SceneManager.DrawScenes(screen, g.Player)
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	audioContext = audio.NewContext(48000)
	grassSound = createLoop("./assets/sfx_step_grass_r.mp3")
	grassSound.SetVolume(0.33)

	spritesheet, _, err := ebitenutil.NewImageFromFile("assets/monochrome_spritesheet.png")
	if err != nil {
		log.Fatal(err)
	}

	player := Player{
		//Sprite: Rect{20, 120, tileSize, tileSize},
		//Bb:     Rect{20, 120, 12, 12},
		//Sprite:       Rect{70, 760, tileSize, tileSize},
		//Bb:           Rect{70, 760, 12, 12},
		//Sprite:       Rect{230, 120, tileSize, tileSize},
		//Bb:           Rect{230, 120, 12, 12},
		Sprite:       Rect{0, 0, tileSize, tileSize},
		Bb:           Rect{0, 0, 12, 12},
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

func createSingle(filepath string) *audio.Player {
	in, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := mp3.DecodeWithoutResampling(in)
	if err != nil {
		log.Fatal(err)
	}
	singlePlayer, err := audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatal(err)
	}
	return singlePlayer
}

func createLoop(filepath string) *audio.Player {
	in, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := mp3.DecodeWithoutResampling(in)
	if err != nil {
		log.Fatal(err)
	}
	loop := audio.NewInfiniteLoop(stream, 123456)
	loopPlayer, err := audioContext.NewPlayer(loop)
	if err != nil {
		log.Fatal(err)
	}
	return loopPlayer
}
