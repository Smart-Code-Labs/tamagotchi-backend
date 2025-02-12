package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
)

const ThinkTicksPerSecond = 12

// ThinkSystem declines the pet's Hy every `ThinkDeclineTicksPerSecond` tick.
func ThinkSystem(world cardinal.WorldContext) error {
	var petActivity *comp.Activity
	var err error

	if world.CurrentTick()%ThinkTicksPerSecond == 0 {

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
				if petActivity.Duration > 0 {
					return true
				}

				// no activity, so lets think
				Think, err := cardinal.GetComponent[comp.Think](world, petId)
				if err != nil {
					return true
				}

				// Health?
				health, err := cardinal.GetComponent[comp.Health](world, petId)
				if err != nil {
					return true
				}
				if health.HP > 30 && health.HP < 40 {
					Think.Think = "I dont feel well."
					if err := cardinal.SetComponent(world, petId, Think); err != nil {
						return true
					}
				}
				if health.HP < 10 {
					Think.Think = "Im gona dye!!!"
					if err := cardinal.SetComponent(world, petId, Think); err != nil {
						return true
					}
				}

				// Hygiene?
				hygiene, err := cardinal.GetComponent[comp.Hygiene](world, petId)
				if err != nil {
					return true
				}
				if hygiene.Hy > 50 && hygiene.Hy < 60 {
					Think.Think = "my whole body itches"
					if err := cardinal.SetComponent(world, petId, Think); err != nil {
						return true
					}
				}
				if hygiene.Hy < 40 {
					Think.Think = "OMG!!! Im really dirty! Someone please Bath me. Please!"
					if err := cardinal.SetComponent(world, petId, Think); err != nil {
						return true
					}
				}

				// Wellness?
				wellness, err := cardinal.GetComponent[comp.Wellness](world, petId)
				if err != nil {
					return true
				}
				if wellness.Wn > 30 && wellness.Wn < 70 {
					Think.Think = "I fell a little depress today."
					if err := cardinal.SetComponent(world, petId, Think); err != nil {
						return true
					}
				}
				if wellness.Wn < 30 {
					Think.Think = "Dont know what to do..."
					if err := cardinal.SetComponent(world, petId, Think); err != nil {
						return true
					}
				}

				// Energy?
				energy, err := cardinal.GetComponent[comp.Energy](world, petId)
				if err != nil {
					return true
				}
				if energy.E > 70 && energy.E < 80 {
					Think.Think = "Im bored. I would kill to go outside."
					if err := cardinal.SetComponent(world, petId, Think); err != nil {
						return true
					}
				}
				if energy.E < 60 {
					Think.Think = "Im bored... to death? Play with me!"
					if err := cardinal.SetComponent(world, petId, Think); err != nil {
						return true
					}
				}

				return true
			})
	} else {
		return nil
	}
}
