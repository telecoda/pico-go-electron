package main

/*
	This is a project to demo the sprite functions

	Copyright 2018 @telecoda

*/

import (
	"github.com/telecoda/pico-go-electron/console"
)

const (
	// define these vars to be used in javascript canvas scaling code
	screenWidth  = 128
	screenHeight = 128
	consoleType  = console.PICO8

	spritesTop = 88
)

type cartridge struct {
	*console.BaseCartridge

	selectedSprite int
	selectedScale  int
}

// Init -  called once
func (c *cartridge) Init() error {
	c.selectedSprite = 0
	c.selectedScale = 1
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
	// Main sprite box
	c.RectFill(7, 11, 73, 76, console.PICO8_BLACK)
	// Palette box
	c.RectFill(79, 11, 121, 48, console.PICO8_BLACK)
	// Colours
	colour := 0
	left := 80
	top := 12
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			c.RectFill(left+(col*10), top, left+(col*10)+10, top+8, console.ColorID(colour))
			colour += 1
		}
		top += 9
	}

	// render all sprites at bottom
	sprite := 0
	left = 0

	c.Sprite(sprite, left, spritesTop, 16, 2, 128, 16)
	c.Sprite(sprite+32, left, spritesTop+16, 16, 2, 128, 16)

	// Render bottom bar
	c.RectFill(0, 120, 128, 127, console.PICO8_RED)

	c.renderSelectSprite()

	// c.SpriteRotated(0, 56, 41, 2, 2, 16, 16, c.rot)

	// c.SpriteFlipped(2, 56, 66, 2, 2, 16, 16, true, false)
	// c.SpriteFlipped(2, 76, 66, 2, 2, 16, 16, false, true)
	// c.SpriteFlipped(2, 96, 66, 2, 2, 16, 16, true, true)

}

func (c *cartridge) renderSelectSprite() {
	// draw box around selected sprite

	// convert spriteNumber into top left coords
	// convert sprite number into x,y pos
	xCell := c.selectedSprite % 16
	yCell := (c.selectedSprite - xCell) / 16

	xPos := xCell * 8
	yPos := yCell * 8

	top := spritesTop + yPos - 1
	left := xPos - 1
	bottom := (top + c.selectedScale*8) + 1
	right := (left + c.selectedScale*8) + 1
	c.Rect(left, top, right, bottom, console.PICO8_WHITE)

	// Render sprite on editor
	// fit into 64 x 64 box
	c.Sprite(c.selectedSprite, 8, 12, c.selectedScale, c.selectedScale, 64, 64)

}
