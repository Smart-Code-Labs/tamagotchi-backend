package main

import (
	"testing"

	"gotest.tools/v3/assert"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/receipt"
	"pkg.world.dev/world-engine/cardinal/types"

	"tamagotchi/component"
	constants "tamagotchi/game"
	"tamagotchi/msg"
)

const (
	playMsgName   = "game.play-pet"
	createMsgName = "game.create-pet"
	sleepMsgName  = "game.sleep-pet"
	bathMsgName   = "game.bath-pet"
	eatMsgname    = "game.eat-pet"
	breedMsgName  = "game.breed-pet"
)

// TestSystem_PlaySystem_ErrorWhenTargetDoesNotExist ensures the play message results in an error when the given
// target does not exist. Note, message errors are stored in receipts; they are NOT returned from the relevant system.
func TestSystem_PlaySystem_ErrorWhenTargetDoesNotExist(t *testing.T) {
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	txHash := tf.AddTransaction(getPlayMsgID(t, tf.World), msg.PlayPetMsg{
		TargetNickname: "does-not-exist",
	})

	tf.DoTick()

	gotReceipt := getReceiptFromPastTick(t, tf.World, txHash)
	if len(gotReceipt.Errs) == 0 {
		t.Fatal("expected error when target does not exist")
	}
}

// TestSystem_PetSpawnerSystem_CanCreatePet ensures the Createpet message can be used to create a new pet
// with the default amount of health. cardinal.NewSearch is used to find the newly created pet.
func TestSystem_PetSpawnerSystem_CanCreatePet(t *testing.T) {
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	const nickname = "Manny"
	createTxHash := tf.AddTransaction(getCreateMsgID(t, tf.World), msg.CreatePetMsg{
		Nickname: nickname,
	})
	tf.DoTick()

	// Make sure the pet creation was successful
	createReceipt := getReceiptFromPastTick(t, tf.World, createTxHash)
	if errs := createReceipt.Errs; len(errs) > 0 {
		t.Fatalf("expected 0 errors when creating a pet, got %v", errs)
	}

	// Make sure the newly created pet has 100 health
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)

	acc := make([]types.EntityID, 0)
	t.Log(len(acc))
	err := cardinal.NewSearch().Entity(filter.Contains(filter.Component[component.Pet]())).
		Each(wCtx, func(id types.EntityID) bool {
			t.Log(id)
			pet, err := cardinal.GetComponent[component.Pet](wCtx, id)
			if err != nil {
				t.Fatalf("failed to get pet component: %v", err)
			}
			if pet.Nickname == nickname {
				acc = append(acc, id)
				return false
			}
			return true
		})
	t.Log(len(acc))

	assert.NilError(t, err)
	assert.Equal(t, len(acc), 1)
	id := acc[0]

	health, err := cardinal.GetComponent[component.Health](wCtx, id)
	if err != nil {
		t.Fatalf("failed to find entity ID: %v", err)
	}
	if health.HP != 100 {
		t.Fatalf("a newly created pet should have 100 health; got %v", health.HP)
	}
}

// TestSystem_PlaySystem_PlayingTargetReducesTheirEnergy ensures an play message can find an existing target the
// reduce the target's health.
func TestSystem_PlaySystem_PlayingTargetReducesTheirEnergy(t *testing.T) {
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	const target = "Manny"

	// Create an initial pet
	_ = tf.AddTransaction(getCreateMsgID(t, tf.World), msg.CreatePetMsg{
		Nickname: target,
	})
	tf.DoTick()

	// Play the pet
	playTxHash := tf.AddTransaction(getPlayMsgID(t, tf.World), msg.PlayPetMsg{
		TargetNickname: target,
	})
	tf.DoTick()

	// Make sure play was successful
	playReceipt := getReceiptFromPastTick(t, tf.World, playTxHash)
	if errs := playReceipt.Errs; len(errs) > 0 {
		t.Fatalf("expected no errors when playing a pet; got %v", errs)
	}

	// Find the played pet and check their health.
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	var found bool
	searchErr := cardinal.NewSearch().Entity(filter.Contains(filter.Component[component.Pet]())).
		Each(wCtx, func(id types.EntityID) bool {
			pet, err := cardinal.GetComponent[component.Pet](wCtx, id)
			if err != nil {
				t.Fatalf("failed to get pet component for %v", id)
			}
			if pet.Nickname != target {
				return true
			}
			// The pet's nickname matches the target. This is the pet we care about.
			found = true
			energy, err := cardinal.GetComponent[component.Energy](wCtx, id)
			if err != nil {
				t.Fatalf("failed to get health component for %v", id)
			}
			// The target started with 100 E, -10 for the play, -1 for energy decline
			expectedEnergy := constants.InitialE - constants.EnergyReduce - 1
			if energy.E != expectedEnergy {
				t.Fatalf("play target should end up with %v hp, got %v", expectedEnergy, energy.E)
			}

			return false
		})
	if searchErr != nil {
		t.Fatalf("error when performing search: %v", searchErr)
	}
	if !found {
		t.Fatalf("failed to find target %q", target)
	}
}

func getCreateMsgID(t *testing.T, world *cardinal.World) types.MessageID {
	return getMsgID(t, world, createMsgName)
}

func getPlayMsgID(t *testing.T, world *cardinal.World) types.MessageID {
	return getMsgID(t, world, playMsgName)
}

func getMsgID(t *testing.T, world *cardinal.World, fullName string) types.MessageID {
	msg, ok := world.GetMessageByFullName(fullName)
	if !ok {
		t.Fatalf("failed to get %q message", fullName)
	}
	return msg.ID()
}

// getReceiptFromPastTick search past ticks for a txHash that matches the given txHash. An error will be returned if
// the txHash cannot be found in Cardinal's history.
func getReceiptFromPastTick(t *testing.T, world *cardinal.World, txHash types.TxHash) receipt.Receipt {
	tick := world.CurrentTick()
	for {
		tick--
		receipts, err := world.GetTransactionReceiptsForTick(tick)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range receipts {
			if r.TxHash == txHash {
				return r
			}
		}
	}
}
