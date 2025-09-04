package main

import "flag"
import "os"
import "fmt"
import "strings"
import "path/filepath"

func main() {
	outDir := flag.String("out", ".", "Path to output directory")
	help := flag.Bool("help", false, "Show usage")

	// Customize usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <input_file_or_directory> [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: missing input")
		flag.Usage()
		return
	}

	info, err := os.Stat(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	if !info.IsDir() {
		templateFile(args[0], outDir)
		return
	}

	files, err := os.ReadDir(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for _, f := range files {
		if !f.IsDir() &&
		strings.HasSuffix(f.Name(), ".html") &&
		!strings.HasSuffix(f.Name(), ".compiled.html") {
			fullPath := filepath.Join(args[0], f.Name())
			templateFile(fullPath, outDir)
		}
	}
}