package msg

type FeedPetMsg struct {
	TargetNickname string `json:"target"`
	ItemId         string `json:"item_id"`
}

type FeedPetMsgReply struct {
	Health   int    `json:"health"`
	Activity string `json:"activity"`
	Duration int    `json:"duration"`
}
