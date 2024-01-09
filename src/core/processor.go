package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Processor struct {
	finders []FinderInterface
}

func NewProcessor() *Processor {
	return &Processor{
		finders: []FinderInterface{},
	}
}

func (p *Processor) AddFinder(finder FinderInterface) {
	p.finders = append(p.finders, finder)
}

func (p *Processor) Process(path string) {
	for _, filePath := range p.getFiles(path) {
		p.findSubtitle(filePath)
	}
}

func (p *Processor) findSubtitle(path string) bool {
	for _, finder := range p.finders {
		found, err := finder.Find(path)
		if err != nil {
			fmt.Println(err)
		}
		if found {
			return true
		}
	}
	return false
}

var validFiles = map[string]bool{
	".avi": true,
	".mkv": true,
	".mp4": true,
}

var invalidNames = []string{"sample", "trailer"}

func isValidFile(filePath string) bool {
	if !validFiles[filepath.Ext(filePath)] {
		return false
	}

	for _, substr := range invalidNames {
		if strings.Contains(filePath, substr) {
			return false
		}
	}

	return true
}

func isSubtitleFilePresent(filePath string) bool {
	base := strings.TrimSuffix(filePath, filepath.Ext(filePath))

	if _, err := os.Stat(base + ".srt"); err != nil {
		return false
	}
	return true
}

func (p *Processor) getFiles(path string) []string {

	var files []string

	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !isValidFile(path) || isSubtitleFilePresent(path) {
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files
}
