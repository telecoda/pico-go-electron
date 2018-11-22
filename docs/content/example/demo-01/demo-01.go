package main

/*
	This is a simple demo project to show you how to use pico-go

	Copyright 2018 @telecoda

*/

import "github.com/telecoda/pico-go-electron/console"

const (
	// set console type to one of the predefined consoles
	consoleType = console.PICO8
)

type cartridge struct {
	*console.BaseCartridge
}

// Init -  called once
func (c *cartridge) Init() error {
	console.ShowFPS()
	return nil
}

// Update -  called once every frame
func (c *cartridge) Update() {
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.ClsWithColor(console.PICO8_BLUE)
	c.PrintAt("Hello", 10, 20)
}
