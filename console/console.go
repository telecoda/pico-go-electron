package console

import (
	"bytes"
	"fmt"
	"image"
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
	//_charWidth      = 4
	//_charHeight     = 6
	_maxCmdLen   = 254
	_cursorFlash = time.Duration(500 * time.Millisecond)
)

const (
	systemSpriteBank = 0
	userSpriteBank1  = 1
	userSpriteMask1  = 2
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
	//modes         map[ModeType]Mode
	hasQuit bool

	// files
	baseDir    string
	currentDir string

	cart Cartridge

	screen *ebiten.Image
	pImage *image.Paletted

	font              font.Face
	logo              *image.RGBA
	sprites           []*image.RGBA
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

	// f, err := ebitenutil.OpenFile(fontPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

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

	// init logo
	img, _, err := image.Decode(bytes.NewReader(images.Logo_png))
	if err != nil {
		return nil, fmt.Errorf("Error loading image: %s", err)
	}
	_console.logo = img.(*image.RGBA)

	// init sprites
	// There are 3 sprite banks
	// 0 = System sprites
	// 1 = User sprite bank 1
	// 2 = User sprite bank 1 mask
	_console.sprites = make([]*image.RGBA, 3)

	_console.palette = newPalette(cfg.consoleType)
	_console.originalPalette = newPalette(cfg.consoleType)

	// create paletted image
	rect := image.Rect(0, 0, cfg.ConsoleWidth, cfg.ConsoleHeight)
	pImage := image.NewPaletted(rect, _console.palette.colors)
	_console.pImage = pImage

	// init icons
	icons, _, err := image.Decode(bytes.NewReader(images.Icons_png))
	if err != nil {
		return nil, fmt.Errorf("Error loading icons: %s", err)
	}
	_console.sprites[systemSpriteBank] = icons.(*image.RGBA)

	// init sprites
	sprites, _, err := image.Decode(bytes.NewReader(images.Sprites_png))
	if err != nil {
		return nil, fmt.Errorf("Error loading sprites: %s", err)
	}
	_console.sprites[userSpriteBank1] = sprites.(*image.RGBA)

	// create a mask
	masks, _, err := image.Decode(bytes.NewReader(images.Sprites_png))
	if err != nil {
		return nil, fmt.Errorf("Error loading sprites: %s", err)
	}
	_console.sprites[userSpriteMask1] = masks.(*image.RGBA)

	// convert all black pixels to zero alpha
	mask := _console.sprites[userSpriteMask1]
	for x := 0; x < mask.Bounds().Dx(); x++ {
		for y := 0; y < mask.Bounds().Dy(); y++ {
			c := mask.RGBAAt(x, y)
			if c.R == 0 && c.G == 0 && c.B == 0 {
				c.A = 0
				mask.SetRGBA(x, y, c)
			} else {
				c.A = 255
				c.R = 255
				c.G = 255
				c.B = 255
				mask.SetRGBA(x, y, c)
			}
		}
	}

	return _console, nil
}

// func (c *console) GetWindow() *sdl.Window {
// 	return c.window
// }

func (c *console) GetBounds() image.Rectangle {
	// TODO
	//	return c.screen.Bounds()
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

	//go c.saveState()

	// poll events
	endFrame = time.Now() // init end frame
	startFrame = time.Now()

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

	screen.ReplacePixels(pb.pixelSurface.Pix)

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

// saveState - saves console state periodically
func (c *console) saveState() {

	// ticker := time.NewTicker(1 * time.Second)

	// for {
	// 	select {
	// 	case <-ticker.C:
	// 		// save state
	// 		c.state.SaveState(c)

	// 	}
	// }
}

func (c *console) Quit() {
	c.hasQuit = true
}

// toggleCLI - toggle between CLI and secondary mode
// func (c *console) toggleCLI() {
// 	switch c.currentMode {
// 	case CLI:
// 		c.SetMode(c.secondaryMode)
// 	case RUNTIME:
// 		if mode, ok := c.modes[c.currentMode]; ok {
// 			runtime := mode.(*runtime)
// 			runtime.Stop()
// 			c.SetMode(CLI)
// 		}

// 	default:
// 		c.secondaryMode = c.currentMode
// 		c.SetMode(CLI)
// 	}
// }

// Destroy cleans up any resources at end
func (c *console) Destroy() {
	//c.window.Destroy()
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
