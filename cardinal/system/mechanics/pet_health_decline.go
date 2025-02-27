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
 * 1. The `HealthDeclineSystem` function is called, which checks if the current tick is a multiple of `game.DeclineTickRate`.
 * 2. If it is, the function queries all entities that have `Pet`, `Health`, and `Hygiene` components.
 * 3. For each entity found, the function retrieves the `Hygiene` component and checks if the hygiene value is less than or equal to `game.HygieneThreshold`.
 * 4. If the hygiene value is less than or equal to `game.HygieneThreshold`, the function retrieves the `Health` component and checks if the health value is greater than zero.
 * 5. If the health value is greater than zero, the function decrements the health value by one.
 * 6. The function updates the `Health` component with the new health value.
 * 7. If any error occurs during the execution of the health decline system, the function logs the error and continues processing other entities.
 *
 * HealthDeclineSystem declines the pet's Hy every `HealthDeclineTicksPerSecond` tick.
 *
 * This system iterates over all entities that have `Pet`, `Health`, and `Hygiene` components,
 * reducing the health value by one each time it is processed if the hygiene value is less than or equal to `game.HygieneThreshold`.
 *
 * The function returns an error if there is a failure during component access or update.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the execution of the health decline system.
 */
func HealthDeclineSystem(world cardinal.WorldContext) error {
	// Step 1: Check if the current tick is a multiple of `game.DeclineTickRate`
	if world.CurrentTick()%game.DeclineTickRate == 0 {
		// Step 2: Query all entities that have `Pet`, `Health`, and `Hygiene` components
		q := cardinal.NewSearch().Entity(
			filter.Contains(
				// Step 2.1: Filter entities with Pet component
				filter.Component[component.Pet](),
				// Step 2.2: Filter entities with Health component
				filter.Component[component.Health](),
				// Step 2.3: Filter entities with Hygiene component
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

				// Step 4: Check if the hygiene value is less than or equal to `game.HygieneThreshold`
				if hygiene.Hy <= game.HygieneThreshold {
					health, err := cardinal.GetComponent[component.Health](world, id)
					if err != nil {
						// Step 4.1: Handle error during component retrieval
						return true
					}
					// Step 5: Check if the health value is greater than zero and decrement it by one
					if health.HP > 0 {
						health.HP--
					}

					// Step 6: Update the Health component with the new health value
					if err := cardinal.SetComponent(world, id, health); err != nil {
						// Step 6.1: Handle error during component update
						return true
					}
				}
				// Step 7: Continue processing other entities
				return true
			})
	} else {
		// Step 8: Return nil if the current tick is not a multiple of `game.DeclineTickRate`
		return nil
	}
}
