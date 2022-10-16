package langs

import (
	"fmt"
	"os/exec"
)

func NewGolangSetupChecker() LanguageChecker {
	return &golangSetupChecker{}
}

type golangSetupChecker struct {
}

func (g *golangSetupChecker) CheckSetup() error {
	_, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("could not find %s in PATH. make sure that %s installed", "go", "go")
	}
	return nil

}
