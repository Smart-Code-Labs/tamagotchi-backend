package msg

type ButItemMsg struct {
	Name string `json:"name"`
}

type BuyItemMsgReply struct {
	Success bool `json:"success"`
}
