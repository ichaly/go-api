package core

import "github.com/go-chi/chi"

type Plugin interface {
	Name() string
	Protected() bool
	Init(r chi.Router)
}
