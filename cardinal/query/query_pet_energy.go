// Package query contains functions to query game data.
package query

import (
	"fmt"
	"tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	"pkg.world.dev/world-engine/cardinal"
)

// Flow:
// 1. Initialize the search query to find entities with Pet and Energy components.
// 2. Iterate over the entities that match the search criteria and check if the Pet component matches the requested nickname.
// 3. If a match is found, retrieve the Energy component and return its value.
// 4. Return an error if no match is found or if the query fails.
type PetEnergyRequest struct {
	// The nickname of the pet to query.
	Nickname string
}

// PetEnergyResponse represents the response to a pet energy query.
type PetEnergyResponse struct {
	// The energy value of the pet.
	E int
}

/**
 * QueryPetEnergy queries the energy of a pet.
 *
 * @param world The game world context.
 * @param req The query request.
 * @return A response containing the energy value of the pet, or an error if the query fails.
 */
func QueryPetEnergy(world cardinal.WorldContext, req *PetEnergyRequest) (*PetEnergyResponse, error) {
	// Step 1: Initialize the search query to find entities with Pet and Energy components.
	var petEnergy *component.Energy
	var err error
	log := world.Logger()
	log.Info().Msgf("Received payload to query-pet-energy: Name[%s]", req.Nickname)

	q := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[component.Pet](), filter.Component[component.Energy]()))

	// Step 2: Iterate over the entities that match the search criteria and check if the Pet component matches the requested nickname.
	// Step 3: If a match is found, retrieve the Energy component and return its value.
	searchErr := q.
		Each(world, func(id types.EntityID) bool {
			var pet *component.Pet
			pet, err = cardinal.GetComponent[component.Pet](world, id)
			if err != nil {
				return false
			}

			// Terminates the search if the pet is found
			if pet.Nickname == req.Nickname {
				petEnergy, err = cardinal.GetComponent[component.Energy](world, id)
				if err != nil {
					return false
				}
				return false
			}

			// Continue searching if the pet is not the target pet
			return true
		})
	if searchErr != nil {
		return nil, searchErr
	}
	if err != nil {
		return nil, err
	}

	// Step 4: Return an error if no match is found or if the query fails.
	if petEnergy == nil {
		return nil, fmt.Errorf("pet %s does not exist", req.Nickname)
	}

	return &PetEnergyResponse{E: petEnergy.E}, nil
}
