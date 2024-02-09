package error

import fiber "github.com/gofiber/fiber/v2"

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).JSON(&FailedResponse{Message: "Resource Not Found"})
}
