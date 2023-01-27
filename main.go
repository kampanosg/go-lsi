package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/kampanosg/go-lsi/clients/db/sqlite"
	"github.com/kampanosg/go-lsi/clients/linnworks"
	"github.com/kampanosg/go-lsi/clients/square"
	"github.com/kampanosg/go-lsi/controllers"
	"github.com/kampanosg/go-lsi/middlewares"
	"github.com/kampanosg/go-lsi/sync"
	"github.com/kampanosg/go-lsi/types"

	"github.com/joho/godotenv"
)

func main() {

	dbPath := getEnv("DB")

	lwAppId := getEnv("LINNWORKS_APP_ID")
	lwAppSecret := getEnv("LINNWORKS_APP_SECRET")
	lwAppToken := getEnv("LINNWORKS_APP_TOKEN")

	sqAccessToken := getEnv("SQUARE_ACCESS_TOKEN")
	sqHost := getEnv("SQUARE_HOST")
	sqApiVersion := getEnv("SQUARE_API_VERSION")
	sqLocationId := getEnv("SQUARE_LOCATION_ID")

	sqliteDb := sqlite.NewSqliteDB(dbPath)
	lwClient := linnworks.NewLinnworksClient(lwAppId, lwAppSecret, lwAppToken)
	sqClient := square.NewSquareClient(sqAccessToken, sqHost, sqApiVersion, sqLocationId)
	syncTool := sync.NewSyncTool(lwClient, sqClient, sqliteDb)
	log.Println(syncTool)

	products := []types.OrderProduct{
		{
			Id:            1,
			SquareOrderId: "9urRtTF6Qwzt01tEiCHs92O3Ri4F",
			SquareVarId:   "XUALJKJOOFNXLUU47H7PWIDL",
			Quantity:      "1",
			ItemNumber:    "5060464363757",
			SKU:           "JC10-BK",
			Title:         "Very Good Coffee",
			PricePerUnit:  7.1,
		},
		{
			Id:            1,
			SquareOrderId: "aurRtTF6Qwzt01tEiCHs92O3Ri4F",
			SquareVarId:   "YUALJKJOOFNXLUU47H7PWIDL",
			Quantity:      "1",
			ItemNumber:    "6060464363757",
			SKU:           "JC11-CK",
			Title:         "Very Medium Coffee",
			PricePerUnit:  7.1,
		},
	}

	orders := []types.Order{
		{
			Id:         1,
			SquareId:   "9urRtTF6Qwzt01tEiCHs92O3Ri4F",
			Products:   products,
			LocationId: "Default",
			State:      "Completed",
			Version:    1,
			TotalMoney: 14.2,
			CreatedAt:  time.Now(),
		},
	}
	lwClient.CreateOrders(orders)
}

func main2() {

	port := getEnv("HTTP_PORT")
	log.Printf("Starting server at port :%s\n", port)

	dbPath := getEnv("DB")

	signingKey := []byte(getEnv("SIGNING_KEY"))

	lwAppId := getEnv("LINNWORKS_APP_ID")
	lwAppSecret := getEnv("LINNWORKS_APP_SECRET")
	lwAppToken := getEnv("LINNWORKS_APP_TOKEN")

	sqAccessToken := getEnv("SQUARE_ACCESS_TOKEN")
	sqHost := getEnv("SQUARE_HOST")
	sqApiVersion := getEnv("SQUARE_API_VERSION")
	sqLocationId := getEnv("SQUARE_LOCATION_ID")

	sqliteDb := sqlite.NewSqliteDB(dbPath)
	lwClient := linnworks.NewLinnworksClient(lwAppId, lwAppSecret, lwAppToken)
	sqClient := square.NewSquareClient(sqAccessToken, sqHost, sqApiVersion, sqLocationId)
	syncTool := sync.NewSyncTool(lwClient, sqClient, sqliteDb)

	authMiddleware := middlewares.NewAuthMiddleware(signingKey)
	authController := controllers.NewAuthController(sqliteDb, signingKey)
	inventoryController := controllers.NewInventoryController(sqliteDb)
	ordersController := controllers.NewOrdersController(sqliteDb)
	pingController := controllers.NewPingController()
	syncController := controllers.NewSyncController(syncTool)

	router := mux.NewRouter()

	router.Handle("/api/v1/ping", authMiddleware.ProtectedEndpoint(http.HandlerFunc(pingController.HandlePingRequest)))
	router.Handle("/api/v1/inventory", authMiddleware.ProtectedEndpoint(http.HandlerFunc(inventoryController.HandleInventoryRequest)))
	router.Handle("/api/v1/orders", authMiddleware.ProtectedEndpoint(http.HandlerFunc(ordersController.HandleOrdersRequest)))
	router.Handle("/api/v1/sync", authMiddleware.ProtectedEndpoint(http.HandlerFunc(syncController.HandleSyncRequest)))
	router.Handle("/api/v1/auth", http.HandlerFunc(authController.HandleAuthRequest))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatalf("Unable to start server. error=%v\n", err.Error())
	}
}

func getEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
