package brother_ql

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/disintegration/imaging"
)

type ConvertOptions struct {
	Cut       bool
	Dither    bool
	Compress  bool
	Red       bool
	Rotate    string // "auto", "0", "90", "180", "270"
	Dpi600    bool
	Hq        bool
	Threshold float64
}

func Convert(qlr *BrotherQLRaster, images []image.Image, labelId string, opts ConvertOptions) ([]byte, error) {
	labelSpecs, ok := GetLabel(labelId)
	if !ok {
		return nil, fmt.Errorf("unknown label size: %s", labelId)
	}

	dotsPrintable := labelSpecs.DotsPrintable
	rightMarginDots := labelSpecs.OffsetR
	rightMarginDots += qlr.Model.AdditionalOffsetR
	devicePixelWidth := qlr.GetPixelWidth()

	if opts.Threshold == 0 {
		opts.Threshold = 70.0 // Default
	}

	if opts.Red && !qlr.Model.TwoColor {
		return nil, fmt.Errorf("printing in red is not supported with the selected model")
	}

	qlr.AddSwitchMode()
	qlr.AddInvalidate()
	qlr.AddInitialize()
	qlr.AddSwitchMode()

	for _, im := range images {
		// Convert to Gray or RGB is handled by imaging
		var processedIm image.Image = im

		dotsExpected := dotsPrintable
		if opts.Dpi600 {
			dotsExpected[0] *= 2
			dotsExpected[1] *= 2
		}

		if labelSpecs.FormFactor == Endless || labelSpecs.FormFactor == PtouchEndless {
			if opts.Rotate != "auto" && opts.Rotate != "0" {
				angle := 0
				fmt.Sscanf(opts.Rotate, "%d", &angle)
				processedIm = rotateImage(processedIm, angle)
			}
			if opts.Dpi600 {
				sz := processedIm.Bounds().Size()
				processedIm = imaging.Resize(processedIm, sz.X/2, sz.Y, imaging.Lanczos)
			}
			sz := processedIm.Bounds().Size()
			if sz.X != dotsPrintable[0] {
				hsize := int(float64(dotsPrintable[0]) / float64(sz.X) * float64(sz.Y))
				processedIm = imaging.Resize(processedIm, dotsPrintable[0], hsize, imaging.Lanczos)
				log.Println("Need to resize the image...")
			}
			if processedIm.Bounds().Dx() < devicePixelWidth {
				bg := imaging.New(devicePixelWidth, processedIm.Bounds().Dy(), color.White)
				pasteX := devicePixelWidth - processedIm.Bounds().Dx() - rightMarginDots
				processedIm = imaging.Paste(bg, processedIm, image.Point{X: pasteX, Y: 0})
			}
		} else if labelSpecs.FormFactor == DieCut || labelSpecs.FormFactor == RoundDieCut {
			if opts.Rotate == "auto" {
				sz := processedIm.Bounds().Size()
				if sz.X == dotsExpected[1] && sz.Y == dotsExpected[0] {
					processedIm = imaging.Rotate90(processedIm)
				}
			} else if opts.Rotate != "0" && opts.Rotate != "" {
				angle := 0
				fmt.Sscanf(opts.Rotate, "%d", &angle)
				processedIm = rotateImage(processedIm, angle)
			}
			sz := processedIm.Bounds().Size()
			if sz.X != dotsExpected[0] || sz.Y != dotsExpected[1] {
				return nil, fmt.Errorf("bad image dimensions: %v. Expecting: %v", sz, dotsExpected)
			}
			if opts.Dpi600 {
				processedIm = imaging.Resize(processedIm, sz.X/2, sz.Y, imaging.Lanczos)
			}
			bg := imaging.New(devicePixelWidth, dotsExpected[1], color.White)
			pasteX := devicePixelWidth - processedIm.Bounds().Dx() - rightMarginDots
			processedIm = imaging.Paste(bg, processedIm, image.Point{X: pasteX, Y: 0})
		}

		var binIm *image.Gray
		if opts.Dither {
			binIm = ditherImage(processedIm)
		} else {
			binIm = binarizeImage(processedIm, opts.Threshold)
		}

		qlr.AddStatusInformation()
		if labelSpecs.FormFactor == DieCut || labelSpecs.FormFactor == RoundDieCut {
			qlr.SetMedia(0x0B, byte(labelSpecs.TapeSize[0]), byte(labelSpecs.TapeSize[1]), opts.Hq)
		} else if labelSpecs.FormFactor == Endless {
			qlr.SetMedia(0x0A, byte(labelSpecs.TapeSize[0]), 0, opts.Hq)
		} else if labelSpecs.FormFactor == PtouchEndless {
			qlr.SetMedia(0x00, byte(labelSpecs.TapeSize[0]), 0, opts.Hq)
		}
		qlr.AddMediaAndQuality(uint32(binIm.Bounds().Dy()))

		if opts.Cut && qlr.Model.Cutting {
			qlr.AddAutocut(true)
			qlr.AddCutEvery(1)
		}

		qlr.CutAtEnd = opts.Cut
		qlr.Dpi600 = opts.Dpi600
		qlr.TwoColorPrinting = opts.Red
		qlr.AddExpandedMode()

		qlr.AddMargins(uint16(labelSpecs.FeedMargin))
		if opts.Compress && qlr.Model.Compression {
			qlr.AddCompression(true)
		}

		qlr.AddRasterData(binIm)
		qlr.AddPrint(true) // Wait, for multiples images we might need false then true on last
	}

	return qlr.Data.Bytes(), nil
}

func rotateImage(img image.Image, angle int) image.Image {
	switch angle {
	case 90:
		return imaging.Rotate90(img)
	case 180:
		return imaging.Rotate180(img)
	case 270:
		return imaging.Rotate270(img)
	}
	return img
}

func binarizeImage(img image.Image, threshold float64) *image.Gray {
	gray := imaging.Grayscale(img)
	res := image.NewGray(gray.Bounds())
	limit := uint8(threshold / 100.0 * 255.0)
	for y := gray.Bounds().Min.Y; y < gray.Bounds().Max.Y; y++ {
		for x := gray.Bounds().Min.X; x < gray.Bounds().Max.X; x++ {
			c := gray.At(x, y).(color.Gray)
			if c.Y <= limit {
				res.SetGray(x, y, color.Gray{Y: 255}) // Dot
			} else {
				res.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
	return res
}

func ditherImage(img image.Image) *image.Gray {
	dst := image.NewPaletted(img.Bounds(), []color.Color{color.Gray{0}, color.Gray{255}})
	draw.FloydSteinberg.Draw(dst, dst.Bounds(), img, image.Point{})
	res := image.NewGray(dst.Bounds())
	for y := dst.Bounds().Min.Y; y < dst.Bounds().Max.Y; y++ {
		for x := dst.Bounds().Min.X; x < dst.Bounds().Max.X; x++ {
			c := dst.At(x, y)
			// dst.At returns color from palette.
			// Index 0 is color.Gray{0} (Black) -> Dot
			// Index 1 is color.Gray{255} (White) -> No Dot
			r, _, _, _ := c.RGBA()
			if r < 32768 { // Black
				res.SetGray(x, y, color.Gray{Y: 255})
			} else {
				res.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
	return res
}
