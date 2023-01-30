package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	gormsqlite "github.com/kampanosg/go-lsi/clients/db/gorm_sqlite"
	"github.com/kampanosg/go-lsi/clients/linnworks"
	"github.com/kampanosg/go-lsi/clients/square"
	"github.com/kampanosg/go-lsi/controllers"
	"github.com/kampanosg/go-lsi/middlewares"
	"github.com/kampanosg/go-lsi/sync"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	port := os.Getenv("PORT")

	logger := logInit()
	logger.Infow("starting http server", "port", port)

	dbPath := os.Getenv("DB")

	signingKey := []byte(os.Getenv("SIGNING_KEY"))

	lwAppId := os.Getenv("LINNWORKS_APP_ID")
	lwAppSecret := os.Getenv("LINNWORKS_APP_SECRET")
	lwAppToken := os.Getenv("LINNWORKS_APP_TOKEN")

	sqAccessToken := os.Getenv("SQUARE_ACCESS_TOKEN")
	sqHost := os.Getenv("SQUARE_HOST")
	sqApiVersion := os.Getenv("SQUARE_API_VERSION")
	sqLocationId := os.Getenv("SQUARE_LOCATION_ID")

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

	sqliteDb, err := gormsqlite.NewSqliteDb(dbPath)
	if err != nil {
		logger.Fatalw("cannot connect to the database", "error", err.Error())
	}

	lwClient := linnworks.NewLinnworksClient(lwAppId, lwAppSecret, lwAppToken, logger)
	sqClient := square.NewSquareClient(sqAccessToken, sqHost, sqApiVersion, sqLocationId, logger)
	syncTool := sync.NewSyncTool(lwClient, sqClient, sqliteDb, logger)

	authMiddleware := middlewares.NewAuthMiddleware(signingKey, logger)
	authController := controllers.NewAuthController(sqliteDb, signingKey, logger)
	inventoryController := controllers.NewInventoryController(sqliteDb, logger)
	ordersController := controllers.NewOrdersController(sqliteDb, logger)
	pingController := controllers.NewPingController()
	syncController := controllers.NewSyncController(syncTool, logger)

	router := mux.NewRouter()

	router.Handle("/api/v1/ping", authMiddleware.ProtectedEndpoint(http.HandlerFunc(pingController.HandlePingRequest)))
	router.Handle("/api/v1/inventory", authMiddleware.ProtectedEndpoint(http.HandlerFunc(inventoryController.HandleInventoryRequest)))
	router.Handle("/api/v1/orders", authMiddleware.ProtectedEndpoint(http.HandlerFunc(ordersController.HandleOrdersRequest)))
	router.Handle("/api/v1/sync/status", authMiddleware.ProtectedEndpoint(http.HandlerFunc(syncController.HandleSyncStatusRequest)))
	router.Handle("/api/v1/sync/recent", authMiddleware.ProtectedEndpoint(http.HandlerFunc(syncController.HandleSyncRecentRequest)))
	router.Handle("/api/v1/sync", authMiddleware.ProtectedEndpoint(http.HandlerFunc(syncController.HandleSyncRequest)))
	router.Handle("/api/v1/auth", http.HandlerFunc(authController.HandleAuthRequest))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		logger.Fatalw("unable to start server", "error", err.Error())
	}
}

func logInit() *zap.SugaredLogger {

	pe := zap.NewProductionEncoderConfig()

	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	level := zap.InfoLevel

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	l := zap.New(core)

	return l.Sugar()
}
