package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
	"tamagotchi/msg"
)

// EatAction reduce pet's Health.
// EatAction reduce pet's Hygiene.
func EatAction(world cardinal.WorldContext) error {
	var petHealth *comp.Health
	var petActivity *comp.Activity

	return cardinal.EachMessage(
		world,
		func(Eat cardinal.TxData[msg.EatPetMsg]) (msg.EatPetMsgReply, error) {

			petId, err := queryPetIdByName(world, Eat.Msg.TargetNickname)
			if err != nil {
				return msg.EatPetMsgReply{}, fmt.Errorf("failed to Eat [get EntityID]: %w", err)
			}

			// check if not activity
			petActivity, err = cardinal.GetComponent[comp.Activity](world, petId)
			if err != nil {
				return msg.EatPetMsgReply{}, fmt.Errorf("failed to play [get Activity]: %w", err)
			}

			if petActivity.Duration > 0 {
				return msg.EatPetMsgReply{}, fmt.Errorf("failed to play [already on Activity]: %w", err)
			}

			// set activity
			petActivity.Activity = "Eating"
			petActivity.Duration = constants.EatDuration

			// get Health
			petHealth, err = cardinal.GetComponent[comp.Health](world, petId)
			if err != nil {
				return msg.EatPetMsgReply{}, fmt.Errorf("failed to Eat [get Health]: %w", err)
			}

			// increase Health
			if petHealth.HP+constants.HealthIncrease <= 100 {
				petHealth.HP += constants.HealthIncrease
			} else {
				petHealth.HP = 100
			}

			if err := cardinal.SetComponent(world, petId, petHealth); err != nil {
				return msg.EatPetMsgReply{}, fmt.Errorf("failed to Eat [set Health]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
				return msg.EatPetMsgReply{}, fmt.Errorf("failed to play [set Activity]: %w", err)
			}

			return msg.EatPetMsgReply{
				Health:   constants.HealthIncrease,
				Activity: petActivity.Activity,
				Duration: petActivity.Duration}, nil
		})
}
