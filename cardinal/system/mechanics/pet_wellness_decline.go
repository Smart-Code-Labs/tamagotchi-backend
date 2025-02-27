// Package system contains game mechanics for the Tamagotchi game.
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
 * 1. The `WellnessDeclineSystem` function is called, which checks if the current tick is a multiple of `game.DeclineTickRate`.
 * 2. If it is, the function queries all entities that have `Pet` and `Wellness` components.
 * 3. For each entity found, the function retrieves the `Wellness` component and checks its current wellness value (`Wn`).
 * 4. If the wellness value is greater than 0, the function decrements the wellness value by 1.
 * 5. The function updates the `Wellness` component with the new wellness value.
 * 6. If any error occurs during component retrieval or update, the function returns an error.
 *
 * WellnessDeclineSystem declines the pet's wellness every `game.DeclineTickRate` tick.
 *
 * This system iterates over all entities that have `Pet` and `Wellness` components,
 * and updates the `Wellness` component by decrementing the wellness value.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the execution of the wellness decline system.
 */
// WellnessDeclineSystem declines the pet's E every `WellnessDeclineTicksPerSecond` tick.
func WellnessDeclineSystem(world cardinal.WorldContext) error {
	// Step 1: Check if the current tick is a multiple of `game.DeclineTickRate`
	if world.CurrentTick()%game.DeclineTickRate == 0 {
		// Step 2: Query all entities that have Pet and Wellness components
		q := cardinal.NewSearch().Entity(
			filter.Contains(
				// Step 2.1: Filter entities with Pet component
				filter.Component[component.Pet](),
				// Step 2.2: Filter entities with Wellness component
				filter.Component[component.Wellness]()))

		return q.
			// Step 3: For each entity found, retrieve the Wellness component and check its current wellness value (Wn)
			Each(world, func(id types.EntityID) bool {
				wellness, err := cardinal.GetComponent[component.Wellness](world, id)
				if err != nil {
					// Step 3.1: Handle error during component retrieval
					return true
				}
				// Step 4: If the wellness value is greater than 0, decrement the wellness value by 1
				if wellness.Wn > 0 {
					wellness.Wn--
				}

				// Step 5: Update the Wellness component with the new wellness value
				if err := cardinal.SetComponent(world, id, wellness); err != nil {
					// Step 5.1: Handle error during component update
					return true
				}
				// Step 6: Continue to the next entity
				return true
			})
	} else {
		// If the current tick is not a multiple of `game.DeclineTickRate`, return nil
		return nil
	}
}
