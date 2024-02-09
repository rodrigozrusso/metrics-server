package main

import (
	"fmt"
	"log"

	"acme.inc/analytics/internal/common"
	internalError "acme.inc/analytics/internal/error"
	"acme.inc/analytics/internal/healthcheck"
	"acme.inc/analytics/internal/metrics"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func main() {
	{
		logger, _ := zap.NewProduction()
		zap.ReplaceGlobals(logger)
		defer logger.Sync()
	}

	zap.L().Info("Starting analytics server")
	// database
	config := common.NewConfiguration()
	// db, err :=
	common.NewTimescaleDB(config.DataBaseConfig)
	db, err := common.NewGormDB(config.DataBaseConfig)
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&metrics.Metric{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// http server
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// app.Use((handler.ErrorHandler))

	// healthcheck
	healthcheck.NewHandler(app, sqlDB)

	v1 := app.Group("/v1")
	// metrics
	repository := metrics.NewRepository(db)
	repository.CreateHyperTable()
	metricService := metrics.NewService(repository)
	metrics.NewHandler(v1, metricService)

	app.Use((internalError.NotFound))

	serverPort := fmt.Sprintf(":%s", config.ServerConfig.Port)
	log.Fatal(app.Listen(serverPort))
}
