package brother_ql

import (
	"testing"
)

func TestGetModel(t *testing.T) {
	m, ok := getModel("QL-800")
	if !ok {
		t.Errorf("expected to find QL-800")
	}
	if m.Identifier != "QL-800" {
		t.Errorf("expected QL-800, got %s", m.Identifier)
	}
	if !m.TwoColor {
		t.Errorf("expected QL-800 to support TwoColor")
	}

	_, ok = getModel("NonExistent")
	if ok {
		t.Errorf("expected to not find NonExistent")
	}
}

func TestGetLabel(t *testing.T) {
	l, ok := getLabel("62")
	if !ok {
		t.Errorf("expected to find label 62")
	}
	if l.Identifier != "62" {
		t.Errorf("expected 62, got %s", l.Identifier)
	}
	if l.FormFactor != Endless {
		t.Errorf("expected 62 to be Endless")
	}

	_, ok = getLabel("NonExistent")
	if ok {
		t.Errorf("expected to not find NonExistent")
	}
}
