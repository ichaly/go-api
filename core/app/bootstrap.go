package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func Bootstrap(lifecycle fx.Lifecycle, e *gin.Engine, c *Config) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				fmt.Printf("Now server is running on port %s\n", c.Port())
				fmt.Printf("Test with Get: curl -g 'http://%s/%s?query={hello}'\n", c.HostPort, c.Endpoint)
				_ = e.Run(fmt.Sprintf(":%v", c.HostPort))
			}()
			return nil
		},
	})
}
