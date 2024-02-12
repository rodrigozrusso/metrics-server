package metrics

import (
	"net/http"
	"time"

	"acme.inc/analytics/internal/internalError"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// NewCommandHandler creates a new command handler that handles HTTP requests for adding metrics.
// It registers a POST route "/metrics" on the provided Fiber app and uses the given service to add the metric.
func NewHandler(app fiber.Router, service Service) {
	metrics := app.Group("/metrics")

	metrics.Post("/", addMetric(service))
	metrics.Get("/", listMetrics(service))
	metrics.Get("/:metricName/:granularity/:startDate/:endDate", getDataByMetricName(service))
}

// addMetric is a Fiber handler function that parses the request body into a Metric struct,
// logs the metric using Zap logger, and adds the metric using the provided service.
// It returns a Fiber error if there is an error parsing the request body or adding the metric.
func addMetric(service Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var metric Metric
		if err := c.BodyParser(&metric); err != nil {
			return c.Status(http.StatusBadRequest).JSON(&internalError.FailedResponse{Message: "Invalid Payload. Please verify the request payload and try again."})
		}
		zap.L().Info("Metric", zap.Any("metric", metric))

		err := service.AddMetric(metric)
		if err != nil {
			zap.L().Error("Metric", zap.Any("metric", metric), zap.Error(err))
			return c.Status(http.StatusInternalServerError).JSON(&internalError.FailedResponse{Message: "Well, This is unexpected. An Error has occurred, and we are working to fix the problem!"})
		}

		return c.SendStatus(fiber.StatusCreated)
	}
}

func listMetrics(service Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		metrics, err := service.ListMetrics()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&internalError.FailedResponse{Message: "Well, This is unexpected. An Error has occurred, and we are working to fix the problem!"})
		}
		return c.JSON(metrics)
	}
}

// curl -X GET "http://localhost:3000/v1/metrics/temperature/1/2021-01-01/2021-01-02"
func getDataByMetricName(service Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		metricName := c.Params("metricName")
		granularity := c.Params("granularity")
		startDateStr := c.Params("startDate")
		endDateStr := c.Params("endDate")

		zap.L().Debug("getDataByMetricName", zap.String("metricName", metricName), zap.String("granularity", granularity), zap.String("startDate", startDateStr), zap.String("endDate", endDateStr))

		const layout = "2006-01-02"
		startDate, err := time.Parse(layout, startDateStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(&internalError.FailedResponse{Message: "Invalid startDate format"})
		}

		endDate, err := time.Parse(layout, endDateStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(&internalError.FailedResponse{Message: "Invalid endDate format"})
		}

		metricsResponse, err := service.GetDataByMetricName(metricName, Granularity(granularity), startDate, endDate)
		if err != nil {
			zap.L().Error("Error to retrieve data",
				zap.Error(err),
				zap.String("metricName", metricName),
				zap.String("granularity", granularity),
				zap.String("startDate", startDateStr),
				zap.String("endDate", startDateStr))
			return c.Status(http.StatusInternalServerError).JSON(&internalError.FailedResponse{Message: "Well, This is unexpected. An Error has occurred, and we are working to fix the problem!"})
		}
		return c.JSON(metricsResponse)
	}
}
