package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/437d5/merch-store/internal/inventory"
	"github.com/437d5/merch-store/internal/user"
)

var	ErrInvalidPassword = errors.New("invalid password")

type UserService struct {
	userRepo user.UserRepo
	logger *slog.Logger
}

func NewUserService(userRepo user.UserRepo, logger *slog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger: logger,
	}
}

func (s *UserService) AuthUser(ctx context.Context, name, password string) (user.User, error) {
	const op = "/internal/service/user_service/AuthUser"

	existingUser, err := s.userRepo.GetUserByName(ctx, name)
	if err == nil {
		if !existingUser.CheckPassword(password) {
			s.logger.Error("Failed to authenticate", "op", op, "errror", ErrInvalidPassword)
			return user.User{}, ErrInvalidPassword
		}

		s.logger.Info("User authenticated succesfully", "op", op, "username", existingUser.Name)
		return existingUser, nil
	}

	newUser := user.User{
		Name: name,
		Coins: 100000,
		Inventory: inventory.Inventory{},
	}

	err = newUser.SetPassword(password)
	if err != nil {
		s.logger.Error("cannot register new user", "op", op, "error", err)
		return user.User{}, fmt.Errorf("cannot set pass: %w", err)
	} 

	id, err := s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		s.logger.Error("Error creating new user", "op", op, "error", err)
	}

	newUser.Id = id
	s.logger.Info("New user authenticated succesfully", "op", op, "username", newUser.Name)
	return newUser, nil
}

func (s *UserService) UserInfo(ctx context.Context, userId int) (user.User, error) {
	const op = "/internal/service/user_service/UserInfo"

	u, err := s.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		s.logger.Error("failed get user", "op", op, "userId", userId)
		return user.User{}, fmt.Errorf("failed get user: %w", err)
	}

	return u, nil
}