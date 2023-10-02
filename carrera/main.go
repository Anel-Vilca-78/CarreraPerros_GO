package main

import (
	"carrera/scenes"
	"github.com/faiface/pixel/pixelgl"
	_ "image/png"
)

func main() {
	pixelgl.Run(scenes.Run)
}
