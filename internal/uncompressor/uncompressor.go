package uncompressor

import "sync"

type Uncompressor interface {
	Uncompress(url, directoryName string)
}

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
