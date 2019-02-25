package console

type BaseCartridge struct {
	cfg          Config // holds details of console config
	PixelBuffer         // ref to console display
	PicoInputAPI        // ref to input logic
}

// NewBaseCart - initialise a struct implementing Cartridge interface
func NewBaseCart() *BaseCartridge {
	cart := &BaseCartridge{}

	return cart
}

// GetConfig - return config need for Cart to run
func (bc *BaseCartridge) GetConfig() Config {
	return bc.cfg
}

func (bc *BaseCartridge) initPb(pb PixelBuffer) {
	// the initPb method receives a PixelBuffer reference
	// hold onto this reference, this is the display that
	// your code will be drawing onto each frame
	bc.PixelBuffer = pb
}

func (bc *BaseCartridge) initInputs(in PicoInputAPI) {
	// the initInputs method receives a PicoInputAPI reference
	// hold onto this reference, this allows the cartridge to access input such as mouse/keyboard
	bc.PicoInputAPI = in
}

// func (bc *BaseCartridge) Btn(id int) bool {
// 	// access runtime button mappings
// 	//return _console.Btn(id)
// 	// TODO
// 	return false
// }

// func (bc *BaseCartridge) IsMousePressed() bool {
// 	return ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
// }

// func (bc *BaseCartridge) MousePosition() (int, int) {
// 	return ebiten.CursorPosition()
// }
