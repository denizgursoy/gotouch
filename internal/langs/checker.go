//go:generate mockgen -source=./checker.go -destination=mockChecker.go -package=langs

package langs

import (
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/store"
	"strings"
	"sync"
)

type Checker interface {
	Init(language string, Logger logger.Logger, str store.Store)
	GetLangChecker() LanguageChecker
}

type LanguageChecker interface {
	CheckSetup() error
	CheckDependency(dependency interface{}) error
	GetDependency(dependency interface{}) error
	CompletePreTask() error
}

type x struct {
	LanguageChecker
}

var (
	main Checker
	once sync.Once
)

func GetInstance() Checker {
	once.Do(func() {
		main = &x{
			NewEmptySetupChecker(),
		}
	})
	return main
}

func (x *x) GetLangChecker() LanguageChecker {
	return x.LanguageChecker
}

func (x *x) Init(language string, Logger logger.Logger, str store.Store) {
	if len(strings.TrimSpace(language)) == 0 ||
		strings.ToLower(language) == "golang" ||
		strings.ToLower(language) == "go" {
		x.LanguageChecker = NewGolangSetupChecker(Logger, str)
	} else {
		x.LanguageChecker = NewEmptySetupChecker()
	}
}
