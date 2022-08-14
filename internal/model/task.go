package model

type Task interface {
	Complete(interface{}) (interface{}, error)
}
