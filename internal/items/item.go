package items

type ItemType struct {
	Name string
	Cost int
}

type ItemRepo interface {
	GetByName(name string) (ItemType, error)
}