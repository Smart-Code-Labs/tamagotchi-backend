package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	comp "tamagotchi/component"
	"tamagotchi/msg"
)

const (
	InitialHP = 100
	InitialE  = 100
	InitialHy = 100
	InitialWn = 100
)

// PetSpawnerAction spawns pets based on `Create-pet` transactions.
// This provides an example of a system that creates a new entity.
func PetSpawnerAction(world cardinal.WorldContext) error {
	rng := world.Rand()
	// rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return cardinal.EachMessage[msg.CreatePetMsg, msg.CreatePetResult](
		world,
		func(create cardinal.TxData[msg.CreatePetMsg]) (msg.CreatePetResult, error) {
			id, err := cardinal.Create(world,
				comp.Pet{Nickname: create.Msg.Nickname},
				comp.Health{HP: InitialHP},
				comp.Energy{E: InitialE},
				comp.Hygiene{Hy: InitialHy},
				comp.Wellness{Wn: InitialWn},
				comp.Dna{
					A: rng.Intn(100),
					C: rng.Intn(100),
					G: rng.Intn(100),
					T: rng.Intn(100),
				},
				comp.Activity{Activity: "None", Duration: 0},
				comp.Think{Think: "..."},
			)
			if err != nil {
				return msg.CreatePetResult{}, fmt.Errorf("error creating pet: %w", err)
			}

			err = world.EmitEvent(map[string]any{
				"event": "new_pet",
				"id":    id,
			})
			if err != nil {
				return msg.CreatePetResult{}, err
			}
			return msg.CreatePetResult{Success: true}, nil
		})
}
