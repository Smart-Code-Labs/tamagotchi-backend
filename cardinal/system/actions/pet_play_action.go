// Package system contains the logic for handling pet play actions.
package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	"tamagotchi/component"
	"tamagotchi/game"
	"tamagotchi/msg"
)

// Function Flow:
// 1. Check if the player exists and is valid.
// 2. Get the player's data, including pets and items.
// 3. Check if the pet is eligible for play (not already doing an activity).
// 4. Get the pet's level and experience.
// 5. Update the pet's experience and level.
// 6. Update the pet's energy, hygiene, and wellness based on play.
// 7. Update the pet's activity and think components.
// 8. Return a reply with the updated pet's status.

/**
 * PetPlayAction handles the pet play action for a given player and pet.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the play action.
 */
func PetPlayAction(world cardinal.WorldContext) error {
	var petEnergy *component.Energy
	var petHygiene *component.Hygiene
	var petWellness *component.Wellness
	var petActivity *component.Activity
	var pet *component.Pet
	var player *component.Player
	var item *component.Wellness

	log := world.Logger()
	return cardinal.EachMessage(
		world,
		func(play cardinal.TxData[msg.PlayPetMsg]) (msg.PlayPetMsgReply, error) {
			// Player sanity check
			playerID, err := component.FindPlayerByPersonaTag(world, play.Tx.PersonaTag)
			if err != nil {
				return msg.PlayPetMsgReply{}, err
			}

			// get player (pets, items, money)
			player, err = cardinal.GetComponent[component.Player](world, playerID)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Player]: %w", err)
			}

			petId, err := player.GetPetNickname(world, play.Msg.TargetNickname)
			if err != nil {
				return msg.PlayPetMsgReply{}, err
			}

			// Pet sanity check
			// check if not activity
			petActivity, err = cardinal.GetComponent[component.Activity](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Activity]: %w", err)
			}

			if petActivity.CountDown > 0 {
				return msg.PlayPetMsgReply{}, fmt.Errorf("pet is already engaged in an activity")
			}

			// get pet lvl
			pet, err = cardinal.GetComponent[component.Pet](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Pet]: %w", err)
			}

			// check if max level
			if pet.Level >= game.MaxLevel {
				return msg.PlayPetMsgReply{}, fmt.Errorf("pet Max lvl, cant grow more")
			}

			// add experience and calculate lvl
			pet.AddXP(game.ExperienceEarn)

			// set activity
			petActivity.Activity = "Playing"
			petActivity.CountDown = game.TickHour
			petActivity.TotalTicks = game.TickHour
			petActivity.Percentage = 100

			// get Energy
			petEnergy, err = cardinal.GetComponent[component.Energy](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Energy]: %w", err)
			}

			// reduce Energy
			if petEnergy.E-game.EnergyReduce > 0 {
				petEnergy.E -= game.EnergyReduce
			} else {
				return msg.PlayPetMsgReply{}, fmt.Errorf("pet energy is insufficient")
			}

			log.Info().Msgf("Playing: reduce pet energy [%d]", petEnergy.E)

			// get Hygiene
			petHygiene, err = cardinal.GetComponent[component.Hygiene](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Hygiene]: %w", err)
			}

			// reduce Hygiene
			if petHygiene.Hy-game.HygieneReduce > 0 {
				petHygiene.Hy -= game.HygieneReduce
			} else {
				petHygiene.Hy = 0
			}

			// get Wellness
			petWellness, err = cardinal.GetComponent[component.Wellness](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Wellness]: %w", err)
			}

			// get `Toy` wellness value
			log.Info().Msgf("Playing: Toy [%s]", play.Msg.ItemName)
			itemId, err := player.GetItemIdByName(world, play.Msg.ItemName)
			if err != nil {
				return msg.PlayPetMsgReply{}, err
			}
			item, err = cardinal.GetComponent[component.Wellness](world, itemId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [%s][%d][get Item Wellness]: %w", play.Msg.ItemName, itemId, err)
			}

			// increase wellness according to `Toy` item
			if petWellness.Wn+item.Wn <= 100 {
				petWellness.Wn += item.Wn
			} else {
				petWellness.Wn = 100
			}

			// update experience
			if err := cardinal.SetComponent(world, petId, pet); err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [set Experience]: %w", err)
			}
			// update energy
			if err := cardinal.SetComponent(world, petId, petEnergy); err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [set Energy]: %w", err)
			}
			// update Hygiene
			if err := cardinal.SetComponent(world, petId, petHygiene); err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [set Hygiene]: %w", err)
			}
			// update wellness
			if err := cardinal.SetComponent(world, petId, petWellness); err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [set Wellness]: %w", err)
			}
			// update activity
			if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [set Activity]: %w", err)
			}

			// consume item
			if err := component.RemoveItem(world, playerID, itemId); err != nil {
				return msg.PlayPetMsgReply{}, err
			}

			log.Info().Msgf("Playing: OK")
			return msg.PlayPetMsgReply{
				Energy:   game.EnergyReduce,
				Hygiene:  game.HygieneReduce,
				Wellness: game.WellnessIncrease,
				Activity: petActivity.Activity,
				Duration: petActivity.CountDown,
			}, nil
		},
	)
}
