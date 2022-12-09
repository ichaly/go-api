package base

import (
	"context"
	eson "encoding/json"
	"fmt"
	gql "github.com/dosco/graphjin/core"
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/ichaly/go-api/core/app/pkg"
	"github.com/ichaly/go-api/core/app/pkg/render"
	"github.com/ichaly/go-api/core/app/pkg/util"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
)

const (
	maxReadBytes       = 100000 // 100Kb
	bigId              = "big_id"
	graphqlEndpoint    = "/api/v1/graphql"
	introspectionQuery = "IntrospectionQuery"
)

type Engine struct {
	Graph *gql.GraphJin
	Cache *cache.Cache[string]
}

func NewEngine(c *Config, d *gorm.DB, s *cache.Cache[string]) (*Engine, error) {
	db, err := d.DB()
	if err != nil {
		return nil, err
	}
	jin, err := gql.NewGraphJin(&c.Core, db)
	if err != nil {
		return nil, err
	}
	return &Engine{Graph: jin, Cache: s}, nil
}

func (my *Engine) Attach(r chi.Router) {
	r.HandleFunc(graphqlEndpoint, my.graphqlHandler())
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
		// 从缓存中获取数据
		var key string
		if req.OpName != introspectionQuery {
			if str, err := json.MarshalToString(req); err == nil {
				key = fmt.Sprintf("gql:%v", util.MD5(str))
				if len(key) > 0 {
					if val, err := my.Cache.Get(r.Context(), key); err == nil {
						reps := &gqlResp{Code: 200, Result: &gql.Result{}}
						if err := json.UnmarshalFromString(val, reps); err == nil {
							_ = render.JSON(w, reps)
							return
						}
					}
				}
			}
		}
		// 存储当前登陆者id
		ctx := context.WithValue(r.Context(), gql.UserIDKey, 101)
		// 配置雪花id生成
		rc := &gql.ReqConfig{Vars: map[string]interface{}{bigId: getBigId}}
		// 执行GraphQL结果
		if res, err := my.Graph.GraphQL(ctx, req.Query, req.Vars, rc); err == nil {
			_ = render.JSON(w, gqlResp{Code: 200, Result: res})
			// 存储到缓存中
			if len(key) > 0 && len(res.Errors) == 0 {
				if gql.OpQuery == res.Operation() {
					if val, err := json.MarshalToString(res); err == nil {
						_ = my.Cache.Set(r.Context(), key, val, store.WithTags(res.Tables()))
					}
				} else {
					_ = my.Cache.Invalidate(r.Context(), store.WithInvalidateTags(res.Tables()))
				}
			}
		} else {
			_ = render.JSON(w, core.ERROR.WithError(err))
		}
	}
}

func getBigId() string {
	if id, err := GenerateID(); err == nil {
		return strconv.FormatUint(id, 10)
	}
	return ""
}

type gqlReq struct {
	OpName string          `json:"operationName"`
	Query  string          `json:"query"`
	Vars   eson.RawMessage `json:"variables,omitempty"`
	Ext    extensions      `json:"extensions,omitempty"`
}

type extensions struct {
	Persisted apqExt `json:"persistedQuery,omitempty"`
}

type apqExt struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

type gqlResp struct {
	Code int `json:"code"`
	*gql.Result
}
