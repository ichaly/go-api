package oauth

import (
	"context"
	"github.com/dosco/graphjin/serv"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ClientStore struct {
	Engine   *serv.Service
	Database *gorm.DB
}

func NewOauthClientStore(d *gorm.DB, s *serv.Service) oauth2.ClientStore {
	return &ClientStore{Engine: s, Database: d}
}

// GetByID http://127.0.0.1:8080/oauth/token?grant_type=client_credentials&client_id=0&client_secret=eV4YeKI484E1nVoioE91Os6eOQip0TFs&scope=read
func (my *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	var c *models.Client
	m := map[string]interface{}{}
	my.Database.Table("clients").Where("id = ?", id).Take(&m)
	err := mapstructure.WeakDecode(m, &c)
	return c, err
}
