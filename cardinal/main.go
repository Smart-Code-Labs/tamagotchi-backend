package main

import (
	"errors"

	"github.com/rs/zerolog/log"
	"pkg.world.dev/world-engine/cardinal"

	"tamagotchi/component"
	"tamagotchi/msg"
	"tamagotchi/query"
	actions "tamagotchi/system/actions"
	game "tamagotchi/system/game"
	mechanics "tamagotchi/system/mechanics"
)

func main() {
	//tickChannel := cardinal.WithTickChannel(time.Tick(constants.TickRate * time.Second)) // 1 tick every 5 seconds
	disableSignature := cardinal.WithDisableSignatureVerification()
	receiptHistory := cardinal.WithReceiptHistorySize(1000)
	messageExpiration := cardinal.WithMessageExpiration(60)
	w, err := cardinal.NewWorld(disableSignature, receiptHistory, messageExpiration)
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
		cardinal.RegisterComponent[component.Leaderboard](w),
		cardinal.RegisterComponent[component.ToyStore](w),
		cardinal.RegisterComponent[component.DrugStore](w),
		cardinal.RegisterComponent[component.FoodStore](w),
		cardinal.RegisterComponent[component.Player](w),
		cardinal.RegisterComponent[component.Pet](w),
		cardinal.RegisterComponent[component.Item](w),
		cardinal.RegisterComponent[component.Dna](w),
		cardinal.RegisterComponent[component.Health](w),
		cardinal.RegisterComponent[component.Energy](w),
		cardinal.RegisterComponent[component.Hygiene](w),
		cardinal.RegisterComponent[component.Wellness](w),
		cardinal.RegisterComponent[component.Activity](w),
		cardinal.RegisterComponent[component.Think](w),
		cardinal.RegisterComponent[component.Magic](w),
		cardinal.RegisterComponent[component.Skill](w),
	)

	// Register messages (user action)
	Must(
		cardinal.RegisterMessage[msg.CreatePlayerMsg, msg.CreatePlayerResult](w, "create-player"),
		cardinal.RegisterMessage[msg.CreatePetMsg, msg.CreatePetResult](w, "create-pet"),
		cardinal.RegisterMessage[msg.PlayPetMsg, msg.PlayPetMsgReply](w, "play-pet"),
		cardinal.RegisterMessage[msg.SleepPetMsg, msg.SleepPetMsgReply](w, "sleep-pet"),
		cardinal.RegisterMessage[msg.BathPetMsg, msg.BathPetMsgReply](w, "bath-pet"),
		cardinal.RegisterMessage[msg.FeedPetMsg, msg.FeedPetMsgReply](w, "feed-pet"),
		cardinal.RegisterMessage[msg.BreedPetMsg, msg.BreedPetMsgReply](w, "breed-pet"),
		cardinal.RegisterMessage[msg.ButItemMsg, msg.BuyItemMsgReply](w, "buy-item"),
	)

	// Register queries
	Must(
		cardinal.RegisterQuery[query.PetHealthRequest, query.PetHealthResponse](w, "pet-health", query.PetHealth),
		cardinal.RegisterQuery[query.PetEnergyRequest, query.PetEnergyResponse](w, "pet-energy", query.PetEnergy),
		cardinal.RegisterQuery[query.CurrentTickMsg, query.CurrentTickReply](w, "current-tick", query.CurrentTick),
		cardinal.RegisterQuery[query.PetsMsg, query.PetsReply](w, "pets-list", query.Pets),
		cardinal.RegisterQuery[query.ToysMsg, query.ToysReply](w, "toystore-list", query.Toys),
		cardinal.RegisterQuery[query.FoodsMsg, query.FoodsReply](w, "foodstore-list", query.Foods),
		cardinal.RegisterQuery[query.DrugsMsg, query.DrugsReply](w, "drugstore-list", query.Drugs),
		cardinal.RegisterQuery[query.ItemListMsg, query.ItemListReply](w, "personaitem-list", query.PlayerItemList),
	)

	// Each system executes deterministically in the order they are added.
	// This is a neat feature that can be strategically used for systems that depends on the order of execution.
	// For example, you may want to run the attack system before the regen system
	// so that the player's HP is subtracted (and player killed if it reaches 0) before HP is regenerated.
	Must(cardinal.RegisterSystems(w,
		game.LeaderboardSystem,
		// Create Actors
		actions.PetSpawnerAction,
		actions.PlayerSpawnerAction,
		// Execute Actions
		actions.PlayAction,
		actions.BathAction,
		actions.SleepAction,
		actions.FeedAction,
		actions.BreedAction,
		actions.BuyItemAction,
		// Execute Game mechanics
		mechanics.EnergyDeclineSystem,
		mechanics.HygieneDeclineSystem,
		mechanics.WellnessDeclineSystem,
		mechanics.HealthDeclineSystem,
		mechanics.ActivityDeclineSystem,
		mechanics.ThinkSystem,
	))

	Must(cardinal.RegisterInitSystems(w,
		game.SpawnDefaultSystem,
	))
}

func Must(err ...error) {
	e := errors.Join(err...)
	if e != nil {
		log.Fatal().Err(e).Msg("")
	}
}
