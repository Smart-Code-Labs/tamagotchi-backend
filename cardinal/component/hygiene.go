package component

type Hygiene struct {
	Hy int `json:"hygiene"`
}

func (Hygiene) Name() string {
	return "Hygiene"
}
