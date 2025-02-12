package transactions

import (
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
	Create(transaction Transaction) error
	GetByUser(userId int) ([]Transaction, error)
}