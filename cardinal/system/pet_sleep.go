package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
	"tamagotchi/msg"
)

// SleepAction reduce pet's Energy.
// SleepAction reduce pet's Hygiene.
func SleepAction(world cardinal.WorldContext) error {
	var petEnergy *comp.Energy
	var petActivity *comp.Activity

	return cardinal.EachMessage(
		world,
		func(Sleep cardinal.TxData[msg.SleepPetMsg]) (msg.SleepPetMsgReply, error) {

			petId, err := queryPetIdByName(world, Sleep.Msg.TargetNickname)
			if err != nil {
				return msg.SleepPetMsgReply{}, fmt.Errorf("failed to Sleep [get EntityID]: %w", err)
			}

			// check if not activity
			petActivity, err = cardinal.GetComponent[comp.Activity](world, petId)
			if err != nil {
				return msg.SleepPetMsgReply{}, fmt.Errorf("failed to play [get Activity]: %w", err)
			}

			if petActivity.Duration > 0 {
				return msg.SleepPetMsgReply{}, fmt.Errorf("failed to play [already on Activity]: %w", err)
			}

			// set activity
			petActivity.Activity = "Sleeping"
			petActivity.Duration = constants.SleepDuration

			// get Energy
			petEnergy, err = cardinal.GetComponent[comp.Energy](world, petId)
			if err != nil {
				return msg.SleepPetMsgReply{}, fmt.Errorf("failed to Sleep [get Energy]: %w", err)
			}

			// increase Energy
			if petEnergy.E+constants.EnergyIncrease <= 100 {
				petEnergy.E += constants.EnergyIncrease
			} else {
				petEnergy.E = 100
			}

			if err := cardinal.SetComponent(world, petId, petEnergy); err != nil {
				return msg.SleepPetMsgReply{}, fmt.Errorf("failed to Sleep [set Energy]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
				return msg.SleepPetMsgReply{}, fmt.Errorf("failed to play [set Activity]: %w", err)
			}

			return msg.SleepPetMsgReply{
				Energy:   constants.EnergyReduce,
				Activity: petActivity.Activity,
				Duration: petActivity.Duration}, nil
		})
}
