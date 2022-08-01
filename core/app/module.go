package app

import (
	_ "github.com/ichaly/go-env/auto"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(
		NewConfig,
		NewEngine,
		NewServer,
	),
	fx.Invoke(Bootstrap),
)
