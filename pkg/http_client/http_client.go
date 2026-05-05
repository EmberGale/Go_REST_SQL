package http_client

import (
	"net/http"
	"time"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

func NewHTTPClient(logger *zap.Logger) HTTPClient {
	baseClient := NewDefaultHTTPClient(5 * time.Second)

	retryConfig := RetryConfig{
		MaxRetries:   5,
		InitialDelay: 200 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		ShouldRetry: func(resp *http.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable
		},
	}

	retryClient := NewRetryClient(baseClient, retryConfig, logger)

	cbSettings := gobreaker.Settings{
		Name:        "ExternalService",
		MaxRequests: 3,
		Interval:    5 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 10 && failureRatio >= 0.3
		},
	}

	client := NewCircuitBreakerClient(retryClient, cbSettings)

	return client
}
