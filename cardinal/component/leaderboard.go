package component

import (
	"sort"

	"pkg.world.dev/world-engine/cardinal"
)

const LeaderboardSize = 2

type Leaderboard struct {
	Pets []*Pet
}

func (Leaderboard) Name() string {
	return "Leaderboard"
}

// Add `pet` to `Leaderboard`
// The array `pets` is a sorted array with `pet.Level` field is the key, top to bottom
// There is a `unique` field `pet.PersonaTag` that cannot be duplicated on
// Flow:
// 1. check if pet `Level` field is greater than the botton element, if not greater then end.
// 2. Then check if this new `pet.PersonaTag` doesnt exist already. If exist then end.
// 3. Add pet to the list according to the sorted position based on `pet.Level`.
func (leaderboard *Leaderboard) AddPetToLeaderboard(world cardinal.WorldContext, addPet Pet) {
	log := world.Logger()
	n := len(leaderboard.Pets)
	log.Info().Msgf("Checking: n[%d]", n)
	// Early exit: leaderboard size
	if n < LeaderboardSize {
		log.Info().Msgf("Early Appending: [%s]", addPet.Nickname)
		leaderboard.Pets = append(leaderboard.Pets, &addPet)
		return
	}

	// check if already exist
	for _, pet := range leaderboard.Pets {
		log.Info().Msgf("Discarding: (already on leaderboard) [%s]", addPet.Nickname)
		if pet != nil && pet.Nickname == addPet.Nickname {
			return
		}
	}

	// lower level than the lowest
	if addPet.Level <= leaderboard.Pets[n-1].Level {
		log.Info().Msgf("Discarding: [%s] Lvl[%d] -> min Lvl[%d] ", addPet.Nickname, addPet.Level, leaderboard.Pets[n-1].Level)
		return
	}

	// Add the pet (unsorted)
	log.Info().Msgf("Appending: [%s] (unsorted)", addPet.Nickname)
	leaderboard.Pets = append(leaderboard.Pets, &addPet)

	log.Info().Msgf("sorting")
	// Sort the leaderboard (most efficient way to maintain sorted order after a single insert)
	sort.Slice(leaderboard.Pets, func(i, j int) bool {
		return leaderboard.Pets[i].Level > leaderboard.Pets[j].Level // Descending order (highest level first)
	})

	// Enforce size limit
	if len(leaderboard.Pets) > LeaderboardSize {
		leaderboard.Pets = leaderboard.Pets[:LeaderboardSize]
	}
}
