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

var lastFrame time.Time
var startFrame time.Time
var endFrame time.Time

type console struct {
	sync.Mutex
	Config

	showFPS bool

	// files
	baseDir    string
	currentDir string

	consoleType ConsoleType
	cart        Cartridge

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

func Run(cart Cartridge) error {

	// init console
	// load cartridge
	// run main loop

	_console = &console{
		//Config:  cfg,
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

	_console.cart = cart

	// use PICO8 as default
	if err := SetType(PICO8); err != nil {
		return err
	}

	// init the cart
	_console.cart.Init()

	// poll events
	endFrame = time.Now() // init end frame
	startFrame = time.Now()

	return ebiten.Run(_console.update, _console.Config.ConsoleWidth, _console.Config.ConsoleHeight, 1, "pico-go")
}

func (c *console) ShowFPS(state bool) {
	c.Lock()
	defer c.Unlock()
	c.showFPS = state
}

func SetType(consoleType ConsoleType) error {
	if _console == nil {
		return fmt.Errorf("Console Init() must be called first")
	}

	// validate type
	if _, ok := ConsoleTypes[consoleType]; !ok {
		return fmt.Errorf("Console type: %s not supported", consoleType)
	}

	cfg := NewConfig(consoleType)

	_console.Config = cfg

	// init sprites
	// There are 2 sprite banks
	// 0 = User sprite bank 1
	// 1 = User sprite bank 1 mask
	_console.sprites = make([]*image.Paletted, 2)

	_console.palette = newPalette(cfg.consoleType)
	_console.originalPalette = newPalette(cfg.consoleType)

	// create paletted image
	rect := image.Rect(0, 0, cfg.ConsoleWidth, cfg.ConsoleHeight)
	pImage := image.NewPaletted(rect, _console.palette.colors)
	_console.pImage = pImage

	// init sprites
	sprites, _, err := image.Decode(bytes.NewReader(images.Sprites_png))
	if err != nil {
		return fmt.Errorf("Error loading sprites: %s", err)
	}
	_console.sprites[userSpriteBank1] = sprites.(*image.Paletted)

	_console.sprites[userSpriteBank1].Palette = _console.palette.colors

	// create a mask
	masks, _, err := image.Decode(bytes.NewReader(images.Sprites_png))
	if err != nil {
		return fmt.Errorf("Error loading sprites: %s", err)
	}
	_console.sprites[userSpriteMask1] = masks.(*image.Paletted)

	// convert all black pixels to zero alpha
	mask := _console.sprites[userSpriteMask1]
	// set palette on mask
	maskPalette := newPalette(PICO8)
	maskPalette.colors[0] = color.RGBA{R: 0, G: 0, B: 0, A: 0}
	for i := 1; i < 16; i++ {
		maskPalette.colors[i] = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	mask.Palette = maskPalette.colors

	// init pixelbuffer
	pb, err := newPixelBuffer(_console.Config)
	if err != nil {
		return fmt.Errorf("Error creating pixel buffer: %s", err)
	}

	_console.cart.initPb(pb)

	return nil
}

func (c *console) update(screen *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}

	c.cart.Update()
	c.cart.Render()

	// record frame
	//			c.recorder.AddFrame(mode.GetFrame(), mode)

	cpb := c.cart.getPb()

	pb := cpb.getPixelBuffer()

	// convert paletted image to RGBA

	pix := make([]uint8, c.ConsoleWidth*c.ConsoleHeight*4)

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

	screen.ReplacePixels(pix)

	if c.showFPS {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))
	}

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

	//	}

	return nil
}

func (c *console) mouseClicked(x, y int32) {
	// transform window x,y coords to pixel buffer coords
	fmt.Printf("Mouse clicked at x: %d y: %d\n", x, y)
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
