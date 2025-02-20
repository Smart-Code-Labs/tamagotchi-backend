package query

import (
	comp "tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
)

type ItemListMsg struct {
	PersonaTag string `json:"personaTag"`
}

type ItemListReply struct {
	ItemList []comp.Item `json:"items"`
}

func PlayerItemList(world cardinal.WorldContext, req *ItemListMsg) (*ItemListReply, error) {
	var err error
	log := world.Logger()
	log.Info().Msgf("Received payload to query-PerosnaItemList")
	list := make([]comp.Item, 0)

	playerID, err := comp.FindPlayerByPersonaTag(world, req.PersonaTag)
	if err != nil {
		log.Info().Msgf("[%s]", req.PersonaTag)
		return &ItemListReply{ItemList: list}, err
	}

	log.Info().Msgf("PlayerId[%d]", playerID)
	player, err := cardinal.GetComponent[comp.Player](world, playerID)
	if err != nil {
		return &ItemListReply{ItemList: list}, err
	}

	log.Info().Msgf("[%f]", player.Money)
	items := make([]comp.Item, 0)

	if len(player.Items) == 0 {
		log.Info().Msgf("Player has no items")
	}
	if len(player.Pets) == 0 {
		log.Info().Msgf("Player has no pets")
	}
	for _, itemID := range player.Items {
		log.Info().Msgf("Looking item[%d]", itemID)
		item, err := cardinal.GetComponent[comp.Item](world, itemID)
		if err != nil {
			log.Info().Msgf("error")
			continue
		}
		log.Info().Msgf("Append [%s]", item.ItemName)
		items = append(items, *item)
	}
	list = append(list, items...)
	return &ItemListReply{ItemList: list}, nil
}
