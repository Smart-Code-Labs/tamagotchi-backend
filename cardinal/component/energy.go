// Package component contains structures and functions for working with game components.
package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// Energy represents the energy of an entity, including the current energy level.
type Energy struct {
	E int `json:"energy"`
}

/**
 * Name returns the name of the Energy component.
 *
 * Code Flow:
 * 1. Return the string "Energy" as the name of the component.
 *
 * Parameters:
 *   None
 *
 * Returns:
 *   (string): The name of the Energy component.
 */
// Name returns the name of the Energy component
func (Energy) Name() string {
	// Step 1: Return the string "Energy" as the name of the component
	//         This method is used to identify the component in the game world.
	return "Energy"
}

/**
 * GetPetEnergy retrieves the pet's energy component.
 *
 * Code Flow:
 * 1. Fetch the pet's energy component from the world context using the provided pet ID.
 * 2. Return the fetched energy component and any error that occurs during the process.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *   petId (types.EntityID): The ID of the pet.
 *
 * Returns:
 *   (*component.Energy, error): The pet's energy component, and any error that occurs during the process.
 */
// GetPetEnergy retrieves the pet's energy component.
func GetPetEnergy(world cardinal.WorldContext, petId types.EntityID) (*Energy, error) {
	// Step 1: Fetch the pet's energy component from the world context
	//         Use the Cardinal GetComponent API to fetch the energy component.
	petEnergy, err := cardinal.GetComponent[Energy](world, petId)
	if err != nil {
		// Step 2: Return the fetched energy component and any error that occurs
		//         If an error occurs during the fetching process, return the error.
		return nil, fmt.Errorf("failed to Sleep [get Energy]: %w", err)
	}
	// Step 2 (continued): Return the fetched energy component
	//                    If no error occurs, return the fetched energy component.
	return petEnergy, nil
}
