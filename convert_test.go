package brother_ql

import (
	"image"
	"image/color"
	"testing"
)

func TestBinarizeImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{0, 0, 0, 255})       // Black -> should be 255 (Dot)
	img.Set(1, 0, color.RGBA{255, 255, 255, 255}) // White -> should be 0
	img.Set(0, 1, color.RGBA{128, 128, 128, 255}) // Gray 50%
	img.Set(1, 1, color.RGBA{200, 200, 200, 255}) // Gray 80%

	// Threshold 70% -> limit = 70/100 * 255 = 178
	res := binarizeImage(img, 70.0)

	if res.GrayAt(0, 0).Y != 255 {
		t.Errorf("expected (0,0) to be 255, got %d", res.GrayAt(0, 0).Y)
	}
	if res.GrayAt(1, 0).Y != 0 {
		t.Errorf("expected (1,0) to be 0, got %d", res.GrayAt(1, 0).Y)
	}
	// 128 < 178 -> Dot (255)
	if res.GrayAt(0, 1).Y != 255 {
		t.Errorf("expected (0,1) to be 255, got %d", res.GrayAt(0, 1).Y)
	}
	// 200 > 178 -> 0
	if res.GrayAt(1, 1).Y != 0 {
		t.Errorf("expected (1,1) to be 0, got %d", res.GrayAt(1, 1).Y)
	}
}

func TestRotateImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255}) // Red
	img.Set(1, 0, color.RGBA{0, 255, 0, 255}) // Green

	// Rotate 90
	res := rotateImage(img, 90)
	if res.Bounds().Dx() != 1 || res.Bounds().Dy() != 2 {
		t.Errorf("expected bounds 1x2, got %dx%d", res.Bounds().Dx(), res.Bounds().Dy())
	}
}
