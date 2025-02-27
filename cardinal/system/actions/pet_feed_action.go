// Package system contains the logic for handling pet feeding actions.
package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"

	"tamagotchi/component"
	"tamagotchi/game"
	"tamagotchi/msg"
	"tamagotchi/system"
)

/**
 * PetFeedAction reduces the pet's hunger and sets its activity to "Eating".
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
func PetFeedAction(world cardinal.WorldContext) error {
	var petActivity *component.Activity
	log := world.Logger()

	return cardinal.EachMessage(
		world,
		func(eat cardinal.TxData[msg.FeedPetMsg]) (msg.FeedPetMsgReply, error) {

			// Player sanity check
			playerID, err := component.FindPlayerByPersonaTag(world, eat.Tx.PersonaTag)
			if err != nil {
				return msg.FeedPetMsgReply{}, err
			}
			// get player (pets, items, money)
			player, err := cardinal.GetComponent[component.Player](world, playerID)
			if err != nil {
				return msg.FeedPetMsgReply{}, fmt.Errorf("failed to bath [get Player]: %w", err)
			}

			// Step 2: Retrieve the pet's ID by its nickname.
			_, petId, err := component.QueryPetIdByName(world, eat.Msg.TargetNickname)
			if err != nil {
				return msg.FeedPetMsgReply{}, err
			}

			// Step 3: Check if the pet is not currently engaged in an activity.
			if system.CheckPetActivity(world, petId) != nil {
				return msg.FeedPetMsgReply{}, err
			}

			// Step 4: Fetch the pet's health component.
			petHealth, err := CheckPetHealth(world, petId)
			if err != nil {
				return msg.FeedPetMsgReply{}, err
			}

			// get `Toy` Health value
			log.Info().Msgf("Eat: Toy [%s]", eat.Msg.ItemName)
			itemId, err := player.GetItemIdByName(world, eat.Msg.ItemName)
			if err != nil {
				return msg.FeedPetMsgReply{}, err
			}

			item, err := cardinal.GetComponent[component.Health](world, itemId)
			if err != nil {
				return msg.FeedPetMsgReply{}, fmt.Errorf("failed to eat [get Item Hygiene]: %w", err)
			}

			// Step 6: Increase the pet's health points.according to `Toy` item
			if petHealth.HP+item.HP <= 100 {
				petHealth.HP += item.HP
			} else {
				petHealth.HP = 100
			}

			// Step 8: Set the pet's activity to "Eating" and initialize the countdown.
			petActivity.Activity = "Eating"
			petActivity.CountDown = game.TickHour
			petActivity.TotalTicks = game.TickHour
			petActivity.Percentage = 100

			// Step 9: Update the pet's components in the world context.
			if err := cardinal.SetComponent(world, petId, petHealth); err != nil {
				return msg.FeedPetMsgReply{}, fmt.Errorf("failed to Eat [set Health]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
				return msg.FeedPetMsgReply{}, fmt.Errorf("failed to Eat [set Activity]: %w", err)
			}

			return msg.FeedPetMsgReply{
				Health:   game.HealthIncrease,
				Activity: petActivity.Activity,
				Duration: petActivity.CountDown}, nil
		})
}

/**
 * CheckPetHealth retrieves the pet's health component.
 *
 * Code Flow:
 * 1. Fetch the pet's health component.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *   petId (types.EntityID): The ID of the pet.
 *
 * Returns:
 *   (*component.Health, error): The pet's health component, and any error that occurs during the process.
 */
func CheckPetHealth(world cardinal.WorldContext, petId types.EntityID) (*component.Health, error) {
	// Step 1: Fetch the pet's health component.
	petHealth, err := cardinal.GetComponent[component.Health](world, petId)
	if err != nil {
		return nil, fmt.Errorf("failed to Eat [get Health]: %w", err)
	}

	return petHealth, nil
}

/**
 * CheckPetThink retrieves the pet's think component.
 *
 * Code Flow:
 * 1. Fetch the pet's think component.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *   petId (types.EntityID): The ID of the pet.
 *
 * Returns:
 *   (*component.Think, error): The pet's think component, and any error that occurs during the process.
 */
func CheckPetThink(world cardinal.WorldContext, petId types.EntityID) (*component.Think, error) {
	// Step 1: Fetch the pet's think component.
	petThink, err := cardinal.GetComponent[component.Think](world, petId)
	if err != nil {
		return nil, fmt.Errorf("failed to Eat [get Think]: %w", err)
	}

	return petThink, nil
}
