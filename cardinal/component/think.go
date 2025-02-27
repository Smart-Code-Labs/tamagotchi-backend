// Package component contains various components for the Tamagotchi game.
package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// Think represents a think component that an entity can have.
type Think struct {
	// Think is the thought of the entity.
	Think string
}

/**
 * Name returns the name of the component.
 *
 * Code Flow:
 * 1. The function simply returns the string "Think" as the name of the component.
 *
 * Returns:
 *   string: The name of the component.
 */
func (Think) Name() string {
	// Step 1: Return the string "Think" as the name of the component.
	return "Think"
}

/**
 * GetPetThink retrieves the pet's think component.
 *
 * Code Flow:
 * 1. Fetch the pet's think component from the world context using the pet's ID.
 * 2. Handle any error that occurs during the fetching process.
 * 3. Return the pet's think component and any error that occurred.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *   petId (types.EntityID): The ID of the pet.
 *
 * Returns:
 *   (*Think, error): The pet's think component, and any error that occurs during the process.
 */
func GetPetThink(world cardinal.WorldContext, petId types.EntityID) (*Think, error) {
	// Step 1: Fetch the pet's think component.
	petThink, err := cardinal.GetComponent[Think](world, petId)
	// Step 2: Handle any error that occurs during the fetching process.
	if err != nil {
		return nil, fmt.Errorf("failed to Sleep [get Think]: %w", err)
	}
	// Step 3: Return the pet's think component and any error that occurred.
	return petThink, nil
}
