package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) HandleInputs() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Inputs["escape"] = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.Inputs["space"] = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.Inputs["up"] = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.Inputs["down"] = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.Inputs["left"] = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.Inputs["right"] = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Inputs["escape"] = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Inputs["r"] = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.Inputs["e"] = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.Inputs["q"] = true
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		g.Inputs["escape"] = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		g.Inputs["space"] = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.Inputs["up"] = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.Inputs["down"] = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.Inputs["left"] = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.Inputs["right"] = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		g.Inputs["escape"] = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		g.Inputs["r"] = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyE) {
		g.Inputs["e"] = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyQ) {
		g.Inputs["q"] = false
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Inputs["leftClick"] = true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.Inputs["rightClick"] = true
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.Inputs["leftClick"] = false
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		g.Inputs["rightClick"] = false
	}
}

func resetInputs(inputs map[string]bool) {
	for k := range inputs {
		inputs[k] = false
	}
}
