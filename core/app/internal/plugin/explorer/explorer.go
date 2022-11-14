package explorer

import (
	"embed"
	"github.com/go-chi/chi"
	"github.com/ichaly/go-api/core/app/internal/base"
	"github.com/ichaly/go-api/core/app/pkg"
	"io/fs"
	"net/http"
)

//go:embed html
var staticBuild embed.FS

type ExplorerService struct {
	Config *base.Config
}

func NewExplorerService(c *base.Config) core.Plugin {
	return &ExplorerService{Config: c}
}

func (my *ExplorerService) Name() string {
	return "ExplorerService"
}

func (my *ExplorerService) Protected() bool {
	return true
}

func (my *ExplorerService) Init(r chi.Router) {
	if my.Config.Engine.Core.Production {
		return
	}

	if webRoot, err := fs.Sub(staticBuild, "html"); err == nil {
		r.Route("/api", func(r chi.Router) {
			r.Handle("/*", http.StripPrefix("/api/", http.FileServer(http.FS(webRoot))))
		})
	}
}
