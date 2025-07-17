package main

import (
	"product_service/internal/app"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		Logger.Errorf("Error loading .env file")
	}
}

func main() {
	app.Run()
}
