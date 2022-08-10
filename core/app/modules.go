package core

import (
	"github.com/ichaly/go-api/core/app/internal"
	"github.com/ichaly/go-api/core/app/internal/plugin"
	_ "github.com/ichaly/go-env/auto"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(
		internal.NewConfig,
		internal.NewDatabase,
		internal.NewEngine,
		internal.NewServer,
	),
	fx.Provide(
		fx.Annotated{
			Group:  "plugin",
			Target: plugin.NewCaptchaService,
		},
	),
	fx.Invoke(internal.Bootstrap),
)
