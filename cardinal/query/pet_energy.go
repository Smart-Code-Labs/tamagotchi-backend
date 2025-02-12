package query

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
)

type PetEnergyRequest struct {
	Nickname string
}

type PetEnergyResponse struct {
	E int
}

func PetEnergy(world cardinal.WorldContext, req *PetEnergyRequest) (*PetEnergyResponse, error) {
	var petEnergy *comp.Energy
	var err error
	log := world.Logger()
	log.Info().Msgf("Received payload to query-pet-energy: Name[%s]", req.Nickname)

	q := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.Pet](), filter.Component[comp.Energy]()))
	searchErr := q.
		Each(world, func(id types.EntityID) bool {
			log.Info().Msgf("Checking: Id[%d]", id)
			var pet *comp.Pet
			pet, err = cardinal.GetComponent[comp.Pet](world, id)
			if err != nil {
				return false
			}

			log.Info().Msgf("Checking: Name[%s]", pet.Nickname)
			// Terminates the search if the pet is found
			if pet.Nickname == req.Nickname {
				petEnergy, err = cardinal.GetComponent[comp.Energy](world, id)
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

	if petEnergy == nil {
		return nil, fmt.Errorf("pet %s does not exist", req.Nickname)
	}

	return &PetEnergyResponse{E: petEnergy.E}, nil
}
