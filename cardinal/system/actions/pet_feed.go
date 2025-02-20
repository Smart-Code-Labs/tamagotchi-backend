package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
	"tamagotchi/msg"
	system "tamagotchi/system"
)

// FeedAction reduce pet's Health.
// FeedAction reduce pet's Hygiene.
func FeedAction(world cardinal.WorldContext) error {
	var petHealth *comp.Health
	var petActivity *comp.Activity
	var petThink *comp.Think

	return cardinal.EachMessage(
		world,
		func(Eat cardinal.TxData[msg.FeedPetMsg]) (msg.FeedPetMsgReply, error) {

			petId, err := system.QueryPetIdByName(world, Eat.Msg.TargetNickname)
			if err != nil {
				return msg.FeedPetMsgReply{}, fmt.Errorf("failed to Eat [get EntityID]: %w", err)
			}

			// check if not activity
			petActivity, err = cardinal.GetComponent[comp.Activity](world, petId)
			if err != nil {
				return msg.FeedPetMsgReply{}, fmt.Errorf("failed to Eat [get Activity]: %w", err)
			}

			if petActivity.CountDown > 0 {
				return msg.FeedPetMsgReply{}, fmt.Errorf("failed to Eat [already on Activity]: %w", err)
			}

			petThink, err = cardinal.GetComponent[comp.Think](world, petId)
			if err != nil {
				return msg.FeedPetMsgReply{}, fmt.Errorf("failed to Eat [get Think]: %w", err)
			}

			petThink.Think = constants.ThinkEat

			// set activity
			petActivity.Activity = "Eating"
			petActivity.CountDown = constants.TickHour
			petActivity.TotalTicks = constants.TickHour
			petActivity.Percentage = 100

			// get Health
			petHealth, err = cardinal.GetComponent[comp.Health](world, petId)
			if err != nil {
				return msg.FeedPetMsgReply{}, fmt.Errorf("failed to Eat [get Health]: %w", err)
			}

			// increase Health
			if petHealth.HP+constants.HealthIncrease <= 100 {
				petHealth.HP += constants.HealthIncrease
			} else {
				petHealth.HP = 100
			}

			if err := cardinal.SetComponent(world, petId, petHealth); err != nil {
				return msg.FeedPetMsgReply{}, fmt.Errorf("failed to Eat [set Health]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
				return msg.FeedPetMsgReply{}, fmt.Errorf("failed to Eat [set Activity]: %w", err)
			}

			if err := cardinal.SetComponent(world, petId, petActivity); err != nil {
				return msg.FeedPetMsgReply{}, fmt.Errorf("failed to Eat [set Activity]: %w", err)
			}

			return msg.FeedPetMsgReply{
				Health:   constants.HealthIncrease,
				Activity: petActivity.Activity,
				Duration: petActivity.CountDown}, nil
		})
}
