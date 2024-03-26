package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

func NewGorm(viper *viper.Viper) *gorm.DB {
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	dbName := viper.GetString("database.name")
	idl := viper.GetInt("database.pooling.idle")
	maxConn := viper.GetInt("database.pooling.max")
	lifetime := viper.GetInt("database.pooling.lifetime")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local", user, password, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to the database %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("Error connecting to the database %v", err)
	}

	connection.SetMaxIdleConns(idl)
	connection.SetConnMaxLifetime(time.Second * time.Duration(lifetime))
	connection.SetMaxOpenConns(maxConn)
	return db
}
