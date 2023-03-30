package databases

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

var config DBConfig

var db *gorm.DB

func initConfig() {
	config.Name = os.Getenv("DB_NAME")
	config.Host = os.Getenv("DB_HOST")
	config.Port = os.Getenv("DB_PORT")
	config.User = os.Getenv("DB_USER")
	config.Password = os.Getenv("DB_PASSWORD")
}

func Connect() *gorm.DB {
	initConfig()

	var dsn string

	if config.Password == "" {
		dsn = fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Host, config.Port, config.Name)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Password, config.Host, config.Port, config.Name)
	}

	if db == nil {
		database, err := gorm.Open(mysql.New(
			mysql.Config{
				DSN: dsn,
			},
		), &gorm.Config{
			Logger:               logger.Default.LogMode(logger.Info),
			FullSaveAssociations: true,
		})
		if err != nil {
			panic(err)
		}
		db = database
	}

	return db
}
