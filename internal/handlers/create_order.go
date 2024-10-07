package handlers

import (
	"encoding/json"
	"net/http"

	"L0/internal/cache"
	"L0/internal/model"
	"L0/internal/storage"
)

func CreateOrderHandler(orderCache *cache.Cache, db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order model.OrderInfo
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err := db.Save(r.Context(), order.OrderUID, order.OrderInfo)
		if err != nil {
			http.Error(w, "Failed to save order", http.StatusInternalServerError)
			return
		}

		orderCache.Set(order.OrderUID, order.OrderInfo)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	}
}
