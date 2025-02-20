package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
)

const LeaderboardTicksPerSecond = 60

// LeaderboardSystem declines the pet's E every `LeaderboardTicksPerSecond` tick.
func LeaderboardSystem(world cardinal.WorldContext) error {
	log := world.Logger()
	if world.CurrentTick()%LeaderboardTicksPerSecond == 0 {
		log.Info().Msgf("Running Leaderbord system")

		// get the `unique` Leaderboard component
		q := cardinal.NewSearch().Entity(
			filter.Exact(filter.Component[comp.Leaderboard]()))

		q.Each(world, func(id types.EntityID) bool {
			leaderboard, err := cardinal.GetComponent[comp.Leaderboard](world, id)
			if err != nil {
				return true
			}

			// Iterate over all pets and add to leaderboard
			q2 := cardinal.NewSearch().Entity(
				filter.Contains(filter.Component[comp.Pet]()))

			q2.Each(world, func(id types.EntityID) bool {
				pet, err := cardinal.GetComponent[comp.Pet](world, id)
				if err != nil {
					return true
				}

				leaderboard.AddPetToLeaderboard(world, *pet)
				return true
			})

			if err := cardinal.SetComponent(world, id, leaderboard); err != nil {
				return true
			}
			return true
		})
		return nil
	} else {
		return nil
	}
}
