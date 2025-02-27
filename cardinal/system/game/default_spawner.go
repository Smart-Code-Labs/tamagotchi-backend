// Package system contains the logic for spawning default game systems.
package system

import (
	"tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
)

/**
 * Function Flow:
 * 1. The `SpawnDefaultSystem` function is called, which creates a Leaderboard entity.
 * 2. The function then attempts to create a Drug Store, Food Store, and Toy Store.
 * 3. For each store, the function calls a corresponding creation function (`createDrugStore`, `createFoodStore`, `createToyStore`).
 * 4. If any of the store creation functions return an error, the `SpawnDefaultSystem` function will return that error.
 *
 * SpawnDefaultSystem creates a Leaderboard and default stores. This System is registered as an
 * Init system, meaning it will be executed exactly one time on tick 0.
 *
 * @param world The WorldContext for the game.
 * @return error if any error occurs during the spawning of default systems.
 */
func SpawnDefaultSystem(world cardinal.WorldContext) error {

	// Step 1: Create LeaderBoard
	//   - Create a new Leaderboard entity using the `Create` function
	//   - If the creation fails, return the error
	_, err := cardinal.Create(world,
		component.Leaderboard{},
	)
	if err != nil {
		return err
	}

	// Step 2: Create Drug Store
	//   - Call the `createDrugStore` function to create a new Drug Store entity
	//   - If the function returns an error, return that error
	shouldReturn, err := createDrugStore(world)
	if shouldReturn {
		return err
	}

	// Step 3: Create Food Store
	//   - Call the `createFoodStore` function to create a new Food Store entity
	//   - If the function returns an error, return that error
	shouldReturn, err = createFoodStore(world)
	if shouldReturn {
		return err
	}

	// Step 4: Create Toy Store
	//   - Call the `createToyStore` function to create a new Toy Store entity
	//   - If the function returns an error, return that error
	shouldReturn, err = createToyStore(world)
	if shouldReturn {
		return err
	}

	return nil
}

/**
 * Function Flow:
 * 1. The `createDrugStore` function is called, which creates a new Drug Store entity.
 * 2. The function retrieves the created Drug Store component and initializes it.
 * 3. The function updates the Drug Store component with the initialized values.
 * 4. If any of the steps fail, the function returns an error.
 *
 * createDrugStore creates a new Drug Store entity and initializes it.
 *
 * @param world The WorldContext for the game.
 * @return bool indicating whether to return immediately, and error if any error occurs during the creation of the Drug Store.
 */
func createDrugStore(world cardinal.WorldContext) (bool, error) {

	// Step 1: Create Drug Store entity
	//   - Create a new Drug Store entity using the `Create` function
	//   - If the creation fails, return an error
	entityId, err := cardinal.Create(world,
		component.DrugStore{},
	)
	if err != nil {
		return true, err
	}

	// Step 2: Init Drug Store
	//   - Retrieve the created Drug Store component using the `GetComponent` function
	//   - Initialize the Drug Store component using the `InitDrugStore` method
	drugStore, err := cardinal.GetComponent[component.DrugStore](world, entityId)
	if err != nil {
		return true, err
	}
	drugStore.InitDrugStore(world)

	// Step 3: Update Drug Store
	//   - Update the Drug Store component with the initialized values using the `SetComponent` function
	//   - If the update fails, return an error
	err = cardinal.SetComponent(world, entityId, drugStore)
	if err != nil {
		return true, err
	}
	return false, nil
}

/**
 * Function Flow:
 * 1. The `createFoodStore` function is called, which creates a new Food Store entity.
 * 2. The function retrieves the created Food Store component and initializes it.
 * 3. The function updates the Food Store component with the initialized values.
 * 4. If any of the steps fail, the function returns an error.
 *
 * createFoodStore creates a new Food Store entity and initializes it.
 *
 * @param world The WorldContext for the game.
 * @return bool indicating whether to return immediately, and error if any error occurs during the creation of the Food Store.
 */
func createFoodStore(world cardinal.WorldContext) (bool, error) {

	// Step 1: Create Food Store entity
	//   - Create a new Food Store entity using the `Create` function
	//   - If the creation fails, return an error
	entityId, err := cardinal.Create(world,
		component.FoodStore{},
	)
	if err != nil {
		return true, err
	}

	// Step 2: Init Food Store
	//   - Retrieve the created Food Store component using the `GetComponent` function
	//   - Initialize the Food Store component using the `InitFoodStore` method
	foodStore, err := cardinal.GetComponent[component.FoodStore](world, entityId)
	if err != nil {
		return true, err
	}
	foodStore.InitFoodStore(world)

	// Step 3: Update Food Store
	//   - Update the Food Store component with the initialized values using the `SetComponent` function
	//   - If the update fails, return an error
	err = cardinal.SetComponent(world, entityId, foodStore)
	if err != nil {
		return true, err
	}
	return false, nil
}

/**
 * Function Flow:
 * 1. The `createToyStore` function is called, which creates a new Toy Store entity.
 * 2. The function retrieves the created Toy Store component and initializes it.
 * 3. The function updates the Toy Store component with the initialized values.
 * 4. If any of the steps fail, the function returns an error.
 *
 * createToyStore creates a new Toy Store entity and initializes it.
 *
 * @param world The WorldContext for the game.
 * @return bool indicating whether to return immediately, and error if any error occurs during the creation of the Toy Store.
 */
func createToyStore(world cardinal.WorldContext) (bool, error) {
	// Step 1: Create Toy Store entity
	//   - Create a new Toy Store entity using the `Create` function
	//   - If the creation fails, return an error
	entityId, err := cardinal.Create(world,
		component.ToyStore{},
	)
	if err != nil {
		return true, err
	}

	// Step 2: Init Toy Store
	//   - Retrieve the created Toy Store component using the `GetComponent` function
	//   - Initialize the Toy Store component using the `InitToyStore` method
	// Note: The variable name 'foodStore' should be 'toyStore' for clarity.
	foodStore, err := cardinal.GetComponent[component.ToyStore](world, entityId)
	if err != nil {
		return true, err
	}
	foodStore.InitToyStore(world)

	// Step 3: Update Toy Store
	//   - Update the Toy Store component with the initialized values using the `SetComponent` function
	//   - If the update fails, return an error
	// Note: The variable name 'foodStore' should be 'toyStore' for clarity.
	err = cardinal.SetComponent(world, entityId, foodStore)
	if err != nil {
		return true, err
	}
	return false, nil
}
