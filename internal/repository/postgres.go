package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/437d5/merch-store/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
	logger *slog.Logger
}

func NewPostgresRepo(db *pgxpool.Pool, logger *slog.Logger) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) GetByID(ctx context.Context, id int) (user.User, error) {
	const op = "/internal/repository/postgres/GetByID"
	
	var u user.User
	var inventoryJSON string

	query := `
		SELECT id, name, password, coins, inventory
		FROM users
		WHERE id = $1;
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.Id, &u.Name, &u.Password, &u.Coins, &inventoryJSON,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn("user not found", "op", op, "name", id)
			return user.User{}, fmt.Errorf("user not found: %d", id)
		}
		r.logger.Error("cannot get user", "op", op, "error", err)
		return user.User{}, fmt.Errorf("cannot get user: %w", err)
	}

	err = json.Unmarshal([]byte(inventoryJSON), &u.Inventory.Items)
	if err != nil {
		r.logger.Error(
			"cannot unmarshal inventory", "op", op, "error", err,
		)
		return user.User{}, fmt.Errorf("cannot unmarshal inventory: %w", err)
	}

	return u, nil
}

func (r *PostgresRepo) GetByName(ctx context.Context, name string) (user.User, error) {
	const op = "/internal/repository/postgres/GetByName"
	
	var u user.User
	var inventoryJSON string

	query := `
		SELECT id, name, password, coins, inventory
		FROM users
		WHERE name = $1;
	`	

	err := r.db.QueryRow(ctx, query, name).Scan(
		&u.Id, &u.Name, &u.Password, &u.Coins, &inventoryJSON,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn("user not found", "op", op, "name", name)
			return user.User{}, fmt.Errorf("user not found: %w", err)
		}
		r.logger.Error("cannot get user", "op", op, "error", err)
		return user.User{}, fmt.Errorf("cannot get user: %w", err)
	}

	err = json.Unmarshal([]byte(inventoryJSON), &u.Inventory.Items)
	if err != nil {
		r.logger.Error(
			"cannot unmarshal inventory", "op", op, "error", err,
		)
		return user.User{}, fmt.Errorf("cannot unmarshal inventory: %w", err)
	}

	return u, nil
}

func (r *PostgresRepo) Create(ctx context.Context, user user.User) (int, error) {
	const op = "/internal/repository/postgres/Create"

	inventoryJSON, err := json.Marshal(user.Inventory.Items)
	if err != nil {
		r.logger.Error("cannot marshal inventory", "op", op, "error", err)
		return 0, fmt.Errorf("cannot marshal inventory: %w", err)
	}

	var id int
	query := `
		INSERT INTO users (name, password, coins, inventory)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	err = r.db.QueryRow(
		ctx, query, user.Name, user.Password,
		user.Coins, inventoryJSON,
	).Scan(&id)
	if err != nil {
		r.logger.Error("cannot create user", "op", op, "error", err)
		return 0, fmt.Errorf("cannot create user: %w", err)
	}

	return id, nil
}

func (r *PostgresRepo) Update(ctx context.Context, user user.User) error {
	const op = "/internal/repository/postgres/Update"

	inventoryJSON, err := json.Marshal(user.Inventory.Items)
	if err != nil {
		r.logger.Error("cannot marshal inventory", "op", op, "error", err)
		return fmt.Errorf("cannot marshal inventory: %w", err)
	}

	query := `
		UPDATE users
		SET coins = $1, inventory = $2
		WHERE id = $3
	`

	_, err = r.db.Exec(ctx, query, user.Coins, inventoryJSON, user.Id)
	if err != nil {
		r.logger.Error("cannot update user", "op", op, "error", err)
		return fmt.Errorf("cannot update user: %w", err)
	}

	return nil
}
