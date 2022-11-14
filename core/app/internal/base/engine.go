package base

import (
	"github.com/dosco/graphjin/core"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/ichaly/go-api/core/app/pkg/render"
	"gorm.io/gorm"
	"io"
	"net/http"
)

const (
	maxReadBytes = 100000 // 100Kb
)

type Engine struct {
	jin *core.GraphJin
}

func NewEngine(c *Config, d *gorm.DB) (*Engine, error) {
	db, err := d.DB()
	if err != nil {
		return nil, err
	}
	jin, err := core.NewGraphJin(&c.Core, db)
	if err != nil {
		return nil, err
	}
	return &Engine{jin: jin}, nil

}

func (my *Engine) Attach(r chi.Router) {
	r.HandleFunc("/api/v1/graphql", my.graphqlHandler())
}

func (my *Engine) graphqlHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if websocket.IsWebSocketUpgrade(r) {
			//TODO: s.apiV1Ws(w, r, ah)
			return
		}
		var req gqlReq

		switch r.Method {
		case http.MethodPost:
			if b, err := io.ReadAll(io.LimitReader(r.Body, maxReadBytes)); err == nil {
				defer r.Body.Close()
				err = json.Unmarshal(b, &req)
			}
		case http.MethodGet:
			q := r.URL.Query()
			req.Query = q.Get("query")
			req.OpName = q.Get("operationName")
			req.Vars = []byte(q.Get("variables"))

			if ext := q.Get("extensions"); ext != "" {
				_ = json.UnmarshalFromString(ext, &req.Ext)
			}
		}

		if res, err := my.jin.GraphQL(r.Context(), req.Query, req.Vars, nil); err == nil {
			_ = render.JSON(w, res)
		}
	}
}

type gqlReq struct {
	OpName string     `json:"operationName"`
	Query  string     `json:"query"`
	Vars   []byte     `json:"variables"`
	Ext    extensions `json:"extensions"`
}

type extensions struct {
	Persisted apqExt `json:"persistedQuery"`
}

type apqExt struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}
