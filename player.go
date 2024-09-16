package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Rect struct {
	X, Y, W, H float64
}

type Animation struct {
	Frames        []*ebiten.Image
	CurrentFrame  int
	FrameDuration int
	ElapsedTicks  int
}

type State interface {
	Update(p *Player, tilemap *Tilemap)
}

type Player struct {
	Sprite           Rect
	Bb               Rect
	Oldbb            Rect
	OffsetX, OffsetY float64
	Vx, Vy           float64
	Image            *ebiten.Image
	Collisions       map[string]bool
	CanJump          bool
	Animations       map[string]*Animation
	CurrentAnimation string
	Flip             bool
	States           map[string]State
	CurrentState     string
}

func (p *Player) Update(tilemap *Tilemap) {
	p.States[p.CurrentState].Update(p, tilemap)
}

func (p *Player) Draw(screen *ebiten.Image, scroll Scroll) {
	opts := &ebiten.DrawImageOptions{}
	if p.Flip {
		opts.GeoM.Scale(-1, 1)
		opts.GeoM.Translate(16, 0)
	}
	opts.GeoM.Translate(p.Sprite.X-scroll.X, p.Sprite.Y-scroll.Y)
	screen.DrawImage(p.Animations[p.CurrentAnimation].Frames[p.Animations[p.CurrentAnimation].CurrentFrame], opts)
}

func HasCollided(r1, r2 *Rect) bool {
	return r1.X < r2.X+r2.W && r1.X+r1.W > r2.X && r1.Y < r2.Y+r2.H && r1.Y+r1.H > r2.Y
}

func (p *Player) AddAnimation(name string, frames []*ebiten.Image, frameDuration int) {
	p.Animations[name] = &Animation{
		Frames:        frames,
		FrameDuration: frameDuration,
	}
}

func (p *Player) SetAnimation(name string) {
	if p.CurrentAnimation != name {
		p.CurrentAnimation = name
		p.Animations[p.CurrentAnimation].CurrentFrame = 0
		p.Animations[p.CurrentAnimation].ElapsedTicks = 0
	}
}

func CreateFramesFromSpritesheetHorizontal(spritesheet *ebiten.Image, frameWidth, frameHeight, frameCount, startX, startY, step int) []*ebiten.Image {
	frames := []*ebiten.Image{}

	for i := 0; i < frameCount; i++ {
		x := startX + (i * (frameWidth + step))
		frame := spritesheet.SubImage(image.Rect(x, startY, x+frameWidth, startY+frameHeight)).(*ebiten.Image)
		frames = append(frames, frame)
	}

	return frames
}

func (p *Player) AddState(name string, state State) {
	p.States[name] = state
}

type PlayerNormalState struct{}

func (s *PlayerNormalState) Update(p *Player, tilemap *Tilemap) {
	p.Oldbb = p.Bb
	p.Bb.X = p.Sprite.X + p.OffsetX
	p.Bb.Y = p.Sprite.Y + p.OffsetY

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.Vx = -1.4
		p.Flip = true
		if p.Collisions["bottom"] {
			p.SetAnimation("run")
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Vx = 1.4
		p.Flip = false
		if p.Collisions["bottom"] {
			p.SetAnimation("run")
		}
	} else {
		p.Vx = 0
		if p.Collisions["bottom"] {
			p.SetAnimation("idle")
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) && p.CanJump {
		p.Vy = -6
		p.CanJump = false
		p.SetAnimation("jump")
	}

	p.Collisions["top"] = false
	p.Collisions["bottom"] = false
	p.Collisions["left"] = false
	p.Collisions["right"] = false

	p.Bb.X += p.Vx
	p.Bb.Y += p.Vy

	tilesAround := tilemap.TilesAround(p.Bb.X, p.Bb.Y)

	for _, t := range tilesAround {
		if HasCollided(&p.Bb, &t.Bb) {
			if t.Type == "spike" {
				p.CurrentState = "death"
			}
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
			if t.Type == "spike" {
				p.CurrentState = "death"
			}
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
	if p.Vy > 1.8 {
		p.SetAnimation("jump")
		p.CanJump = false
	}

	if p.Collisions["bottom"] || p.Collisions["top"] {
		p.Vy = 0
	}
	if p.Collisions["bottom"] {
		p.CanJump = true
	}

	p.Animations[p.CurrentAnimation].ElapsedTicks++
	if p.Animations[p.CurrentAnimation].ElapsedTicks >= p.Animations[p.CurrentAnimation].FrameDuration {
		p.Animations[p.CurrentAnimation].CurrentFrame = (p.Animations[p.CurrentAnimation].CurrentFrame + 1) % len(p.Animations[p.CurrentAnimation].Frames)
		p.Animations[p.CurrentAnimation].ElapsedTicks = 0
	}
}

type PlayerDeathState struct{}

func (s *PlayerDeathState) Update(p *Player, tilemap *Tilemap) {
	p.SetAnimation("death")
	p.Vx = 0
	p.Vy = 0

	p.Oldbb = p.Bb
	p.Bb.X = p.Sprite.X + p.OffsetX
	p.Bb.Y = p.Sprite.Y + p.OffsetY

	p.Collisions["top"] = false
	p.Collisions["bottom"] = false
	p.Collisions["left"] = false
	p.Collisions["right"] = false

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
}
