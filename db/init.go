package db

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mrjoshlab/pepe.bot/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Connection *gorm.DB
)

func Configure() (err error) {
	dbLogger := logger.New(log.New(ioutil.Discard, "\r\n", log.LstdFlags), logger.Config{})
	Connection, err = gorm.Open(sqlite.Open(config.Map.DB.Path), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return fmt.Errorf("Error opening sqlite3 database CLI: %s", err)
	}
	return nil
}
