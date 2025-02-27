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
type ItemListMsg struct {
	// The persona tag of the player to query.
	PersonaTag string `json:"personaTag"`
}

// ItemListReply represents the response to a player items query.
type ItemListReply struct {
	// The list of items belonging to the player.
	ItemList []component.Item `json:"items"`
}

/**
 * QueryPlayerItems queries the items of a player.
 *
 * @param world The game world context.
 * @param req The query request.
 * @return A response containing the list of items, or an error if the query fails.
 */
func QueryPlayerItems(world cardinal.WorldContext, req *ItemListMsg) (*ItemListReply, error) {
	// Step 1: Find the player entity with the given persona tag.
	var err error
	log := world.Logger()
	log.Info().Msgf("Received payload to query-PerosnaItemList")
	list := make([]component.Item, 0)

	playerID, err := component.FindPlayerByPersonaTag(world, req.PersonaTag)
	if err != nil {
		log.Info().Msgf("QueryPlayerItems Error [%s]", req.PersonaTag)
		return &ItemListReply{ItemList: list}, err
	}

	// Step 2: Retrieve the player's items component.
	player, err := cardinal.GetComponent[component.Player](world, playerID)
	if err != nil {
		return &ItemListReply{ItemList: list}, err
	}

	// Step 3: Iterate over the items and retrieve each item's component.
	items := make([]component.Item, 0)

	if len(player.Items) == 0 {
		log.Info().Msgf("Player has no items")
	}
	if len(player.Pets) == 0 {
		log.Info().Msgf("Player has no pets")
	}
	for _, itemID := range player.Items {
		item, err := cardinal.GetComponent[component.Item](world, itemID)
		if err != nil {
			log.Info().Msgf("QueryPlayerItems Error [%s]", err)
			continue
		}
		items = append(items, *item)
	}
	list = append(list, items...)

	// Step 4: Return a list of item components.
	return &ItemListReply{ItemList: list}, nil
}
