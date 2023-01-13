package main

import (
	"fmt"
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
    cats, _ := c.GetCategories()

    fmt.Printf("%v\n", cats)
}

func getEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
