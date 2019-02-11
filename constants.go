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

const genMainSrc1 = `/*
	This is the generated bootstrap code for running a pico-go cartridge.

	Don't change any of this code, you shouldn't even really be seeing it, I guess you're the curious type.
*/
`

// we had to split this into 3 const or the //go generate line would trigger when we didn't want it to
const genMainSrc2 = `//go:generate file2byteslice -package=main -input=./sprites.gif -output=./gen-sprites.go -var=sprites_png
`

const genMainSrc3 = `package main

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

var default_sprite_png = []byte("GIF89a\x80\x00\x80\x00\xe3\x10\x00\x00\x00\x00\".S})S\x00\x85Q\xa9R8_WP\xc0\xc1\xc5\xff\xf0\xe7\xff\aN\xff\xa1\b\xfe\xeb,\x00\xe39,\xab\xfe\xffv\xa6\x82u\x9a\xffʨ,\x00\x00\x00\x00\x80\x00\x80\x00\x00\x04\xfe\x10\xc8I\x8b\xb54ˋu\xe5\x1eЅdir(jn\x17\x89\x81\xe1\xdbƅ8\u007fc\xbd\xee\xf4u\xa40\xdaaxp\xa5\x8c\xaa\x1eLg\xeb\xdcxЦe\b\x1cy\n\xc4!2\xa8\x01n[\xb3\xa5\xf5\x13\xbdN\xb3\xe7߸\xf2\xcb\"w\xcf7\x8e\xe5\\K˙4\x1a\xcd\xcc\x13\xb1ZJv\x13G_:a\x88}\x84qe\x80Yj\x8f\x8a6\x8fD\"y\\\x968c\x8aI\x8b\x87\x89f9x\x8e\x80(\u007f]\x94\x959\x85,\x9e\\\xab\x9dMt\x9fL\xb5\xa0Q\xa4j\x1c\xa7\x93z\x8f6\xb3\xaf\xb4\xb1u\xac\xb2RA\x89\xb4\x84d\xa8T@\xa7\xbem\u007fEUq֛ؗ\xc3O\xcat\xdfw\x1b\x1b\u007f\xd0ϥ\x8e\xba\xca\xc6\xc1\xb1\xec\x98Ȍ\xc7b\xcb\xe1`\xe4)\xbc\x8d^\xa1\xed\xae\xef\xee\xa1\xdc\xd9\x1a\xf6\r]\xaa|\xfa\xfe\xf9;ѯK<u\xd9\n\x811\x98j\xd0\n\x06\x16\x18`,\xc0 \xc4F\x8c\xfe\x1aI\x80\xd4ȱ\xe3\x16A\xf1\x04&:h\x0e\x8f\x84\x91\x1aCz\x18YRdƑ\x1e7\xde\xe4w!f\x86\x98:9\x98Ķ\x0e\xcaG\xa0\x1e3f\x14\tt\xe9̠2\u007f\x92\xec\x19\xf5eL\xa1>U6tɵ\xabWM_\xc3\x02\x13K\x96\xab\x1e\x8beK\xb4䁖\xe1\xa2\x13+(\xeaJ\v%R\xdbQ\x92\xe0B\xc1\"w-]\x12\x91\xa8\xb8\xb8d\xb6\x82Yj\x15\x05\xff\r\xb1\xc7\xefb\xb1\xd2(\x91\xba+\x96\xa5\xe2\xc7a\xfbZ\xbe\x8cy\xf3\x0f\xccf\xf9j\xf6\x9c7\xad\xe7ϠqQD7\x1a\rh\xa2)S\x9by\xc6\avQٸs\xeb\xdeͻ\xb7\xef\xdf\xc0\x83\v\x1fN\xbc\xb8\xf1\xe3ȓ+_μ\xb9\xf3\xe7УK\x9fN\xbd\xba\xf5\xebسk\xdfν\xbb\xf7\xef\xe0Ë\x1fO\xbe\xbc\xf9\xf3\xe8ӫ_Ͼ\xbd\xfb\xf7\xf0\xe3˟O\xbf\xbe\xfd\xfb\xf8\xf3\xeb\xdfϿ\xbf\xff\xff\x00\x06H(\xe0\x80\x04\x16h\xe0\x81\b&\xa8\xe0\x82\f6\xe8\xe0\x83\x10F(\xe1\x84\x14Vh\xe1\x85\x18f\xa8\xe1\x86\x1cv\xe8\xe1\x87 \x86(\xe2\x88$\x96h\xe2\x89(\xa6\xa8\xe2\x8a,\xb6\xe8\xe2\x8b0\xc6(\xe3\x8c4\xd6h\xe3\x8d8\xba\x14\x01\x00;")
