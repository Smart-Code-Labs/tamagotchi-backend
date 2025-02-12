package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
	"tamagotchi/msg"
)

// BathAction reduce pet's Hygiene.
// BathAction reduce pet's Hygiene.
func BathAction(world cardinal.WorldContext) error {
	var petHygiene *comp.Hygiene
	var petActivity *comp.Activity

	return cardinal.EachMessage(
		world,
		func(Bath cardinal.TxData[msg.BathPetMsg]) (msg.BathPetMsgReply, error) {

			petId, err := queryPetIdByName(world, Bath.Msg.TargetNickname)
			if err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [get EntityID]: %w", err)
			}

			// check if not activity
			petActivity, err = cardinal.GetComponent[comp.Activity](world, petId)
			if err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to play [get Activity]: %w", err)
			}

			if petActivity.Duration > 0 {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to play [already on Activity]: %w", err)
			}

			// set activity
			petActivity.Activity = "Bathing"
			petActivity.Duration = constants.BathDuration

			// get Hygiene
			petHygiene, err = cardinal.GetComponent[comp.Hygiene](world, petId)
			if err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [get Hygiene]: %w", err)
			}

			// increase Hygiene
			if petHygiene.Hy+constants.HygieneIncrease <= 100 {
				petHygiene.Hy += constants.HygieneIncrease
			} else {
				petHygiene.Hy = 100
			}

			if err := cardinal.SetComponent(world, petId, petHygiene); err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [set Hygiene]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to play [set Activity]: %w", err)
			}

			return msg.BathPetMsgReply{
				Hygiene:  constants.HygieneIncrease,
				Activity: petActivity.Activity,
				Duration: petActivity.Duration}, nil
		})
}
