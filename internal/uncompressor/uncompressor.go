//go:generate mockgen -source=./uncompressor.go -destination=mock-uncompressor.go -package=uncompressor

package uncompressor

import "sync"

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
		extractor = newZipUncompressor()
	})
	return extractor
}
