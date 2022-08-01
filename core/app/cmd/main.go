package main

import (
	"github.com/ichaly/go-api/core/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		//禁用fx 默认logger
		fx.NopLogger,
		core.Modules,
	).Run()
}
