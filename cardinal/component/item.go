// Package component contains structures and functions for working with game components.
package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

/**
 * Item represents a game item with properties such as name, kind, description, and price.
 *
 * Code Flow:
 *   This struct has no specific code flow as it is a simple data structure.
 *   However, it is used in conjunction with other components and functions to manage items in the game.
 */
type Item struct {
	/**
	 * ItemName is the name of the item.
	 */
	ItemName string `json:"name"`
	/**
	 * Kind is the type or category of the item.
	 */
	Kind string `json:"kind"`
	/**
	 * Description is a brief description of the item.
	 */
	Description string `json:"description"`
	/**
	 * Price is the cost of the item.
	 */
	Price float64 `json:"price"`
}

/**
 * Name returns the name of the Item component.
 *
 * Code Flow:
 * 1. Return the string "Item" as the name of the component.
 *
 * Parameters:
 *   None
 *
 * Returns:
 *   (string): The name of the Item component.
 *
 * Step-by-Step Explanation:
 *   Step 1: This method returns the name of the component as a string. The name is used to identify the component in the game world.
 */
func (Item) Name() string {
	// Step 1: Return the string "Item" as the name of the component
	//         This method is used to identify the component in the game world.
	return "Item"
}

// Items
/**
 * ItemKind represents the type or category of an item.
 *
 * Code Flow:
 *   This enum has no specific code flow as it is a simple data structure.
 *   However, it is used in conjunction with other components and functions to manage items in the game.
 */
type ItemKind int

const (
	ItemNone ItemKind = iota // 0
	ItemFood                 // 1
	ItemToy                  // 2
	ItemCare                 // 3
)

/**
 * String method for ItemKind returns a human-readable string representation of the ItemKind.
 *
 * Code Flow:
 * 1. Check the ItemKind value and return the corresponding string.
 *
 * Parameters:
 *   itemKind (ItemKind): The ItemKind to be converted to a string.
 *
 * Returns:
 *   (string): A human-readable string representation of the ItemKind.
 *
 * Step-by-Step Explanation:
 *   Step 1: This method checks the ItemKind value and returns the corresponding string. If the ItemKind is not recognized, it returns a default string.
 */
func (itemKind ItemKind) String() string { // String method for nice printing
	// Step 1: Check the ItemKind value and return the corresponding string
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

/**
 * FindItemByName searches for an item with the given name and returns its EntityID.
 *
 * Code Flow:
 * 1. Create a new search query that searches for entities with the Item component.
 * 2. Iterate over the entities that match the search query.
 * 3. For each entity, check if the ItemName matches the given name.
 * 4. If a match is found, return the EntityID of the item.
 * 5. If no match is found, return an error.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The game world context.
 *   itemName (string): The name of the item to search for.
 *
 * Returns:
 *   (types.EntityID, error): The EntityID of the item if found, or an error if not found.
 *
 * Step-by-Step Explanation:
 *   Step 1: Create a new search query that searches for entities with the Item component.
 *   Step 2: Iterate over the entities that match the search query.
 *   Step 3: For each entity, check if the ItemName matches the given name.
 *   Step 4: If a match is found, return the EntityID of the item.
 *   Step 5: If no match is found, return an error.
 */
func FindItemByName(world cardinal.WorldContext, itemName string) (types.EntityID, error) {
	// Step 1: Create a new search query that searches for entities with the Item component
	q := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[Item]()),
	)

	var itemID types.EntityID
	var found bool = false

	// Step 2: Iterate over the entities that match the search query
	q.Each(world, func(id types.EntityID) bool {
		// Step 3: For each entity, check if the ItemName matches the given name
		item, err := cardinal.GetComponent[Item](world, id)
		if err != nil {
			return true
		}

		if item.ItemName == itemName {
			// Step 4: If a match is found, return the EntityID of the item
			itemID = id
			found = true
			return false // Stop searching once we find the item
		}

		return true
	})

	// Step 5: If no match is found, return an error
	if !found {
		return 0, fmt.Errorf("Item with name [%s] does not exist", itemName)
	}

	return itemID, nil
}
