package core

type Plugin interface {
	Name() string
	Init()
}
