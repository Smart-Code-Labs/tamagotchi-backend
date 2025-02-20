package system

import (
	"fmt"

	comp "tamagotchi/component"
	"tamagotchi/msg"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

const (
	InitialHP = 100
)

// PlayerSpawnerAction spawns players based on `CreatePlayer` transactions.
// This provides an example of a system that creates a new entity.
func PlayerSpawnerAction(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.CreatePlayerMsg, msg.CreatePlayerResult](
		world,
		func(create cardinal.TxData[msg.CreatePlayerMsg]) (msg.CreatePlayerResult, error) {

			_, err := comp.FindPlayerByPersonaTag(world, create.Tx.PersonaTag)
			if err == nil {
				return msg.CreatePlayerResult{}, fmt.Errorf("PLAYER ALREADY EXIST")
			}

			id, err := cardinal.Create(world,
				comp.Player{
					PersonaTag: create.Tx.PersonaTag,
					Pets:       make([]types.EntityID, 0),
					Items:      make([]types.EntityID, 0),
					Money:      1000,
				},
			)
			if err != nil {
				return msg.CreatePlayerResult{}, fmt.Errorf("error creating player: %w", err)
			}

			err = world.EmitEvent(map[string]any{
				"event": "new_player",
				"id":    id,
			})
			if err != nil {
				return msg.CreatePlayerResult{}, err
			}
			return msg.CreatePlayerResult{Success: true}, nil
		})
}
