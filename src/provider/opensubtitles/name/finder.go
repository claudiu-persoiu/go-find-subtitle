package name

import (
	"path/filepath"
)

func Search(path string) (string, error) {
	return "query=" + filepath.Base(path), nil
}
