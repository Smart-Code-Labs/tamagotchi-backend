package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
)

// HealthDeclineSystem declines the pet's Hy every `HealthDeclineTicksPerSecond` tick.
func HealthDeclineSystem(world cardinal.WorldContext) error {
	if world.CurrentTick()%constants.HealthDeclineTicksPerSecond == 0 {

		q := cardinal.NewSearch().Entity(
			filter.Contains(
				filter.Component[comp.Pet](),
				filter.Component[comp.Health](),
				filter.Component[comp.Hygiene](),
			))

		return q.
			Each(world, func(id types.EntityID) bool {
				hygiene, err := cardinal.GetComponent[comp.Hygiene](world, id)
				if err != nil {
					return true
				}

				if hygiene.Hy <= constants.HygieneThresshold {
					health, err := cardinal.GetComponent[comp.Health](world, id)
					if err != nil {
						return true
					}
					if health.HP > 0 {
						health.HP--
					}

					if err := cardinal.SetComponent(world, id, health); err != nil {
						return true
					}
				}
				return true
			})
	} else {
		return nil
	}
}
