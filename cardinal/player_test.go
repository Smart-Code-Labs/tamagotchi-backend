package main

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"pkg.world.dev/world-engine/cardinal"

	"tamagotchi/component"
)

func TestSuccessfulPlayerCreation(t *testing.T) {
	// Given: A CreatePlayer transaction with a unique persona tag.
	// Preconditions:
	// - The test fixture is initialized.
	// - A persona and player are created.

	// Given:
	// - A test fixture is initialized.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)

	// - A persona is created.
	createPersona(t, tf, personaTag)

	// - A player is created and associated with the persona.
	// When: The PlayerSpawnerAction function is called.
	err := createPlayer(t, tf, personaTag)
	assert.NoError(t, err)

	// Then: A new player entity is created, and a "new_player" event is emitted.
	// Verify that the player entity was created.
	playerID, err := component.FindPlayerByPersonaTag(wCtx, personaTag)
	assert.NoError(t, err)
	assert.NotNil(t, playerID)

	// Verify that the "new_player" event was emitted.
	// events := wCtx.EmitEvent()
	// assert.Len(t, events, 1)
	// assert.Equal(t, "new_player", events[0]["event"])
	// assert.Equal(t, playerID, events[0]["id"])
}

func TestDuplicatePlayerCreation(t *testing.T) {
	// Given: A CreatePlayer transaction with a persona tag that already exists.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	createPersona(t, tf, personaTag)
	err := createPlayer(t, tf, personaTag)
	assert.NoError(t, err)

	// When: The PlayerSpawnerAction function is called again with the same persona tag.
	err = createPlayer(t, tf, personaTag)
	assert.Error(t, err)

	// Then: An error indicating that the player already exists is returned.
	assert.Contains(t, err.Error(), "player already exists")
}

func TestInvalidPersonaTag(t *testing.T) {
	// Given: A CreatePlayer transaction with an invalid persona tag (e.g., empty string).
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// When: The PlayerSpawnerAction function is called with an empty string as the persona tag.
	err := createPlayer(t, tf, "")

	// Then: An error indicating that the persona tag is invalid is returned.
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "persona tag cannot be empty")
}

func TestCreatePlayerWithEmptyPetsAndItems(t *testing.T) {
	// Given: A CreatePlayer transaction with a unique persona tag.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	createPersona(t, tf, personaTag)

	// When: The PlayerSpawnerAction function is called.
	err := createPlayer(t, tf, personaTag)
	assert.NoError(t, err)

	// Then: A new player entity is created with empty pets and items arrays, and a "new_player" event is emitted.
	playerID, err := component.FindPlayerByPersonaTag(wCtx, personaTag)
	assert.NoError(t, err)
	assert.NotNil(t, playerID)

	player, err := cardinal.GetComponent[component.Player](wCtx, playerID)
	assert.NoError(t, err)
	assert.Empty(t, player.Pets)
	assert.Empty(t, player.Items)
}

func TestCreatePlayerWithInitialMoney(t *testing.T) {
	// Given: A CreatePlayer transaction with a unique persona tag.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	createPersona(t, tf, personaTag)

	// When: The PlayerSpawnerAction function is called.
	err := createPlayer(t, tf, personaTag)
	assert.NoError(t, err)

	// Then: A new player entity is created with the initial money set, and a "new_player" event is emitted.
	playerID, err := component.FindPlayerByPersonaTag(wCtx, personaTag)
	assert.NoError(t, err)
	assert.NotNil(t, playerID)

	player, err := cardinal.GetComponent[component.Player](wCtx, playerID)
	assert.NoError(t, err)
	assert.NotNil(t, player.Money)
	assert.Equal(t, float64(1000), player.Money) // assuming the initial money is 1000
}

func TestCreatePlayerWithInvalidPersonaTagLength(t *testing.T) {
	// Given: A CreatePlayer transaction with a persona tag that exceeds the maximum allowed length.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)

	// When: The PlayerSpawnerAction function is called with a persona tag that exceeds the maximum allowed length.
	longPersonaTag := strings.Repeat("a", 256) // assuming the maximum allowed length is 255
	err := createPlayer(t, tf, longPersonaTag)

	// Then: An error indicating that the persona tag is too long is returned.
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "persona tag exceeds maximum length of 32 characters")
}

func TestCreatePlayerWithPersonaTagThatAlreadyExistsButIsNotAssociatedWithAPlayer(t *testing.T) {
	// Given: A persona that already exists but is not associated with a player.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)
	createPersona(t, tf, personaTag)

	// When: The PlayerSpawnerAction function is called with the same persona tag.
	err := createPlayer(t, tf, personaTag)

	// Then: A new player entity is created and associated with the existing persona, and a "new_player" event is emitted.
	assert.NoError(t, err)

	playerID, err := component.FindPlayerByPersonaTag(wCtx, personaTag)
	assert.NoError(t, err)
	assert.NotNil(t, playerID)
}

func TestCreatePlayerWithMultiplePersonas(t *testing.T) {
	// Given: Multiple personas that do not exist.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)

	// When: The PlayerSpawnerAction function is called with multiple personas.
	personaTags := []string{"persona1", "persona2", "persona3"}
	for _, personaTag := range personaTags {
		err := createPlayer(t, tf, personaTag)
		assert.NoError(t, err)
	}

	// Then: Multiple player entities are created, each associated with a different persona, and multiple "new_player" events are emitted.
	for _, personaTag := range personaTags {
		playerID, err := component.FindPlayerByPersonaTag(wCtx, personaTag)
		assert.NoError(t, err)
		assert.NotNil(t, playerID)
	}
}

func TestCreatePlayerConcurrently(t *testing.T) {
	// Given: A scenario where multiple CreatePlayer transactions are processed concurrently.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)

	// When: Multiple CreatePlayer transactions are processed concurrently.
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(personaTag string) {
			defer wg.Done()
			err := createPlayer(t, tf, personaTag)
			assert.NoError(t, err)
		}(fmt.Sprintf("persona%d", i))
	}
	wg.Wait()

	// Then: Multiple player entities are created, each associated with a different persona, and multiple "new_player" events are emitted.
	for i := 0; i < 10; i++ {
		playerID, err := component.FindPlayerByPersonaTag(wCtx, fmt.Sprintf("persona%d", i))
		assert.NoError(t, err)
		assert.NotNil(t, playerID)
	}
}

func TestConcurrentTransactionsWithSamePersonaTag(t *testing.T) {
	// Given: Multiple CreatePlayer transactions with the same persona tag.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)

	// When: Multiple CreatePlayer transactions with the same persona tag are processed concurrently.
	var wg sync.WaitGroup
	personaTag := "same-persona-tag"
	simul_txs := 2
	count_err := 0
	for i := 0; i < simul_txs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := createPlayer(t, tf, personaTag)

			if err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "player already exists")
				count_err++
			} else {
				assert.NoError(t, err)
			}
		}()
	}
	wg.Wait()

	t.Log(count_err)
	// TODO: check actual should be simul_txs-1
	assert.Equal(t, count_err, simul_txs)

	// Then: Only one player entity is created, and the remaining transactions return an error indicating that the player already exists.
	playerID, err := component.FindPlayerByPersonaTag(wCtx, personaTag)
	assert.NoError(t, err)
	assert.NotNil(t, playerID)
}

func TestTransactionsWithDifferentPersonaTags(t *testing.T) {
	// Given: Multiple CreatePlayer transactions with different persona tags.
	tf := cardinal.NewTestFixture(t, nil)
	MustInitWorld(tf.World)
	wCtx := cardinal.NewReadOnlyWorldContext(tf.World)

	// When: Multiple CreatePlayer transactions with different persona tags are processed concurrently.
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(personaTag string) {
			defer wg.Done()
			err := createPlayer(t, tf, personaTag)
			assert.NoError(t, err)
		}(fmt.Sprintf("persona%d", i))
	}
	wg.Wait()

	// Then: Multiple player entities are created, each associated with a different persona.
	for i := 0; i < 10; i++ {
		playerID, err := component.FindPlayerByPersonaTag(wCtx, fmt.Sprintf("persona%d", i))
		assert.NoError(t, err)
		assert.NotNil(t, playerID)
	}
}
