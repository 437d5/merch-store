package user

import (
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
	GetByID(id int) (User, error)
	GetByName(name string) (User, error)
	Create(user User) error
	Update(user User) error
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