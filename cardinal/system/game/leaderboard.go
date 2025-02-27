// Package system contains the logic for handling leaderboard-related actions.
package system

import (
	"tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

/**
 * Constant representing the number of ticks per second for the leaderboard system.
 */
const LeaderboardTicksPerSecond = 60

/**
 * Function Flow:
 * 1. The `LeaderboardSystem` function is called, which checks if the current tick is a multiple of `LeaderboardTicksPerSecond`.
 * 2. If it is, the function logs a message indicating that the leaderboard system is running.
 * 3. The function then searches for the unique `Leaderboard` component in the world.
 * 4. For each `Leaderboard` component found, the function iterates over all `Pet` components in the world and adds them to the leaderboard.
 * 5. The function updates the `Leaderboard` component with the new list of pets.
 *
 * LeaderboardSystem declines the pet's E every `LeaderboardTicksPerSecond` tick.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the execution of the leaderboard system.
 */
func LeaderboardSystem(world cardinal.WorldContext) error {
	// Step 1: Check if the current tick is a multiple of `LeaderboardTicksPerSecond`

	if world.CurrentTick()%LeaderboardTicksPerSecond == 0 {
		// Step 2: Log a message indicating that the leaderboard system is running

		// Step 3: Search for the unique `Leaderboard` component in the world
		q := cardinal.NewSearch().Entity(
			filter.Exact(filter.Component[component.Leaderboard]()))

		// Step 4: Iterate over all `Leaderboard` components found and add pets to leaderboard
		q.Each(world, func(id types.EntityID) bool {
			// Step 4.1: Get the `Leaderboard` component
			leaderboard, err := cardinal.GetComponent[component.Leaderboard](world, id)
			if err != nil {
				// Step 4.2: Return true to continue to the next entity if an error occurs
				return true
			}

			// Step 4.3: Iterate over all `Pet` components in the world and add them to the leaderboard
			q2 := cardinal.NewSearch().Entity(
				filter.Contains(filter.Component[component.Pet]()))

			q2.Each(world, func(id types.EntityID) bool {
				// Step 4.3.1: Get the `Pet` component
				pet, err := cardinal.GetComponent[component.Pet](world, id)
				if err != nil {
					// Step 4.3.2: Return true to continue to the next entity if an error occurs
					return true
				}

				// Step 4.3.3: Add the pet to the leaderboard
				leaderboard.AddPetToLeaderboard(world, *pet)
				return true
			})

			// Step 4.4: Update the `Leaderboard` component with the new list of pets
			if err := cardinal.SetComponent(world, id, leaderboard); err != nil {
				// Step 4.5: Return true to continue to the next entity if an error occurs
				return true
			}
			return true
		})
		return nil
	} else {
		// Step 5: Return nil if the current tick is not a multiple of `LeaderboardTicksPerSecond`
		return nil
	}
}
