package models

import (
	"github.com/faiface/pixel/pixelgl"
	"sync"
	"time"
)

type Dog struct {
	X, Y  float64
	Speed float64
}

var DogScale float64 = 0.3
var DogMutex sync.Mutex

func UpdateDogPositions(dogs []Dog, playerDogIndex int, win *pixelgl.Window) {
	for {
		DogMutex.Lock()
		if win.JustPressed(pixelgl.KeySpace) {
			dogs[playerDogIndex].X += dogs[playerDogIndex].Speed
		}
		for i := range dogs {
			if i == playerDogIndex {
				continue
			}
			dogs[i].X += dogs[i].Speed
		}
		DogMutex.Unlock()
		time.Sleep(time.Millisecond * 16)
	}
}
