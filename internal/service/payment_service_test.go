package service

import (
	"GoRestSQL/internal/model"
	repomocks "GoRestSQL/internal/repository/mocks"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupPaymentService(t *testing.T) (*PaymentServiceImpl, *repomocks.MockPaymentRepository, *miniredis.Miniredis) {
	t.Helper()

	mockRepo := repomocks.NewMockPaymentRepository(t)
	mr := miniredis.RunT(t)
	cache := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() {
		_ = cache.Close()
	})

	svc := NewPaymentServiceImpl(context.Background(), mockRepo, nil, zap.NewNop().Sugar(), cache)
	return svc, mockRepo, mr
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
	svc, mockRepo, mr := setupPaymentService(t)
	payment := samplePayment(5)
	ctx := context.Background()

	data, err := json.Marshal(payment)
	require.NoError(t, err)
	mr.Set("payment_id:5", string(data))

	result, err := svc.GetPaymentById(5)

	assert.NoError(t, err)
	assert.Equal(t, payment.Id, result.Id)
	assert.Equal(t, payment.Person, result.Person)
	assert.Equal(t, payment.Amount, result.Amount)
	mockRepo.AssertNotCalled(t, "GetById", ctx, int64(5))
}

func TestGetPaymentById_CacheMiss(t *testing.T) {
	svc, mockRepo, mr := setupPaymentService(t)
	payment := samplePayment(3)
	ctx := context.Background()

	mockRepo.EXPECT().GetById(ctx, int64(3)).Return(payment, nil)

	result, err := svc.GetPaymentById(3)

	assert.NoError(t, err)
	assert.Equal(t, payment, result)

	cached, err := mr.Get("payment_id:3")
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
