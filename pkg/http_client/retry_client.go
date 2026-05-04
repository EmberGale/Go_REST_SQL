package http_client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type RetryConfig struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration

	ShouldRetry func(*http.Response, error) bool
}

type RetryClient struct {
	delegate HTTPClient
	config   RetryConfig
	logger   *zap.Logger
}

func NewRetryClient(delegate HTTPClient, config RetryConfig, logger *zap.Logger) *RetryClient {
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.InitialDelay == 0 {
		config.InitialDelay = 100 * time.Millisecond
	}
	if config.MaxDelay == 0 {
		config.MaxDelay = 5 * time.Second
	}
	if config.ShouldRetry == nil {
		config.ShouldRetry = func(resp *http.Response, err error) bool {
			return err != nil || resp.StatusCode >= 500
		}
	}
	return &RetryClient{
		delegate: delegate,
		config:   config,
		logger:   logger,
	}
}

func (c *RetryClient) Do(req *http.Request) (*http.Response, error) {
	var (
		resp  *http.Response
		err   error
		delay = c.config.InitialDelay
	)

	for i := 0; i < c.config.MaxRetries; i++ {
		if req.Body != nil {
			if seeker, ok := req.Body.(io.Seeker); ok {
				_, SeekerErr := seeker.Seek(0, io.SeekStart)
				if SeekerErr != nil {
					return nil, SeekerErr
				}
			} else {
				bodyBytes, readErr := io.ReadAll(req.Body)
				if readErr != nil {
					return nil, fmt.Errorf("failed to read request body for retry: %w", readErr)
				}
				req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}

		resp, err = c.delegate.Do(req)

		if c.config.ShouldRetry(resp, err) {
			c.logger.Info("Request failed, retrying",
				zap.Int("attempt", i+1),
				zap.Int("max_retries", c.config.MaxRetries),
				zap.Duration("retry_delay", delay),
				zap.Error(err),
			)

			time.Sleep(delay)
			delay *= 2
			if delay > c.config.MaxDelay {
				delay = c.config.MaxDelay
			}
			continue
		}
		return resp, err
	}
	return resp, err
}
