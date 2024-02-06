package main

import (
	"fmt"
	"log"

	handler "acme.inc/analytics/internal/api/v1"
	"acme.inc/analytics/internal/common"
	"acme.inc/analytics/internal/domain"
	"acme.inc/analytics/internal/repository"
	"acme.inc/analytics/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
)

func main() {
	{
		logger, _ := zap.NewProduction()
		zap.ReplaceGlobals(logger)
		defer logger.Sync()
	}

	log.Printf("Starting analytics server")
	// database
	config := common.NewConfiguration()
	// db, err := common.NewTimescaleDB(config.DataBaseConfig)
	db, err := common.NewGormDB(config.DataBaseConfig)
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&domain.Metric{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// http server
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(logger.New())

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// services
	repository := repository.NewRepository(db)
	repository.CreateHyperTable()
	metricService := service.NewService(repository)
	handler.NewHandler(app, metricService)
	// vinPositionEventHandler := service.NewEventHandler(broker, stateRepo)
	// stateService := service.NewService(stateRepo, vinPositionEventHandler)
	// handler.NewCommandHandler(app, stateService)
	// handler.NewQuerydHandler(app, stateService)

	// go vinPositionEventHandler.Consume()

	serverPort := fmt.Sprintf(":%s", config.ServerConfig.Port)
	log.Fatal(app.Listen(serverPort))
}
