package brother_ql

import (
	"context"
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

// NewLabelPrinter creates a new LabelPrinter. It attempts an initial connection
// but does not return an error if the connection fails (e.g., printer is off),
// allowing the printer object to be used once it's powered on.
func NewLabelPrinter(ctx context.Context, model, backend, id string) (*LabelPrinter, error) {
	if _, ok := getModel(model); !ok {
		return nil, fmt.Errorf("unknown model: %s", model)
	}

	p := &LabelPrinter{
		model:   model,
		backend: backend,
		id:      id,
	}

	// Attempt initial connection, but don't fail if the printer is off.
	// We use the provided context for the initial attempt.
	p.conn, _ = backends.Connect(ctx, backend, id)

	return p, nil
}

// IsLive checks if the printer connection is available.
func (p *LabelPrinter) IsLive(ctx context.Context) bool {
	return backends.IsLive(ctx, p.backend, p.id)
}

// Reconnect attempts to re-establish the connection to the printer.
func (p *LabelPrinter) Reconnect(ctx context.Context) error {
	if p.conn != nil {
		p.conn.Close()
	}

	conn, err := backends.Connect(ctx, p.backend, p.id)
	if err != nil {
		return fmt.Errorf("reconnect failed: %w", err)
	}
	p.conn = conn
	return nil
}

// Close closes the connection to the printer after a short delay to allow
// the printer to process buffered commands.
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
func (p *LabelPrinter) Print(ctx context.Context, images []image.Image, opts PrintOptions) error {
	qlr, err := newBrotherQLRaster(p.model)
	if err != nil {
		return fmt.Errorf("failed to initialize raster: %w", err)
	}

	data, err := convert(qlr, images, opts.Label, opts.ConvertOptions)
	if err != nil {
		return fmt.Errorf("failed to convert images: %w", err)
	}

	if p.conn == nil {
		if err := p.Reconnect(ctx); err != nil {
			return fmt.Errorf("printer not connected: %w", err)
		}
	}

	_, err = p.conn.Write(data)
	if err != nil {
		// If write fails, try to reconnect and retry once.
		if reconnectErr := p.Reconnect(ctx); reconnectErr == nil {
			_, err = p.conn.Write(data)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to write data to printer: %w", err)
	}

	// Character devices (like /dev/usb/lpX) may return EINVAL on fsync.
	// Since write flush is sufficient, we can ignore the Sync entirely
	// or specifically ignore syscall.EINVAL, but simply avoiding Sync is safer.
	if f, ok := p.conn.(*os.File); ok {
		// Just a no-op to keep the type assertion if ever needed,
		// but we do not call f.Sync() anymore.
		_ = f
	}

	return nil
}
