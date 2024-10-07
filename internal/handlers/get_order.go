package handlers

import (
	"encoding/json"
	"net/http"

	"L0/internal/cache"
	"L0/internal/storage"
	"github.com/go-chi/chi/v5"
)

func GetOrderByIDHandler(orderCache *cache.Cache, db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "id")

		if orderData, exists := orderCache.Get(orderID); exists {
			w.Header().Set("Content-Type", "application/json")
			w.Write(orderData)
			return
		}

		orderData, err := db.Get(r.Context(), orderID)
		if err != nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		orderCache.Set(orderID, orderData)

		w.Header().Set("Content-Type", "application/json")
		w.Write(orderData)
	}
}

func GetAllOrdersHandler(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := db.GetAll(r.Context())
		if err != nil {
			http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	}
}
