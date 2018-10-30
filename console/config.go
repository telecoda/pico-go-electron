package console

type Config struct {
	BorderWidth     int
	ConsoleWidth    int
	ConsoleHeight   int
	WindowWidth     int
	WindowHeight    int
	FPS             int
	Verbose         bool
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
	errColor    ColorID
	cursorColor ColorID
}

var optVerbose bool
var screenshotScale int
var gifScale int
var gifLength int

// electron vars we never use
var noSandbox bool
var servicePipeToken string
var lang string

func NewConfig(consoleType ConsoleType) Config {
	switch consoleType {
	case PICO8:
		return newPico8Config()
	case TIC80:
		return newTic80Config()
	case ZX_SPECTRUM:
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
		WindowWidth:     400,
		WindowHeight:    400,
		FPS:             60,
		Verbose:         optVerbose,
		ScreenshotScale: screenshotScale,
		GifScale:        gifScale,
		GifLength:       gifLength,
		consoleType:     PICO8,
		fontWidth:       4,
		fontHeight:      8,
		BgColor:         PICO8_BLACK,
		FgColor:         PICO8_WHITE,
		errColor:        PICO8_PINK,
		BorderColor:     PICO8_BLACK,
		cursorColor:     PICO8_RED,
	}
	return config
}

func newTic80Config() Config {
	config := Config{
		ConsoleWidth:    240,
		ConsoleHeight:   136,
		WindowWidth:     480,
		WindowHeight:    272,
		FPS:             60,
		Verbose:         optVerbose,
		ScreenshotScale: screenshotScale,
		GifScale:        gifScale,
		GifLength:       gifLength,
		consoleType:     TIC80,
		fontWidth:       8,
		fontHeight:      8,
		BgColor:         TIC80_BLACK,
		FgColor:         TIC80_WHITE,
		errColor:        TIC80_YELLOW,
		BorderColor:     TIC80_BLACK,
		cursorColor:     TIC80_RED,
	}
	return config
}

func newZXSpectrumConfig() Config {
	config := Config{
		BorderWidth:     25,
		ConsoleWidth:    256,
		ConsoleHeight:   192,
		WindowWidth:     512,
		WindowHeight:    384,
		FPS:             60,
		Verbose:         optVerbose,
		ScreenshotScale: screenshotScale,
		GifScale:        gifScale,
		GifLength:       gifLength,
		consoleType:     ZX_SPECTRUM,
		fontWidth:       8,
		fontHeight:      8,
		BgColor:         ZX_WHITE,
		FgColor:         ZX_BLACK,
		errColor:        ZX_BLUE,
		BorderColor:     ZX_WHITE,
		cursorColor:     ZX_RED,
	}
	return config
}

func newCBM64Config() Config {
	config := Config{
		BorderWidth:     25,
		ConsoleWidth:    320,
		ConsoleHeight:   200,
		WindowWidth:     640,
		WindowHeight:    400,
		FPS:             60,
		Verbose:         optVerbose,
		ScreenshotScale: screenshotScale,
		GifScale:        gifScale,
		GifLength:       gifLength,
		consoleType:     CBM64,
		fontWidth:       8,
		fontHeight:      8,
		BgColor:         C64_BLUE,
		FgColor:         C64_LIGHT_BLUE,
		errColor:        C64_WHITE,
		BorderColor:     C64_LIGHT_BLUE,
		cursorColor:     C64_LIGHT_BLUE,
	}
	return config
}
