//go:generate mockgen -source=./compressor.go -destination=mockCompressor.go -package=compressor

package compressor

import (
	"sync"
)

type (
	Compressor interface {
		UncompressFromUrl(url string) error
		CompressDirectory(source, target string) error
	}
)

var (
	compressor Compressor
	once       sync.Once
)

func GetInstance() Compressor {
	once.Do(func() {
		compressor = newZipCompressor()
	})
	return compressor
}
