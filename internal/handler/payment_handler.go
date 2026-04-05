package handler

import (
	"GoRestSQL/internal/model"
	"GoRestSQL/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// PaymentHandler обрабатывает HTTP-запросы для платежей
type PaymentHandler struct {
	service service.PaymentService
	logger  *zap.Logger
}

// NewPaymentHandler создаёт новый хендлер с логгером
func NewPaymentHandler(service service.PaymentService, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{
		service: service,
		logger:  logger,
	}
}

func (pH *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	pH.logger.Debug("creating payment", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	var payment model.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		pH.logger.Error("failed to decode payment", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := pH.service.CreatePayment(&payment)
	if err != nil {
		pH.logger.Error("failed to create payment", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pH.logger.Info("payment created", zap.Int64("id", id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payment " + strconv.FormatInt(id, 10)))
}

func (pH *PaymentHandler) GetById(w http.ResponseWriter, r *http.Request) {
	pH.logger.Debug("getting payment by id", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		pH.logger.Error("id query parameter is missing")
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		pH.logger.Error("invalid id parameter", zap.String("id", idStr), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment, err := pH.service.GetPaymentById(int64(id))
	if err != nil {
		pH.logger.Error("failed to get payment", zap.Int64("id", int64(id)), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pH.logger.Info("payment retrieved", zap.Int64("id", int64(id)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

func (pH *PaymentHandler) GetByPerson(w http.ResponseWriter, r *http.Request) {
	pH.logger.Debug("getting payments by person", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	person := r.URL.Query().Get("person")
	if person == "" {
		pH.logger.Error("person query parameter is missing")
		http.Error(w, "person query parameter is required", http.StatusBadRequest)
		return
	}

	payments, err := pH.service.GetPaymentByPerson(person)
	if err != nil {
		pH.logger.Error("failed to get payments by person", zap.String("person", person), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pH.logger.Info("payments retrieved", zap.String("person", person), zap.Int("count", len(payments)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payments)
}

func (pH *PaymentHandler) Update(w http.ResponseWriter, r *http.Request) {
	pH.logger.Debug("updating payment", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	var payment model.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		pH.logger.Error("failed to decode payment", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		pH.logger.Error("id url parameter is missing")
		http.Error(w, "id url parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pH.logger.Error("invalid id parameter", zap.String("id", idStr), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment.Id = int(id)
	updatedID, err := pH.service.UpdatePayment(&payment)
	if err != nil {
		pH.logger.Error("failed to update payment", zap.Int64("id", id), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pH.logger.Info("payment updated", zap.Int64("id", updatedID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedID)
}

func (pH *PaymentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	pH.logger.Debug("deleting payment", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		pH.logger.Error("id url parameter is missing")
		http.Error(w, "id url parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pH.logger.Error("invalid id parameter", zap.String("id", idStr), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deletedID, err := pH.service.DeletePayment(id)
	if err != nil {
		pH.logger.Error("failed to delete payment", zap.Int64("id", id), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pH.logger.Info("payment deleted", zap.Int64("id", deletedID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payment " + strconv.FormatInt(deletedID, 10)))
}
