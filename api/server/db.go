package server

import (
	"log"
	"strings"
	"sync"
	"weatherapi/configs"

	"github.com/glebarez/sqlite"
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

		var dialector gorm.Dialector
		if strings.HasPrefix(conf.DbConnectionString, "postgresql://") {
			dialector = postgres.Open(conf.DbConnectionString)
		} else if strings.HasPrefix(conf.DbConnectionString, "sqlite://") {
			sqliteConnection, _ := strings.CutPrefix(conf.DbConnectionString, "sqlite://")
			dialector = sqlite.Open(sqliteConnection)
		}

		if dialector == nil {
			log.Fatalf("unsupported database connection string: %s", conf.DbConnectionString)
		}

		dbInstance, err = gorm.Open(dialector, &gorm.Config{})

		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}
	})
	return dbInstance
}

func init() {
	GetDb()
}
