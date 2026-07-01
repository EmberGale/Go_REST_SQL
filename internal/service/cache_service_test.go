package service

import (
	"GoRestSQL/internal/model"
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFromCache_Hit(t *testing.T) {
	svc, _, mockCache := setupPaymentService(t)
	ctx := context.Background()
	_, err := mockCache.Set(ctx, "payment_id:1", `{"id":1,"person":"Alice"}`, 0).Result()
	require.NoError(t, err)

	data, err := svc.getFromCache(ctx, "payment_id:1")

	assert.NoError(t, err)
	assert.Equal(t, `{"id":1,"person":"Alice"}`, data)
}

func TestGetFromCache_Miss(t *testing.T) {
	svc, _, _ := setupPaymentService(t)
	ctx := context.Background()

	data, err := svc.getFromCache(ctx, "payment_id:missing")

	assert.ErrorIs(t, err, redis.Nil)
	assert.Nil(t, data)
}

func TestWarmCache_Success(t *testing.T) {
	svc, mockRepo, _ := setupPaymentService(t)
	ctx := context.Background()
	payments := []model.Payment{*samplePayment(1), *samplePayment(2)}

	mockRepo.EXPECT().GetAllPayments(ctx).Return(payments, nil)

	err := svc.WarmCache(ctx)

	assert.NoError(t, err)
}

func TestWarmCache_RepoError(t *testing.T) {
	svc, mockRepo, _ := setupPaymentService(t)
	ctx := context.Background()
	repoErr := errors.New("db unavailable")

	mockRepo.EXPECT().GetAllPayments(ctx).Return(nil, repoErr)

	err := svc.WarmCache(ctx)

	assert.ErrorIs(t, err, repoErr)
}

func TestGetWithSingleFlight_Success(t *testing.T) {
	svc, mockRepo, _ := setupPaymentService(t)
	payment := samplePayment(10)
	ctx := context.Background()

	mockRepo.EXPECT().GetById(ctx, int64(10)).Return(payment, nil).Once()

	result, err := svc.GetWithSingleFlight(ctx, 10)

	require.NoError(t, err)
	assert.Equal(t, payment.Id, result.Id)
	assert.Equal(t, payment.Person, result.Person)
}

func TestGetWithSingleFlight_DeduplicatesConcurrentCalls(t *testing.T) {
	svc, mockRepo, _ := setupPaymentService(t)
	payment := samplePayment(11)
	ctx := context.Background()

	mockRepo.EXPECT().GetById(ctx, int64(11)).Return(payment, nil).Once()

	start := make(chan struct{})
	results := make(chan struct {
		payment *model.Payment
		err     error
	}, 2)

	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			<-start
			result, err := svc.GetWithSingleFlight(ctx, 11)
			results <- struct {
				payment *model.Payment
				err     error
			}{payment: result, err: err}
		}()
	}

	close(start)
	wg.Wait()
	close(results)

	var received []struct {
		payment *model.Payment
		err     error
	}
	for res := range results {
		received = append(received, res)
	}

	require.Len(t, received, 2)
	for _, res := range received {
		require.NoError(t, res.err)
		assert.Equal(t, payment.Id, res.payment.Id)
	}
}
