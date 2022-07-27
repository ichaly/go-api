package app

import (
	"context"
	"encoding/json"
	"github.com/dosco/graphjin/core"
	"github.com/gin-gonic/gin"
	"github.com/ichaly/go-api/core/app/util"
	"io"
	"net/http"
)

type request struct {
	Query         string          `json:"query"`
	OperationName string          `json:"operationName"`
	Variables     json.RawMessage `json:"variables"`
	Extensions    json.RawMessage `json:"extensions"`
	Headers       http.Header     `json:"headers"`
}

func readJson(r io.Reader, val interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(val)
}

func handler(gj *core.GraphJin) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *request
		if err := readJson(c.Request.Body, &req); err != nil {
			return
		}
		ctx := context.WithValue(c, core.UserIDKey, 1)
		res, err := gj.GraphQL(ctx, req.Query, req.Variables, nil)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func NewServer(g *core.GraphJin, c *Config) *gin.Engine {
	r := gin.New()
	r.POST(util.String(c.Endpoint, "/graphql"), handler(g))
	return r
}
