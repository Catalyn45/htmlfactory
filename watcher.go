package main

import (
	"os"
	"time"
)

type File struct {
	path             string
	size             int64
	lastTimeModified time.Time
}

type Watcher struct {
	timeInSeconds int
	input         string
	output        string
	filesWatched  map[string]File
}

func newWatcher(input string, output string) *Watcher {
	return &Watcher{
		timeInSeconds: 1,
		input:         input,
		output:        output,
		filesWatched:  make(map[string]File),
	}
}

func (self *Watcher) run() error {
	for {
		files, err := getFiles(self.input)
		if err != nil {
			return err
		}

		var filesToTempalte []string
		for _, file := range files {
			info, err := os.Stat(file)
			if err != nil {
				delete(self.filesWatched, file)
				continue
			}

			watchedFile, ok := self.filesWatched[file]
			if !ok ||
			info.Size() != watchedFile.size ||
			info.ModTime() != watchedFile.lastTimeModified {
				filesToTempalte = append(filesToTempalte, file)

				self.filesWatched[file] = File{
					path: file,
					size: info.Size(),
					lastTimeModified: info.ModTime(),
				}
			}
		}

		if len(filesToTempalte) > 0 {
			templateFiles(filesToTempalte, self.output)
		}

		time.Sleep(time.Duration(self.timeInSeconds) * time.Second)
	}
}