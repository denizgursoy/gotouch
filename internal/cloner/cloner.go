//go:generate mockgen -source=./cloner.go -destination=mockCloner.go -package=cloner

package cloner

import "sync"

type Cloner interface {
	CloneFromUrl(url, branchName string) error
}

var (
	compressorInstance Cloner
	once               sync.Once
)

func GetInstance() Cloner {
	once.Do(func() {
		compressorInstance = newCloner()
	})
	return compressorInstance
}
