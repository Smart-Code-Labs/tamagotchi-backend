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
	var pet *comp.Pet
	var player *comp.Player
	var item *comp.Wellness

	log := world.Logger()
	return cardinal.EachMessage(
		world,
		func(play cardinal.TxData[msg.PlayPetMsg]) (msg.PlayPetMsgReply, error) {

			// Check if player exists
			playerID, err := comp.FindPlayerByPersonaTag(world, play.Tx.PersonaTag)
			if err != nil {
				return msg.PlayPetMsgReply{}, err
			}

			player, err = cardinal.GetComponent[comp.Player](world, playerID)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Player]: %w", err)
			}

			petId, err := player.GetPetNickname(world, play.Msg.TargetNickname)
			if err != nil {
				return msg.PlayPetMsgReply{}, err
			}

			itemId, err := player.GetItemByName(world, play.Msg.ItemName)
			if err != nil {
				return msg.PlayPetMsgReply{}, err
			}

			// get pet lvl
			pet, err = cardinal.GetComponent[comp.Pet](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Pet]: %w", err)
			}

			// add experince and calculate lvl
			pet.AddXP(constants.ExperienceEarn)

			// check if not activity
			petActivity, err = cardinal.GetComponent[comp.Activity](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Activity]: %w", err)
			}

			if petActivity.CountDown > 0 {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [already on Activity]: %w", err)
			}

			petThink, err = cardinal.GetComponent[comp.Think](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Think]: %w", err)
			}

			petThink.Think = constants.ThinkPlay

			// set activity
			petActivity.Activity = "Playing"
			petActivity.CountDown = constants.TickHour
			petActivity.TotalTicks = constants.TickHour
			petActivity.Percentage = 100

			// get Energy
			petEnergy, err = cardinal.GetComponent[comp.Energy](world, petId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Energy]: %w", err)
			}

			// reduce Energy
			log.Info().Msgf("Playing: Energy Before[%d]", petEnergy.E)

			if petEnergy.E-constants.EnergyReduce > 0 {
				petEnergy.E -= constants.EnergyReduce
			} else {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [not enought energy]: %w", err)
			}

			log.Info().Msgf("Playing: Energy After[%d]", petEnergy.E)

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

			item, err = cardinal.GetComponent[comp.Wellness](world, itemId)
			if err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [get Item Wellness]: %w", err)
			}

			// increase wellness
			if petWellness.Wn+item.Wn <= 100 {
				petWellness.Wn += item.Wn
			} else {
				petWellness.Wn = 100
			}

			// increase experience
			if err := cardinal.SetComponent(world, petId, pet); err != nil {
				return msg.PlayPetMsgReply{}, fmt.Errorf("failed to play [set Experince]: %w", err)
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
				Duration: petActivity.CountDown,
			}, nil
		})
}
