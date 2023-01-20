package main

import (
	// "fmt"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/kampanosg/go-lsi/controllers"
	middleware "github.com/kampanosg/go-lsi/middlewares"

	// "strings"

	"github.com/kampanosg/go-lsi/clients/db/sqlite"
	// "github.com/kampanosg/go-lsi/clients/linnworks"
	// "github.com/kampanosg/go-lsi/clients/square"
	// "github.com/kampanosg/go-lsi/sync"

	"github.com/joho/godotenv"
)

func main() {

	port := 8080 // TODO: Bring from config
	log.Printf("Starting server at port :%d\n", port)

	dbPath := getEnv("DB")
	signingKey := []byte(getEnv("SIGNING_KEY"))
	sqliteDb := sqlite.NewSqliteDB(dbPath)

	authMiddleware := middleware.NewAuthMiddleware(signingKey)
	authController := controllers.NewAuthController(sqliteDb, signingKey)
	pingController := controllers.NewPingController()

	router := mux.NewRouter()

	router.Handle("/ping", authMiddleware.ProtectedEndpoint(http.HandlerFunc(pingController.HandlePingRequest)))
	router.Handle("/auth", http.HandlerFunc(authController.HandleAuthRequest))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		log.Fatalf("Unable to start server. error=%v\n", err.Error())
	}
	// lwAppId := getEnv("LINNWORKS_APP_ID")
	// lwAppSecret := getEnv("LINNWORKS_APP_SECRET")
	// lwAppToken := getEnv("LINNWORKS_APP_TOKEN")

	// sqAccessToken := getEnv("SQUARE_ACCESS_TOKEN")
	// sqHost := getEnv("SQUARE_HOST")

	// lwClient := linnworks.NewLinnworksClient(lwAppId, lwAppSecret, lwAppToken)
	// sqClient := square.NewSquareClient(sqAccessToken, sqHost)
	// s := sync.NewSyncTool(lwClient, sqClient, sqliteDb)
	// s.SyncCategories()
	// s.SyncProducts()
}

func handleAuthRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Ok",
	})
}

func getEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
