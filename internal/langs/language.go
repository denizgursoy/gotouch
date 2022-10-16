package langs

import (
	"github.com/denizgursoy/gotouch/internal/model"
	"strings"
)

type LanguageChecker interface {
	CheckSetup() error
}

func NewLanguageChecker(projectStructureData *model.ProjectStructureData) LanguageChecker {
	language := projectStructureData.Language
	if len(strings.TrimSpace(language)) == 0 ||
		strings.ToLower(language) == "golang" ||
		strings.ToLower(language) == "go" {
		return NewGolangSetupChecker()
	} else {
		return NewEmptySetupChecker()
	}
}
