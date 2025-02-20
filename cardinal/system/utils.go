package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "tamagotchi/component"
	constants "tamagotchi/game"
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

func QueryPetIdByName(world cardinal.WorldContext, name string) (types.EntityID, error) {
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

// func QueryPersonaItemIdList(world cardinal.WorldContext, personaTag string) ([]types.EntityID, error) {
// 	var err error
// 	var item *comp.Item

// 	list := make([]types.EntityID, 0)

// 	q := cardinal.NewSearch().Entity(
// 		filter.Contains(filter.Component[comp.Item]()))
// 	searchErr := q.Each(world,
// 		func(id types.EntityID) bool {
// 			item, err = cardinal.GetComponent[comp.Item](world, id)
// 			if item.PersonaTag == personaTag {
// 				list = append(list, id)
// 			}

// 			// Continue searching
// 			return true
// 		})
// 	if searchErr != nil {
// 		return list, err
// 	}
// 	return list, err
// }

func GetMessageForRange(value int, rangeMessages map[constants.Range]constants.Message) (bool, string) {
	for r, m := range rangeMessages {
		if value >= r.Min && value <= r.Max {
			return true, m.Text // Return the message and true (success)
		}
	}
	return false, "" // Return an empty string and false (no match)
}
