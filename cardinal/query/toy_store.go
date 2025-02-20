package query

import (
	comp "tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

type ToysMsg struct{}

type ToysReply struct {
	Toys []*comp.Item `json:"toys"`
}

// Toys queries the available toys from the ToyStore component
func Toys(world cardinal.WorldContext, req *ToysMsg) (*ToysReply, error) {
	log := world.Logger()
	toys := make([]*comp.Item, 0)
	log.Info().Msgf("Received payload to query-Foods")

	// Search for the unique ToyStore component
	q := cardinal.NewSearch().Entity(
		filter.Exact(filter.Component[comp.ToyStore]()),
	)

	// Process each entity that matches the search criteria
	searchError := q.Each(world, func(id types.EntityID) bool {
		toyStore, err := cardinal.GetComponent[comp.ToyStore](world, id)
		if err != nil {
			return true
		}

		for _, entityId := range toyStore.Toys {
			// Assuming foodItemId is of type types.EntityID
			toyItem, err := cardinal.GetComponent[comp.Item](world, entityId)
			if err != nil {
				return true
			}
			toys = append(toys, toyItem)
		}

		return true
	})

	if searchError != nil {
		return nil, searchError
	}

	return &ToysReply{Toys: toys}, nil
}
