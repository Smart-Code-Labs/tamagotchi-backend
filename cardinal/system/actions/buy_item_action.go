// Package system contains the logic for handling item buying actions.
package system

import (
	"pkg.world.dev/world-engine/cardinal"

	"tamagotchi/component"
	"tamagotchi/msg"
)

/**
 * Function Flow:
 * 1. Check if the player exists and is valid.
 * 2. Check if the item to be bought exists.
 * 3. Get the player's and item's data.
 * 4. Check if the player has enough balance to buy the item.
 * 5. Reduce the player's balance by the item's price.
 * 6. Add the item to the player's inventory.
 * 7. Return a reply indicating the success of the buy action.
 *
 * BuyItemAction handles the item buying action for a given player and item.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the buy action.
 */
func BuyItemAction(world cardinal.WorldContext) error {
	log := world.Logger()
	return cardinal.EachMessage(
		world,
		func(buyItem cardinal.TxData[msg.ButItemMsg]) (msg.BuyItemMsgReply, error) {
			log.Info().Msgf("buyItem: n[%s]", buyItem.Msg.Name)

			// Player sanity check
			playerID, err := component.FindPlayerByPersonaTag(world, buyItem.Tx.PersonaTag)
			if err != nil {
				return msg.BuyItemMsgReply{}, err
			}

			// Item sanity check
			itemId, err := component.FindItemByName(world, buyItem.Msg.Name)
			if err != nil {
				return msg.BuyItemMsgReply{}, err
			}

			// Get item
			item, err := cardinal.GetComponent[component.Item](world, itemId)
			if err != nil {
				return msg.BuyItemMsgReply{}, err
			}

			// Reduce player's balance
			err = component.ReducePlayerMoney(world, playerID, item.Price)
			if err != nil {
				return msg.BuyItemMsgReply{}, err
			}

			// Buy item
			err = component.AddPlayerItem(world, playerID, itemId)
			if err != nil {
				return msg.BuyItemMsgReply{}, err
			}

			return msg.BuyItemMsgReply{
				Success: true,
			}, nil
		},
	)
}
