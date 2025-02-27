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
		cardinal.RegisterMessage[msg.CreatePlayerMsg, msg.CreatePlayerReply](w, "create-player"),
		cardinal.RegisterMessage[msg.CreatePetMsg, msg.CreatePetReply](w, "create-pet"),
		cardinal.RegisterMessage[msg.CurePetMsg, msg.CurePetMsgReply](w, "cure-pet"),
		cardinal.RegisterMessage[msg.PlayPetMsg, msg.PlayPetMsgReply](w, "play-pet"),
		cardinal.RegisterMessage[msg.SleepPetMsg, msg.SleepPetMsgReply](w, "sleep-pet"),
		cardinal.RegisterMessage[msg.BathPetMsg, msg.BathPetMsgReply](w, "bath-pet"),
		cardinal.RegisterMessage[msg.FeedPetMsg, msg.FeedPetMsgReply](w, "feed-pet"),
		cardinal.RegisterMessage[msg.BreedPetMsg, msg.BreedPetMsgReply](w, "breed-pet"),
		cardinal.RegisterMessage[msg.ButItemMsg, msg.BuyItemMsgReply](w, "buy-item"),
	)

	// Register queries
	Must(
		cardinal.RegisterQuery[query.PetHealthRequest, query.PetHealthResponse](w, "pet-health", query.QueryPetHealth),
		cardinal.RegisterQuery[query.PetEnergyRequest, query.PetEnergyResponse](w, "pet-energy", query.QueryPetEnergy),
		cardinal.RegisterQuery[query.CurrentTickMsg, query.CurrentTickReply](w, "current-tick", query.QueryCurrentTick),
		cardinal.RegisterQuery[query.PetsMsg, query.PetsReply](w, "pets-list", query.GamePets),
		cardinal.RegisterQuery[query.ToysMsg, query.ToysReply](w, "toystore-list", query.QueryToyStore),
		cardinal.RegisterQuery[query.FoodsMsg, query.FoodsReply](w, "foodstore-list", query.QueryFoodStore),
		cardinal.RegisterQuery[query.DrugsMsg, query.DrugsReply](w, "drugstore-list", query.QueryDrugStore),
		cardinal.RegisterQuery[query.ItemListMsg, query.ItemListReply](w, "personaItem-list", query.QueryPlayerItems),
		cardinal.RegisterQuery[query.PlayerExistMsg, query.PlayerExistReply](w, "player-exist", query.QueryPlayerExist),
		cardinal.RegisterQuery[query.LeaderboardMsg, query.LeaderboardReply](w, "leaderboard", query.QueryLeaderboard),
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
		actions.PetPlayAction,
		actions.PetCureAction,
		actions.PetBathAction,
		actions.PetSleepAction,
		actions.PetFeedAction,
		actions.PetBreedAction,
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
