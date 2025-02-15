package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/437d5/merch-store/internal/items"
	"github.com/437d5/merch-store/internal/transactions"
	"github.com/437d5/merch-store/internal/user"
	"github.com/jackc/pgx/v5"
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

// UserRepo implementation
func (r *PostgresRepo) GetUserByID(ctx context.Context, id int) (user.User, error) {
	const op = "/internal/repository/postgres/GetUserByID"
	
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
		if errors.Is(err, pgx.ErrNoRows) {
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

func (r *PostgresRepo) GetUserByName(ctx context.Context, name string) (user.User, error) {
	const op = "/internal/repository/postgres/GetUserByName"
	
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

func (r *PostgresRepo) CreateUser(ctx context.Context, user user.User) (int, error) {
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

func (r *PostgresRepo) UpdateUser(ctx context.Context, user user.User) error {
	const op = "/internal/repository/postgres/UpdateUser"

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

// TransactionRepo implementation
func (r *PostgresRepo) CreateTransaction(ctx context.Context, t transactions.Transaction) error {
	const op = "/internal/repository/postgres/CreateTransaction"

	query := `
		INSERT INTO transactions (from_user, to_user, amount)
		VALUES ($1, $2, $3);
	`

	_, err := r.db.Exec(ctx, query, t.FromUser, t.ToUser, t.Amount)
	if err != nil {
		r.logger.Error("cannot create transaction", "op", op, "error", err)
		return fmt.Errorf("cannot create transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepo) GetTransactionByUser(ctx context.Context, userId int) ([]transactions.Transaction, error) {
	const op = "/internal/repository/postgres/GetTransactionByUser"

	query := `
		SELECT from_user, to_user, amount
		FROM transactions
		WHERE from_user = $1 OR to_user = $1
		ORDER BY timestamp DESC;
	`

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn("no transactions found", "op", op, "userID", userId)
			return nil, fmt.Errorf("no transaction found: %w", err)
		}

		r.logger.Error("failed to get transactions", "op", op, "error", err)
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	var tList []transactions.Transaction

	for rows.Next() {
		var t transactions.Transaction
		err := rows.Scan(
			&t.FromUser, &t.ToUser, &t.Amount,
		)
		if err != nil {
			r.logger.Error("failed to scan transaction", "op", op, "error", err)
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		tList = append(tList, t)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("rows iteration error", "op", op, "error", err)
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tList, nil
}

// ItemRepo implementation
func (r *PostgresRepo) GetItemByName(ctx context.Context, name string) (items.ItemType, error) {
	const op = "/internal/repository/postgres/GetItemByName"

	var item items.ItemType
	item.Name = name

	query := `
		SELECT cost FROM items
		WHERE name = $1;
	`

	err := r.db.QueryRow(ctx, query, name).Scan(&item.Cost)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn("item not found", "op", op, "name", name)
			return items.ItemType{}, fmt.Errorf("item not found: %w", err)
		}

		r.logger.Error("failed to get item", "op", op, "error", err)
		return items.ItemType{}, fmt.Errorf("failed to get item: %w", err)
	}

	return item, nil
}