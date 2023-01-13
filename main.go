package main

import (
	"log"
	"os"

	"kev/client"

	"github.com/joho/godotenv"
)

func main() {
	appId := getEnv("APP_ID")
	secret := getEnv("APP_SECRET")
	token := getEnv("APP_TOKEN")

	c := client.NewLinnworksClient(appId, secret, token)
	c.GetCategories()
}

func getEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
