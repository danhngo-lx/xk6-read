package read

import (
	"os"

	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/read", new(READ))
}

// READ is the k6 extension
type READ struct{}

// Item is the extension file system item representation (file or directory).
type Item interface {
	GetPath() string
	GetContent() any
}

// Directory is the extension directory representation
type Directory struct {
	Item
	Path    string
	Content []Item
}

func (d *Directory) GetPath() string {
	return d.Path
}

func (d *Directory) GetContent() any {
	return d.Content
}

// File is the extension file representation
type File struct {
	Item
	Path    string
	Content any
}

func (f *File) GetPath() string {
	return f.Path
}

func (f *File) GetContent() any {
	return f.Content
}

func (r *READ) ReadDirectory(path string) (Directory, error) {
	directoryEntries, readError := os.ReadDir(path)
	if readError != nil {
		return Directory{}, readError
	} else {
		directory := Directory{Path: path, Content: make([]Item, len(directoryEntries))}

		for i := 0; i < len(directoryEntries); i++ {
			if directoryEntries[i].IsDir() {
				newDirectory, _ := r.ReadDirectory(path + "/" + directoryEntries[i].Name())
				directory.Content[i] = &newDirectory
			} else {
				newFile, _ := r.ReadFile(path + "/" + directoryEntries[i].Name())
				directory.Content[i] = &newFile
			}
		}

		return directory, readError
	}
}

func (*READ) ReadFile(path string, args ...string) (File, error) {
	fileContent, readError := os.ReadFile(path)
	rt := sobek.New()
	ab := rt.NewArrayBuffer(fileContent)

	if readError != nil {
		return File{}, readError
	}

	if len(args) > 0 && args[0] == "b" {
		return File{Path: path, Content: rt.ToValue(&ab)}, nil
	}
	return File{Path: path, Content: string(fileContent)}, nil
}
