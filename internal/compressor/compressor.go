//go:generate mockgen -source=./compressor.go -destination=mockCompressor.go -package=compressor

package compressor

import (
	"github.com/denizgursoy/gotouch/internal/manager"
	"sync"
)

type (
	Compressor interface {
		UncompressFromUrl(url, directoryName string) error
	}
)

var (
	compressor Compressor
	once       sync.Once
)

func GetInstance() Compressor {
	once.Do(func() {
		compressor = newZipCompressor(manager.GetInstance())
	})
	return compressor
}
