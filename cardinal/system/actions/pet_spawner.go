package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	comp "tamagotchi/component"
	"tamagotchi/msg"
)

// PetSpawnerAction spawns pets based on `Create-pet` transactions.
// This provides an example of a system that creates a new entity.
func PetSpawnerAction(world cardinal.WorldContext) error {
	log := world.Logger()
	log.Info().Msgf("called")
	return cardinal.EachMessage[msg.CreatePetMsg, msg.CreatePetResult](
		world,
		func(create cardinal.TxData[msg.CreatePetMsg]) (msg.CreatePetResult, error) {

			// Check if player exists
			playerID, err := comp.FindPlayerByPersonaTag(world, create.Tx.PersonaTag)
			if err != nil {
				log.Error().Msgf("Player with PersonaTag %s not found: %v", create.Tx.PersonaTag, err)
				return msg.CreatePetResult{}, fmt.Errorf("error creating pet: %w", err)
			}

			found, err := comp.CheckPetNicknameExists(world, create.Msg.Nickname)
			if found {
				return msg.CreatePetResult{}, err
			}

			petID, err := comp.CreateRandomPet(world, create.Tx.PersonaTag, create.Msg.Nickname)
			if err != nil {
				return msg.CreatePetResult{}, err
			}
			log.Info().Msgf("Created: Pet[%d]", petID)

			err = comp.AddPet(world, playerID, petID)
			if err != nil {
				return msg.CreatePetResult{}, err
			}

			err = world.EmitEvent(map[string]any{
				"event": "new_pet",
				"id":    petID,
			})
			if err != nil {
				log.Error().Msgf("Failed to emit new_pet event for pet %s (ID: %v): %v", create.Msg.Nickname, petID, err)
				return msg.CreatePetResult{}, err
			}
			return msg.CreatePetResult{Success: true}, nil
		})
}
