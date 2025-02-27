// Package query contains functions to query game data.
package query

import (
	"tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// Flow:
// 1. Find the player entity with the given persona tag.
// 2. Retrieve the player's items component.
// 3. Iterate over the items and retrieve each item's component.
// 4. Return a list of item components.
type LeaderboardMsg struct{}

// ItemListReply represents the response to a player items query.
type LeaderboardReply struct {
	// The list of items belonging to the player.
	Pets []component.Pet `json:"pets"`
}

/**
 * QueryPlayerItems queries the items of a player.
 *
 * @param world The game world context.
 * @param req The query request.
 * @return A response containing the list of items, or an error if the query fails.
 */
func QueryLeaderboard(world cardinal.WorldContext, req *LeaderboardMsg) (*LeaderboardReply, error) {
	// Step 1: Initialize the search query to find entities with Pet component.
	log := world.Logger()
	var pets []component.Pet = make([]component.Pet, 0)
	log.Info().Msgf("Received payload to query-leaderboard")
	q := cardinal.NewSearch().Entity(
		filter.Exact(filter.Component[component.Leaderboard]()))

	// Step 4: Iterate over all `Leaderboard` components found and add pets to leaderboard
	q.Each(world, func(id types.EntityID) bool {
		// Step 4.1: Get the `Leaderboard` component
		leaderboard, err := cardinal.GetComponent[component.Leaderboard](world, id)
		if err != nil {
			// Step 4.2: Return true to continue to the next entity if an error occurs
			return true
		}
		for _, pet := range leaderboard.Pets {
			pets = append(pets, *pet)
		}
		return true
	})
	return &LeaderboardReply{Pets: pets}, nil
}
