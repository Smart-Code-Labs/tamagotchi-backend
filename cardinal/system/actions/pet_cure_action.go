// Package system contains the logic for handling pet feeding actions.
package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	"tamagotchi/component"
	"tamagotchi/game"
	"tamagotchi/msg"
)

/**
 * PetCureAction Cure instantly a Pet.
 *
 * Code Flow:
 * 1. Process each incoming message using `cardinal.EachMessage`.
 * 2. Retrieve the pet's ID by its nickname using `system.QueryPetIdByName`.
 * 3. Check if the pet is not currently engaged in an activity using `CheckPetActivity`.
 * 4. Fetch the pet's health component using `CheckPetHealth`.
 * 5. Fetch the pet's think component using `CheckPetThink`.
 * 6. Increase the pet's health points.
 * 7. Set the pet's think to "Eating".
 * 8. Set the pet's activity to "Eating" and initialize the countdown.
 * 9. Update the pet's components in the world context.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *
 * Returns:
 *   error: Any error that occurs during the process.
 */
func PetCureAction(world cardinal.WorldContext) error {
	log := world.Logger()

	return cardinal.EachMessage(
		world,
		func(eat cardinal.TxData[msg.CurePetMsg]) (msg.CurePetMsgReply, error) {

			// Player sanity check
			playerID, err := component.FindPlayerByPersonaTag(world, eat.Tx.PersonaTag)
			if err != nil {
				return msg.CurePetMsgReply{}, err
			}
			// get player (pets, items, money)
			player, err := cardinal.GetComponent[component.Player](world, playerID)
			if err != nil {
				return msg.CurePetMsgReply{}, fmt.Errorf("failed to bath [get Player]: %w", err)
			}

			// Step 2: Retrieve the pet's ID by its nickname.
			_, petId, err := component.QueryPetIdByName(world, eat.Msg.TargetNickname)
			if err != nil {
				return msg.CurePetMsgReply{}, err
			}

			petHealth, err := CheckPetHealth(world, petId)
			if err != nil {
				return msg.CurePetMsgReply{}, err
			}

			// get `Toy` Health value
			log.Info().Msgf("Cure: Toy [%s]", eat.Msg.ItemName)
			itemId, err := player.GetItemIdByName(world, eat.Msg.ItemName)
			if err != nil {
				return msg.CurePetMsgReply{}, err
			}

			item, err := cardinal.GetComponent[component.Health](world, itemId)
			if err != nil {
				return msg.CurePetMsgReply{}, fmt.Errorf("failed to eat [get Item Hygiene]: %w", err)
			}

			// Step 6: Increase the pet's health points.according to `Toy` item
			if petHealth.HP+item.HP <= 100 {
				petHealth.HP += item.HP
			} else {
				petHealth.HP = 100
			}

			if err := cardinal.SetComponent(world, petId, petHealth); err != nil {
				return msg.CurePetMsgReply{}, fmt.Errorf("failed to Eat [set Health]: %w", err)
			}

			return msg.CurePetMsgReply{
				Health: game.HealthIncrease}, nil
		})
}
