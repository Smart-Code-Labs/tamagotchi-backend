// Package system contains the logic for handling bath pet actions.
package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	"tamagotchi/component"
	"tamagotchi/game"
	"tamagotchi/msg"
)

/**
 * Function Flow:
 * 1. Get the pet ID by its nickname.
 * 2. Check if the pet is not currently engaged in an activity.
 * 3. Set the pet's think state to "Bath".
 * 4. Set the pet's activity state to "Bathing" with a countdown timer.
 * 5. Get the pet's current hygiene level and increase it.
 * 6. Update the pet's hygiene and activity components in the game state.
 * 7. Return a reply with the updated hygiene, activity, and duration.
 *
 * PetBathAction handles the bath pet action for a given pet.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the bath action.
 */
func PetBathAction(world cardinal.WorldContext) error {
	log := world.Logger()

	return cardinal.EachMessage(
		world,
		func(bath cardinal.TxData[msg.BathPetMsg]) (msg.BathPetMsgReply, error) {

			// Player sanity check
			playerID, err := component.FindPlayerByPersonaTag(world, bath.Tx.PersonaTag)
			if err != nil {
				return msg.BathPetMsgReply{}, err
			}
			// get player (pets, items, money)
			player, err := cardinal.GetComponent[component.Player](world, playerID)
			if err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to bath [get Player]: %w", err)
			}

			// Get the pet ID by its nickname.
			petId, err := player.GetPetNickname(world, bath.Msg.TargetNickname)
			if err != nil {
				return msg.BathPetMsgReply{}, err
			}

			// Pet sanity check
			// check if not activity
			petActivity, err := cardinal.GetComponent[component.Activity](world, petId)
			if err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to bath [get Activity]: %w", err)
			}

			if petActivity.CountDown > 0 {
				return msg.BathPetMsgReply{}, fmt.Errorf("pet is already engaged in an activity")
			}

			// Step 4: Set the pet's activity state to "Bathing" with a countdown timer.
			petActivity.Activity = "Bathing"
			petActivity.CountDown = game.TickHour
			petActivity.TotalTicks = game.TickHour
			petActivity.Percentage = 100

			// Step 5: Get the pet's current hygiene level and increase it.
			petHygiene, err := cardinal.GetComponent[component.Hygiene](world, petId)
			if err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [get Hygiene]: %w", err)
			}

			// get `Toy` Hygiene value
			log.Info().Msgf("Bath: Toy [%s]", bath.Msg.ItemName)
			itemId, err := player.GetItemIdByName(world, bath.Msg.ItemName)
			if err != nil {
				return msg.BathPetMsgReply{}, err
			}

			item, err := cardinal.GetComponent[component.Hygiene](world, itemId)
			if err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to bath [get Item Hygiene]: %w", err)
			}

			// increase hygiene according to `Toy` item
			if petHygiene.Hy+item.Hy <= 100 {
				petHygiene.Hy += item.Hy
			} else {
				petHygiene.Hy = 100
			}

			// get Energy
			petEnergy, err := cardinal.GetComponent[component.Energy](world, petId)
			if err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to play [get Energy]: %w", err)
			}

			// reduce Energy
			if petEnergy.E-game.EnergyReduce > 0 {
				petEnergy.E -= game.EnergyReduce
			} else {
				return msg.BathPetMsgReply{}, fmt.Errorf("pet energy is insufficient")
			}

			// Step 6: Update the pet's hygiene and activity components in the game state.
			if err := cardinal.SetComponent(world, petId, petHygiene); err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [set Hygiene]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [set Activity]: %w", err)
			}

			// update energy
			if err := cardinal.SetComponent(world, petId, petEnergy); err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [set Energy]: %w", err)
			}

			// update activity
			if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [set Activity]: %w", err)
			}

			// consume item
			if err := component.RemoveItem(world, playerID, itemId); err != nil {
				return msg.BathPetMsgReply{}, err
			}

			// Step 7: Return a reply with the updated hygiene, activity, and duration.
			return msg.BathPetMsgReply{
				Hygiene:  game.HygieneIncrease,
				Activity: petActivity.Activity,
				Duration: petActivity.CountDown}, nil
		})
}
