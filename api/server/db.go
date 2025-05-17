package server

import (
	"log"
	"sync"
	"weatherapi/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

func GetDb() *gorm.DB {
	dbOnce.Do(func() {
		conf := configs.Get()
		var err error
		dbInstance, err = gorm.Open(postgres.Open(conf.DbConnectionString), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}
	})
	return dbInstance
}

func init() {
	GetDb()
}
