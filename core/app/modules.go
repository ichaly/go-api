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
		internal.NewEngine,
		internal.NewServer,
	),
	fx.Provide(
		fx.Annotated{
			Target: plugin.NewCaptchaService,
			Group:  "plugin",
		},
	),
	fx.Invoke(internal.Bootstrap),
)
