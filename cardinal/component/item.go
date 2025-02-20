package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

type Item struct {
	ItemName    string  `json:"name"`
	Kind        string  `json:"kind"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (Item) Name() string {
	return "Item"
}

// Items
type ItemKind int

const (
	ItemNone ItemKind = iota // 0
	ItemFood                 // 1
	ItemToy                  // 2
	ItemCare                 // 3
)

func (itemKind ItemKind) String() string { // String method for nice printing
	switch itemKind {
	case ItemNone:
		return "Unknown"
	case ItemFood:
		return "Food"
	case ItemToy:
		return "Toy"
	case ItemCare:
		return "Care"
	default:
		return fmt.Sprintf("ItemKind(%d)", itemKind) // Handle unexpected values
	}
}

// FindItemByName searches for an item with the given name and returns its EntityID
func FindItemByName(world cardinal.WorldContext, itemName string) (types.EntityID, error) {
	q := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[Item]()),
	)

	var itemID types.EntityID
	var found bool = false

	q.Each(world, func(id types.EntityID) bool {
		item, err := cardinal.GetComponent[Item](world, id)
		if err != nil {
			return true
		}

		if item.ItemName == itemName {
			itemID = id
			found = true
			return false // Stop searching once we find the item
		}

		return true
	})

	if !found {
		return 0, fmt.Errorf("item with name %s does not exist", itemName)
	}

	return itemID, nil
}
