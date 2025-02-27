// Package component contains structures and functions for working with game components.
package component

import (
	"fmt"
	"math"
	"tamagotchi/game"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// nextLevelXP = baseXP * (level ^ growthRate)
const (
	baseXP     = int64(100)
	growthRate = 1.1
)

/**
 * Pet represents a pet in the game.
 *
 * Code Flow:
 *   This struct has no specific code flow as it is a simple data structure.
 *   However, it is used in conjunction with other components and functions to manage pets in the game.
 *
 * The following are the steps to understand the Pet struct:
 *   Step 1: The Pet struct has several fields that describe its characteristics, such as PersonaTag, Nickname, Gender, Level, XP, and NextLevelXP.
 *   Step 2: Each field has a specific purpose, such as identifying the pet, storing its name, determining its gender, tracking its level and experience points, and calculating the experience points required to reach the next level.
 */
type Pet struct {
	/**
	 * PersonaTag is a unique identifier for the pet.
	 */
	PersonaTag string `json:"personaTag"`
	/**
	 * Nickname is the name given to the pet.
	 */
	Nickname string `json:"nickname"`
	/**
	 * Gender is a boolean value indicating the pet's gender.
	 */
	Gender bool
	/**
	 * Level is the current level of the pet.
	 */
	Level int64 `json:"lvl"`
	/**
	 * XP is the current experience points of the pet.
	 */
	XP int64 `json:"exp"`
	/**
	 * NextLevelXP is the experience points required to reach the next level.
	 */
	NextLevelXP int64

	BornTick uint64 `json:"born_tick"`
}

/**
 * Name returns the name of the Pet component.
 *
 * Code Flow:
 *   Step 1: Return the string "Pet" as the name of the component.
 *
 * Parameters:
 *   None
 *
 * Returns:
 *   (string): The name of the Pet component.
 *
 * Step-by-Step Explanation:
 *   Step 1: This method returns the name of the component as a string. The name is used to identify the component in the game world.
 */
func (Pet) Name() string {
	return "Pet"
}

/**
 * AddXP adds experience points to the pet.
 *
 * Code Flow:
 *   Step 1: Add the given experience points to the pet's current experience points.
 *   Step 2: Check if the pet's experience points are greater than or equal to the experience points required to reach the next level.
 *   Step 3: If the pet's experience points are sufficient, call the LevelUp method to advance the pet to the next level.
 *
 * Parameters:
 *   xp (int64): The experience points to add to the pet.
 *
 * Returns:
 *   None
 *
 * Step-by-Step Explanation:
 *   Step 1: This method increases the pet's experience points by the given amount.
 *   Step 2: It then checks if the pet has enough experience points to advance to the next level.
 *   Step 3: If the pet has sufficient experience points, it calls the LevelUp method to advance the pet to the next level.
 */
func (h *Pet) AddXP(xp int64) {
	h.XP += xp
	for h.XP >= h.NextLevelXP {
		h.LevelUp()
	}
}

/**
 * LevelUp advances the pet to the next level.
 *
 * Code Flow:
 *   Step 1: Increment the pet's level by 1.
 *   Step 2: Calculate the excess experience points after advancing to the next level.
 *   Step 3: Update the pet's experience points to the excess experience points.
 *   Step 4: Calculate the experience points required to reach the next level using the CalculateNextLevelXP method.
 *
 * Parameters:
 *   None
 *
 * Returns:
 *   None
 *
 * Step-by-Step Explanation:
 *   Step 1: This method increments the pet's level by 1.
 *   Step 2: It calculates the excess experience points after advancing to the next level.
 *   Step 3: It updates the pet's experience points to the excess experience points.
 *   Step 4: It calculates the experience points required to reach the next level using the CalculateNextLevelXP method.
 */
func (h *Pet) LevelUp() {
	h.Level += 1
	excessXP := h.XP - h.NextLevelXP
	h.XP = excessXP
	h.NextLevelXP = CalculateNextLevelXP(h.Level, baseXP, growthRate)
}

/**
 * CalculateNextLevelXP calculates the experience points required to reach the next level.
 *
 * Code Flow:
 *   Step 1: Calculate the next level's experience points using the formula: nextLevelXP = baseXP * (level ^ growthRate).
 *
 * Parameters:
 *   level (int64): The current level of the pet.
 *   baseXP (int64): The base experience points required to reach the next level.
 *   growthRate (float64): The growth rate of the experience points required to reach the next level.
 *
 * Returns:
 *   (int64): The experience points required to reach the next level.
 *
 * Step-by-Step Explanation:
 *   Step 1: This method calculates the experience points required to reach the next level using the given formula.
 */
func CalculateNextLevelXP(level int64, baseXP int64, growthRate float64) int64 {
	nextLevelXP := int64(math.Pow(float64(level+1), growthRate)) * baseXP
	return nextLevelXP
}

/*
*
GetPetByNickname Function Flow:
1. Query the world for the entity ID of the pet with the given nickname.
2. Get the Pet component from the entity.
3. Return the entity ID and the Pet component.
*/
func GetPetByNickname(world cardinal.WorldContext, nickname string) (types.EntityID, *Pet, error) {
	// 1. Query the world for the entity ID of the pet with the given nickname.
	_, petId, err := QueryPetIdByName(world, nickname)
	if err != nil {
		return 0, nil, fmt.Errorf("error creating pet: [get %s EntityID]: %w", nickname, err)
	}
	// 2. Get the Pet component from the entity.
	pet, err := cardinal.GetComponent[Pet](world, petId)
	if err != nil {
		return 0, nil, fmt.Errorf("error creating pet [get %s]: %w", nickname, err)
	}
	// 3. Return the entity ID and the Pet component.
	return petId, pet, nil
}

func QueryPetIdByName(world cardinal.WorldContext, name string) (bool, types.EntityID, error) {
	var petID types.EntityID
	var err error
	var found = false
	log := world.Logger()

	q := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[Pet]()))

	count, err := q.Count(world)
	if count > 0 {
		searchErr := q.Each(world,
			func(id types.EntityID) bool {
				var pet *Pet
				pet, err = cardinal.GetComponent[Pet](world, id)
				if err != nil {
					return false
				}

				// Terminates the search if the pet is found
				if pet.Nickname == name {
					log.Info().Msgf("QueryPetIdByName Found it [%d]", petID)
					petID = id
					found = true
					return false
				}

				// Continue searching if the pet is not the target pet
				return true
			})
		if searchErr != nil {
			return found, 0, fmt.Errorf("error searching pet: %w", err)
		}
	}
	return found, petID, err
}

/**
 * CreateRandomPet creates a new pet with random characteristics.
 *
 * Code Flow:
 *   Step 1: Create a new entity with the Pet component using the cardinal.Create method.
 *   Step 2: Initialize the pet's characteristics, such as PersonaTag, Nickname, Level, XP, and NextLevelXP.
 *   Step 3: Generate random values for the pet's Gender and other characteristics.
 *   Step 4: Add the pet's components, such as Health, Energy, Hygiene, Wellness, Dna, Activity, and Think.
 *   Step 5: Return the entity ID of the newly created pet.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The game world context.
 *   personaTag (string): The persona tag of the pet.
 *   nickname (string): The nickname of the pet.
 *
 * Returns:
 *   (types.EntityID, error): The entity ID of the newly created pet, and an error if any.
 *
 * Step-by-Step Explanation:
 *   Step 1: This method creates a new entity with the Pet component using the cardinal.Create method.
 *   Step 2: It initializes the pet's characteristics, such as PersonaTag, Nickname, Level, XP, and NextLevelXP.
 *   Step 3: It generates random values for the pet's Gender and other characteristics.
 *   Step 4: It adds the pet's components, such as Health, Energy, Hygiene, Wellness, Dna, Activity, and Think.
 *   Step 5: It returns the entity ID of the newly created pet.
 */
func CreateRandomPet(world cardinal.WorldContext, personaTag string, nickname string) (types.EntityID, error) {
	rng := world.Rand()
	log := world.Logger()

	petID, err := cardinal.Create(world,
		Pet{PersonaTag: personaTag, Nickname: nickname, Level: 0, XP: 0, NextLevelXP: 0, Gender: rng.Intn(2) > 0, BornTick: world.CurrentTick()},
		Health{HP: game.MaxHP},
		Energy{E: game.MaxEnergy},
		Hygiene{Hy: game.MaxHygiene},
		Wellness{Wn: game.MaxWellness},
		Dna{
			A: rng.Intn(100),
			C: rng.Intn(100),
			G: rng.Intn(100),
			T: rng.Intn(100),
		},
		Activity{Activity: game.InitialActivity, CountDown: 0},
		Think{Think: game.InitialThink},
	)
	if err != nil {
		log.Error().Msgf("Failed to create pet with nickname %s: %v", nickname, err)
		return 0, fmt.Errorf("error creating pet: %w", err)
	}
	log.Info().Msgf("Created: Pet[%d] [%s]", petID, nickname)

	return petID, nil
}
