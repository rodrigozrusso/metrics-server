package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"acme.inc/analytics/internal/metrics"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Position struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func main() {

	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	err := godotenv.Load()
	if err != nil {
		zap.L().Fatal("Error loading .env file")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// random metric names
	metrics := []string{
		"metric_a", "metric_b", "metric_c", "metric_d", "metric_e", "metric_f", "metric_g", "metric_h", "metric_i", "metric_j",
		"metric_k", "metric_l", "metric_m", "metric_n", "metric_o", "metric_p", "metric_q", "metric_r", "metric_s", "metric_t",
	}

	for i := 0; i < len(metrics); i++ {
		go send(metrics[i])
	}
	// <-ctx.Done()
	// zap.L().Info("got interruption signal")
	// shutdown = true
	// zap.L().Info("final")

	select {
	case <-ctx.Done():
		zap.L().Info("got interruption signal")
		return
	}
}

func generateData(metricName string) metrics.Metric {
	metric := metrics.Metric{
		Timestamp: time.Date(gofakeit.Number(2022, 2023), time.Month(gofakeit.Number(1, 12)), gofakeit.Number(1, 28), gofakeit.Number(0, 23), gofakeit.Number(0, 59), gofakeit.Number(0, 59), 0, time.UTC),
		Name:      metricName,
		Value:     gofakeit.Float64Range(20, 40),
	}
	return metric
}

func send(metricName string) {
	baseURL := os.Getenv("ANALYTICS_SERVER_URL")
	url := fmt.Sprintf("%s/v1/metrics", baseURL)

	for true {
		metric := generateData(metricName)
		json_data, err := json.Marshal(metric)
		if err != nil {
			zap.L().Fatal("Error marshallin data",
				zap.Any("metric", metric),
				zap.Error(err))
		}
		makePost(url, json_data)
	}
}

func makePost(url string, json_data []byte) bool {
	zap.L().Info("Sending position", zap.String("url", url))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		zap.L().Fatal("Error posting vin position",
			zap.Error(err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		zap.L().Fatal("Position could not be processed on server")
		return false
	}
	// time.Sleep(1 * time.Second)
	time.Sleep(300 * time.Millisecond)
	return true
}
