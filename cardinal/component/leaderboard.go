// Package component contains structures and functions for working with game components.
package component

import (
	"sort"

	"pkg.world.dev/world-engine/cardinal"
)

// LeaderboardSize is the maximum number of pets that can be stored in the leaderboard.
const LeaderboardSize = 2

/**
 * Leaderboard represents a ranked list of pets.
 *
 * Code Flow:
 *   This struct has no specific code flow as it is a simple data structure.
 *   However, it is used in conjunction with other components and functions to manage the leaderboard in the game.
 */
type Leaderboard struct {
	/**
	 * Pets is a slice of pointers to Pet structs that are stored in the leaderboard.
	 */
	Pets []*Pet
}

/**
 * Name returns the name of the Leaderboard component.
 *
 * Code Flow:
 * 1. Return the string "Leaderboard" as the name of the component.
 *
 * Parameters:
 *   None
 *
 * Returns:
 *   (string): The name of the Leaderboard component.
 *
 * Step-by-Step Explanation:
 *   Step 1: This method returns the name of the component as a string. The name is used to identify the component in the game world.
 */
func (Leaderboard) Name() string {
	// Step 1: Return the string "Leaderboard" as the name of the component
	//         This method is used to identify the component in the game world.
	return "Leaderboard"
}

/**
 * AddPetToLeaderboard adds a pet to the leaderboard if it meets certain conditions.
 *
 * Code Flow:
 * 1. Check if the pet's level is greater than the lowest level in the leaderboard.
 * 2. Then check if the pet's persona tag does not already exist in the leaderboard.
 * 3. Add the pet to the leaderboard according to its level, maintaining the sorted order.
 * 4. Enforce the size limit of the leaderboard by removing the lowest level pet if necessary.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The game world context.
 *   addPet (Pet): The pet to be added to the leaderboard.
 *
 * Returns:
 *   None
 *
 * Step-by-Step Explanation:
 *   Step 1: Check if the pet's level is greater than the lowest level in the leaderboard.
 *           If not, the pet is not added to the leaderboard.
 *   Step 2: Check if the pet's persona tag does not already exist in the leaderboard.
 *           If it does, the pet is not added to the leaderboard.
 *   Step 3: Add the pet to the leaderboard according to its level, maintaining the sorted order.
 *   Step 4: Enforce the size limit of the leaderboard by removing the lowest level pet if necessary.
 */
func (leaderboard *Leaderboard) AddPetToLeaderboard(world cardinal.WorldContext, addPet Pet) {
	log := world.Logger()
	n := len(leaderboard.Pets)
	// Step 1: Check if the pet's level is greater than the lowest level in the leaderboard
	//         If not, the pet is not added to the leaderboard
	// Early exit: leaderboard size
	if n < LeaderboardSize {
		log.Info().Msgf("Early Appending: [%s]", addPet.Nickname)
		leaderboard.Pets = append(leaderboard.Pets, &addPet)
		return
	}

	// Step 2: Check if the pet's persona tag does not already exist in the leaderboard
	//         If it does, the pet is not added to the leaderboard
	for _, pet := range leaderboard.Pets {
		log.Info().Msgf("Discarding: (already on leaderboard) [%s]", addPet.Nickname)
		if pet != nil && pet.Nickname == addPet.Nickname {
			return
		}
	}

	// Step 1: Check if the pet's level is greater than the lowest level in the leaderboard
	//         If not, the pet is not added to the leaderboard
	// lower level than the lowest
	if addPet.Level <= leaderboard.Pets[n-1].Level {
		log.Info().Msgf("Discarding: [%s] Lvl[%d] -> min Lvl[%d] ", addPet.Nickname, addPet.Level, leaderboard.Pets[n-1].Level)
		return
	}

	// Step 3: Add the pet to the leaderboard according to its level, maintaining the sorted order
	// Add the pet (unsorted)
	leaderboard.Pets = append(leaderboard.Pets, &addPet)

	// Sort the leaderboard (most efficient way to maintain sorted order after a single insert)
	sort.Slice(leaderboard.Pets, func(i, j int) bool {
		return leaderboard.Pets[i].Level > leaderboard.Pets[j].Level // Descending order (highest level first)
	})

	// Step 4: Enforce the size limit of the leaderboard by removing the lowest level pet if necessary
	if len(leaderboard.Pets) > LeaderboardSize {
		leaderboard.Pets = leaderboard.Pets[:LeaderboardSize]
	}
}
