package main

import (
	"fmt"
	"net/http"
	"os"

	// "time"

	"github.com/gorilla/mux"
	"github.com/kampanosg/go-lsi/clients/db/sqlite"
	"github.com/kampanosg/go-lsi/clients/linnworks"
	"github.com/kampanosg/go-lsi/clients/square"
	"github.com/kampanosg/go-lsi/controllers"
	"github.com/kampanosg/go-lsi/middlewares"
	"github.com/kampanosg/go-lsi/sync"

	// "github.com/kampanosg/go-lsi/types"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// func main() {

// 	dbPath := getEnv("DB")

// 	lwAppId := getEnv("LINNWORKS_APP_ID")
// 	lwAppSecret := getEnv("LINNWORKS_APP_SECRET")
// 	lwAppToken := getEnv("LINNWORKS_APP_TOKEN")

// 	sqAccessToken := getEnv("SQUARE_ACCESS_TOKEN")
// 	sqHost := getEnv("SQUARE_HOST")
// 	sqApiVersion := getEnv("SQUARE_API_VERSION")
// 	sqLocationId := getEnv("SQUARE_LOCATION_ID")

// 	sqliteDb := sqlite.NewSqliteDB(dbPath)
// 	lwClient := linnworks.NewLinnworksClient(lwAppId, lwAppSecret, lwAppToken)
// 	sqClient := square.NewSquareClient(sqAccessToken, sqHost, sqApiVersion, sqLocationId)
// 	syncTool := sync.NewSyncTool(lwClient, sqClient, sqliteDb)
// 	log.Println(syncTool)

// 	products := []types.OrderProduct{
// 		{
// 			Id:            1,
// 			SquareOrderId: "9urRtTF6Qwzt01tEiCHs92O3Rj4C",
// 			SquareVarId:   "XUALJKJOOFNXLUU47H7PWIDL",
// 			Quantity:      "1",
// 			ItemNumber:    "5060464363757",
// 			SKU:           "JC10-BK",
// 			Title:         "1/4\" Mono Output Jack Socket - Black",
// 			PricePerUnit:   2.99,
// 		},
// 		{
// 			Id:            1,
// 			SquareOrderId: "aurRtTF6Qwzt01tEiCHs92O3Rj4C",
// 			SquareVarId:   "YUALJKJOOFNYLUU47H7PWIDL",
// 			Quantity:      "1",
// 			ItemNumber:    "5060464363764",
// 			SKU:           "JC10-CR",
// 			Title:         "1/4\" Mono Output Jack Socket - Chrome",
// 			PricePerUnit:  2.99,
// 		},
// 	}

// 	orders := []types.Order{
// 		{
// 			Id:         1,
// 			SquareId:   "9utRtTF6Qwzt01tEiCHs92O3Rj4C",
// 			Products:   products,
// 			LocationId: "Default",
// 			State:      "Completed",
// 			Version:    1,
// 			TotalMoney: 5.98,
// 			CreatedAt:  time.Now(),
// 		},
// 	}
// 	lwClient.CreateOrders(orders)
// }

func main() {

	port := getEnv("HTTP_PORT")

	logger := logInit()
	logger.Infow("starting http server", "port", port)

	dbPath := getEnv("DB")

	signingKey := []byte(getEnv("SIGNING_KEY"))

	lwAppId := getEnv("LINNWORKS_APP_ID")
	lwAppSecret := getEnv("LINNWORKS_APP_SECRET")
	lwAppToken := getEnv("LINNWORKS_APP_TOKEN")

	sqAccessToken := getEnv("SQUARE_ACCESS_TOKEN")
	sqHost := getEnv("SQUARE_HOST")
	sqApiVersion := getEnv("SQUARE_API_VERSION")
	sqLocationId := getEnv("SQUARE_LOCATION_ID")

	logger.Debugw("loaded application config",
		"linnworksAppId", lwAppId,
		"linnworksAppSecret", lwAppSecret,
		"linnworksAppToken", lwAppToken,
		"squareAccessToken", sqAccessToken,
		"sqHost", sqHost,
		"sqApiVersion", sqApiVersion,
		"sqLocationId", sqLocationId,
		"signingKey", signingKey,
		"db", dbPath,
	)

	sqliteDb := sqlite.NewSqliteDB(dbPath)
	lwClient := linnworks.NewLinnworksClient(lwAppId, lwAppSecret, lwAppToken)
	sqClient := square.NewSquareClient(sqAccessToken, sqHost, sqApiVersion, sqLocationId)
	syncTool := sync.NewSyncTool(lwClient, sqClient, sqliteDb)

	authMiddleware := middlewares.NewAuthMiddleware(signingKey, logger)
	authController := controllers.NewAuthController(sqliteDb, signingKey, logger)
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
		logger.Fatalw("unable to start server", "error", err.Error())
	}
}

func getEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		panic("cannot find .env file")
	}

	return os.Getenv(key)
}

func logInit() *zap.SugaredLogger {

	f, err := os.Create("logs/test.log")
	if err != nil {
		panic("unable to open log file")
	}

	pe := zap.NewProductionEncoderConfig()

	fileEncoder := zapcore.NewJSONEncoder(pe)

	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	level := zap.DebugLevel

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(f), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	l := zap.New(core)

	return l.Sugar()
}
