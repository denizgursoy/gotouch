package model

import "context"

type Requirement interface {
	AskForInput() ([]Task, []Requirement, error)
}

type Task interface {
	Complete(ctx context.Context) error
}
