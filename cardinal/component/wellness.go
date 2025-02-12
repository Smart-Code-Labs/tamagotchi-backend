package component

type Wellness struct {
	Wn int `json:"wellness"`
}

func (Wellness) Name() string {
	return "Wellness"
}
