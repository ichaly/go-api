package core

import "go.uber.org/fx"

type Middleware interface {
	Name() string
	Init()
}

type Group struct {
	fx.In
	middlewares []*Middleware `group:"middleware"`
}
