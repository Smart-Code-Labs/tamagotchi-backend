package query

import (
	comp "tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

type DrugsMsg struct{}

type DrugsReply struct {
	Drugs []comp.Item `json:"drugs"`
}

// Drugs queries the available Drugs from the DrugStore component
func Drugs(world cardinal.WorldContext, req *DrugsMsg) (*DrugsReply, error) {
	log := world.Logger()
	drugs := make([]comp.Item, 0)
	log.Info().Msgf("Received payload to query-Drugs")

	// Search for the unique DrugStore component
	q := cardinal.NewSearch().Entity(
		filter.Exact(filter.Component[comp.DrugStore]()),
	)

	// Process each entity that matches the search criteria
	searchError := q.Each(world, func(id types.EntityID) bool {
		drugStore, err := cardinal.GetComponent[comp.DrugStore](world, id)
		if err != nil {
			return true
		}
		for _, entityId := range drugStore.Drugs {
			drugItem, err := cardinal.GetComponent[comp.Item](world, entityId)
			if err != nil {
				return true
			}
			drugs = append(drugs, *drugItem)
		}
		return true
	})
	if searchError != nil {
		return nil, searchError
	}

	return &DrugsReply{Drugs: drugs}, nil
}
