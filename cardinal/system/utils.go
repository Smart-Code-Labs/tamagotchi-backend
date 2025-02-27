package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	"tamagotchi/component"
	"tamagotchi/game"
)

// queryTargetHealthPet queries for the target pet's entity ID and health component.
func queryTargetHealthPet(world cardinal.WorldContext, targetNickname string) (types.EntityID, *component.Health, error) {
	var petID types.EntityID
	var petHealth *component.Health
	var err error
	searchErr := cardinal.NewSearch().Entity(
		filter.Exact(filter.Component[component.Pet](), filter.Component[component.Health]())).Each(world,
		func(id types.EntityID) bool {
			var pet *component.Pet
			pet, err = cardinal.GetComponent[component.Pet](world, id)
			if err != nil {
				return false
			}

			// Terminates the search if the pet is found
			if pet.Nickname == targetNickname {
				petID = id
				petHealth, err = cardinal.GetComponent[component.Health](world, id)
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
func queryTargetEnergyPet(world cardinal.WorldContext, targetNickname string) (types.EntityID, *component.Energy, error) {
	var petID types.EntityID
	var petEnergy *component.Energy
	var err error
	searchErr := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[component.Pet](), filter.Component[component.Energy]())).Each(world,
		func(id types.EntityID) bool {
			var pet *component.Pet
			pet, err = cardinal.GetComponent[component.Pet](world, id)
			if err != nil {
				return false
			}

			// Terminates the search if the pet is found
			if pet.Nickname == targetNickname {
				petID = id
				petEnergy, err = cardinal.GetComponent[component.Energy](world, id)
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

// func QueryPersonaItemIdList(world cardinal.WorldContext, personaTag string) ([]types.EntityID, error) {
// 	var err error
// 	var item *component.Item

// 	list := make([]types.EntityID, 0)

// 	q := cardinal.NewSearch().Entity(
// 		filter.Contains(filter.Component[component.Item]()))
// 	searchErr := q.Each(world,
// 		func(id types.EntityID) bool {
// 			item, err = cardinal.GetComponent[component.Item](world, id)
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

func GetMessageForRange(value int, rangeMessages map[game.Range]game.Message) (bool, string) {
	for r, m := range rangeMessages {
		if value >= r.Min && value <= r.Max {
			return true, m.Text // Return the message and true (success)
		}
	}
	return false, "" // Return an empty string and false (no match)
}
