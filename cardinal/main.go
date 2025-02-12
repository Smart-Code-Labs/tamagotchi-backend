package main

import (
	"errors"

	"github.com/rs/zerolog/log"
	"pkg.world.dev/world-engine/cardinal"

	"tamagotchi/component"
	"tamagotchi/msg"
	"tamagotchi/query"
	"tamagotchi/system"
)

func main() {
	w, err := cardinal.NewWorld(cardinal.WithDisableSignatureVerification())
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	MustInitWorld(w)

	Must(w.StartGame())
}

// MustInitWorld registers all components, messages, queries, and systems. This initialization happens in a helper
// function so that this can be used directly in tests.
func MustInitWorld(w *cardinal.World) {
	// Register components
	Must(
		cardinal.RegisterComponent[component.Pet](w),
		cardinal.RegisterComponent[component.Dna](w),
		cardinal.RegisterComponent[component.Health](w),
		cardinal.RegisterComponent[component.Energy](w),
		cardinal.RegisterComponent[component.Hygiene](w),
		cardinal.RegisterComponent[component.Wellness](w),
		cardinal.RegisterComponent[component.Activity](w),
		cardinal.RegisterComponent[component.Think](w),
	)

	// Register messages (user action)
	Must(
		cardinal.RegisterMessage[msg.CreatePetMsg, msg.CreatePetResult](w, "create-pet"),
		cardinal.RegisterMessage[msg.PlayPetMsg, msg.PlayPetMsgReply](w, "play-pet"),
		cardinal.RegisterMessage[msg.SleepPetMsg, msg.SleepPetMsgReply](w, "sleep-pet"),
		cardinal.RegisterMessage[msg.BathPetMsg, msg.BathPetMsgReply](w, "bath-pet"),
		cardinal.RegisterMessage[msg.EatPetMsg, msg.EatPetMsgReply](w, "eat-pet"),
	)

	// Register queries
	Must(
		cardinal.RegisterQuery[query.PetHealthRequest, query.PetHealthResponse](w, "pet-health", query.PetHealth),
		cardinal.RegisterQuery[query.PetEnergyRequest, query.PetEnergyResponse](w, "pet-energy", query.PetEnergy),
		cardinal.RegisterQuery[query.CurrentTickMsg, query.CurrentTickReply](w, "current-tick", query.CurrentTick),
		cardinal.RegisterQuery[query.PetsMsg, query.PetsReply](w, "pets-list", query.Pets),
	)

	// Each system executes deterministically in the order they are added.
	// This is a neat feature that can be strategically used for systems that depends on the order of execution.
	// For example, you may want to run the attack system before the regen system
	// so that the player's HP is subtracted (and player killed if it reaches 0) before HP is regenerated.
	Must(cardinal.RegisterSystems(w,
		system.PlayAction,
		system.PetSpawnerAction,
		system.BathAction,
		system.SleepAction,
		system.EatAction,
		system.EnergyDeclineSystem,
		system.HygieneDeclineSystem,
		system.WellnessDeclineSystem,
		system.HealthDeclineSystem,
		system.ActivityDeclineSystem,
		system.ThinkSystem,
	))
}

func Must(err ...error) {
	e := errors.Join(err...)
	if e != nil {
		log.Fatal().Err(e).Msg("")
	}
}
