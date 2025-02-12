package component

type Health struct {
	HP int `json:"health"`
}

func (Health) Name() string {
	return "Health"
}
