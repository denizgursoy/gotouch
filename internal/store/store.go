//go:generate mockgen -source=./store.go -destination=mockStore.go -package=store

package store

import "sync"

const (
	ModuleName       = "ModuleName"
	ProjectName      = "ProjectName"
	WorkingDirectory = "WorkingDirectory"
	ProjectFullPath  = "ProjectFullPath"
	Dependencies     = "Dependencies"
)

type (
	Store interface {
		SetValue(key, value string)
		GetValue(key string) string
		AddValues(key map[string]any)
		GetValues() map[string]any
		AddDependency(dependency any)
	}

	storeImpl struct{}
)

var (
	store          map[string]any
	questionValues map[string]any
	dependencies   []any
	keyValueStore  Store
	once           = sync.Once{}
)

func init() {
}

func GetInstance() Store {
	once.Do(func() {
		keyValueStore = newStore()
		dependencies = make([]any, 0)
		questionValues = map[string]any{}
		store = map[string]any{}
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
	return store[key].(string)
}

func (s *storeImpl) AddValues(key map[string]any) {
	if key != nil {
		for key, value := range key {
			questionValues[key] = value
		}
	}
}

func (s *storeImpl) GetValues() map[string]any {
	questionValues[Dependencies] = dependencies
	return questionValues
}

func (s *storeImpl) AddDependency(dependency any) {
	dependencies = append(dependencies, dependency)
}
