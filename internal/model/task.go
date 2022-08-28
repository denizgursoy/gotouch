package model

type Requirement interface {
	AskForInput() ([]Task, []Requirement, error)
}

type Task interface {
	Complete() error
}
