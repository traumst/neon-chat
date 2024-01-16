// package models

type HTML interface {
	GetHTML() (string, error)
}

type Loggable interface {
	Log() string
}
