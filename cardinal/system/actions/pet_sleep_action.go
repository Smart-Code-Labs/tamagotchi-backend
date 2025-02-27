// Package system contains the logic for handling pet sleep actions.
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
 * PetSleepAction reduces the pet's energy and sets its activity to "Sleeping".
 *
 * Code Flow:
 * 1. Process each incoming message using `cardinal.EachMessage`.
 * 2. Retrieve the pet's ID by its nickname using `QueryPetIdByName`.
 * 3. Check if the pet is not currently engaged in an activity using `CheckPetActivity`.
 * 4. Fetch the pet's think component using `GetPetThink`.
 * 5. Fetch the pet's energy component using `GetPetEnergy`.
 * 6. Increase the pet's energy points.
 * 7. Set the pet's think to "Sleeping".
 * 8. Set the pet's activity to "Sleeping" and initialize the countdown.
 * 9. Update the pet's components in the world context.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *
 * Returns:
 *   error: Any error that occurs during the process.
 */
func PetSleepAction(world cardinal.WorldContext) error {
	var petEnergy *component.Energy
	var petActivity *component.Activity
	var petThink *component.Think

	return cardinal.EachMessage(
		world,
		func(Sleep cardinal.TxData[msg.SleepPetMsg]) (msg.SleepPetMsgReply, error) {
			_, petId, err := component.QueryPetIdByName(world, Sleep.Msg.TargetNickname)
			if err != nil {
				return msg.SleepPetMsgReply{}, err
			}

			if err := system.CheckPetActivity(world, petId); err != nil {
				return msg.SleepPetMsgReply{}, err
			}

			petThink, err = component.GetPetThink(world, petId)
			if err != nil {
				return msg.SleepPetMsgReply{}, err
			}

			petThink.Think = game.ThinkSleep

			// set activity
			petActivity, err = component.GetPetActivity(world, petId)
			if err != nil {
				return msg.SleepPetMsgReply{}, err
			}

			petActivity.Activity = "Sleeping"
			petActivity.CountDown = game.TickEightHours
			petActivity.TotalTicks = game.TickEightHours
			petActivity.Percentage = 100

			petEnergy, err = component.GetPetEnergy(world, petId)
			if err != nil {
				return msg.SleepPetMsgReply{}, err
			}

			// increase Energy
			if petEnergy.E+game.EnergyIncrease <= 100 {
				petEnergy.E += game.EnergyIncrease
			} else {
				petEnergy.E = 100
			}

			if err := UpdatePetComponents(world, petId, petEnergy, petActivity, petThink); err != nil {
				return msg.SleepPetMsgReply{}, err
			}

			return msg.SleepPetMsgReply{
				Energy:   game.EnergyReduce,
				Activity: petActivity.Activity,
				Duration: petActivity.CountDown}, nil
		})
}

/**
 * UpdatePetComponents updates the pet's components in the world context.
 *
 * Code Flow:
 * 1. Update the pet's energy component.
 * 2. Update the pet's activity component.
 * 3. Update the pet's think component.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *   petId (types.EntityID): The ID of the pet.
 *   petEnergy (*component.Energy): The pet's energy component.
 *   petActivity (*component.Activity): The pet's activity component.
 *   petThink (*component.Think): The pet's think component.
 *
 * Returns:
 *   error: Any error that occurs during the process.
 */
func UpdatePetComponents(world cardinal.WorldContext, petId types.EntityID, petEnergy *component.Energy, petActivity *component.Activity, petThink *component.Think) error {
	if err := cardinal.SetComponent(world, petId, petEnergy); err != nil {
		return fmt.Errorf("failed to Sleep [set Energy]: %w", err)
	}

	if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
		return fmt.Errorf("failed to Sleep [set Activity]: %w", err)
	}

	if err := cardinal.SetComponent(world, petId, petThink); err != nil {
		return fmt.Errorf("failed to Sleep [set Think]: %w", err)
	}

	return nil
}
