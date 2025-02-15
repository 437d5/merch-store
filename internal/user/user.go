package user

import (
	"context"
	"fmt"

	"github.com/437d5/merch-store/internal/inventory"
	hash "github.com/437d5/merch-store/pkg/password"
)

type User struct {
	Id int
	Name string
	Password string
	Coins int
	Inventory inventory.Inventory
}

type UserRepo interface {
	GetUserByID(ctx context.Context, id int) (User, error)
	GetUserByName(ctx context.Context, name string) (User, error)
	CreateUser(ctx context.Context, user User) (int, error)
	UpdateUser(ctx context.Context, user User) error
}

func (u *User) SetPassword(password string) error {
	res, err := hash.HashPassword(password)
	if err != nil {
		return fmt.Errorf("cannot generate password")
	}

	u.Password = res
	return nil
}

func (u *User) CheckPassword(password string) bool {
	ok := hash.VerifyPassword(password, u.Password)
	return ok
}