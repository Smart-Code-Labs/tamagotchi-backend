package component

type Pet struct {
	Nickname string `json:"nickname"`
}

func (Pet) Name() string {
	return "Pet"
}
