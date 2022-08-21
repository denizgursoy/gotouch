package lister

import "io"

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
	panic("implement me")
}
