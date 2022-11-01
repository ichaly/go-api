package base

import (
	"github.com/dosco/graphjin/serv"
	"gorm.io/gorm"
)

func NewEngine(c *Config, d *gorm.DB) (*serv.Service, error) {
	db, err := d.DB()
	if err != nil {
		return nil, err
	}
	return serv.NewGraphJinService(&c.Engine, serv.OptionSetDB(db))
}
