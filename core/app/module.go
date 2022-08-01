package app

import (
	_ "github.com/ichaly/go-env/auto"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(NewConfig),
	fx.Provide(NewEngine),
	fx.Provide(NewServer),
	fx.Invoke(Bootstrap),
)
