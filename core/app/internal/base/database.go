package base

import (
	"fmt"
	"github.com/dosco/graphjin/serv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(c *Config) (*gorm.DB, error) {
	return gorm.Open(buildDialect(c.Engine.DB), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
}

func buildDialect(ds serv.Database) gorm.Dialector {
	args := []interface{}{ds.User, ds.Password, ds.Host, ds.Port, ds.DBName}
	if ds.Type == "mysql" {
		return mysql.Open(fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", args...,
		))
	} else {
		return postgres.Open(fmt.Sprintf(
			"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable TimeZone=Asia/Shanghai", args...,
		))
	}
}
