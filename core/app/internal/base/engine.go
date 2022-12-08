package base

import (
	"fmt"
	"github.com/dosco/graphjin/core"
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/ichaly/go-api/core/app/pkg/render"
	"github.com/ichaly/go-api/core/app/pkg/util"
	"gorm.io/gorm"
	"io"
	"net/http"
)

const (
	maxReadBytes = 100000 // 100Kb
)

type Engine struct {
	Graph *core.GraphJin
	Cache *cache.Cache[string]
}

func NewEngine(c *Config, d *gorm.DB, s *cache.Cache[string]) (*Engine, error) {
	db, err := d.DB()
	if err != nil {
		return nil, err
	}
	jin, err := core.NewGraphJin(&c.Core, db)
	if err != nil {
		return nil, err
	}
	return &Engine{Graph: jin, Cache: s}, nil
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

		var key string
		if str, err := json.MarshalToString(req); err == nil {
			key = fmt.Sprintf("gql:%s", util.MD5(str))
			if len(key) > 0 {
				if val, err := my.Cache.Get(r.Context(), key); err == nil {
					res := &core.Result{}
					if err := json.UnmarshalFromString(val, res); err == nil {
						_ = render.JSON(w, res)
						return
					}
				}
			}
		}

		if res, err := my.Graph.GraphQL(r.Context(), req.Query, req.Vars, nil); err == nil {
			_ = render.JSON(w, res)
			if len(key) > 0 && len(res.Errors) == 0 {
				if core.OpQuery == res.Operation() {
					if val, err := json.MarshalToString(res); err == nil {
						_ = my.Cache.Set(r.Context(), key, val, store.WithTags(res.Tables()))
					}
				} else {
					_ = my.Cache.Invalidate(r.Context(), store.WithInvalidateTags(res.Tables()))
				}
			}
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
