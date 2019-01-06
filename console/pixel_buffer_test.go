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

func TestRotatedWithCache(t *testing.T) {
	Init(PICO8)
	for i := 0; i < 720; i++ {
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, i%360, false, false)
	}

	// key cache size
	if len(_console.pb.spriteCache) != 359 {
		// we expect a size of 359 because angle 0 is never transformed
		t.Errorf("Expected size: %d got: %d", 359, len(_console.pb.spriteCache))
	}

	// sprite in cache should differ
	keys := make([]spriteTx, len(_console.pb.spriteCache))
	sprites := make([]spriteCached, len(_console.pb.spriteCache))
	i := 0
	for key, cached := range _console.pb.spriteCache {
		keys[i] = key
		sprites[i] = cached
		i++
	}

	if keys[0] == keys[1] {
		t.Errorf("Keys should not match: %#v vs %#v", keys[0], keys[1])
	}

	if &sprites[0].txImage == &sprites[1].txImage {
		t.Errorf("Tx image addresses should not match: %#v vs %#v", sprites[0].txImage, sprites[1].txImage)
	}

	if &sprites[0].maskImage == &sprites[1].maskImage {
		t.Errorf("Mask image addresses should not match: %#v vs %#v", sprites[0].maskImage, sprites[1].maskImage)
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
		_console.pb.sprite(1, 0, 0, 16, 16, 16, 16, 0.0, false, false)
	}
}

func BenchmarkDrawSpriteUnscaledWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 16, 16, 0.0, false, false)
	}
}

func BenchmarkDrawSpriteScaled(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.sprite(1, 0, 0, 16, 16, 32, 32, 0.0, false, false)
	}
}

func BenchmarkDrawSpriteScaledWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 32, 32, 0.0, false, false)
	}
}
func BenchmarkDrawSpriteScaledWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, 0.0, false, false)
	}
}
func BenchmarkDrawSpriteXFlipped(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.sprite(1, 0, 0, 16, 16, 16, 16, 0.0, true, false)
	}
}

func BenchmarkDrawSpriteXFlippedWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 16, 16, 0.0, true, false)
	}
}

func BenchmarkDrawSpriteXFlippedWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 16, 16, 0.0, true, false)
	}
}
func BenchmarkDrawSpriteYFlippedWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 16, 16, 0.0, false, true)
	}
}
func BenchmarkDrawSpriteYFlipped(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.sprite(1, 0, 0, 16, 16, 16, 16, 0.0, false, true)
	}
}

func BenchmarkDrawSpriteYFlippedWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 16, 16, 0.0, false, true)
	}
}

func BenchmarkDrawSpriteXYFlipped(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.sprite(1, 0, 0, 16, 16, 16, 16, 0.0, true, true)
	}
}

func BenchmarkDrawSpriteXYFlippedWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 16, 16, 0.0, true, true)
	}
}

func BenchmarkDrawSpriteXYFlippedWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 16, 16, 0.0, true, true)
	}
}

func BenchmarkDrawSpriteXFlippedScaled(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.sprite(1, 0, 0, 16, 16, 32, 32, 0.0, true, false)
	}
}

func BenchmarkDrawSpriteXFlippedScaledWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 32, 32, 0.0, true, false)
	}
}

func BenchmarkDrawSpriteXFlippedScaledWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, 0.0, true, false)
	}
}

func BenchmarkDrawSpriteYFlippedScaled(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.sprite(1, 0, 0, 16, 16, 32, 32, 0.0, false, true)
	}
}

func BenchmarkDrawSpriteYFlippedScaledWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 32, 32, 0.0, false, true)
	}
}

func BenchmarkDrawSpriteYFlippedScaledWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, 0.0, false, true)
	}
}
func BenchmarkDrawSpriteXYFlippedScaled(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.sprite(1, 0, 0, 16, 16, 32, 32, 0.0, true, true)
	}
}

func BenchmarkDrawSpriteXYFlippedScaledWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 32, 32, 0.0, true, true)
	}
}

func BenchmarkDrawSpriteXYFlippedScaledWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, 0.0, true, true)
	}
}

func BenchmarkDrawSpriteRotated(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.sprite(1, 0, 0, 16, 16, 32, 32, i%360, false, false)
	}
}

func BenchmarkDrawSpriteRotatedWithMaps(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 32, 32, i%360, false, false)
	}
}

func BenchmarkDrawSpriteRotatedWithCache(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, i%360, false, false)
	}
}

func BenchmarkDrawSpriteRotatedSameAngle(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.sprite(1, 0, 0, 16, 16, 32, 32, 45, false, false)
	}
}

func BenchmarkDrawSpriteRotatedWithMapsSameAngle(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithMaps(1, 0, 0, 16, 16, 32, 32, 45, false, false)
	}
}

func BenchmarkDrawSpriteRotatedWithCacheSameAngle(b *testing.B) {
	Init(PICO8)
	for i := 0; i < b.N; i++ {
		_console.pb.spriteWithCache(1, 0, 0, 16, 16, 32, 32, 45, false, false)
	}
}
