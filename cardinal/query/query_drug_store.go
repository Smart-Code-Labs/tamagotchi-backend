// Package query contains functions for querying the state of the Tamagotchi game world.
package query

import (
	"tamagotchi/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// DrugsMsg is a request message for querying the available drugs.
type DrugsMsg struct{}

// DrugsReply is a response message containing a list of available drugs.
type DrugsReply struct {
	Drugs []component.Item `json:"drugs"`
}

/**
 * QueryDrugStore queries the available drug store from the DrugStore component.
 *
 * Flow:
 * 1. Log the receipt of the query request.
 * 2. Search for the unique DrugStore component in the game world.
 * 3. Process each entity that matches the search criteria.
 * 4. Retrieve the drug items from the DrugStore component.
 * 5. Return a response message containing the list of drug items.
 */
// QueryDrugStore queries the available QueryDrugStore from the DrugStore component
func QueryDrugStore(world cardinal.WorldContext, req *DrugsMsg) (*DrugsReply, error) {
	// Step 1: Log the receipt of the query request
	//         Log a message to indicate that the query request has been received.
	log := world.Logger()
	drugs := make([]component.Item, 0)
	log.Info().Msgf("Received payload to query-Drugs")

	// Step 2: Search for the unique DrugStore component
	//         Search for the DrugStore component using the Cardinal search API.
	q := cardinal.NewSearch().Entity(
		filter.Exact(filter.Component[component.DrugStore]()),
	)

	// Step 3: Process each entity that matches the search criteria
	//         Use the Cardinal Each API to process each entity that matches the search criteria.
	// Step 4: Retrieve the drug items from the DrugStore component
	//         Get the drug items from the DrugStore component and append them to the response list.
	searchError := q.Each(world, func(id types.EntityID) bool {
		// Step 4 (continued): Retrieve the DrugStore component
		//                     Get the DrugStore component for the current entity.
		drugStore, err := cardinal.GetComponent[component.DrugStore](world, id)
		if err != nil {
			return true
		}
		// Step 4 (continued): Retrieve the drug items
		//                     Get the drug items from the DrugStore component and append them to the response list.
		for _, entityId := range drugStore.Drugs {
			drugItem, err := cardinal.GetComponent[component.Item](world, entityId)
			if err != nil {
				return true
			}
			drugs = append(drugs, *drugItem)
		}
		return true
	})
	// Step 5: Return a response message containing the list of drug items
	//         Return a response message containing the list of drug items, or an error if the search failed.
	if searchError != nil {
		return nil, searchError
	}

	return &DrugsReply{Drugs: drugs}, nil
}
