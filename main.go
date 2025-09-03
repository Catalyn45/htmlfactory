package main

import "flag"

func main() {
	flag.Parse()

	args := flag.Args()
	templateFile(args[0])
}