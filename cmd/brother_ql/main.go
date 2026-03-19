package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	brother_ql "github.com/suapapa/go_brother-ql"
)

var (
	backend string
	model   string
	printer string
	debug   bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "brother_ql",
		Short: "Command line interface for the brother_ql Go package",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				log.Println("Debug mode enabled")
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&model, "model", "m", os.Getenv("BROTHER_QL_MODEL"), "Choices: QL-500, QL-550, ...")
	rootCmd.PersistentFlags().StringVarP(&backend, "backend", "b", os.Getenv("BROTHER_QL_BACKEND"), "Choices: network, linux_kernel")
	rootCmd.PersistentFlags().StringVarP(&printer, "printer", "p", os.Getenv("BROTHER_QL_PRINTER"), "The identifier for the printer")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")

	// discover
	var discoverCmd = &cobra.Command{
		Use:   "discover",
		Short: "find connected label printers",
		Run: func(cmd *cobra.Command, args []string) {
			discoverAndList(backend)
		},
	}
	rootCmd.AddCommand(discoverCmd)

	// info
	var infoCmd = &cobra.Command{
		Use:   "info",
		Short: "list available labels, models etc.",
	}
	rootCmd.AddCommand(infoCmd)

	var modelsCmd = &cobra.Command{
		Use:   "models",
		Short: "List the choices for --model",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Supported models:")
			for _, m := range brother_ql.AllModels {
				fmt.Println(" ", m.Identifier)
			}
		},
	}
	infoCmd.AddCommand(modelsCmd)

	var labelsCmd = &cobra.Command{
		Use:   "labels",
		Short: "List the choices for --label",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Supported labels:")
			for _, l := range brother_ql.AllLabels {
				fmt.Printf("  %s (%dx%d mm)\n", l.Identifier, l.TapeSize[0], l.TapeSize[1])
			}
		},
	}
	infoCmd.AddCommand(labelsCmd)

	var envCmd = &cobra.Command{
		Use:   "env",
		Short: "print debug info about running environment",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("\n##################\n\n")
			fmt.Println("Information about the running environment of brother_ql.")
			fmt.Println("About the computer:")
			fmt.Printf("  * OS: %s\n", runtime.GOOS)
			fmt.Printf("  * Arch: %s\n", runtime.GOARCH)
			fmt.Printf("  * Go Version: %s\n", runtime.Version())
			fmt.Print("\n##################\n\n")
		},
	}
	infoCmd.AddCommand(envCmd)

	// print
	var labelArg string
	var rotate string
	var threshold float64
	var dither bool
	var ditherAlgo string
	var compress bool
	var red bool
	var dpi600 bool
	var hq bool
	var noCut bool

	var printCmd = &cobra.Command{
		Use:   "print [IMAGE ...]",
		Short: "Print a label",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runPrint(args, labelArg, rotate, threshold, dither, ditherAlgo, compress, red, dpi600, hq, noCut)
		},
	}
	printCmd.Flags().StringVarP(&labelArg, "label", "l", os.Getenv("BROTHER_QL_LABEL"), "The label size")
	printCmd.Flags().StringVarP(&rotate, "rotate", "r", "auto", "Rotate image by degrees (0, 90, 180, 270, auto)")
	printCmd.Flags().Float64VarP(&threshold, "threshold", "t", 70.0, "Threshold value (%) to discriminate black & white")
	printCmd.Flags().BoolVarP(&dither, "dither", "d", false, "Enable dithering")
	printCmd.Flags().StringVar(&ditherAlgo, "dither-algo", "floyd_steinberg", "Dithering algorithm (atkinson, burkes, stucki, sierra2, sierra3, sierralite, floyd_steinberg)")
	printCmd.Flags().BoolVarP(&compress, "compress", "c", false, "Enable compression")
	printCmd.Flags().BoolVar(&red, "red", false, "Create label for black/red/white tape")
	printCmd.Flags().BoolVar(&dpi600, "600dpi", false, "Print with 600x300 dpi available on some models")
	printCmd.Flags().BoolVar(&hq, "hq", true, "Print with high quality")
	printCmd.Flags().BoolVar(&noCut, "no-cut", false, "Don't cut tape after printing")

	rootCmd.AddCommand(printCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func discoverAndList(backend string) {
	if backend == "" {
		backend = "linux_kernel"
	}
	if backend == "linux_kernel" {
		matches, err := filepath.Glob("/dev/usb/lp*")
		if err != nil {
			fmt.Println("Error reading device paths:", err)
			return
		}
		for _, m := range matches {
			fmt.Printf("file://%s\n", m)
		}
	} else if backend == "network" {
		fmt.Println("Discover for network is not implemented yet.")
	} else {
		fmt.Println("Unknown backend:", backend)
	}
}

func runPrint(images []string, label, rotate string, threshold float64, dither bool, ditherAlgo string, compress, red, dpi600, hq, noCut bool) {
	var parsedImages []image.Image
	for _, imgPath := range images {
		f, err := os.Open(imgPath)
		if err != nil {
			fmt.Printf("Error opening image %s: %v\n", imgPath, err)
			return
		}
		defer f.Close()
		img, _, err := image.Decode(f)
		if err != nil {
			fmt.Printf("Error decoding image %s: %v\n", imgPath, err)
			return
		}
		parsedImages = append(parsedImages, img)
	}

	if printer == "" {
		fmt.Println("Error: --printer flag is required to send to device")
		return
	}

	brd, err := brother_ql.NewLabelPrinter(model, backend, printer)
	if err != nil {
		fmt.Println("Error creating printer:", err)
		return
	}
	defer brd.Close()

	opts := brother_ql.PrintOptions{
		Label: label,
		ConvertOptions: brother_ql.ConvertOptions{
			Cut:        !noCut,
			Dither:     dither,
			DitherAlgo: ditherAlgo,
			Compress:   compress,
			Red:        red,
			Rotate:     rotate,
			Dpi600:     dpi600,
			Hq:         hq,
			Threshold:  threshold,
		},
	}

	if err := brd.Print(parsedImages, opts); err != nil {
		fmt.Println("Error printing:", err)
		return
	}

	if debug {
		fmt.Println("Label printed successfully")
	}
}
