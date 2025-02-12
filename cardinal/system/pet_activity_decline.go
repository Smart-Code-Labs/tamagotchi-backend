package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
)

// ActivityDeclineSystem declines the pet's E every `ActivityDeclineTicksPerSecond` tick.
func ActivityDeclineSystem(world cardinal.WorldContext) error {
	q := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.Pet](), filter.Component[comp.Activity]()))

	return q.
		Each(world, func(id types.EntityID) bool {
			Activity, err := cardinal.GetComponent[comp.Activity](world, id)
			if err != nil {
				return true
			}
			if Activity.Duration > 0 {
				Activity.Duration--
				if err := cardinal.SetComponent(world, id, Activity); err != nil {
					return true
				}
			}

			return true
		})
}
