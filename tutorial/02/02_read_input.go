package main

import (
	"image"
	"log"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Starfighter",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	pic, err := loadPicture("./pig.png")
	if err != nil {
		log.Fatal(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())
	spriteX := win.Bounds().Center().X
	spriteY := win.Bounds().Center().Y
	for !win.Closed() {
		win.Clear(colornames.Black)
		if win.Pressed(pixelgl.KeyLeft) {
			spriteX -= 1
		}
		if win.Pressed(pixelgl.KeyRight) {
			spriteX += 1
		}
		sprite.Draw(win, pixel.IM.Moved(pixel.V(spriteX, spriteY)))
		win.Update()
	}

}
