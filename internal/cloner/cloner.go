//go:generate mockgen -source=$GOFILE -destination=mockCloner.go -package=cloner

package cloner

import (
	"context"
	"sync"
)

type Cloner interface {
	CloneFromUrl(ctx context.Context, url, branchName string) error
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
