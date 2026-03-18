package brother_ql

import (
	"bytes"
	"testing"
)

func TestPackbits(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{
			name:  "Empty",
			input: []byte{},
			want:  nil,
		},
		{
			name:  "Single byte",
			input: []byte{0x01},
			want:  []byte{0x00, 0x01}, // 0x00 = literal run of length 1
		},
		{
			name:  "Two non-repeating",
			input: []byte{0x01, 0x02},
			want:  []byte{0x01, 0x01, 0x02}, // 0x01 = literal run of length 2
		},
		{
			name:  "Three repeating",
			input: []byte{0x01, 0x01, 0x01},
			want:  []byte{0xFE, 0x01}, // 257 - 3 = 254 (0xFE)
		},
		{
			name:  "Mixed",
			input: []byte{0x01, 0x02, 0x02, 0x02, 0x03},
			want:  []byte{0x00, 0x01, 0xFE, 0x02, 0x00, 0x03},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := packbits(tt.input)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("packbits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackbitsMaxRun(t *testing.T) {
	input := make([]byte, 128)
	for i := range input {
		input[i] = 0xAA
	}
	got := packbits(input)
	want := []byte{0x81, 0xAA} // 257 - 128 = 129 (0x81)
	if !bytes.Equal(got, want) {
		t.Errorf("packbits() = %v, want %v", got, want)
	}
}

func TestPackbitsMaxNonRun(t *testing.T) {
	input := make([]byte, 128)
	for i := range input {
		input[i] = byte(i)
	}
	got := packbits(input)
	want := append([]byte{0x7F}, input...) // 128 - 1 = 127 (0x7F)
	if !bytes.Equal(got, want) {
		t.Errorf("packbits() = %v, want %v", got, want)
	}
}
