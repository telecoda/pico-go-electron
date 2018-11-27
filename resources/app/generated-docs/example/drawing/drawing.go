package main

/*
	This is a project to demo the drawing primitives

	Copyright 2018 @telecoda

*/

import (
	"github.com/telecoda/pico-go-electron/console"
)

const (
	// set console type to one of the predefined consoles
	consoleType = console.PICO8
)

type cartridge struct {
	*console.BaseCartridge
}

// Init - called once when cart is initialised
func (c *cartridge) Init() error {
	return nil
}

// Update -  called once every frame
func (c *cartridge) Update() {
}

// Render - called once every frame
func (c *cartridge) Render() {

	c.Cls(console.PICO8_BLACK)
	c.PrintAt("DRAWING:", 50, 5, console.PICO8_WHITE)

	c.Line(0, 12, 128, 12, console.PICO8_WHITE)
	c.PrintAt("RECTS:", 10, 32, console.PICO8_WHITE)
	c.Rect(45, 30, 55, 40)
	c.SetColor(console.PICO8_GREEN)
	c.RectFill(65, 30, 75, 40)
	c.RectFill(85, 25, 105, 45, console.PICO8_RED)
	c.PrintAt("CIRCLE:", 10, 55, console.PICO8_WHITE)
	c.Circle(50, 57, 5)
	c.SetColor(console.PICO8_BLUE)
	c.CircleFill(70, 57, 5)
	c.CircleFill(95, 57, 10, console.PICO8_BROWN)
	c.PrintAt("LINES:", 10, 77, console.PICO8_WHITE)
	c.SetColor(console.PICO8_LIGHT_GRAY)
	c.Line(45, 77, 105, 77)
	c.Line(45, 79, 105, 79, console.PICO8_YELLOW)
	c.PrintAt("POINTS:", 10, 99, console.PICO8_WHITE)
	c.PSet(50, 99)
	c.PSet(70, 99, console.PICO8_PINK)
	// get color of point // earlier rect
	pointColor := c.PGet(85, 25)
	c.PSet(95, 99, pointColor)
}
