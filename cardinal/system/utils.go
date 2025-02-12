package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
)

// queryTargetHealthPet queries for the target pet's entity ID and health component.
func queryTargetHealthPet(world cardinal.WorldContext, targetNickname string) (types.EntityID, *comp.Health, error) {
	var petID types.EntityID
	var petHealth *comp.Health
	var err error
	searchErr := cardinal.NewSearch().Entity(
		filter.Exact(filter.Component[comp.Pet](), filter.Component[comp.Health]())).Each(world,
		func(id types.EntityID) bool {
			var pet *comp.Pet
			pet, err = cardinal.GetComponent[comp.Pet](world, id)
			if err != nil {
				return false
			}

			// Terminates the search if the pet is found
			if pet.Nickname == targetNickname {
				petID = id
				petHealth, err = cardinal.GetComponent[comp.Health](world, id)
				if err != nil {
					return false
				}
				return false
			}

			// Continue searching if the pet is not the target pet
			return true
		})
	if searchErr != nil {
		return 0, nil, err
	}
	if err != nil {
		return 0, nil, err
	}
	if petHealth == nil {
		return 0, nil, fmt.Errorf("pet %q does not exist", targetNickname)
	}

	return petID, petHealth, err
}

// queryTargetEnergyPet queries for the target pet's entity ID and energy component.
func queryTargetEnergyPet(world cardinal.WorldContext, targetNickname string) (types.EntityID, *comp.Energy, error) {
	var petID types.EntityID
	var petEnergy *comp.Energy
	var err error
	searchErr := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.Pet](), filter.Component[comp.Energy]())).Each(world,
		func(id types.EntityID) bool {
			var pet *comp.Pet
			pet, err = cardinal.GetComponent[comp.Pet](world, id)
			if err != nil {
				return false
			}

			// Terminates the search if the pet is found
			if pet.Nickname == targetNickname {
				petID = id
				petEnergy, err = cardinal.GetComponent[comp.Energy](world, id)
				if err != nil {
					return false
				}
				return false
			}

			// Continue searching if the pet is not the target pet
			return true
		})
	if searchErr != nil {
		return 0, nil, err
	}
	if err != nil {
		return 0, nil, err
	}
	if petEnergy == nil {
		return 0, nil, fmt.Errorf("pet %q does not exist", targetNickname)
	}

	return petID, petEnergy, err
}

func queryPetIdByName(world cardinal.WorldContext, name string) (types.EntityID, error) {
	var petID types.EntityID
	var err error
	searchErr := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.Pet]())).Each(world,
		func(id types.EntityID) bool {
			var pet *comp.Pet
			pet, err = cardinal.GetComponent[comp.Pet](world, id)
			if err != nil {
				return false
			}

			// Terminates the search if the pet is found
			if pet.Nickname == name {
				petID = id
				return false
			}

			// Continue searching if the pet is not the target pet
			return true
		})
	if searchErr != nil {
		return 0, err
	}
	return petID, err
}
