package query

import (
	comp "tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

type PetsMsg struct{}

type PetsReply struct {
	Pets []comp.Pet `json:"pets"`
}

func Pets(world cardinal.WorldContext, req *PetsMsg) (*PetsReply, error) {
	var err error
	log := world.Logger()
	var pets []comp.Pet
	log.Info().Msgf("Received payload to query-pets")
	q := cardinal.NewSearch().Entity(filter.Contains(filter.Component[comp.Pet]()))
	searchErr := q.
		Each(world, func(id types.EntityID) bool {
			log.Info().Msgf("Checking: Id[%d]", id)
			var pet *comp.Pet
			pet, err = cardinal.GetComponent[comp.Pet](world, id)
			if err != nil {
				return false
			}
			pets = append(pets, *pet)
			return true
		})

	if searchErr != nil {
		return nil, searchErr
	}
	if err != nil {
		return nil, err
	}

	return &PetsReply{Pets: pets}, nil
}
