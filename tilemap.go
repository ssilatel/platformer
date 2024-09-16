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
	Collidable bool
}

type Tilemap struct {
	Tiles [][]Tile
}

func NewTilemap(sizeX, sizeY int) *Tilemap {
	t := &Tilemap{}
	t.Tiles = make([][]Tile, sizeY)
	for i := range t.Tiles {
		t.Tiles[i] = make([]Tile, sizeX)
	}
	for y, row := range t.Tiles {
		for x := range row {
			t.Tiles[y][x] = Tile{
				Image: nil,
			}
		}
	}

	return t
}

func (t *Tilemap) LoadTiles(spritesheet *ebiten.Image, filepath string, gridWidth, tileSize, separation int) {
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
				t.Tiles[y][x].Image = nil
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
				t.Tiles[y][x] = Tile{
					Type: "spike",
					Sprite: Rect{
						X: float64(x * tileSize),
						Y: float64(y * tileSize),
						W: float64(tileSize),
						H: float64(tileSize),
					},
					Bb: Rect{
						X: float64(x*tileSize) + 3,
						Y: float64(y*tileSize) + 7,
						W: float64(tileSize) - 3,
						H: float64(tileSize) - 7,
					},
					Image:      img,
					Collidable: true,
				}
			} else if num == 116 {
				t.Tiles[y][x] = Tile{
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
					Collidable: true,
				}
			} else {
				t.Tiles[y][x] = Tile{
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
					Collidable: true,
				}
			}
		}
	}
}

func (t *Tilemap) Draw(screen *ebiten.Image, scroll Scroll) {
	for _, row := range t.Tiles {
		for _, tile := range row {
			if tile.Image != nil {
				opts := ebiten.DrawImageOptions{}
				opts.GeoM.Translate(tile.Sprite.X-scroll.X, tile.Sprite.Y-scroll.Y)
				screen.DrawImage(tile.Image, &opts)
			}
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

			tile := t.Tiles[tileY+j][tileX+i]
			if tile.Image != nil && tile.Collidable {
				tilesAround = append(tilesAround, tile)
			}
		}
	}

	return tilesAround
}
