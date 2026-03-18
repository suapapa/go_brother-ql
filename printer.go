package brother_ql

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/suapapa/go_brother-ql/backends"
)

// PrintOptions contains options for printing.
type PrintOptions struct {
	Label string // The label size identifier (e.g., "62", "29x90")
	ConvertOptions
}

// NewDefaultOptions creates a default PrintOptions with recommended settings.
func NewDefaultOptions(label string) PrintOptions {
	return PrintOptions{
		Label: label,
		ConvertOptions: ConvertOptions{
			Cut:       true,
			Rotate:    "auto",
			Threshold: 70.0,
			Hq:        true,
		},
	}
}

// LabelPrinter manages connection and printing to a Brother QL printer.
type LabelPrinter struct {
	model   string
	backend string
	id      string
}

// NewLabelPrinter creates a new LabelPrinter.
func NewLabelPrinter(model, backend, id string) (*LabelPrinter, error) {
	if _, ok := getModel(model); !ok {
		return nil, fmt.Errorf("unknown model: %s", model)
	}
	return &LabelPrinter{
		model:   model,
		backend: backend,
		id:      id,
	}, nil
}

// Print converts and sends the images to the printer.
func (p *LabelPrinter) Print(images []image.Image, opts PrintOptions) error {
	qlr, err := newBrotherQLRaster(p.model)
	if err != nil {
		return err
	}

	data, err := convert(qlr, images, opts.Label, opts.ConvertOptions)
	if err != nil {
		return err
	}

	conn, err := backends.Connect(p.backend, p.id)
	if err != nil {
		return fmt.Errorf("failed to connect to %s backend at %s: %v", p.backend, p.id, err)
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to printer: %v", err)
	}

	if f, ok := conn.(*os.File); ok {
		f.Sync()
	}
	// Give the printer some time to process or buffer the cut command
	// before the file descriptor is closed.
	time.Sleep(500 * time.Millisecond)

	return nil
}
