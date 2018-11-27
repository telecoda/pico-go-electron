package main

import (
	"github.com/telecoda/pico-go-electron/console"
)

/*
	This is a project to demo the palette manipulation

	Copyright 2018 @telecoda

*/

const (
	// define these vars to be used in javascript canvas scaling code
	screenWidth  = 128
	screenHeight = 128
	consoleType  = console.PICO8
)

// Code must implement console.Cartridge interface

type cartridge struct {
	*console.BaseCartridge

	// example vars below
	mapAnim      bool
	frameCount   int
	totalFrames  int
	currentColor int
}

// Init - called once when cart is initialised
func (c *cartridge) Init() error {
	// the Init method receives a PixelBuffer reference
	// hold onto this reference, this is the display that
	// your code will be drawing onto each frame
	c.frameCount = 0
	c.totalFrames = 25
	c.currentColor = 0
	c.mapAnim = false
	return nil
}

// Update -  called once every frame
func (c *cartridge) Update() {
	c.frameCount++
	if c.frameCount > c.totalFrames {
		// trigger update

		if c.mapAnim {
			c.MapColor(console.ColorID(c.currentColor), console.PICO8_RED)
		} else {
			c.SetTransparent(console.ColorID(c.currentColor), true)
		}
		// reset counters
		c.frameCount = 0
		c.currentColor++
		if c.currentColor > 15 {
			// all colors have been swapped reset
			c.currentColor = 0
			c.PaletteReset()
			c.mapAnim = !c.mapAnim
		}
	}
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.Cls()
	c.RectFill(0, 0, 32, 32, 0)
	c.RectFill(32, 0, 64, 32, 1)
	c.RectFill(64, 0, 96, 32, 2)
	c.RectFill(96, 0, 128, 32, 3)

	c.RectFill(0, 32, 32, 64, 4)
	c.RectFill(32, 32, 64, 64, 5)
	c.RectFill(64, 32, 96, 64, 6)
	c.RectFill(96, 32, 128, 64, 7)

	c.RectFill(0, 64, 32, 96, 8)
	c.RectFill(32, 64, 64, 96, 9)
	c.RectFill(64, 64, 96, 96, 10)
	c.RectFill(96, 64, 128, 96, 11)

	c.RectFill(0, 96, 32, 128, 12)
	c.RectFill(32, 96, 64, 128, 13)
	c.RectFill(64, 96, 96, 128, 14)
	c.RectFill(96, 96, 128, 128, 15)

	c.PrintAt("PALETTE:", 46, 5, 15)
	c.Line(0, 12, 128, 12)

	if c.mapAnim {
		c.PrintAt("COLORS CAN BE SWAPPED.", 20, 20, 15)
	} else {
		c.PrintAt("COLORS CAN BE TRANSPARENT.", 12, 20, 15)
	}
}
