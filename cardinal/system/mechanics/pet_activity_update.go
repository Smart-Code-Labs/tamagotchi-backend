/*
Package system provides core game mechanics for the Tamagotchi backend.

This package contains various system-level functionalities that manage and update the state of game entities.
It includes systems for handling activities, declining pet statistics over time, and other essential game mechanics.
*/
package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
)

// ActivityDeclineSystem periodically decreases the duration of a pet's current activity.
//
// This system iterates over all entities that have both Pet and Activity components,
// reducing the activity duration by one each time it is processed. If the duration reaches zero,
// the activity is effectively completed.
//
// The function returns an error if there is a failure during component access or update.
func ActivityDeclineSystem(world cardinal.WorldContext) error {
	log := world.Logger()
	if world.CurrentTick()%constants.ActivityUpdateTickRate == 0 {
		// Query all entities that have both Pet and Activity components
		q := cardinal.NewSearch().Entity(
			filter.Contains(filter.Component[comp.Pet](), filter.Component[comp.Activity]()))

		return q.Each(world, func(id types.EntityID) bool {
			// Retrieve the Activity component for the current entity
			activity, err := cardinal.GetComponent[comp.Activity](world, id)
			if err != nil {
				// Log the error and continue processing other entities
				return true
			}

			if activity.Activity != "None" {
				activity.CountDown--
				// Decrement the activity duration if it is greater than zero
				if activity.CountDown > 0 {
					if activity.TotalTicks != 0 {
						activity.Percentage = int((float64(activity.CountDown) / float64(activity.TotalTicks)) * 100)
					} else {
						activity.Percentage = 0
					}
					if err := cardinal.SetComponent(world, id, activity); err != nil {
						// Log the error and continue processing other entities
						log.Printf("Error updating activity component for entity %v: %v", id, err)
						return true
					}
					if err := cardinal.SetComponent(world, id, activity); err != nil {
						// Log the error and continue processing other entities
						return true
					}
				} else {
					activity.Activity = "None"
					activity.CountDown = 0
					activity.Percentage = 0
					activity.TotalTicks = 0
					if err := cardinal.SetComponent(world, id, activity); err != nil {
						// Log the error and continue processing other entities
						return true
					}
				}
			}

			return true
		})
	}
	return nil
}
