package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
	"tamagotchi/msg"
)

// PlayAction reduce pet's Energy.
// PlayAction reduce pet's Hygiene.
func PlayAction(world cardinal.WorldContext) error {
	var petEnergy *comp.Energy
	var petHygiene *comp.Hygiene
	var petWellness *comp.Wellness
	var petActivity *comp.Activity
	var petThink *comp.Think

	return cardinal.EachMessage(
		world,
		func(play cardinal.TxData[msg.PlayPetMsg]) (msg.PlayPetMsgReply, error) {

			petId, err := queryPetIdByName(world, play.Msg.TargetNickname)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get EntityID]: %w", err)
			}

			// check if not activity
			petActivity, err = cardinal.GetComponent[comp.Activity](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Activity]: %w", err)
			}

			if petActivity.Duration > 0 {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [already on Activity]: %w", err)
			}

			// set activity
			petActivity.Activity = "Playing"
			petActivity.Duration = constants.PlayDuration

			petThink, err = cardinal.GetComponent[comp.Think](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Think]: %w", err)
			}

			petThink.Think = "Love to play!"

			// get Energy
			petEnergy, err = cardinal.GetComponent[comp.Energy](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Energy]: %w", err)
			}

			// reduce Energy
			if petEnergy.E-constants.EnergyReduce > 0 {
				petEnergy.E -= constants.EnergyReduce
			} else {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [not enought energy]: %w", err)
			}

			// get Hygiene
			petHygiene, err = cardinal.GetComponent[comp.Hygiene](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Hygiene]: %w", err)
			}

			// reduce Hygiene
			if petHygiene.Hy-constants.HygieneReduce > 0 {
				petHygiene.Hy -= constants.HygieneReduce
			} else {
				petHygiene.Hy = 0
			}

			// get Wellness
			petWellness, err = cardinal.GetComponent[comp.Wellness](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Wellness]: %w", err)
			}

			// increase wellness
			if petWellness.Wn+constants.WellnessIncrease <= 100 {
				petWellness.Wn += constants.WellnessIncrease
			} else {
				petWellness.Wn = 100
			}

			if err := cardinal.SetComponent(world, petId, petEnergy); err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [set Energy]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petHygiene); err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [set Hygiene]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petWellness); err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [set Wellness]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [set Activity]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petThink); err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [set Think]: %w", err)
			}

			return msg.PlayPetMsgReply{
				Energy:   constants.EnergyReduce,
				Hygiene:  constants.HygieneReduce,
				Wellness: constants.WellnessIncrease,
				Activity: petActivity.Activity,
				Duration: petActivity.Duration,
			}, nil
		})
}
