package console

import (
	"image"
	"image/color"
)

type ColorID uint8

/*
	This package tries to replicate the pico8 API as closely as possible
	During development I will be trying to implement more an more of the API
	To achieve feature parity with pico8
	Documented extensively here http://pico-8.wikia.com/wiki/Category:API
*/

type PicoGraphicsAPI interface {
	Clearer
	Drawer
	Paletter
	Printer
	Spriter
}

type PicoInputAPI interface {
	Btn(id int) bool
}

type Clearer interface {
	Cls()                         // Clear screen
	ClsWithColor(colorID ColorID) // Clear screen with color
}

type Drawer interface {
	Color(colorID ColorID) // Set drawing color (colour!!!)
	// drawing primitives
	Circle(x, y, r int)
	CircleWithColor(x, y, r int, colorID ColorID)
	CircleFill(x, y, r int)
	CircleFillWithColor(x, y, r int, colorID ColorID)
	Line(x0, y0, x1, y1 int)
	LineWithColor(x0, y0, x1, y1 int, colorID ColorID)
	PGet(x, y int) ColorID
	PSet(x, y int)
	PSetWithColor(x, y int, colorID ColorID)
	Rect(x0, y0, x1, y1 int)
	RectWithColor(x0, y0, x1, y1 int, colorID ColorID)
	RectFill(x0, y0, x1, y1 int)
	RectFillWithColor(x0, y0, x1, y1 int, colorID ColorID)
}

type Paletter interface {
	PaletteReset()
	PaletteCopy() Paletter
	GetColorID(rgba rgba) ColorID
	GetColor(colorID ColorID) color.Color
	GetRGBA(color ColorID) (rgba, uint32)
	GetColors() []color.Color
	MapColor(fromColor ColorID, toColor ColorID) error
	SetTransparent(color ColorID, enabled bool) error
}

type Printer interface {
	// Text/Printing
	Cursor(x, y int) // Set text cursor
	GetCursor() pos
	Print(str string)                                       // Print a string of characters to the screen at default pos
	PrintAt(str string, x, y int)                           // Print a string of characters to the screen at position
	PrintAtWithColor(str string, x, y int, colorID ColorID) // Print a string of characters to the screen at position with color
	ScrollUpLine()
}

type Spriter interface {
	Sprite(n, x, y, w, h, dw, dh int)
	SpriteFlipped(n, x, y, w, h, dw, dh int, flipX, flipY bool)
	SpriteRotated(n, x, y, w, h, dw, dh int, rot float64)
}

type ConsoleType string

const (
	PICO8       = "pico8"
	TIC80       = "tic80"
	ZX_SPECTRUM = "zxspectrum"
	CBM64       = "cbm64"
)

var ConsoleTypes = map[ConsoleType]string{
	PICO8:       "PICO8",
	TIC80:       "TIC80",
	ZX_SPECTRUM: "ZX_SPECTRUM",
	CBM64:       "CBM64",
}

const TOTAL_COLORS = 16

type Configger interface {
	GetConfig() Config
}

type Cartridge interface {
	// BaseCartridge methods already implemented
	Configger
	initPb(pb PixelBuffer)
	getPb() PixelBuffer
	PicoInputAPI
	// User implemented methods below
	Init() error
	Render()
	Update()
}

type Runtime interface {
	PicoInputAPI
	LoadCart(cart Cartridge) error
}

type PixelBuffer interface {
	Flip() error // Copy graphics buffer to screen
	Destroy()
	GetFrame() *image.Paletted
	PicoGraphicsAPI
	getPixelBuffer() *pixelBuffer
	GetWidth() int
	GetHeight() int
}

var title = "pico-go virtual games console"

type size struct {
	width  int
	height int
}
