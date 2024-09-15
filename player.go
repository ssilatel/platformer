package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Rect struct {
	X, Y, W, H float64
}

type Player struct {
	Sprite           Rect
	Bb               Rect
	Oldbb            Rect
	OffsetX, OffsetY float64
	Vx, Vy           float64
	Image            *ebiten.Image
	Collisions       map[string]bool
}

func (p *Player) Update(tilemap *Tilemap) {
	p.Oldbb = p.Bb
	p.Bb.X = p.Sprite.X + p.OffsetX
	p.Bb.Y = p.Sprite.Y + p.OffsetY

	p.Collisions["top"] = false
	p.Collisions["bottom"] = false
	p.Collisions["left"] = false
	p.Collisions["right"] = false

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.Vx = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Vx = 1
	} else {
		p.Vx = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		p.Vy = -5.5
	}

	p.Bb.X += p.Vx
	p.Bb.Y += p.Vy

	tilesAround := tilemap.TilesAround(p.Bb.X, p.Bb.Y)

	for _, t := range tilesAround {
		if HasCollided(&p.Bb, &t.Bb) {
			if p.Bb.X <= t.Bb.X+t.Bb.W && p.Oldbb.X >= t.Bb.X+t.Bb.W {
				p.Bb.X = t.Bb.X + t.Bb.W
				p.Collisions["left"] = true
			}
			if p.Bb.X+p.Bb.W >= t.Bb.X && p.Oldbb.X+p.Oldbb.W <= t.Bb.X {
				p.Bb.X = t.Bb.X - p.Bb.W
				p.Collisions["right"] = true
			}
		}
	}
	p.Sprite.X = p.Bb.X - p.OffsetX

	for _, t := range tilesAround {
		if HasCollided(&p.Bb, &t.Bb) {
			if p.Bb.Y <= t.Bb.Y+t.Bb.H && p.Oldbb.Y >= t.Bb.Y+t.Bb.H {
				p.Bb.Y = t.Bb.Y + t.Bb.H
				p.Collisions["top"] = true
			}
			if p.Bb.Y+p.Bb.H >= t.Bb.Y && p.Oldbb.Y+p.Oldbb.H <= t.Bb.Y {
				p.Bb.Y = t.Bb.Y - p.Bb.H
				p.Collisions["bottom"] = true
			}
		}
	}
	p.Sprite.Y = p.Bb.Y - p.OffsetY

	p.Vy = math.Min(5, p.Vy+0.4)
	if p.Collisions["bottom"] || p.Collisions["top"] {
		p.Vy = 0
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.Sprite.X, p.Sprite.Y)
	screen.DrawImage(p.Image, opts)
}

func HasCollided(r1, r2 *Rect) bool {
	return r1.X < r2.X+r2.W && r1.X+r1.W > r2.X && r1.Y < r2.Y+r2.H && r1.Y+r1.H > r2.Y
}
