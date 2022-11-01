package core

import (
	"github.com/ichaly/go-api/core/app/internal/base"
	"github.com/ichaly/go-api/core/app/internal/plugin/captcha"
	"github.com/ichaly/go-api/core/app/internal/plugin/oauth"
	_ "github.com/ichaly/go-env/auto"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(
		base.NewConfig,
		base.NewDatabase,
		base.NewEngine,
		base.NewServer,
		base.NewCache,
	),
	fx.Provide(fx.Annotated{
		Group:  "plugin",
		Target: captcha.NewCaptchaService,
	}),
	fx.Provide(
		oauth.NewOauthServer,
		oauth.NewOauthTokenStore,
		oauth.NewOauthClientStore,
		fx.Annotated{
			Group:  "plugin",
			Target: oauth.NewOauthService,
		},
		//fx.Annotated{
		//	Group:  "middleware",
		//	Target: oauth.NewOauthTokenVerify,
		//},
	),
	fx.Invoke(base.Bootstrap),
)
