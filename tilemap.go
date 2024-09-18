package main

import (
	"encoding/csv"
	"image"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	Type       string
	Sprite     Rect
	Bb         Rect
	Image      *ebiten.Image
	Colour     string
	Collidable bool
}

type Tilemap struct {
	SizeX int
	SizeY int
	Tiles [][]*Tile
}

func NewTilemap(sizeX, sizeY int) *Tilemap {
	t := &Tilemap{
		SizeX: sizeX,
		SizeY: sizeY,
	}
	t.Tiles = make([][]*Tile, t.SizeY)
	for i := range t.Tiles {
		t.Tiles[i] = make([]*Tile, t.SizeX)
	}
	for y, row := range t.Tiles {
		for x := range row {
			t.Tiles[y][x] = &Tile{
				Image: nil,
			}
		}
	}

	return t
}

func (t *Tilemap) LoadTiles(spritesheet *ebiten.Image, filepath string, gridWidth, tileSize, separation int, colour string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("level load:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal("csv read:", err)
	}

	for y, row := range rows {
		for x, cell := range row {
			if cell == "-1" {
				continue
			}

			num, err := strconv.Atoi(cell)
			if err != nil {
				log.Fatal("cell to num:", err)
			}

			tileRow := num / gridWidth
			tileCol := num % gridWidth

			tileX := tileCol * (tileSize + separation)
			tileY := tileRow * (tileSize + separation)

			img := spritesheet.SubImage(image.Rect(tileX, tileY, tileX+tileSize, tileY+tileSize)).(*ebiten.Image)

			if num == 183 {
				t.Tiles[y][x] = &Tile{
					Type: "spike",
					Sprite: Rect{
						X: float64(x * tileSize),
						Y: float64(y * tileSize),
						W: float64(tileSize),
						H: float64(tileSize),
					},
					Bb: Rect{
						X: float64(x*tileSize) + 3,
						Y: float64(y*tileSize) + 9,
						W: float64(tileSize) - 6,
						H: float64(tileSize) - 7,
					},
					Image:      img,
					Colour:     colour,
					Collidable: true,
				}
			} else if num == 116 {
				t.Tiles[y][x] = &Tile{
					Type: "ledge",
					Sprite: Rect{
						X: float64(x * tileSize),
						Y: float64(y * tileSize),
						W: float64(tileSize),
						H: float64(tileSize),
					},
					Bb: Rect{
						X: float64(x * tileSize),
						Y: float64(y * tileSize),
						W: float64(tileSize),
						H: float64(tileSize) - 11,
					},
					Image:      img,
					Colour:     colour,
					Collidable: true,
				}
			} else {
				t.Tiles[y][x] = &Tile{
					Type: "tile",
					Sprite: Rect{
						X: float64(x * tileSize),
						Y: float64(y * tileSize),
						W: float64(tileSize),
						H: float64(tileSize),
					},
					Bb: Rect{
						X: float64(x * tileSize),
						Y: float64(y * tileSize),
						W: float64(tileSize),
						H: float64(tileSize),
					},
					Image:      img,
					Colour:     colour,
					Collidable: true,
				}
			}
		}
	}
}

func (t *Tilemap) Draw(screen *ebiten.Image, scroll Scroll) {
	for _, tile := range t.TilesVisible(scroll) {
		if tile.Image != nil {
			opts := ebiten.DrawImageOptions{}
			opts.GeoM.Translate(tile.Sprite.X-scroll.X, tile.Sprite.Y-scroll.Y)
			screen.DrawImage(tile.Image, &opts)
		}
	}
}

func (t *Tilemap) TilesAround(x, y float64) []Tile {
	tileX := int(x / tileSize)
	tileY := int(y / tileSize)

	var tilesAround []Tile
	for i := -1; i <= 1; i += 1 {
		for j := -1; j <= 1; j += 1 {
			currentY := tileY + j
			currentX := tileX + i

			if currentX < 0 || currentX >= len(t.Tiles[0]) || currentY < 0 || currentY >= len(t.Tiles) {
				continue
			}

			tile := *t.Tiles[tileY+j][tileX+i]
			if tile.Image != nil && tile.Collidable {
				tilesAround = append(tilesAround, tile)
			}
		}
	}

	return tilesAround
}

func (t *Tilemap) TilesVisible(scroll Scroll) []Tile {
	var tilesVisible []Tile
	for x := int(scroll.X / tileSize); x < int((scroll.X+screenWidth)/tileSize)+1; x++ {
		for y := int(scroll.Y / tileSize); y < int((scroll.Y+screenHeight)/tileSize)+1; y++ {
			if y >= 0 && x >= 0 && y < t.SizeY && x < t.SizeX {
				tile := *t.Tiles[y][x]
				if tile.Image != nil {
					tilesVisible = append(tilesVisible, tile)
				}
			}
		}
	}
	return tilesVisible
}
