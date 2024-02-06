package handler

import (
	"fmt"
	"net/http"
	"time"

	"acme.inc/analytics/internal/domain"
	"acme.inc/analytics/internal/service"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// NewCommandHandler creates a new command handler that handles HTTP requests for adding metrics.
// It registers a POST route "/metrics" on the provided Fiber app and uses the given service to add the metric.
func NewHandler(app fiber.Router, service service.Service) {
	app.Post("/v1/metrics", addMetric(service))
	app.Get("/v1/metrics", listMetrics(service))
	app.Get("/v1/metrics/:metricName/:granularity/:startDate/:endDate", getDataByMetricName(service))
}

// addMetric is a Fiber handler function that parses the request body into a Metric struct,
// logs the metric using Zap logger, and adds the metric using the provided service.
// It returns a Fiber error if there is an error parsing the request body or adding the metric.
func addMetric(service service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var metric domain.Metric
		// metric := new(domain.Metric)
		if err := c.BodyParser(&metric); err != nil {
			return c.Status(503).Send([]byte(err.Error()))
		}
		fmt.Println(metric.Name)
		zap.L().Info("Metric", zap.Any("metric", metric))

		if err := service.AddMetric(metric); err != nil {
			c.Status(http.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusCreated)
	}
}

func listMetrics(service service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		metrics, err := service.ListMetrics()
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(metrics)
	}
}

// curl -X GET "http://localhost:3000/v1/metrics/temperature/1/2021-01-01/2021-01-02"
func getDataByMetricName(service service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		metricName := c.Params("metricName")
		granularity := c.Params("granularity")
		startDateStr := c.Params("startDate")
		endDateStr := c.Params("endDate")

		const layout = "2006-01-02"
		startDate, err := time.Parse(layout, startDateStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("Invalid startDate format")
		}

		endDate, err := time.Parse(layout, endDateStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("Invalid endDate format")
		}

		metrics, err := service.GetDataByMetricName(metricName, domain.Granularity(granularity), startDate, endDate)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(metrics)
	}
}
