package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"pkg.world.dev/world-engine/cardinal"

	"tamagotchi/component"
	"tamagotchi/game"
	"tamagotchi/query"
)

const playToyName = "Ball"

// TestSystem_PetPlayAction_HappyPath tests the happy path of the PetPlayAction function
func TestSystem_PetPlayAction_HappyPath(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	// When:
	// - The PetPlayAction function is called.
	err = PetPlayAction(t, tf, petName, playToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's components are updated correctly.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	petEnergy, err := cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	assert.Less(t, petEnergy.E, game.MaxEnergy)

	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.Less(t, petHygiene.Hy, game.MaxHygiene)

	petWellness, err := cardinal.GetComponent[component.Wellness](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, petWellness.Wn, game.MaxWellness)

	petActivity, err := cardinal.GetComponent[component.Activity](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, "Playing", petActivity.Activity)
}

// This function tests playing with a pet.
// Flow:
// 1. Initialize a new test fixture.
// 2. Create a persona.
// 3. Create a player associated with the persona.
// 4. Create a pet that belongs to the player.
// 5. Buy a toy to play with the pet.
// 6. Play with the pet using the toy.
// 7. Verify that the pet's energy was reduced.
func TestSystem_PlaySystem_PlayingTargetReducesTheirEnergy(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// When:
	// - A toy is bought to play with the pet.
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	toyList := query.GetAllItemsFromToyStore(wCtx)

	for _, toy := range toyList {
		t.Log(toy.ItemName)
	}

	err := buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	// - The pet is played with using the toy.
	err = PetPlayAction(t, tf, petName, playToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's energy is verified to be reduced.

	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	if err != nil || petId == 0 {
		t.Fatalf("failed to get pet component for %v", petId)
	}

	energy, err := component.GetPetEnergy(wCtx, petId)
	if err != nil {
		t.Fatalf("failed to get energy component for %v", petId)
	}
	// The target started with 100 E, -10 for the play
	expectedEnergy := game.MaxEnergy - game.EnergyReduce
	if energy.E != expectedEnergy {
		t.Fatalf("play target should end up with %v hp, got %v", expectedEnergy, energy.E)
	}

	// Asserts:
	// - No errors should occur when playing the pet.
	// - The pet's energy should be reduced by the expected amount.
}

// Test non-existent player
func TestSystem_PetPlayAction_NonExistentPlayer(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// When:
	// - The PetPlayAction function is called with a non-existent player ID.
	err := PetPlayAction(t, tf, "non-existent-pet", "toy")
	assert.Error(t, err)

	// Then:
	// - The function returns an error message indicating that the player doesn't exist.
	assert.Contains(t, err.Error(), "player not found")
}

// Test non-existent pet
func TestSystem_PetPlayAction_NonExistentPet(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// When:
	// - The PetPlayAction function is called with a non-existent pet ID.
	err := PetPlayAction(t, tf, "non-existent-pet", "toy")
	assert.Error(t, err)

	// Then:
	// - The function returns an error message indicating that the pet doesn't exist.
	assert.Contains(t, err.Error(), "pet not found")
}

// Test pet already engaged in activity
func TestSystem_PetPlayAction_PetEngagedInActivity(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	err := buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	// - The pet is played with using the toy.
	err = PetPlayAction(t, tf, petName, playToyName)
	assert.NoError(t, err)

	// When:
	// - The PetPlayAction function is called with the pet that is already engaged in an activity.
	err = buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	// - The pet is played with using the toy.
	err = PetPlayAction(t, tf, petName, playToyName)
	assert.Error(t, err)

	// Then:
	// - The function returns an error message indicating that the pet is already engaged in an activity.
	assert.Contains(t, err.Error(), "pet is already engaged in an activity")
}

// Test pet experience and level update
func TestSystem_PetPlayAction_PetExperienceAndLevelUpdate(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buy a toy.
	err := buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	_, pet, err := component.GetPetByNickname(wCtx, petName)
	assert.NoError(t, err)
	assert.Equal(t, pet.XP, int64(0))
	assert.Equal(t, pet.Level, int64(0))

	// When:
	// - The PetPlayAction function is called with the pet.
	err = PetPlayAction(t, tf, petName, playToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's experience and level are updated correctly.
	_, pet, err = component.GetPetByNickname(wCtx, petName)
	assert.NoError(t, err)
	assert.Greater(t, pet.XP, int64(0))
	assert.Greater(t, pet.Level, int64(0))
}

// Test pet energy, hygiene, and wellness update
func TestSystem_PetPlayAction_PetEnergyHygieneAndWellnessUpdate(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The player buys a toy.
	err := buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	// - The pet's energy, hygiene, and wellness are set to initial values.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	assert.NotEqual(t, petId, 0)

	petEnergy, err := cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, petEnergy.E, game.MaxEnergy)

	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, petHygiene.Hy, game.MaxHygiene)

	petWellness, err := cardinal.GetComponent[component.Wellness](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, petWellness.Wn, game.MaxWellness)

	// When:
	// - The PetPlayAction function is called with the pet.
	err = PetPlayAction(t, tf, petName, playToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's energy, hygiene, and wellness are updated correctly.
	petEnergy, err = cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	assert.Less(t, petEnergy.E, 100)

	petHygiene, err = cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.Less(t, petHygiene.Hy, 100)

	petWellness, err = cardinal.GetComponent[component.Wellness](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, petWellness.Wn, game.MaxWellness)
}

// Test insufficient pet energy
func TestSystem_PetPlayAction_InsufficientPetEnergy(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	world := cardinal.NewWorldContext(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// - The pet's energy is set to a low value.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	assert.NotEqual(t, petId, 0)

	petEnergy, err := cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	petEnergy.E = 0
	err = cardinal.SetComponent(world, petId, petEnergy)
	assert.NoError(t, err)

	// - The player buys a toy.
	err = buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	// When:
	// - The PetPlayAction function is called with the pet.
	err = PetPlayAction(t, tf, petName, playToyName)
	assert.Error(t, err)

	// Then:
	// - The function returns an error message indicating that the pet's energy is insufficient.
	assert.Contains(t, err.Error(), "pet energy is insufficient")
}

// Test exceeding maximum hygiene or wellness
func TestSystem_PetPlayAction_ExceedingMaximumHygieneOrWellness(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	// Given:
	// - A test fixture is initialized.
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

	// - The pet's hygiene and wellness are set to values that exceed the maximum.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	assert.NotEqual(t, petId, 0)

	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = 100 + 1
	err = cardinal.SetComponent(world, petId, petHygiene)

	assert.NoError(t, err)

	petWellness, err := cardinal.GetComponent[component.Wellness](wCtx, petId)
	assert.NoError(t, err)
	petWellness.Wn = 100 + 1
	err = cardinal.SetComponent(world, petId, petWellness)
	assert.NoError(t, err)

	// - The player buys a toy.
	err = buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	// When:
	// - The PetPlayAction function is called with the pet.
	err = PetPlayAction(t, tf, petName, playToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's hygiene and wellness values are capped at the maximum.
	petHygiene, err = cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.LessOrEqual(t, petHygiene.Hy, 100)

	petWellness, err = cardinal.GetComponent[component.Wellness](wCtx, petId)
	assert.NoError(t, err)
	assert.LessOrEqual(t, petWellness.Wn, 100)
}

// Test non-existent item
func TestSystem_PetPlayAction_NonExistentItem(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	// - A pet is created that belongs to the player.
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	createPet(t, tf, petName, personaTag)

	// When:
	// - The PetPlayAction function is called with a non-existent item.
	err := PetPlayAction(t, tf, petName, "non-existent-item")
	assert.Error(t, err)

	// Then:
	// - The function returns an error message indicating that the item doesn't exist.
	assert.Contains(t, err.Error(), "item not found")
}

// Test multiple updates
func TestSystem_PetPlayAction_MultipleUpdates(t *testing.T) {
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
	err := buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	// - The pet's components are set to initial values.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	assert.NotEqual(t, petId, 0)

	petEnergy, err := cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	petEnergy.E = game.MaxEnergy
	err = cardinal.SetComponent(world, petId, petEnergy)
	assert.NoError(t, err)

	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = game.MaxHygiene
	err = cardinal.SetComponent(world, petId, petHygiene)
	assert.NoError(t, err)

	petWellness, err := cardinal.GetComponent[component.Wellness](wCtx, petId)
	assert.NoError(t, err)
	petWellness.Wn = game.MaxWellness
	err = cardinal.SetComponent(world, petId, petWellness)
	assert.NoError(t, err)

	petActivity, err := cardinal.GetComponent[component.Activity](wCtx, petId)
	assert.NoError(t, err)
	petActivity.Activity = game.InitialActivity
	err = cardinal.SetComponent(world, petId, petActivity)
	assert.NoError(t, err)

	petThink, err := cardinal.GetComponent[component.Think](wCtx, petId)
	assert.NoError(t, err)
	petThink.Think = game.InitialThink
	err = cardinal.SetComponent(world, petId, petThink)
	assert.NoError(t, err)

	// When:
	// - The PetPlayAction function is called with multiple updates.
	err = PetPlayAction(t, tf, petName, playToyName)
	assert.NoError(t, err)

	// Then:
	// - The pet's components are updated correctly.
	petEnergy, err = cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	assert.Less(t, petEnergy.E, game.MaxEnergy)

	petHygiene, err = cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.Less(t, petHygiene.Hy, game.MaxHygiene)

	petWellness, err = cardinal.GetComponent[component.Wellness](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, petWellness.Wn, game.MaxWellness)

	petActivity, err = cardinal.GetComponent[component.Activity](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, "Playing", petActivity.Activity)
}

// Test concurrent updates
// TODO: sometimes this test fail, concurrency problem?
func TestSystem_PetPlayAction_ConcurrentUpdates(t *testing.T) {
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
	err := buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	// - The pet's components are set to initial values.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	assert.NotEqual(t, petId, 0)

	petEnergy, err := cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	petEnergy.E = game.MaxEnergy
	err = cardinal.SetComponent(world, petId, petEnergy)
	assert.NoError(t, err)

	petHygiene, err := cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	petHygiene.Hy = game.MaxHygiene
	err = cardinal.SetComponent(world, petId, petHygiene)
	assert.NoError(t, err)

	petWellness, err := cardinal.GetComponent[component.Wellness](wCtx, petId)
	assert.NoError(t, err)
	petWellness.Wn = game.MaxWellness
	err = cardinal.SetComponent(world, petId, petWellness)
	assert.NoError(t, err)

	petActivity, err := cardinal.GetComponent[component.Activity](wCtx, petId)
	assert.NoError(t, err)
	petActivity.Activity = game.InitialActivity
	err = cardinal.SetComponent(world, petId, petActivity)
	assert.NoError(t, err)

	petThink, err := cardinal.GetComponent[component.Think](wCtx, petId)
	assert.NoError(t, err)
	petThink.Think = game.InitialThink
	err = cardinal.SetComponent(world, petId, petThink)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	simul_txs := 1 // TODO: increase this to reproduce the error
	count_err := 0
	count_ok := 0
	for i := 0; i < simul_txs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := PetPlayAction(t, tf, petName, playToyName)
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
	// - The pet's components are updated correctly.
	petEnergy, err = cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	assert.Less(t, petEnergy.E, game.MaxEnergy)

	petHygiene, err = cardinal.GetComponent[component.Hygiene](wCtx, petId)
	assert.NoError(t, err)
	assert.Less(t, petHygiene.Hy, game.MaxHygiene)

	petWellness, err = cardinal.GetComponent[component.Wellness](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, petWellness.Wn, game.MaxWellness)

	petActivity, err = cardinal.GetComponent[component.Activity](wCtx, petId)
	assert.NoError(t, err)
	assert.Equal(t, "Playing", petActivity.Activity)
}

// Test pet's energy, hygiene, or wellness reaches 0
func TestSystem_PetPlayAction_ZeroEnergy(t *testing.T) {
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
	err := buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	// - The pet's components are set to initial values.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	assert.NotEqual(t, petId, 0)

	petEnergy, err := cardinal.GetComponent[component.Energy](wCtx, petId)
	assert.NoError(t, err)
	petEnergy.E = 0
	err = cardinal.SetComponent(world, petId, petEnergy)
	assert.NoError(t, err)

	// When:
	// - The PetPlayAction function is called.
	err = PetPlayAction(t, tf, petName, playToyName)
	// Then:
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pet energy is insufficient")

}

// Test pet's level exceeds the maximum value
func TestSystem_PetPlayAction_MaxLevel(t *testing.T) {
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
	err := buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	// - The pet's components are set to initial values.
	petId, pet, err := component.GetPetByNickname(wCtx, petName)
	assert.NoError(t, err)
	assert.NotEqual(t, petId, 0)

	pet.Level = game.MaxLevel
	err = cardinal.SetComponent(world, petId, pet)
	assert.NoError(t, err)

	// When:
	// - The PetPlayAction function is called.
	err = PetPlayAction(t, tf, petName, playToyName)
	// Then:
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pet Max lvl, cant grow more")

}

// Test pet's data is empty or null
// Interested border case that a `Pet` is removed by the `system`, but the reference in `Player` array is not updated
// As result since `Pet` object is not accessed, but other components like Hygiene, health, etc are changed on the function call
// The function execute successful, despite of being an inconsistency in data
func TestSystem_PetPlayAction_NullPetData(t *testing.T) {
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.
	tf := cardinal.NewTestFixture(t, nil)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	world := cardinal.NewWorldContext(tf.World)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// - A pet is created that belongs to the player.
	err := createPet(t, tf, petName, personaTag)
	assert.NoError(t, err)

	// - The player buys a toy.
	err = buyToy(t, tf, playToyName)
	assert.NoError(t, err)

	// - The pet's is removed.
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)

	err = cardinal.Remove(world, petId)
	assert.NoError(t, err)

	// When:
	// - The PetPlayAction function is called.
	err = PetPlayAction(t, tf, petName, playToyName)
	assert.NoError(t, err)

}
