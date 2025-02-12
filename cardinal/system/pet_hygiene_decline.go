package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
)

const HygieneDeclineTicksPerSecond = 12

// HygieneDeclineSystem declines the pet's Hy every `HygieneDeclineTicksPerSecond` tick.
func HygieneDeclineSystem(world cardinal.WorldContext) error {
	if world.CurrentTick()%HygieneDeclineTicksPerSecond == 0 {

		q := cardinal.NewSearch().Entity(
			filter.Contains(filter.Component[comp.Pet](), filter.Component[comp.Hygiene]()))

		return q.
			Each(world, func(id types.EntityID) bool {
				hygiene, err := cardinal.GetComponent[comp.Hygiene](world, id)
				if err != nil {
					return true
				}
				if hygiene.Hy > 0 {
					hygiene.Hy--
				}

				if err := cardinal.SetComponent(world, id, hygiene); err != nil {
					return true
				}
				return true
			})
	} else {
		return nil
	}
}
