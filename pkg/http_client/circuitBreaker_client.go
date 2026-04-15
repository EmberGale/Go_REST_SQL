package http_client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type CircuitBreaker struct {
	delegate HTTPClient
	breaker  *gobreaker.CircuitBreaker
}

func NewCircuitBreakerClient(delegate HTTPClient, settings gobreaker.Settings) *CircuitBreaker {
	if settings.Name == "" {
		settings.Name = "default-circuit-breaker"
	}
	if settings.Timeout == 0 {
		settings.Timeout = 60 * time.Second
	}
	if settings.MaxRequests == 0 {
		settings.MaxRequests = 1
	}

	if settings.Interval == 0 {
		settings.Interval = 5 * time.Second
	}
	if settings.ReadyToTrip == nil {
		settings.ReadyToTrip = func(counts gobreaker.Counts) bool {
			failureRation := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRation >= 0.6
		}
	}

	return &CircuitBreaker{
		delegate: delegate,
		breaker:  gobreaker.NewCircuitBreaker(settings),
	}
}

func (c *CircuitBreaker) Do(req *http.Request) (*http.Response, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		resp, err := c.delegate.Do(req)
		if err != nil {
			// Timeouts, Retries from delegate
			return nil, err
		}

		if resp.StatusCode >= 500 {
			return resp, fmt.Errorf("Server error: %d", resp.StatusCode)
		}
		return resp, nil
	})

	if err != nil {
		if err == gobreaker.ErrOpenState {
			return nil, fmt.Errorf("circuit breaker is open: %w", err)
		}
		if resp, ok := result.(*http.Response); ok && resp != nil {
			return resp, err
		}
		return nil, err
	}

	return result.(*http.Response), nil
}
