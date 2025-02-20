package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
	"tamagotchi/msg"
	system "tamagotchi/system"
)

// PetSpawnerAction spawns pets based on `Create-pet` transactions.
// This provides an example of a system that creates a new entity.
func BreedAction(world cardinal.WorldContext) error {
	rng := world.Rand()
	var pet *comp.Pet
	var father *comp.Pet
	var mother *comp.Pet
	var err error
	var found bool = false
	log := world.Logger()
	// rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return cardinal.EachMessage[msg.BreedPetMsg, msg.BreedPetMsgReply](
		world,
		func(create cardinal.TxData[msg.BreedPetMsg]) (msg.BreedPetMsgReply, error) {

			q := cardinal.NewSearch().Entity(filter.Contains(filter.Component[comp.Pet]()))
			q.Each(world, func(petId types.EntityID) bool {
				pet, err = cardinal.GetComponent[comp.Pet](world, petId)
				if err != nil {
					return true
				}
				if pet.Nickname == create.Msg.BornName {
					found = true
				}
				return true
			})
			log.Info().Msgf("Checking: n[%t]", found)
			if found {
				return msg.BreedPetMsgReply{}, fmt.Errorf("error creating pet: Name already exist")
			}

			fatherId, err := system.QueryPetIdByName(world, create.Msg.FatherName)
			if err != nil {
				return msg.BreedPetMsgReply{}, fmt.Errorf("error creating pet: [get father EntityID]: %w", err)
			}

			father, err = cardinal.GetComponent[comp.Pet](world, fatherId)
			if err != nil {
				return msg.BreedPetMsgReply{}, fmt.Errorf("error creating pet [get Father]: %w", err)
			}

			motherId, err := system.QueryPetIdByName(world, create.Msg.MotherName)
			if err != nil {
				return msg.BreedPetMsgReply{}, fmt.Errorf("error creating pet: [get mother EntityID]: %w", err)
			}

			mother, err = cardinal.GetComponent[comp.Pet](world, motherId)
			if err != nil {
				return msg.BreedPetMsgReply{}, fmt.Errorf("error creating pet [get Mother]: %w", err)
			}

			if mother.Gender == father.Gender {
				return msg.BreedPetMsgReply{}, fmt.Errorf("error creating pet [Mother and father have the same Gender]: %w", err)
			}

			if mother.PersonaTag != create.Tx.PersonaTag {
				return msg.BreedPetMsgReply{}, fmt.Errorf("error creating pet [You are not the owner of Mother]: %w", err)
			}

			if father.PersonaTag != create.Tx.PersonaTag {
				return msg.BreedPetMsgReply{}, fmt.Errorf("error creating pet [You are not the owner of father]: %w", err)
			}

			// if father.Level < 5 || mother.Level < 5 {
			// 	return msg.BreedPetMsgReply{}, fmt.Errorf("error creating pet [Mother or Father not grown enought]: %w", err)
			// }

			// Assign "water", "fire", or "earth" based on the random number
			var element string
			switch rng.Intn(4) {
			case 0:
				element = "wynd"
			case 1:
				element = "water"
			case 2:
				element = "fire"
			case 3:
				element = "earth"
			}

			// Assign "water", "fire", or "earth" based on the random number
			var skill string
			switch rng.Intn(3) {
			case 0:
				skill = "Intellect"
			case 1:
				skill = "Force"
			case 2:
				skill = "skilled"
			}

			id, err := cardinal.Create(world,
				comp.Pet{PersonaTag: create.Tx.PersonaTag, Nickname: create.Msg.BornName, Level: 0, XP: 0, NextLevelXP: 0},
				comp.Health{HP: constants.InitialHP},
				comp.Energy{E: constants.InitialE},
				comp.Hygiene{Hy: constants.InitialHy},
				comp.Wellness{Wn: constants.InitialWn},
				comp.Dna{ // use Mother and father DNA here
					A: rng.Intn(100),
					C: rng.Intn(100),
					G: rng.Intn(100),
					T: rng.Intn(100),
				},
				comp.Activity{Activity: "None", CountDown: 0},
				comp.Think{Think: "..."},
				comp.Magic{Kind: element, Level: 0, XP: 0, NextLevelXP: 0},
				comp.Skill{Kind: skill, Level: 0, XP: 0, NextLevelXP: 0},
			)
			if err != nil {
				return msg.BreedPetMsgReply{}, fmt.Errorf("error creating pet: %w", err)
			}

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
