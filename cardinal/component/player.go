// Package component contains structures and functions for working with game components.
package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// Player represents a player in the game.
type Player struct {
	PersonaTag string           `json:"personaTag"`
	Pets       []types.EntityID `json:"pets"`
	Items      []types.EntityID `json:"items"`
	Money      float64          `json:"money"`
}

// Name returns the name of the component.
//
// Code Flow:
// 1. Return the string "Player" as the name of the component.
func (Player) Name() string {
	return "Player"
}

// HasPet checks if the given EntityID exists in the Player's Pets slice.
//
// Code Flow:
// 1. Iterate over the Player's Pets slice.
// 2. Check if the given EntityID matches any of the IDs in the slice.
// 3. Return true if a match is found, false otherwise.
//
// Returns true if the ItemId is found, false otherwise.
func (p Player) HasPet(id types.EntityID) bool {
	for _, petID := range p.Pets {
		if petID == id {
			return true
		}
	}
	return false
}

// GetPetNickname gets the EntityID of a Pet with the specified nickname.
//
// Code Flow:
// 1. Iterate over the Player's Pets slice.
// 2. For each Pet, retrieve the Pet component from the world.
// 3. Check if the Pet's nickname matches the specified nickname.
// 4. Return the EntityID of the Pet if a match is found, or an error if not.
func (p Player) GetPetNickname(world cardinal.WorldContext, nickname string) (types.EntityID, error) {
	for _, petID := range p.Pets {
		pet, err := cardinal.GetComponent[Pet](world, petID)
		if err != nil {
			return 0, err
		}
		if pet.Nickname == nickname {
			return petID, nil
		}
	}
	return 0, fmt.Errorf("pet not found")
}

// HasItem checks if the given EntityID exists in the Player's Items slice.
//
// Code Flow:
// 1. Iterate over the Player's Items slice.
// 2. Check if the given EntityID matches any of the IDs in the slice.
// 3. Return true if a match is found, false otherwise.
//
// Returns true if the ItemId is found, false otherwise.
func (p Player) HasItem(id types.EntityID) bool {
	for _, itemID := range p.Items {
		if itemID == id {
			return true
		}
	}
	return false
}

// GetItemByName gets an Item with the specified name.
//
// Code Flow:
// 1. Iterate over the Player's Items slice.
// 2. For each Item, retrieve the Item component from the world.
// 3. Check if the Item's name matches the specified name.
// 4. Return the Item if a match is found, or an error if not.
func (p Player) GetItemByName(world cardinal.WorldContext, itemName string) (*Item, error) {
	for _, itemID := range p.Items {
		item, err := cardinal.GetComponent[Item](world, itemID)
		if err != nil {
			return nil, err
		}
		if item.ItemName == itemName {
			return item, nil
		}
	}
	return nil, fmt.Errorf("item not found")
}

// GetItemIdByName gets the EntityID of an Item with the specified name.
//
// Code Flow:
// 1. Iterate over the Player's Items slice.
// 2. For each Item, retrieve the Item component from the world.
// 3. Check if the Item's name matches the specified name.
// 4. Return the EntityID of the Item if a match is found, or an error if not.
func (p Player) GetItemIdByName(world cardinal.WorldContext, itemName string) (types.EntityID, error) {
	for _, itemID := range p.Items {
		item, err := cardinal.GetComponent[Item](world, itemID)
		if err != nil {
			return itemID, err
		}
		if item.ItemName == itemName {
			return itemID, nil
		}
	}
	return 0, fmt.Errorf("item not found")
}

// AddPlayerPet adds a new Pet ID to the Player's Pets array.
//
// Code Flow:
// 1. Retrieve the Player component from the world.
// 2. Append the new Pet ID to the Player's Pets array.
// 3. Update the Player component in the world.
func AddPlayerPet(world cardinal.WorldContext, playerID types.EntityID, petID types.EntityID) error {
	// Add the new pet to the player's Pets array
	player, err := cardinal.GetComponent[Player](world, playerID)
	if err != nil {
		return fmt.Errorf("error adding pet to player: %w", err)
	}

	player.Pets = append(player.Pets, petID)
	err = cardinal.SetComponent(world, playerID, player)
	if err != nil {
		return fmt.Errorf("error updating player pets: %w", err)
	}
	return nil
}

// AddPlayerItem adds a new Item ID to the Player's Items array.
//
// Code Flow:
// 1. Retrieve the Player component from the world.
// 2. Append the new Item ID to the Player's Items array.
// 3. Update the Player component in the world.
func AddPlayerItem(world cardinal.WorldContext, playerID types.EntityID, itemID types.EntityID) error {
	// Append the new item to the player's Items array
	player, err := cardinal.GetComponent[Player](world, playerID)
	if err != nil {
		return fmt.Errorf("error adding item to player: %w", err)
	}

	player.Items = append(player.Items, itemID)
	err = cardinal.SetComponent(world, playerID, player)
	if err != nil {
		return fmt.Errorf("error updating player items: %w", err)
	}
	return nil
}

// ReducePlayerMoney reduces the Player's money by a specified amount.
//
// Code Flow:
// 1. Retrieve the Player component from the world.
// 2. Subtract the specified amount from the Player's money.
// 3. Update the Player component in the world.
func ReducePlayerMoney(world cardinal.WorldContext, playerID types.EntityID, itemPrice float64) error {
	// Append the new item to the player's Items array
	player, err := cardinal.GetComponent[Player](world, playerID)
	if err != nil {
		return err
	}
	if player.Money-itemPrice < 0 {
		return fmt.Errorf("error Buying, no enough balance [%f] [%f]", player.Money, itemPrice)
	}

	player.Money -= itemPrice
	err = cardinal.SetComponent(world, playerID, player)
	if err != nil {
		return fmt.Errorf("error updating player money: %w", err)
	}
	return nil
}

// IncreasePlayerMoney increases the Player's money by a specified amount.
//
// Code Flow:
// 1. Retrieve the Player component from the world.
// 2. Add the specified amount to the Player's money.
// 3. Update the Player component in the world.
func IncreasePlayerMoney(world cardinal.WorldContext, playerID types.EntityID, quantity float64) error {
	// Append the new item to the player's Items array
	player, err := cardinal.GetComponent[Player](world, playerID)
	if err != nil {
		return fmt.Errorf("error adding pet to player: %w", err)
	}

	player.Money += quantity
	err = cardinal.SetComponent(world, playerID, player)
	if err != nil {
		return fmt.Errorf("error updating player money: %w", err)
	}
	return nil
}

// RemoveItem removes an Item ID from the Player's Items array.
//
// Code Flow:
// 1. Retrieve the Player component from the world.
// 2. Iterate over the Player's Items array and remove the specified Item ID.
// 3. Update the Player component in the world.
func RemoveItem(world cardinal.WorldContext, playerID types.EntityID, itemID types.EntityID) error {
	// Retrieve the player component
	player, err := cardinal.GetComponent[Player](world, playerID)
	if err != nil {
		return fmt.Errorf("error getting player: %w", err)
	}

	// Check if Items is nil and initialize if necessary
	if player.Items == nil {
		player.Items = make([]types.EntityID, 0)
	}

	// Check if the item exists
	var itemFound bool
	for i, existingItemID := range player.Items {
		if existingItemID == itemID {
			// Remove the item by slicing the array
			player.Items = append(player.Items[:i], player.Items[i+1:]...)
			itemFound = true
			break
		}
	}

	if !itemFound {
		return fmt.Errorf("item %d does not belong to the player", itemID)
	}

	// Save the updated player state
	err = cardinal.SetComponent(world, playerID, player)
	if err != nil {
		return fmt.Errorf("error updating player items: %w", err)
	}

	return nil
}

// FindPlayerByPersonaTag finds a Player with the given PersonaTag and returns their EntityID.
//
// Code Flow:
// 1. Create a search query for Players.
// 2. Iterate over the search results and check if the PersonaTag matches the specified PersonaTag.
// 3. Return the EntityID of the Player if a match is found, or an error if not.

// TODO: refactor this, dont like how the find overlaps with return 0 EntityID works if not found.
func FindPlayerByPersonaTag(world cardinal.WorldContext, personaTag string) (types.EntityID, error) {
	q := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[Player]()),
	)

	var playerId types.EntityID
	var found bool = false

	q.Each(world, func(id types.EntityID) bool {
		player, err := cardinal.GetComponent[Player](world, id)
		if err != nil {
			return true
		}

		if player.PersonaTag == personaTag {
			playerId = id
			found = true
			return false // Stop searching once we find the player
		}

		return true
	})

	if !found {
		return 0, fmt.Errorf("player not found")
	}

	return playerId, nil
}

// GetPlayerByPersonaTag gets a Player with the given PersonaTag.
//
// Code Flow:
// 1. Create a search query for Players.
// 2. Iterate over the search results and check if the PersonaTag matches the specified PersonaTag.
// 3. Return the Player if a match is found, or an error if not.
func GetPlayerByPersonaTag(world cardinal.WorldContext, personaTag string) (*Player, error) {
	q := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[Player]()),
	)

	var player *Player
	var err error

	q.Each(world, func(id types.EntityID) bool {
		player, err = cardinal.GetComponent[Player](world, id)
		if err != nil {
			return true
		}

		if player.PersonaTag == personaTag {
			return false // Stop searching once we find the player
		}

		return true
	})

	if player == nil {
		return nil, fmt.Errorf("player not found")
	}

	return player, nil
}
