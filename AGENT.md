# AGENT.md

This file provides context and guidelines for AI coding assistants working on the `go_brother-ql` codebase.

## 📁 Package Architecture

The code is partitioned across a small set of files in the main package with a simple `backends` sub-package.

- **`models.go` / `labels.go`**:
    - Predefined list of supported Printer models and Label types.
    - **To add support for a new model or label**: Simply append a new item in the `initModels()` or `initLabels()` functions.

- **`printer.go`**:
    - High-level API structure containing `LabelPrinter` and `PrintOptions`.
    - Handles managing connections and sending formatted image streams.

- **`raster.go` / `raster_data.go`**:
    - Implements forming the raster instruction bytecode in `BrotherQLRaster`.
    - `AddRasterData` flips the image left-to-right and packs 1-bit pixels into bytes, most-significant bit first.
    - In the inner bytes structure: 255 (or 1-bit) stands for Black / "Print Dot".

- **`convert.go`**:
    - Workflow translating standard images to binarized dots (Padding, Rotation and Thresholding over Gray model mapping).
    - Relies on `github.com/disintegration/imaging` for pixel-level transformations.

- **`packbits.go`**:
    - Native implementation of PackBits RLE encoding used for optionally compressing raster rows before transmission.

- **`backends/backends.go`**:
    - Wrap connection addresses into pure generic `io.ReadWriteCloser` endpoints.
    - Go's building blocks (e.g. `net.Conn`, `*os.File`) satisfy this with no extra boilerplate.

- **`cmd/brother_ql/`**:
    - CLI wrapper replicating the original index interface command behaviors using standard Go options.

---

## 🛠 Guidelines for AI Agents

1. **Keep Imports Simple**:
    - The core relies heavily on Go's standard library. Avoid installing new dependencies unless absolutely required for heavy feature triggers.

2. **Wait for Inversion rules**:
    - Custom threshold setup passes values where **255** stands for Print Dot (Black). Ensure this state is preserved prior to calling `AddRasterData`.

3. **Backend limits**:
    - Respect the restriction to only allow simple byte write-triggers to standard endpoints. Do not enforce heavier dependencies that the Python library had (libusb).

4. **Context and Concurrency**:
    - Pass `context.Context` to all high-level methods (`Print`, `Reconnect`) and backend operations. Use it for timeouts and cancellation.

5. **Error Handling**:
    - Wrap all errors with `fmt.Errorf("%w", err)` to preserve causality. NEVER ignore errors (no `_` assignments) when calling library or system functions.

6. **Maintain Documentation**:
    - Whenever you make modifications to the code, check if high-level structures, CLI behaviors, or backend supports have shifted. If so, you **MUST** update both `AGENT.md` and `README.md` to reflect those changes instantly to keep context up-to-date for future workflows.
