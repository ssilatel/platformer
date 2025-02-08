package main

const (
	screenWidth  = 320
	screenHeight = 240
	tileSize     = 16
)

type Inputs map[string]bool

type Rect struct {
	X, Y, W, H float64
}
