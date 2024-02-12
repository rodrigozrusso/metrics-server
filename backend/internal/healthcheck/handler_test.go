package healthcheck

import (
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetPing(t *testing.T) {

	app := newHTTPServer()
	NewHandler(app, nil)
	apitest.New().
		HandlerFunc(FiberToHandlerFunc(app)).
		Get("/ping").
		Expect(t).
		// Body(`{"status":"up"}`).
		Status(200).
		End()
}

func TestGetHealth(t *testing.T) {

	t.Run("Database is up", func(t *testing.T) {

		mockDb, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDb,
			DriverName: "postgres",
		})
		defer mockDb.Close()
		db, _ := gorm.Open(dialector, &gorm.Config{})
		sqlDB, _ := db.DB()
		defer sqlDB.Close()
		mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))

		app := newHTTPServer()
		NewHandler(app, sqlDB)

		apitest.New().
			HandlerFunc(FiberToHandlerFunc(app)).
			Get("/health").
			Expect(t).
			Assert(jsonpath.Equal(`status`, "up")).
			Status(200).
			End()
	})

	t.Run("Database is down", func(t *testing.T) {

		mockDb, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDb,
			DriverName: "postgres",
		})
		defer mockDb.Close()
		db, _ := gorm.Open(dialector, &gorm.Config{})
		sqlDB, _ := db.DB()
		defer sqlDB.Close()
		mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))

		app := newHTTPServer()
		NewHandler(app, nil)

		apitest.New().
			HandlerFunc(FiberToHandlerFunc(app)).
			Get("/health").
			Expect(t).
			Assert(jsonpath.Equal(`status`, "down")).
			Status(200).
			End()
	})
}

func newHTTPServer() *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

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
