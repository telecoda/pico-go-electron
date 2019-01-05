package console

import (
	"testing"
)

type testPos struct {
	testValue     pos
	expectedValue pos
}

func TestPixelToChar(t *testing.T) {
	tests := []testPos{
		{testValue: pos{x: 10, y: 10}, expectedValue: pos{x: 40, y: 80}},
		{testValue: pos{x: 0, y: 0}, expectedValue: pos{x: 0, y: 0}},
	}

	Init(PICO8)

	for _, test := range tests {
		pixelPos := charToPixel(test.testValue)
		if pixelPos != test.expectedValue {
			t.Error(
				"For", test.testValue,
				"expected", test.expectedValue,
				"got", pixelPos,
			)
		}
	}
}

func TestCharToPixel(t *testing.T) {
	tests := []testPos{
		{testValue: pos{x: 64, y: 64}, expectedValue: pos{x: 16, y: 8}},
		{testValue: pos{x: 0, y: 0}, expectedValue: pos{x: 0, y: 0}},
	}
	for _, test := range tests {
		charPos := pixelToChar(test.testValue)
		if charPos != test.expectedValue {
			t.Error(
				"For", test.testValue,
				"expected", test.expectedValue,
				"got", charPos,
			)
		}
	}
}

func BenchmarkCopyPixels(b *testing.B) {
	// this benchmark measures the performance of the code the copies the offset pixelbuffer into an array of RGBA pixels every frame
	cfg := newPico8Config()
	cfg.palette = newPalette(cfg.consoleType)

	_console.Config = cfg

	pb, err := newPixelBuffer(cfg)
	if err != nil {
		b.Fatalf("Failed to create pixelBuffer: %s", err)
	}

	for i := 0; i < b.N; i++ {
		pb.copyIndexedToRGBA()
	}
}

func BenchmarkDrawSpriteUnscaled(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.Sprite(1, 0, 0, 16, 16, 16, 16)
		_console.pb.sprite(1, 0, 0, 16, 16, 16, 16, 0.0, false, false)
	}
}

func BenchmarkDrawSpriteUnscaledWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.Sprite(1, 0, 0, 16, 16, 16, 16)
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 16, 16, 0.0, false, false)
	}
}

func BenchmarkDrawSpriteScaled(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.Sprite(1, 0, 0, 16, 16, 32, 32)
		_console.pb.sprite(1, 0, 0, 16, 16, 32, 32, 0.0, false, false)
	}
}

func BenchmarkDrawSpriteScaledWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.Sprite(1, 0, 0, 16, 16, 32, 32)
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 32, 32, 0.0, false, false)
	}
}
func BenchmarkDrawSpriteScaledWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.Sprite(1, 0, 0, 16, 16, 32, 32)
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, 0.0, false, false)
	}
}
func BenchmarkDrawSpriteXFlipped(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 32, 32, true, false)
		_console.pb.sprite(1, 0, 0, 16, 16, 16, 16, 0.0, true, false)
	}
}

func BenchmarkDrawSpriteXFlippedWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 32, 32, true, false)
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 16, 16, 0.0, true, false)
	}
}

func BenchmarkDrawSpriteXFlippedWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 32, 32, true, false)
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 16, 16, 0.0, true, false)
	}
}
func BenchmarkDrawSpriteYFlippedWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 16, 16, false, true)
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 16, 16, 0.0, false, true)
	}
}
func BenchmarkDrawSpriteYFlipped(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 16, 16, false, true)
		_console.pb.sprite(1, 0, 0, 16, 16, 16, 16, 0.0, false, true)
	}
}

func BenchmarkDrawSpriteYFlippedWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 16, 16, false, true)
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 16, 16, 0.0, false, true)
	}
}

func BenchmarkDrawSpriteXYFlipped(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 16, 16, true, true)
		_console.pb.sprite(1, 0, 0, 16, 16, 16, 16, 0.0, true, true)
	}
}

func BenchmarkDrawSpriteXYFlippedWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 16, 16, true, true)
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 16, 16, 0.0, true, true)
	}
}

func BenchmarkDrawSpriteXYFlippedWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 16, 16, true, true)
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 16, 16, 0.0, true, true)
	}
}

func BenchmarkDrawSpriteXFlippedScaled(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 32, 32, true, false)
		_console.pb.sprite(1, 0, 0, 16, 16, 32, 32, 0.0, true, false)
	}
}

func BenchmarkDrawSpriteXFlippedScaledWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 32, 32, true, false)
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 32, 32, 0.0, true, false)
	}
}

func BenchmarkDrawSpriteXFlippedScaledWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 32, 32, true, false)
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, 0.0, true, false)
	}
}

func BenchmarkDrawSpriteYFlippedScaled(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 32, 32, false, true)
		_console.pb.sprite(1, 0, 0, 16, 16, 32, 32, 0.0, false, true)
	}
}

func BenchmarkDrawSpriteYFlippedScaledWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 32, 32, false, true)
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 32, 32, 0.0, false, true)
	}
}

func BenchmarkDrawSpriteYFlippedScaledWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 32, 32, false, true)
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, 0.0, false, true)
	}
}
func BenchmarkDrawSpriteXYFlippedScaled(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 32, 32, true, true)
		_console.pb.sprite(1, 0, 0, 16, 16, 32, 32, 0.0, true, true)
	}
}

func BenchmarkDrawSpriteXYFlippedScaledWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 32, 32, true, true)
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 32, 32, 0.0, true, true)
	}
}

func BenchmarkDrawSpriteXYFlippedScaledWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteFlipped(1, 0, 0, 16, 16, 32, 32, true, true)
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, 0.0, true, true)
	}
}

func BenchmarkDrawSpriteRotated(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteRotated(1, 0, 0, 16, 16, 32, 32, float64(i%360))
		_console.pb.sprite(1, 0, 0, 16, 16, 32, 32, i%360, false, false)
	}
}

func BenchmarkDrawSpriteRotatedWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteRotated(1, 0, 0, 16, 16, 32, 32, float64(i%360))
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 32, 32, i%360, false, false)
	}
}

func BenchmarkDrawSpriteRotatedWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteRotated(1, 0, 0, 16, 16, 32, 32, float64(i%360))
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, i%360, false, false)
	}
}

func BenchmarkDrawSpriteRotatedSameAngle(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteRotated(1, 0, 0, 16, 16, 32, 32, float64(i%360))
		_console.pb.sprite(1, 0, 0, 16, 16, 32, 32, 45, false, false)
	}
}

func BenchmarkDrawSpriteRotatedWithMapsSameAngle(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteRotated(1, 0, 0, 16, 16, 32, 32, float64(i%360))
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 32, 32, 45, false, false)
	}
}

func BenchmarkDrawSpriteRotatedWithCacheSameAngle(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		// _console.pb.SpriteRotated(1, 0, 0, 16, 16, 32, 32, float64(i%360))
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, 45, false, false)
	}
}
