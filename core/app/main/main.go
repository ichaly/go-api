package main

import (
	"context"
	"encoding/json"
	"github.com/dosco/graphjin/core"
	"github.com/dosco/graphjin/serv"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"io"
	"log"
	"net/http"
	"path"
	"path/filepath"
)

type GqlRequest struct {
	Query         string          `json:"query"`
	OperationName string          `json:"operationName"`
	Variables     json.RawMessage `json:"variables"`
	Extensions    json.RawMessage `json:"extensions"`
	Headers       http.Header     `json:"headers"`
}

func ReadJson(r io.Reader, val interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(val)
}

func main() {
	cp, err := filepath.Abs("./config")
	if err != nil {
		log.Fatal(err)
	}
	conf, err := serv.ReadInConfig(path.Join(cp, serv.GetConfigName()))
	if err != nil {
		log.Fatal(err)
	}
	gj, err := serv.NewGraphJinService(conf)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.POST("/graphql", func(c *gin.Context) {
		var params *GqlRequest
		if err := ReadJson(c.Request.Body, &params); err != nil {
			log.Print(err)
			return
		}
		ctx := context.Background()
		ctx = context.WithValue(ctx, core.UserIDKey, 1)
		res, err := gj.GraphQL(ctx, params.Query, params.Variables, nil)
		if err != nil {
			log.Print(err)
			return
		}
		c.JSON(http.StatusOK, res)
	})
	_ = r.Run()
}
