package brother_ql

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// BrotherQLRaster builds the raster data stream for Brother QL printers.
type BrotherQLRaster struct {
	Model            Model
	Data             bytes.Buffer
	pageNumber       int
	CutAtEnd         bool
	Dpi600           bool
	TwoColorPrinting bool
	compression      bool

	// Media properties
	mtype    *byte
	mwidth   *byte
	mlength  *byte
	pquality bool

	// Callback function to report warnings
	onWarning func(string)
}

func newBrotherQLRaster(modelId string) (*BrotherQLRaster, error) {
	m, ok := getModel(modelId)
	if !ok {
		return nil, fmt.Errorf("unknown model: %s", modelId)
	}
	return &BrotherQLRaster{
		Model:    m,
		CutAtEnd: true,
		pquality: true,
	}, nil
}

func (r *BrotherQLRaster) addInitialize() error {
	r.pageNumber = 0
	r.Data.Write([]byte{0x1B, 0x40}) // ESC @
	return nil
}

func (r *BrotherQLRaster) addStatusInformation() error {
	r.Data.Write([]byte{0x1B, 0x69, 0x53}) // ESC i S
	return nil
}

func (r *BrotherQLRaster) addSwitchMode() error {
	if !r.Model.ModeSetting {
		r.warn("Trying to switch the operating mode on a printer that doesn't support the command.")
		return nil
	}
	r.Data.Write([]byte{0x1B, 0x69, 0x61, 0x01}) // ESC i a
	return nil
}

func (r *BrotherQLRaster) addInvalidate() error {
	r.Data.Write(make([]byte, r.Model.NumInvalidateBytes))
	return nil
}

func (r *BrotherQLRaster) setMedia(mtype, width, length byte, quality bool) error {
	r.mtype = &mtype
	r.mwidth = &width
	r.mlength = &length
	r.pquality = quality
	return nil
}

func (r *BrotherQLRaster) addMediaAndQuality(rnumber uint32) error {
	r.Data.Write([]byte{0x1B, 0x69, 0x7A})

	validFlags := byte(0x80)
	if r.mtype != nil {
		validFlags |= 1 << 1
	}
	if r.mwidth != nil {
		validFlags |= 1 << 2
	}
	if r.mlength != nil {
		validFlags |= 1 << 3
	}
	if r.pquality {
		validFlags |= 1 << 6
	}
	r.Data.WriteByte(validFlags)

	if r.mtype != nil {
		r.Data.WriteByte(*r.mtype)
	} else {
		r.Data.WriteByte(0)
	}
	if r.mwidth != nil {
		r.Data.WriteByte(*r.mwidth)
	} else {
		r.Data.WriteByte(0)
	}
	if r.mlength != nil {
		r.Data.WriteByte(*r.mlength)
	} else {
		r.Data.WriteByte(0)
	}

	if err := binary.Write(&r.Data, binary.LittleEndian, rnumber); err != nil {
		return fmt.Errorf("failed to write rnumber: %w", err)
	}

	if r.pageNumber == 0 {
		r.Data.WriteByte(0)
	} else {
		r.Data.WriteByte(1)
	}
	r.Data.WriteByte(0)
	return nil
}

func (r *BrotherQLRaster) addAutocut(autocut bool) error {
	if !r.Model.Cutting {
		r.warn("Trying to call AddAutocut with a printer that doesn't support it")
		return nil
	}
	r.Data.Write([]byte{0x1B, 0x69, 0x4D})
	var val byte
	if autocut {
		val = 1 << 6
	}
	r.Data.WriteByte(val)
	return nil
}

func (r *BrotherQLRaster) addCutEvery(n byte) error {
	if !r.Model.Cutting {
		r.warn("Trying to call AddCutEvery with a printer that doesn't support it")
		return nil
	}
	r.Data.Write([]byte{0x1B, 0x69, 0x41})
	r.Data.WriteByte(n)
	return nil
}

func (r *BrotherQLRaster) addExpandedMode() error {
	if !r.Model.ExpandedMode {
		r.warn("Trying to set expanded mode on a printer that doesn't support it")
		return nil
	}
	r.Data.Write([]byte{0x1B, 0x69, 0x4B})
	var flags byte
	if r.CutAtEnd {
		flags |= 1 << 3
	}
	if r.Dpi600 {
		flags |= 1 << 6
	}
	if r.TwoColorPrinting {
		flags |= 1 << 0
	}
	r.Data.WriteByte(flags)
	return nil
}

func (r *BrotherQLRaster) addMargins(dots uint16) error {
	r.Data.Write([]byte{0x1B, 0x69, 0x64})
	if err := binary.Write(&r.Data, binary.LittleEndian, dots); err != nil {
		return fmt.Errorf("failed to write margins: %w", err)
	}
	return nil
}

func (r *BrotherQLRaster) addCompression(enable bool) error {
	if !r.Model.Compression {
		r.warn("Trying to set compression on a printer that doesn't support it")
		return nil
	}
	r.compression = enable
	r.Data.WriteByte(0x4D) // M
	var val byte
	if enable {
		val = 1 << 1
	}
	r.Data.WriteByte(val)
	return nil
}

func (r *BrotherQLRaster) addPrint(lastPage bool) error {
	if lastPage {
		r.Data.WriteByte(0x1A) // EOF
	} else {
		r.Data.WriteByte(0x0C) // Form Feed
	}
	return nil
}

func (r *BrotherQLRaster) getPixelWidth() int {
	return r.Model.NumberBytesPerRow * 8
}

func (r *BrotherQLRaster) warn(msg string) {
	if r.onWarning != nil {
		r.onWarning(msg)
	}
}
