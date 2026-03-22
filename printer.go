package brother_ql

import (
	"fmt"
	"image"
	"io"
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
	conn    io.ReadWriteCloser
}

// NewLabelPrinter creates a new LabelPrinter.
func NewLabelPrinter(model, backend, id string) (*LabelPrinter, error) {
	if _, ok := getModel(model); !ok {
		return nil, fmt.Errorf("unknown model: %s", model)
	}

	conn, err := backends.Connect(backend, id)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s backend at %s: %v", backend, id, err)
	}

	return &LabelPrinter{
		model:   model,
		backend: backend,
		id:      id,
		conn:    conn,
	}, nil
}

// IsLive checks if the printer connection is available.
func (p *LabelPrinter) IsLive() bool {
	return backends.IsLive(p.backend, p.id)
}

// Reconnect attempts to re-establish the connection to the printer.
func (p *LabelPrinter) Reconnect() error {
	if p.conn != nil {
		p.conn.Close()
	}

	conn, err := backends.Connect(p.backend, p.id)
	if err != nil {
		return err
	}
	p.conn = conn
	return nil
}


// Close closes the connection to the printer.
func (p *LabelPrinter) Close() error {
	if p.conn != nil {
		// Give the printer some time to process or buffer the cut command
		// before the file descriptor is closed.
		time.Sleep(500 * time.Millisecond)
		return p.conn.Close()
	}
	return nil
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

	_, err = p.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to printer: %v", err)
	}

	if f, ok := p.conn.(*os.File); ok {
		f.Sync()
	}

	return nil
}
