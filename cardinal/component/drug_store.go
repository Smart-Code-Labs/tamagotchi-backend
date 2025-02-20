package component

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// Define a struct to hold drug properties including description
type DrugProperties struct {
	Price       float64
	Value       int
	Description string
}

// Initialize the DrugKinds map with DrugProperties including descriptions
var DrugKinds = map[string]DrugProperties{
	"Vaccine": {Price: 5.0, Value: 80, Description: "A vaccine to boost your health!"},
	"Pill":    {Price: 1.0, Value: 20, Description: "A small pill to help you recover."},
	"Vitamin": {Price: 0.5, Value: 15, Description: "Essential vitamins for daily health."},
	"Mineral": {Price: 0.1, Value: 10, Description: "Important minerals to keep you strong."},
}

type DrugStore struct {
	Drugs []types.EntityID `json:"drug_list"`
}

func (DrugStore) Name() string {
	return "DrugStore"
}

func (shop *DrugStore) InitDrugStore(world cardinal.WorldContext) {
	log := world.Logger()
	var entityId types.EntityID
	var err error
	for drugName, properties := range DrugKinds {

		// Since we're iterating directly, we can use the current drugName
		selectedDrug := drugName

		item := Item{
			ItemName:    selectedDrug,
			Kind:        ItemCare.String(),
			Description: properties.Description,
			Price:       properties.Price,
		}

		entityId, err = cardinal.Create(world, item,
			Health{HP: properties.Value},
		)

		if err != nil {
			log.Error().Msgf("Failed to create item %s: %v", item.ItemName, err)
			return
		}

		shop.Drugs = append(shop.Drugs, entityId)
	}

}
