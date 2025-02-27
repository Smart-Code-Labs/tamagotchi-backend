// Package query contains functions for querying the state of the Tamagotchi game world.
package query

import (
	"tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// FoodsMsg is a request message for querying the available foods.
type FoodsMsg struct{}

// FoodsReply is a response message containing a list of available foods.
type FoodsReply struct {
	Foods []*component.Item `json:"foods"`
}

/**
 * QueryFoodStore queries the available food store from the FoodStore component.
 *
 * Flow:
 * 1. Search for the unique FoodStore component in the game world.
 * 2. Process each entity that matches the search criteria.
 * 3. Retrieve the food items from the FoodStore component.
 * 4. Return a response message containing the list of food items.
 */
// QueryFoodStore queries the available QueryFoodStore from the FoodStore component
func QueryFoodStore(world cardinal.WorldContext, req *FoodsMsg) (*FoodsReply, error) {
	// Step 1: Search for the unique FoodStore component
	//         Search for the FoodStore component using the Cardinal search API.
	log := world.Logger()
	foods := make([]*component.Item, 0)
	log.Info().Msgf("Received payload to query-Foods")

	// Step 2: Process each entity that matches the search criteria
	//         Use the Cardinal Each API to process each entity that matches the search criteria.
	q := cardinal.NewSearch().Entity(
		filter.Exact(filter.Component[component.FoodStore]()),
	)

	// Step 3: Retrieve the food items from the FoodStore component
	//         Get the food items from the FoodStore component and append them to the response list.
	searchError := q.Each(world, func(id types.EntityID) bool {
		foodStore, err := cardinal.GetComponent[component.FoodStore](world, id)
		if err != nil {
			return true
		}

		for _, entityId := range foodStore.Foods {
			// Assuming foodItemId is of type types.EntityID
			foodItem, err := cardinal.GetComponent[component.Item](world, entityId)
			if err != nil {
				return true
			}
			foods = append(foods, foodItem)
		}

		return true
	})

	// Step 4: Return a response message containing the list of food items
	//         Return a response message containing the list of food items.
	if searchError != nil {
		return nil, searchError
	}

	return &FoodsReply{Foods: foods}, nil
}
