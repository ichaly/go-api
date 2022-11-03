package oauth

import (
	"context"
	"github.com/dosco/graphjin/serv"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/ichaly/go-api/core/app/pkg"
	"github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ClientStore struct {
	Engine *serv.Service
}

func NewOauthClientStore(s *serv.Service) oauth2.ClientStore {
	return &ClientStore{Engine: s}
}

// GetByID http://127.0.0.1:8080/oauth/token?grant_type=client_credentials&client_id=0&client_secret=eV4YeKI484E1nVoioE91Os6eOQip0TFs&scope=read
func (my *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	gql := `query getClientByID($id: ID) {
	  clientsByID(id: $id) {
		id
		secret
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
	raw := gjson.GetBytes(ql.Data, "clientsByID").Raw
	if err = json.UnmarshalFromString(raw, &c); err != nil {
		return nil, err
	}
	return c, nil
}
