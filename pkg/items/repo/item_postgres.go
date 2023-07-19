package repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/litepubl/test-game/pkg/entity"
	"github.com/litepubl/test-game/pkg/postgres"
)

const defaultEntityCap = 512

// ItemRepo -.
type ItemRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *ItemRepo {
	return &ItemRepo{pg}
}

// List возвращает слайс всех не удаленных записей БД
func (r *ItemRepo) List(ctx context.Context) ([]entity.Item, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("items").
		Where("removed = false").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ItemRepo - List- r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ItemRepo - List- r.Pool.Query: %w", err)
	}

	defer rows.Close()

	entities := make([]entity.Item, 0, defaultEntityCap)

	for rows.Next() {
		e := entity.Item{}

		err = rows.Scan(&e.Id, &e.CampaignId, &e.Name, &e.Description, &e.Priority, &e.Removed, &e.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("ItemRepo - List- rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

// MaxPriority -.
func (r *ItemRepo) MaxPriority(ctx context.Context) (int, error) {
	sql := "Select max(priority) FROM items"
	row := r.Pool.QueryRow(ctx, sql)

	var p *int
	err := row.Scan(&p)
	if err != nil {
		return 0, fmt.Errorf("ItemRepo- MaxPriority- row.Scan: %w", err)
	}

	if p == nil {
		return 0, nil
	}

	return *p, nil
}

// Create -.
func (r *ItemRepo) Create(ctx context.Context, item *entity.Item) error {
	sql, args, err := r.Builder.
		Insert("items").
		Columns("campaign_id, name, description, priority, removed, created_at").
		Values(item.CampaignId, item.Name, item.Description, item.Priority, item.Removed, item.CreatedAt).
		ToSql()
	sql += " RETURNING id"

	if err != nil {
		return fmt.Errorf("ItemRepo - create- r.Builder: %w", err)
	}

	o := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}
	tx, err := r.Pool.BeginTx(ctx, o)
	if err != nil {
		return fmt.Errorf("ItemRepo - create- r.Pool.BeginEx error %w", err)
	}

	defer tx.Rollback(ctx) //nolint

	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(&item.Id)
	if err != nil {
		return fmt.Errorf("ItemRepo - Create- rows.Scan: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("ItemRepo - Create- tx.Commit: %w", err)
	}

	return nil
}

// Update записывает персону вБД
func (r *ItemRepo) Update(ctx context.Context, u entity.UpdateData) (entity.Item, error) {
	e := entity.Item{}
	b := r.Builder.
		Update("Items").
		Set("campaign_id", u.CampaignId).
		Set("name", u.Name).
		Set("priority", u.Priority)

	if u.Description != nil {
		b = b.Set("description", *u.Description)
	}

	updateSQL, args, err := b.
		Where(sq.Eq{"id": u.Id}).
		ToSql()
	if err != nil {
		return e, fmt.Errorf("ItemRepo - Update - r.Builder: %w", err)
	}

	o := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}

	tx, err := r.Pool.BeginTx(ctx, o)
	if err != nil {
		return e, fmt.Errorf("ItemRepo - Update- r.Pool.BeginEx error %w", err)
	}

	defer tx.Rollback(ctx) //nolint

	err = r.selectItem(ctx, tx, u.Id, &e)
	if err != nil {
		return e, fmt.Errorf("ItemRepo - Update - row.Scan: %w", err)
	}

	rows, err := tx.Query(ctx, updateSQL, args...)
	if err != nil {
		return e, fmt.Errorf("ItemRepo - Update- tx.QueryRow: %w", err)
	}

	rows.Close()
	err = tx.Commit(ctx)
	if err != nil {
		return e, fmt.Errorf("ItemRepo - Update - tx.Commit: %w", err)
	}

	e.Name = u.Name
	e.CampaignId = u.CampaignId
	e.Priority = u.Priority
	if u.Description != nil {
		e.Description = *u.Description
	}

	return e, nil
}

// Delete помечает как удаленный
func (r *ItemRepo) Delete(ctx context.Context, id int) (entity.Item, error) {
	item := entity.Item{}
	p, err := r.MaxPriority(ctx)
	if err != nil {
		return item, fmt.Errorf("itemRepo.Delete error get MaxPriority %w", err)
	}

	sql, args, err := r.Builder.
		Update("Items").
		Set("removed", true).
		Set("priority", p+1).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return item, fmt.Errorf("ItemRepo - Update - r.Builder: %w", err)
	}

	o := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}

	tx, err := r.Pool.BeginTx(ctx, o)
	if err != nil {
		return item, fmt.Errorf("ItemRepo - Update- r.Pool.BeginEx error %w", err)
	}

	defer tx.Rollback(ctx) //nolint

	err = r.selectItem(ctx, tx, id, &item)
	if err != nil {
		return item, fmt.Errorf("ItemRepo- Delete - tx.select for update: %w", err)
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return item, fmt.Errorf("ItemRepo - Delete - tx.QueryRow: %w", err)
	}

	rows.Close()
	err = tx.Commit(ctx)
	if err != nil {
		return item, fmt.Errorf("ItemRepo - Delete - tx.Commit: %w", err)
	}

	item.Priority = p + 1
	item.Removed = true

	return item, nil
}

func (r *ItemRepo) selectItem(ctx context.Context, tx pgx.Tx, id int, e *entity.Item) error {
	sql, args, err := r.Builder.
		Select("*").
		From("items").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ItemRepo - selectItem - r.Builder: %w", err)
	}

	sql += " FOR UPDATE"
	row := tx.QueryRow(ctx, sql, args...)

	return row.Scan(&e.Id, &e.CampaignId, &e.Name, &e.Description, &e.Priority, &e.Removed, &e.CreatedAt)
}
