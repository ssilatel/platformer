package main

import "github.com/hajimehoshi/ebiten/v2"

const (
	screenWidth  = 320
	screenHeight = 240
	tileSize     = 16
)

type Camera struct {
	X, Y float64
}

type Rect struct {
	X, Y, W, H float64
}

type Animation struct {
	Frames        []*ebiten.Image
	CurrentFrame  int
	FrameDuration int
	ElapsedTicks  int
}
