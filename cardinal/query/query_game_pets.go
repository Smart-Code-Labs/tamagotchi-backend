// Package query contains functions to query game data.
package query

import (
	"tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// Flow:
// 1. Initialize the search query to find entities with Pet component.
// 2. Iterate over the entities that match the search criteria and retrieve their Pet component.
// 3. Append the retrieved Pet components to a list.
// 4. Return the list of Pet components.
type PetsMsg struct{}

// PetsReply represents the response to a game pets query.
type PetsReply struct {
	// List of pets in the game.
	Pets []component.Pet `json:"pets"`
}

/**
 * GamePets queries the pets in the game.
 *
 * @param world The game world context.
 * @param req The query request.
 * @return A response containing the list of pets, or an error if the query fails.
 */
func GamePets(world cardinal.WorldContext, req *PetsMsg) (*PetsReply, error) {
	// Step 1: Initialize the search query to find entities with Pet component.
	var err error
	log := world.Logger()
	var pets []component.Pet
	log.Info().Msgf("Received payload to query-pets")
	q := cardinal.NewSearch().Entity(filter.Contains(filter.Component[component.Pet]()))

	// Step 2: Iterate over the entities that match the search criteria and retrieve their Pet component.
	// Step 3: Append the retrieved Pet components to a list.
	searchError := q.
		Each(world, func(id types.EntityID) bool {
			var pet *component.Pet
			pet, err = cardinal.GetComponent[component.Pet](world, id)
			if err != nil {
				return false
			}
			pets = append(pets, *pet)
			return true
		})

	// Step 4: Return the list of Pet components.
	if searchError != nil {
		return nil, searchError
	}
	if err != nil {
		return nil, err
	}

	return &PetsReply{Pets: pets}, nil
}
