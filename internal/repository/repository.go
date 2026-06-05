package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"subscription-service/internal/model"
)

type Repository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, logger *zap.Logger) *Repository {
	return &Repository{pool: pool, logger: logger}
}

func (r *Repository) Create(ctx context.Context, req model.CreateSubscriptionRequest) (*model.Subscription, error) {
	query := `
  INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING id, created_at, updated_at`

	var sub model.Subscription
	err := r.pool.QueryRow(ctx, query,
		req.ServiceName, req.Price, req.UserID, req.StartDate, req.EndDate,
	).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)

	if err != nil {
		r.logger.Error("Failed to create subscription", zap.Error(err))
		return nil, err
	}

	sub.ServiceName = req.ServiceName
	sub.Price = req.Price
	sub.UserID = req.UserID
	sub.StartDate = req.StartDate
	sub.EndDate = req.EndDate

	return &sub, nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*model.Subscription, error) {
	query := `
  SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
  FROM subscriptions WHERE id = $1`

	var sub model.Subscription
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
		&sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *Repository) Update(ctx context.Context, id int, req model.UpdateSubscriptionRequest) (*model.Subscription, error) {
	setParts := []string{}
	args := []interface{}{}
	counter := 1

	if req.ServiceName != nil {
		setParts = append(setParts, fmt.Sprintf("service_name = $%d", counter))
		args = append(args, *req.ServiceName)
		counter++
	}
	if req.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price = $%d", counter))
		args = append(args, *req.Price)
		counter++
	}
	if req.StartDate != nil {
		setParts = append(setParts, fmt.Sprintf("start_date = $%d", counter))
		args = append(args, *req.StartDate)
		counter++
	}
	if req.EndDate != nil {
		setParts = append(setParts, fmt.Sprintf("end_date = $%d", counter))
		args = append(args, *req.EndDate)
		counter++
	}

	if len(setParts) == 0 {
		return r.GetByID(ctx, id)
	}

	query := fmt.Sprintf(`
  UPDATE subscriptions SET %s
  WHERE id = $%d
  RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`,
		strings.Join(setParts, ", "), counter)

	args = append(args, id)

	var sub model.Subscription
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
		&sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *Repository) List(ctx context.Context, userID, serviceName string, limit, offset int) ([]model.Subscription, int64, error) {
	where := []string{}
	args := []interface{}{}
	counter := 1

	if userID != "" {
		where = append(where, fmt.Sprintf("user_id = $%d", counter))
		args = append(args, userID)
		counter++
	}
	if serviceName != "" {
		where = append(where, fmt.Sprintf("service_name = $%d", counter))
		args = append(args, serviceName)
		counter++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM subscriptions %s", whereClause)
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
  SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
  FROM subscriptions %s
  ORDER BY id DESC
  LIMIT $%d OFFSET $%d`, whereClause, counter, counter+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subscriptions []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
			&sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt)
		if err != nil {
			continue
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, total, nil
}

func (r *Repository) GetTotalCost(ctx context.Context, userID, serviceName, periodFrom, periodTo string) (int, error) {
	query := `
  SELECT COALESCE(SUM(price), 0)
  FROM subscriptions
  WHERE user_id = $1
    AND start_date <= $2
    AND (end_date IS NULL OR end_date >= $3)`

	args := []interface{}{userID, periodTo, periodFrom}

	if serviceName != "" {
		query += " AND service_name = $4"
		args = append(args, serviceName)
	}

	var total int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&total)
	return total, err
}
