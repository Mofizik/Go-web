package repository

import (
	"context"
	"fmt"
	"order/internal/order/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Save(ctx context.Context, order *model.Order) (*model.Order, error) {
	query, args, err := sq.Insert("orders").
		Columns("id", "item", "quantity").
		Values(order.ID, order.Item, order.Quantity).
		PlaceholderFormat(sq.Dollar).
		ToSql()
		if err != nil {
			return nil, fmt.Errorf("repository save %w", err)
		}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("repository save %w", err)
	}
	return order, nil

}

func (r *Repository) FindByID(ctx context.Context, id string) (*model.Order, error) {
	query, args, err := sq.Select("id", "item", "quantity").
		From("orders").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("repository find by id %w", err)
	}

	row := r.pool.QueryRow(ctx, query, args...)
	order := &model.Order{}
	err = row.Scan(&order.ID, &order.Item, &order.Quantity)
	if err != nil {
		return nil, fmt.Errorf("repository find by id %w", err)
	}
	return order, nil
}

func (r *Repository) Update(ctx context.Context, order *model.Order) (*model.Order, error) {
	query, args, err := sq.Update("orders").
		Set("item", order.Item).
		Set("quantity", order.Quantity).
		Where(sq.Eq{"id": order.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("repository update %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("repository update %w", err)
	}
	return order, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query, args, err := sq.Delete("orders").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("repository delete %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("repository delete %w", err)
	}
	return nil
}

func (r *Repository) List(ctx context.Context) ([]*model.Order, error) {
	query, args, err := sq.Select("id", "item", "quantity").
		From("orders").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("repository list %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("repository list %w", err)
	}
	defer rows.Close()
	orders := make([]*model.Order, 0)
	for rows.Next() {
		order := &model.Order{}
		err = rows.Scan(&order.ID, &order.Item, &order.Quantity)
		if err != nil {
			return nil, fmt.Errorf("repository list %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}
