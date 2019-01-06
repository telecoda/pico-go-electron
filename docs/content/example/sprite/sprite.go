package main

/*
	This is a project to demo the sprite functions

	Copyright 2018 @telecoda

*/

import (
	"fmt"

	"github.com/telecoda/pico-go-electron/console"
)

const (
	// define these vars to be used in javascript canvas scaling code
	screenWidth  = 128
	screenHeight = 128
	consoleType  = console.PICO8
)

type cartridge struct {
	*console.BaseCartridge

	// example vars below
	running bool
	rot     int
	barY    int
	scaleY  float64
	scaleX  float64
}

// Init -  called once
func (c *cartridge) Init() error {
	console.ShowFPS()
	return nil
}

// Update -  called once every frame
func (c *cartridge) Update() {
	c.rot += 4
	if c.rot > 360 {
		c.rot = 0
	}
	c.barY += 1
	if c.barY > 128 {
		c.barY = 0
	}
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.Cls(console.PICO8_WHITE)
	c.RectFill(0, c.barY, 128, c.barY+48, console.PICO8_LIGHT_GRAY)
	c.MapColor(console.PICO8_BLUE, console.PICO8_WHITE)
	c.PrintAt("SPRITES:", 50, 5, console.PICO8_BLACK)
	c.Line(0, 12, 128, 12)
	c.PrintAt("SPRITE:", 10, 20, console.PICO8_BLACK)
	c.Sprite(0, 56, 16, 2, 2, 16, 16)

	c.PrintAt("ROTATED:", 10, 45, console.PICO8_BLACK)
	c.PrintAt(fmt.Sprintf("%d", int(c.rot)), 80, 45, console.PICO8_BLACK)
	c.SpriteRotated(0, 56, 41, 2, 2, 16, 16, c.rot)

	c.PrintAt("FLIPPED:", 10, 70, console.PICO8_BLACK)
	c.SpriteFlipped(2, 56, 66, 2, 2, 16, 16, true, false)
	c.PrintAt("X", 62, 83, console.PICO8_BLACK)
	c.SpriteFlipped(2, 76, 66, 2, 2, 16, 16, false, true)
	c.PrintAt("Y", 82, 83, console.PICO8_BLACK)
	c.SpriteFlipped(2, 96, 66, 2, 2, 16, 16, true, true)
	c.PrintAt("XY", 100, 83, console.PICO8_BLACK)

	c.PrintAt("SCALED:", 10, 95, console.PICO8_BLACK)
	c.Sprite(40, 40, 95, 4, 2, 64, 32)

}
