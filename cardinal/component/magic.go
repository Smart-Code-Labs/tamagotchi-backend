// Package component contains structures and functions for working with game components.
package component

/**
 * Magic represents a magic component.
 *
 * Code Flow:
 *   This struct has no specific code flow as it is a simple data structure.
 *   However, it is used in conjunction with other components and functions to manage magic in the game.
 */
type Magic struct {
	/**
	 * Kind is the type of magic.
	 */
	Kind string
	/**
	 * Level is the current level of the magic.
	 */
	Level int64 `json:"lvl"`
	/**
	 * XP is the current experience points of the magic.
	 */
	XP int64 `json:"exp"`
	/**
	 * NextLevelXP is the experience points required to reach the next level.
	 */
	NextLevelXP int64
}

/**
 * Name returns the name of the Magic component.
 *
 * Code Flow:
 * 1. Return the string "Magic" as the name of the component.
 *
 * Parameters:
 *   None
 *
 * Returns:
 *   (string): The name of the Magic component.
 *
 * Step-by-Step Explanation:
 *   Step 1: This method returns the name of the component as a string. The name is used to identify the component in the game world.
 */
func (Magic) Name() string {
	// Step 1: Return the string "Magic" as the name of the component
	//         This method is used to identify the component in the game world.
	return "Magic"
}
