package component

type Energy struct {
	E int `json:"energy"`
}

func (Energy) Name() string {
	return "Energy"
}
