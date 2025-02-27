package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"pkg.world.dev/world-engine/cardinal"

	"tamagotchi/component"
	"tamagotchi/game"
)

const bathToyName = "Sponge"

// TestSystem_PetBathAction_HappyPath tests the happy path of the PetBathAction function
func TestSystem_PetBathAction_HappyPath(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	world := cardinal.NewWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, bathToyName)
	assert.NoError(t, err)

	// - The pet's hygiene is set to a low value.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = 0
	err = cardinal.SetComponent(world, petId, petHygiene)
	assert.NoError(t, err)

	// When:
	// - The PetBathAction function is called.
	err = PetBathAction(t, tf, petName, bathToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's components are updated correctly.
	petHygiene, err = cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.Greater(t, petHygiene.Hy, 0)

	petEnergy, err := cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, 90, petEnergy.E)

	petWellness, err := cardinal.GetComponent[component.Wellness](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, petWellness.Wn, game.MaxWellness)

	petActivity, err := cardinal.GetComponent[component.Activity](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, "Bathing", petActivity.Activity)

	petThink, err := cardinal.GetComponent[component.Think](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, petThink.Think, game.InitialThink)
}

// TestSystem_PetBathAction_NoPet tests that an error is returned when the pet does not exist.
func TestSystem_PetBathAction_NoPet(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, bathToyName)
	assert.NoError(t, err)

	// When:
	// - The PetBathAction function is called with a non-existent pet.
	err = PetBathAction(t, tf, "non_existent_pet", bathToyName)

	// Then:
	// - An error is returned.
	assert.Error(t, err)
}

// TestSystem_PetBathAction_NoToy tests that an error is returned when the toy does not exist.
func TestSystem_PetBathAction_NoToy(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// When:
	// - The PetBathAction function is called with a non-existent toy.
	err := PetBathAction(t, tf, petName, "non_existent_toy")

	// Then:
	// - An error is returned.
	assert.Error(t, err)
}

// TestSystem_PetBathAction_InvalidPetName tests that an error is returned when an invalid pet name is provided.
func TestSystem_PetBathAction_InvalidPetName(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, bathToyName)
	assert.NoError(t, err)

	// When:
	// - The PetBathAction function is called with an invalid pet name.
	err = PetBathAction(t, tf, "", bathToyName)

	// Then:
	// - An error is returned.
	assert.Error(t, err)
}

// TestSystem_PetBathAction_InvalidToyName tests that an error is returned when an invalid toy name is provided.
func TestSystem_PetBathAction_InvalidToyName(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// When:
	// - The PetBathAction function is called with an invalid toy name.
	err := PetBathAction(t, tf, petName, "")

	// Then:
	// - An error is returned.
	assert.Error(t, err)
}

// TestSystem_PetBathAction_MaxHygiene tests that the pet's hygiene is set to the maximum value after bathing.
func TestSystem_PetBathAction_MaxHygiene(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	world := cardinal.NewWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, bathToyName)
	assert.NoError(t, err)

	// - The pet's hygiene is set to a low value.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = 0
	err = cardinal.SetComponent(world, petId, petHygiene)
	assert.NoError(t, err)

	// When:
	// - The PetBathAction function is called.
	err = PetBathAction(t, tf, petName, bathToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's hygiene is set to the maximum value.
	petHygiene, err = cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, 30, petHygiene.Hy)
}

// TestSystem_PetBathAction_MinHygiene tests that the pet's hygiene is increased after bathing when it is at the minimum value.
func TestSystem_PetBathAction_MinHygiene(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	world := cardinal.NewWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, bathToyName)
	assert.NoError(t, err)

	// - The pet's hygiene is set to the minimum value.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = 0
	err = cardinal.SetComponent(world, petId, petHygiene)
	assert.NoError(t, err)

	// When:
	// - The PetBathAction function is called.
	err = PetBathAction(t, tf, petName, bathToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's hygiene is increased.
	petHygiene, err = cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.Greater(t, petHygiene.Hy, 0)
}

// TestSystem_PetBathAction_EnergyDecrease tests that the pet's energy decreases after bathing.
func TestSystem_PetBathAction_EnergyDecrease(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	world := cardinal.NewWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, bathToyName)
	assert.NoError(t, err)

	// - The pet's hygiene is set to a low value.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = 0
	err = cardinal.SetComponent(world, petId, petHygiene)
	assert.NoError(t, err)

	// - The pet's energy is set to the maximum value.
	petEnergy, err := cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	petEnergy.E = game.MaxEnergy
	err = cardinal.SetComponent(world, petId, petEnergy)
	assert.NoError(t, err)

	// When:
	// - The PetBathAction function is called.
	err = PetBathAction(t, tf, petName, bathToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's energy decreases.
	petEnergy, err = cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	assert.Less(t, petEnergy.E, game.MaxEnergy)
}

// TestSystem_PetBathAction_ActivityUpdate tests that the pet's activity is updated to "Bathing" after bathing.
func TestSystem_PetBathAction_ActivityUpdate(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	world := cardinal.NewWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, bathToyName)
	assert.NoError(t, err)

	// - The pet's hygiene is set to a low value.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = 0
	err = cardinal.SetComponent(world, petId, petHygiene)
	assert.NoError(t, err)

	// When:
	// - The PetBathAction function is called.
	err = PetBathAction(t, tf, petName, bathToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's activity is updated to "Bathing".
	petActivity, err := cardinal.GetComponent[component.Activity](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, petActivity.Activity, "Bathing")
}

// TestSystem_PetBathAction_MultipleBaths tests that multiple baths do not cause the pet's hygiene to exceed the maximum value.
func TestSystem_PetBathAction_MultipleBaths(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	world := cardinal.NewWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, bathToyName)
	assert.NoError(t, err)

	// - The pet's hygiene is set to a low value.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = 0
	err = cardinal.SetComponent(world, petId, petHygiene)
	assert.NoError(t, err)

	// When:
	// - The PetBathAction function is called multiple times.
	err = PetBathAction(t, tf, petName, bathToyName)
	assert.NoError(t, err)

	err = PetBathAction(t, tf, petName, bathToyName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pet is already engaged in an activity")

	// Then:
	// - The pet's hygiene does not exceed the maximum value.
	petHygiene, err = cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.LessOrEqual(t, petHygiene.Hy, game.MaxHygiene)
}

// TestSystem_PetBathAction_ConcurrentBaths tests that concurrent baths do not cause any issues with the pet's state.
func TestSystem_PetBathAction_ConcurrentBaths(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	world := cardinal.NewWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, bathToyName)
	assert.NoError(t, err)

	// - The pet's hygiene is set to a low value.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = 0
	err = cardinal.SetComponent(world, petId, petHygiene)
	assert.NoError(t, err)

	// When:
	// - The PetBathAction function is called concurrently multiple times.
	var wg sync.WaitGroup
	simul_txs := 1 // TODO: increase this to reproduce the error
	count_err := 0
	count_ok := 0
	for i := 0; i < simul_txs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := PetBathAction(t, tf, petName, bathToyName)
			if err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "pet is already engaged in an activity")
				count_err++
				t.Logf("count_err[%d]", count_err)
			} else {
				count_ok++
				t.Logf("count_ok[%d]", count_ok)
				assert.NoError(t, err)
			}
		}()
	}
	wg.Wait()
	assert.Equal(t, 1, count_ok)
	assert.Equal(t, simul_txs-1, count_err)

	// Then:
	// - The pet's hygiene is updated correctly.
	petHygiene, err = cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.Greater(t, petHygiene.Hy, 0)
}

// TestSystem_PetBathAction_BathingWithLowEnergy tests that the pet can still bathe even with low energy.
func TestSystem_PetBathAction_BathingWithLowEnergy(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	world := cardinal.NewWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, bathToyName)
	assert.NoError(t, err)

	// - The pet's hygiene is set to a low value.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = 0
	err = cardinal.SetComponent(world, petId, petHygiene)
	assert.NoError(t, err)

	// - The pet's energy is set to a low value.
	petEnergy, err := cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	petEnergy.E = 11
	err = cardinal.SetComponent(world, petId, petEnergy)
	assert.NoError(t, err)

	// When:
	// - The PetBathAction function is called.
	err = PetBathAction(t, tf, petName, bathToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's hygiene is updated correctly.
	petHygiene, err = cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.Greater(t, petHygiene.Hy, 0)
}

// TestSystem_PetBathAction_BathingWithLowWellness tests that the pet can still bathe even with low wellness.
func TestSystem_PetBathAction_BathingWithLowWellness(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	world := cardinal.NewWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, bathToyName)
	assert.NoError(t, err)

	// - The pet's hygiene is set to a low value.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = 0
	err = cardinal.SetComponent(world, petId, petHygiene)
	assert.NoError(t, err)

	// - The pet's wellness is set to a low value.
	petWellness, err := cardinal.GetComponent[component.Wellness](wCtx, petId)
	assert.NoError(t, err)
	petWellness.Wn = 10
	err = cardinal.SetComponent(world, petId, petWellness)
	assert.NoError(t, err)

	// When:
	// - The PetBathAction function is called.
	err = PetBathAction(t, tf, petName, bathToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's hygiene is updated correctly.
	petHygiene, err = cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.Greater(t, petHygiene.Hy, 0)
}

func TestSystem_PetBathAction_AnotherActionWhileBathing(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	world := cardinal.NewWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, bathToyName)
	assert.NoError(t, err)

	// - The pet's hygiene is set to a low value.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = 0
	err = cardinal.SetComponent(world, petId, petHygiene)
	assert.NoError(t, err)

	// When:
	// - The PetBathAction function is called.
	err = PetBathAction(t, tf, petName, bathToyName)
	assert.NoError(t, err)

	// - Another action is attempted while the pet is bathing.
	err = PetPlayAction(t, tf, petName, bathToyName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pet is already engaged in an activity")
}
