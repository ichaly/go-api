package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ichaly/go-env"
	"go.uber.org/fx"
)

func Bootstrap(lifecycle fx.Lifecycle, engine *gin.Engine) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				port, _ := env.String(":${PORT:=8080}")
				fmt.Printf("Now server is running on port %s\n", port)
				fmt.Printf("Connect to http://localhost%s/ for GraphQL playground\n", port)
				fmt.Printf("Test with Get: curl -g 'http://localhost%s/graphql?query={hello}'\n", port)
				_ = engine.Run(port)
			}()
			return nil
		},
	})
}
