package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
)

const WellnessDeclineTicksPerSecond = 6

// WellnessDeclineSystem declines the pet's E every `WellnessDeclineTicksPerSecond` tick.
func WellnessDeclineSystem(world cardinal.WorldContext) error {
	if world.CurrentTick()%WellnessDeclineTicksPerSecond == 0 {

		q := cardinal.NewSearch().Entity(
			filter.Contains(filter.Component[comp.Pet](), filter.Component[comp.Wellness]()))

		return q.
			Each(world, func(id types.EntityID) bool {
				wellness, err := cardinal.GetComponent[comp.Wellness](world, id)
				if err != nil {
					return true
				}
				if wellness.Wn > 0 {
					wellness.Wn--
				}

				if err := cardinal.SetComponent(world, id, wellness); err != nil {
					return true
				}
				return true
			})
	} else {
		return nil
	}
}
