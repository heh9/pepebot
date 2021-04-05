package db

import (
	"fmt"

	"github.com/mrjoshlab/pepe.bot/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Connection *gorm.DB
)

func Configure() (err error) {
	Connection, err = gorm.Open(sqlite.Open(config.Map.DB.Path), nil)
	if err != nil {
		return fmt.Errorf("Error opening sqlite3 database CLI: %s", err)
	}
	return nil
}
