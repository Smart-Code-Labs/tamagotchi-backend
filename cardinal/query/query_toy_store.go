// Package query contains functions to query game data.
package query

import (
	"tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

/**
 * Flow:
 * 1. Initialize a search query to find the ToyStore component.
 * 2. Iterate over the entities that match the search criteria and retrieve their ToyStore component.
 * 3. Retrieve the list of toys from the ToyStore component.
 * 4. Return a list of toy components.
 *
 * ToysMsg represents a request to query the toy store.
 */
type ToysMsg struct{}

/**
 * Flow:
 * 1. Initialize with the list of toys from the toy store.
 *
 * ToysReply represents the response to a toy store query.
 */
type ToysReply struct {
	// The list of toys in the toy store.
	Toys []*component.Item `json:"toys"`
}

/**
 * Flow:
 * 1. Initialize a search query to find the ToyStore component.
 * 2. Log the query request.
 * 3. Retrieve the list of toys from the toy store.
 * 4. Return a response containing the list of toys.
 *
 * QueryToyStore queries the toys in the toy store.
 *
 * @param world The game world context.
 * @param req The query request.
 * @return A response containing the list of toys, or an error if the query fails.
 */
func QueryToyStore(world cardinal.WorldContext, req *ToysMsg) (*ToysReply, error) {
	// Step 1: Initialize a search query to find the ToyStore component.
	// Step 2: Log the query request.
	log := world.Logger()

	log.Info().Msgf("Received payload to query-Foods")

	// Step 3: Retrieve the list of toys from the toy store.
	toys := GetAllItemsFromToyStore(world)

	// Step 4: Return a list of toy components.
	return &ToysReply{Toys: toys}, nil
}

/**
 * Flow:
 * 1. Initialize a search query to find entities with the ToyStore component.
 * 2. Iterate over the entities that match the search criteria and retrieve their ToyStore component.
 * 3. Retrieve the list of toys from the ToyStore component and append them to the result list.
 * 4. Return the list of toy components.
 *
 * getAllItemsFromToyStore retrieves all items from the toy store.
 *
 * @param world The game world context.
 * @return A list of toy components.
 */
func GetAllItemsFromToyStore(world cardinal.WorldContext) []*component.Item {
	// Step 1: Initialize a search query to find entities with the ToyStore component.
	toys := make([]*component.Item, 0)
	q := cardinal.NewSearch().Entity(
		filter.Exact(filter.Component[component.ToyStore]()),
	)

	// Step 2 and 3: Iterate over the entities that match the search criteria,
	// retrieve their ToyStore component, and retrieve the list of toys.
	// Step 4 is implicit: the function returns the list of toy components.
	q.Each(world, func(id types.EntityID) bool {
		// Retrieve the ToyStore component.
		toyStore, err := cardinal.GetComponent[component.ToyStore](world, id)
		if err != nil {
			return true
		}

		// Retrieve the list of toys from the ToyStore component and append them to the result list.
		for _, entityId := range toyStore.Toys {
			// Assuming foodItemId is of type types.EntityID
			toyItem, err := cardinal.GetComponent[component.Item](world, entityId)
			if err != nil {
				return true
			}
			toys = append(toys, toyItem)
		}

		return true
	})

	return toys
}
