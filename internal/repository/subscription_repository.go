package repository

import (
	"database/sql"
	"test_task/internal/entity"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func (r *SubscriptionRepository) CreateSubscription(e entity.Subscription) error {
	query := `
		INSERT INTO subscription(service_name, price, user_id, start_date)
		VALUES($1, $2, $3, $4)
		RETURNING id
	`

	var id int

	err := r.db.QueryRow(query, e.ServiceName, e.Price, e.UserId, e.StartDate).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SubscriptionRepository) GetSubscriptionById(id int) (*entity.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date
		FROM subsciption
		WHERE id = $1
	`

	var sub entity.Subscription

	err := r.db.QueryRow(query, id).Scan(
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

func (r *SubscriptionRepository) GetAllSubsctiptions() ([]entity.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date
		FROM subsciption
	`

	rows, err := r.db.Query(query)
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
			continue
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

func (r *SubscriptionRepository) DeleteSubById(id int) error {
	query := `
		DELETE FROM subscription
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SubscriptionRepository) UpdateSubById(e entity.Subscription) error {
	query := `
		UPDATE subscription 
		SET service_name = $1, price = $2, user_id = $3, start_date = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(query, e.ServiceName, e.Price, e.UserId, e.StartDate, e.Id)
	if err != nil {
		return err
	}

	return nil
}
