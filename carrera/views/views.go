package views

import (
	"github.com/faiface/pixel"
	"image"
	_ "image/png"
	"os"
	"sync"
)

func LoadPicture(path string) (pixel.Picture, error) {
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

func LoadRabbitPics(wg *sync.WaitGroup, rabbitPics []string, rabbitSprites []*pixel.Sprite) {
	defer wg.Done()
	for i, pic := range rabbitPics {
		img, _ := LoadPicture(pic)
		rabbitSprites[i] = pixel.NewSprite(img, img.Bounds())
	}
}

func LoadDogPics(wg *sync.WaitGroup, i int, pics []string, dogSprites [][]*pixel.Sprite) {
	defer wg.Done()
	for _, pic := range pics {
		img, _ := LoadPicture(pic)
		dogSprites[i] = append(dogSprites[i], pixel.NewSprite(img, img.Bounds()))
	}
}
