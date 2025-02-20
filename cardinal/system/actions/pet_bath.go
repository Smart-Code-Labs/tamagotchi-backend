package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
	"tamagotchi/msg"
	system "tamagotchi/system"
)

// BathAction reduce pet's Hygiene.
// BathAction reduce pet's Hygiene.
func BathAction(world cardinal.WorldContext) error {
	var petHygiene *comp.Hygiene
	var petActivity *comp.Activity
	var petThink *comp.Think

	return cardinal.EachMessage(
		world,
		func(Bath cardinal.TxData[msg.BathPetMsg]) (msg.BathPetMsgReply, error) {

			petId, err := system.QueryPetIdByName(world, Bath.Msg.TargetNickname)
			if err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [get EntityID]: %w", err)
			}

			// check if not activity
			petActivity, err = cardinal.GetComponent[comp.Activity](world, petId)
			if err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [get Activity]: %w", err)
			}

			if petActivity.CountDown > 0 {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [already on Activity]: %w", err)
			}

			petThink, err = cardinal.GetComponent[comp.Think](world, petId)
			if err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [get Think]: %w", err)
			}

			petThink.Think = constants.ThinkBath

			// set activity
			petActivity.Activity = "Bathing"
			petActivity.CountDown = constants.TickHour
			petActivity.TotalTicks = constants.TickHour
			petActivity.Percentage = 100

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
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [set Activity]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
				return msg.BathPetMsgReply{}, fmt.Errorf("failed to Bath [set Activity]: %w", err)
			}

			return msg.BathPetMsgReply{
				Hygiene:  constants.HygieneIncrease,
				Activity: petActivity.Activity,
				Duration: petActivity.CountDown}, nil
		})
}
