package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"pkg.world.dev/world-engine/cardinal"

	"tamagotchi/component"
	"tamagotchi/query"
)

// This function tests the creation of a pet.
// Flow:
// 1. Initialize a new test fixture.
// 2. Create a persona.
// 3. Create a player associated with the persona.
// 4. Create a pet that belongs to the player.
// 5. Verify that the pet was created successfully.
func TestSystem_PetSpawnerSystem_CanCreatePet2(t *testing.T) {
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

// Test 2: Duplicate Pet Nickname
func TestPet_PetSpawnerSystem_DuplicatePetNickname(t *testing.T) {
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
	// - Another pet is created with the same nickname.
	err := createPet(t, tf, petName, personaTag)

	// Then:
	// - The error is verified to be "Pet nickname already exists".
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating pet: Name already exist")
}

// Test 5: Pet Creation with Empty Nickname
func TestSystem_PetSpawnerSystem_PetCreationWithEmptyNickname(t *testing.T) {
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// When:
	// - A pet is created with an empty nickname.
	err := createPet(t, tf, "", personaTag)

	// Then:
	// - The error is verified to be "pet nickname cannot be empty".
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pet nickname cannot be empty")
}

// Test 6: Pet Creation with Null Persona Tag
func TestSystem_PetSpawnerSystem_PetCreationWithNullPersonaTag(t *testing.T) {
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is not created.

	// When:
	// - A pet is created with a null persona tag.
	err := createPet(t, tf, petName, "")

	// Then:
	// - The error is verified to be "persona tag cannot be null".
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "persona tag cannot be null")
}

// Test 7: Multiple Pet Creations
func TestSystem_PetSpawnerSystem_MultiplePetCreations(t *testing.T) {
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// When:
	// - Multiple pets are created with unique nicknames.
	petNames := []string{"pet1", "pet2", "pet3"}
	for _, petName := range petNames {
		err := createPet(t, tf, petName, personaTag)
		assert.NoError(t, err)
	}

	// Then:
	// - All pets are verified to exist and be owned by the player.
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	player, err := component.GetPlayerByPersonaTag(wCtx, personaTag)
	assert.NoError(t, err)
	for _, petName := range petNames {
		_, petId, err := component.QueryPetIdByName(wCtx, petName)
		assert.NoError(t, err)
		assert.NotZero(t, petId)
		playerPetId, err := player.GetPetNickname(wCtx, petName)
		assert.NoError(t, err)
		assert.Equal(t, playerPetId, petId, "Expected pet ID %v to be in the player's Pets array, but got %v", petId, playerPetId)
	}
}

// Test 10: Edge Case - Nickname with Special Characters
func TestSystem_PetSpawnerSystem_PetCreationWithSpecialCharacters(t *testing.T) {
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// When:
	// - A pet is created with a nickname containing special characters.
	petNameWithSpecialChars := "pet!@#$"
	err := createPet(t, tf, petNameWithSpecialChars, personaTag)

	// Then:
	// - The pet is verified to exist and be owned by the player.
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nickname can only contain letters and numbers")
}

// Test 11: Edge Case - Nickname with Numbers
func TestSystem_PetSpawnerSystem_PetCreationWithNumbers(t *testing.T) {
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	err := createPersona(t, tf, personaTag)
	assert.NoError(t, err)
	// - A player is created and associated with the persona.
	err = createPlayer(t, tf, personaTag)
	assert.NoError(t, err)
	// When:
	// - A pet is created with a nickname containing numbers.
	petNameWithNumbers := "pet123"
	err = createPet(t, tf, petNameWithNumbers, personaTag)

	// Then:
	// - The pet is verified to exist and be owned by the player.
	assert.NoError(t, err)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	_, petId, err := component.QueryPetIdByName(wCtx, petNameWithNumbers)
	assert.NoError(t, err)
	assert.NotZero(t, petId)
}

// Test 12: Edge Case - Nickname with Uppercase Letters
func TestSystem_PetSpawnerSystem_PetCreationWithUppercaseLetters(t *testing.T) {
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	err := createPlayer(t, tf, personaTag)
	assert.NoError(t, err)
	// When:
	// - A pet is created with a nickname containing uppercase letters.
	petNameWithUppercaseLetters := "PetName"
	err = createPet(t, tf, petNameWithUppercaseLetters, personaTag)

	// Then:
	// - The pet is verified to exist and be owned by the player.
	assert.NoError(t, err)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	_, petId, err := component.QueryPetIdByName(wCtx, petNameWithUppercaseLetters)
	assert.NoError(t, err)
	assert.NotZero(t, petId)
}

// Test 18: Concurrency - Create multiple pets concurrently
func TestSystem_PetSpawnerSystem_Concurrency_CreateMultiplePetsConcurrently(t *testing.T) {
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	err := createPersona(t, tf, personaTag)
	assert.NoError(t, err)
	// - A player is created and associated with the persona.
	err = createPlayer(t, tf, personaTag)
	assert.NoError(t, err)

	// When:
	// - Multiple pets are created concurrently.
	var wg sync.WaitGroup
	petNames := []string{"pet1", "pet2", "pet3"}
	for _, petName := range petNames {
		wg.Add(1)
		go func(petName string) {
			defer wg.Done()
			err := createPet(t, tf, petName, personaTag)
			assert.NoError(t, err)
		}(petName)
	}
	wg.Wait()

	// Then:
	// - All pets are verified to exist and be owned by the player.
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	player, err := component.GetPlayerByPersonaTag(wCtx, personaTag)
	assert.NoError(t, err)
	for _, petName := range petNames {
		_, petId, err := component.QueryPetIdByName(wCtx, petName)
		assert.NoError(t, err)
		assert.NotZero(t, petId)
		playerPetId, err := player.GetPetNickname(wCtx, petName)
		assert.NoError(t, err)
		assert.Equal(t, playerPetId, petId, "Expected pet ID %v to be in the player's Pets array, but got %v", petId, playerPetId)
	}
}

// Test 19: Concurrency - Create multiple pets concurrently with same nickname
func TestSystem_PetSpawnerSystem_Concurrency_CreateMultiplePetsConcurrentlyWithSameNickname(t *testing.T) {
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	createPlayer(t, tf, personaTag)

	// When:
	// - Multiple pets are created concurrently with the same nickname.
	var wg sync.WaitGroup
	petName := "pet"
	numPets := 3
	for i := 0; i < numPets; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := createPet(t, tf, petName, personaTag)
			if err != nil {
				t.Logf("Error creating pet: %v", err)
			}
		}()
	}
	wg.Wait()

	// Then:
	// - Only one pet is verified to exist and be owned by the player.
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	player, err := component.GetPlayerByPersonaTag(wCtx, personaTag)
	assert.NoError(t, err)
	_, petId, err := component.QueryPetIdByName(wCtx, petName)
	assert.NoError(t, err)
	assert.NotZero(t, petId)
	playerPetId, err := player.GetPetNickname(wCtx, petName)
	assert.NoError(t, err)
	assert.Equal(t, playerPetId, petId, "Expected pet ID %v to be in the player's Pets array, but got %v", petId, playerPetId)

	// Verify that only one pet is created.
	assert.Len(t, player.Pets, 1, "Expected only one pet to be created, but got %v", len(player.Pets))
}

// Test 20: Concurrency - Create multiple pets concurrently with non-existent player
func TestSystem_PetSpawnerSystem_Concurrency_CreateMultiplePetsConcurrentlyWithNonExistentPlayer(t *testing.T) {
	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// - A persona is not created.

	// When:
	// - Multiple pets are created concurrently with a non-existent player persona tag.
	var wg sync.WaitGroup
	petName := "pet"
	numPets := 3
	for i := 0; i < numPets; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := createPet(t, tf, petName, personaTag)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "player not found")
		}()
	}
	wg.Wait()

	// Then:
	// - No pets are verified to exist.
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	msg := query.PetsMsg{}
	response, err := query.GamePets(wCtx, &msg)
	assert.NoError(t, err)
	assert.Empty(t, response.Pets, "Expected no pets to be created, but got %v", len(response.Pets))
}
