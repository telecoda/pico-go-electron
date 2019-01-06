package console

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	drawx "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/f64"
	"golang.org/x/image/math/fixed"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type mode struct {
	pixelBuffer
}

type pixelBuffer struct {
	textCursor pos // note print pos in char/line pos not pixel pos
	fgColor    ColorID
	bgColor    ColorID
	palette    *palette
	r          []uint8 // lookup for red color component
	g          []uint8 // lookup for green color component
	b          []uint8 // lookup for blue color component
	a          []uint8 // lookup for alpha color component

	charCols     int
	charRows     int
	pixelSurface *image.Paletted // offscreen pixel buffer
	rgbaPixels   []uint8
	screen       *ebiten.Image
	psRect       image.Rectangle // rect of pixelSurface
	renderRect   image.Rectangle // rect on main window that pixelbuffer is rendered into

	// these are temp paletted images stored by size for reuse
	copySpritesMap map[image.Rectangle]*image.Paletted
	txSpritesMap   map[image.Rectangle]*image.Paletted
	maskSpritesMap map[image.Rectangle]*image.Paletted

	// these are cached versions of previously transformed sprites
	spriteCache map[spriteTx]spriteCached

	flipReady bool
}

type spriteTx struct {
	number      int
	width       int
	height      int
	scaleWidth  int
	scaleHeight int
	rot         int
	flipX       bool
	flipY       bool
}

type spriteCached struct {
	txImage   *image.Paletted
	maskImage *image.Paletted
	lastUsed  time.Time
}

type pos struct {
	x int
	y int
}

func newPixelBuffer(cfg Config) (*pixelBuffer, error) {
	p := &pixelBuffer{}

	p.psRect = image.Rect(0, 0, cfg.ConsoleWidth, cfg.ConsoleHeight)
	p.renderRect = image.Rect(0, 0, cfg.ConsoleWidth, cfg.ConsoleHeight)

	ps := image.NewPaletted(p.psRect, cfg.palette.colors)

	if ps == nil {
		return nil, fmt.Errorf("Surface is nil")
	}

	p.palette = cfg.palette
	p.r = make([]uint8, len(_console.palette.colorMap), len(_console.palette.colorMap))
	p.g = make([]uint8, len(_console.palette.colorMap), len(_console.palette.colorMap))
	p.b = make([]uint8, len(_console.palette.colorMap), len(_console.palette.colorMap))
	p.a = make([]uint8, len(_console.palette.colorMap), len(_console.palette.colorMap))

	p.spriteCache = make(map[spriteTx]spriteCached)

	if err := setSurfacePalette(p.palette, ps); err != nil {
		return nil, err
	}

	var err error
	screen, err := ebiten.NewImageFromImage(ps, ebiten.FilterNearest)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ebiten image %s", err)
	}
	p.screen = screen
	p.pixelSurface = ps

	p.textCursor.x = 0
	p.textCursor.y = 0
	p.fgColor = 7
	p.bgColor = 0

	p.charCols = cfg.ConsoleWidth / cfg.fontWidth
	p.charRows = cfg.ConsoleHeight / cfg.fontHeight

	p.rgbaPixels = make([]uint8, _console.Config.ConsoleWidth*_console.Config.ConsoleHeight*4)
	// init temp sprite maps
	p.copySpritesMap = make(map[image.Rectangle]*image.Paletted)
	p.txSpritesMap = make(map[image.Rectangle]*image.Paletted)
	p.maskSpritesMap = make(map[image.Rectangle]*image.Paletted)

	return p, nil
}

func (p *pixelBuffer) GetFrame() *image.Paletted {
	return p.pixelSurface
}

func (p *pixelBuffer) Render() error {

	// this is never called, always locally implemented

	return nil

}

// API

// Cls - clears pixel buffer
func (p *pixelBuffer) Cls(colorID ...ColorID) {
	// clear buffer with background color
	if len(colorID) != 0 {
		p.bgColor = colorID[0]
	}

	bg := uint8(p.bgColor)

	// fill every pixel with same color
	for i, _ := range p.pixelSurface.Pix {
		p.pixelSurface.Pix[i] = bg
	}
}

func (p *pixelBuffer) Cursor(x, y int) {
	p.textCursor.x = x
	p.textCursor.y = y
}

var delay = time.Duration(1 * time.Second / 60)

// Flip - copy offscreen buffer to onscreen buffer
func (p *pixelBuffer) Flip() error {

	if p.pixelSurface == nil {
		return fmt.Errorf("No pixelsurface")
	}

	// if _console.screen == nil {
	// 	fmt.Printf("TEMP: no screen\n")
	// 	return nil
	// }

	if !p.flipReady {
		time.Sleep(delay)
		return nil
	} else {
	}

	p.flipReady = false
	// record frame
	//_console.recorder.AddFrame(p.GetFrame(), p)

	// at end of frame delay start timing for next one
	startFrame = time.Now()

	p.copyIndexedToRGBA()

	_console.screen.ReplacePixels(p.rgbaPixels)
	if _console.showFPS {
		ebitenutil.DebugPrint(_console.screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))
	}

	return nil
}

// copyIndexedToRGBA - convert the paletted (indexed) image into a set of RGBA pixels for rendering on display
func (p *pixelBuffer) copyIndexedToRGBA() {
	for id, rgba := range _console.palette.colorMap {
		p.r[id] = rgba.R
		p.g[id] = rgba.G
		p.b[id] = rgba.B
		p.a[id] = rgba.A
	}

	i := 0
	for _, palPix := range p.pixelSurface.Pix {
		p.rgbaPixels[i] = p.r[palPix]
		i++
		p.rgbaPixels[i] = p.g[palPix]
		i++
		p.rgbaPixels[i] = p.b[palPix]
		i++
		p.rgbaPixels[i] = p.a[palPix]
		i++
	}
}

func (p *pixelBuffer) GetCursor() pos {
	return p.textCursor
}

func charToPixel(charPos pos) pos {
	return pos{
		x: charPos.x * _console.Config.fontWidth,
		y: charPos.y * _console.Config.fontHeight,
	}
}

func pixelToChar(pixelPos pos) pos {
	return pos{
		x: pixelPos.x / _console.Config.fontWidth,
		y: pixelPos.y / _console.Config.fontHeight,
	}
}

// ScrollUpLine - scrolls display up a single line
func (p *pixelBuffer) ScrollUpLine() {

	// to scroll the screen up a line we need to copy the bottom part of the image
	// to the top part of the image.
	// Rect size will be total height - 1 line
	// then we blank out the bottom line
	// fromRect := &sdl.Rect{X: 0, Y: int32(_console.Config.fontHeight), W: p.pixelSurface.W, H: p.pixelSurface.H - int32(_console.Config.fontHeight)}

	srcRect := image.Rect(0, _console.Config.fontHeight, _console.ConsoleHeight, _console.Config.ConsoleWidth)
	drawx.Draw(p.pixelSurface, srcRect, p.pixelSurface, image.Point{}, drawx.Over)

	// TODO
	// fromRect := &sdl.Rect{X: 0, Y: int32(_console.Config.fontHeight), W: p.pixelSurface.W, H: p.pixelSurface.H - int32(_console.Config.fontHeight)}
	// toRect := &sdl.Rect{X: 0, Y: 0, W: p.pixelSurface.W, H: p.pixelSurface.H - int32(_console.Config.fontHeight)}
	// p.pixelSurface.Blit(fromRect, p.pixelSurface, toRect)
	// p.textCursor.y = p.charRows - 2
}

// Print - prints string of characters to the screen with drawing color
func (p *pixelBuffer) Print(str string) {
	pixelPos := charToPixel(p.textCursor)

	p.PrintAt(str, int(pixelPos.x), int(pixelPos.y), p.fgColor)

	// increase printPos by 1 line
	p.textCursor.y++

	if p.textCursor.y > p.charRows-2 {
		p.ScrollUpLine()
	}
}

// PrintAt - prints a string of characters to the screen at position with drawing color
func (p *pixelBuffer) PrintAt(str string, x, y int, colorID ...ColorID) {
	if len(colorID) == 0 {
		p.printAtWithColor(str, x, y, p.fgColor)
	} else {
		p.printAtWithColor(str, x, y, colorID[0])
	}
}

// PrintAtWithColor - prints a string of characters to the screen at position with color
func (p *pixelBuffer) printAtWithColor(str string, x, y int, colorID ColorID) {
	p.fgColor = colorID

	if str != "" {
		y += (_console.fontHeight / 2) + 1
		col := _console.palette.colors[colorID]
		point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

		d := &font.Drawer{
			Dst:  p.pixelSurface,
			Src:  image.NewUniform(col),
			Face: _console.font,
			Dot:  point,
		}
		d.DrawString(str)

	}

	// save print pos
	p.textCursor = pixelToChar(pos{x, y})

}

// Drawer methods

// Circle - draw circle with drawing color
func (p *pixelBuffer) Circle(x, y, r int, colorID ...ColorID) {
	if len(colorID) == 0 {
		p.circleWithColor(x, y, r, p.fgColor)
	} else {
		p.circleWithColor(x, y, r, colorID[0])
	}
}

// CircleWithColor - draw circle with color
func (p *pixelBuffer) circleWithColor(x0, y0, r int, colorID ColorID) {
	p.fgColor = colorID

	x := 0
	y := r
	p0 := (5 - r*4) / 4

	/* circle calcs from this blog
	http://xiaohuiliucuriosity.blogspot.co.uk/2015/03/draw-circle-using-integer-arithmetic.html
	*/

	p.circlePoints(x0, y0, x, y)
	for x < y {
		x++
		if p0 < 0 {
			p0 += 2*x + 1
		} else {
			y--
			p0 += 2*(x-y) + 1
		}
		p.circlePoints(x0, y0, x, y)
	}

}

func (p *pixelBuffer) circlePoints(cx, cy, x, y int) {

	if x == 0 {
		p.PSet(cx, cy+y)
		p.PSet(cx, cy-y)
		p.PSet(cx+y, cy)
		p.PSet(cx-y, cy)
	} else if x == y {
		p.PSet(cx+x, cy+y)
		p.PSet(cx-x, cy+y)
		p.PSet(cx+x, cy-y)
		p.PSet(cx-x, cy-y)
	} else if x < y {
		p.PSet(cx+x, cy+y)
		p.PSet(cx-x, cy+y)
		p.PSet(cx+x, cy-y)
		p.PSet(cx-x, cy-y)
		p.PSet(cx+y, cy+x)
		p.PSet(cx-y, cy+x)
		p.PSet(cx+y, cy-x)
		p.PSet(cx-y, cy-x)
	}
}

// CircleFill - fill circle with drawing color
func (p *pixelBuffer) CircleFill(x, y, r int, colorID ...ColorID) {
	if len(colorID) == 0 {
		p.circleFillWithColor(x, y, r, p.fgColor)
	} else {
		p.circleFillWithColor(x, y, r, colorID[0])
	}
}

// CircleFillWithColor - fill circle with color
func (p *pixelBuffer) circleFillWithColor(x0, y0, r int, colorID ColorID) {
	p.fgColor = colorID

	x := 0
	y := r
	p0 := (5 - r*4) / 4

	/* circle calcs from this blog
	http://groups.csail.mit.edu/graphics/classes/6.837/F98/Lecture6/circle.html
	*/

	p.circleLines(x0, y0, x, y)
	for x < y {
		x++
		if p0 < 0 {
			p0 += 2*x + 1
		} else {
			y--
			p0 += 2*(x-y) + 1
		}
		p.circleLines(x0, y0, x, y)
	}
}

func (p *pixelBuffer) circleLines(cx, cy, x, y int) {
	p.Line(cx-x, cy+y, cx+x, cy+y)
	p.Line(cx-x, cy-y, cx+x, cy-y)
	p.Line(cx-y, cy+x, cx+y, cy+x)
	p.Line(cx-y, cy-x, cx+y, cy-x)
}

// Line - line in drawing color
func (p *pixelBuffer) Line(x0, y0, x1, y1 int, colorID ...ColorID) {
	if len(colorID) == 0 {
		p.lineWithColor(x0, y0, x1, y1, p.fgColor)
	} else {
		p.lineWithColor(x0, y0, x1, y1, colorID[0])
	}
}

// LineWithColor - line with color
func (p *pixelBuffer) lineWithColor(x1, y1, x2, y2 int, colorID ColorID) {
	p.setFGColor(colorID)

	col := p.palette.GetColor(colorID)

	/* Code from
	https://github.com/StephaneBunel/bresenham/blob/master/drawline.go#L12-L22
	*/
	// draw line

	var dx, dy, e, slope int

	// Because drawing p1 -> p2 is equivalent to draw p2 -> p1,
	// I sort points in x-axis order to handle only half of possible cases.
	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dx, dy = x2-x1, y2-y1
	// Because point is x-axis ordered, dx cannot be negative
	if dy < 0 {
		dy = -dy
	}

	switch {

	// Is line a point ?
	case x1 == x2 && y1 == y2:
		p.pixelSurface.Set(x1, y1, col)

	// Is line an horizontal ?
	case y1 == y2:
		for ; dx != 0; dx-- {
			p.pixelSurface.Set(x1, y1, col)
			x1++
		}
		p.pixelSurface.Set(x1, y1, col)

	// Is line a vertical ?
	case x1 == x2:
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for ; dy != 0; dy-- {
			p.pixelSurface.Set(x1, y1, col)
			y1++
		}
		p.pixelSurface.Set(x1, y1, col)

	// Is line a diagonal ?
	case dx == dy:
		if y1 < y2 {
			for ; dx != 0; dx-- {
				p.pixelSurface.Set(x1, y1, col)
				x1++
				y1++
			}
		} else {
			for ; dx != 0; dx-- {
				p.pixelSurface.Set(x1, y1, col)
				x1++
				y1--
			}
		}
		p.pixelSurface.Set(x1, y1, col)

	// wider than high ?
	case dx > dy:
		if y1 < y2 {
			// BresenhamDxXRYD(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				p.pixelSurface.Set(x1, y1, col)
				x1++
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}
			}
		} else {
			// BresenhamDxXRYU(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				p.pixelSurface.Set(x1, y1, col)
				x1++
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}
			}
		}
		p.pixelSurface.Set(x2, y2, col)

	// higher than wide.
	default:
		if y1 < y2 {
			// BresenhamDyXRYD(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				p.pixelSurface.Set(x1, y1, col)
				y1++
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		} else {
			// BresenhamDyXRYU(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				p.pixelSurface.Set(x1, y1, col)
				y1--
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		}
		p.pixelSurface.Set(x2, y2, col)
	}
}

func (p *pixelBuffer) setFGColor(colorID ColorID) {
	p.fgColor = colorID
	// c := p.palette.GetColor(colorID)
	// r, g, b, _ := c.RGBA()
	// if p.gc != nil {
	// 	p.gc.SetRGB255(int(r), int(g), int(b))
	// }
}

// PGet - pixel get
func (p *pixelBuffer) PGet(x, y int) ColorID {

	c := p.pixelSurface.At(x, y)
	r, g, b, a := c.RGBA()
	color := rgba{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}

	return p.palette.GetColorID(color)
}

// PSet - pixel set in drawing color
func (p *pixelBuffer) PSet(x, y int, colorID ...ColorID) {
	if len(colorID) == 0 {
		p.pSetWithColor(x, y, p.fgColor)
	} else {
		p.pSetWithColor(x, y, colorID[0])
	}
}

// PSetWithColor - pixel set with color
func (p *pixelBuffer) pSetWithColor(x0, y0 int, colorID ColorID) {
	p.setFGColor(colorID)
	p.pixelSurface.Set(x0, y0, p.palette.GetColor(colorID))
}

// Rect - draw rectangle with drawing color
func (p *pixelBuffer) Rect(x0, y0, x1, y1 int, colorID ...ColorID) {
	if len(colorID) == 0 {
		p.rectWithColor(x0, y0, x1, y1, p.fgColor)
	} else {
		p.rectWithColor(x0, y0, x1, y1, colorID[0])
	}
}

// RectWithColor - draw rectangle with color
func (p *pixelBuffer) rectWithColor(x0, y0, x1, y1 int, colorID ColorID) {
	p.fgColor = colorID
	p.Line(x0, y0, x1, y0)
	p.Line(x1, y0, x1, y1)
	p.Line(x1, y1, x0, y1)
	p.Line(x0, y1, x0, y0)
}

// RectFill - fill rectangle with drawing color
func (p *pixelBuffer) RectFill(x0, y0, x1, y1 int, colorID ...ColorID) {
	if len(colorID) == 0 {
		p.rectFillWithColor(x0, y0, x1, y1, p.fgColor)
	} else {
		p.rectFillWithColor(x0, y0, x1, y1, colorID[0])
	}
}

// RectFillWithColor - fill rectangle with color
func (p *pixelBuffer) rectFillWithColor(x0, y0, x1, y1 int, colorID ColorID) {
	p.fgColor = colorID
	for x := x0; x < x1; x++ {
		p.Line(x, y0, x, y1)
	}
}

// Spriter methods
func (p *pixelBuffer) Sprite(n, x, y, w, h, dw, dh int) {
	p.sprite(n, x, y, w, h, dw, dh, 0, false, false)
}

func (p *pixelBuffer) SpriteFlipped(n, x, y, w, h, dw, dh int, flipX, flipY bool) {
	p.spriteWithCache(n, x, y, w, h, dw, dh, 0, flipX, flipY)
}

func (p *pixelBuffer) SpriteRotated(n, x, y, w, h, dw, dh int, rot int) {
	p.spriteWithCache(n, x, y, w, h, dw, dh, rot, false, false)
}

func (p *pixelBuffer) sprite(n, x, y, w, h, dw, dh, rot int, flipX, flipY bool) {

	_console.currentSpriteBank = userSpriteBank1

	sw := w * _spriteWidth
	sh := h * _spriteHeight

	// convert sprite number into x,y pos
	xCell := n % _spritesPerLine
	yCell := (n - xCell) / _spritesPerLine

	xPos := xCell * _spriteWidth
	yPos := yCell * _spriteHeight

	// this is the rect to copy from sprite sheet
	spriteSrcRect := image.Rect(xPos, yPos, xPos+sw, yPos+sh)
	// this rect is where the sprite will be copied to
	screenRect := image.Rect(x, y, x+dw, y+dh)

	if flipX || flipY || rot != 0 {
		// we need to transform & mask

		// setup transform matrix

		// do nothing matrix
		// f64.Aff3{1, -0, 0, 0, 1, 0}
		// matrix  = {xscale, rotx, xOffset, rotY, yscale, yOffset}

		var matrix f64.Aff3

		if flipX && !flipY {
			// flip x
			matrix = f64.Aff3{-1, 0, float64(sw), 0, 1, 0}
		}
		if flipY && !flipX {
			// flip y
			matrix = f64.Aff3{1, 0, 0, 0, -1, float64(sh)}
		}
		if flipX && flipY {
			// flip xy
			matrix = f64.Aff3{-1, 0, float64(sw), 0, -1, float64(sh)}
		}

		if rot != 0 {
			angle := float64(rot)
			a := math.Pi * angle / 180
			xf, yf := float64(sw/2), float64(sh/2)
			sin := math.Sin(a)
			cos := math.Cos(a)
			matrix = f64.Aff3{
				cos, -sin, xf - xf*cos + yf*sin,
				sin, cos, yf - xf*sin - yf*cos,
			}
		}

		// create a copy of the sprite
		copyRect := image.Rect(0, 0, sw, sh)
		copyImage := image.NewPaletted(copyRect, _console.sprites[userSpriteBank1].Palette)
		point := image.Point{X: 0, Y: 0}
		options := &drawx.Options{
			SrcMask:  _console.sprites[userSpriteMask1],
			SrcMaskP: image.Point{0, 0},
		}
		drawx.Copy(copyImage, point, _console.sprites[userSpriteBank1], spriteSrcRect, drawx.Src, options)

		// rotate it
		txImage := image.NewPaletted(copyRect, _console.sprites[userSpriteBank1].Palette)
		drawx.NearestNeighbor.Transform(txImage, matrix, copyImage, copyRect, drawx.Src, nil)

		// create a copy of the sprite as a mask
		maskRect := image.Rect(0, 0, sw, sh)
		maskImage := image.NewPaletted(maskRect, _console.sprites[userSpriteMask1].Palette)
		drawx.Copy(maskImage, point, txImage, maskRect, drawx.Src, nil)

		options = &drawx.Options{
			SrcMask:  maskImage,
			SrcMaskP: image.Point{0, 0},
		}

		drawx.NearestNeighbor.Scale(p.pixelSurface, screenRect, txImage, maskRect, drawx.Over, options)

		return
	}

	options := &drawx.Options{
		SrcMask:  _console.sprites[userSpriteMask1],
		SrcMaskP: image.Point{0, 0},
	}

	drawx.NearestNeighbor.Scale(p.pixelSurface, screenRect, _console.sprites[userSpriteBank1], spriteSrcRect, drawx.Over, options)

}

func (p *pixelBuffer) spriteWithMaps(n, x, y, w, h, dw, dh, rot int, flipX, flipY bool) {

	_console.currentSpriteBank = userSpriteBank1

	sw := w * _spriteWidth
	sh := h * _spriteHeight

	// convert sprite number into x,y pos
	xCell := n % _spritesPerLine
	yCell := (n - xCell) / _spritesPerLine

	xPos := xCell * _spriteWidth
	yPos := yCell * _spriteHeight

	// this is the rect to copy from sprite sheet
	spriteSrcRect := image.Rect(xPos, yPos, xPos+sw, yPos+sh)
	// this rect is where the sprite will be copied to
	screenRect := image.Rect(x, y, x+dw, y+dh)

	if flipX || flipY || rot != 0 {
		// we need to transform & mask

		// setup transform matrix

		// do nothing matrix
		// f64.Aff3{1, -0, 0, 0, 1, 0}
		// matrix  = {xscale, rotx, xOffset, rotY, yscale, yOffset}

		var matrix f64.Aff3

		if flipX && !flipY {
			// flip x
			matrix = f64.Aff3{-1, 0, float64(sw), 0, 1, 0}
		}
		if flipY && !flipX {
			// flip y
			matrix = f64.Aff3{1, 0, 0, 0, -1, float64(sh)}
		}
		if flipX && flipY {
			// flip xy
			matrix = f64.Aff3{-1, 0, float64(sw), 0, -1, float64(sh)}
		}

		if rot != 0 {
			angle := float64(rot)
			a := math.Pi * angle / 180
			xf, yf := float64(sw/2), float64(sh/2)
			sin := math.Sin(a)
			cos := math.Cos(a)
			matrix = f64.Aff3{
				cos, -sin, xf - xf*cos + yf*sin,
				sin, cos, yf - xf*sin - yf*cos,
			}
		}

		// create a copy of the sprite
		copyRect := image.Rect(0, 0, sw, sh)
		copyImage := p.getCopyImage(copyRect)
		point := image.Point{X: 0, Y: 0}
		options := &drawx.Options{
			SrcMask:  _console.sprites[userSpriteMask1],
			SrcMaskP: image.Point{0, 0},
		}
		drawx.Copy(copyImage, point, _console.sprites[userSpriteBank1], spriteSrcRect, drawx.Src, options)

		// rotate it
		txImage := p.getTxImage(copyRect)

		drawx.NearestNeighbor.Transform(txImage, matrix, copyImage, copyRect, drawx.Src, nil)

		// create a copy of the sprite as a mask
		maskRect := image.Rect(0, 0, sw, sh)
		maskImage := p.getMaskImage(maskRect)
		drawx.Copy(maskImage, point, txImage, maskRect, drawx.Src, nil)

		options = &drawx.Options{
			SrcMask:  maskImage,
			SrcMaskP: image.Point{0, 0},
		}

		drawx.NearestNeighbor.Scale(p.pixelSurface, screenRect, txImage, maskRect, drawx.Over, options)

		return
	}

	options := &drawx.Options{
		SrcMask:  _console.sprites[userSpriteMask1],
		SrcMaskP: image.Point{0, 0},
	}

	drawx.NearestNeighbor.Scale(p.pixelSurface, screenRect, _console.sprites[userSpriteBank1], spriteSrcRect, drawx.Over, options)

}

func (p *pixelBuffer) spriteWithCache(n, x, y, w, h, dw, dh, rot int, flipX, flipY bool) {

	_console.currentSpriteBank = userSpriteBank1

	sw := w * _spriteWidth
	sh := h * _spriteHeight

	// convert sprite number into x,y pos
	xCell := n % _spritesPerLine
	yCell := (n - xCell) / _spritesPerLine

	xPos := xCell * _spriteWidth
	yPos := yCell * _spriteHeight

	// this is the rect to copy from sprite sheet
	spriteSrcRect := image.Rect(xPos, yPos, xPos+sw, yPos+sh)
	// this rect is where the sprite will be copied to
	screenRect := image.Rect(x, y, x+dw, y+dh)

	if flipX || flipY || rot != 0 {

		tx := spriteTx{
			number:      n,
			width:       w,
			height:      h,
			scaleWidth:  dw,
			scaleHeight: dh,
			rot:         rot,
			flipX:       flipX,
			flipY:       flipY,
		}

		cached, ok := p.spriteCache[tx]

		if !ok {
			// do transform

			// we need to transform & mask

			// setup transform matrix

			// do nothing matrix
			// f64.Aff3{1, -0, 0, 0, 1, 0}
			// matrix  = {xscale, rotx, xOffset, rotY, yscale, yOffset}

			var matrix f64.Aff3
			if flipY && !flipX {
				// flip y
				matrix = f64.Aff3{1, 0, 0, 0, -1, float64(sh)}
			}
			if flipX && !flipY {
				// flip x
				matrix = f64.Aff3{-1, 0, float64(sw), 0, 1, 0}
			}
			if flipX && flipY {
				// flip xy
				matrix = f64.Aff3{-1, 0, float64(sw), 0, -1, float64(sh)}
			}

			if rot != 0 {
				angle := float64(rot)
				a := math.Pi * angle / 180
				xf, yf := float64(sw/2), float64(sh/2)
				sin := math.Sin(a)
				cos := math.Cos(a)
				matrix = f64.Aff3{
					cos, -sin, xf - xf*cos + yf*sin,
					sin, cos, yf - xf*sin - yf*cos,
				}
			}

			// create a copy of the sprite
			copyRect := image.Rect(0, 0, sw, sh)
			// this fetches and empty image of the correct dimensions
			copyImage := image.NewPaletted(copyRect, _console.sprites[userSpriteBank1].Palette)

			point := image.Point{X: 0, Y: 0}
			options := &drawx.Options{
				SrcMask:  _console.sprites[userSpriteMask1],
				SrcMaskP: image.Point{0, 0},
			}
			drawx.Copy(copyImage, point, _console.sprites[userSpriteBank1], spriteSrcRect, drawx.Src, options)

			// rotate it
			txImage := image.NewPaletted(copyRect, _console.sprites[userSpriteBank1].Palette)

			drawx.NearestNeighbor.Transform(txImage, matrix, copyImage, copyRect, drawx.Src, nil)

			// create a copy of the sprite as a mask
			maskRect := image.Rect(0, 0, sw, sh)
			maskImage := image.NewPaletted(maskRect, _console.sprites[userSpriteMask1].Palette)

			drawx.Copy(maskImage, point, txImage, maskRect, drawx.Src, nil)

			options = &drawx.Options{
				SrcMask:  maskImage,
				SrcMaskP: image.Point{0, 0},
			}

			drawx.NearestNeighbor.Scale(p.pixelSurface, screenRect, txImage, maskRect, drawx.Over, options)

			// store in cache

			p.spriteCache[tx] = spriteCached{
				txImage:   txImage,
				maskImage: maskImage,
				lastUsed:  time.Now(),
			}

			if len(p.spriteCache) > MaxSpriteCache {
				// too many sprites in cache, lets delete the old ones.
				now := time.Now()
				delKeys := make([]spriteTx, 0)
				for key, cached := range p.spriteCache {
					if now.Sub(cached.lastUsed) > MaxCacheAge {
						// out of date
						delKeys = append(delKeys, key)
					}
				}
				// do we have any sprites to delete?
				if len(delKeys) > 0 {
					for _, key := range delKeys {
						delete(p.spriteCache, key)
					}
				}
			}

			return
		} else {
			//			fmt.Printf("TEMP: tx: %#v\n", tx)
			//fmt.Printf("TEMP: cached: %#v\n", cached)

			options := &drawx.Options{
				SrcMask:  cached.maskImage,
				SrcMaskP: image.Point{0, 0},
			}
			maskRect := image.Rect(0, 0, sw, sh)

			drawx.NearestNeighbor.Scale(p.pixelSurface, screenRect, cached.txImage, maskRect, drawx.Over, options)

			// update last used time
			cached.lastUsed = time.Now()
			p.spriteCache[tx] = cached

			return

		}

	}

	options := &drawx.Options{
		SrcMask:  _console.sprites[userSpriteMask1],
		SrcMaskP: image.Point{0, 0},
	}

	drawx.NearestNeighbor.Scale(p.pixelSurface, screenRect, _console.sprites[userSpriteBank1], spriteSrcRect, drawx.Over, options)
}

// getCopyImage returns an empty image with the correct dimensions
func (p *pixelBuffer) getCopyImage(r image.Rectangle) *image.Paletted {

	copyImage, ok := p.copySpritesMap[r]
	if ok {
		return copyImage
	}

	copyImage = image.NewPaletted(r, _console.sprites[userSpriteBank1].Palette)

	p.copySpritesMap[r] = copyImage
	return copyImage
}

func (p *pixelBuffer) getTxImage(r image.Rectangle) *image.Paletted {

	txImage, ok := p.txSpritesMap[r]
	if ok {
		// fill pixels
		for i, _ := range txImage.Pix {
			txImage.Pix[i] = 0
		}
		return txImage
	}

	txImage = image.NewPaletted(r, _console.sprites[userSpriteBank1].Palette)

	p.txSpritesMap[r] = txImage
	return txImage
}

func (p *pixelBuffer) getMaskImage(r image.Rectangle) *image.Paletted {

	maskImage, ok := p.maskSpritesMap[r]
	if ok {
		return maskImage
	}

	maskImage = image.NewPaletted(r, _console.sprites[userSpriteMask1].Palette)

	p.maskSpritesMap[r] = maskImage
	return maskImage
}

// SetColor - Set current drawing color
func (p *pixelBuffer) SetColor(colorID ColorID) {
	p.fgColor = colorID
}

// Paletter methods

// getRGBA - returns color as Color and uint32
func (p *pixelBuffer) GetRGBA(color ColorID) (rgba, uint32) {
	return p.palette.GetRGBA(color)

}

// GetColor - get color by id
func (p *pixelBuffer) GetColor(colorID ColorID) color.Color {
	return p.palette.GetColor(colorID)
}

// GetColorID - find color from rgba
func (p *pixelBuffer) GetColorID(rgba rgba) ColorID {
	return p.palette.GetColorID(rgba)
}

func (p *pixelBuffer) PaletteReset() {
	p.palette.PaletteReset()
}

func (p *pixelBuffer) PaletteCopy() Paletter {
	return p.palette.PaletteCopy()
}

func (p *pixelBuffer) GetColors() []color.Color {
	return p.palette.GetColors()
}

func (p *pixelBuffer) MapColor(fromColor ColorID, toColor ColorID) error {
	if err := p.palette.MapColor(fromColor, toColor); err != nil {
		return err
	}
	// update palette for surface
	return setSurfacePalette(p.palette, p.pixelSurface)
}

func (p *pixelBuffer) SetTransparent(color ColorID, enabled bool) error {
	if err := p.palette.SetTransparent(color, enabled); err != nil {
		return err
	}
	// update palette for surface
	return setSurfacePalette(p.palette, p.pixelSurface)
}

// Destroy cleans up any resources at end
func (p *pixelBuffer) Destroy() {
	p.pixelSurface = nil
}

func (p *pixelBuffer) GetWidth() int {
	if p.pixelSurface == nil {
		return 0
	}
	return p.pixelSurface.Bounds().Dx()
}

func (p *pixelBuffer) GetHeight() int {
	if p.pixelSurface == nil {
		return 0
	}
	return p.pixelSurface.Bounds().Dy()
}
