// Package system contains the logic for handling player spawning actions.
package system

import (
	"fmt"

	"tamagotchi/component"
	"tamagotchi/game"
	"tamagotchi/msg"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

/**
 * Constants for player spawning.
 */
const (
	InitialHP           = 100
	MaxPersonaTagLength = 32
)

/**
 * Function Flow:
 * 1. The `PlayerSpawnerAction` function is called, which processes `CreatePlayer` transactions.
 * 2. For each transaction, the function checks if a player with the provided persona tag already exists.
 * 3. If the player does not exist, the function creates a new player entity with the provided persona tag.
 * 4. The function then emits a "new_player" event to notify other systems of the new player.
 *
 * PlayerSpawnerAction spawns players based on `CreatePlayer` transactions.
 * This provides an example of a system that creates a new entity.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the player spawning action.
 */
func PlayerSpawnerAction(world cardinal.WorldContext) error {

	return cardinal.EachMessage[msg.CreatePlayerMsg, msg.CreatePlayerReply](
		world,
		func(create cardinal.TxData[msg.CreatePlayerMsg]) (msg.CreatePlayerReply, error) {

			if create.Tx.PersonaTag == "" {
				return msg.CreatePlayerReply{}, fmt.Errorf("persona tag cannot be empty")
			}

			// Check if persona tag exceeds maximum length
			if len(create.Tx.PersonaTag) > MaxPersonaTagLength {
				return msg.CreatePlayerReply{}, fmt.Errorf("persona tag exceeds maximum length of %d characters", MaxPersonaTagLength)
			}

			// Step 2: Check if player already exists
			//   - Retrieve the player's ID using their persona tag
			//   - If the player is found, return an error indicating that the player already exists
			_, err := component.FindPlayerByPersonaTag(world, create.Tx.PersonaTag)
			if err == nil {
				// Player already exists, return an error
				return msg.CreatePlayerReply{}, fmt.Errorf("player already exists")
			}

			// Step 3: Create a new player entity
			//   - Use the `Create` function to create a new player entity with the provided persona tag
			//   - Initialize the player's properties, such as pets, items, and money
			id, err := cardinal.Create(world,
				component.Player{
					PersonaTag: create.Tx.PersonaTag,
					Pets:       make([]types.EntityID, 0),
					Items:      make([]types.EntityID, 0),
					Money:      game.PlayerInitialMoney,
				},
			)
			if err != nil {
				// Error creating player, return an error
				return msg.CreatePlayerReply{}, fmt.Errorf("error creating player: %w", err)
			}

			// Step 4: Emit a "new_player" event
			//   - Use the `EmitEvent` function to emit a "new_player" event with the new player's ID
			//   - If the event emission fails, return an error
			err = world.EmitEvent(map[string]any{
				"event": "new_player",
				"id":    id,
			})
			if err != nil {
				// Error emitting event, return an error
				return msg.CreatePlayerReply{}, err
			}

			// Return a successful result
			return msg.CreatePlayerReply{Success: true}, nil
		})
}
