package scenes

import (
	"carrera/models"
	"carrera/views"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type GameResources struct {
	bgSprite      *pixel.Sprite
	rabbitSprites []*pixel.Sprite
	dogSprites    [][]*pixel.Sprite
	atlas         *text.Atlas
	bgScale       float64
}

func Run() {
	resources := loadResources()
	win := setupWindow()
	initializeGameVariables(win, resources)
}

func setupWindow() *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:  "Carrera de perritos",
		Bounds: pixel.R(0, 0, 1200, 800),
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	return win
}

func loadResources() *GameResources {
	var wg sync.WaitGroup

	bgPic, _ := views.LoadPicture("Assets/background.png")
	bgSprite := pixel.NewSprite(bgPic, bgPic.Bounds())

	rabbitPics := []string{"Assets/rabbit1.png", "Assets/rabbit2.png"}
	rabbitSprites := make([]*pixel.Sprite, 2)

	allDogPics := [][]string{
		{"Assets/dog1.png", "Assets/dog12.png"},
		{"Assets/dog2.png", "Assets/dog22.png"},
		{"Assets/dog3.png", "Assets/dog32.png"},
		{"Assets/dog4.png", "Assets/dog42.png"},
	}
	dogSprites := make([][]*pixel.Sprite, 4)

	wg.Add(1)
	go views.LoadRabbitPics(&wg, rabbitPics, rabbitSprites)

	for i, pics := range allDogPics {
		wg.Add(1)
		go views.LoadDogPics(&wg, i, pics, dogSprites)
	}

	wg.Wait()

	rand.Seed(time.Now().UnixNano())
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	return &GameResources{
		bgSprite:      bgSprite,
		rabbitSprites: rabbitSprites,
		dogSprites:    dogSprites,
		atlas:         atlas,
		bgScale:       2.0,
	}
}

func initializeGameVariables(win *pixelgl.Window, resources *GameResources) {
	frame := 0
	toggle := 0
	isGameFinished := false
	winner := -1

	for !win.Closed() {
		dogs := []models.Dog{
			{X: 0, Y: 100, Speed: 60},
			{X: 0, Y: 300, Speed: rand.Float64()*5 + 5},
			{X: 0, Y: 500, Speed: rand.Float64()*5 + 5},
			{X: 0, Y: 650, Speed: rand.Float64()*5 + 5},
		}

		playerDogIndex := 0
		isGameFinished = false
		winner = -1

		go models.UpdateDogPositions(dogs, playerDogIndex, win)

		runGameLoop(win, resources, &frame, &toggle, &isGameFinished, &winner, dogs, playerDogIndex)
	}
}

func runGameLoop(win *pixelgl.Window, resources *GameResources, frame *int, toggle *int, isGameFinished *bool, winner *int, dogs []models.Dog, playerDogIndex int) {
	for !*isGameFinished && !win.Closed() {
		(*frame)++

		drawScene(win, resources, frame, toggle, dogs)

		handleGameLogic(win, resources, isGameFinished, winner, dogs, playerDogIndex)
		win.Update()
		time.Sleep(time.Millisecond * 16)
	}

	handleGameEnd(win, resources, isGameFinished, winner)
}

func drawScene(win *pixelgl.Window, resources *GameResources, frame *int, toggle *int, dogs []models.Dog) {
	win.Clear(colornames.Greenyellow)
	resources.bgSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, resources.bgScale).Moved(pixel.V(600, 400)))

	if *frame%20 == 0 {
		*toggle = 1 - *toggle
	}

	models.DogMutex.Lock()
	for i, dog := range dogs {
		resources.dogSprites[i][*toggle].Draw(win, pixel.IM.Scaled(pixel.ZV, models.DogScale).Moved(pixel.V(dog.X, dog.Y)))
	}
	models.DogMutex.Unlock()

	resources.rabbitSprites[*toggle].Draw(win, pixel.IM.Scaled(pixel.ZV, models.RabbitScale).Moved(pixel.V(1100, 750)))

	txt := text.New(pixel.V(50, 750), resources.atlas)
	txt.Color = colornames.Black
	fmt.Fprintf(txt, "Presiona espacio para mover a tu perro")
	txt.Draw(win, pixel.IM.Scaled(txt.Orig, 3))
}

func handleGameLogic(win *pixelgl.Window, resources *GameResources, isGameFinished *bool, winner *int, dogs []models.Dog, playerDogIndex int) {
	models.DogMutex.Lock()
	for i, dog := range dogs {
		if dog.X >= 1150 && *winner == -1 {
			*isGameFinished = true
			*winner = i
		}
	}
	models.DogMutex.Unlock()

	if *winner != -1 {
		win.Clear(colornames.Greenyellow)
		txt := text.New(pixel.V(100, 400), resources.atlas)
		if *winner == playerDogIndex {
			txt.Color = colornames.Blue
			fmt.Fprintf(txt, "Ganaste! reinicia el juego con la tecla R")
		} else {
			txt.Color = colornames.Red
			fmt.Fprintf(txt, "Oh no, parece que has perdido! reinicia el juego con la tecla R")
		}
		txt.Draw(win, pixel.IM.Scaled(txt.Orig, 2))
	}
}

func handleGameEnd(win *pixelgl.Window, resources *GameResources, isGameFinished *bool, winner *int) {
	for *isGameFinished && !win.Closed() {
		if win.JustPressed(pixelgl.KeyR) {
			break
		}
		win.Update()
		time.Sleep(time.Millisecond * 16)
	}
}
