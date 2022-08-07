package model

type Requirement interface {
	AskForInput() Task
}
