package component

import (
	"fmt"
	"math"

	constants "tamagotchi/game"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// nextLevelXP = baseXP * (level ^ growthRate)
const (
	baseXP     = int64(100)
	growthRate = 1.1
)

type Pet struct {
	PersonaTag  string `json:"personaTag"`
	Nickname    string `json:"nickname"`
	Gender      bool
	Level       int64 `json:"lvl"`
	XP          int64 `json:"exp"`
	NextLevelXP int64
}

func (Pet) Name() string {
	return "Pet"
}

func (h *Pet) AddXP(xp int64) {
	h.XP += xp
	for h.XP >= h.NextLevelXP {
		h.LevelUp()
	}
}

func (h *Pet) LevelUp() {
	h.Level += 1
	excessXP := h.XP - h.NextLevelXP
	h.XP = excessXP
	h.NextLevelXP = CalculateNextLevelXP(h.Level, baseXP, growthRate)
}

func CalculateNextLevelXP(level int64, baseXP int64, growthRate float64) int64 {
	nextLevelXP := int64(math.Pow(float64(level+1), growthRate)) * baseXP
	return nextLevelXP
}

// TODO check all `comp.Pet` if a pet with `Nickname` already exist. If exist return true, if not false

// CheckPetNicknameExists checks if a pet with the given nickname already exists in the game
// Returns true if the nickname exists, false otherwise
func CheckPetNicknameExists(world cardinal.WorldContext, nickname string) (bool, error) {
	var found bool = false
	var pet *Pet
	var err error
	log := world.Logger()
	log.Info().Msgf("called")
	// Find all entities with the Pet component
	q := cardinal.NewSearch().Entity(filter.Contains(filter.Component[Pet]()))
	q.Each(world, func(petId types.EntityID) bool {
		pet, err = cardinal.GetComponent[Pet](world, petId)
		if err != nil {
			return true
		}
		if pet.Nickname == nickname {
			found = true
		}
		return true
	})
	log.Info().Msgf("Checking if nickname already exist[%t]", found)
	if found {
		return false, fmt.Errorf("error creating pet: Name already exist")
	}

	return false, nil

}

func CreateRandomPet(world cardinal.WorldContext, personaTag string, nickname string) (types.EntityID, error) {
	rng := world.Rand()
	log := world.Logger()
	log.Info().Msgf("called")
	petID, err := cardinal.Create(world,
		Pet{PersonaTag: personaTag, Nickname: nickname, Level: 0, XP: 0, NextLevelXP: 0, Gender: rng.Intn(2) > 0},
		Health{HP: constants.InitialHP},
		Energy{E: constants.InitialE},
		Hygiene{Hy: constants.InitialHy},
		Wellness{Wn: constants.InitialWn},
		Dna{
			A: rng.Intn(100),
			C: rng.Intn(100),
			G: rng.Intn(100),
			T: rng.Intn(100),
		},
		Activity{Activity: "None", CountDown: 0},
		Think{Think: "..."},
	)
	if err != nil {
		log.Error().Msgf("Failed to create pet with nickname %s: %v", nickname, err)
		return 0, fmt.Errorf("error creating pet: %w", err)
	}
	log.Info().Msgf("Created: Pet[%d]", petID)

	return petID, nil
}
