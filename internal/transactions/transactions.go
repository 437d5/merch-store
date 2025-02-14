package transactions

import (
	"context"
	"time"
)

type Transaction struct {
	Id int
	FromUser int
	ToUser int
	Amount int
	Timestamp time.Time
}

type TransactionRepo interface {
	Create(ctx context.Context, transaction Transaction) error
	GetByUser(ctx context.Context, userId int) ([]Transaction, error)
}