package healthcheck

import (
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"
)

func NewHandler(app fiber.Router, sqlDB *sql.DB) {
	app.Get("/ping", ping())
	app.Get("/health", health(sqlDB))
}

func ping() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": Ok})
	}
}

func health(sqlDB *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		response := HealthCheckResponse{
			Status: Ok,
		}

		dbHealthCheckResp := ResourceHealthCheckResponse{
			Name: "database",
		}

		// check db health
		err := sqlDB.Ping()
		if err != nil {
			dbHealthCheckResp.Status = Fail
			dbHealthCheckResp.Message = err.Error()
			response.Status = Warm // warm that db is not healthy
		} else {
			dbHealthCheckResp.Status = Ok
		}

		response.Resources = append(response.Resources, dbHealthCheckResp)

		return c.Status(200).JSON(&response)
	}
}
