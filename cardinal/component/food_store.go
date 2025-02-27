// Package component contains structures and functions for working with game components.
package component

import (
	"tamagotchi/game"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// FoodStore represents a store that holds a list of food items, identified by their entity IDs.
type FoodStore struct {
	Foods []types.EntityID `json:"food_list"`
}

/**
 * Name returns the name of the FoodStore component.
 *
 * Code Flow:
 * 1. Return the string "FoodStore" as the name of the component.
 *
 * Parameters:
 *   None
 *
 * Returns:
 *   (string): The name of the FoodStore component.
 */
// Name returns the name of the FoodStore component
func (FoodStore) Name() string {
	// Step 1: Return the string "FoodStore" as the name of the component
	//         This method is used to identify the component in the game world.
	return "FoodStore"
}

/**
 * InitFoodStore initializes the FoodStore component by creating and adding food items to it.
 *
 * Code Flow:
 * 1. Iterate over the available food kinds and their properties.
 * 2. For each food kind, create a new item with the corresponding properties.
 * 3. Create a new entity for the item with the item and its components (Health and Energy).
 * 4. Add the entity ID of the item to the FoodStore's list of foods.
 * 5. Handle any errors that occur during the creation process.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *
 * Returns:
 *   None
 *
 * Step-by-Step Explanation:
 *   Step 1: Iterate over the available food kinds and their properties. This is done using a range loop to access each food kind and its properties.
 *   Step 2: For each food kind, create a new item with the corresponding properties. This includes setting the item's name, kind, description, and price.
 *   Step 3: Create a new entity for the item with the item and its components. This involves using the cardinal.Create function to create a new entity with the item and its components (Health and Energy).
 *   Step 4: Add the entity ID of the item to the FoodStore's list of foods. This is done by appending the entity ID to the FoodStore's Foods slice.
 *   Step 5: Handle any errors that occur during the creation process. If an error occurs, log the error and return from the function.
 */
// InitFoodStore initializes the FoodStore component
func (shop *FoodStore) InitFoodStore(world cardinal.WorldContext) {
	// Step 1: Iterate over the available food kinds and their properties
	//         Get the logger from the world context to log any errors.
	log := world.Logger()

	// Step 1 (continued): Iterate over the available food kinds and their properties
	for foodName, properties := range game.FoodKinds {
		// Step 2: For each food kind, create a new item with the corresponding properties
		//         Create an item with the properties from the map.
		item := Item{
			ItemName:    foodName,
			Kind:        ItemFood.String(),
			Description: properties.Description,
			Price:       properties.Price,
		}

		// Step 3: Create a new entity for the item with the item and its components
		//         Create the entity with the item and its components (Health and Energy).
		entityId, err := cardinal.Create(world, item,
			Health{HP: properties.Health},
			Energy{E: properties.Energy},
		)

		// Step 5: Handle any errors that occur during the creation process
		//         If an error occurs during the creation process, log the error and return.
		if err != nil {
			log.Error().Msgf("Failed to create item %s: %v", item.ItemName, err)
			return
		}

		// Step 4: Add the entity ID of the item to the FoodStore's list of foods
		//         Append the entity ID of the item to the FoodStore's list of foods.
		shop.Foods = append(shop.Foods, entityId)
	}
}
