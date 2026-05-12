//go:generate mockgen -source=./compressor.go -destination=mockCompressor.go -package=compressor --typed

package compressor

import (
	"context"
	"sync"
)

type (
	Compressor interface {
		UncompressFromUrl(ctx context.Context, url string, directoryToCreateProject string) error
		UncompressFromLocalPath(ctx context.Context, path string, directoryToCreateProject string) error
		CopyDirectory(path string, directoryToCreateProject string) error
		CheckIfFileExtensionIsSupported(path string) error
		CompressDirectory(source, target string) error
	}

	ZipStrategy interface {
		UnCompressDirectory(source, target string) error
		CompressDirectory(source, target string) error
		GetExtension() string
	}
)

var (
	compressorInstance Compressor
	once               sync.Once
)

func GetInstance() Compressor {
	once.Do(func() {
		compressorInstance = newCompressor()
	})
	return compressorInstance
}
