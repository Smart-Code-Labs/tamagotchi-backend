package component

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// ToyProperties holds all necessary properties for a toy
type ToyProperties struct {
	Name        string
	Description string
	Price       float64
	Wellness    int
}

// ToyKinds map initializes each toy with its properties
var ToyKinds = map[string]ToyProperties{
	"Ball":    {Name: "Ball", Description: "Yuuju!", Price: 5.0, Wellness: 15},
	"Frisbee": {Name: "Frisbee", Description: "Will be back?", Price: 1.0, Wellness: 10},
	"Rope":    {Name: "Rope", Description: "Grrrr", Price: 0.5, Wellness: 10},
	"Stick":   {Name: "Stick", Description: "Throw it! Throw it!", Price: 0.1, Wellness: 5},
}

type ToyStore struct {
	Toys []types.EntityID `json:"toy_list"`
}

func (ToyStore) Name() string {
	return "ToyStore"
}

func (shop *ToyStore) InitToyStore(world cardinal.WorldContext) {
	log := world.Logger()

	for toyName, properties := range ToyKinds {
		// Create an item with the properties from the map
		item := Item{
			ItemName:    toyName,
			Kind:        ItemToy.String(),
			Description: properties.Description,
			Price:       properties.Price,
		}

		// Create the entity with the item and its components
		entityId, err := cardinal.Create(world, item,
			Wellness{Wn: properties.Wellness},
		)

		if err != nil {
			log.Error().Msgf("Failed to create toy %s: %v", toyName, err)
			return
		}

		shop.Toys = append(shop.Toys, entityId)
	}
}
