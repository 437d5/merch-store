package items

import "context"

type ItemType struct {
	Name string
	Cost int
}

type ItemRepo interface {
	GetItemByName(ctx context.Context, name string) (ItemType, error)
}