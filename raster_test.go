package brother_ql

import (
	"bytes"
	"testing"
)

func TestBrotherQLRaster_Commands(t *testing.T) {
	r, err := newBrotherQLRaster("QL-580N")
	if err != nil {
		t.Fatalf("failed to create raster: %v", err)
	}

	// Test addInitialize
	if err := r.addInitialize(); err != nil {
		t.Errorf("addInitialize() failed: %v", err)
	}
	want := []byte{0x1B, 0x40}
	if !bytes.Equal(r.Data.Bytes(), want) {
		t.Errorf("addInitialize() = %v, want %v", r.Data.Bytes(), want)
	}
	r.Data.Reset()

	// Test addStatusInformation
	if err := r.addStatusInformation(); err != nil {
		t.Errorf("addStatusInformation() failed: %v", err)
	}
	want = []byte{0x1B, 0x69, 0x53}
	if !bytes.Equal(r.Data.Bytes(), want) {
		t.Errorf("addStatusInformation() = %v, want %v", r.Data.Bytes(), want)
	}
	r.Data.Reset()

	// Test addSwitchMode
	if err := r.addSwitchMode(); err != nil {
		t.Errorf("addSwitchMode() failed: %v", err)
	}
	want = []byte{0x1B, 0x69, 0x61, 0x01}
	if !bytes.Equal(r.Data.Bytes(), want) {
		t.Errorf("addSwitchMode() = %v, want %v", r.Data.Bytes(), want)
	}
	r.Data.Reset()

	// Test addInvalidate
	if err := r.addInvalidate(); err != nil {
		t.Errorf("addInvalidate() failed: %v", err)
	}
	if len(r.Data.Bytes()) != r.Model.NumInvalidateBytes {
		t.Errorf("addInvalidate() length = %d, want %d", len(r.Data.Bytes()), r.Model.NumInvalidateBytes)
	}
	r.Data.Reset()

	// Test setMedia and addMediaAndQuality
	if err := r.setMedia(11, 62, 0, true); err != nil {
		t.Errorf("setMedia() failed: %v", err)
	}
	if err := r.addMediaAndQuality(0); err != nil {
		t.Errorf("addMediaAndQuality() failed: %v", err)
	}
	want = []byte{0x1B, 0x69, 0x7A, 0xCE, 11, 62, 0, 0, 0, 0, 0, 0, 0}
	if !bytes.Equal(r.Data.Bytes(), want) {
		t.Errorf("addMediaAndQuality() = %x, want %x", r.Data.Bytes(), want)
	}
	r.Data.Reset()

	// Test addAutocut
	if err := r.addAutocut(true); err != nil {
		t.Errorf("addAutocut() failed: %v", err)
	}
	want = []byte{0x1B, 0x69, 0x4D, 0x40}
	if !bytes.Equal(r.Data.Bytes(), want) {
		t.Errorf("addAutocut(true) = %x, want %x", r.Data.Bytes(), want)
	}
	r.Data.Reset()

	// Test addCutEvery
	if err := r.addCutEvery(3); err != nil {
		t.Errorf("addCutEvery() failed: %v", err)
	}
	want = []byte{0x1B, 0x69, 0x41, 0x03}
	if !bytes.Equal(r.Data.Bytes(), want) {
		t.Errorf("addCutEvery() = %x, want %x", r.Data.Bytes(), want)
	}
	r.Data.Reset()

	// Test addExpandedMode
	r.CutAtEnd = true
	r.Dpi600 = true
	r.TwoColorPrinting = false
	if err := r.addExpandedMode(); err != nil {
		t.Errorf("addExpandedMode() failed: %v", err)
	}
	want = []byte{0x1B, 0x69, 0x4B, 0x48}
	if !bytes.Equal(r.Data.Bytes(), want) {
		t.Errorf("addExpandedMode() = %x, want %x", r.Data.Bytes(), want)
	}
	r.Data.Reset()

	// Test addMargins
	if err := r.addMargins(35); err != nil {
		t.Errorf("addMargins() failed: %v", err)
	}
	want = []byte{0x1B, 0x69, 0x64, 0x23, 0x00}
	if !bytes.Equal(r.Data.Bytes(), want) {
		t.Errorf("addMargins() = %x, want %x", r.Data.Bytes(), want)
	}
	r.Data.Reset()

	// Test addCompression
	if err := r.addCompression(true); err != nil {
		t.Errorf("addCompression() failed: %v", err)
	}
	want = []byte{0x4D, 0x02}
	if !bytes.Equal(r.Data.Bytes(), want) {
		t.Errorf("addCompression() = %x, want %x", r.Data.Bytes(), want)
	}
	r.Data.Reset()

	// Test addPrint
	if err := r.addPrint(true); err != nil {
		t.Errorf("addPrint() failed: %v", err)
	}
	want = []byte{0x1A}
	if !bytes.Equal(r.Data.Bytes(), want) {
		t.Errorf("addPrint(true) = %x, want %x", r.Data.Bytes(), want)
	}
	r.Data.Reset()
}

func TestBrotherQLRaster_OnWarning(t *testing.T) {
	r, err := newBrotherQLRaster("QL-500")
	if err != nil {
		t.Fatalf("failed to create raster: %v", err)
	}

	var warningMsg string
	r.onWarning = func(msg string) {
		warningMsg = msg
	}

	if err := r.addCompression(true); err != nil {
		t.Errorf("addCompression() failed: %v", err)
	}

	if warningMsg == "" {
		t.Fatalf("expected warning callback to be called")
	}
	wantMsg := "Trying to set compression on a printer that doesn't support it"
	if warningMsg != wantMsg {
		t.Errorf("warningMsg = %q, want %q", warningMsg, wantMsg)
	}
}
