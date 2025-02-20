package component

type Magic struct {
	Kind        string
	Level       int64 `json:"lvl"`
	XP          int64 `json:"exp"`
	NextLevelXP int64
}

func (Magic) Name() string {
	return "Magic"
}
