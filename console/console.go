package console

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/image/font"

	"sync"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/telecoda/pico-go-electron/console/resources/fonts"
	"github.com/telecoda/pico-go-electron/console/resources/images"
)

// Global var

var _console *console

const (
	_version        = "v0.1"
	_logoWidth      = 57
	_logoHeight     = 24
	_spriteWidth    = 8
	_spriteHeight   = 8
	_spritesPerLine = 16
	_maxCmdLen      = 254
	_cursorFlash    = time.Duration(500 * time.Millisecond)
)

const (
	userSpriteBank1 = 0
	userSpriteMask1 = 1
)

type Console interface {
	LoadCart(cart Cartridge) error
	Run() error
	Destroy()
	GetBounds() image.Rectangle
	SetMode(newMode ModeType)
	//Inputter
}

type console struct {
	sync.Mutex
	Config

	currentMode   ModeType
	secondaryMode ModeType
	hasQuit       bool

	// files
	baseDir    string
	currentDir string

	cart Cartridge

	screen *ebiten.Image
	pImage *image.Paletted

	font              font.Face
	sprites           []*image.Paletted
	currentSpriteBank int
	originalPalette   *palette

	//state    Persister
	//recorder Recorder
	//Inputter
}

func NewConsole(consoleType ConsoleType) (Console, error) {

	// validate type
	if _, ok := ConsoleTypes[consoleType]; !ok {
		return nil, fmt.Errorf("Console type: %s not supported", consoleType)
	}

	cfg := NewConfig(consoleType)

	_console = &console{
		Config:        cfg,
		currentMode:   CLI,
		secondaryMode: CODE_EDITOR,
		hasQuit:       false,
	}

	// TODO screen recorder
	//	_console.recorder = NewRecorder(cfg.FPS, cfg.GifLength)
	//	_console.Inputter = NewInputter()

	// init font
	f := bytes.NewReader(fonts.Font_ttf)

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	tt, err := truetype.Parse(b)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 48
	mplusNormalFont := truetype.NewFace(tt, &truetype.Options{
		Size:    6,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	_console.font = mplusNormalFont

	// init sprites
	// There are 2 sprite banks
	// 0 = User sprite bank 1
	// 1 = User sprite bank 1 mask
	_console.sprites = make([]*image.Paletted, 3)

	_console.palette = newPalette(cfg.consoleType)
	_console.originalPalette = newPalette(cfg.consoleType)

	// create paletted image
	rect := image.Rect(0, 0, cfg.ConsoleWidth, cfg.ConsoleHeight)
	pImage := image.NewPaletted(rect, _console.palette.colors)
	_console.pImage = pImage

	// init sprites
	sprites, _, err := image.Decode(bytes.NewReader(images.Sprites_png))
	if err != nil {
		return nil, fmt.Errorf("Error loading sprites: %s", err)
	}
	_console.sprites[userSpriteBank1] = sprites.(*image.Paletted)

	_console.sprites[userSpriteBank1].Palette = _console.palette.colors

	// create a mask
	masks, _, err := image.Decode(bytes.NewReader(images.Sprites_png))
	if err != nil {
		return nil, fmt.Errorf("Error loading sprites: %s", err)
	}
	_console.sprites[userSpriteMask1] = masks.(*image.Paletted)

	// convert all black pixels to zero alpha
	mask := _console.sprites[userSpriteMask1]
	// set palette on mask
	maskPalette := newPalette(PICO8)
	maskPalette.colors[0] = color.RGBA{R: 0, G: 0, B: 0, A: 0}
	mask.Palette = maskPalette.colors

	return _console, nil
}

func (c *console) GetBounds() image.Rectangle {
	return image.Rect(0, 0, 0, 0)
}

func (c *console) SetMode(newMode ModeType) {
	c.Lock()
	defer c.Unlock()
	c.currentMode = newMode
}

func (c *console) LoadCart(cart Cartridge) error {
	c.cart = cart
	return nil
}

var lastFrame time.Time
var startFrame time.Time
var endFrame time.Time

// Run is the main run loop
func (c *console) Run() error {

	// init pixelbuffer
	pb, _ := newPixelBuffer(c.Config)

	c.cart.initPb(pb)

	// poll events
	endFrame = time.Now() // init end frame
	startFrame = time.Now()

	// init the cart
	c.cart.Init()

	return ebiten.Run(c.update, c.Config.ConsoleWidth, c.Config.ConsoleHeight, 1, "pico-go")
}

func (c *console) update(screen *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}

	c.cart.Update()
	c.cart.Render()

	// record frame
	//			c.recorder.AddFrame(mode.GetFrame(), mode)

	//mode.Flip()

	cpb := c.cart.getPb()

	pb := cpb.getPixelBuffer()

	// convert paletted image to RGBA

	pix := make([]uint8, 65536)

	b := 0
	for _, palPix := range pb.pixelSurface.Pix {
		// lookup color
		rgba := c.palette.colorMap[ColorID(palPix)]
		pix[b] = rgba.R
		b++
		pix[b] = rgba.G
		b++
		pix[b] = rgba.B
		b++
		pix[b] = rgba.A
		b++
	}

	//screen.ReplacePixels(pb.pixelSurface.Pix)
	screen.ReplacePixels(pix)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))

	return nil
}

var escapePressed bool

func (c *console) handleInput() error {

	// This method is called every iteration
	// it checks system wide input events such as the escape key

	// 	// TODO keys to implement
	// 	// F7 Capture cartridge label image
	// 	// F8 Start recording a video
	// 	// F9 Save GIF video to desktop (max: 8 seconds by default)

	// 	switch t := event.(type) {
	// 	case *sdl.QuitEvent:
	// 		fmt.Printf("Quit event..\n")
	// 		c.state.SaveState(c)
	// 		c.hasQuit = true
	// 	case *sdl.KeyDownEvent:
	// 		switch t.Keysym.Sym {
	// 		case sdl.K_ESCAPE:
	// 			c.toggleCLI()

	// if not handled pass event to mode event handler
	// if err := c.modes[c.currentMode].HandleInput(); err != nil {
	// 	return err
	// }

	// 		case sdl.K_F6:
	// 			if err := c.saveScreenshot(); err != nil {
	// 				return err
	// 			}
	// 		case sdl.K_F9:
	// 			if err := c.saveVideo(); err != nil {
	// 				return err
	// 			}
	// 		default:
	// 			// pass keydown events to mode handle
	// 			if err := mode.HandleEvent(event); err != nil {
	// 				return err
	// 			}
	// 		}
	// 	case *sdl.MouseButtonEvent:
	// 		// we only care about mouse clicks
	// 		if t.State == 1 && t.Button == 1 {
	// 			c.mouseClicked(t.X, t.Y)
	// 		}
	// 	default:
	// 		// if not handled pass event to mode event handler
	// 		if err := mode.HandleEvent(event); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	//	}

	return nil
}

func (c *console) mouseClicked(x, y int32) {
	// transform window x,y coords to pixel buffer coords

	fmt.Printf("Mouse clicked at x: %d y: %d\n", x, y)

	// get mode
	// if mode, ok := c.modes[c.currentMode]; ok {
	// 	pb := mode.getPixelBuffer()
	// 	fmt.Printf("RenderRect: %#v\n", pb.renderRect)
	// 	fmt.Printf("pixelBuffer: %#v\n", pb.psRect)
	// 	// subtract top left offset
	// 	x -= int32(pb.renderRect.Min.X)
	// 	y -= int32(pb.renderRect.Min.Y)
	// 	fmt.Printf("[adjusted] Mouse clicked at x: %d y: %d\n", x, y)
	// 	// scale to match pixelbuffer
	// 	scale := float32(pb.renderRect.Max.X) / float32(pb.pixelSurface.Bounds().Max.X)
	// 	scaledX := float32(x) / scale
	// 	scaledY := float32(y) / scale
	// 	x = int32(scaledX)
	// 	y = int32(scaledY)
	// 	fmt.Printf("[scaled] Mouse clicked at x: %d y: %d\n", x, y)
	// }

}

func (c *console) Quit() {
	c.hasQuit = true
}

// Destroy cleans up any resources at end
func (c *console) Destroy() {
}

// saveScreenshot - saves a screenshot of current frame
func (c *console) saveScreenshot() error {

	// return c.recorder.SaveScreenshot("out.png", c.Config.ScreenshotScale)

	return nil
}

// saveVideo - saves a video of last x seconds
func (c *console) saveVideo() error {
	// return c.recorder.SaveVideo("out.gif", c.Config.GifScale)
	return nil
}
