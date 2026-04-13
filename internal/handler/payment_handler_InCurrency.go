package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func (pH *PaymentHandler) GetPaymentInCurrency(w http.ResponseWriter, r *http.Request) {
	pH.logger.Debug("getting payment by id in currency",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path))

	idStr := chi.URLParam(r, "id")
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

	currency := r.URL.Query().Get("currency")
	if currency == "" {
		pH.logger.Error("currency query parameter is missing")
		http.Error(w, "currency query parameter is required", http.StatusBadRequest)
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
