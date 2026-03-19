# brother-ql : Go pkg for Brother QL label printers

[![Go Reference](https://pkg.go.dev/badge/github.com/suapapa/go_brother-ql.svg)](https://pkg.go.dev/github.com/suapapa/go_brother-ql)

Go port of the popular [brother_ql](https://github.com/pklaus/brother_ql) Python package for creating and sending raster instructions to Brother QL-series label printers.

## Features

- **Multiple Models Support**: Configurable model sizes, compression rules, and features (e.g., QL-500, QL-700, QL-800 series).
- **Multiple Label support**: Supports continuous (endless) tapes and pre-cut (die-cut) labels.
- **Image handling**: Built-in support for resizing, rotation, and dithering for correct raster formatting.
- **Limited backends**: Fully supporting standard `network` (TCP) and `linux_kernel` (`/dev/usb/lpX`) backends without heavy external library dependencies.

---

## Installation

To get the command line tool:
```bash
go install github.com/suapapa/go_brother-ql/cmd/brother_ql@latest
```

To use as a package in your Go application:
```bash
go get github.com/suapapa/go_brother-ql
```

---

## Command Line Interface (CLI)

The CLI aims to replicate the original `brother_ql` Python toolkit commands.

### Global Flags

```bash
brother_ql [command]

Flags:
  -b, --backend string   Backends supporting: network, linux_kernel
  -m, --model string     Select target model (e.g., QL-570, QL-820NWB)
  -p, --printer string   Direct connection identifier (tcp://... or file:///...)
```

### Commands

#### 1. Information Helper
List all supported items by the library:
```bash
brother_ql info models
brother_ql info labels
brother_ql info env
```

#### 2. Print a Label
Convert standard PNG or JPEG images into raster streams and print:
```bash
# Print with network-connected QL-820NWB
brother_ql -b network -p tcp://192.168.1.50 -m QL-820NWB print -l 62 image.png

# Options for print command:
#   -l, --label string    The label size (e.g., 62, 29, 62red)
#   -r, --rotate string   Rotate image ('auto', 0, 90, 180, 270)
#   -t, --threshold float Threshold value percentage (default 70.0)
#   -d, --dither          Enable dithering
#   -c, --compress        Enable output stream compression
#   --no-cut              Don't cut tape after printing
```

---

## Code Example

Example of converting an image and writing out using Go API:

```go
package main

import (
	"fmt"
	"image"
	_ "image/png" // import image formats as needed
	"os"

	"github.com/suapapa/go_brother-ql"
)

func main() {
	// 1. Load image
	f, err := os.Open("label.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	// 2. Create printer
	printer, err := brother_ql.NewLabelPrinter("QL-820NWB", "network", "192.168.1.50")
	if err != nil {
		panic(err)
	}
	defer printer.Close()

	// 3. Setup print options
	opts := brother_ql.NewDefaultOptions("62")

	// 4. Print images
	err = printer.Print([]image.Image{img}, opts)
	if err != nil {
		panic(err)
	}
}
```

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
