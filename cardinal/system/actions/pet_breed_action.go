// Package system contains the logic for handling bath pet actions.
package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	"tamagotchi/component"
	"tamagotchi/game"
	"tamagotchi/msg"
)

/**
PetBreedAction Function Flow:
1. Spawn pets based on 'Create-pet' transactions.
2. Iterate through each message of type BreedPetMsg.
3. Perform sanity checks:
   - Check if the pet name already exists.
   - Get the father and mother pets.
   - Check if the mother and father are of different genders.
   - Check if the persona is the owner of the mother and father pets.
4. Create a new pet entity with initial characteristics:
   - Generate random element and skill.
   - Create a new entity with Pet, Health, Energy, Hygiene, Wellness, Dna, Activity, Think, Magic, and Skill components.
5. Emit a 'new_pet' event with the new pet's ID.
*/
// PetBreedAction spawns pets based on `Create-pet` transactions.
// This provides an example of a system that creates a new entity.
func PetBreedAction(world cardinal.WorldContext) error {
	// 1. Spawn pets based on 'Create-pet' transactions.
	rng := world.Rand()
	// rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return cardinal.EachMessage[msg.BreedPetMsg, msg.BreedPetMsgReply](
		world,
		func(create cardinal.TxData[msg.BreedPetMsg]) (msg.BreedPetMsgReply, error) {
			// 3. Perform sanity checks:
			//    - Check if the pet name already exists.
			if _, _, err := component.QueryPetIdByName(world, create.Msg.BornName); err != nil {
				return msg.BreedPetMsgReply{}, err
			}

			//    - Get the father and mother pets.
			_, father, err := component.GetPetByNickname(world, create.Msg.FatherName)
			if err != nil {
				return msg.BreedPetMsgReply{}, err
			}

			_, mother, err := component.GetPetByNickname(world, create.Msg.MotherName)
			if err != nil {
				return msg.BreedPetMsgReply{}, err
			}

			// Check Pets lvl disable for video
			// Checks Pets gender disable for video
			//    - Check if the mother and father are of different genders.
			// if err := CheckPetsAreDifferentGenders(world, father, mother); err != nil {
			// 	return msg.BreedPetMsgReply{}, err
			// }

			//    - Check if the persona is the owner of the mother and father pets.
			if err := CheckOwnerOfPets(world, create.Tx.PersonaTag, father, mother); err != nil {
				return msg.BreedPetMsgReply{}, err
			}

			// 4. Create a new pet entity with initial characteristics:
			//    - Generate random element and skill.
			element := game.Elements[rng.Intn(len(game.Elements))]
			skill := game.Skills[rng.Intn(len(game.Skills))]

			//    - Create a new entity with Pet, Health, Energy, Hygiene, Wellness, Dna, Activity, Think, Magic, and Skill components.
			id, err := cardinal.Create(world,
				component.Pet{PersonaTag: create.Tx.PersonaTag, Nickname: create.Msg.BornName, Level: 0, XP: 0, NextLevelXP: 0},
				component.Health{HP: game.MaxHP},
				component.Energy{E: game.MaxEnergy},
				component.Hygiene{Hy: game.MaxHygiene},
				component.Wellness{Wn: game.MaxWellness},
				component.Dna{ // use Mother and father DNA here
					A: rng.Intn(100),
					C: rng.Intn(100),
					G: rng.Intn(100),
					T: rng.Intn(100),
				},
				component.Activity{Activity: "None", CountDown: 0},
				component.Think{Think: "..."},
				component.Magic{Kind: element, Level: 0, XP: 0, NextLevelXP: 0},
				component.Skill{Kind: skill, Level: 0, XP: 0, NextLevelXP: 0},
			)
			if err != nil {
				return msg.BreedPetMsgReply{}, fmt.Errorf("error creating pet: %w", err)
			}

			// 5. Emit a 'new_pet' event with the new pet's ID.
			err = world.EmitEvent(map[string]any{
				"event": "new_pet",
				"id":    id,
			})
			if err != nil {
				return msg.BreedPetMsgReply{}, err
			}
			return msg.BreedPetMsgReply{Success: true}, nil
		})
}

/**
CheckPetsAreDifferentGenders Function Flow:
1. Compare the genders of the two pets.
2. If the genders are the same, return an error indicating the pets have the same gender.
*/
// CheckPetsAreDifferentGenders checks if the given two pets are of different genders.
// It returns an error if they are of the same gender.
func CheckPetsAreDifferentGenders(world cardinal.WorldContext, father *component.Pet, mother *component.Pet) error {
	// 1. Compare the genders of the two pets.
	if mother.Gender == father.Gender {
		// 2. If the genders are the same, return an error indicating the pets have the same gender.
		return fmt.Errorf("error creating pet [Mother and father have the same Gender]")
	}
	return nil
}

/**
CheckOwnerOfPets Function Flow:
1. Compare the persona tag of the given pets with the given persona tag.
2. If the persona tag does not match for either pet, return an error indicating the persona is not the owner.
*/
// CheckOwnerOfPets checks if the given persona is the owner of the given pets.
// It returns an error if the persona is not the owner of either pet.
func CheckOwnerOfPets(world cardinal.WorldContext, personaTag string, father *component.Pet, mother *component.Pet) error {
	// 1. Compare the persona tag of the given pets with the given persona tag.
	if mother.PersonaTag != personaTag {
		// 2. If the persona tag does not match for either pet, return an error indicating the persona is not the owner.
		return fmt.Errorf("error creating pet [You are not the owner of Mother]")
	}
	if father.PersonaTag != personaTag {
		return fmt.Errorf("error creating pet [You are not the owner of father]")
	}
	return nil
}
