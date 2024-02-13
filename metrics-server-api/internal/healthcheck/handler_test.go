package healthcheck

import (
	"testing"

	"acme.inc/analytics/internal/common"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetPing(t *testing.T) {

	app := common.NewHTTPTestServer()
	NewHandler(app, nil)
	apitest.New().
		HandlerFunc(common.FiberToHandlerFunc(app)).
		Get("/ping").
		Expect(t).
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

		app := common.NewHTTPTestServer()
		NewHandler(app, sqlDB)

		apitest.New().
			HandlerFunc(common.FiberToHandlerFunc(app)).
			Get("/health").
			Expect(t).
			Assert(jsonpath.Equal(`status`, "up")).
			Status(200).
			End()
	})

	t.Run("Database is down", func(t *testing.T) {

		app := common.NewHTTPTestServer()
		NewHandler(app, nil)

		apitest.New().
			HandlerFunc(common.FiberToHandlerFunc(app)).
			Get("/health").
			Expect(t).
			Assert(jsonpath.Equal(`status`, "down")).
			Status(200).
			End()
	})
}
