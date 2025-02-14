package service

import (
	"fmt"
	"log/slog"

	"github.com/437d5/merch-store/internal/inventory"
	"github.com/437d5/merch-store/internal/items"
	"github.com/437d5/merch-store/internal/user"
)

type MarketService struct {
	userRepo user.UserRepo
	itemRepo items.ItemRepo
	logger *slog.Logger
}

func NewMarketService(
	userRepo user.UserRepo, logger *slog.Logger, itemRepo items.ItemRepo,
) *MarketService {
	return &MarketService{
		userRepo: userRepo,
		itemRepo: itemRepo,
		logger: logger,
	}
}

func (s *MarketService) BuyMerch(userId int, item inventory.Item) error {
	const op = "/internal/service/market_service/BuyMerch"

	user, err := s.userRepo.GetByID(userId)
	if err != nil {
		s.logger.Error("cannot find user", "op", op, "error", err)
		return fmt.Errorf("cannot find user: %s", err)
	}

	itemCard, err := s.itemRepo.GetByName(item.ItemType)
	if err != nil {
		s.logger.Error("cannot find item", "op", op, "error", err)
		return fmt.Errorf("cannot find item: %s", err)
	}

	if user.Coins < itemCard.Cost {
		s.logger.Error("cannot buy item", "op", op, "error", ErrNotEnoughCoins)
		return fmt.Errorf("cannot buy item: %s", ErrNotEnoughCoins)
	}

	user.Coins = user.Coins - itemCard.Cost
	user.Inventory.AddItem(item)

	err = s.userRepo.Update(user)
	if err != nil {
		s.logger.Error("cannot update user", "op", op, "error", err)
		return fmt.Errorf("cannot update user: %s", err)
	}

	return nil
}