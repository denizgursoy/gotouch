//go:generate mockgen -source=./store.go -destination=mockStore.go -package=store

package store

import "sync"

const (
	ModuleName       = "ModuleName"
	ProjectName      = "ProjectName"
	WorkingDirectory = "WorkingDirectory"
	ProjectFullPath  = "ProjectFullPath"
	Dependencies     = "Dependencies"
	Inline           = "Inline"
)

type (
	Store interface {
		SetValue(key, value string)
		GetValue(key string) string
		AddValues(key map[string]any)
		GetValues() map[string]any
		AddCustomValues(key map[string]any)
		GetCustomValues() map[string]any
		AddDependency(dependency any)
	}

	storeImpl struct{}
)

var (
	store                map[string]any
	questionValues       map[string]any
	customQuestionValues map[string]any
	dependencies         []any
	keyValueStore        Store
	once                 = sync.Once{}
)

func init() {
}

func GetInstance() Store {
	once.Do(func() {
		keyValueStore = newStore()
		dependencies = make([]any, 0)
		questionValues = map[string]any{}
		customQuestionValues = map[string]any{}
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

func (s *storeImpl) AddValues(values map[string]any) {
	appendToMap(questionValues, values)
}

func (s *storeImpl) GetValues() map[string]any {
	questionValues[Dependencies] = dependencies
	return questionValues
}

func (s *storeImpl) AddDependency(dependency any) {
	dependencies = append(dependencies, dependency)
}

func (s *storeImpl) AddCustomValues(values map[string]any) {
	appendToMap(customQuestionValues, values)
}

func (s *storeImpl) GetCustomValues() map[string]any {
	return customQuestionValues
}

func appendToMap(target, source map[string]any) {
	if len(source) != 0 {
		for key, value := range source {
			target[key] = value
		}
	}
}
