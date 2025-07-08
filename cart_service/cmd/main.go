package main

import (
	"cart_service/internal/app"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
}
func main() {
	app.Run()
}
