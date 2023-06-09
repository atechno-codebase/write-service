package service

import (
	"write/config"
	"write/models"
)

func Init() {
	dbName := config.Configuration.Database
	mongoUrl := config.Configuration.MongoURL
	models.Init(mongoUrl, dbName)
	return

}
