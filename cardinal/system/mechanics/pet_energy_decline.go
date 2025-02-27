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
 * 1. The `EnergyDeclineSystem` function is called, which checks if the current tick is a multiple of `game.DeclineTickRate`.
 * 2. If it is, the function queries all entities that have both `Pet` and `Energy` components.
 * 3. For each entity found, the function retrieves the `Energy` component and checks if the energy value is greater than zero.
 * 4. If the energy value is greater than zero, the function decrements the energy value by one.
 * 5. The function logs the energy values before and after the decrement operation.
 * 6. The function updates the `Energy` component with the new energy value.
 * 7. If any error occurs during the execution of the energy decline system, the function logs the error and continues processing other entities.
 *
 * EnergyDeclineSystem declines the pet's E every `EnergyDeclineTicksPerSecond` tick.
 *
 * This system iterates over all entities that have both Pet and Energy components,
 * reducing the energy value by one each time it is processed. If the energy value reaches zero,
 * the pet's energy is effectively depleted.
 *
 * The function returns an error if there is a failure during component access or update.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the execution of the energy decline system.
 */
func EnergyDeclineSystem(world cardinal.WorldContext) error {
	// Step 1: Check if the current tick is a multiple of `game.DeclineTickRate`
	log := world.Logger()
	if world.CurrentTick()%game.DeclineTickRate == 0 {
		// Step 2: Query all entities that have both Pet and Energy components
		q := cardinal.NewSearch().Entity(
			filter.Contains(filter.Component[component.Pet](), filter.Component[component.Energy]()))

		return q.
			// Step 3: For each entity found, retrieve the Energy component
			Each(world, func(id types.EntityID) bool {
				energy, err := cardinal.GetComponent[component.Energy](world, id)
				if err != nil {
					// Step 3.1: Handle error during component retrieval
					return true
				}

				// Step 4: Log the energy value before the decrement operation
				log.Info().Msgf("Energy Decline: Energy Before[%d]", energy.E)
				// Step 5: Decrement the energy value if it is greater than zero
				if energy.E > 0 {
					energy.E--
				}
				// Step 6: Log the energy value after the decrement operation
				log.Info().Msgf("Energy Decline: Energy After[%d]", energy.E)

				// Step 7: Update the Energy component with the new energy value
				if err := cardinal.SetComponent(world, id, energy); err != nil {
					// Step 7.1: Handle error during component update
					return true
				}
				// Step 8: Continue processing other entities
				return true
			})
	} else {
		// Step 9: Return nil if the current tick is not a multiple of `game.DeclineTickRate`
		return nil
	}
}
