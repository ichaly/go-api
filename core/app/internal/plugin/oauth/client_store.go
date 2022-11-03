package oauth

import (
	"context"
	"encoding/json"
	"github.com/dosco/graphjin/serv"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/ichaly/go-api/core/app/pkg"
)

type ClientStore struct {
	Engine *serv.Service
}

func NewOauthClientStore() oauth2.ClientStore {
	return store.NewClientStore()
}

func (my *ClientStore) funGetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	gql := `query getClientByID($id: ID) {
	  clientsByID(id: $id) {
		id
		domain
	  }
	}`
	ql, err := my.Engine.GetGraphJin().GraphQL(ctx, gql, core.Variable{
		"id": id,
	}.Marshal(), nil)
	if err != nil {
		return nil, err
	}
	var c *models.Client
	if err = json.Unmarshal(ql.Data, c); err != nil {
		return nil, err
	}
	return c, nil
}
