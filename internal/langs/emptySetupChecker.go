package langs

func NewEmptySetupChecker() LanguageChecker {
	return &emptySetupChecker{}
}

type emptySetupChecker struct {
}

func (e *emptySetupChecker) CheckSetup() error {
	return nil
}
