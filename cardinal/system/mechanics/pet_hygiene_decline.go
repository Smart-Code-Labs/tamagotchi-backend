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
 * 1. The `HygieneDeclineSystem` function is called, which checks if the current tick is a multiple of `game.DeclineTickRate`.
 * 2. If it is, the function queries all entities that have both `Pet` and `Hygiene` components.
 * 3. For each entity found, the function retrieves the `Hygiene` component and checks if the hygiene value is greater than zero.
 * 4. If the hygiene value is greater than zero, the function decrements the hygiene value by one.
 * 5. The function updates the `Hygiene` component with the new hygiene value.
 * 6. If any error occurs during the execution of the hygiene decline system, the function logs the error and continues processing other entities.
 *
 * HygieneDeclineSystem declines the pet's Hy every `HygieneDeclineTicksPerSecond` tick.
 *
 * This system iterates over all entities that have both `Pet` and `Hygiene` components,
 * reducing the hygiene value by one each time it is processed if the hygiene value is greater than zero.
 *
 * The function returns an error if there is a failure during component access or update.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the execution of the hygiene decline system.
 */
func HygieneDeclineSystem(world cardinal.WorldContext) error {
	// Step 1: Check if the current tick is a multiple of `game.DeclineTickRate`
	if world.CurrentTick()%game.DeclineTickRate == 0 {
		// Step 2: Query all entities that have both Pet and Hygiene components
		q := cardinal.NewSearch().Entity(
			filter.Contains(
				// Step 2.1: Filter entities with Pet component
				filter.Component[component.Pet](),
				// Step 2.2: Filter entities with Hygiene component
				filter.Component[component.Hygiene](),
			))

		return q.
			// Step 3: For each entity found, retrieve the Hygiene component
			Each(world, func(id types.EntityID) bool {
				hygiene, err := cardinal.GetComponent[component.Hygiene](world, id)
				if err != nil {
					// Step 3.1: Handle error during component retrieval
					return true
				}
				// Step 4: Check if the hygiene value is greater than zero and decrement it by one
				if hygiene.Hy > 0 {
					hygiene.Hy--
				}

				// Step 5: Update the Hygiene component with the new hygiene value
				if err := cardinal.SetComponent(world, id, hygiene); err != nil {
					// Step 5.1: Handle error during component update
					return true
				}
				// Step 6: Continue processing other entities
				return true
			})
	} else {
		// Step 7: Return nil if the current tick is not a multiple of `game.DeclineTickRate`
		return nil
	}
}
