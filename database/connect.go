package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
)

const (
	DefaultDatabaseFilename = "db.sqlite3"
)

var (
	Db *gorm.DB
)

func Connect() (shouldInitialize bool) {
	var err error
	dbFilename := os.Getenv("db_filename")
	if dbFilename == "" {
		dbFilename = DefaultDatabaseFilename
	}
	shouldInitialize = false
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) {
		shouldInitialize = true
	}
	Db, err = gorm.Open("sqlite3", dbFilename)
	if err != nil {
		log.Panic(err)
	}
	Db.SingularTable(true)
	return
}
