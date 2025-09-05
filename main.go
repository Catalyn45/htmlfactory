package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

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

	var wg sync.WaitGroup
	defer wg.Wait()

	if !info.IsDir() {
		templater := Templater {
			fileName: args[0],
			outdir: *outDir,
			wg: &wg,
		}

		go templater.templateFile()
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

			templater := Templater {
				fileName: fullPath,
				outdir: *outDir,
				wg: &wg,
			}

			wg.Add(1)
			go templater.templateFile()
		}
	}
}