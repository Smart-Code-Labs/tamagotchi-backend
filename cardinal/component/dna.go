// Package component contains structures and functions for working with game components.
package component

// Dna represents a DNA structure, including the counts of adenine (A), cytosine (C), guanine (G), and thymine (T).
type Dna struct {
	A int
	C int
	G int
	T int
}

/**
 * Name returns the name of the Dna component.
 *
 * Code Flow:
 * 1. Return the string "Dna" as the name of the component.
 *
 * Parameters:
 *   None
 *
 * Returns:
 *   (string): The name of the Dna component.
 */
// Name returns the name of the Dna component
func (Dna) Name() string {
	// Step 1: Return the string "Dna" as the name of the component
	//         This method is used to identify the component in the game world.
	return "Dna"
}
