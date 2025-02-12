package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
)

const EnergyDeclineTicksPerSecond = 6

// EnergyDeclineSystem declines the pet's E every `EnergyDeclineTicksPerSecond` tick.
func EnergyDeclineSystem(world cardinal.WorldContext) error {
	if world.CurrentTick()%EnergyDeclineTicksPerSecond == 0 {

		q := cardinal.NewSearch().Entity(
			filter.Contains(filter.Component[comp.Pet](), filter.Component[comp.Energy]()))

		return q.
			Each(world, func(id types.EntityID) bool {
				energy, err := cardinal.GetComponent[comp.Energy](world, id)
				if err != nil {
					return true
				}

				if energy.E > 0 {
					energy.E--
				}

				if err := cardinal.SetComponent(world, id, energy); err != nil {
					return true
				}
				return true
			})
	} else {
		return nil
	}
}
