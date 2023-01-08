package langs

import (
	"github.com/denizgursoy/gotouch/internal/logger"
)

type emptySetupChecker struct {
	Logger logger.Logger
}

func NewEmptySetupChecker() Checker {
	return &emptySetupChecker{}
}

func (e *emptySetupChecker) Setup() error {
	return nil
}

func (e *emptySetupChecker) CheckDependency(dependency any) error {
	return nil
}

func (e *emptySetupChecker) CleanUp() error {
	return nil
}

func (e *emptySetupChecker) GetDependency(dependency any) error {
	return nil
}

func (e *emptySetupChecker) CheckSetup() error {
	return nil
}
