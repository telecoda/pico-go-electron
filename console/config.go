package console

type Config struct {
	BorderWidth     int
	ConsoleWidth    int
	ConsoleHeight   int
	ScreenshotScale int
	GifScale        int
	GifLength       int
	// private vars
	palette     *palette
	consoleType ConsoleType
	fontWidth   int
	fontHeight  int
	BgColor     ColorID
	FgColor     ColorID
	BorderColor ColorID
}

var optVerbose bool
var screenshotScale int
var gifScale int
var gifLength int

func NewConfig(consoleType ConsoleType) Config {
	switch consoleType {
	case PICO8:
		return newPico8Config()
	case TIC80:
		return newTic80Config()
	case ZXSPECTRUM:
		return newZXSpectrumConfig()
	case CBM64:
		return newCBM64Config()
	}
	return newPico8Config() // always default to PICO8
}

// Default configs for different console types
func newPico8Config() Config {
	config := Config{
		ConsoleWidth:    128,
		ConsoleHeight:   128,
		ScreenshotScale: screenshotScale,
		GifScale:        gifScale,
		GifLength:       gifLength,
		consoleType:     PICO8,
		fontWidth:       4,
		fontHeight:      8,
		BgColor:         PICO8_BLACK,
		FgColor:         PICO8_WHITE,
		BorderColor:     PICO8_BLACK,
	}
	return config
}

func newTic80Config() Config {
	config := Config{
		ConsoleWidth:    240,
		ConsoleHeight:   136,
		ScreenshotScale: screenshotScale,
		GifScale:        gifScale,
		GifLength:       gifLength,
		consoleType:     TIC80,
		fontWidth:       8,
		fontHeight:      8,
		BgColor:         TIC80_BLACK,
		FgColor:         TIC80_WHITE,
		BorderColor:     TIC80_BLACK,
	}
	return config
}

func newZXSpectrumConfig() Config {
	config := Config{
		BorderWidth:     25,
		ConsoleWidth:    256,
		ConsoleHeight:   192,
		ScreenshotScale: screenshotScale,
		GifScale:        gifScale,
		GifLength:       gifLength,
		consoleType:     ZXSPECTRUM,
		fontWidth:       8,
		fontHeight:      8,
		BgColor:         ZX_WHITE,
		FgColor:         ZX_BLACK,
		BorderColor:     ZX_WHITE,
	}
	return config
}

func newCBM64Config() Config {
	config := Config{
		BorderWidth:     25,
		ConsoleWidth:    320,
		ConsoleHeight:   200,
		ScreenshotScale: screenshotScale,
		GifScale:        gifScale,
		GifLength:       gifLength,
		consoleType:     CBM64,
		fontWidth:       8,
		fontHeight:      8,
		BgColor:         C64_BLUE,
		FgColor:         C64_LIGHT_BLUE,
		BorderColor:     C64_LIGHT_BLUE,
	}
	return config
}
