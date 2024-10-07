package storage

import (
	"L0/internal/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewStorage(ctx context.Context, databaseURL string, logger *slog.Logger) (*Storage, error) {
	dbpool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	logger.Info("Connected to PostgreSQL")
	return &Storage{
		db:  dbpool,
		log: logger,
	}, nil
}

func (s *Storage) Close() {
	s.log.Info("Closing PostgreSQL connection")
	s.db.Close()
}

func (s *Storage) Save(ctx context.Context, orderUID string, orderInfo []byte) error {
	const query = `
		INSERT INTO orders (order_uid, order_info) 
		VALUES ($1, $2) 
		ON CONFLICT (order_uid) DO NOTHING`

	_, err := s.db.Exec(ctx, query, orderUID, orderInfo)
	if err != nil {
		s.log.Error("Failed to save order", slog.String("order_uid", orderUID), slog.String("error", err.Error()))
		return fmt.Errorf("error saving order to the database: %w", err)
	}

	s.log.Info("Order saved successfully", slog.String("order_uid", orderUID))
	return nil
}

func (s *Storage) Get(ctx context.Context, orderUID string) ([]byte, error) {
	const query = "SELECT order_info FROM orders WHERE order_uid = $1"

	var orderInfo []byte
	err := s.db.QueryRow(ctx, query, orderUID).Scan(&orderInfo)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found: %w", err)
		}
		s.log.Error("Failed to retrieve order", slog.String("order_uid", orderUID), slog.String("error", err.Error()))
		return nil, fmt.Errorf("error retrieving order from the database: %w", err)
	}

	s.log.Info("Order retrieved successfully", slog.String("order_uid", orderUID))
	return orderInfo, nil
}

func (s *Storage) GetAll(ctx context.Context) ([]model.OrderInfo, error) {
	const query = "SELECT order_uid, order_info FROM orders"

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		s.log.Error("Failed to retrieve all orders", slog.String("error", err.Error()))
		return nil, fmt.Errorf("error retrieving orders: %w", err)
	}
	defer rows.Close()

	var orders []model.OrderInfo
	for rows.Next() {
		var order model.OrderInfo
		if err := rows.Scan(&order.OrderUID, &order.OrderInfo); err != nil {
			s.log.Error("Error scanning order", slog.String("error", err.Error()))
			return nil, fmt.Errorf("error scanning order: %w", err)
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		s.log.Error("Error during rows iteration", slog.String("error", err.Error()))
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	s.log.Info("All orders retrieved successfully", slog.Int("count", len(orders)))
	return orders, nil
}
