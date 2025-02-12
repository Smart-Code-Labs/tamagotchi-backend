package msg

type CreatePetMsg struct {
	Nickname string `json:"nickname"`
}

type CreatePetResult struct {
	Success bool `json:"success"`
}
