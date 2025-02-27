package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"pkg.world.dev/world-engine/cardinal"
	persona_msg "pkg.world.dev/world-engine/cardinal/persona/msg"
	"pkg.world.dev/world-engine/cardinal/receipt"
	"pkg.world.dev/world-engine/cardinal/types"
	"pkg.world.dev/world-engine/sign"

	"tamagotchi/component"
	"tamagotchi/msg"
)

const (
	createMsgName        = "game.create-pet"
	createPlayerMsgName  = "game.create-player"
	createPersonaMsgName = "persona.create-persona"
	buyItemMsgName       = "game.buy-item"
	playMsgName          = "game.play-pet"
	sleepMsgName         = "game.sleep-pet"
	bathMsgName          = "game.bath-pet"
	eatMsgName           = "game.eat-pet"
	breedMsgName         = "game.breed-pet"
	personaTag           = "_test_persona"
	signerAddress        = "0xa1D239A61908FaC55Ca95Cd112698623bD36bC4f"
	petName              = "Manny"
)

// This function tests the creation of a pet.
// Flow:
// 1. Initialize a new test fixture.
// 2. Create a persona.
// 3. Create a player associated with the persona.
// 4. Create a pet that belongs to the player.
// 5. Verify that the pet was created successfully.
func TestSystem_PetSpawnerSystem_CanCreatePet(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	var err error

	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// When:
	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// Then:
	// - The pet is verified to exist.
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)

	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	assert.NotZero(t, petId)

	// - The newly created pet should have 100 health.
	health, err := cardinal.GetComponent[component.Health](wCtx, petId)
	assert.NoError(t, err)
	assert.EqualValues(t, health.HP, int(100))

	// Verify that the pet's ID is in the player's Pets array.
	player, err := component.GetPlayerByPersonaTag(wCtx, personaTag)
	assert.NoError(t, err)

	playerPetId, err := player.GetPetNickname(wCtx, petName)
	assert.NoError(t, err)
	assert.Equal(t, playerPetId, petId, "Expected pet ID %v to be in the player's Pets array, but got %v", petId, playerPetId)
}

// TODO: Test 2: Duplicate Pet Nickname Sent: Create pet with duplicate nickname Received: Error: Pet nickname already exists Verify: No new pet created, error logged, and no "new_pet" event emitted

func executeTx[T any](t *testing.T, tf *cardinal.TestFixture, msgName string, req interface{}, personaTag string) (*T, error) {
	msgFullName, ok := tf.World.GetMessageByFullName(msgName)
	if !ok {
		t.Fatalf("failed to get %q message", msgName)
	}
	TxHash := tf.AddTransaction(msgFullName.ID(), req, &sign.Transaction{PersonaTag: personaTag})
	tf.DoTick()

	// Then:
	// - The pet's energy is verified to be reduced.
	t.Log("Before getReceiptFromPastTick")
	response := getReceiptFromPastTick(t, tf.World, TxHash)
	t.Log("After getReceiptFromPastTick")
	if errs := response.Errs; len(errs) > 0 {
		return nil, errs[0]
	}
	t.Log("RESPONSE; ", response)
	resp, ok := response.Result.(T)
	if !ok {
		t.Fatalf("failed to cast for %v", response)
	}
	t.Log("casted: ", response)
	return &resp, nil
}

// This function gets the receipt from a past tick.
// Flow:
// 1. Search past ticks for the transaction hash.
// 2. Return the receipt for the transaction.
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

// This function creates a pet.
// Flow:
// 1. Create a new pet message.
// 2. Add the transaction to the test fixture.
// 3. Verify that the pet was created successfully.
func createPet(t *testing.T, tf *cardinal.TestFixture, petName string, personaTag string) error {
	// Preconditions:
	// - The test fixture is initialized.
	createMsg := msg.CreatePetMsg{
		Nickname: petName,
	}
	_, err := executeTx[msg.CreatePetReply](t, tf, createMsgName, createMsg, personaTag)
	return err
}

// This function creates a player.
// Flow:
// 1. Get the message type for creating a player.
// 2. Add the transaction to the test fixture.
// 3. Verify that the player was created successfully.
func createPlayer(t *testing.T, tf *cardinal.TestFixture, personaTag string) error {
	// Preconditions:
	// - The test fixture is initialized.
	createPlayerMsg := msg.CreatePlayerMsg{}
	_, err := executeTx[msg.CreatePlayerReply](t, tf, createPlayerMsgName, createPlayerMsg, personaTag)
	return err
}

// This function creates a persona.
// Flow:
// 1. Get the message type for creating a persona.
// 2. Add the transaction to the test fixture.
// 3. Verify that the persona was created successfully.
func createPersona(t *testing.T, tf *cardinal.TestFixture, personaTag string) error {
	// Preconditions:
	// - The test fixture is initialized.
	createPersonaMsg := persona_msg.CreatePersona{
		PersonaTag:    personaTag,
		SignerAddress: signerAddress,
	}
	_, err := executeTx[persona_msg.CreatePersonaResult](t, tf, createPersonaMsgName, createPersonaMsg, personaTag)
	return err
}

// This function buys a toy.
// Flow:
// 1. Get the message type for buying a toy.
// 2. Add the transaction to the test fixture.
// 3. Verify that the toy was bought successfully.
func buyToy(t *testing.T, tf *cardinal.TestFixture, toyName string) error {
	// Preconditions:
	// - The test fixture is initialized.
	t.Log("buyToy")
	buyItemMsg := msg.ButItemMsg{
		Name: toyName,
	}
	_, err := executeTx[msg.BuyItemMsgReply](t, tf, buyItemMsgName, buyItemMsg, personaTag)
	return err
}

// This function buys a toy.
// Flow:
// 1. Get the message type for buying a toy.
// 2. Add the transaction to the test fixture.
// 3. Verify that the toy was bought successfully.
func PetPlayAction(t *testing.T, tf *cardinal.TestFixture, nickName string, itemName string) error {
	// Preconditions:
	// - The test fixture is initialized.
	petPlayMsg := msg.PlayPetMsg{
		TargetNickname: nickName,
		ItemName:       itemName,
	}
	_, err := executeTx[msg.PlayPetMsgReply](t, tf, playMsgName, petPlayMsg, personaTag)
	return err
}

// This function buys a toy.
// Flow:
// 1. Get the message type for buying a toy.
// 2. Add the transaction to the test fixture.
// 3. Verify that the toy was bought successfully.
func PetBathAction(t *testing.T, tf *cardinal.TestFixture, nickName string, itemName string) error {
	// Preconditions:
	// - The test fixture is initialized.
	petBathMsg := msg.BathPetMsg{
		TargetNickname: nickName,
		ItemName:       itemName,
	}
	_, err := executeTx[msg.BathPetMsgReply](t, tf, bathMsgName, petBathMsg, personaTag)
	return err
}
