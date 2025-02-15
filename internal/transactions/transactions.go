package transactions

import (
	"context"
)

type Transaction struct {
	FromUser int
	ToUser int
	Amount int
}

type TransactionRepo interface {
	CreateTransaction(ctx context.Context, transaction Transaction) error
	GetTransactionByUser(ctx context.Context, userId int) ([]Transaction, error)
}