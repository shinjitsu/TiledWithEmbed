package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"log"
	"path"
)

//go:embed assets/*
var EmbeddedAssets embed.FS

type mapGame struct {
	Level    *tiled.Map
	tileHash map[uint32]*ebiten.Image
}

func (m mapGame) Update() error {
	return nil
}

func (game mapGame) Draw(screen *ebiten.Image) {
	drawOptions := ebiten.DrawImageOptions{}
	for tileY := 0; tileY < game.Level.Height; tileY += 1 {
		for tileX := 0; tileX < game.Level.Width; tileX += 1 {
			drawOptions.GeoM.Reset()
			TileXpos := float64(game.Level.TileWidth * tileX)
			TileYpos := float64(game.Level.TileHeight * tileY)
			drawOptions.GeoM.Translate(TileXpos, TileYpos)
			tileToDraw := game.Level.Layers[0].Tiles[tileY*game.Level.Width+tileX]
			ebitenTileToDraw := game.tileHash[tileToDraw.ID]
			screen.DrawImage(ebitenTileToDraw, &drawOptions)
		}
	}
}

func (m mapGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	gameMap := loadMapFromEmbedded(path.Join("assets", "demoMap.tmx"))
	ebiten.SetWindowSize(gameMap.TileWidth*gameMap.Width, gameMap.TileHeight*gameMap.Height)
	ebiten.SetWindowTitle("Maps Embedded")
	ebitenImageMap := makeEbiteImagesFromMap(*gameMap)
	oneLevelGame := mapGame{
		Level:    gameMap,
		tileHash: ebitenImageMap,
	}
	err := ebiten.RunGame(&oneLevelGame)
	if err != nil {
		fmt.Println("Couldn't run game:", err)
	}
}

func makeEbiteImagesFromMap(tiledMap tiled.Map) map[uint32]*ebiten.Image {
	idToImage := make(map[uint32]*ebiten.Image)
	for _, tile := range tiledMap.Tilesets[0].Tiles {
		embeddedFile, err := EmbeddedAssets.Open(path.Join("assets", tile.Image.Source))
		if err != nil {
			log.Fatal("failed to load embedded image ", embeddedFile, err)
		}
		ebitenImageTile, _, err := ebitenutil.NewImageFromReader(embeddedFile)
		if err != nil {
			fmt.Println("Error loading tile image:", tile.Image.Source, err)
		}
		idToImage[tile.ID] = ebitenImageTile
	}
	return idToImage
}

func loadMapFromEmbedded(name string) *tiled.Map {
	embeddedMap, err := tiled.LoadFile(name, tiled.WithFileSystem(EmbeddedAssets))
	if err != nil {
		fmt.Println("Error loading embedded map:", err)
	}
	return embeddedMap
}
