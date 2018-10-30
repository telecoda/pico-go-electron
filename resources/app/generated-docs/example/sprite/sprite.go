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
)

type cartridge struct {
	*console.BaseCartridge

	// example vars below
	running bool
	rot     float64
}

// NewCart - initialise a struct implementing Cartridge interface
func NewCart() console.Cartridge {
	return &cartridge{
		BaseCartridge: console.NewBaseCart(),
	}
}

// Init -  called once
func (c *cartridge) Init() {
}

// Update -  called once every frame
func (c *cartridge) Update() {
	c.rot -= 4
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.ClsWithColor(console.PICO8_WHITE)
	//c.SetTransparent(5, true)
	//c.MapColor(console.PICO8_BLUE, console.PICO8_WHITE)
	// for i := 20; i < 40; i++ {
	// 	c.LineWithColor(0, i, 128, i, console.PICO8_DARK_GRAY)
	// }

	//c.MapColor(console.PICO8_DARK_GRAY, console.PICO8_WHITE)
	c.PrintAtWithColor("SPRITES:", 50, 5, console.PICO8_BLACK)
	c.Line(0, 12, 128, 12)
	c.PrintAtWithColor("FLIPX: false", 10, 18, console.PICO8_BLACK)
	c.PrintAtWithColor("FLIPY: false", 10, 26, console.PICO8_BLACK)
	c.PrintAtWithColor(fmt.Sprintf("R: %d", int(c.rot)), 100, 22, console.PICO8_BLACK)
	c.Sprite(0, 70, 16, 2, 2, 16, 16, c.rot, false, false)
	c.PrintAtWithColor("FLIPX: true", 10, 38, console.PICO8_BLACK)
	c.PrintAtWithColor("FLIPY: false", 10, 46, console.PICO8_BLACK)
	c.Sprite(2, 70, 36, 2, 2, 16, 16, 0, true, false)
	c.PrintAtWithColor("FLIPX: false", 10, 58, console.PICO8_BLACK)
	c.PrintAtWithColor("FLIPY: true", 10, 66, console.PICO8_BLACK)
	c.Sprite(4, 70, 56, 2, 2, 16, 16, 0, false, true)
	c.PrintAtWithColor("FLIPX: true", 10, 78, console.PICO8_BLACK)
	c.PrintAtWithColor("FLIPY: true", 10, 86, console.PICO8_BLACK)
	c.Sprite(0, 70, 76, 2, 2, 16, 16, 0, true, true)
	c.PrintAtWithColor("Scaled:", 10, 98, console.PICO8_BLACK)
	c.Sprite(40, 40, 96, 4, 2, 64, 32, 0, false, false)

}

func main() {

	// Create virtual console - based on cart config
	con, err := console.NewConsole(console.PICO8)
	if err != nil {
		panic(err)
	}
	defer con.Destroy()

	cart := NewCart()

	if err := con.LoadCart(cart); err != nil {
		panic(err)
	}

	if err := con.Run(); err != nil {
		panic(err)
	}
}
