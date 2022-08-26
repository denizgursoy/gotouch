package lister

import (
	"bytes"
	"io"
	"os"
)

func NewFileReader(path string) ReadStrategy {
	return &fileReader{
		path: path,
	}
}

type (
	fileReader struct {
		path string
	}
)

func (f *fileReader) ReadProjectStructures() (io.ReadCloser, error) {
	file, err := os.ReadFile(f.path)
	if err != nil {
		return nil, nil
	}
	return io.NopCloser(bytes.NewReader(file)), nil
}
