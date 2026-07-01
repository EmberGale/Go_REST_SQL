package service

import (
	"GoRestSQL/internal/model"
	repomocks "GoRestSQL/internal/repository/mocks"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockCache struct {
	data map[string]string
}

func newMockCache() *mockCache {
	return &mockCache{data: make(map[string]string)}
}

func (m *mockCache) Get(ctx context.Context, key string) *redis.StringCmd {
	if value, ok := m.data[key]; ok {
		return redis.NewStringResult(value, nil)
	}
	return redis.NewStringResult("", redis.Nil)
}

func (m *mockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	switch v := value.(type) {
	case []byte:
		m.data[key] = string(v)
	case string:
		m.data[key] = v
	default:
		bytes, _ := json.Marshal(v)
		m.data[key] = string(bytes)
	}
	return redis.NewStatusResult("OK", nil)
}

func (m *mockCache) Del(ctx context.Context, key ...string) *redis.IntCmd {
	for _, k := range key {
		delete(m.data, k)
	}
	return redis.NewIntResult(int64(len(key)), nil)
}

func setupPaymentService(t *testing.T) (*PaymentServiceImpl, *repomocks.MockPaymentRepository, *mockCache) {
	t.Helper()

	mockRepo := repomocks.NewMockPaymentRepository(t)
	mockCache := newMockCache()

	svc := NewPaymentServiceImpl(context.Background(), mockRepo, nil, zap.NewNop().Sugar(), mockCache)
	return svc, mockRepo, mockCache
}

func samplePayment(id int64) *model.Payment {
	return &model.Payment{
		Id:       id,
		Person:   "Alice",
		Amount:   100.50,
		Currency: "USD",
		Status:   model.PAYMENT_STATUS_CREATED,
		Time:     time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
		UserID:   42,
	}
}

func TestCreatePayment_Success(t *testing.T) {
	svc, mockRepo, _ := setupPaymentService(t)
	payment := samplePayment(0)
	ctx := context.Background()

	mockRepo.EXPECT().CreateWithOutbox(ctx, payment).Return(int64(1), nil)

	id, err := svc.CreatePayment(payment)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestCreatePayment_RepoError(t *testing.T) {
	svc, mockRepo, _ := setupPaymentService(t)
	payment := samplePayment(0)
	repoErr := errors.New("db error")
	ctx := context.Background()

	mockRepo.EXPECT().CreateWithOutbox(ctx, payment).Return(int64(0), repoErr)

	id, err := svc.CreatePayment(payment)

	assert.ErrorIs(t, err, repoErr)
	assert.Equal(t, int64(0), id)
}

func TestUpdatePayment_Success(t *testing.T) {
	svc, mockRepo, _ := setupPaymentService(t)
	payment := samplePayment(1)
	ctx := context.Background()

	mockRepo.EXPECT().Update(ctx, payment).Return(int64(1), nil)

	rows, err := svc.UpdatePayment(payment)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
}

func TestDeletePayment_Success(t *testing.T) {
	svc, mockRepo, _ := setupPaymentService(t)
	ctx := context.Background()

	mockRepo.EXPECT().Delete(ctx, int64(7)).Return(int64(1), nil)

	rows, err := svc.DeletePayment(7)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
}

func TestGetPaymentByPerson_Success(t *testing.T) {
	svc, mockRepo, _ := setupPaymentService(t)
	ctx := context.Background()
	expected := []model.Payment{*samplePayment(1), *samplePayment(2)}

	mockRepo.EXPECT().GetByPerson(ctx, "Alice").Return(expected, nil)

	payments, err := svc.GetPaymentByPerson("Alice")

	assert.NoError(t, err)
	assert.Equal(t, expected, payments)
}

func TestGetPaymentById_CacheHit(t *testing.T) {
	svc, mockRepo, mockCache := setupPaymentService(t)
	payment := samplePayment(5)
	ctx := context.Background()

	data, err := json.Marshal(payment)
	require.NoError(t, err)
	_, err = mockCache.Set(ctx, "payment_id:5", data, 0).Result()
	require.NoError(t, err)

	result, err := svc.GetPaymentById(5)

	assert.NoError(t, err)
	assert.Equal(t, payment.Id, result.Id)
	assert.Equal(t, payment.Person, result.Person)
	assert.Equal(t, payment.Amount, result.Amount)
	mockRepo.AssertNotCalled(t, "GetById", ctx, int64(5))
}

func TestGetPaymentById_CacheMiss(t *testing.T) {
	svc, mockRepo, mockCache := setupPaymentService(t)
	payment := samplePayment(3)
	ctx := context.Background()

	mockRepo.EXPECT().GetById(ctx, int64(3)).Return(payment, nil)

	result, err := svc.GetPaymentById(3)

	assert.NoError(t, err)
	assert.Equal(t, payment, result)

	cached, err := mockCache.Get(ctx, "payment_id:3").Result()
	assert.NoError(t, err)

	var cachedPayment model.Payment
	require.NoError(t, json.Unmarshal([]byte(cached), &cachedPayment))
	assert.Equal(t, payment.Id, cachedPayment.Id)
}

func TestGetPaymentById_RepoError(t *testing.T) {
	svc, mockRepo, _ := setupPaymentService(t)
	repoErr := errors.New("not found")
	ctx := context.Background()

	mockRepo.EXPECT().GetById(ctx, int64(99)).Return(nil, repoErr)

	result, err := svc.GetPaymentById(99)

	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, result)
}
