package brother_ql

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
)

func (r *BrotherQLRaster) addRasterData(images ...image.Image) error {
	var frames [][]byte
	var width, height int

	if len(images) > 0 {
		bounds := images[0].Bounds()
		width = bounds.Dx()
		height = bounds.Dy()
	}

	if width != r.getPixelWidth() {
		return fmt.Errorf("wrong pixel width: %d, expected %d", width, r.getPixelWidth())
	}

	for _, img := range images {
		bounds := img.Bounds()
		if bounds.Dx() != width || bounds.Dy() != height {
			return fmt.Errorf("images must have same dimensions")
		}

		rowLen := width / 8
		frame := make([]byte, rowLen*height)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := 0; x < width; x++ {
				// Flip left-to-right
				pixelX := bounds.Max.X - 1 - x

				c := img.At(pixelX, y)
				var v byte
				if g, ok := c.(color.Gray); ok {
					if g.Y > 127 { // Assumes inverted image: 255 means Dot
						v = 1
					}
				} else {
					r, g, b, _ := c.RGBA()
					// Simplistic: if any color has high value, it might be dot
					// Invert logic should be in convert.go
					if r > 32768 || g > 32768 || b > 32768 {
						v = 1
					}
				}

				if v == 1 {
					byteIdx := (y-bounds.Min.Y)*rowLen + x/8
					bitIdx := 7 - (x % 8) // MSB first
					frame[byteIdx] |= (1 << bitIdx)
				}
			}
		}
		frames = append(frames, frame)
	}

	rowLen := width / 8
	var fileStr bytes.Buffer

	for start := 0; start < len(frames[0]); start += rowLen {
		for i, frame := range frames {
			row := frame[start : start+rowLen]
			if r.compression {
				row = packbits(row)
			}
			translen := len(row)

			if r.Model.Identifier[:2] == "PT" {
				fileStr.WriteByte(0x47)
				if err := binary.Write(&fileStr, binary.LittleEndian, uint16(translen)); err != nil {
					return fmt.Errorf("failed to write PT translen: %w", err)
				}
			} else {
				if len(images) > 1 {
					if i == 0 {
						fileStr.Write([]byte{0x77, 0x01})
					} else {
						fileStr.Write([]byte{0x77, 0x02})
					}
				} else {
					fileStr.Write([]byte{0x67, 0x00})
				}
				fileStr.WriteByte(byte(translen))
			}
			fileStr.Write(row)
		}
	}

	r.Data.Write(fileStr.Bytes())
	return nil
}
