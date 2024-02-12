package metrics_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"acme.inc/analytics/internal/common"
	"acme.inc/analytics/internal/metrics"
	mocks "acme.inc/analytics/internal/metrics/mock"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"go.uber.org/mock/gomock"
)

func TestGetDataByMetricName(t *testing.T) {

	t.Run("Valid Request", func(t *testing.T) {
		app := common.NewHTTPTestServer()
		// Prepare test data
		metricName := "visitors"
		granularity := "daily"
		startDateStr := "2024-02-01"
		endDateStr := "2024-02-07"
		avgMetricResponse := metrics.AVGMetricResponse{
			MetricName:  "visitors",
			Granularity: "daily",
			StartTime: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2024-02-01T01:00:00.000+0100")
				return t
			}(),
			EndTime: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2024-02-07T01:00:00.000+0100")
				return t
			}(),
			Data: []metrics.MetricAVGResult{
				{TimeFrame: func() time.Time {
					t, _ := time.Parse(time.RFC3339, "2024-02-01T01:00:00.000+0100")
					return t
				}(), Avg: 100},
				{TimeFrame: func() time.Time {
					t, _ := time.Parse(time.RFC3339, "2024-02-02T01:00:00.000+0100")
					return t
				}(), Avg: 100},
				{TimeFrame: func() time.Time {
					t, _ := time.Parse(time.RFC3339, "2024-02-03T01:00:00.000+0100")
					return t
				}(), Avg: 150},
				{TimeFrame: func() time.Time {
					t, _ := time.Parse(time.RFC3339, "2024-02-04T01:00:00.000+0100")
					return t
				}(), Avg: 200},
				{TimeFrame: func() time.Time {
					t, _ := time.Parse(time.RFC3339, "2024-02-05T01:00:00.000+0100")
					return t
				}(), Avg: 180},
				{TimeFrame: func() time.Time {
					t, _ := time.Parse(time.RFC3339, "2024-02-06T01:00:00.000+0100")
					return t
				}(), Avg: 220},
				{TimeFrame: func() time.Time {
					t, _ := time.Parse(time.RFC3339, "2024-02-07T01:00:00.000+0100")
					return t
				}(), Avg: 190},
			},
		}

		const layout = "2006-01-02"

		ctrl := gomock.NewController(t)
		mockMetricService := mocks.NewMockService(ctrl)
		mockMetricService.EXPECT().GetDataByMetricName(
			metricName,
			metrics.Granularity(granularity),
			func() time.Time {
				t, _ := time.Parse(layout, startDateStr)
				return t
			}(),
			func() time.Time {
				t, _ := time.Parse(layout, endDateStr)
				return t
			}()).
			Return(avgMetricResponse, nil)

		url := fmt.Sprintf("/metrics/%s/%s/%s/%s", metricName, granularity, startDateStr, endDateStr)

		metrics.NewHandler(app, mockMetricService)
		apitest.New().
			HandlerFunc(common.FiberToHandlerFunc(app)).
			Get(url).
			Expect(t).
			Assert(jsonpath.Equal(`$.metricName`, metricName)).
			Status(200).
			End()
	})

	t.Run("Invalid Start Date", func(t *testing.T) {
		metricName := "visitors"
		granularity := "daily"
		startDateStr := "2024-Feb-01"
		endDateStr := "2024-02-07"

		const layout = "2006-01-02"

		ctrl := gomock.NewController(t)
		mockMetricService := mocks.NewMockService(ctrl)
		// mockMetricService.EXPECT().GetDataByMetricName(
		// 	metricName,
		// 	metrics.Granularity(granularity),
		// 	func() time.Time {
		// 		t, _ := time.Parse(layout, startDateStr)
		// 		return t
		// 	}(),
		// 	func() time.Time {
		// 		t, _ := time.Parse(layout, endDateStr)
		// 		return t
		// 	}()).
		// 	Return(&[]metrics.MetricsResponse{}, nil)

		url := fmt.Sprintf("/metrics/%s/%s/%s/%s", metricName, granularity, startDateStr, endDateStr)

		app := common.NewHTTPTestServer()
		metrics.NewHandler(app, mockMetricService)
		apitest.New().
			HandlerFunc(common.FiberToHandlerFunc(app)).
			Get(url).
			Expect(t).
			Assert(jsonpath.Equal(`$.message`, "Invalid startDate format")).
			Status(400).
			End()

	})

	t.Run("Invalid End Date", func(t *testing.T) {
		metricName := "visitors"
		granularity := "daily"
		startDateStr := "2024-02-01"
		endDateStr := "2024-"

		ctrl := gomock.NewController(t)
		mockMetricService := mocks.NewMockService(ctrl)

		url := fmt.Sprintf("/metrics/%s/%s/%s/%s", metricName, granularity, startDateStr, endDateStr)

		app := common.NewHTTPTestServer()
		metrics.NewHandler(app, mockMetricService)
		apitest.New().
			HandlerFunc(common.FiberToHandlerFunc(app)).
			Get(url).
			Expect(t).
			Assert(jsonpath.Equal(`$.message`, "Invalid endDate format")).
			Status(400).
			End()
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		metricName := "visitors"
		granularity := "daily"
		startDateStr := "2024-02-01"
		endDateStr := "2024-02-07"
		avgMetricResponse := metrics.AVGMetricResponse{}
		expectedError := errors.New("database error")

		const layout = "2006-01-02"

		ctrl := gomock.NewController(t)
		mockMetricService := mocks.NewMockService(ctrl)
		mockMetricService.EXPECT().GetDataByMetricName(
			metricName,
			metrics.Granularity(granularity),
			func() time.Time {
				t, _ := time.Parse(layout, startDateStr)
				return t
			}(),
			func() time.Time {
				t, _ := time.Parse(layout, endDateStr)
				return t
			}()).
			Return(avgMetricResponse, expectedError)

		url := fmt.Sprintf("/metrics/%s/%s/%s/%s", metricName, granularity, startDateStr, endDateStr)

		app := common.NewHTTPTestServer()
		metrics.NewHandler(app, mockMetricService)
		apitest.New().
			HandlerFunc(common.FiberToHandlerFunc(app)).
			Get(url).
			Expect(t).
			Assert(jsonpath.Equal(`$.message`, "Well, This is unexpected. An Error has occurred, and we are working to fix the problem!")).
			Status(500).
			End()
	})
}

func TestAddMetric(t *testing.T) {
	t.Run("Valid Request", func(t *testing.T) {
		app := common.NewHTTPTestServer()
		// Prepare test data
		timestamp, _ := time.Parse(time.RFC3339, "2024-02-01T15:04:05Z07:00")
		metric := metrics.Metric{
			Name:      "visitors",
			Value:     100,
			Timestamp: timestamp,
		}

		ctrl := gomock.NewController(t)
		mockMetricService := mocks.NewMockService(ctrl)
		mockMetricService.EXPECT().AddMetric(metric).Return(nil)

		url := fmt.Sprintf("/metrics")

		metrics.NewHandler(app, mockMetricService)
		apitest.New().
			HandlerFunc(common.FiberToHandlerFunc(app)).
			Post(url).
			JSON(metric).
			Expect(t).
			Status(201).
			End()
	})

	t.Run("Invalid Request", func(t *testing.T) {

		// Prepare test data
		timestamp, _ := time.Parse(time.RFC3339, "2024-02-01T15:04:05Z07:00")
		metric := metrics.Metric{
			Name:      "visitors",
			Value:     100,
			Timestamp: timestamp,
		}

		ctrl := gomock.NewController(t)
		mockMetricService := mocks.NewMockService(ctrl)
		mockMetricService.EXPECT().AddMetric(metric).Return(errors.New("database error"))

		url := fmt.Sprintf("/metrics")

		app := common.NewHTTPTestServer()
		metrics.NewHandler(app, mockMetricService)
		apitest.New().
			HandlerFunc(common.FiberToHandlerFunc(app)).
			Post(url).
			JSON(metric).
			Expect(t).
			Status(500).
			End()
	})
}

func TestListMetrics(t *testing.T) {
	t.Run("Valid Request", func(t *testing.T) {

		// Prepare test data
		metricsResponse := []string{"visitors", "temperature"}

		ctrl := gomock.NewController(t)
		mockMetricService := mocks.NewMockService(ctrl)
		mockMetricService.EXPECT().ListMetrics().Return(metricsResponse, nil)

		url := fmt.Sprintf("/metrics")

		app := common.NewHTTPTestServer()
		metrics.NewHandler(app, mockMetricService)
		apitest.New().
			HandlerFunc(common.FiberToHandlerFunc(app)).
			Get(url).
			Expect(t).
			Status(200).
			End()
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		app := common.NewHTTPTestServer()
		// Prepare test data
		metricsResponse := []string{"visitors", "temperature"}

		ctrl := gomock.NewController(t)
		mockMetricService := mocks.NewMockService(ctrl)
		mockMetricService.EXPECT().ListMetrics().Return(metricsResponse, errors.New("database error"))

		url := fmt.Sprintf("/metrics")

		metrics.NewHandler(app, mockMetricService)
		apitest.New().
			HandlerFunc(common.FiberToHandlerFunc(app)).
			Get(url).
			Expect(t).
			Status(500).
			End()
	})
}
