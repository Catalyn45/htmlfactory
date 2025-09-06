package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type StringQueue struct {
	container []string
}

func (self *StringQueue) push(item string) bool {
	if self.contains(item) {
		fmt.Println("importing fragment: ", item, " will cause circular dependency")
		return false
	}

	self.container = append(self.container, item)

	return true
}

func (self *StringQueue) pop() string {
	length := len(self.container)
	if length == 0 {
		return ""
	}

	lastIndex := length - 1
	item := self.container[lastIndex]
	self.container = self.container[:lastIndex]

	return item
}

func (self *StringQueue) contains(item string) bool {
	for _, internalItem := range self.container {
		if internalItem == item {
			return true
		}
	}

	return false
}

func getFiles(input string) ([]string, error){
	info, err := os.Stat(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return nil, err
	}

	if !info.IsDir() {
		path, err := filepath.Abs(input)
		return []string{path}, err
	}

	files, err := os.ReadDir(input)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, f := range files {
		if !f.IsDir() &&
		strings.HasSuffix(f.Name(), ".html") &&
		!strings.HasSuffix(f.Name(), ".compiled.html") {
			fullPath := filepath.Join(input, f.Name())
			paths = append(paths, fullPath)
		}
	}

	return paths, nil
}

func templateFiles(files []string, outDir string) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		fmt.Println("done templating all the files")
		fmt.Println("------------------------------")
	}()

	for _, file := range files {
		templater := Templater {
			fileName: file,
			outdir: outDir,
			wg: &wg,
		}

		wg.Add(1)

		fmt.Println("starting templating: ", file)
		go templater.templateFile()
	}
}
