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
}

// Init -  called once
func (c *cartridge) Init() error {

	// override resources with local versions
	err := console.InitSprites(sprites_gif)
	if err != nil {
		panic(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return nil
}

// Update -  called once every frame
func (c *cartridge) Update() {

}

// Render - called once every frame
func (c *cartridge) Render() {
	c.Cls(console.PICO8_BLACK)
	// Render top bar
	c.RectFill(0, 0, 128, 7, console.PICO8_RED)
	// Render edit box
	c.RectFill(0, 8, 128, 86, console.PICO8_DARK_GRAY)
	// Render bottom bar
	c.RectFill(0, 120, 128, 127, console.PICO8_RED)

	// c.SpriteRotated(0, 56, 41, 2, 2, 16, 16, c.rot)

	// c.SpriteFlipped(2, 56, 66, 2, 2, 16, 16, true, false)
	// c.SpriteFlipped(2, 76, 66, 2, 2, 16, 16, false, true)
	// c.SpriteFlipped(2, 96, 66, 2, 2, 16, 16, true, true)

	// c.Sprite(40, 40, 95, 4, 2, 64, 32)

}
