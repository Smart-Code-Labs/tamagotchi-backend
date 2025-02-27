package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	"tamagotchi/component"
	"tamagotchi/game"
	"tamagotchi/system"
)

/**
 * Function Flow:
 * 1. The `ThinkSystem` function is called, which checks if the current tick is a multiple of `game.ThinkTickRate`.
 * 2. If it is, the function queries all entities that have `Pet`, `Activity`, and `Think` components.
 * 3. For each entity found, the function checks if the entity has an activity.
 * 4. If the entity has an activity, the function skips the thinking process.
 * 5. If the entity does not have an activity, the function retrieves the `Think` component and checks the entity's health, hygiene, wellness, and energy.
 * 6. For each checked component, the function generates a message based on the component's value and updates the `Think` component with the message.
 * 7. The function returns an error if there is a failure during component access or update.
 *
 * ThinkSystem generates a thought for the pet based on its current health, hygiene, wellness, and energy.
 *
 * This system iterates over all entities that have `Pet`, `Activity`, and `Think` components,
 * and updates the `Think` component based on the entity's current state.
 *
 * The function returns an error if there is a failure during component access or update.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the execution of the think system.
 */
func ThinkSystem(world cardinal.WorldContext) error {
	// Step 1: Check if the current tick is a multiple of `game.ThinkTickRate`
	if world.CurrentTick()%game.ThinkTickRate == 0 {
		// Step 2: Query all entities that have Pet, Activity, and Think components
		q := cardinal.NewSearch().Entity(
			filter.Contains(
				// Step 2.1: Filter entities with Pet component
				filter.Component[component.Pet](),
				// Step 2.2: Filter entities with Activity component
				filter.Component[component.Activity](),
				// Step 2.3: Filter entities with Think component
				filter.Component[component.Think]()))

		return q.
			// Step 3: For each entity found, check if the entity has an activity
			Each(world, func(petId types.EntityID) bool {
				// Step 3.1: Retrieve the Activity component
				petActivity, err := cardinal.GetComponent[component.Activity](world, petId)
				if err != nil {
					// Step 3.1.1: Handle error during component retrieval
					return true
				}

				// Step 4: If the entity has an activity, skip the thinking process
				if petActivity.CountDown > 0 {
					// Step 4.1: Skip to the next entity
					return true
				}

				// Step 5: If the entity does not have an activity, retrieve the Think component
				petThink, err := cardinal.GetComponent[component.Think](world, petId)
				if err != nil {
					// Step 5.1: Handle error during component retrieval
					return true
				}

				// Step 6: Check the entity's health and generate a message
				health, err := cardinal.GetComponent[component.Health](world, petId)
				if err != nil {
					// Step 6.1: Handle error during component retrieval
					return true
				}
				found, message := system.GetMessageForRange(health.HP, game.HealthMessages)

				if found {
					// Step 6.2: Update the Think component with the health message
					petThink.Think = message
					if err := cardinal.SetComponent(world, petId, petThink); err != nil {
						// Step 6.2.1: Handle error during component update
						return true
					}
				}

				// Step 7: Check the entity's hygiene and generate a message
				hygiene, err := cardinal.GetComponent[component.Hygiene](world, petId)
				if err != nil {
					// Step 7.1: Handle error during component retrieval
					return true
				}
				found, message = system.GetMessageForRange(hygiene.Hy, game.HygieneMessages)

				if found {
					// Step 7.2: Update the Think component with the hygiene message
					petThink.Think = message
					if err := cardinal.SetComponent(world, petId, petThink); err != nil {
						// Step 7.2.1: Handle error during component update
						return true
					}
				}

				// Step 8: Check the entity's wellness and generate a message
				wellness, err := cardinal.GetComponent[component.Wellness](world, petId)
				if err != nil {
					// Step 8.1: Handle error during component retrieval
					return true
				}
				found, message = system.GetMessageForRange(wellness.Wn, game.WellnessMessages)

				if found {
					// Step 8.2: Update the Think component with the wellness message
					petThink.Think = message
					if err := cardinal.SetComponent(world, petId, petThink); err != nil {
						// Step 8.2.1: Handle error during component update
						return true
					}
				}

				// Step 9: Check the entity's energy and generate a message
				energy, err := cardinal.GetComponent[component.Energy](world, petId)
				if err != nil {
					// Step 9.1: Handle error during component retrieval
					return true
				}
				found, message = system.GetMessageForRange(energy.E, game.EnergyMessages)

				if found {
					// Step 9.2: Update the Think component with the energy message
					petThink.Think = message
					if err := cardinal.SetComponent(world, petId, petThink); err != nil {
						// Step 9.2.1: Handle error during component update
						return true
					}
				}

				// Step 10: Continue to the next entity
				return true
			})
	} else {
		// Step 11: Return nil if the current tick is not a multiple of `game.ThinkTickRate`
		return nil
	}
}
