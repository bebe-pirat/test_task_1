package repository

import (
	"context"
	"database/sql"
	"test_task/internal/entity"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, e entity.Subscription) error {
	query := `
		INSERT INTO subscription(service_name, price, user_id, start_date)
		VALUES($1, $2, $3, $4)
		RETURNING id
		`

	var id int

	err := r.db.QueryRowContext(ctx, query, e.ServiceName, e.Price, e.UserId, e.StartDate).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SubscriptionRepository) GetSubscriptionById(ctx context.Context, id int) (*entity.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date
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
	)

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *SubscriptionRepository) GetAllSubsctiptions(ctx context.Context) ([]entity.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date
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
			&sub.StartDate)

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
		SET service_name = $1, price = $2, user_id = $3, start_date = $4
		WHERE id = $5
	`

	res, err := r.db.ExecContext(ctx, query, e.ServiceName, e.Price, e.UserId, e.StartDate, e.Id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SubscriptionRepository) GetTotalCost(ctx context.Context) (int, error) {
	query := `
		SELECT SUM(price)
		FROM subscription
	`

	var sum int

	err := r.db.QueryRowContext(ctx, query).Scan(
		&sum,
	)

	if err != nil {
		return -1, err
	}

	return sum, nil
}
