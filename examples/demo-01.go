/*
	This is a simple demo project to show you how to use pico-go

	Copyright 2018 @telecoda

*/
package main

import (
	"fmt"
	"log"

	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	// define these vars to be used in javascript canvas scaling code
	screenWidth  = 320
	screenHeight = 240
)

var blue = color.RGBA{R: 33, G: 174, B: 255, A: 255}

// update - this method is called 60 time a second
func update(screen *ebiten.Image) error {

	// this is here to skip frames when things struggle
	if ebiten.IsRunningSlowly() {
		return nil
	}

	// screen screen with Gopher color
	screen.Fill(blue)
	// show some frame rate text to impress your friends
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))
	return nil
}

func main() {
	// run app
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "pico-go (demo project)"); err != nil {
		log.Fatal(err)
	}
}
