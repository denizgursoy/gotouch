package extractor

import "sync"

type Extractor interface {
	Extract(url, directoryName string)
}

var (
	extractor Extractor
	once      sync.Once
)

func GetInstance() Extractor {
	once.Do(func() {
		extractor = newGzipExtractor()
	})
	return extractor
}
