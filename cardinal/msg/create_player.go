package msg

type CreatePlayerMsg struct{}

type CreatePlayerResult struct {
	Success bool `json:"success"`
}
