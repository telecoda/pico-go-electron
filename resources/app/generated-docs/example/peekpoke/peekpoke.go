package main

/*
	This is a project to demo the peek poke functions

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
)

type cartridge struct {
	*console.BaseCartridge
}

// Init - called once when cart is initialised
func (c *cartridge) Init() error {
	c.Cls(console.PICO8_BLACK)
	for i := 0; i < 128; i++ {
		c.Line(0, i, 128, i, console.ColorID(i)%16)
	}
	return nil
}

// Update -  called once every frame
func (c *cartridge) Update() {
}

// Render - called once every frame
func (c *cartridge) Render() {

	c.PrintAt("PEEK POKE:", 40, 5, console.PICO8_WHITE)

	for i := 0; i < 128*128; i++ {
		pos0Value := c.Peek(i)
		pos0Value++
		c.Poke(i, pos0Value)
	}
}
