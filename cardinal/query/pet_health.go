package query

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
)

type PetHealthRequest struct {
	Nickname string
}

type PetHealthResponse struct {
	HP int
}

func PetHealth(world cardinal.WorldContext, req *PetHealthRequest) (*PetHealthResponse, error) {
	var petHealth *comp.Health
	var err error

	q := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.Pet](), filter.Component[comp.Health]()))

	searchErr := q.
		Each(world, func(id types.EntityID) bool {
			var pet *comp.Pet
			pet, err = cardinal.GetComponent[comp.Pet](world, id)
			if err != nil {
				return false
			}

			// Terminates the search if the pet is found
			if pet.Nickname == req.Nickname {
				petHealth, err = cardinal.GetComponent[comp.Health](world, id)
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

	if petHealth == nil {
		return nil, fmt.Errorf("pet %s does not exist", req.Nickname)
	}

	return &PetHealthResponse{HP: petHealth.HP}, nil
}
