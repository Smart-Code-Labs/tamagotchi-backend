// Package component contains various components for the Tamagotchi game.
package component

// Wellness represents a wellness component that tracks the wellness of an entity.
type Wellness struct {
	// Wn is the wellness value of the entity.
	Wn int `json:"wellness"`
}

/**
 * Name returns the name of the component.
 *
 * Code Flow:
 * 1. The function simply returns the string "Wellness" as the name of the component.
 *
 * Returns:
 *   string: The name of the component.
 */
func (Wellness) Name() string {
	// Step 1: Return the string "Wellness" as the name of the component.
	// This step is necessary to identify the component in the system.
	// The function returns the string "Wellness" to indicate that this is the wellness component.
	return "Wellness"
}
