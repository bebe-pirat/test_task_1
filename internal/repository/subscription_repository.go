package repository

import (
	"context"
	"database/sql"
	"fmt"
	"test_task/internal/entity"

	"github.com/google/uuid"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, e entity.Subscription) (int, error) {
	query := `
		INSERT INTO subscription(service_name, price, user_id, start_date, end_date)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id
		`

	var id int

	err := r.db.QueryRowContext(ctx, query, e.ServiceName, e.Price, e.UserId, e.StartDate, e.EndDate).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *SubscriptionRepository) GetSubscriptionById(ctx context.Context, id int) (*entity.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, formatted_start_date, formatted_end_date
		FROM subscription
		WHERE id = $1
	`

	var sub entity.Subscription

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sub.Id,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserId,
		&sub.StartDate,
		&sub.EndDate,
	)

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *SubscriptionRepository) GetAllSubscriptions(ctx context.Context) ([]entity.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, formatted_start_date, formatted_end_date
		FROM subscription
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription
	for rows.Next() {
		var sub entity.Subscription
		err := rows.Scan(
			&sub.Id,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserId,
			&sub.StartDate,
			&sub.EndDate,
		)

		if err != nil {
			return nil, err
		}

		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *SubscriptionRepository) DeleteSubById(ctx context.Context, id int) error {
	query := `
		DELETE FROM subscription
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SubscriptionRepository) UpdateSubById(ctx context.Context, e entity.Subscription) error {
	query := `
		UPDATE subscription 
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6
	`

	res, err := r.db.ExecContext(ctx, query, e.ServiceName, e.Price, e.UserId, e.StartDate, e.EndDate, e.Id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *SubscriptionRepository) GetTotalCost(ctx context.Context, userId uuid.UUID, serviceName string, fromDate string, toDate *string) (int, error) {
	query := `
        SELECT COALESCE(SUM(price), 0)
        FROM subscription
        WHERE 1=1
    `
	args := []interface{}{}
	argCounter := 1

	if userId != uuid.Nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCounter)
		args = append(args, userId)
		argCounter++
	}

	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argCounter)
		args = append(args, serviceName)
		argCounter++
	}

	if fromDate != "" {
		query += fmt.Sprintf(" AND start_date >= $%d", argCounter)
		args = append(args, fromDate)
		argCounter++
	}

	if toDate != nil {
		query += fmt.Sprintf(" AND end_date <= $%d", argCounter)
		args = append(args, toDate)
		argCounter++
	}

	var sum int

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&sum)
	if err != nil {
		return 0, err
	}

	return sum, nil
}
