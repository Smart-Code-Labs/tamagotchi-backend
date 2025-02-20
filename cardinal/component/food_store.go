package component

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// Define a struct to hold food properties including description
type FoodProperties struct {
	Price       float64
	Health      int
	Energy      int
	Description string
}

// Initialize the FoodKinds map with FoodProperties including descriptions
var FoodKinds = map[string]FoodProperties{
	"Apple":   {Price: 0.1, Health: 10, Energy: 10, Description: "Yuumy Red Food"},
	"Banana":  {Price: 0.3, Health: 5, Energy: 15, Description: "What is this Yellow Food?"},
	"Soup":    {Price: 0.5, Health: 15, Energy: 20, Description: "Spicy!!!"},
	"Carrots": {Price: 0.1, Health: 5, Energy: 25, Description: "Cheap, but powerful"},
}

type FoodStore struct {
	Foods []types.EntityID `json:"food_list"`
}

func (FoodStore) Name() string {
	return "FoodStore"
}

func (shop *FoodStore) InitFoodStore(world cardinal.WorldContext) {
	log := world.Logger()

	for foodName, properties := range FoodKinds {
		// Create an item with the properties from the map
		item := Item{
			ItemName:    foodName,
			Kind:        ItemFood.String(),
			Description: properties.Description,
			Price:       properties.Price,
		}

		// Create the entity with the item and its components
		entityId, err := cardinal.Create(world, item,
			Health{HP: properties.Health},
			Energy{E: properties.Energy},
		)

		if err != nil {
			log.Error().Msgf("Failed to create item %s: %v", item.ItemName, err)
			return
		}

		shop.Foods = append(shop.Foods, entityId)
	}
}
