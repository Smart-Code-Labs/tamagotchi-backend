package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	comp "tamagotchi/component"
	"tamagotchi/msg"
)

// TODO: Here we need to update the differents stores

// ItemAction reduce pet's Hygiene.
// ItemAction reduce pet's Hygiene.
func BuyItemAction(world cardinal.WorldContext) error {
	log := world.Logger()
	log.Info().Msgf("Called BuyItemAction")
	return cardinal.EachMessage(
		world,
		func(buyItem cardinal.TxData[msg.ButItemMsg]) (msg.BuyItemMsgReply, error) {
			log.Info().Msgf("buyItem: n[%s]", buyItem.Msg.Name)

			// Check if player exists
			playerID, err := comp.FindPlayerByPersonaTag(world, buyItem.Tx.PersonaTag)
			if err != nil {
				return msg.BuyItemMsgReply{}, err
			}

			// Check if item exist
			itemId, err := comp.FindItemByName(world, buyItem.Msg.Name)
			if err != nil {
				return msg.BuyItemMsgReply{}, err
			}

			log.Info().Msgf("Found: Player[%d] item[%d]", playerID, itemId)

			player, err := cardinal.GetComponent[comp.Player](world, playerID)
			if err != nil {
				return msg.BuyItemMsgReply{}, err
			}
			item, err := cardinal.GetComponent[comp.Item](world, itemId)
			if err != nil {
				return msg.BuyItemMsgReply{}, err
			}

			// Chek if user has enough funds
			if player.Money < item.Price {
				return msg.BuyItemMsgReply{}, fmt.Errorf("error Buying,no enought balance : %w", err)
			}

			// Buy item
			comp.AddItem(world, playerID, itemId)
			err = comp.ReducePlayerMoney(world, playerID, item.Price)
			if err != nil {
				return msg.BuyItemMsgReply{}, err
			}
			return msg.BuyItemMsgReply{
				Success: true}, nil
		})
}
