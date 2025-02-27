// Package system contains the logic for handling pet spawning actions.
package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	"tamagotchi/component"
	"tamagotchi/game"
	"tamagotchi/msg"
)

/**
 * Function Flow:
 * 1. The `PetSpawnerAction` function is called, which logs a message indicating that it has been called.
 * 2. The function utilizes the `cardinal.EachMessage` function to process `CreatePetMsg` transactions.
 * 3. For each transaction, the function checks if the player exists by their persona tag.
 * 4. If the player exists, the function checks if a pet with the provided nickname already exists.
 * 5. If the nickname is unique, the function creates a new pet and assigns it to the player.
 * 6. The function logs a message indicating that a new pet has been created and emits a "new_pet" event.
 *
 * PetSpawnerAction spawns pets based on `Create-pet` transactions.
 * This provides an example of a system that creates a new entity.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the pet spawning action.
 */
func PetSpawnerAction(world cardinal.WorldContext) error {
	log := world.Logger()

	return cardinal.EachMessage[msg.CreatePetMsg, msg.CreatePetReply](
		world,
		func(create cardinal.TxData[msg.CreatePetMsg]) (msg.CreatePetReply, error) {

			log.Info().Msgf("PetSpawnerAction -> PersonaTag [%v]", create.Tx.PersonaTag)

			if create.Tx.PersonaTag == "" {
				return msg.CreatePetReply{}, fmt.Errorf("persona tag cannot be null")
			}
			err := create.Msg.Validate()
			if err != nil {
				return msg.CreatePetReply{}, err
			}
			// Step 3: Check if player exists
			//   - Retrieve the player's ID using their persona tag
			//   - If the player is not found, log an error and return an error
			playerID, err := component.FindPlayerByPersonaTag(world, create.Tx.PersonaTag)
			if err != nil {
				return msg.CreatePetReply{}, err
			}

			// Step 4: Check if pet nickname already exists
			//   - Use the `CheckPetNicknameExists` function to check if a pet with the provided nickname exists
			//   - If the nickname is not unique, return an error
			found, _, err := component.QueryPetIdByName(world, create.Msg.Nickname)

			if err != nil {
				return msg.CreatePetReply{}, err
			}

			if found {
				return msg.CreatePetReply{}, fmt.Errorf("error creating pet: Name already exist")
			}

			// Step 5: Create a new pet and assign it to the player
			//   - Use the `CreateRandomPet` function to create a new pet
			//   - Add the pet to the player's pets
			petID, err := component.CreateRandomPet(world, create.Tx.PersonaTag, create.Msg.Nickname)
			if err != nil {
				return msg.CreatePetReply{}, err
			}

			if err := component.ReducePlayerMoney(world, playerID, game.PetCost); err != nil {
				return msg.CreatePetReply{}, err
			}

			// Add pet to player's pets
			err = component.AddPlayerPet(world, playerID, petID)
			if err != nil {
				return msg.CreatePetReply{}, err
			}

			// Step 6: Emit a "new_pet" event
			//   - Use the `EmitEvent` function to emit a "new_pet" event
			//   - If the event emission fails, log an error and return an error
			err = world.EmitEvent(map[string]any{
				"event": "new_pet",
				"id":    petID,
			})
			if err != nil {
				log.Error().Msgf("Failed to emit new_pet event for pet %s (ID: %v): %v", create.Msg.Nickname, petID, err)
				return msg.CreatePetReply{}, err
			}

			// Return a successful result
			return msg.CreatePetReply{Success: true}, nil
		})
}
