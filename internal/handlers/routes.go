package handlers

import (
	"L0/internal/cache"
	"L0/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"path/filepath"
)

func SetupRoutes(orderCache *cache.Cache, db *storage.Storage) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/order", CreateOrderHandler(orderCache, db))
	r.Get("/order/{id}", GetOrderByIDHandler(orderCache, db))
	r.Get("/orders", GetAllOrdersHandler(db))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("static", "index.html"))
	})

	return r
}
