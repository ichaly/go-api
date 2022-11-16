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
var html embed.FS

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
	if my.Config.Core.Production {
		return
	}

	if root, err := fs.Sub(html, "html"); err == nil {
		r.Route("/api", func(r chi.Router) {
			r.Handle("/*", http.FileServer(http.FS(root)))
		})
	}
}
