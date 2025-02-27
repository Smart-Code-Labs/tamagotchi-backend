// Package component contains various components for the Tamagotchi game.
package component

// Skill represents a skill that a player or pet can have.
type Skill struct {
	// Kind is the type of skill, such as "fighting" or "magic".
	Kind string
	// Level is the current level of the skill.
	Level int64 `json:"lvl"`
	// XP is the amount of experience points the skill has.
	XP int64 `json:"exp"`
	// NextLevelXP is the amount of experience points needed to reach the next level.
	NextLevelXP int64
}

// Name returns the name of the component.
//
// Code Flow:
// 1. The function simply returns the string "Skill" as the name of the component.
func (Skill) Name() string {
	// The function returns the string "Skill" as the name of the component.
	return "Skill"
}
