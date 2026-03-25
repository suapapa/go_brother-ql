package backends

import (
	"context"
	"strings"
	"testing"
)

func TestConnect_Unsupported(t *testing.T) {
	_, err := Connect(context.TODO(), "unsupported", "address")
	if err == nil {
		t.Errorf("expected error for unsupported backend")
	}
	if !strings.Contains(err.Error(), "unsupported backend") {
		t.Errorf("expected error to contain 'unsupported backend', got '%s'", err.Error())
	}
}
