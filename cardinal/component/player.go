package component

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

type Player struct {
	PersonaTag string           `json:"personaTag"`
	Pets       []types.EntityID `json:"pets"`
	Items      []types.EntityID `json:"items"`
	Money      float64          `json:"money"`
}

func (Player) Name() string {
	return "Player"
}

// HasPet checks if the given EntityID exists in the Player's Pets slice.
// Returns true if the ItemId is found, false otherwise.
func (p Player) HasPet(id types.EntityID) bool {
	for _, petID := range p.Pets {
		if petID == id {
			return true
		}
	}
	return false
}

// HasPetNickname checks if the Player has a Pet with the specified nickname
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
	return 0, nil
}

// HasItem checks if the given EntityID exists in the Player's Items slice.
// Returns true if the ItemId is found, false otherwise.
func (p Player) HasItem(id types.EntityID) bool {
	for _, itemID := range p.Items {
		if itemID == id {
			return true
		}
	}
	return false
}

func (p Player) GetItemByName(world cardinal.WorldContext, itemName string) (types.EntityID, error) {
	for _, itemID := range p.Items {
		item, err := cardinal.GetComponent[Item](world, itemID)
		if err != nil {
			return 0, err
		}
		if item.ItemName == itemName {
			return itemID, nil
		}
	}
	return 0, nil
}

// AddPet adds a new Pet ID to the Player's Pets array.
func AddPet(world cardinal.WorldContext, playerID types.EntityID, petID types.EntityID) error {
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

// AddItem adds a new Item ID to the Player's Items array.
func AddItem(world cardinal.WorldContext, playerID types.EntityID, itemID types.EntityID) error {
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

// AddItem adds a new Item ID to the Player's Items array.
func ReducePlayerMoney(world cardinal.WorldContext, playerID types.EntityID, quantity float64) error {
	// Append the new item to the player's Items array
	player, err := cardinal.GetComponent[Player](world, playerID)
	if err != nil {
		return fmt.Errorf("error adding pet to player: %w", err)
	}

	player.Money -= quantity
	err = cardinal.SetComponent(world, playerID, player)
	if err != nil {
		return fmt.Errorf("error updating player money: %w", err)
	}
	return nil
}

// AddItem adds a new Item ID to the Player's Items array.
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

// findPlayerByPersonaTag searches for a player with the given PersonaTag and returns their EntityID
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
		return 0, fmt.Errorf("player with PersonaTag %s does not exist", personaTag)
	}

	return playerId, nil
}
