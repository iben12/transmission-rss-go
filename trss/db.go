package transmissionrss

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var once sync.Once

type DB struct{}

var dbConnection *gorm.DB

func (d *DB) getConnection() (db *gorm.DB) {
	if dbConnection == nil {
		once.Do(
			func() {
				Logger.Info().
					Str("action", "DB connect").
					Msg("Creating new DB connection")
				dbConnection = d.connect()
			})
	} else {
		Logger.Info().
			Str("action", "DB connect").
			Msg("DB is already connected")
	}

	return dbConnection
}

func (d *DB) connect() (db *gorm.DB) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				IgnoreRecordNotFoundError: true,
			},
		),
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Episode{})

	return db
}
