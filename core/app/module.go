package core

import (
	"github.com/ichaly/go-api/core/app/internal"
	_ "github.com/ichaly/go-env/auto"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(
		internal.NewConfig,
		internal.NewEngine,
		internal.NewServer,
	),
	fx.Invoke(internal.Bootstrap),
)
