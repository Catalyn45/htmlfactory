package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	outDir := flag.String("out", ".", "Path to output directory")
	help := flag.Bool("help", false, "Show usage")
	watch := flag.Bool("watch", false, "keep watching for modified files and template them")

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

	if *watch {
		watcher := newWatcher(args[0], *outDir)
		watcher.run()
	} else {
		files, err := getFiles(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		templateFiles(files, *outDir)
	}
}