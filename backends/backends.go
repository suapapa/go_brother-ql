package backends

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func Connect(backendType, address string) (io.ReadWriteCloser, error) {
	switch backendType {
	case "network":
		if strings.HasPrefix(address, "tcp://") {
			address = address[6:]
		}
		if !strings.Contains(address, ":") {
			address += ":9100"
		}
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return nil, err
		}
		return conn, nil
	case "linux_kernel":
		if strings.HasPrefix(address, "file://") {
			address = address[7:]
		}
		f, err := os.OpenFile(address, os.O_RDWR, 0)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return nil, fmt.Errorf("unsupported backend: %s", backendType)
	}
}
