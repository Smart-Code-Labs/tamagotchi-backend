/*
Package system provides core game mechanics for the Tamagotchi backend.

This package contains various system-level functionalities that manage and update the state of game entities.
It includes systems for handling activities, declining pet statistics over time, and other essential game mechanics.
*/
package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	"tamagotchi/component"
	"tamagotchi/game"
)

/**
 * Function Flow:
 * 1. The `ActivityDeclineSystem` function is called, which checks if the current tick is a multiple of `game.ActivityUpdateTickRate`.
 * 2. If it is, the function queries all entities that have both `Pet` and `Activity` components.
 * 3. For each entity found, the function retrieves the `Activity` component and checks if the activity is not "None".
 * 4. If the activity is not "None", the function decrements the activity duration by one.
 * 5. If the activity duration is greater than zero, the function updates the activity percentage and sets the updated `Activity` component.
 * 6. If the activity duration reaches zero, the function resets the activity to "None" and sets the updated `Activity` component.
 *
 * ActivityDeclineSystem periodically decreases the duration of a pet's current activity.
 *
 * This system iterates over all entities that have both Pet and Activity components,
 * reducing the activity duration by one each time it is processed. If the duration reaches zero,
 * the activity is effectively completed.
 *
 * The function returns an error if there is a failure during component access or update.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the execution of the activity decline system.
 */
func ActivityDeclineSystem(world cardinal.WorldContext) error {
	log := world.Logger()
	// Step 1: Check if the current tick is a multiple of `game.ActivityUpdateTickRate`
	if world.CurrentTick()%game.ActivityUpdateTickRate == 0 {
		// Step 2: Query all entities that have both Pet and Activity components
		q := cardinal.NewSearch().Entity(
			filter.Contains(filter.Component[component.Pet](), filter.Component[component.Activity]()))

		return q.Each(world, func(petId types.EntityID) bool {
			// Step 3: Retrieve the Activity component for the current entity
			activity, err := cardinal.GetComponent[component.Activity](world, petId)
			if err != nil {
				// Log the error and continue processing other entities
				// Step 3.1: Handle error during component retrieval
				return true
			}

			// Step 4: Check if the activity is not "None"
			if activity.Activity != "None" {
				// Step 5: Decrement the activity duration by one
				activity.CountDown--
				// Decrement the activity duration if it is greater than zero
				if activity.CountDown > 0 {
					// Update player money
					pet, err := cardinal.GetComponent[component.Pet](world, petId)
					if err != nil {
						// Log the error and continue processing other entities
						// Step 3.1: Handle error during component retrieval
						return true
					}
					playerId, err := component.FindPlayerByPersonaTag(world, pet.PersonaTag)
					if err != nil {
						return true
					}

					// Step 6: Update the activity percentage
					if activity.TotalTicks != 0 {
						activity.Percentage = int((float64(activity.CountDown) / float64(activity.TotalTicks)) * 100)
					} else {
						activity.Percentage = 0
					}
					// Step 7: Update the `Activity` component
					if err := cardinal.SetComponent(world, petId, activity); err != nil {
						// Log the error and continue processing other entities
						// Step 7.1: Handle error during component update
						log.Printf("Error updating activity component for entity %v: %v", petId, err)
						return true
					}

					// Step 8: Update the `Player` Money component again
					if err := component.IncreasePlayerMoney(world, playerId, game.PetEarnMoney); err != nil {
						// Log the error and continue processing other entities
						// Step 8.1: Handle error during component update
						log.Printf("Error updating player component for entity %v: %v", petId, err)
						return true
					}

				} else {
					// Step 9: Reset the activity to "None" when duration reaches zero
					activity.Activity = "None"
					activity.CountDown = 0
					activity.Percentage = 0
					activity.TotalTicks = 0
					// Step 10: Update the `Activity` component
					if err := cardinal.SetComponent(world, petId, activity); err != nil {
						// Log the error and continue processing other entities
						// Step 10.1: Handle error during component update
						return true
					}
				}
			}

			// Step 11: Continue processing other entities
			return true
		})
	}
	// Step 12: Return nil if the current tick is not a multiple of `game.ActivityUpdateTickRate`
	return nil
}
