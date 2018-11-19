package console

type BaseCartridge struct {
	cfg         Config // holds details of console config
	PixelBuffer        // ref to console display
	PicoInputAPI
}

// NewBaseCart - initialise a struct implementing Cartridge interface
func NewBaseCart() *BaseCartridge {
	cart := &BaseCartridge{
		//		cfg: _console.Config,
	}

	return cart
}

// GetConfig - return config need for Cart to run
func (bc *BaseCartridge) GetConfig() Config {
	return bc.cfg
}

// Init - called once when cart is initialised
func (bc *BaseCartridge) initPb(pb PixelBuffer) {
	// the Init method receives a PixelBuffer reference
	// hold onto this reference, this is the display that
	// your code will be drawing onto each frame
	bc.PixelBuffer = pb
}

func (bc *BaseCartridge) getPb() PixelBuffer {
	return bc.PixelBuffer
}

func (bc *BaseCartridge) Btn(id int) bool {
	// access runtime button mappings
	//return _console.Btn(id)
	// TODO
	return false
}
