package common

import (
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGormDB create a new connection with database
func NewGormDB(dataBaseConfig *DataBaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d", dataBaseConfig.Host, dataBaseConfig.Port, dataBaseConfig.User, dataBaseConfig.Password, dataBaseConfig.Name, dataBaseConfig.Timeout)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			utc, _ := time.LoadLocation("")
			return time.Now().In(utc)
		}})

	if err != nil {
		zap.L().Error("Error to connect to the database", zap.String("dsn", dsn))
		return nil, err
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(dataBaseConfig.MinConnections)
	sqlDB.SetMaxOpenConns(dataBaseConfig.MaxConnections)
	sqlDB.SetConnMaxLifetime(time.Duration(30 * time.Minute))

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
