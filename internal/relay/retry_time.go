package relay

import (
	"math"
	"math/rand"
	"time"
)

// CalculateNextRetryTime вычисляет время следующей попытки с экспоненциальной задержкой и джиттером
func CalculateNextRetryTime(attempt int) time.Time {
	baseDelay := 1 * time.Second
	maxDelay := 30 * time.Second

	// Экспоненциальный рост: 1s, 2s, 4s, 8s...
	delay := float64(baseDelay) * math.Pow(2, float64(attempt))

	if delay > float64(maxDelay) {
		delay = float64(maxDelay)
	}

	// Jitter: рандомизация +-20% от delay
	// rand.Float64() возвращает [0.0, 1.0).
	// (rand.Float64()*2 - 1) возвращает [-1.0, 1.0).
	// Умножаем на 0.2 (20% отклонения).
	jitterFactor := rand.Float64()*2 - 1
	jitter := delay * 0.2 * jitterFactor

	finalDelay := time.Duration(delay + jitter)

	// Защита от отрицательной задержки из-за джиттера
	if finalDelay < 0 {
		finalDelay = 0
	}

	return time.Now().Add(finalDelay)
}
