package component

type Skill struct {
	Kind        string
	Level       int64 `json:"lvl"`
	XP          int64 `json:"exp"`
	NextLevelXP int64
}

func (Skill) Name() string {
	return "Skill"
}
