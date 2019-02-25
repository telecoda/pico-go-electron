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

	// UI positions
	mainSpriteTop    = 12
	mainSpriteBottom = 75
	mainSpriteLeft   = 8
	mainSpriteRight  = 72

	colourWidth   = 10
	colourHeight  = 9
	coloursTop    = 12
	coloursLeft   = 80
	coloursRight  = coloursLeft + 4*colourWidth
	coloursBottom = coloursTop + 4*colourHeight

	spritesTop    = 88
	spritesBottom = 88 + 31
	spritesLeft   = 0
	spritesRight  = 127
)

type cartridge struct {
	*console.BaseCartridge

	selectedSprite int
	selectedScale  int
	selectedColour console.ColorID
}

// Init -  called once
func (c *cartridge) Init() error {
	c.selectedSprite = 0
	c.selectedScale = 1
	c.selectedColour = 5
	return nil
}

// Update -  called once every frame
func (c *cartridge) Update() {
	if c.MouseClicked() {
		// check where on the screen click was?
		x, y := c.MousePosition()

		// inside main sprite?
		if x >= mainSpriteLeft && x <= mainSpriteRight && y >= mainSpriteTop && y <= mainSpriteBottom {
			x -= mainSpriteLeft
			y -= mainSpriteTop
			c.spriteClick(x, y)
			return
		}

		// inside colour picker?
		if x >= coloursLeft && x <= coloursRight && y >= coloursTop && y <= coloursBottom {
			x -= coloursLeft
			y -= coloursTop
			c.colourPick(x, y)
			return
		}

		// inside sprite picker?
		if x >= spritesLeft && x <= spritesRight && y >= spritesTop && y <= spritesBottom {
			x -= spritesLeft
			y -= spritesTop
			c.spritePick(x, y)
			return
		}

		// other?

		fmt.Printf("Mouse was clicked at %d, %d\n", x, y)
	}
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.Cls(console.PICO8_BLACK)
	// Render top bar
	c.RectFill(0, 0, 128, 7, console.PICO8_BLUE)
	// Render edit box
	c.RectFill(0, 8, 128, 86, console.PICO8_DARK_GRAY)
	// Main sprite box
	c.RectFill(mainSpriteLeft-1, mainSpriteTop-1, mainSpriteRight+1, mainSpriteBottom+1, console.PICO8_BLACK)
	// Palette box
	c.RectFill(coloursLeft-1, coloursTop-1, coloursRight+1, coloursBottom+1, console.PICO8_BLACK)
	// Colours
	colour := 0
	top := coloursTop
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			c.RectFill(coloursLeft+(col*colourWidth), top, coloursLeft+(col*colourWidth)+colourWidth, top+colourHeight, console.ColorID(colour))
			colour += 1
		}
		top += colourHeight
	}

	// render all sprites at bottom
	sprite := 0

	c.Sprite(sprite, 0, spritesTop, 16, 2, 128, 16)
	c.Sprite(sprite+32, 0, spritesTop+16, 16, 2, 128, 16)

	// Render bottom bar
	c.RectFill(0, 120, 128, 127, console.PICO8_BLUE)

	c.renderSelectedSprite()
	c.renderSelectedColour()
}

func (c *cartridge) colourPick(x, y int) {
	// convert x,y pos into selected colour
	xCell := x / colourWidth
	yCell := y / colourHeight

	c.selectedColour = console.ColorID(xCell + yCell*4)
}

func (c *cartridge) spriteClick(x, y int) {
	fmt.Printf("Sprite was clicked at %d, %d\n", x, y)
}

func (c *cartridge) spritePick(x, y int) {
	// convert x,y pos into selected sprite
	xCell := x / 8
	yCell := y / 8

	c.selectedSprite = xCell + (yCell * 16)

	fmt.Printf("Sprite was picked at %d, %d\n", x, y)
}

func (c *cartridge) renderSelectedSprite() {
	// draw box around selected sprite

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

func (c *cartridge) renderSelectedColour() {
	// draw box around selected colour

	// convert colour number into x,y pos
	xCell := int(c.selectedColour) % 4
	yCell := (int(c.selectedColour) - xCell) / 4

	xPos := xCell * colourWidth
	yPos := yCell * colourHeight

	top := coloursTop + yPos - 1
	left := coloursLeft + xPos - 1
	bottom := (top + colourHeight) + 1
	right := (left + colourWidth) + 1
	c.Rect(left, top, right, bottom, console.PICO8_WHITE)

}
