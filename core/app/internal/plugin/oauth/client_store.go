package oauth

import (
	"context"
	"github.com/dosco/graphjin/serv"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"gorm.io/gorm"
)

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
	my.Database.Table("clients").Where("id = ?", id).Take(&c)
	return c, nil
}
