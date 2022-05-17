package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"pipelineService/env"
	"pipelineService/utils"
)

var db *gorm.DB

func init() {
	logger := utils.GetLogger()

	pgDB, err := GetClient(env.Env.DbHost, env.Env.DbUsername, env.Env.DbPassword, env.Env.DbName, env.Env.DbPort, "disable")

	if err != nil {
		logger.Error(err.Error())
	}

	db = pgDB
}

func GetConnection() *gorm.DB {
	return db
}

func GetClient(host string, user string, password string, dbName string, port string, sslmode string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbName, port, sslmode)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func CloseConnection(db *gorm.DB) {
	logger := utils.GetLogger()

	dbSQL, err := db.DB()

	if err != nil {
		logger.Error(err.Error())

		return
	}

	if err = dbSQL.Close(); err != nil {
		logger.Error(err.Error())

		return
	}
}
