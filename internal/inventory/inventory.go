package inventory

type Item struct {
	ItemType string
	Quantity int
}

type Inventory struct {
	Items []Item
}

func (i *Inventory) AddItem(item Item) {
	for idx, savedItem := range i.Items {
		if savedItem.ItemType == item.ItemType {
			i.Items[idx].Quantity += item.Quantity
			return
		}
	}

	i.Items = append(i.Items, item)
}
