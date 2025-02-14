package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/437d5/merch-store/internal/transactions"
	"github.com/437d5/merch-store/internal/user"
)

var (
	ErrNotEnoughCoins = errors.New("not enough coins")
	ErrInvalidAmount = errors.New("invalid amount of coins")
)

type TransactionService struct {
	transactionRepo transactions.TransactionRepo
	userRepo user.UserRepo
	logger *slog.Logger
}

func NewTransactionService(
	transactionRepo transactions.TransactionRepo, userRepo user.UserRepo,
	logger *slog.Logger,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		userRepo: userRepo,
		logger: logger,
	}
}

func (s *TransactionService) TransferCoins(
	ctx context.Context, 
	fromUserId, toUserId, amount int,
) error {
	const op = "/internal/service/transaction_service/TransferCoins"
	
	fromUser, err := s.userRepo.GetByID(ctx, fromUserId)
	if err != nil {
		s.logger.Error("Error transfering coins", "op", op, "error", err)
		return fmt.Errorf("cannot transfer coins: %w", err)
	}

	toUser, err := s.userRepo.GetByID(ctx, toUserId)
	if err != nil {
		s.logger.Error("Error transfering coins", "op", op, "error", err)
		return fmt.Errorf("cannot transfer coins: %w", err)
	}

	if fromUser.Coins < amount {
		s.logger.Error("Not enough coins to transfer", "op", op, "error", ErrNotEnoughCoins)
		return ErrNotEnoughCoins
	}

	if amount <= 0 {
		s.logger.Error("Error transfer coins", "op", op, "errors", ErrInvalidAmount)
		return ErrInvalidAmount
	}

	fromUser.Coins -= amount
	toUser.Coins += amount

	err = s.userRepo.Update(ctx, fromUser)
	if err != nil {
		s.logger.Error("Cannot update 'from' user", "op", op, "error", err)
		return err
	}
	s.logger.Debug("FromUser updated", "op", op)

	err = s.userRepo.Update(ctx, toUser)
	if err != nil {
		s.logger.Error("Cannot update 'to' user", "op", op, "error", err)
		return err
	}
	s.logger.Debug("ToUser updated", "op", op)

	transaction := transactions.Transaction {
		FromUser: fromUserId,
		ToUser: toUserId,
		Amount: amount,
		Timestamp: time.Now(),
	}

	s.logger.Debug("Trying create transaction", "op", op)
	return s.transactionRepo.Create(ctx, transaction)
}