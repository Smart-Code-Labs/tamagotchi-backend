package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
	system "tamagotchi/system"
)

// ThinkSystem declines the pet's Hy every `ThinkDeclineTicksPerSecond` tick.
func ThinkSystem(world cardinal.WorldContext) error {
	var petActivity *comp.Activity
	var err error
	var found bool
	var message string

	if world.CurrentTick()%constants.ThinkTickRate == 0 {

		q := cardinal.NewSearch().Entity(
			filter.Contains(
				filter.Component[comp.Pet](),
				filter.Component[comp.Activity](),
				filter.Component[comp.Think]()))

		return q.
			Each(world, func(petId types.EntityID) bool {
				// check if not activity
				petActivity, err = cardinal.GetComponent[comp.Activity](world, petId)
				if err != nil {
					return true
				}

				// If he has activity, he doesnt think, just do the activity
				if petActivity.CountDown > 0 {
					return true
				}

				// no activity, so lets think
				petThink, err := cardinal.GetComponent[comp.Think](world, petId)
				if err != nil {
					return true
				}

				// Health?
				health, err := cardinal.GetComponent[comp.Health](world, petId)
				if err != nil {
					return true
				}
				found, message = system.GetMessageForRange(health.HP, constants.HealthMessages)

				if found {
					petThink.Think = message
					if err := cardinal.SetComponent(world, petId, petThink); err != nil {
						return true
					}
				}

				// Hygiene?
				hygiene, err := cardinal.GetComponent[comp.Hygiene](world, petId)
				if err != nil {
					return true
				}
				found, message = system.GetMessageForRange(hygiene.Hy, constants.HygieneMessages)

				if found {
					petThink.Think = message
					if err := cardinal.SetComponent(world, petId, petThink); err != nil {
						return true
					}
				}

				// Wellness?
				wellness, err := cardinal.GetComponent[comp.Wellness](world, petId)
				if err != nil {
					return true
				}
				found, message = system.GetMessageForRange(wellness.Wn, constants.WellnessMessages)

				if found {
					petThink.Think = message
					if err := cardinal.SetComponent(world, petId, petThink); err != nil {
						return true
					}
				}

				// Energy?
				energy, err := cardinal.GetComponent[comp.Energy](world, petId)
				if err != nil {
					return true
				}
				found, message = system.GetMessageForRange(energy.E, constants.EnergyMessages)

				if found {
					petThink.Think = message
					if err := cardinal.SetComponent(world, petId, petThink); err != nil {
						return true
					}
				}

				return true
			})
	} else {
		return nil
	}
}
