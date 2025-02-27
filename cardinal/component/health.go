// Package component contains structures and functions for working with game components.
package component

/**
 * Health represents a component that stores the health points (HP) of an entity.
 *
 * Code Flow:
 *   This struct has no specific code flow as it is a simple data structure.
 *   However, it is used in conjunction with other components and functions to manage the health of entities in the game.
 */
type Health struct {
	HP int `json:"health"`
}

/**
 * Name returns the name of the Health component.
 *
 * Code Flow:
 * 1. Return the string "Health" as the name of the component.
 *
 * Parameters:
 *   None
 *
 * Returns:
 *   (string): The name of the Health component.
 *
 * Step-by-Step Explanation:
 *   Step 1: This method returns the name of the component as a string. The name is used to identify the component in the game world.
 */
func (Health) Name() string {
	// Step 1: Return the string "Health" as the name of the component
	//         This method is used to identify the component in the game world.
	return "Health"
}
