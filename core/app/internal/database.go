package internal

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

func NewDatabase(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(buildDialect(cfg.Database), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return nil, err
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDb.SetMaxIdleConns(10)                   //最大空闲连接数
	sqlDb.SetMaxOpenConns(30)                   //最大连接数
	sqlDb.SetConnMaxLifetime(time.Second * 300) //设置连接空闲超时
	return db, nil
}

func buildDialect(ds Database) gorm.Dialector {
	dsn := strings.Trim(ds.Url, "")
	args := []interface{}{ds.Username, ds.Password, ds.Host, ds.Port, ds.Name}
	if ds.Type == "mysql" {
		dsn = map[bool]string{true: dsn, false: fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", args...,
		)}[len(dsn) > 0]
		return mysql.Open(dsn)
	} else if ds.Type == "pgsql" {
		dsn = map[bool]string{true: dsn, false: fmt.Sprintf(
			"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable TimeZone=Asia/Shanghai", args...,
		)}[len(dsn) > 0]
		return postgres.Open(dsn)
	}
	return nil
}
