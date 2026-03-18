package brother_ql

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type BrotherQLRaster struct {
	Model              Model
	Data               bytes.Buffer
	pageNumber         int
	CutAtEnd           bool
	Dpi600             bool
	TwoColorPrinting   bool
	compression        bool
	ExceptionOnWarning bool

	// Media properties
	mtype   *byte
	mwidth  *byte
	mlength *byte
	pquality bool
}

func NewBrotherQLRaster(modelId string) (*BrotherQLRaster, error) {
	m, ok := GetModel(modelId)
	if !ok {
		return nil, fmt.Errorf("unknown model: %s", modelId)
	}
	return &BrotherQLRaster{
		Model:    m,
		CutAtEnd: true,
		pquality: true,
	}, nil
}

func (r *BrotherQLRaster) AddInitialize() {
	r.pageNumber = 0
	r.Data.Write([]byte{0x1B, 0x40}) // ESC @
}

func (r *BrotherQLRaster) AddStatusInformation() {
	r.Data.Write([]byte{0x1B, 0x69, 0x53}) // ESC i S
}

func (r *BrotherQLRaster) AddSwitchMode() {
	if !r.Model.ModeSetting {
		r.warn("Trying to switch the operating mode on a printer that doesn't support the command.")
		return
	}
	r.Data.Write([]byte{0x1B, 0x69, 0x61, 0x01}) // ESC i a
}

func (r *BrotherQLRaster) AddInvalidate() {
	r.Data.Write(make([]byte, r.Model.NumInvalidateBytes))
}

func (r *BrotherQLRaster) SetMedia(mtype, width, length byte, quality bool) {
	r.mtype = &mtype
	r.mwidth = &width
	r.mlength = &length
	r.pquality = quality
}

func (r *BrotherQLRaster) AddMediaAndQuality(rnumber uint32) {
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

	binary.Write(&r.Data, binary.LittleEndian, rnumber)

	if r.pageNumber == 0 {
		r.Data.WriteByte(0)
	} else {
		r.Data.WriteByte(1)
	}
	r.Data.WriteByte(0)
}

func (r *BrotherQLRaster) AddAutocut(autocut bool) {
	if !r.Model.Cutting {
		r.warn("Trying to call AddAutocut with a printer that doesn't support it")
		return
	}
	r.Data.Write([]byte{0x1B, 0x69, 0x4D})
	var val byte
	if autocut {
		val = 1 << 6
	}
	r.Data.WriteByte(val)
}

func (r *BrotherQLRaster) AddCutEvery(n byte) {
	if !r.Model.Cutting {
		r.warn("Trying to call AddCutEvery with a printer that doesn't support it")
		return
	}
	r.Data.Write([]byte{0x1B, 0x69, 0x41})
	r.Data.WriteByte(n)
}

func (r *BrotherQLRaster) AddExpandedMode() {
	if !r.Model.ExpandedMode {
		r.warn("Trying to set expanded mode on a printer that doesn't support it")
		return
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
}

func (r *BrotherQLRaster) AddMargins(dots uint16) {
	r.Data.Write([]byte{0x1B, 0x69, 0x64})
	binary.Write(&r.Data, binary.LittleEndian, dots)
}

func (r *BrotherQLRaster) AddCompression(enable bool) {
	if !r.Model.Compression {
		r.warn("Trying to set compression on a printer that doesn't support it")
		return
	}
	r.compression = enable
	r.Data.WriteByte(0x4D) // M
	var val byte
	if enable {
		val = 1 << 1
	}
	r.Data.WriteByte(val)
}

func (r *BrotherQLRaster) AddPrint(lastPage bool) {
	if lastPage {
		r.Data.WriteByte(0x1A) // EOF
	} else {
		r.Data.WriteByte(0x0C) // Form Feed
	}
}

func (r *BrotherQLRaster) GetPixelWidth() int {
	return r.Model.NumberBytesPerRow * 8
}

func (r *BrotherQLRaster) warn(msg string) {
	if r.ExceptionOnWarning {
		panic(msg)
	} else {
		// Log warning
	}
}
