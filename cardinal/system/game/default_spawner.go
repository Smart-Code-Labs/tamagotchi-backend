package system

import (
	comp "tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
)

// SpawnDefaultSystem creates a Leaderboard. This System is registered as an
// Init system, meaning it will be executed exactly one time on tick 0.
func SpawnDefaultSystem(world cardinal.WorldContext) error {
	// Create LeaderBoard
	_, err := cardinal.Create(world,
		comp.Leaderboard{},
	)
	if err != nil {
		return err
	}

	// Create Drug Store
	shouldReturn, err := createDrugStore(world)
	if shouldReturn {
		return err
	}

	// Create Drug Store
	shouldReturn, err = createFoodStore(world)
	if shouldReturn {
		return err
	}

	// Create Drug Store
	shouldReturn, err = createToyStore(world)
	if shouldReturn {
		return err
	}

	return nil
}

func createDrugStore(world cardinal.WorldContext) (bool, error) {
	entityId, err := cardinal.Create(world,
		comp.DrugStore{},
	)
	if err != nil {
		return true, err
	}

	// Init Drug Store
	drugStore, err := cardinal.GetComponent[comp.DrugStore](world, entityId)
	if err != nil {
		return true, err
	}
	drugStore.InitDrugStore(world)

	// Update
	err = cardinal.SetComponent(world, entityId, drugStore)
	if err != nil {
		return true, err
	}
	return false, nil
}

func createFoodStore(world cardinal.WorldContext) (bool, error) {
	entityId, err := cardinal.Create(world,
		comp.FoodStore{},
	)
	if err != nil {
		return true, err
	}

	// Init Drug Store
	foodStore, err := cardinal.GetComponent[comp.FoodStore](world, entityId)
	if err != nil {
		return true, err
	}
	foodStore.InitFoodStore(world)

	// Update
	err = cardinal.SetComponent(world, entityId, foodStore)
	if err != nil {
		return true, err
	}
	return false, nil
}

func createToyStore(world cardinal.WorldContext) (bool, error) {
	entityId, err := cardinal.Create(world,
		comp.ToyStore{},
	)
	if err != nil {
		return true, err
	}

	// Init Drug Store
	foodStore, err := cardinal.GetComponent[comp.ToyStore](world, entityId)
	if err != nil {
		return true, err
	}
	foodStore.InitToyStore(world)

	// Update
	err = cardinal.SetComponent(world, entityId, foodStore)
	if err != nil {
		return true, err
	}
	return false, nil
}
