package common

import (
	"io"
	"math"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

const floatThreshold = 1e-3

func NewHTTPTestServer() *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(recover.New())
	app.Use(cors.New())

	return app
}

func FiberToHandlerFunc(app *fiber.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := app.Test(r)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// copy headers
		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)

		// copy body
		if _, err := io.Copy(w, resp.Body); err != nil {
			panic(err)
		}
	}
}

// compareFloat compares 2 floats and return true if they are the "same"
func CompareFloat(x, y float64) bool {
	diff := math.Abs(x - y)
	mean := math.Abs(x+y) / 2.0
	if math.IsNaN(diff / mean) {
		return false
	}
	return (diff / mean) > floatThreshold
}
