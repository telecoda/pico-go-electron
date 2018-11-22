package main

const (
	defaultCodeDir    = "gosrc"
	defaultSourceFile = "main.go"
	defaultCompileDir = "js"
	defaultOutputFile = "cart.js"
	ebitenRepo        = "github.com/hajimehoshi/ebiten"
)

const aboutBody = "Welcome on to `pico-go`\n\nThe golang fantasy console.\n\nby @telecoda\n"

const demoSrc = `package main

/*
	This is a simple demo project to show you how to use pico-go
	Copyright 2018 @telecoda
*/

import "github.com/telecoda/pico-go-electron/console"

const (
	// set console type to one of the predefined consoles
	consoleType = console.PICO8
	// define these vars to be used in javascript canvas scaling code
	screenWidth  = 128
	screenHeight = 128
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
}`

const genMainSrc = `/*
	This is the generated bootstrap code for running a pico-go cartridge.

	Don't change any of this code, you shouldn't even really be seeing it, I guess you're the curious type.
*/
package main

import (
	"fmt"

	"github.com/telecoda/pico-go-electron/console"
)

func NewCart() console.Cartridge {
	return &cartridge{
		BaseCartridge: console.NewBaseCart(),
	}
}

func main() {
	cart := NewCart()
	if err := console.Init(consoleType); err != nil {
		fmt.Printf("Failed to init console: %s\n", err)
		return
	}
	if err := console.Run(cart); err != nil {
		fmt.Printf("Failed to run cartridge: %s\n", err)
		return
	}
}`
