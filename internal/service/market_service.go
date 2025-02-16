package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/437d5/merch-store/internal/inventory"
	"github.com/437d5/merch-store/internal/items"
	"github.com/437d5/merch-store/internal/user"
)

type MarketService struct {
	userRepo user.UserRepo
	itemRepo items.ItemRepo
	logger   *slog.Logger
}

func NewMarketService(
	userRepo user.UserRepo, logger *slog.Logger, itemRepo items.ItemRepo,
) *MarketService {
	return &MarketService{
		userRepo: userRepo,
		itemRepo: itemRepo,
		logger:   logger,
	}
}

func (s *MarketService) BuyMerch(ctx context.Context, userId int, itemType string) error {
	const op = "/internal/service/market_service/BuyMerch"

	user, err := s.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		s.logger.Error("cannot find user", "op", op, "error", err)
		return fmt.Errorf("cannot find user: %w", err)
	}

	itemCard, err := s.itemRepo.GetItemByName(ctx, itemType)
	if err != nil {
		s.logger.Error("cannot find item", "op", op, "error", err)
		return fmt.Errorf("cannot find item: %w", err)
	}

	if user.Coins < itemCard.Cost {
		s.logger.Error("cannot buy item", "op", op, "error", ErrNotEnoughCoins)
		return fmt.Errorf("cannot buy item: %w", ErrNotEnoughCoins)
	}

	user.Coins = user.Coins - itemCard.Cost
	user.Inventory.AddItem(inventory.Item{
		ItemType: itemType,
		Quantity: 1,
	})

	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		s.logger.Error("cannot update user", "op", op, "error", err)
		return fmt.Errorf("cannot update user: %w", err)
	}

	return nil
}
