package system

import (
	"fmt"
	"tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

/**
 * CheckPetActivity checks if the pet is not currently engaged in an activity.
 *
 * Code Flow:
 * 1. Fetch the pet's activity component.
 * 2. Check if the activity countdown is greater than 0.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *   petId (types.EntityID): The ID of the pet.
 *
 * Returns:
 *   error: Any error that occurs during the process.
 */
func CheckPetActivity(world cardinal.WorldContext, petId types.EntityID) error {
	petActivity, err := cardinal.GetComponent[component.Activity](world, petId)
	if err != nil {
		return fmt.Errorf("failed to check activity [get Activity]: %w", err)
	}

	if petActivity.CountDown > 0 {
		return fmt.Errorf("failed to Sleep [already on Activity]: %w", err)
	}

	return nil
}

/**
 * CheckPlayerExistence checks if the player exists and is valid.
 *
 * Code Flow:
 * 1. Use `component.FindPlayerByPersonaTag` to retrieve the player's ID by their persona tag.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *   personaTag (string): The player's persona tag.
 *
 * Returns:
 *   (types.EntityID, error): The player's ID, and any error that occurs during the process.
 */
func CheckPlayerExistence(world cardinal.WorldContext, personaTag string) (types.EntityID, error) {
	playerID, err := component.FindPlayerByPersonaTag(world, personaTag)
	if err != nil {
		return 0, err
	}
	return playerID, nil
}

/**
 * CheckItemExistence checks if the item to be bought exists.
 *
 * Code Flow:
 * 1. Use `component.FindItemByName` to retrieve the item's ID by its name.
 *
 * Parameters:
 *   world (cardinal.WorldContext): The world context.
 *   itemName (string): The name of the item.
 *
 * Returns:
 *   (types.EntityID, error): The item's ID, and any error that occurs during the process.
 */
func CheckItemExistence(world cardinal.WorldContext, itemName string) (types.EntityID, error) {
	itemId, err := component.FindItemByName(world, itemName)
	if err != nil {
		return 0, err
	}
	return itemId, nil
}
