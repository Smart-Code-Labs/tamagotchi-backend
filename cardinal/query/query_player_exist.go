// Package query contains functions to query game data.
package query

import (
	"tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
)

// Flow:
// 1. Find the player entity with the given persona tag.
// 2. Retrieve the player's items component.
// 3. Iterate over the items and retrieve each item's component.
// 4. Return a list of item components.
type PlayerExistMsg struct {
	// The persona tag of the player to query.
	PersonaTag string `json:"personaTag"`
}

// ItemListReply represents the response to a player items query.
type PlayerExistReply struct {
	// The list of items belonging to the player.
	exist bool
}

/**
 * QueryPlayerItems queries the items of a player.
 *
 * @param world The game world context.
 * @param req The query request.
 * @return A response containing the list of items, or an error if the query fails.
 */
func QueryPlayerExist(world cardinal.WorldContext, req *PlayerExistMsg) (*PlayerExistReply, error) {
	// Step 1: Find the player entity with the given persona tag.
	log := world.Logger()
	log.Info().Msgf("Received payload to player-exist [%s]", req.PersonaTag)

	playerId, err := component.FindPlayerByPersonaTag(world, req.PersonaTag)
	if err != nil || playerId == 0 {
		log.Info().Msgf("QueryPlayerExist Error [%s]", err.Error())
		return &PlayerExistReply{exist: false}, nil
	}

	log.Info().Msgf("player-exist true.")
	return &PlayerExistReply{exist: true}, nil
}
