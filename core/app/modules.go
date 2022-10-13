package core

import (
	"github.com/ichaly/go-api/core/app/internal/base"
	"github.com/ichaly/go-api/core/app/internal/plugin"
	_ "github.com/ichaly/go-env/auto"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(
		base.NewConfig,
		base.NewEngine,
		base.NewServer,
		base.NewCache,
	),
	fx.Provide(fx.Annotated{
		Group:  "plugin",
		Target: plugin.NewCaptchaService,
	}),
	fx.Provide(fx.Annotated{
		Group:  "plugin",
		Target: plugin.NewOauthService,
	}),
	fx.Invoke(base.Bootstrap),
)
