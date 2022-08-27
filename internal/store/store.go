//go:generate mockgen -source=./store.go -destination=mockStore.go -package=store

package store

import (
	"sync"
)

const (
	ModuleName       = "ModuleName"
	ProjectName      = "ProjectName"
	WorkingDirectory = "WorkingDirectory"
	ProjectFullPath  = "ProjectFullPath"
)

type (
	Store interface {
		SetValue(key, value string)
		GetValue(key string) string
	}

	storeImpl struct {
	}
)

var (
	store         = map[string]string{}
	keyValueStore Store
	once          = sync.Once{}
)

func init() {

}

func GetInstance() Store {
	once.Do(func() {
		keyValueStore = newStore()
	})
	return keyValueStore
}

func newStore() Store {
	return &storeImpl{}
}

func (s *storeImpl) SetValue(key, value string) {
	store[key] = value
}

func (s *storeImpl) GetValue(key string) string {
	return store[key]
}
