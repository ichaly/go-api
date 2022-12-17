package core

import (
	"github.com/ichaly/go-api/core/app/internal/base"
	"github.com/ichaly/go-api/core/app/internal/plugin/captcha"
	"github.com/ichaly/go-api/core/app/internal/plugin/explorer"
	"github.com/ichaly/go-api/core/app/internal/plugin/oauth"
	_ "github.com/ichaly/go-env/auto"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(
		base.NewConfig,
		base.NewStore,
		base.NewCache,
		base.NewDatabase,
		base.NewEngine,
		base.NewServer,
		base.NewCasbin,
	),
	fx.Provide(fx.Annotated{
		Group:  "plugin",
		Target: captcha.NewCaptchaService,
	}),
	fx.Provide(fx.Annotated{
		Group:  "plugin",
		Target: explorer.NewExplorerService,
	}),
	fx.Provide(
		oauth.NewOauthServer,
		oauth.NewOauthTokenStore,
		oauth.NewOauthClientStore,
		fx.Annotated{
			Group:  "plugin",
			Target: oauth.NewOauthService,
		},
		fx.Annotated{
			Group:  "middleware",
			Target: oauth.NewOauthVerify,
		},
	),
	fx.Invoke(base.Bootstrap),
)
