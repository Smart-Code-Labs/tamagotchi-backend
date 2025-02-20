package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
	"tamagotchi/msg"
	system "tamagotchi/system"
)

// SleepAction reduce pet's Energy.
// SleepAction reduce pet's Hygiene.
func SleepAction(world cardinal.WorldContext) error {
	var petEnergy *comp.Energy
	var petActivity *comp.Activity
	var petThink *comp.Think

	return cardinal.EachMessage(
		world,
		func(Sleep cardinal.TxData[msg.SleepPetMsg]) (msg.SleepPetMsgReply, error) {

			petId, err := system.QueryPetIdByName(world, Sleep.Msg.TargetNickname)
			if err != nil {
				return msg.SleepPetMsgReply{}, fmt.Errorf("failed to Sleep [get EntityID]: %w", err)
			}

			// check if not activity
			petActivity, err = cardinal.GetComponent[comp.Activity](world, petId)
			if err != nil {
				return msg.SleepPetMsgReply{}, fmt.Errorf("failed to Sleep [get Activity]: %w", err)
			}

			if petActivity.CountDown > 0 {
				return msg.SleepPetMsgReply{}, fmt.Errorf("failed to Sleep [already on Activity]: %w", err)
			}

			petThink, err = cardinal.GetComponent[comp.Think](world, petId)
			if err != nil {
				return msg.SleepPetMsgReply{}, fmt.Errorf("failed to Sleep [get Think]: %w", err)
			}

			petThink.Think = constants.ThinkSleep

			// set activity
			petActivity.Activity = "Sleeping"
			petActivity.CountDown = constants.TickEightHours
			petActivity.TotalTicks = constants.TickEightHours
			petActivity.Percentage = 100

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
				return msg.SleepPetMsgReply{}, fmt.Errorf("failed to Sleep [set Activity]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
				return msg.SleepPetMsgReply{}, fmt.Errorf("failed to Sleep [set Activity]: %w", err)
			}

			return msg.SleepPetMsgReply{
				Energy:   constants.EnergyReduce,
				Activity: petActivity.Activity,
				Duration: petActivity.CountDown}, nil
		})
}
