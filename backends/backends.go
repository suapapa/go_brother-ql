package backends

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// Connect establishes a connection to the printer based on the backend type.
// Supported backend types are "network" (e.g., "192.168.1.100" or "tcp://192.168.1.100:9100")
// and "linux_kernel" (e.g., "/dev/usb/lp0").
func Connect(backendType, address string) (io.ReadWriteCloser, error) {
	switch backendType {
	case "network":
		address = strings.TrimPrefix(address, "tcp://")
		if !strings.Contains(address, ":") {
			address += ":9100"
		}
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return nil, err
		}
		return conn, nil
	case "linux_kernel":
		address = strings.TrimPrefix(address, "file://")
		f, err := os.OpenFile(address, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return nil, fmt.Errorf("unsupported backend: %s", backendType)
	}
}
