//go:generate mockgen -source=./compressor.go -destination=mockCompressor.go -package=compressor

package compressor

import (
	"github.com/denizgursoy/gotouch/internal/manager"
	"sync"
)

type (
	Uncompressor interface {
		UncompressFromUrl(url, directoryName string)
	}
)

var (
	extractor Uncompressor
	once      sync.Once
)

func GetInstance() Uncompressor {
	once.Do(func() {
		extractor = newZipUncompressor(manager.GetInstance())
	})
	return extractor
}
