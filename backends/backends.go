package backends

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

// normalizeAddress handles prefix trimming and default port appending for network/file addresses.
func normalizeAddress(backendType, address string) string {
	switch backendType {
	case "network":
		address = strings.TrimPrefix(address, "tcp://")
		if !strings.Contains(address, ":") {
			address += ":9100"
		}
	case "linux_kernel":
		address = strings.TrimPrefix(address, "file://")
	}
	return address
}

// Connect establishes a connection to the printer based on the backend type.
// Supported backend types are "network" (e.g., "192.168.1.100" or "tcp://192.168.1.100:9100")
// and "linux_kernel" (e.g., "/dev/usb/lp0").
func Connect(ctx context.Context, backendType, address string) (io.ReadWriteCloser, error) {
	address = normalizeAddress(backendType, address)

	switch backendType {
	case "network":
		var d net.Dialer
		conn, err := d.DialContext(ctx, "tcp", address)
		if err != nil {
			return nil, fmt.Errorf("network connection failed to %s: %w", address, err)
		}
		return conn, nil
	case "linux_kernel":
		// os.OpenFile doesn't support context directly; we'd need to use FD or polling if async is needed.
		// For a simple local file backend, we check if the context is already cancelled.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		f, err := os.OpenFile(address, os.O_RDWR, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to open device %s: %w", address, err)
		}
		return f, nil
	default:
		return nil, fmt.Errorf("unsupported backend type: %q", backendType)
	}
}

// IsLive checks if the printer is reachable at the given address.
func IsLive(ctx context.Context, backendType, address string) bool {
	address = normalizeAddress(backendType, address)

	switch backendType {
	case "network":
		var d net.Dialer
		// Use the shorter of the context-provided timeout or a default 2s.
		deadline := time.Now().Add(2 * time.Second)
		if ddl, ok := ctx.Deadline(); ok && ddl.Before(deadline) {
			deadline = ddl
		}
		
		subCtx, cancel := context.WithDeadline(ctx, deadline)
		defer cancel()

		conn, err := d.DialContext(subCtx, "tcp", address)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	case "linux_kernel":
		f, err := os.OpenFile(address, os.O_RDWR, 0)
		if err != nil {
			// If the file doesn't exist, the printer is off or disconnected.
			if os.IsNotExist(err) {
				return false
			}
			// If it exists but we can't open it (e.g., permission, busy), 
			// it's considered "on" but unavailable.
			return true
		}
		f.Close()
		return true
	default:
		return false
	}
}

