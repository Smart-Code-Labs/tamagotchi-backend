// Package component contains structures and functions for working with game components.
package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// Activity represents a pet's activity, including the type, total ticks, countdown, and percentage.
type Activity struct {
	Activity   string
	TotalTicks int
	CountDown  int
	Percentage int
}

/**
 * Name returns the name of the Activity component.
 *
 * Code Flow:
 * 1. Return the string "Activity" as the name of the component.
 *
 * Returns:
 *   (string): The name of the Activity component.
 */
// Name returns the name of the Activity component
func (Activity) Name() string {
	// Step 1: Return the string "Activity" as the name of the component
	//         This method is used to identify the component in the game world.
	return "Activity"
}

/**
 * GetPetActivity retrieves the pet's activity component.
 *
 * Code Flow:
 * 1. Fetch the pet's activity component from the world context using the provided pet ID.
 * 2. Return the fetched activity component and any error that occurs during the process.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *   petId (types.EntityID): The ID of the pet.
 *
 * Returns:
 *   (*component.Activity, error): The pet's activity component, and any error that occurs during the process.
 */
// GetPetActivity retrieves the pet's activity component.
func GetPetActivity(world cardinal.WorldContext, petId types.EntityID) (*Activity, error) {
	// Step 1: Fetch the pet's activity component from the world context
	//         Use the Cardinal GetComponent API to fetch the activity component.
	petActivity, err := cardinal.GetComponent[Activity](world, petId)
	if err != nil {
		// Step 2: Return the fetched activity component and any error that occurs
		//         If an error occurs during the fetching process, return the error.
		return nil, fmt.Errorf("failed to Sleep [get Activity]: %w", err)
	}
	// Step 2 (continued): Return the fetched activity component
	//                    If no error occurs, return the fetched activity component.
	return petActivity, nil
}
