// Package component contains various components for the Tamagotchi game.
package component

import (
	"tamagotchi/game"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// ToyStore represents a toy store component that contains a list of toy entities.
type ToyStore struct {
	// Toys is a list of entity IDs representing the toys in the store.
	Toys []types.EntityID `json:"toy_list"`
}

/**
 * Name returns the name of the component.
 *
 * Code Flow:
 * 1. The function simply returns the string "ToyStore" as the name of the component.
 *
 * Returns:
 *   string: The name of the component.
 */
func (ToyStore) Name() string {
	// Step 1: Return the string "ToyStore" as the name of the component.
	return "ToyStore"
}

/**
 * InitToyStore initializes the toy store by creating toy entities based on the game's toy kinds.
 *
 * Code Flow:
 * 1. Fetch the logger from the world context.
 * 2. Iterate over the game's toy kinds and create an item for each kind.
 * 3. Create an entity for each item with the item and wellness components.
 * 4. Append the entity ID to the toy store's list of toys.
 * 5. Handle any errors that occur during the creation process.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 */
func (shop *ToyStore) InitToyStore(world cardinal.WorldContext) {
	// Step 1: Fetch the logger from the world context.
	log := world.Logger()

	// Step 2: Iterate over the game's toy kinds and create an item for each kind.
	for toyName, properties := range game.ToyKinds {
		// Create an item with the properties from the map
		item := Item{
			ItemName:    toyName,
			Kind:        ItemToy.String(),
			Description: properties.Description,
			Price:       properties.Price,
		}

		// Step 3: Create the entity with the item and its components
		entityId, err := cardinal.Create(world, item,
			Wellness{Wn: properties.Wellness},
		)

		// Step 5: Handle any errors that occur during the creation process.
		if err != nil {
			// Log the error and return.
			log.Error().Msgf("Failed to create toy %s: %v", toyName, err)
			return
		}

		// Step 4: Append the entity ID to the toy store's list of toys.
		shop.Toys = append(shop.Toys, entityId)
	}
}
