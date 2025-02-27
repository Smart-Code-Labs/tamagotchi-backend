// Package component contains structures and functions for working with game components.
package component

import (
	"tamagotchi/game"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// DrugStore represents a store that holds a list of drugs, identified by their entity IDs.
type DrugStore struct {
	Drugs []types.EntityID `json:"drug_list"`
}

/**
 * Name returns the name of the DrugStore component.
 *
 * Code Flow:
 * 1. Return the string "DrugStore" as the name of the component.
 *
 * Parameters:
 *   None
 *
 * Returns:
 *   (string): The name of the DrugStore component.
 */
// Name returns the name of the DrugStore component
func (DrugStore) Name() string {
	// Step 1: Return the string "DrugStore" as the name of the component
	//         This method is used to identify the component in the game world.
	return "DrugStore"
}

/**
 * InitDrugStore initializes the DrugStore component by creating and adding items to it.
 *
 * Code Flow:
 * 1. Iterate over the available drug kinds and their properties.
 * 2. For each drug kind, create a new item with the corresponding properties.
 * 3. Create a new entity for the item and add it to the DrugStore.
 * 4. Handle any errors that occur during the creation process.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *
 * Returns:
 *   None
 */
// InitDrugStore initializes the DrugStore component
func (shop *DrugStore) InitDrugStore(world cardinal.WorldContext) {
	// Step 1: Iterate over the available drug kinds and their properties
	//         Get the logger from the world context to log any errors.
	log := world.Logger()
	var entityId types.EntityID
	var err error
	for drugName, properties := range game.DrugKinds {

		// Step 2: For each drug kind, create a new item with the corresponding properties
		//         Since we're iterating directly, we can use the current drugName
		selectedDrug := drugName

		// Step 2 (continued): Create a new item with the corresponding properties
		//                    Create a new item with the selected drug name, kind, description, and price.
		item := Item{
			ItemName:    selectedDrug,
			Kind:        ItemCare.String(),
			Description: properties.Description,
			Price:       properties.Price,
		}

		// Step 3: Create a new entity for the item and add it to the DrugStore
		//         Use the Cardinal Create API to create a new entity for the item.
		entityId, err = cardinal.Create(world, item,
			Health{HP: properties.Value},
		)

		// Step 4: Handle any errors that occur during the creation process
		//         If an error occurs during the creation process, log the error and return.
		if err != nil {
			log.Error().Msgf("Failed to create item %s: %v", item.ItemName, err)
			return
		}

		// Step 3 (continued): Add the item to the DrugStore
		//                    Append the entity ID of the item to the DrugStore's list of drugs.
		shop.Drugs = append(shop.Drugs, entityId)
	}

	for bathName, properties := range game.BathKinds {

		// Step 2: For each drug kind, create a new item with the corresponding properties
		//         Since we're iterating directly, we can use the current drugName
		selectedDrug := bathName

		// Step 2 (continued): Create a new item with the corresponding properties
		//                    Create a new item with the selected drug name, kind, description, and price.
		item := Item{
			ItemName:    selectedDrug,
			Kind:        ItemCare.String(),
			Description: properties.Description,
			Price:       properties.Price,
		}

		// Step 3: Create a new entity for the item and add it to the DrugStore
		//         Use the Cardinal Create API to create a new entity for the item.
		entityId, err = cardinal.Create(world, item,
			Hygiene{Hy: properties.Value},
		)

		// Step 4: Handle any errors that occur during the creation process
		//         If an error occurs during the creation process, log the error and return.
		if err != nil {
			log.Error().Msgf("Failed to create item %s: %v", item.ItemName, err)
			return
		}

		// Step 3 (continued): Add the item to the DrugStore
		//                    Append the entity ID of the item to the DrugStore's list of drugs.
		shop.Drugs = append(shop.Drugs, entityId)
	}

}
