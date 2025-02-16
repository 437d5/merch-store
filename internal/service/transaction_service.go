package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/437d5/merch-store/internal/transactions"
	"github.com/437d5/merch-store/internal/user"
)

var (
	ErrNotEnoughCoins = errors.New("not enough coins")
	ErrInvalidAmount  = errors.New("invalid amount of coins")
)

type TransactionService struct {
	transactionRepo transactions.TransactionRepo
	userRepo        user.UserRepo
	logger          *slog.Logger
}

func NewTransactionService(
	transactionRepo transactions.TransactionRepo, userRepo user.UserRepo,
	logger *slog.Logger,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		logger:          logger,
	}
}

func (s *TransactionService) TransferCoins(
	ctx context.Context,
	fromUserId, amount int, toUsername string,
) error {
	const op = "/internal/service/transaction_service/TransferCoins"

	fromUser, err := s.userRepo.GetUserByID(ctx, fromUserId)
	if err != nil {
		s.logger.Error("Error transfering coins", "op", op, "error", err)
		return fmt.Errorf("cannot transfer coins: %w", err)
	}

	toUser, err := s.userRepo.GetUserByName(ctx, toUsername)
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

	err = s.userRepo.UpdateUser(ctx, fromUser)
	if err != nil {
		s.logger.Error("Cannot update 'from' user", "op", op, "error", err)
		return err
	}
	s.logger.Debug("FromUser updated", "op", op)

	err = s.userRepo.UpdateUser(ctx, toUser)
	if err != nil {
		s.logger.Error("Cannot update 'to' user", "op", op, "error", err)
		return err
	}
	s.logger.Debug("ToUser updated", "op", op)

	transaction := transactions.Transaction{
		FromUser: fromUserId,
		ToUser:   toUser.Id,
		Amount:   amount,
	}

	s.logger.Debug("Trying create transaction", "op", op)
	return s.transactionRepo.CreateTransaction(ctx, transaction)
}

func (s *TransactionService) GetTransactionsByUser(
	ctx context.Context, userId int,
) ([]transactions.Transaction, error) {
	const op = "/internal/service/transaction_service/GetTransactionsByUser"

	tList, err := s.transactionRepo.GetTransactionByUser(ctx, userId)
	if err != nil {
		s.logger.Error("failed get transactions", "op", op, "error", err)
		return nil, fmt.Errorf("failed to get transactions")
	}

	return tList, nil
}
