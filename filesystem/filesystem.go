package filesystem

import (
	"net/http"
	"os"
	"strings"
)

type mf struct {
	http.File
}

func (f mf) Readdir(n int) (fis []os.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			fis = append(fis, file)
		}
	}
	return
}

// Fs is a custom filesystem
type FS struct {
	http.FileSystem
}

func (fs FS) Open(path string) (http.File, error) {
	file, err := fs.FileSystem.Open(path)
	return mf{file}, err
}
