package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
)

// EnergyDeclineSystem declines the pet's E every `EnergyDeclineTicksPerSecond` tick.
func EnergyDeclineSystem(world cardinal.WorldContext) error {
	log := world.Logger()
	if world.CurrentTick()%constants.DeclineTickRate == 0 {

		q := cardinal.NewSearch().Entity(
			filter.Contains(filter.Component[comp.Pet](), filter.Component[comp.Energy]()))

		return q.
			Each(world, func(id types.EntityID) bool {
				energy, err := cardinal.GetComponent[comp.Energy](world, id)
				if err != nil {
					return true
				}

				log.Info().Msgf("Energy Decline: Energy Before[%d]", energy.E)
				if energy.E > 0 {
					energy.E--
				}
				log.Info().Msgf("Energy Decline: Energy After[%d]", energy.E)

				if err := cardinal.SetComponent(world, id, energy); err != nil {
					return true
				}
				return true
			})
	} else {
		return nil
	}
}
