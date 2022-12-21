package resolver

import (
	"context"
	"encoding/json"
	gql "github.com/dosco/graphjin/core"
	"github.com/go-chi/chi"
	"github.com/ichaly/go-api/core/app/internal/base"
	core "github.com/ichaly/go-api/core/app/pkg"
)

type Test struct {
	Config *base.Config
}

func NewTest(c *base.Config) core.Plugin {
	t := &Test{Config: c}
	_ = c.Core.SetResolver("test", func(v gql.ResolverProps) (gql.Resolver, error) {
		return t, nil
	})
	return t
}

func (my *Test) Name() string {
	return "TestResolver"
}

func (my *Test) Protected() bool {
	return true
}

func (my *Test) Init(r chi.Router) {

}

func (my *Test) Resolve(ctx context.Context, rr gql.ResolverReq) ([]byte, error) {
	m := map[string]interface{}{
		"desc":  "Test 1 for payment_id_1002" + rr.ID,
		"count": "1",
	}
	return json.Marshal(m)
}
