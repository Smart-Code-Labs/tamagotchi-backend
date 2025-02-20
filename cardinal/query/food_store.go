package query

import (
	comp "tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

type FoodsMsg struct{}

type FoodsReply struct {
	Foods []*comp.Item `json:"foods"`
}

// Foods queries the available Foods from the FoodStore component
func Foods(world cardinal.WorldContext, req *FoodsMsg) (*FoodsReply, error) {
	log := world.Logger()
	foods := make([]*comp.Item, 0)
	log.Info().Msgf("Received payload to query-Foods")

	// Search for the unique FoodStore component
	q := cardinal.NewSearch().Entity(
		filter.Exact(filter.Component[comp.FoodStore]()),
	)

	// Process each entity that matches the search criteria
	searchError := q.Each(world, func(id types.EntityID) bool {
		foodStore, err := cardinal.GetComponent[comp.FoodStore](world, id)
		if err != nil {
			return true
		}

		for _, entityId := range foodStore.Foods {
			// Assuming foodItemId is of type types.EntityID
			foodItem, err := cardinal.GetComponent[comp.Item](world, entityId)
			if err != nil {
				return true
			}
			foods = append(foods, foodItem)
		}

		return true
	})

	if searchError != nil {
		return nil, searchError
	}

	return &FoodsReply{Foods: foods}, nil
}
